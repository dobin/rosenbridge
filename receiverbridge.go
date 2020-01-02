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
	fmt.Printf("Download code: %s\n", receiveBridge.Code())

	// Check if code is valid
	msg, err := wormholeConnect(receiveBridge.Code())
	if err != nil {
		// log.Fatal(err)
		fmt.Printf("Could not connect, wrong code?")
		return
	}

	receiveTableModel.addNative(msg.Name, strconv.Itoa(msg.TransferBytes), "0", "Added")

	// Start downloading it in the background
	switch msg.Type {
	case wormhole.TransferText:
		wormholeTransferText(msg)
	case wormhole.TransferFile:
		wormholeTransferFile(msg)
	case wormhole.TransferDirectory:
		wormholeTransferDirectory(msg)
	}

	receiveTableModel.edit(
		msg.Name,
		strconv.Itoa(msg.TransferBytes),
		strconv.Itoa(msg.TransferBytes),
		"Done")
}

func (b *ReceiveBridge) codeTextUpdate(s string) {
	receiveBridge.SetCode(s)
}
