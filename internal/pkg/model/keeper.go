package model

import (
	"time"
)

type Caver struct {
	First string
	Last  string
	Club  string
}

type Row struct {
	Date 	*time.Time // REVIEW: This might not be correct examine the fmt returned from db
	Cave 	string
	Names []*Caver // REVIEW: Could I make this a list of pointers ?
	Notes string
}
