package model

type Cave struct {
	ID      string
	Name    string
	Region  string
	Country string
	SRT     bool
	Visits  int64
	Notes   string
}

type Caver struct {
	ID    string
	Name  string
	Club  string
	Count int64
	Notes string
}

type Log struct {
	ID     string
	Date   string
	Cave   string
	Names  string // `, ` sep
	Notes  string
}
