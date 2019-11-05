package db

import (
	"context"
	"strings"
	"time"

	"github.com/bvinc/go-sqlite-lite/sqlite3"
	"github.com/sirupsen/logrus"

	"github.com/idlephysicist/cave-logger/internal/model"
)

const datetime = `2006-01-02T15:04:05Z`
const date = `2006-01-02`

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

	d, err := time.Parse(datetime, strings.Join([]string{date,`12:00:00Z`},`T`))
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
	query := `INSERT INTO cavers (first, last, club) VALUES (?,?)`
	params := []interface{}{name, club}

	newID, err := db.insert(query, params)
	if err != nil {
		return -1, err
	}
	return newID, nil
}

func (db *Database) GetAllLogs() ([]*model.Log, error) {
	// Build query
	var query string
	query = `SELECT
			entries.id AS 'id',
			date AS 'date',
			name AS 'cave',
			caverids AS 'caverids',
			notes AS 'notes'
		FROM entries JOIN caves 
		WHERE entries.caveid == caves.id`

	result, err := db.conn.Prepare(query)
	if err != nil {
		db.log.Errorf("db.prepare: Failed to query database", err)
		return nil, err
	}
	defer result.Close()
	
	trips := make([]*model.Log, 0)
	for {
		var caverIDstr string
		var stamp int64
		var trip model.Log

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

		trip.Date = time.Unix(stamp, 0).Format(date)
		trip.Names = db.getCaverFirstNames(caverIDstr)
		
		// Add this formatted row to the rows map
		trips = append(trips, &trip)  
	}

	return trips, err
}

func (db *Database) GetLog(logID string) (*model.Log, error) { //FIXME: 
	// Build query
	var query string
	query = `SELECT 
			entries.id AS 'id',
			date AS 'date',
			name AS 'cave',
			caverids AS 'caverids',
			notes AS 'notes'
		FROM entries JOIN caves 
		WHERE entries.caveid == caves.id AND entries.id = ?`

	result, err := db.conn.Prepare(query, logID)
	if err != nil {
		db.log.Errorf("db.prepare: Failed to query database", err)
		return nil, err
	}
	defer result.Close()
	
	trips := make([]*model.Log, 0)
	for {
		var caverIDstr string
		var stamp int64
		var trip model.Log

		rowExists, err := result.Step()
		if err != nil {
			db.log.Errorf("db.get: Step error: %s", err)
			return trips[0], err
		}

		if !rowExists {
			break
		}
		
		err = result.Scan(&trip.ID, &stamp, &trip.Cave, &caverIDstr, &trip.Notes)
		if err != nil {
			db.log.Error(err)
			return trips[0], err
		}

		trip.Date = time.Unix(stamp, 0).Format(date)
		trip.Names = db.getFullCaverNames(caverIDstr)
		
		// Add this formatted row to the rows map
		trips = append(trips, &trip)  
	}

	return trips[0], err
}

func (db *Database) GetAllCavers() ([]*model.Caver, error) {
	result, err := db.conn.Prepare("SELECT `id`,`first`,`last`,`club` FROM cavers")
	if err != nil {
		db.log.Errorf("db.getcaverlist: Failed to get cavers", err)
	}

	cavers := make([]*model.Caver, 0)
	for {
		var c model.Caver
		
		rowExists, err := result.Step()
		if err != nil {
			db.log.Errorf("db.get: Step error: %s", err)
			return cavers, err
		}

		if !rowExists {
			break
		}
		
		err = result.Scan(&c.ID, &c.First, &c.Last, &c.Club)
		if err != nil {
			db.log.Errorf("Scan: %v", err)
		}
		cavers = append(cavers, &c)
		//cavers[id] = c
	}
	return cavers, err
}

func (db *Database) GetAllCaves() ([]*model.Cave, error) {
	result, err := db.conn.Prepare(
		"SELECT `id`,`name`,`region`,`country`, `srt`, `visits` FROM caves",
	)
	if err != nil {
		db.log.Errorf("db.getcaverlist: Failed to get cavers", err)
	}

	caves := make([]*model.Cave, 0)
	for {
		var c model.Cave
		
		rowExists, err := result.Step()
		if err != nil {
			db.log.Errorf("db.get: Step error: %s", err)
			return caves, err
		}

		if !rowExists {
			break
		}
		
		err = result.Scan(&c.ID, &c.Name, &c.Region, &c.Country, &c.SRT, &c.Visits)
		if err != nil {
			db.log.Errorf("Scan: %v", err)
		}
		caves = append(caves, &c)
		//cavers[id] = c
	}
	return caves, err
}

func (db *Database) GetCave(caveID string) (*model.Cave, error) {
	// Build query
	var query string
	query = "SELECT `id`,`name`,`region`,`country`,`srt`,`visits` FROM caves WHERE id = ?"

	result, err := db.conn.Prepare(query, caveID)
	if err != nil {
		db.log.Errorf("db.prepare: Failed to query database", err)
		return nil, err
	}
	defer result.Close()
	
	caves := make([]*model.Cave, 0)
	for {
		//var caverIDstr string
		//var stamp int64
		var cave model.Cave

		rowExists, err := result.Step()
		if err != nil {
			db.log.Errorf("db.get: Step error: %s", err)
			return caves[0], err
		}

		if !rowExists {
			break
		}
		
		err = result.Scan(&cave.ID, &cave.Name, &cave.Region, &cave.Country, &cave.SRT, &cave.Visits)
		if err != nil {
			db.log.Error(err)
			return caves[0], err
		}

		// Add this formatted row to the rows map
		caves = append(caves, &cave)  
	}

	return caves[0], err
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
// For formatting the ids for a new Log
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
			if fullName == caver.First + `+` + caver.Last {
				caverIDs = append(caverIDs, caver.ID)
			}
		}
	}

	return strings.Join(caverIDs, `|`)
}

//
// For retrieving the names given a str of ids 
func (db *Database) getFullCaverNames(idStr string) string {
	// Get the IDs
	cavers, err := db.GetAllCavers()
	if err != nil {
		db.log.Errorf("Database.Query: Failed to fetch list of cavers")
		return ``
	}

	var names []string
	caverIDs := strings.Split(idStr, "|")

	for _, caver_id := range caverIDs {
		for _, caver := range cavers {
			if caver_id == caver.ID {
				fullName := caver.First + `+` + caver.Last
				names = append(names, fullName)
			}
		}	
	}
	
	return strings.Join(names, `, `)
}

//
// For retrieving the names given a str of ids 
func (db *Database) getCaverFirstNames(idStr string) string {
	// Get the IDs
	cavers, err := db.GetAllCavers()
	if err != nil {
		db.log.Errorf("Database.Query: Failed to fetch list of cavers")
		return ``
	}

	var names []string
	caverIDs := strings.Split(idStr, "|")

	for _, caver_id := range caverIDs {
		for _, caver := range cavers {
			if caver_id == caver.ID {
				names = append(names, caver.First)
			}
		}	
	}
	
	return strings.Join(names, `, `)
}

/*//
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
}*/
