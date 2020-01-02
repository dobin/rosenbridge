package main

import (
	"os"

	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/quick"
	"github.com/therecipe/qt/widgets"
)

var (
	receiveBridge *ReceiveBridge
)

func init() {
	FileTableModel_QmlRegisterType2("CustomQmlTypes", 1, 0, "FileTableModel")
	ReceiveBridge_QmlRegisterType2("CustomQmlTypes", 1, 0, "ReceiveBridge") // Download
}

func main() {
	core.QCoreApplication_SetAttribute(core.Qt__AA_EnableHighDpiScaling, true)

	app := widgets.NewQApplication(len(os.Args), os.Args)
	view := quick.NewQQuickView(nil)

	view.SetTitle("Rosen")
	view.SetResizeMode(quick.QQuickView__SizeRootObjectToView)
	view.SetSource(core.NewQUrl3("qrc:/qml/rosen.qml", 0))
	view.Show()

	app.Exec()
}