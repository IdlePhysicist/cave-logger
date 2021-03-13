package register

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

type Register struct {
	log  *logrus.Logger
	db   *sql.DB
	ctx  context.Context
}

func New(log *logrus.Logger, dbFN string) *Register {
	var reg Register

	reg.log = log
	reg.ctx = context.Background()

	db, err := sql.Open("sqlite", dbFN)
	if err != nil {
		log.Fatalf("database.new: Cannot establish connection to %s: %v", dbFN, err)
	}
	reg.db = db
	return &reg
}

func (reg *Register) Close() {
	reg.db.Close()
}

//
// MAIN FUNCTIONS --------------------------------------------------------------
//

//
// ADD FUNCS ---- ----

func (reg *Register) AddTrip(date, location, names, notes string) error {
	query := `INSERT INTO trips (date, caveid, notes) VALUES (?,?,?)`

	params, caverIDs, err := reg.verifyTrip(date, location, names, notes)
	if err != nil {
		return err
	}

	trans, err := reg.db.Begin()
	if err != nil {
		return err
	}

	// Insert the trip itself
	tripID, err := reg.execute(query, params, trans)
	if err != nil {
		if rb_err := trans.Rollback(); rb_err != nil {
			panic(rb_err)
		}
		return err
	}

	// Insert the group of people
	tgQuery, tgParams := reg.addTripGroups(tripID, caverIDs)
	_, err = reg.execute(tgQuery, tgParams, trans)
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

func (reg *Register) AddCave(name, region, country, notes string, srt bool) error {
	query := `INSERT INTO locations (name, region, country, srt, notes) VALUES (?,?,?,?,?)`
	params := []interface{}{name, region, country, srt, notes}

	trans, err := reg.db.Begin()
	if err != nil {
		return err
	}

	_, err = reg.execute(query, params, trans)
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

func (reg *Register) AddCaver(name, club, notes string) error {
	query := `INSERT INTO people (name, club, notes) VALUES (?,?,?)`
	params := []interface{}{name, club, notes}

	trans, err := reg.db.Begin()
	if err != nil {
		return err
	}

	_, err = reg.execute(query, params, trans)
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

func (reg *Register) GetAllTrips() ([]*model.Log, error) {
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

	result, err := reg.db.Query(query)
	if err != nil {
		reg.log.Errorf("reg.prepare: Failed to query database", err)
		return nil, err
	}
	defer result.Close()

	trips := make([]*model.Log, 0)
	for result.Next() {
		var stamp int64
		var trip model.Log

		err = result.Scan(&trip.ID, &stamp, &trip.Cave, &trip.Names, &trip.Notes)
		if err != nil {
			reg.log.Error(err)
			return trips, err
		}

		trip.Date = time.Unix(stamp, 0).Format(date)

		// Add this formatted row to the rows map
		trips = append(trips, &trip)
	}
	if err = result.Err(); err != nil {
		reg.log.Errorf("reg.get: Step error: %s", err)
		return trips, err
	}

	return trips, err
}

func (reg *Register) GetTrip(id string) (*model.Log, error) {
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

	result, err := reg.db.Query(query, id)
	if err != nil {
		reg.log.Errorf("reg.prepare: Failed to query database", err)
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
			reg.log.Error(err)
			return nil, err
		}

		trip.Date = time.Unix(stamp, 0).Format(date)
	}
	if err = result.Err(); err != nil {
		reg.log.Errorf("reg.get: Step error: %s", err)
		return nil, err
	}

	return &trip, err
}

//
// ---- PEOPLE FUNCS

func (reg *Register) GetAllCavers() ([]*model.Caver, error) {
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

	result, err := reg.db.Query(query)
	if err != nil {
		reg.log.Errorf("reg.getcaverlist: Failed to get cavers", err)
	}

	cavers := make([]*model.Caver, 0)
	for result.Next() {
		var c model.Caver

		err = result.Scan(&c.ID, &c.Name, &c.Club, &c.Count)
		if err != nil {
			reg.log.Errorf("Scan: %v", err)
		}
		cavers = append(cavers, &c)
	}
	if err = result.Err(); err != nil {
		reg.log.Errorf("reg.get: Step error: %s", err)
		return cavers, err
	}
	return cavers, err
}

/*
func (reg *Register) GetTopPeople() ([]*model.Statistic, error) {
	query := `
    SELECT 
        people.name AS 'name',
        (
            SELECT COUNT(1) FROM trip_groups WHERE trip_groups.caverid = people.id
        ) AS count
    FROM people
    ORDER BY count DESC LIMIT 15`
	result, err := reg.db.Prepare(query)
	if err != nil {
		reg.log.Errorf("reg.GetTopPeople: Failed to get cavers", err)
	}

	cavers := make([]*model.Statistic, 0)
	for result.Next() {
		var c model.Statistic

		err = result.Scan(&c.Name, &c.Value)
		if err != nil {
			reg.log.Errorf("Scan: %v", err)
		}
		cavers = append(cavers, &c)
	}
	if err = result.Err(); err != nil {
		reg.log.Errorf("reg.get: Step error: %s", err)
		return cavers, err
	}
	return cavers, err
}
*/

func (reg *Register) GetCaver(id string) (*model.Caver, error) {
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

	result, err := reg.db.Query(query, id)
	if err != nil {
		reg.log.Errorf("reg.prepare: Failed to query database", err)
		return nil, err
	}
	defer result.Close()

	var caver model.Caver
	for result.Next() {

		err = result.Scan(&caver.ID, &caver.Name, &caver.Club, &caver.Count, &caver.Notes)
		if err != nil {
			reg.log.Errorf("reg.scan", err)
			return nil, err
		}
	}
	if err = result.Err(); err != nil {
		reg.log.Errorf("reg.get: Step error: %s", err)
		return nil, err
	}

	return &caver, err
}

//
// ---- LOCATION FUNCS

func (reg *Register) GetAllCaves() ([]*model.Cave, error) {
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
	result, err := reg.db.Query(query)
	if err != nil {
		reg.log.Errorf("reg.getcaverlist: Failed to get cavers", err)
	}

	caves := make([]*model.Cave, 0)
	for result.Next() {
		var c model.Cave

		err = result.Scan(&c.ID, &c.Name, &c.Region, &c.Country, &c.SRT, &c.Visits)
		if err != nil {
			reg.log.Errorf("Scan: %v", err)
		}
		caves = append(caves, &c)
	}
	if err = result.Err(); err != nil {
		reg.log.Errorf("reg.get: Step error: %s", err)
		return caves, err
	}
	return caves, err
}

/*
func (reg *Register) GetTopLocations() ([]*model.Statistic, error) {
	query := `
    SELECT
        locations.name AS 'name',
        (
            SELECT COUNT(1) FROM trips WHERE trips.caveid = locations.id
        ) AS visits
    FROM locations
    ORDER BY visits DESC LIMIT 15`
	result, err := reg.db.Prepare(query)
	if err != nil {
		reg.log.Errorf("reg.GetTopLocations: Failed to get caves", err)
	}

	stats := make([]*model.Statistic, 0)
	for result.Next() {
		var s model.Statistic

		err = result.Scan(&s.Name, &s.Value)
		if err != nil {
			reg.log.Errorf("Scan: %v", err)
		}
		stats = append(stats, &s)
	}
	if err = result.Err(); err != nil {
		reg.log.Errorf("reg.get: Step error: %s", err)
		return stats, err
	}
	return stats, err
}
*/

func (reg *Register) GetCave(id string) (*model.Cave, error) {
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

	result, err := reg.db.Query(query, id)
	if err != nil {
		reg.log.Errorf("reg.prepare: Failed to query database", err)
		return nil, err
	}
	defer result.Close()

	var cave model.Cave
	for result.Next() {
		err = result.Scan(&cave.ID, &cave.Name, &cave.Region, &cave.Country, &cave.SRT, &cave.Visits, &cave.Notes)
		if err != nil {
			reg.log.Error(err)
			return nil, err
		}
	}
	if err = result.Err(); err != nil {
		reg.log.Errorf("reg.get: Step error: %s", err)
		return nil, err
	}

	return &cave, err
}


//
// DELETE FUNCS ---- ----

func (reg *Register) RemoveTrip(id string) error {
	trans, err := reg.db.Begin()
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

func (reg *Register) RemoveCaver(id string) error {
	trans, err := reg.db.Begin()
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

func (reg *Register) RemoveCave(id string) error {
	trans, err := reg.db.Begin()
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

func (reg *Register) ModifyTrip(id, date, location, names, notes string) error {
	query := `UPDATE trips SET date = ?, caveid = ?, notes = ? WHERE id = ?`

	params, caverIDs, err := reg.verifyTrip(date, location, names, notes)
	if err != nil {
		return err
	}

	params = append(params, id)

	trans, err := reg.db.Begin()
	if err != nil {
		return err
	}

	// Update the trip itself
	_, err = reg.execute(query, params, trans)
	if err != nil {
		if rb_err := trans.Rollback(); rb_err != nil {
			panic(rb_err)
		}
		return err
	}

	// Update the group of people
	_, err = reg.execute(`DELETE FROM trip_groups WHERE tripid = ?`, []interface{}{id}, trans)
	if err != nil {
		if rb_err := trans.Rollback(); rb_err != nil {
			panic(rb_err)
		}
		return err
	}

	tgQuery, tgParams := reg.addTripGroups(id, caverIDs)
	_, err = reg.execute(tgQuery, tgParams, trans)
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

func (reg *Register) ModifyCaver(id, name, club, notes string) error {
	query := `UPDATE people SET name = ?, club = ?, notes = ? WHERE id = ?`
	params := []interface{}{name, club, notes, id}

	trans, err := reg.db.Begin()
	if err != nil {
		return err
	}

	_, err = reg.execute(query, params, trans)
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

func (reg *Register) ModifyCave(id, name, region, country, notes string, srt bool) error {
	query := `UPDATE locations SET name = ?, region = ?, country = ?, srt = ?, notes = ? WHERE id = ?`
	params := []interface{}{name, region, country, srt, notes, id}

	trans, err := reg.db.Begin()
	if err != nil {
		return err
	}

	_, err = reg.execute(query, params, trans)
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
func (reg *Register) addTripGroups(tripID interface{}, caverIDs []string) (string, []interface{}) {
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
func (reg *Register) execute(query string, params []interface{}, trans *sql.Tx) (int64, error) {
	var (
		result sql.Result
		err    error
	)

	if trans == nil {
		result, err = trans.Exec(query, params...)
	} else {
		result, err = reg.db.Exec(query, params...)
	}
	if err != nil {
		return -1, err
	}

	return result.LastInsertId()
}
