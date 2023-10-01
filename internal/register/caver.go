package register

import (
	"database/sql"

	"github.com/idlephysicist/cave-logger/internal/model"
)

func (reg *Register) AddCaver(name, club, notes string) error {
	query := `INSERT INTO caver (name, club, notes) VALUES (?,?,?)`
	params := []any{name, club, notes}

	_, err := reg.executeTx(query, params)
	return errorWrapper("addcaver", err)
}

func (reg *Register) GetAllCavers() ([]*model.Caver, error) {
	query := `
    SELECT
        caver.id AS 'id',
        caver.name AS 'name',
        caver.club AS 'club',
        (
            SELECT COUNT(1)
            FROM trip_group
            WHERE trip_group.caverid = caver.id
        ),
        (
            SELECT printf("%s in %s", trip.date, cave.name)
            FROM trip, trip_group, cave
            WHERE trip.caveid == cave.id
              AND trip.id == trip_group.tripid
              AND trip_group.caverid == caver.id
            ORDER BY trip.date DESC, trip.id DESC LIMIT 1
        ) AS 'last_trip'
    FROM caver
    ORDER BY name`

	result, err := reg.db.Query(query)
	if err != nil {
		return nil, errorWrapper("getallcavers", err)
	}

	cavers := make([]*model.Caver, 0)
	for result.Next() {
		var (
			c        model.Caver
			lastTrip sql.NullString
		)

		err = result.Scan(&c.ID, &c.Name, &c.Club, &c.Count, &lastTrip)
		if err != nil {
			return nil, errorWrapper("getallcavers", err)
		}

		if lastTrip.Valid {
			c.LastTrip = lastTrip.String
		}

		cavers = append(cavers, &c)
	}
	if err = result.Err(); err != nil {
		return cavers, errorWrapper("getallcavers", err)
	}
	return cavers, nil
}

func (reg *Register) GetCaver(id string) (*model.Caver, error) {
	query := `
    SELECT
        caver.id AS 'id',
        caver.name AS 'name',
        caver.club AS 'club',
        (
            SELECT COUNT(1) FROM trip_group WHERE trip_group.caverid = caver.id
        ) AS 'count',
        (
            SELECT printf("%s in %s", trip.date, cave.name)
            FROM trip, trip_group, cave
            WHERE trip.caveid == cave.id
              AND trip.id == trip_group.tripid
              AND trip_group.caverid == caver.id
            ORDER BY trip.date DESC, trip.id DESC LIMIT 1
        ) AS 'last_trip',
        caver.notes AS 'notes'
    FROM caver
    WHERE caver.id = ?`

	result, err := reg.db.Query(query, id)
	if err != nil {
		return nil, errorWrapper("getcaver", err)
	}
	defer result.Close()

	var (
		c        model.Caver
		lastTrip sql.NullString
	)
	for result.Next() {

		err = result.Scan(&c.ID, &c.Name, &c.Club, &c.Count, &lastTrip, &c.Notes)
		if err != nil {
			return nil, errorWrapper("getcaver", err)
		}

		if lastTrip.Valid {
			c.LastTrip = lastTrip.String
		}

	}
	if err = result.Err(); err != nil {
		return nil, errorWrapper("getcaver", err)
	}

	return &c, nil
}

func (reg *Register) ModifyCaver(id, name, club, notes string) error {
	query := `UPDATE caver SET name = ?, club = ?, notes = ? WHERE id = ?`
	params := []any{name, club, notes, id}

	_, err := reg.executeTx(query, params)
	return errorWrapper("modifycaver", err)
}

func (reg *Register) RemoveCaver(id string) error {
	_, err := reg.executeTx(`DELETE FROM caver WHERE id = ?`, []any{id})
	return errorWrapper("removecaver", err)
}
