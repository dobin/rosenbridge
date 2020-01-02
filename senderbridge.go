package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/therecipe/qt/core"
)

type SenderBridge struct {
	core.QObject

	_ func() `constructor:"init"`
	_ func() `signal:"clickAddFile,auto"`
	_ string `property:"code"`

	_ *FileTableModel `property:"TableModel"`
}

func (l *SenderBridge) init() {
	senderBridge = l
	l.SetTableModel(senderTableModel)
}

func (b *SenderBridge) clickAddFile() {
	fmt.Printf("Sending\n")
	filename := "meh"

	// Check which type it is
	stat, err := os.Stat(filename)
	if err != nil {
		bail("Failed to read %s: %s", filename, err)
	}

	senderTableModel.addNative(
		filename,
		strconv.FormatInt(stat.Size(), 10),
		"0",
		"Added")

	if stat.IsDir() {
		sendDir(filename)
	} else {
		sendFile(filename)
	}

	senderTableModel.edit(
		filename,
		strconv.FormatInt(stat.Size(), 10),
		strconv.FormatInt(stat.Size(), 10),
		"Done")
}
