package main

import (
	"fmt"

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
	fmt.Printf("Click click 2\n")
	senderTableModel.addNative("a", "b", "c", "d")
}
