package db

import (
	"context"
	"database/sql"
	"strings"
	"time"

	_ "modernc.org/sqlite"
	"github.com/sirupsen/logrus"

	"github.com/idlephysicist/cave-logger/internal/model"
)

const datetime = `2006-01-02T15:04:05Z`
const date = `2006-01-02`

type Database struct {
	log  *logrus.Logger
	conn *sql.DB
	ctx  context.Context
}

func New(log *logrus.Logger, dbFN string) *Database {
	var db Database

	db.log = log
	db.ctx = context.Background()

	conn, err := sql.Open("sqlite", dbFN)
	if err != nil {
		log.Fatalf("database.new: Cannot establish connection to %s: %v", dbFN, err)
	}
	db.conn = conn
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

	trans, err := db.conn.Begin()
	if err != nil {
		return err
	}

	// Insert the trip itself
	tripID, err := db.execute(query, params, trans)
	if err != nil {
		if rb_err := trans.Rollback(); rb_err != nil {
			panic(rb_err)
		}
		return err
	}

	// Insert the group of people
	tgQuery, tgParams := db.addTripGroups(tripID, caverIDs)
	_, err = db.execute(tgQuery, tgParams, trans)
	if err != nil {
		if rb_err := trans.Rollback(); rb_err != nil {
			panic(rb_err)
		}
		return err
	}

	// If there are no errors commit changes
	if err = trans.Commit(); err != nil {
		if rb_err := trans.Rollback(); rb_err != nil {
			panic(rb_err)
		}
		return err
	}

	return nil
}

func (db *Database) AddLocation(name, region, country, notes string, srt bool) error {
	query := `INSERT INTO locations (name, region, country, srt, notes) VALUES (?,?,?,?,?)`
	params := []interface{}{name, region, country, srt, notes}

	trans, err := db.conn.Begin()
	if err != nil {
		return err
	}

	_, err = db.execute(query, params, trans)
	if err != nil {
		if rb_err := trans.Rollback(); rb_err != nil {
			panic(rb_err)
		}
		return err
	}

	// If there are no errors commit changes
	if err = trans.Commit(); err != nil {
		if rb_err := trans.Rollback(); rb_err != nil {
			panic(rb_err)
		}
		return err
	}

	return nil
}

