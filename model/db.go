package model

import (
	"bytes"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"reflect"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"jason.go/util"
)

const dbfile string = "sa_track.db"

var DbHandle *sql.DB

func OpenDB() {
	defer util.TimeMe(time.Now(), "OpenDB")
	var err error
	DbHandle, err = sql.Open("sqlite3", dbfile)
	if err != nil {
		log.Fatal(err)
	}
}

func newTx() *sql.Tx {
	tx, err := DbHandle.Begin()
	if err != nil {
		log.Fatal(err)
	}
	return tx
}

func execInsertQuery(query string, params ...any) error {
	tx := newTx()
	stmt, err := tx.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(params...)
	if err != nil {
		return err
	}
	err = tx.Commit()
	if err != nil {
		return err
	}
	return nil
}

func execSelectQuery(query string, params ...any) *sql.Rows {
	stmt, err := DbHandle.Prepare(query)
	if err != nil {
		log.Fatal(err)
	}
	rows, err := stmt.Query(params...)
	if err != nil {
		log.Fatal(err)
	}

	return rows
}

func execSelectQueryRow(query string, params ...any) *sql.Row {
	stmt, err := DbHandle.Prepare(query)
	if err != nil {
		log.Fatal(err)
	}
	row := stmt.QueryRow(params...)

	return row
}

func byte2string(b []byte) string {
	return bytes.NewBuffer(b).String()
}

// Scan rows into struct
func structScan(rows *sql.Rows, model interface{}) error {
	v := reflect.ValueOf(model)
	if v.Kind() != reflect.Ptr {
		return errors.New("must pass a pointer, not a value, to StructScan destination") // @todo add new error message
	}

	v = reflect.Indirect(v)
	t := v.Type()

	cols, _ := rows.Columns()

	var m map[string]interface{}
	for rows.Next() {
		columns := make([]interface{}, len(cols))
		columnPointers := make([]interface{}, len(cols))
		for i := range columns {
			columnPointers[i] = &columns[i]
		}

		if err := rows.Scan(columnPointers...); err != nil {
			return err
		}

		m = make(map[string]interface{})
		for i, colName := range cols {
			val := columnPointers[i].(*interface{})
			m[colName] = *val
		}

	}

	for i := 0; i < v.NumField(); i++ {
		field := strings.Split(t.Field(i).Tag.Get("json"), ",")[0]

		if item, ok := m[field]; ok {
			if v.Field(i).CanSet() {
				if item != nil {
					switch v.Field(i).Kind() {
					case reflect.String:
						v.Field(i).SetString(byte2string(item.([]uint8)))
					case reflect.Float32, reflect.Float64:
						v.Field(i).SetFloat(item.(float64))
					case reflect.Ptr:
						if reflect.ValueOf(item).Kind() == reflect.Bool {
							itemBool := item.(bool)
							v.Field(i).Set(reflect.ValueOf(&itemBool))
						}
					case reflect.Struct:
						v.Field(i).Set(reflect.ValueOf(item))
					default:
						fmt.Println(t.Field(i).Name, ": ", v.Field(i).Kind(), " - > - ", reflect.ValueOf(item).Kind()) // @todo remove after test out the Get methods
					}
				}
			}
		}
	}

	return nil
}

// Boat
func NewBoat(b *Boat) {
	err := execInsertQuery(q_NEWBOAT, b.Ubtnr, b.Name, b.Type)
	if err != nil {
		log.Fatal(fmt.Sprintf("model.NewBoat: %s", err))
	}
}

func GetBoats() []*Boat {
	rows := execSelectQuery(q_GETBOATS)
	defer rows.Close()
	var boats []*Boat
	for rows.Next() {
		boat := new(Boat)
		err := rows.Scan(boat.Ubtnr, boat.Name, boat.Type)
		if err != nil {
			log.Fatal(err)
		}
		boats = append(boats, boat)
	}
	return boats
}

