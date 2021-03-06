import QtQuick 2.4
import QtQuick 2.9
import QtQuick.Window 2.3
import QtQuick.Controls 1.4
import QtQuick.Layouts 1.3
import QtQuick.Window 2.1
import QtQuick.Dialogs 1.1

Window {
    width: 800
    height: 400
    id: settingsview

    FileDialog {
            id: fileDialog
            property string dir

            title: "Choose download directory"
            selectExisting: true
            selectMultiple: false
            selectFolder: true
            selectedNameFilter: "All directories (*)"
            sidebarVisible: false

            onAccepted: {
                dir = fileDialog.fileUrls[0].slice(7)
                console.log("Accepted: " + dir)
                textFieldDownloadDirectory.text = dir
                root.settingsbridge.onDownloadDirectoryUpdate(dir)
            }
            onRejected: { console.log("Rejected") }
        }


    ColumnLayout {
        id: columnLayout
        anchors.fill: parent

        RowLayout {
            id: rowLayout
            height: 32

            Label {
                id: label
                text: qsTr("Download Directory:")
                verticalAlignment: Text.AlignVCenter
                Layout.alignment: Qt.AlignLeft | Qt.AlignTop
                Layout.fillHeight: true
            }

            TextField {
                id: textFieldDownloadDirectory
                text: root.settingsbridge.downloadDirectory
                Layout.fillWidth: true
                placeholderText: qsTr("Download Directory Path (empty = working directory")
                onTextChanged: root.settingsbridge.onDownloadDirectoryUpdate(text)
            }

            Button {
                text: "Open"
                onClicked: fileDialog.open()
            }
        }

        RowLayout {
            id: rowLayout1
            width: 100
            height: 32

            Label {
                id: label1
                text: qsTr("Server:")
                verticalAlignment: Text.AlignVCenter
                Layout.fillHeight: true
            }

            TextField {
                id: textField
                Layout.fillHeight: true
                Layout.fillWidth: true
                placeholderText: qsTr("Server Address (empty = default)")
                onTextChanged: root.settingsbridge.onServerAddressUpdate(text)
            }
        }

        RowLayout {
            id: rowLayout2
            width: 100
            height: 100

            TextArea {
                id: textArea
                Layout.fillHeight: true
                Layout.fillWidth: true
            }
        }

        RowLayout {
            id: rowLayout3
            width: 100
            height: 32
            Layout.fillWidth: true

            Button {
                id: buttonSave
                text: qsTr("Save")
                Layout.fillHeight: true
                Layout.fillWidth: true
            }

            Button {
                id: buttonClose
                text: qsTr("Close")
                Layout.fillHeight: true
                Layout.fillWidth: true
            }
        }
    }
}
