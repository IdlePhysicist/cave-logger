package model

type Cave struct {
	ID      int
	Name    string
	Region  string
	Country string
	SRT     bool
	Visits  int
}

type Caver struct {
	ID    int
	First string
	Last  string
	Club  string
}

type Entry struct {
	ID       int
	Date 	   string // REVIEW: This might not be correct examine the fmt returned from db
	Cave 	   string
	Names    string // REVIEW: Could I make this a list of pointers ?
	CaverIDs []int
	Notes    string
}
