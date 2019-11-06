package model

type Cave struct {
	ID      string
	Name    string
	Region  string
	Country string
	SRT     bool
	Visits  int64
}

type Caver struct {
	ID    string
	Name  string
	Club  string
	Count int64
}

/*type Entry struct {
	ID     int
	TripID int
	CaveID int
	CaverID int
}

type Trip struct {
	ID    int
	Date 	int
	Notes string
}*/

type Log struct {
	ID     string
	Date   string
	Cave   string
	Names  string // `, ` sep
	Notes  string
}
