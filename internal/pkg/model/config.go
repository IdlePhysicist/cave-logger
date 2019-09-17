package model

type Config struct {
	Database struct {
		Created string `json:"created"`
		Filename string `json:"filename"`
	} `json:"database"`
}

