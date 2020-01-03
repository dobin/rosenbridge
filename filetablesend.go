package main

import "github.com/therecipe/qt/core"

const (
	sendFilename = int(core.Qt__UserRole) + 1<<iota
	sendCode
	sendFilesize
	sendTransmitted
	sendStatus
)

type SendFileTableItem struct {
	filename    string
	code        string
	size        string
	transmitted string
	status      string
}

type SendFileTableModel struct {
	core.QAbstractTableModel

	_ func() `constructor:"init"`

	_ func()                                                                                           `signal:"remove,auto"`
	_ func(item []*core.QVariant)                                                                      `signal:"add,auto"`
	_ func(filename string, code string, size string, transmitted string, status string)               `signal:"editLast,auto"`
	_ func(tableIdx int, filename string, code string, size string, transmitted string, status string) `signal:"editIdx,auto"`

	modelData []SendFileTableItem
}

func (m *SendFileTableModel) init() {
	m.modelData = []SendFileTableItem{
		//		{"test1.txt", "1000", "1000", "Done"},
		//		{"test2.txt", "1000", "0", "Started"},
	}

	m.ConnectRoleNames(m.roleNames)
	m.ConnectRowCount(m.rowCount)
	m.ConnectColumnCount(m.columnCount)
	m.ConnectData(m.data)
}

func (m *SendFileTableModel) roleNames() map[int]*core.QByteArray {
	return map[int]*core.QByteArray{
		sendFilename:    core.NewQByteArray2("sendFilename", -1),
		sendCode:        core.NewQByteArray2("sendCode", -1),
		sendFilesize:    core.NewQByteArray2("sendFilesize", -1),
		sendTransmitted: core.NewQByteArray2("sendTransmitted", -1),
		sendStatus:      core.NewQByteArray2("sendStatus", -1),
	}
}

func (m *SendFileTableModel) rowCount(*core.QModelIndex) int {
	return len(m.modelData)
}

func (m *SendFileTableModel) columnCount(*core.QModelIndex) int {
	return 5
}

func (m *SendFileTableModel) data(index *core.QModelIndex, role int) *core.QVariant {
	item := m.modelData[index.Row()]
	switch role {
	case sendFilename:
		return core.NewQVariant1(item.filename)
	case sendCode:
		return core.NewQVariant1(item.code)
	case sendFilesize:
		return core.NewQVariant1(item.size)
	case sendTransmitted:
		return core.NewQVariant1(item.transmitted)
	case sendStatus:
		return core.NewQVariant1(item.status)

	}
	return core.NewQVariant()
}

func (m *SendFileTableModel) remove() {
	if len(m.modelData) == 0 {
		return
	}
	m.BeginRemoveRows(core.NewQModelIndex(), len(m.modelData)-1, len(m.modelData)-1)
	m.modelData = m.modelData[:len(m.modelData)-1]
	m.EndRemoveRows()
}

func (m *SendFileTableModel) add(item []*core.QVariant) {
	m.BeginInsertRows(core.NewQModelIndex(), len(m.modelData), len(m.modelData))
	m.modelData = append(
		m.modelData,
		SendFileTableItem{
			item[0].ToString(),
			item[1].ToString(),
			item[2].ToString(),
			item[3].ToString(),
			item[4].ToString(),
		})
	m.EndInsertRows()
}

func (m *SendFileTableModel) addNative(
	filename string, code string, size string, transmitted string, status string) int {
	m.BeginInsertRows(core.NewQModelIndex(), len(m.modelData), len(m.modelData))
	m.modelData = append(
		m.modelData,
		SendFileTableItem{
			filename,
			code,
			size,
			transmitted,
			status,
		})
	m.EndInsertRows()

	return len(m.modelData) - 1
}

func (m *SendFileTableModel) editLast(
	filename string,
	code string,
	size string,
	transmitted string,
	status string) {
	if len(m.modelData) == 0 {
		return
	}
	m.modelData[len(m.modelData)-1] = SendFileTableItem{
		filename, code, size, transmitted, status}
	m.DataChanged(
		m.Index(len(m.modelData)-1, 0, core.NewQModelIndex()),
		m.Index(len(m.modelData)-1, 1, core.NewQModelIndex()),
		[]int{sendFilename, sendCode, sendFilesize, sendTransmitted, sendStatus})
}

func (m *SendFileTableModel) editIdx(
	tableIdx int,
	filename string,
	code string,
	size string,
	transmitted string,
	status string) {
	if len(m.modelData) == 0 {
		return
	}
	m.modelData[tableIdx] = SendFileTableItem{
		filename, code, size, transmitted, status}
	m.DataChanged(
		m.Index(tableIdx, 0, core.NewQModelIndex()),
		m.Index(tableIdx, 1, core.NewQModelIndex()),
		[]int{sendFilename, sendCode, sendFilesize, sendTransmitted, sendStatus})
}
