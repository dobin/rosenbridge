package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/psanford/wormhole-william/wormhole"
	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/widgets"
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

/***/

func showError(s string) {
	a := widgets.NewQMessageBox(nil)
	a.SetText(s)
	a.Show()
	fmt.Printf("Error: %s", s)
}

func getWorkingDirectory() string {
	if len(settingsBridge.DownloadDirectory()) <= 0 {
		wd, err := os.Getwd()
		if err != nil {
			return ""
		}
		fmt.Printf("No download directory set, using working dir: %s\n", wd)
		return wd
	} else {
		wd := settingsBridge.DownloadDirectory()
		fmt.Printf("Using download directory: %s\n", wd)
		return wd
	}
}

/***/

func (b *ReceiveBridge) clickDownload(s string) { // Download
	jobtotal := new(int)
	jobdone := new(int)
	feedbackstr := new(string)
	var tableIndex int = 0

	// Check if code is valid
	msg, err := wormholeConnect(receiveBridge.Code())
	if err != nil {
		// log.Fatal(err)
		//fmt.Printf("Could not connect, wrong code?")
		//showError(err.Error())
		showError(fmt.Sprintf("Could not connect, wrong code %s?", receiveBridge.Code()))
		return
	}
	*jobtotal = msg.TransferBytes
	*jobdone = 0

	// Check if file already exists. Cancel if it does
	// This is necessary, as os.Rename() doesnt generate an error when file already
	// exists, even though it should...
	// Also better to handle it here before the download starts
	wd := getWorkingDirectory()
	filePath := fmt.Sprintf("%s/%s", wd, msg.Name)
	if _, err := os.Stat(filePath); err == nil {
		showError(fmt.Sprintf("File already exists: %s", filePath))
		return
	}

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

			if *feedbackstr == "Done" {
				receiverTableModel.editIdx(
					tableIndex,
					msg.Name,
					strconv.Itoa(*jobtotal),
					strconv.Itoa(*jobdone),
					"Done")
			} else {
				receiverTableModel.editIdx(
					tableIndex,
					msg.Name,
					strconv.Itoa(*jobtotal),
					strconv.Itoa(*jobdone),
					"Error")

				showError(*feedbackstr)
			}

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
		/// ???
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
	}()
}

func (b *ReceiveBridge) codeTextUpdate(s string) {
	receiveBridge.SetCode(s)
}
