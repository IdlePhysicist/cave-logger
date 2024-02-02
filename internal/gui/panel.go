package gui

type panel interface {
	name() string
	entries(*Gui)
	setEntries(*Gui)
	updateEntries(*Gui)
	focus(*Gui)
	unfocus()
	setFilter(string, string, string)
}
