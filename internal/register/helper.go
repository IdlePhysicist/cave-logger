package register

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

// For retrieving the ID of a cave from the `locations` table.
func (reg *Register) getCaveID(cave string) (int, error) {
	result, err := reg.db.Query(`SELECT id FROM cave WHERE name == ?`, cave)
	if err != nil {
		return 0, errorWrapper("getcaveid", err)
	}
	defer result.Close()

	var caveID int
	for result.Next() {
		if err = result.Scan(&caveID); err != nil {
			return caveID, errorWrapper("getcaveid", err)
		}
	}

	// Check for errors from iterating over rows.
	if err := result.Err(); err != nil {
		return 0, errorWrapper("getcaveid", err)
	}
	return caveID, nil
}

// For formatting the ids for a new
func (reg *Register) getCaverIDs(names string) ([]string, error) {
	var caverIDs []string

	cavers, err := reg.GetAllCavers()
	if err != nil {
		return caverIDs, errorWrapper("getcaverids", err)
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
		return caverIDs, errorWrapper("getcaverids", errors.New(`≥1 unknown cavers`))
	}

	return caverIDs, nil
}

func (reg *Register) verifyTrip(date, location, names, notes string) (
	[]any, []string, error,
) {
	var peopleIDs []string

	// Conv the date to unix time REVIEW: do I need this now?
	_, err := time.Parse(dateFormat, date)
	if err != nil {
		return []any{}, peopleIDs, err
	}

	locationID, err := reg.getCaveID(location)
	if err != nil {
		return []any{}, peopleIDs, fmt.Errorf("verifyTrip: Cave not known - %w", err)
	}

	peopleIDs, err = reg.getCaverIDs(names)
	if err != nil {
		return []any{}, peopleIDs, err
	}

	return []any{date, locationID, notes}, peopleIDs, nil
}
