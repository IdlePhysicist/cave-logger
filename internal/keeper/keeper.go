package keeper

import (
	"context"
	"database/sql"
	"strings"
	"time"
	
	_ "github.com/mattn/go-sqlite3"
	"github.com/sirupsen/logrus"
)

/*type Caver struct {
	First string
	Last  string
	Club  string
}

type Row struct {
	Date 	*time.Time // REVIEW: This might not be correct examine the fmt returned from db
	Cave 	string
	Names []*Caver // REVIEW: Could I make this a list of pointers ?
	Notes string
}*/

type Keeper struct {
	log 			*logrus.Logger
	db 				*sql.DB
	ctx 			*context.Context
	CaverList [int]*Caver
}

func NewKeeper(log *logrus.Logger, dbFN string) *Keeper {
	var k Keeper

	k.log = log
	k.ctx = context.Background()

	db, err := sql.Open("sqlite3", dbFN)
	if err != nil {
		log.Errorf("keeper.newkeeper: Cannot establish database connection", err)
	}
	k.db = db
	
	k.CaverList = func () [int]*Caver {
		result, err := k.db.QueryContext(k.ctx, "SELECT `id`,`first`,`last`,`club` FROM cavers")
		if err != nil {
			k.log.Errorf("keeper.newkeeper: Failed to get cavers", err)
		}
		cavers := make(map[int]*Caver) 
		for result.Next() {
			var id int
			var c Caver
			err = rows.Scan(&id, &c.First, &c.Last, &c.Club)
			if err != nil {
				k.log.Errorf("Scan: %v", err)
			}
			//cavers = append(cavers, c)
			cavers[id] = c
		}
		return cavers	
	}

	return &k
}

/* 
REVIEW: IDEA: Perhaps as we come across the name ids we look up the ones on the 
row and add them to the CaverList as we go.

At first this approach might be more memory and cpu efficient. 
*/
func (k *Keeper) QueryLogs(query string, args interface{}) ([int]Row, error) {
	result, err := k.db.QueryContext(k.ctx, query, args)
	if err != nil {
		k.log.Errorf("keeper.Query: Failed to query database", err)
		return nil, err
	}
	defer result.close()

	err = result.Err()
	if err != nil {
		k.log.Errorf("keeper.Query: An error occurred querying database", err)
		return nil, err
	}

	rows = make(map[int]*Row) // NOTE: The int here is the db index
	for result.Next() {
		var id int
		var names string
		var row Row
		
		err = result.Scan(&id, &row.Date, &row.Cave, &names, &row.Notes)
		if err != nil {
			k.log.Error(err)
		}

		caverIDs := strings.Split(names, ",")
		for _, id := range caverIDs {
			if keyExists(k.CaverList, id) { 
				_, caver := k.CaverList[id]
			} else {
				k.log.Errorf("keeper.querylogs: id not in table of cavers", id)
				// TODO: Do something else here too; break?
			}
			row.Names = append(row.Names, caver) 
		}
		// Add this formatted row to the rows map
		rows[id] = row  
	}

	return rows, err
}

func (k *Keeper) InsertLog(args interface{}) error {
	query := "INSERT INTO logs ('date','cave','cavers','notes') VALUES (?,?,?,?,?)"
	// NOTE: The cavers field should be a cs string of indices for `cavers`
	caverIDs, err := k.fetchCaverIDs(args.Names)
	args.Names = strings.Join(caverIDs, ", ") // NOTE: The space

	result, err := k.db.ExecContext(k.ctx, query, args)
	if err != nil {
		k.log.Errorf("keeper.insert: Failed to exec statement", err)
		return err
	}

	rows, err := result.RowsAffected()
	if rows != 1 {
		k.log.WithFields(logrus.Fields{"rows":rows, "error":err}).Errorf("keeper.insert.rows")
	}
	return nil
}

// INTERNAL FUNCTIONS ----------------------------------------------------------

func (k *Keeper) fetchCaverIDs(names string) ([]string, error) {
	var caverIDs []string
	for id, name := range CaverList {
	
		fullNames := strings.Split(names, ",")
		for fullName, id := range fullNames {

			nameSplit := strings.Split(fullName, " ")
			if (nameSplit[0] == name.First && nameSplit[1] == name.Last) {
				caverIDs = append(caverIDs, id)
			}
		}
	}
	return caverIDs
}

func keyExists(array map[int]interface{}, key int) bool {
	_, grand := array[key]
	return grand
}