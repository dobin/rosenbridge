import QtQuick 2.9
import QtQuick.Window 2.3
import QtQuick.Controls 1.4
import QtQuick.Layouts 1.3
import CustomQmlTypes 1.0

Item {
    id: root
    property ReceiveBridge receivebridge: ReceiveBridge{}
    property SenderBridge senderbridge: SenderBridge{}

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
                id: buttonReceive
                text: qsTr("Receive")
                Layout.fillHeight: true
                Layout.fillWidth: true
                enabled: false

                onClicked: {
                    buttonSend.enabled = true
                    buttonReceive.enabled = false
                    tabView.currentIndex = 0
                }
            }

            Button {
                id: buttonSend
                text: qsTr("Send")
                Layout.fillHeight: true
                Layout.fillWidth: true

                onClicked: {
                    buttonSend.enabled = false
                    buttonReceive.enabled = true
                    tabView.currentIndex = 1
                }
            }
        }


        TabView {
            id: tabView
            width: 360
            height: 300
            Layout.fillHeight: true
            Layout.fillWidth: true
            tabsVisible: false

            Tab {
                id: viewReceiver
                anchors.fill: parent

                ColumnLayout {
                    anchors.fill: parent
                    Layout.alignment: Qt.AlignLeft | Qt.AlignTop
                    Layout.fillHeight: true
                    Layout.fillWidth: true

                    RowLayout {
                        id: rowLayout1
                        width: 100
                        height: 100

                        TextField {
                            id: ti
                            width: 80
                            height: 20
                            font.pixelSize: 12
                            Layout.fillWidth: true
                            Layout.fillHeight: true

                            text: root.receivebridge.code
                            readOnly: false
                            placeholderText: "Wormhole code"
                            onTextChanged: root.receivebridge.codeTextUpdate(text)
                        }

                        Button {
                            id: buttonDownload
                            text: qsTr("Download")
                            Layout.fillHeight: true

                            onClicked: {
                                root.receivebridge.clickDownload(ti.text)
                            }
                        }
                    }

                    TableView {
                        id: tableview

                        Layout.fillWidth: true
                        Layout.fillHeight: true

                        model: root.receivebridge.TableModel

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

            Tab {
               id: viewSend
                anchors.fill: parent
                
                ColumnLayout {
                    id: columnLayout
                    anchors.fill: parent
                    
                    RowLayout {
                        id: rowLayout2
                        y: 0
                        width: 100
                        height: 64
                        Layout.fillHeight: false
                        Layout.fillWidth: true
                        Layout.alignment: Qt.AlignLeft | Qt.AlignTop
                        
                        Button {
                            id: buttonSendAdd
                            text: qsTr("Add")
                            Layout.fillHeight: true
                            Layout.fillWidth: true

                            onClicked: {
                                root.senderbridge.clickAddFile()
                            }
                        }
                        
                        Button {
                            id: buttonSendRemove
                            text: qsTr("Button")
                            Layout.fillHeight: true
                            Layout.fillWidth: true
                        }
                    }
                    
                    TableView {
                        id: tableviewSend
                        
                        Layout.fillWidth: true
                        Layout.fillHeight: true
                        
                        model: root.senderbridge.TableModel
                        
                        TableViewColumn {
                            role: "sendFilename"
                            title: role
                        }

                        TableViewColumn {
                            role: "sendCode"
                            title: role
                        }

                        TableViewColumn {
                            role: "sendFilesize"
                            title: role
                        }
                        
                        TableViewColumn {
                            role: "sendTransmitted"
                            title: role
                        }
                        
                        TableViewColumn {
                            role: "sendStatus"
                            title: role
                        }
                    }
                    
                    
                }
            }
        }
    }
}
