package main

import (
	"fmt"

	"github.com/therecipe/qt/core"
)

type SettingsBridge struct {
	core.QObject

	_ func()         `constructor:"init"`
	_ func(s string) `signal:"onServerAddressUpdate,auto"`
	_ func(s string) `signal:"onDownloadDirectoryUpdate,auto"`
	_ func(s string) `signal:"onServerAddressEnabledUpdate,auto"`

	_ string `property:"ServerAddress"`
	_ string `property:"ServerAddressEnabled"`
	_ string `property:"DownloadDirectory"`
}

func (b *SettingsBridge) init() {
	fmt.Printf("Settings init")

	settingsBridge = b
}

func (b *SettingsBridge) onServerAddressUpdate(s string) {
	//fmt.Printf("Server Address: %s", s)
	b.SetServerAddress(s)
}

func (b *SettingsBridge) onServerAddressEnabledUpdate() {
	fmt.Printf("Server Address enabled: ??")
	//b.SetServerAddressEnabled(s)
}

func (b *SettingsBridge) onDownloadDirectoryUpdate(s string) {
	//fmt.Printf("Download directory: %s", s)
	b.SetDownloadDirectory(s)
}