func (db *Database) AddPerson(name, club, notes string) error {
	query := `INSERT INTO people (name, club, notes) VALUES (?,?,?)`
	params := []interface{}{name, club, notes}

	trans, err := db.conn.Begin()
	if err != nil {
		return err
	}

	_, err = db.execute(query, params, trans)
	if err != nil {
		if rb_err := trans.Rollback(); rb_err != nil {
			panic(rb_err)
		}
		return err
	}

	// If there are no errors commit changes
	if err = trans.Commit(); err != nil {
		if rb_err := trans.Rollback(); rb_err != nil {
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
	query := `
    SELECT
        trips.id AS 'id',
        trips.date AS 'date',
        locations.name AS 'cave',
        (
            SELECT GROUP_CONCAT(people.name, ', ')
            FROM trip_groups, people
            WHERE trip_groups.caverid = people.id AND trip_groups.tripid = trips.id
        ) AS 'names',
        trips.notes AS 'notes'
    FROM trips, locations
    WHERE trips.caveid = locations.id`

	result, err := db.conn.Query(query)
	if err != nil {
		db.log.Errorf("db.prepare: Failed to query database", err)
		return nil, err
	}
	defer result.Close()

	trips := make([]*model.Log, 0)
	for result.Next() {
		var stamp int64
		var trip model.Log

		err = result.Scan(&trip.ID, &stamp, &trip.Cave, &trip.Names, &trip.Notes)
		if err != nil {
			db.log.Error(err)
			return trips, err
		}

		trip.Date = time.Unix(stamp, 0).Format(date)

		// Add this formatted row to the rows map
		trips = append(trips, &trip)
	}
	if err = result.Err(); err != nil {
		db.log.Errorf("db.get: Step error: %s", err)
		return trips, err
	}

	return trips, err
}

func (db *Database) GetTrip(tripID string) (*model.Log, error) {
	query := `
    SELECT trips.id AS 'id',
        trips.date AS 'date',
        locations.name AS 'cave',
        (
            SELECT GROUP_CONCAT(people.name, ', ')
            FROM trip_groups, people
            WHERE trip_groups.caverid = people.id AND trip_groups.tripid = trips.id
        ) AS 'names',
        trips.notes AS 'notes'
    FROM trips, locations
    WHERE trips.caveid = locations.id AND trips.id = ?`

	result, err := db.conn.Query(query, tripID)
	if err != nil {
		db.log.Errorf("db.prepare: Failed to query database", err)
		return nil, err
	}
	defer result.Close()

	var (
		trip  model.Log
		stamp int64
	)
	for result.Next() {
		err = result.Scan(&trip.ID, &stamp, &trip.Cave, &trip.Names, &trip.Notes)
		if err != nil {
			db.log.Error(err)
			return nil, err
		}

		trip.Date = time.Unix(stamp, 0).Format(date)
	}
	if err = result.Err(); err != nil {
		db.log.Errorf("db.get: Step error: %s", err)
		return nil, err
	}

	return &trip, err
}

//
// ---- PEOPLE FUNCS

func (db *Database) GetAllPeople() ([]*model.Caver, error) {
	query := `
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

	result, err := db.conn.Query(query)
	if err != nil {
		db.log.Errorf("db.getcaverlist: Failed to get cavers", err)
	}

	cavers := make([]*model.Caver, 0)
	for result.Next() {
		var c model.Caver

		err = result.Scan(&c.ID, &c.Name, &c.Club, &c.Count)
		if err != nil {
			db.log.Errorf("Scan: %v", err)
		}
		cavers = append(cavers, &c)
	}
	if err = result.Err(); err != nil {
		db.log.Errorf("db.get: Step error: %s", err)
		return cavers, err
	}
	return cavers, err
}

/*
func (db *Database) GetTopPeople() ([]*model.Statistic, error) {
	query := `
    SELECT 
        people.name AS 'name',
        (
            SELECT COUNT(1) FROM trip_groups WHERE trip_groups.caverid = people.id
        ) AS count
    FROM people
    ORDER BY count DESC LIMIT 15`
	result, err := db.conn.Prepare(query)
	if err != nil {
		db.log.Errorf("db.GetTopPeople: Failed to get cavers", err)
	}

	cavers := make([]*model.Statistic, 0)
	for result.Next() {
		var c model.Statistic

		err = result.Scan(&c.Name, &c.Value)
		if err != nil {
			db.log.Errorf("Scan: %v", err)
		}
		cavers = append(cavers, &c)
	}
	if err = result.Err(); err != nil {
		db.log.Errorf("db.get: Step error: %s", err)
		return cavers, err
	}
	return cavers, err
}
*/

func (db *Database) GetPerson(caverID string) (*model.Caver, error) {
	query := `
    SELECT 
        people.id AS 'id',
        people.name AS 'name',
        people.club AS 'club',
        (
            SELECT COUNT(1) FROM trip_groups WHERE trip_groups.caverid = people.id
        ) AS 'count',
        people.notes AS 'notes'
    FROM people
    WHERE people.id = ?`

	result, err := db.conn.Query(query, caverID)
	if err != nil {
		db.log.Errorf("db.prepare: Failed to query database", err)
		return nil, err
	}
	defer result.Close()

	var caver model.Caver
	for result.Next() {

		err = result.Scan(&caver.ID, &caver.Name, &caver.Club, &caver.Count, &caver.Notes)
		if err != nil {
			db.log.Errorf("db.scan", err)
			return nil, err
		}
	}
	if err = result.Err(); err != nil {
		db.log.Errorf("db.get: Step error: %s", err)
		return nil, err
	}

	return &caver, err
}

//
// ---- LOCATION FUNCS

func (db *Database) GetAllLocations() ([]*model.Cave, error) {
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
    ORDER BY name`
	result, err := db.conn.Query(query)
	if err != nil {
		db.log.Errorf("db.getcaverlist: Failed to get cavers", err)
	}

	caves := make([]*model.Cave, 0)
	for result.Next() {
		var c model.Cave

		err = result.Scan(&c.ID, &c.Name, &c.Region, &c.Country, &c.SRT, &c.Visits)
		if err != nil {
			db.log.Errorf("Scan: %v", err)
		}
		caves = append(caves, &c)
	}
	if err = result.Err(); err != nil {
		db.log.Errorf("db.get: Step error: %s", err)
		return caves, err
	}
	return caves, err
}

/*
func (db *Database) GetTopLocations() ([]*model.Statistic, error) {
	query := `
    SELECT
        locations.name AS 'name',
        (
            SELECT COUNT(1) FROM trips WHERE trips.caveid = locations.id
        ) AS visits
    FROM locations
    ORDER BY visits DESC LIMIT 15`
	result, err := db.conn.Prepare(query)
	if err != nil {
		db.log.Errorf("db.GetTopLocations: Failed to get caves", err)
	}

	stats := make([]*model.Statistic, 0)
	for result.Next() {
		var s model.Statistic

		err = result.Scan(&s.Name, &s.Value)
		if err != nil {
			db.log.Errorf("Scan: %v", err)
		}
		stats = append(stats, &s)
	}
	if err = result.Err(); err != nil {
		db.log.Errorf("db.get: Step error: %s", err)
		return stats, err
	}
	return stats, err
}
*/

func (db *Database) GetLocation(caveID string) (*model.Cave, error) {
	query := `
    SELECT
        locations.id AS 'id',
        locations.name AS 'name',
        locations.region AS 'region',
        locations.country AS 'country',
        locations.srt AS 'srt',
        (
            SELECT COUNT(1) FROM trips WHERE trips.caveid = locations.id
        ) AS 'visits',
        locations.notes AS 'notes'
    FROM locations
    WHERE id = ?`

	result, err := db.conn.Query(query, caveID)
	if err != nil {
		db.log.Errorf("db.prepare: Failed to query database", err)
		return nil, err
	}
	defer result.Close()

	var cave model.Cave
	for result.Next() {
		err = result.Scan(&cave.ID, &cave.Name, &cave.Region, &cave.Country, &cave.SRT, &cave.Visits, &cave.Notes)
		if err != nil {
			db.log.Error(err)
			return nil, err
		}
	}
	if err = result.Err(); err != nil {
		db.log.Errorf("db.get: Step error: %s", err)
		return nil, err
	}

	return &cave, err
}


//
// DELETE FUNCS ---- ----

func (db *Database) RemoveTrip(id string) error {
	trans, err := db.conn.Begin()
	if err != nil {
		return err
	}

	// Delete trip entry
	_, err = trans.Exec(`DELETE FROM trips WHERE id = ?`, id)
	if err != nil {
		if rb_err := trans.Rollback(); rb_err != nil {
			panic(rb_err)
		}
		return err
	}

	_, err = trans.Exec(`DELETE FROM trip_groups WHERE tripid = ?`, id)
	if err != nil {
		if rb_err := trans.Rollback(); rb_err != nil {
			panic(rb_err)
		}
		return err
	}

	// If there are no errors commit changes
	if err = trans.Commit(); err != nil {
		if rb_err := trans.Rollback(); rb_err != nil {
			panic(rb_err)
		}
		return err
	}

	return nil
}

func (db *Database) RemovePerson(id string) error {
	trans, err := db.conn.Begin()
	if err != nil {
		return err
	}

	// Delete trip entry
	_, err = trans.Exec(`DELETE FROM people WHERE id = ?`, id)
	if err != nil {
		if rb_err := trans.Rollback(); rb_err != nil {
			panic(rb_err)
		}
		return err
	}

	// If there are no errors commit changes
	if err = trans.Commit(); err != nil {
		if rb_err := trans.Rollback(); rb_err != nil {
			panic(rb_err)
		}
		return err
	}

	return nil
}

func (db *Database) RemoveLocation(id string) error {
	trans, err := db.conn.Begin()
	if err != nil {
		return err
	}

	// Delete trip entry
	_, err = trans.Exec(`DELETE FROM locations WHERE id = ?`, id)
	if err != nil {
		if rb_err := trans.Rollback(); rb_err != nil {
			panic(rb_err)
		}
		return err
	}

	// If there are no errors commit changes
	if err = trans.Commit(); err != nil {
		if rb_err := trans.Rollback(); rb_err != nil {
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

	trans, err := db.conn.Begin()
	if err != nil {
		return err
	}

	// Update the trip itself
	_, err = db.execute(query, params, trans)
	if err != nil {
		if rb_err := trans.Rollback(); rb_err != nil {
			panic(rb_err)
		}
		return err
	}

	// Update the group of people
	_, err = db.execute(`DELETE FROM trip_groups WHERE tripid = ?`, []interface{}{id}, trans)
	if err != nil {
		if rb_err := trans.Rollback(); rb_err != nil {
			panic(rb_err)
		}
		return err
	}

	tgQuery, tgParams := db.addTripGroups(id, caverIDs)
	_, err = db.execute(tgQuery, tgParams, trans)
	if err != nil {
		if rb_err := trans.Rollback(); rb_err != nil {
			panic(rb_err)
		}
		return err
	}

	// If there are no errors commit changes
	if err = trans.Commit(); err != nil {
		if rb_err := trans.Rollback(); rb_err != nil {
			panic(rb_err)
		}
		return err
	}

	return nil
}

func (db *Database) ModifyPerson(id, name, club, notes string) error {
	query := `UPDATE people SET name = ?, club = ?, notes = ? WHERE id = ?`
	params := []interface{}{name, club, notes, id}

	trans, err := db.conn.Begin()
	if err != nil {
		return err
	}

	_, err = db.execute(query, params, trans)
	if err != nil {
		if rb_err := trans.Rollback(); rb_err != nil {
			panic(rb_err)
		}
		return err
	}

	// If there are no errors commit changes
	if err = trans.Commit(); err != nil {
		if rb_err := trans.Rollback(); rb_err != nil {
			panic(rb_err)
		}
		return err
	}

	return nil
}

func (db *Database) ModifyLocation(id, name, region, country, notes string, srt bool) error {
	query := `UPDATE locations SET name = ?, region = ?, country = ?, srt = ?, notes = ? WHERE id = ?`
	params := []interface{}{name, region, country, srt, notes, id}

	trans, err := db.conn.Begin()
	if err != nil {
		return err
	}

	_, err = db.execute(query, params, trans)
	if err != nil {
		if rb_err := trans.Rollback(); rb_err != nil {
			panic(rb_err)
		}
		return err
	}

	// If there are no errors commit changes
	if err = trans.Commit(); err != nil {
		if rb_err := trans.Rollback(); rb_err != nil {
			panic(rb_err)
		}
		return err
	}

	return nil
}
//
// INTERNAL FUNCTIONS ----------------------------------------------------------
//

// addTripGroups builds a SQL statement and slice of parameters representing the
// trip group.
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
// Execute database query
func (db *Database) execute(query string, params []interface{}, trans *sql.Tx) (int64, error) {
	var (
		result sql.Result
		err    error
	)

	if trans == nil {
		result, err = trans.Exec(query, params...)
	} else {
		result, err = db.conn.Exec(query, params...)
	}
	if err != nil {
		return -1, err
	}

	return result.LastInsertId()
}
