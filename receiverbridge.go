package main

import (
	"fmt"
	"strconv"

	"github.com/psanford/wormhole-william/wormhole"
	"github.com/therecipe/qt/core"
)

type ReceiveBridge struct {
	core.QObject

	_ func()         `constructor:"init"`
	_ func(s string) `signal:"clickDownload,auto"`
	_ func(s string) `signal:"codeTextUpdate,auto"`
	_ string         `property:"code"`

	_ *RecvFileTableModel `property:"TableModel"`
}

func (l *ReceiveBridge) init() {
	receiveBridge = l
	l.SetTableModel(receiverTableModel)
}

func (b *ReceiveBridge) clickDownload(s string) { // Download
	jobtotal := new(int)
	jobdone := new(int)
	feedbackstr := new(string)
	var tableIndex int = 0

	// Check if code is valid
	msg, err := wormholeConnect(receiveBridge.Code())
	if err != nil {
		// log.Fatal(err)
		fmt.Printf("Could not connect, wrong code?")
		return
	}
	*jobtotal = msg.TransferBytes
	*jobdone = 0

	tableIndex = receiverTableModel.addNative(
		msg.Name,
		strconv.Itoa(*jobtotal),
		"0",
		"Added")

	// Handling results of changes in
	// jobdone, feedbackstr
	t := core.NewQTimer(nil)
	t.ConnectEvent(func(e *core.QEvent) bool {
		// Handle errors and finish
		if len(*feedbackstr) > 0 {
			t.DisconnectEvent()
			receiverTableModel.editIdx(
				tableIndex,
				msg.Name,
				strconv.Itoa(*jobtotal),
				strconv.Itoa(*jobdone),
				"Done")

			/*a := widgets.NewQMessageBox(nil)
			a.SetText(*feedbackstr)
			a.Show()*/

			return true
		}

		// hand start and upload updates
		if *jobdone > 0 && *jobdone < *jobtotal {
			receiverTableModel.editIdx(
				tableIndex,
				msg.Name,
				strconv.Itoa(*jobtotal),
				strconv.Itoa(*jobdone),
				"Downloading")
		}

		return true

	})
	t.Start(200) // Every x ms

	// Thread doing the actual file receiving
	// It communicates with parent via variables
	// jobdone, feedbackstr
	go func() {
		/*
			defer func() {
				if err := recover(); err != nil {
					*errstr = fmt.Sprintf("%v", err)
				}
			}()
		*/

		// Start downloading it in the background
		switch msg.Type {
		case wormhole.TransferText:
			wormholeTransferText(msg, jobtotal, jobdone, feedbackstr)
		case wormhole.TransferFile:
			wormholeTransferFile(msg, jobtotal, jobdone, feedbackstr)
		case wormhole.TransferDirectory:
			wormholeTransferDirectory(msg, jobtotal, jobdone, feedbackstr)
		}

		//*feedbackstr = "ok"
	}()
}

func (b *ReceiveBridge) codeTextUpdate(s string) {
	receiveBridge.SetCode(s)
}
