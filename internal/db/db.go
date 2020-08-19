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
		log.Fatalf("database.new: Cannot establish connection to %s: %v", dbFN, err)
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
// ADD FUNCS ---- ----

func (db *Database) AddTrip(date, location, names, notes string) error {
	query := `INSERT INTO trips (date, caveid, notes) VALUES (?,?,?)`

	params, caverIDs, err := db.verifyTrip(date, location, names, notes)
	if err != nil {
		return err
	}

	if err = db.conn.Begin(); err != nil {
		return err
	}

	// Insert the trip itself
	tripID, err := db.execute(query, params)
	if err != nil {
		if rb_err := db.rollback(); rb_err != nil {
			panic(rb_err)
		}
		return err
	}

	// Insert the group of people
	_, err = db.execute(db.addTripGroups(tripID, caverIDs))
	if err != nil {
		if rb_err := db.rollback(); rb_err != nil {
			panic(rb_err)
		}
		return err
	}

	// If there are no errors commit changes
	if err = db.conn.Commit(); err != nil {
		if rb_err := db.rollback(); rb_err != nil {
			panic(rb_err)
		}
		return err
	}

	return nil
}

func (db *Database) AddLocation(name, region, country string, srt bool) error {
	query := `INSERT INTO locations (name, region, country, srt) VALUES (?,?,?,?)`
	params := []interface{}{name, region, country, srt}

	if err := db.conn.Begin(); err != nil {
		return err
	}

	_, err := db.execute(query, params)
	if err != nil {
		if rb_err := db.rollback(); rb_err != nil {
			panic(rb_err)
		}
		return err
	}

	// If there are no errors commit changes
	if err = db.conn.Commit(); err != nil {
		if rb_err := db.rollback(); rb_err != nil {
			panic(rb_err)
		}
		return err
	}

	return nil
}

func (db *Database) AddPerson(name, club string) error {
	query := `INSERT INTO people (name, club) VALUES (?,?)`
	params := []interface{}{name, club}

	if err := db.conn.Begin(); err != nil {
		return err
	}

	_, err := db.execute(query, params)
	if err != nil {
		if rb_err := db.rollback(); rb_err != nil {
			panic(rb_err)
		}
		return err
	}

	// If there are no errors commit changes
	if err = db.conn.Commit(); err != nil {
		if rb_err := db.rollback(); rb_err != nil {
			panic(rb_err)
		}
		return err
	}

	return nil
}

//
// GET FUNCS ---- ----

//
// ---- TRIPS FUNCS

