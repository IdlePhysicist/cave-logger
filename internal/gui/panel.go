package gui

type panel interface {
	name() string
	entries(*Gui)
	setEntries(*Gui)
	updateEntries(*Gui)
	setKeybinding(*Gui)
	focus(*Gui)
	unfocus()
	setFilter(string, string)
	getSortedCol() int
	setSortedCol(int)
	getColumnCount() int
}
