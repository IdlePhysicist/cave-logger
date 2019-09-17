package tui

type panel interface {
	name() string
	entries(*Tui)
	setEntries(*Tui)
	updateEntries(*Tui)
	setKeybinding(*Tui)
	focus(*Tui)
	unfocus()
	setFilterWord(string)
}