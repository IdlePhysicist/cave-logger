package register

import (
	"strings"

	"github.com/idlephysicist/cave-logger/internal/model"
)

func (reg *Register) AddTrip(date, location, names, notes string) error {
	query := `INSERT INTO trip (date, caveid, notes) VALUES (?,?,?)`

	params, caverIDs, err := reg.verifyTrip(date, location, names, notes)
	if err != nil {
		return err
	}

	// Insert the trip itself
	tripID, err := reg.executeTx(query, params)
	if err != nil {
		return err
	}

	// Insert the group of people
	tgQuery, tgParams := reg.addTripGroups(tripID, caverIDs)
	_, err = reg.executeTx(tgQuery, tgParams)
	return err
}

func (reg *Register) GetAllTrips() ([]*model.Log, error) {
	query := `
    SELECT
        trip.id AS 'id',
        trip.date AS 'date',
        cave.name AS 'cave',
        (
            SELECT GROUP_CONCAT(caver.name, ', ')
            FROM trip_group, caver
            WHERE trip_group.caverid = caver.id AND trip_group.tripid = trip.id
        ) AS 'names',
        trip.notes AS 'notes'
    FROM trip, cave
    WHERE trip.caveid = cave.id
    ORDER BY date ASC`

	result, err := reg.db.Query(query)
	if err != nil {
		// reg.log.Errorf("reg.prepare: Failed to query database %w", err)
		return nil, err
	}
	defer result.Close()

	trips := make([]*model.Log, 0)
	for result.Next() {
		var trip model.Log

		err = result.Scan(&trip.ID, &trip.Date, &trip.Cave, &trip.Names, &trip.Notes)
		if err != nil {
			// reg.log.Error(err)
			return trips, err
		}

		trips = append(trips, &trip)
	}
	if err = result.Err(); err != nil { // REVEIW: Should I be checking this earlier?
		// reg.log.Errorf("reg.get: Step error: %s", err)
		return trips, err
	}

	return trips, err
}

func (reg *Register) GetTrip(id string) (*model.Log, error) {
	query := `
    SELECT trip.id AS 'id',
        trip.date AS 'date',
        cave.name AS 'cave',
        (
            SELECT GROUP_CONCAT(caver.name, ', ')
            FROM trip_group, caver
            WHERE trip_group.caverid = caver.id AND trip_group.tripid = trip.id
        ) AS 'names',
        trip.notes AS 'notes'
    FROM trip, cave
    WHERE trip.caveid = cave.id AND trip.id = ?`

	result, err := reg.db.Query(query, id)
	if err != nil {
		// reg.log.Errorf("reg.prepare: Failed to query database", err)
		return nil, err
	}
	defer result.Close()

	var trip model.Log
	for result.Next() {
		err = result.Scan(&trip.ID, &trip.Date, &trip.Cave, &trip.Names, &trip.Notes)
		if err != nil {
			// reg.log.Error(err)
			return nil, err
		}

	}
	if err = result.Err(); err != nil {
		// reg.log.Errorf("reg.get: Step error: %s", err)
		return nil, err
	}

	return &trip, err
}

func (reg *Register) ModifyTrip(id, date, location, names, notes string) error {
	query := `UPDATE trip SET date = ?, caveid = ?, notes = ? WHERE id = ?`

	params, caverIDs, err := reg.verifyTrip(date, location, names, notes)
	if err != nil {
		return err
	}

	params = append(params, id)

	// trans, err := reg.db.Begin()
	// if err != nil {
	// 	return err
	// }

	// Update the trip itself
	_, err = reg.executeTx(query, params)
	if err != nil {
		return err
	}

	// Update the group of people
	_, err = reg.executeTx(`DELETE FROM trip_group WHERE tripid = ?`, []interface{}{id})
	if err != nil {
		return err
	}

	tgQuery, tgParams := reg.addTripGroups(id, caverIDs)
	_, err = reg.executeTx(tgQuery, tgParams)
	return err
}

func (reg *Register) RemoveTrip(id string) error {
	_, err := reg.executeTx(`DELETE FROM trip WHERE id = ?`, []any{id})
	if err != nil {
		return err
	}

	_, err = reg.executeTx(`DELETE FROM trip_group WHERE tripid = ?`, []any{id})
	return err
}

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

	query := `INSERT INTO trip_group (tripid, caverid) VALUES x`
	query = strings.Replace(query, `x`, paramPlaceholder, 1)

	// Build the parameters
	var params []interface{}
	for _, caverID := range caverIDs {
		params = append(params, tripID, caverID)
	}

	return query, params
}
