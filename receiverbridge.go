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

	_ *FileTableModel `property:"TableModel"`
}

func (l *ReceiveBridge) init() {
	receiveBridge = l
	l.SetTableModel(receiveTableModel)
}

func (b *ReceiveBridge) clickDownload(s string) { // Download
	jobtotal := new(int)
	jobdone := new(int)
	feedbackstr := new(string)

	//fmt.Printf("Download code: %s\n", receiveBridge.Code())

	// Check if code is valid
	msg, err := wormholeConnect(receiveBridge.Code())
	if err != nil {
		// log.Fatal(err)
		fmt.Printf("Could not connect, wrong code?")
		return
	}

	*jobtotal = msg.TransferBytes
	*jobdone = 0

	// Add file to table
	receiveTableModel.addNative(
		msg.Name,
		strconv.Itoa(*jobtotal),
		"0",
		"Added")

	t := core.NewQTimer(nil)
	t.ConnectEvent(func(e *core.QEvent) bool {
		receiveTableModel.edit(
			msg.Name,
			strconv.Itoa(*jobtotal),
			strconv.Itoa(*jobdone),
			"Downloading")

		if len(*feedbackstr) > 0 {
			t.DisconnectEvent()
			receiveTableModel.edit(
				msg.Name,
				strconv.Itoa(*jobtotal),
				strconv.Itoa(*jobdone),
				"Done")

			/*a := widgets.NewQMessageBox(nil)
			a.SetText(*feedbackstr)
			a.Show()*/

			return true
		}
		return true

	})
	t.Start(100)

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
