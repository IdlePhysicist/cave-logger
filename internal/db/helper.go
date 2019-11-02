package db


func (db *Database) CheckCave(inputText string) []string {
  result, err := db.conn.Prepare(`SELECT name FROM caves WHERE name LIKE ?`, inputText)
  if err != nil {
    db.log.Errorf("db.checkcave: Failed to get caves", err)
  }
  defer result.Close()

  var list []string
  for {
    var name string

    exists, err := result.Step()
    if err != nil {
      db.log.Errorf("db.checkcave: Step Error", err)
    }

    if !exists {
      break
    }

    err = result.Scan(&name)
    if err != nil {
			db.log.Errorf("Scan: %v", err)
    }
    list = append(list, name)
  }

  return list
}


//
// For retrieving the ID of a cave
func (db *Database) getCaveID(cave string) (int, error) {
	result, err := db.conn.Prepare(`SELECT id FROM caves WHERE name == ?`, cave)
	if err != nil {
		db.log.Errorf("db.getcaverlist: Failed to get cavers", err)
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
