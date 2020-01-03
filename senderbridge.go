package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/widgets"
)

type SenderBridge struct {
	core.QObject

	_ func() `constructor:"init"`
	_ func() `signal:"clickAddFile,auto"`
	_ string `property:"code"`

	_ *SendFileTableModel `property:"TableModel"`
}

func (l *SenderBridge) init() {
	senderBridge = l
	l.SetTableModel(senderTableModel)
}

func (b *SenderBridge) clickAddFile() {
	jobtotal := new(int64)
	jobdone := new(int64)
	feedbackstr := new(string)
	mycode := new(string) // The real code
	code := new(string)   // Transport the code from thread to copy into mycode
	var tableIndex int = 0

	filenames := widgets.QFileDialog_GetOpenFileNames(nil, "some caption", "", "", "", 0)
	if len(filenames) < 1 {
		return
	}
	filename := filenames[0]
	fmt.Printf("Sending: %s\n", filename)

	// Check which type it is
	stat, err := os.Stat(filename)
	if err != nil {
		bail("Failed to read %s: %s", filename, err)
	}

	*jobtotal = stat.Size()
	*jobdone = 0

	// Add file to table
	tableIndex = senderTableModel.addNative(
		filename,
		*mycode,
		strconv.FormatInt(*jobtotal, 10),
		"0",
		"Added")

	t := core.NewQTimer(nil)
	t.ConnectEvent(func(e *core.QEvent) bool {
		senderTableModel.editIdx(
			tableIndex,
			filename,
			*mycode,
			strconv.FormatInt(*jobtotal, 10),
			strconv.FormatInt(*jobdone, 10),
			"Uploading")

		if len(*feedbackstr) > 0 {
			t.DisconnectEvent()
			senderTableModel.editIdx(
				tableIndex,
				filename,
				*mycode,
				strconv.FormatInt(*jobtotal, 10),
				strconv.FormatInt(*jobdone, 10),
				"Done")

			/*a := widgets.NewQMessageBox(nil)
			a.SetText(*feedbackstr)
			a.Show()*/

			return true
		}

		if len(*code) > 0 {
			mycode = code
			senderTableModel.editIdx(
				tableIndex,
				filename,
				*mycode,
				strconv.FormatInt(*jobtotal, 10),
				"0",
				"Added")
			*code = ""
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
		if stat.IsDir() {
			sendDir(filename)
		} else {
			code, status, err := sendFile(filename, jobdone, feedbackstr)

			fmt.Printf("Code 1: %s", code)
			if err != nil {
				bail("Error sending message: %s", err)
			}

			// Cant update TableModel here in this go func()
			// Return it to the QTimer
			*mycode = code

			// Wait till its finished
			s := <-status

			if s.OK {
				fmt.Println("file sent")
				*feedbackstr = "Ok"
			} else {
				bail("Send error: %s", s.Error)
				*feedbackstr = "Error"
			}
		}

		//*feedbackstr = "ok"
	}()

	/*senderTableModel.edit(
	filename,
	strconv.FormatInt(stat.Size(), 10),
	strconv.FormatInt(stat.Size(), 10),
	"Done")*/
}
