package register

import (
	"database/sql"

	"github.com/idlephysicist/cave-logger/internal/model"
)

func (reg *Register) AddCaver(name, club, notes string) error {
	query := `INSERT INTO caver (name, club, notes) VALUES (?,?,?)`
	params := []any{name, club, notes}

	_, err := reg.executeTx(query, params)
	return err
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
		// reg.log.Errorf("reg.getcaverlist: Failed to get cavers", err)
		return nil, err
	}

	cavers := make([]*model.Caver, 0)
	for result.Next() {
		var (
			c        model.Caver
			lastTrip sql.NullString
		)

		err = result.Scan(&c.ID, &c.Name, &c.Club, &c.Count, &lastTrip)
		if err != nil {
			reg.log.Errorf("Scan: %v", err)
		}

		if lastTrip.Valid {
			c.LastTrip = lastTrip.String
		}

		cavers = append(cavers, &c)
	}
	if err = result.Err(); err != nil {
		reg.log.Errorf("reg.get: Step error: %s", err)
		return cavers, err
	}
	return cavers, err
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
		// reg.qg.Errorf("reg.prepare: Failed to query database", err)
		return nil, err
	}
	defer result.Close()

	var (
		caver    model.Caver
		lastTrip sql.NullString
	)
	for result.Next() {

		err = result.Scan(&caver.ID, &caver.Name, &caver.Club, &caver.Count, &lastTrip, &caver.Notes)
		if err != nil {
			// reg.log.Errorf("reg.scan", err)
			return nil, err
		}

		if lastTrip.Valid {
			caver.LastTrip = lastTrip.String
		}

	}
	if err = result.Err(); err != nil {
		reg.log.Errorf("reg.get: Step error: %s", err)
		return nil, err
	}

	return &caver, err
}

func (reg *Register) ModifyCaver(id, name, club, notes string) error {
	query := `UPDATE caver SET name = ?, club = ?, notes = ? WHERE id = ?`
	params := []any{name, club, notes, id}

	_, err := reg.executeTx(query, params)
	return err
}

func (reg *Register) RemoveCaver(id string) error {
	_, err := reg.executeTx(`DELETE FROM caver WHERE id = ?`, []any{id})
	return err
}