func GetBoat(ubtnr int) (b *Boat, err error) {
	row := execSelectQueryRow(q_GETBOAT, ubtnr)
	boat := new(Boat)
	err = row.Scan(&boat.Ubtnr, &boat.Name, &boat.Type)
	if err == sql.ErrNoRows {
		return nil, err
	} else if err != nil {
		log.Fatal(fmt.Sprintf("model.GetBoat: %s", err))
	}

	return boat, nil
}

func FindBoat(name string) (b Boat, err error) {
	return Boat{}, nil
}

// BoatType
func NewBoatType(id int, name string) {
	err := execInsertQuery(q_NEWBOATTYPE, id, name)
	if err != nil {
		log.Fatal(fmt.Sprintf("model.NewBoatType: %s", err))
	}
}

// Check if all boat and sail types are set in the DB
func PopulateDB() {
	log.Println("PopulateDB: start")
	for boattype, typename := range BoatTypes {
		log.Printf("PopulateDB-Boats: %s...", typename)
		var bt int
		var btn string
		row := execSelectQueryRow("SELECT id, name FROM boat_type where id=?", boattype)
		err := row.Scan(&bt, &btn)
		if err == sql.ErrNoRows {
			log.Printf("inserting...")
			NewBoatType(int(boattype), typename)
		} else if err != nil {
			log.Fatal(fmt.Sprintf("model.PopulateDB: %s", err))
		}
		log.Println("OK!")
	}
	for sailtype, sailname := range SailTypes {
		log.Printf("PopulateDB-Sails: %s...", sailname)
		var st int
		var sn string
		row := execSelectQueryRow("SELECT id, name FROM sail where id=?", sailtype)
		err := row.Scan(&st, &sn)
		if err == sql.ErrNoRows {
			log.Printf("inserting...")
			NewSail(int(sailtype), sailname)
		} else if err != nil {
			log.Fatal(fmt.Sprintf("model.PopulateDB: %s", err))
		}
		log.Println("OK!")
	}
	log.Println("PopulateDB: finished!")
}

// Voyage
func NewVoyage()       {}
func GetVoyage(id int) {}

// MVP : voyage are only challenges/races. User defined voyage for later
//func FindVoyage(name string) {}

// Sail
func NewSail(id int, name string) {
	err := execInsertQuery(q_NEWSAILTYPE, id, name)
	if err != nil {
		log.Fatal(fmt.Sprintf("model.NewSailType: %s", err))
	}
}
func GetSails() {} // Load all sails at once, it's only for display

// BoatState
func NewState(b *Boat) {
	err := execInsertQuery(q_NEWSTATE,
		b.Timestamp, b.Ubtnr, b.Voyage.Id, b.Latitude, b.Longitude, b.Sog, b.Cog, b.Spd, b.Hdg, b.Awa, b.Aws,
		b.Twa, b.Tws, b.Divedegrees, b.Drift, b.Foilleft, b.Foilright, b.Heeldegrees, b.Keelangle,
		b.Waterballast, b.Weatherhelm, SailsToDB(b.ActiveSails))
	if err != nil {
		log.Fatal(fmt.Sprintf("model.NewState: %s", err))
	}
}
func BoatRefresh(boat *Boat) error {
	row := execSelectQueryRow(q_LASTBOATSTATE, boat.Ubtnr)
	b := new(Boat)
	var dbsails string
	var dbid int
	err := row.Scan(&dbid, &b.Ubtnr, &b.Latitude, &b.Longitude, &b.Sog, &b.Cog, &b.Spd, &b.Hdg, &b.Awa, &b.Aws, &b.Twa, &b.Tws, &dbsails)
	if err == sql.ErrNoRows {
		return err
	} else if err != nil {
		log.Fatal(fmt.Sprintf("model.GetBoat: %s", err))
	}
	if boat.Ubtnr != b.Ubtnr {
		// Db data is not corrupted somehow ?!
		log.Fatal(fmt.Sprintf("Boat id mismatch in db : %s(%s) has id %d, but got %s(%s) with id %d\n",
			b.Name, BoatTypes[b.Type], b.Ubtnr,
			boat.Name, BoatTypes[boat.Type], boat.Ubtnr))
	}
	boat.BoatState = b.BoatState
	return nil
}
