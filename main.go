package main

import (
	"os"

	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/quick"
	"github.com/therecipe/qt/widgets"
)

var (
	receiveBridge  *ReceiveBridge
	senderBridge   *SenderBridge
	settingsBridge *SettingsBridge

	receiverTableModel *RecvFileTableModel
	senderTableModel   *SendFileTableModel
)

const BufferSize = 1024

func init() {
	RecvFileTableModel_QmlRegisterType2("CustomQmlTypes", 1, 0, "RecvFileTableModel")
	SendFileTableModel_QmlRegisterType2("CustomQmlTypes", 1, 0, "SendFileTableModel")
	ReceiveBridge_QmlRegisterType2("CustomQmlTypes", 1, 0, "ReceiveBridge")
	SenderBridge_QmlRegisterType2("CustomQmlTypes", 1, 0, "SenderBridge")
	SettingsBridge_QmlRegisterType2("CustomQmlTypes", 1, 0, "SettingsBridge")
}

func main() {
	core.QCoreApplication_SetAttribute(core.Qt__AA_EnableHighDpiScaling, true)

	receiverTableModel = NewRecvFileTableModel(nil)
	senderTableModel = NewSendFileTableModel(nil)

	app := widgets.NewQApplication(len(os.Args), os.Args)
	view := quick.NewQQuickView(nil)

	view.SetTitle("Rosen")
	view.SetResizeMode(quick.QQuickView__SizeRootObjectToView)
	view.SetSource(core.NewQUrl3("qrc:/qml/rosen.qml", 0))
	view.Show()

	app.Exec()
}
