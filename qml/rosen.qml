import QtQuick 2.9
import QtQuick.Window 2.3
import QtQuick.Controls 1.4
import QtQuick.Layouts 1.3
import CustomQmlTypes 1.0

Item {
    id: root

    property ReceiveBridge template: ReceiveBridge{}

    width: 640
    height: 480

    ColumnLayout {
        anchors.fill: parent

        RowLayout {
            id: rowLayout
            width: 100
            height: 64
            Layout.alignment: Qt.AlignLeft | Qt.AlignTop
            Layout.fillWidth: true

            Button {
                id: button
                text: qsTr("Send")
                Layout.fillHeight: true
                Layout.fillWidth: true
            }

            Button {
                id: button1
                text: qsTr("Receive")
                Layout.fillHeight: true
                Layout.fillWidth: true
            }

        }


        ColumnLayout {
            Layout.alignment: Qt.AlignLeft | Qt.AlignTop
            Layout.fillHeight: true
            Layout.fillWidth: true

            RowLayout {
                id: rowLayout1
                width: 100
                height: 100

                Label {
                    Layout.fillHeight: true
                    Layout.fillWidth: false
                    Layout.alignment: Qt.AlignLeft | Qt.AlignTop

                    text: qsTr("Code")
                }

                TextField {
                    id: ti
                    width: 80
                    height: 20
                    font.pixelSize: 12
                    Layout.fillWidth: true
                    Layout.fillHeight: true

                    text: root.template.code
                    onTextChanged: root.template.codeTextUpdate(text)
                }

                Button {
                    id: buttonDownload
                    text: qsTr("Download")
                    Layout.fillHeight: true

                    onClicked: {
                        root.template.clickDownload(ti.text)
                    }
                }
            }

            TableView {
                id: tableview

                Layout.fillWidth: true
                Layout.fillHeight: true

                model: FileTableModel{}

                TableViewColumn {
                    role: "Filename"
                    title: role
                }

                TableViewColumn {
                    role: "Filesize"
                    title: role
                }

                TableViewColumn {
                    role: "Transmitted"
                    title: role
                }

                TableViewColumn {
                    role: "Status"
                    title: role
                }
            }
        }
    }
}
