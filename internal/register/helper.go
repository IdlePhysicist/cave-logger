package register

import (
	"errors"
	"strings"
	"time"
)

// For retrieving the ID of a cave from the `locations` table.
func (reg *Register) getCaveID(cave string) (int, error) {
	result, err := reg.db.Query(`SELECT id FROM cave WHERE name == ?`, cave)
	if err != nil {
		return 0, err
	}
	defer result.Close()

	var caveID int
	for result.Next() {
		if err = result.Scan(&caveID); err != nil {
			return caveID, err
		}
	}

	// Check for errors from iterating over rows.
	if err := result.Err(); err != nil {
		return 0, err
	}
	return caveID, err
}

// For formatting the ids for a new
func (reg *Register) getCaverIDs(names string) ([]string, error) {
	var caverIDs []string

	cavers, err := reg.GetAllCavers()
	if err != nil {
		return caverIDs, err
	}

	namesList := strings.Split(names, ", ")

	for _, caver := range cavers {
		for _, name := range namesList {
			if strings.TrimSpace(name) == strings.TrimSpace(caver.Name) {
				caverIDs = append(caverIDs, caver.ID)
			}
		}
	}

	if len(caverIDs) != len(namesList) {
		return caverIDs, errors.New(`≥1 unknown cavers`)
	}

	return caverIDs, nil
}

func (reg *Register) verifyTrip(date, location, names, notes string) (
	[]interface{}, []string, error,
) {
	var peopleIDs []string

	// Conv the date to unix time
	_, err := time.Parse(dateFormat, date)
	if err != nil {
		return []interface{}{}, peopleIDs, err
	}

	locationID, err := reg.getCaveID(location)
	if err != nil {
		return []interface{}{}, peopleIDs, err
	} else if locationID == 0 {
		return []interface{}{}, peopleIDs, errors.New(`verifyTrip: Cave not known`)
	}

	peopleIDs, err = reg.getCaverIDs(names)
	if err != nil {
		return []interface{}{}, peopleIDs, err
	}

	return []interface{}{date, locationID, notes}, peopleIDs, nil
}
