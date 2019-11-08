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

//
// ADD FUNCS

func (db *Database) AddLog(date, cave, names, notes string) error {
	query := `INSERT INTO trips (date, caveid, notes) VALUES (?,?,?)`

	d, err := time.Parse(datetime, strings.Join([]string{date,`12:00:00Z`},`T`))
	if err != nil {
		return err
	}
	dateStamp := d.Unix()

	caveID, err := db.getCaveID(cave)
	if err != nil {
		return err
	}

	caverIDs, err := db.getCaverIDs(names)
	if err != nil {
		return err
	}

	params := []interface{}{dateStamp, caveID, notes}
	
	if err = db.conn.Begin(); err != nil {
		return err
	}
	
	// Insert the trip itself
	tripID, err := db.insert(query, params)
	if err != nil {
		_ = db.rollback()
		return err
	}

	// Insert the group of people
	_, err = db.insert(db.addTripGroups(tripID, caverIDs))
	if err != nil {
		_ = db.rollback()
		return err
	}

	// If there are no errors commit changes
	if err = db.conn.Commit(); err != nil {
		_ = db.rollback()
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

//
// GET FUNCS

//
// TRIPS FUNCS

func (db *Database) GetAllLogs() ([]*model.Log, error) {
	// Build query
	var query string
	query = `
	SELECT
		trips.id AS 'id',
		trips.date AS 'date',
		locations.name AS 'cave',
		(
			SELECT GROUP_CONCAT(people.name, ', ')
			FROM trip_groups, people
			WHERE trip_groups.caverid = people.id
				AND trip_groups.tripid = trips.id
		) AS 'names',
		trips.notes AS 'notes'
	FROM trips, locations
	WHERE trips.caveid = locations.id`
	
	result, err := db.conn.Prepare(query)
	if err != nil {
		db.log.Errorf("db.prepare: Failed to query database", err)
		return nil, err
	}
	defer result.Close()
	
	trips := make([]*model.Log, 0)
	for {
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
		
		err = result.Scan(&trip.ID, &stamp, &trip.Cave, &trip.Names, &trip.Notes)
		if err != nil {
			db.log.Error(err)
			return trips, err
		}

		trip.Date = time.Unix(stamp, 0).Format(date)
		
		// Add this formatted row to the rows map
		trips = append(trips, &trip)  
	}

	return trips, err
}

func (db *Database) GetLog(logID string) (*model.Log, error) { //FIXME: 
	// Build query
	var query string
	query = `
	SELECT trips.id AS 'id',
		trips.date AS 'date',
		locations.name AS 'cave',
		(
			SELECT GROUP_CONCAT(people.name, ', ')
			FROM trip_groups, people
			WHERE trip_groups.caverid = people.id
				AND trip_groups.tripid = trips.id
		) AS 'names',
		trips.notes AS 'notes'
	FROM trips, locations
	WHERE trips.caveid = locations.id
		AND trips.id = ?`

	result, err := db.conn.Prepare(query, logID)
	if err != nil {
		db.log.Errorf("db.prepare: Failed to query database", err)
		return nil, err
	}
	defer result.Close()
	
	trips := make([]*model.Log, 0)
	for {
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
		
		err = result.Scan(&trip.ID, &stamp, &trip.Cave, &trip.Names, &trip.Notes)
		if err != nil {
			db.log.Error(err)
			return trips[0], err
		}

		trip.Date = time.Unix(stamp, 0).Format(date)
		
		// Add this formatted row to the rows map
		trips = append(trips, &trip)  
	}

	return trips[0], err
}

//
// PEOPLE FUNCS

func (db *Database) GetAllCavers() ([]*model.Caver, error) {
	var query string
	query = `
	SELECT 
		people.id AS 'id',
		people.name AS 'name',
		people.club AS 'club',
		(
			SELECT COUNT(1)
				FROM trip_groups
			 WHERE trip_groups.caverid = people.id
		)
	FROM people
	ORDER BY name`
	result, err := db.conn.Prepare(query)
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
		
		err = result.Scan(&c.ID, &c.Name, &c.Club, &c.Count)
		if err != nil {
			db.log.Errorf("Scan: %v", err)
		}
		cavers = append(cavers, &c)
	}
	return cavers, err
}

func (db *Database) GetTopCavers() ([]*model.Statistic, error) {
	var query string
	query = `
	SELECT 
		people.name AS 'name',
		(
			SELECT COUNT(1)
				FROM trip_groups
			 WHERE trip_groups.caverid = people.id
		) AS count
	FROM people
	ORDER BY count DESC LIMIT 15`
	result, err := db.conn.Prepare(query)
	if err != nil {
		db.log.Errorf("db.gettopcavers: Failed to get cavers", err)
	}

	cavers := make([]*model.Statistic, 0)
	for {
		var c model.Statistic
		
		rowExists, err := result.Step()
		if err != nil {
			db.log.Errorf("db.get: Step error: %s", err)
			return cavers, err
		}

		if !rowExists {
			break
		}
		
		err = result.Scan(&c.Name, &c.Value)
		if err != nil {
			db.log.Errorf("Scan: %v", err)
		}
		cavers = append(cavers, &c)
	}
	return cavers, err
}

func (db *Database) GetCaver(personID string) (*model.Caver, error) {
	// Build query
	var query string
	query = `
	SELECT 
		people.id AS 'id',
		people.name AS 'name',
		people.club AS 'club',
		(
			SELECT COUNT(1)
				FROM trip_groups
			WHERE trip_groups.caverid = people.id
		) AS 'count'
	FROM people
	WHERE people.id = ?`

	result, err := db.conn.Prepare(query, personID)
	if err != nil {
		db.log.Errorf("db.prepare: Failed to query database", err)
		return nil, err
	}
	defer result.Close()
	
	people := make([]*model.Caver, 0)
	for {
		var person model.Caver

		rowExists, err := result.Step()
		if err != nil {
			db.log.Errorf("db.get: Step error: %s", err)
			return people[0], err
		}

		if !rowExists {
			break
		}
		
		err = result.Scan(&person.ID, &person.Name, &person.Club, &person.Count)
		if err != nil {
			db.log.Error(err)
			return people[0], err
		}

		// Add this formatted row to the rows map
		people = append(people, &person)  
	}

	return people[0], err
}

//
// LOCATION FUNCS

func (db *Database) GetAllCaves() ([]*model.Cave, error) {
	var query string
	query = `
	SELECT
		locations.id AS 'id',
		locations.name AS 'name',
		locations.region AS 'region',
		locations.country AS 'country',
		locations.srt AS 'srt',
		(
			SELECT COUNT(1)
			FROM trips
			WHERE trips.caveid = locations.id
		) AS 'visits'
	FROM locations
	ORDER BY name`
	result, err := db.conn.Prepare(query)
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
	}
	return caves, err
}

