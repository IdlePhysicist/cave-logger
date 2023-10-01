package register

import (
	"database/sql"

	"github.com/idlephysicist/cave-logger/internal/model"
)

func (reg *Register) AddCave(name, region, country, notes string, srt bool) error {
	query := `INSERT INTO cave (name, region, country, srt, notes) VALUES (?,?,?,?,?)`
	params := []any{name, region, country, srt, notes}

	_, err := reg.executeTx(query, params)
	return errorWrapper("addcave", err)
}

func (reg *Register) GetAllCaves() ([]*model.Cave, error) {
	query := `
    SELECT
        cave.id AS 'id',
        cave.name AS 'name',
        cave.region AS 'region',
        cave.country AS 'country',
        cave.srt AS 'srt',
        (
            SELECT COUNT(1)
            FROM trip
            WHERE trip.caveid = cave.id
        ) AS 'visits',
        (
            SELECT trip.date
            FROM trip
            WHERE trip.caveid = cave.id
            ORDER BY trip.date DESC LIMIT 1
        ) AS 'last_visit'
    FROM cave
    ORDER BY name`
	result, err := reg.db.Query(query)
	if err != nil {
		return nil, errorWrapper("getallcaves", err)
	}

	caves := make([]*model.Cave, 0)
	for result.Next() {
		var (
			c         model.Cave
			lastVisit sql.NullString
		)

		err = result.Scan(&c.ID, &c.Name, &c.Region, &c.Country, &c.SRT, &c.Visits, &lastVisit)
		if err != nil {
			return nil, errorWrapper("getallcaves", err)
		}

		if lastVisit.Valid {
			c.LastVisit = lastVisit.String
		}

		caves = append(caves, &c)
	}
	if err = result.Err(); err != nil {
		return caves, errorWrapper("getallcaves", err)
	}
	return caves, nil
}

func (reg *Register) GetCave(id string) (*model.Cave, error) {
	query := `
    SELECT
        cave.id AS 'id',
        cave.name AS 'name',
        cave.region AS 'region',
        cave.country AS 'country',
        cave.srt AS 'srt',
        (
            SELECT COUNT(1) FROM trip WHERE trip.caveid = cave.id
        ) AS 'visits',
        (
            SELECT trip.date
            FROM trip
            WHERE trip.caveid = cave.id
            ORDER BY trip.date DESC LIMIT 1
        ) AS 'last_visit',
        cave.notes AS 'notes'
    FROM cave
    WHERE id = ?`

	result, err := reg.db.Query(query, id)
	if err != nil {
		return nil, errorWrapper("getcave", err)
	}
	defer result.Close()

	var (
		cave      model.Cave
		lastVisit sql.NullString
	)
	for result.Next() {
		err = result.Scan(&cave.ID, &cave.Name, &cave.Region, &cave.Country, &cave.SRT, &cave.Visits, &lastVisit, &cave.Notes)
		if err != nil {
			return nil, errorWrapper("getcave", err)
		}

		if lastVisit.Valid {
			cave.LastVisit = lastVisit.String
		}
	}
	if err = result.Err(); err != nil {
		return nil, errorWrapper("getcave", err)
	}

	return &cave, nil
}

func (reg *Register) RemoveCave(id string) error {
	_, err := reg.executeTx(`DELETE FROM cave WHERE id = ?`, []any{id})
	return errorWrapper("removecave", err)
}

//
// MODIFY FUNCS ---- ----

func (reg *Register) ModifyCave(id, name, region, country, notes string, srt bool) error {
	query := `UPDATE cave SET name = ?, region = ?, country = ?, srt = ?, notes = ? WHERE id = ?`
	params := []any{name, region, country, srt, notes, id}

	_, err := reg.executeTx(query, params)
	return errorWrapper("modifycave", err)
}
