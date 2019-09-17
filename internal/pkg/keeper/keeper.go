package keeper

import (
	"context"
	"database/sql"
	"strings"
	"strconv"

	_ "github.com/mattn/go-sqlite3"
	"github.com/sirupsen/logrus"

	"github.com/idlephysicist/cave-logger/internal/pkg/model"
)

type Keeper struct {
	log			*logrus.Logger
	db			*sql.DB
	ctx			context.Context
}

func New(log *logrus.Logger, dbFN string) *Keeper {
	var k Keeper

	k.log = log
	k.ctx = context.Background()

	db, err := sql.Open("sqlite3", dbFN)
	if err != nil {
		log.Errorf("keeper.newkeeper: Cannot establish database connection", err)
	}
	k.db = db
	return &k
}

func (k *Keeper) QueryLogs(arg string) ([]*model.Entry, error) {
	var query string // Go seems to complain if something is defined in if blocks
	if arg == `-1` {
		query = "SELECT id, date, cave, names, notes FROM `logs`"
	} else {
		query = "SELECT id, date, cave, names, notes FROM `logs` WHERE id = ?"
	}

	result, err := k.db.QueryContext(k.ctx, query, arg)
	if err != nil {
		k.log.Errorf("keeper.Query: Failed to query database", err)
		return nil, err
	}
	defer result.Close()

	err = result.Err()
	if err != nil {
		k.log.Errorf("keeper.Query: An error occurred querying database", err)
		return nil, err
	}

	rows := make([]*model.Entry, 0)
	for result.Next() {
		var caverIDstr string
		var row *model.Entry
		
		err = result.Scan(&row.ID, &row.Date, &row.Cave, &caverIDstr, &row.Notes)
		if err != nil {
			k.log.Error(err)
		}

		// Break up cavers into ids.
		brokenIDs := strings.Split(caverIDstr, `|`)
		for _, id_str := range brokenIDs {
			id, _ := strconv.Atoi(id_str)
			row.CaverIDs = append(row.CaverIDs, id)
		}
		
		row.Names = k.getCaverNames(caverIDstr)
		if row.Names == `` {
			continue
		}
		
		// Add this formatted row to the rows map
		rows = append(rows, row)  
	}

	return rows, err
}

func (k *Keeper) AddLog(params []string) error {
	return nil
}
// INTERNAL FUNCTIONS ----------------------------------------------------------

/*func (k *Keeper) fetchCaverIDs(caverList names string) ([]string, error) {
	var caverIDs []string
	for id, name := range caverList {
	
		fullNames := strings.Split(names, ",")
		for fullName, id := range fullNames {

			nameSplit := strings.Split(fullName, " ")
			if (nameSplit[0] == name.First && nameSplit[1] == name.Last) {
				caverIDs = append(caverIDs, id)
			}
		}
	}
	return caverIDs
}*/

func (k *Keeper) getCaverNames(idStr string) string {
	// Get the IDs
	caverList, err := k.getCaverList()
	if err != nil {
		k.log.Errorf("keeper.Query: Failed to fetch list of cavers")
		return ``
	}

	var names []string
	caverIDs := strings.Split(idStr, "|")

	for _, caver_id_str := range caverIDs {
		caver_id, _ := strconv.Atoi(caver_id_str) // Convert the string to an int...
	
		if keyExists(caverList, caver_id) {
			caver, _ := caverList[caver_id]

			fullName := caver.First + ` ` + caver.Last

			names = append(names, fullName)
		} else {
			continue
		}
	}

	
	return strings.Join(names, `, `)
}

func (k *Keeper) getCaverList() (map[int]*model.Caver, error) {
	result, err := k.db.QueryContext(k.ctx, "SELECT `id`,`first`,`last`,`club` FROM cavers")
	if err != nil {
		k.log.Errorf("keeper.newkeeper: Failed to get cavers", err)
	}
	cavers := make(map[int]*model.Caver) 
	for result.Next() {
		var id int
		var c model.Caver
		err = result.Scan(&id, &c.First, &c.Last, &c.Club)
		if err != nil {
			k.log.Errorf("Scan: %v", err)
		}
		//cavers = append(cavers, c)
		cavers[id] = &c
	}
	return cavers, err
}

func keyExists(array map[int]*model.Caver, key int) bool {
	_, grand := array[key]
	return grand
}