func (db *Database) GetTopCaves() ([]*model.Statistic, error) {
	var query string
	query = `
	SELECT
		locations.name AS 'name',
		(
			SELECT COUNT(1)
			FROM trips
			WHERE trips.caveid = locations.id
		) AS visits
	FROM locations
	ORDER BY visits DESC LIMIT 15`
	result, err := db.conn.Prepare(query)
	if err != nil {
		db.log.Errorf("db.gettopcaves: Failed to get caves", err)
	}

	stats := make([]*model.Statistic, 0)
	for {
		var s model.Statistic
		
		rowExists, err := result.Step()
		if err != nil {
			db.log.Errorf("db.get: Step error: %s", err)
			return stats, err
		}

		if !rowExists {
			break
		}
		
		err = result.Scan(&s.Name, &s.Value)
		if err != nil {
			db.log.Errorf("Scan: %v", err)
		}
		stats = append(stats, &s)
	}
	return stats, err
}

func (db *Database) GetCave(caveID string) (*model.Cave, error) {
	// Build query
	query := `
	SELECT
		locations.id AS 'id',
		locations.name AS 'name',
		locations.region AS 'region',
		locations.country AS 'country',
		locations.srt AS 'srt',
		(
			SELECT COUNT(1)
			FROM trips
			WHERE trips.caveid = locations.id
		) AS 'visits'
	FROM locations
	WHERE id = ?`
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


//
// DELETE FUNCS

func (db *Database) RemoveLog(logID string) error {
	return nil
}

//
// INTERNAL FUNCTIONS ----------------------------------------------------------
//

func (db *Database) addTripGroups(tripID int64, caverIDs []string) (string, []interface{}) {
	// Build the query
	paramPlaceholderTemplate := `(?,?)`
	var paramPlaceholderGroup []string

	for i := 0; i < len(caverIDs); i++ {
		paramPlaceholderGroup = append(paramPlaceholderGroup, paramPlaceholderTemplate)
	}
	paramPlaceholder := strings.Join(paramPlaceholderGroup, `,`)

	query := `INSERT INTO trip_groups (tripid, caverid) VALUES x`
	query = strings.Replace(query, `x`, paramPlaceholder, 1)

	// Build the parameters
	var params []interface{}
	for _, caverID := range caverIDs {
		params = append(params, tripID, caverID)
	}

	return query, params
}

func (db *Database) rollback() error {
	return db.conn.Rollback()
}

func (db *Database) insert(query string, params []interface{}) (int64, error) {
	err := db.conn.Exec(query, params...)
	if err != nil {
		return -1, err
	}

	return db.conn.LastInsertRowID(), nil
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