func (db *Database) GetAllTrips() ([]*model.Log, error) {
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

func (db *Database) GetTrip(logID string) (*model.Log, error) { //FIXME: 
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
// ---- PEOPLE FUNCS

func (db *Database) GetAllPeople() ([]*model.Caver, error) {
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

func (db *Database) GetTopPeople() ([]*model.Statistic, error) {
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
		db.log.Errorf("db.GetTopPeople: Failed to get cavers", err)
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

func (db *Database) GetPerson(personID string) (*model.Caver, error) {
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
// ---- LOCATION FUNCS

func (db *Database) GetAllLocations() ([]*model.Cave, error) {
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

func (db *Database) GetTopLocations() ([]*model.Statistic, error) {
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
		db.log.Errorf("db.GetTopLocations: Failed to get caves", err)
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

func (db *Database) GetLocation(caveID string) (*model.Cave, error) {
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
// DELETE FUNCS ---- ----

func (db *Database) RemoveTrip(id string) error {
	if err := db.conn.Begin(); err != nil {
		return err
	}

	// Delete trip entry
	err := db.conn.Exec(`DELETE FROM trips WHERE id = ?`, id)
	if err != nil {
		if rb_err := db.rollback(); rb_err != nil {
			panic(rb_err)
		}
		return err
	}

	err = db.conn.Exec(`DELETE FROM trip_groups WHERE tripid = ?`, id)
	if err != nil {
		if rb_err := db.rollback(); rb_err != nil {
			panic(rb_err)
		}
		return err
	}

	// If there are no errors commit changes
	if err = db.conn.Commit(); err != nil {
		if rb_err := db.rollback(); rb_err != nil {
			panic(rb_err)
		}
		return err
	}

	return nil
}

func (db *Database) RemovePerson(id string) error {
	if err := db.conn.Begin(); err != nil {
		return err
	}

	// Delete trip entry
	err := db.conn.Exec(`DELETE FROM people WHERE id = ?`, id)
	if err != nil {
		if rb_err := db.rollback(); rb_err != nil {
			panic(rb_err)
		}
		return err
	}

	// If there are no errors commit changes
	if err = db.conn.Commit(); err != nil {
		if rb_err := db.rollback(); rb_err != nil {
			panic(rb_err)
		}
		return err
	}

	return nil
}

func (db *Database) RemoveLocation(id string) error {
	if err := db.conn.Begin(); err != nil {
		return err
	}

	// Delete trip entry
	err := db.conn.Exec(`DELETE FROM locations WHERE id = ?`, id)
	if err != nil {
		if rb_err := db.rollback(); rb_err != nil {
			panic(rb_err)
		}
		return err
	}

	// If there are no errors commit changes
	if err = db.conn.Commit(); err != nil {
		if rb_err := db.rollback(); rb_err != nil {
			panic(rb_err)
		}
		return err
	}

	return nil
}

//
// MODIFY FUNCS ---- ----

func (db *Database) ModifyTrip(id, date, location, names, notes string) error {
	query := `UPDATE trips SET date = ?, caveid = ?, notes = ? WHERE id = ?`

	params, caverIDs, err := db.verifyTrip(date, location, names, notes)
	if err != nil {
		return err
	}

	params = append(params, id)

	if err = db.conn.Begin(); err != nil {
		return err
	}

	// Update the trip itself
	_, err = db.execute(query, params)
	if err != nil {
		if rb_err := db.rollback(); rb_err != nil {
			panic(rb_err)
		}
		return err
	}

	// Update the group of people
	_, err = db.execute(`DELETE FROM trip_groups WHERE tripid = ?`, []interface{}{id})
	if err != nil {
		if rb_err := db.rollback(); rb_err != nil {
			panic(rb_err)
		}
		return err
	}

	_, err = db.execute(db.addTripGroups(id, caverIDs))
	if err != nil {
		if rb_err := db.rollback(); rb_err != nil {
			panic(rb_err)
		}
		return err
	}

	// If there are no errors commit changes
	if err = db.conn.Commit(); err != nil {
		if rb_err := db.rollback(); rb_err != nil {
			panic(rb_err)
		}
		return err
	}

	return nil
}

func (db *Database) ModifyPerson(id, name, club string) error {
	query := `UPDATE people SET name = ?, club = ? WHERE id = ?`
	params := []interface{}{name, club, id}

	if err := db.conn.Begin(); err != nil {
		return err
	}

	_, err := db.execute(query, params)
	if err != nil {
		if rb_err := db.rollback(); rb_err != nil {
			panic(rb_err)
		}
		return err
	}

	// If there are no errors commit changes
	if err = db.conn.Commit(); err != nil {
		if rb_err := db.rollback(); rb_err != nil {
			panic(rb_err)
		}
		return err
	}

	return nil
}

func (db *Database) ModifyLocation(id, name, region, country string, srt bool) error {
	query := `UPDATE locations SET name = ?, region = ?, country = ?, srt = ? WHERE id = ?`
	params := []interface{}{name, region, country, srt, id}

	if err := db.conn.Begin(); err != nil {
		return err
	}

	_, err := db.execute(query, params)
	if err != nil {
		if rb_err := db.rollback(); rb_err != nil {
			panic(rb_err)
		}
		return err
	}

	// If there are no errors commit changes
	if err = db.conn.Commit(); err != nil {
		if rb_err := db.rollback(); rb_err != nil {
			panic(rb_err)
		}
		return err
	}

	return nil
}
//
// INTERNAL FUNCTIONS ----------------------------------------------------------
//

func (db *Database) addTripGroups(tripID interface{}, caverIDs []string) (string, []interface{}) {
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

//
// Rollback database changes
func (db *Database) rollback() error {
	return db.conn.Rollback()
}

//
// Execute database query
func (db *Database) execute(query string, params []interface{}) (int64, error) {
	err := db.conn.Exec(query, params...)
	if err != nil {
		return -1, err
	}

	return db.conn.LastInsertRowID(), nil
}
