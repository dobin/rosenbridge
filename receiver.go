package main

import (
	"fmt"

	"github.com/therecipe/qt/core"
)

type ReceiveBridge struct {
	core.QObject

	_ func()         `constructor:"init"`
	_ func(s string) `signal:"clickDownload,auto"`

	_ func(s string) `signal:"codeTextUpdate,auto"`
	_ string         `property:"code"`
}

func (l *ReceiveBridge) init() {
	receiveBridge = l
}

func (b *ReceiveBridge) clickDownload(s string) { // Download
	fmt.Printf("Click click 1: %s\n", receiveBridge.Code())
	//fmt.Printf("Click click 2: %s\n", s)
	tableModel.addNative("a", "b", "c", "d")
}

func (b *ReceiveBridge) codeTextUpdate(s string) {
	receiveBridge.SetCode(s)
}
