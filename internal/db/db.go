package db

import (
	"context"
	"strings"
	"time"

	"github.com/bvinc/go-sqlite-lite/sqlite3"
	"github.com/sirupsen/logrus"

	"github.com/idlephysicist/cave-logger/internal/model"
)

const date_layout = `2006-01-02`

type Database struct {
	log	 *logrus.Logger
	conn *sqlite3.Conn
	ctx	 context.Context
}

func New(log *logrus.Logger, dbFN string) *Database {
	var db Database

	db.log = log
	db.ctx = context.Background()

	hldr, err := sqlite3.Open(dbFN)
	if err != nil {
		log.Fatalf("database.new: Cannot establish database connection", err)
	}
	db.conn = hldr
	return &db
}

func (db *Database) Close() {
	db.conn.Close()
}

//
// MAIN FUNCTIONS --------------------------------------------------------------
//


func (db *Database) AddLog(date, cave, names, notes string) error {
	query := `INSERT INTO entries (date, caveid, caverids, notes) VALUES (?,?,?,?)`

	d, err := time.Parse(date_layout, date)
	if err != nil {
		return err
	}
	dateStamp := d.Unix()

	caveID, err := db.getCaveID(cave)
	if err != nil {
		return err
	} else if caveID == 0 {
		// return err prompting dialog.
	}

	caverIDs := db.getCaverIDs(names)

	params := []interface{}{dateStamp, caveID, caverIDs, notes}

	_, err = db.insert(query, params)
	if err != nil {
		return err
	}
	return nil
}

func (db *Database) AddCave(name, region, country string, srt bool) (int64, error) {
	query := `INSERT INTO caves (name, region, country, srt) VALUES (?,?,?,?)`
	params := []interface{}{name, region, country, srt}

	newID, err := db.insert(query, params)
	if err != nil {
		return -1, err
	}
	return newID, nil
}

func (db *Database) AddCaver(name, club string) (int64, error) {
	query := `INSERT INTO cavers (name, club) VALUES (?,?)`
	params := []interface{}{name, club}

	newID, err := db.insert(query, params)
	if err != nil {
		return -1, err
	}
	return newID, nil
}

func (db *Database) GetLogs(logID string) ([]*model.Entry, error) {
	// Build query
	var query string // Go seems to complain if something is defined in if blocks
	if logID == `-1` {
		query = `SELECT
			entries.id AS 'id',
			date AS 'date',
			name AS 'cave',
			caverids AS 'caverids',
			notes AS 'notes'
		FROM entries JOIN caves 
		WHERE entries.caveid == caves.id`
		
		//result, err := db.conn.Prepare(query)
	} else {
		query = `SELECT 
			entries.id AS 'id',
			date AS 'date',
			name AS 'cave',
			caverids AS 'caverids',
			notes AS 'notes'
		FROM entries JOIN caves 
		WHERE entries.caveid == caves.id AND entries.id = ?`

	}

	result, err := db.conn.Prepare(query)
	if err != nil {
		db.log.Errorf("db.prepare: Failed to query database", err)
		return nil, err
	}
	defer result.Close()
	
	trips := make([]*model.Entry, 0)
	for {
		//var caverIDs []int
		var caverIDstr string
		var stamp int64
		var trip model.Entry

		rowExists, err := result.Step()
		if err != nil {
			db.log.Errorf("db.get: Step error: %s", err)
			return trips, err
		}

		if !rowExists {
			break
		}
		
		err = result.Scan(&trip.ID, &stamp, &trip.Cave, &caverIDstr, &trip.Notes)
		if err != nil {
			db.log.Error(err)
			return trips, err
		}

		trip.Date = time.Unix(stamp, 0).Format(date_layout)

		// Break up cavers into ids.
		/*brokenIDs := strings.Split(caverIDstr, `|`)
		for _, id_str := range brokenIDs {
			id, _ := strconv.Atoi(id_str)
			caverIDs = append(caverIDs, id)
		}*/
		
		trip.Names = db.getCaverNames(caverIDstr)
		//if row.Names == `` {
		//	continue
		//}
		
		// Add this formatted row to the rows map
		trips = append(trips, &trip)  
	}

	return trips, err
}

func (db *Database) GetAllCavers() ([]*model.Caver, error) {
	result, err := db.conn.Prepare("SELECT `id`,`name`,`club` FROM cavers")
	if err != nil {
		db.log.Errorf("db.getcaverlist: Failed to get cavers", err)
	}

	cavers := make([]*model.Caver, 0)
	for {
		var id int
		var c model.Caver
		
		rowExists, err := result.Step()
		if err != nil {
			db.log.Errorf("db.get: Step error: %s", err)
			return cavers, err
		}

		if !rowExists {
			break
		}
		
		err = result.Scan(&id, &c.Name, &c.Club)
		if err != nil {
			db.log.Errorf("Scan: %v", err)
		}
		cavers = append(cavers, &c)
		//cavers[id] = c
	}
	return cavers, err
}

func (db *Database) RemoveLog(logID string) error {
	return nil
}

//
// INTERNAL FUNCTIONS ----------------------------------------------------------
//


func (db *Database) insert(query string, params []interface{}) (int64, error) {
	err := db.conn.Exec(query, params...)
	if err != nil {
		return -1, err
	}

	return db.conn.LastInsertRowID(), nil
}

//
// For formatting the ids for a new entry
func (db *Database) getCaverIDs(names string) string {
	var caverIDs []string

	cavers, err := db.GetAllCavers()
	if err != nil {
		db.log.Errorf(`db.getcaverids: `)
		return ``
	}

	for _, caver := range cavers {
		namesList := strings.Split(names, ", ")

		for _, fullName := range namesList {
			if fullName == caver.Name {
				caverIDs = append(caverIDs, caver.ID)
			}
		}
	}

	return strings.Join(caverIDs, `|`)
}

//
// For retrieving the names given a str of ids 
func (db *Database) getCaverNames(idStr string) string {
	// Get the IDs
	cavers, err := db.GetAllCavers()
	if err != nil {
		db.log.Errorf("Database.Query: Failed to fetch list of cavers")
		return ``
	}

	// Split up the str arg we were given
	var names []string
	caverIDs := strings.Split(idStr, "|")

	for _, caver_id := range caverIDs {
		//caver_id, _ := strconv.Atoi(_id_str) // Convert the string to an int...

		for _, caver := range cavers {
			if caver_id == caver.ID {
				names = append(names, caver.Name)
			} else {
				continue
			}
		}	
	}
	
	return strings.Join(names, `, `)
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


/*
func keyExists(array map[int]*model.Caver, key int) bool {
	_, grand := array[key]
	return grand
}
*/