package db

import (
	"errors"
	"strings"
)

func (db *Database) CheckCave(inputText string) ([]string, error) {
	var list []string

  result, err := db.conn.Prepare(`SELECT name FROM locations WHERE name LIKE ?`, inputText)
  if err != nil {
		//db.log.Errorf("db.checkcave: Failed to get caves", err)
		return list, err
  }
  defer result.Close()

  for {
    var name string

    exists, err := result.Step()
    if err != nil {
			//db.log.Errorf("db.checkcave: Step Error", err)
			return list, err
    }

    if !exists {
      break
    }

    err = result.Scan(&name)
    if err != nil {
			//db.log.Errorf("Scan: %v", err)
			return list, err
    }
    list = append(list, name)
  }

  return list, nil
}


//
// For retrieving the ID of a cave
func (db *Database) getCaveID(cave string) (int, error) {
	result, err := db.conn.Prepare(`SELECT id FROM locations WHERE name == ?`, cave)
	if err != nil {
		//db.log.Errorf("db.getcaverlist: Failed to get cavers", err)
		return 0, err
	}
	defer result.Close()

	var caveID int
	for {
		rowExists, err := result.Step()
		if err != nil {
			db.log.Errorf("db.get: Step error: %s", err)
			return 0, err
		}

		if !rowExists {
			break
		}

		err = result.Scan(&caveID)
		if err != nil {
			db.log.Errorf("Scan: %v", err)
			return caveID, err
		}
	}
	return caveID, err
}


//
// For formatting the ids for a new Log
func (db *Database) getCaverIDs(names string) ([]string, error) {
	var caverIDs []string

	cavers, err := db.GetAllCavers()
	if err != nil {
		//db.log.Errorf(`db.getcaverids: `)
		return caverIDs, err
	}

	namesList := strings.Split(names, ", ")

	for _, caver := range cavers {
		for _, name := range namesList {
			if strings.TrimSpace(name) == strings.TrimSpace(caver.Name) {
				caverIDs = append(caverIDs, caver.ID)
			}
		}
	}

	if len(caverIDs) != len(namesList) {
		return caverIDs, errors.New(`length mismatch`)
	}

	return caverIDs, nil
}
