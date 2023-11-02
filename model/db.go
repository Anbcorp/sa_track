package model

import (
	"bytes"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"reflect"
	"strings"

	_ "github.com/mattn/go-sqlite3"

	"jason.go/nmea"
)

const dbfile string = "sa_track.db"

var dbHandle *sql.DB

func OpenDB() {
	var err error
	dbHandle, err = sql.Open("sqlite3", dbfile)
	if err != nil {
		log.Fatal(err)
	}
}

func newTx() *sql.Tx {
	tx, err := dbHandle.Begin()
	if err != nil {
		log.Fatal(err)
	}
	return tx
}

func execInsertQuery(query string, params ...any) {
	tx := newTx()
	stmt, err := tx.Prepare(query)
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()
	_, err = stmt.Exec(params)
	if err != nil {
		log.Fatal(err)
	}
	err = tx.Commit()
	if err != nil {
		log.Fatal(err)
	}
}

func execSelectQuery(query string, params ...any) *sql.Rows {
	stmt, err := dbHandle.Prepare(query)
	if err != nil {
		log.Fatal(err)
	}
	rows, err := stmt.Query(params)
	if err != nil {
		log.Fatal(err)
	}

	return rows
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
func NewBoat(id int, name string, boat_type int) {
	execInsertQuery(q_NEWBOAT, id, name, boat_type)
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

func GetBoat(id int) {}
func FindBoat(name string) nmea.SABoat {
	return nmea.SABoat{}
}

// BoatType
func NewBoatType() {

}

// Voyage
func NewVoyage()       {}
func GetVoyage(id int) {}

// MVP : voyage are only challenges/races. User defined voyage for later
//func FindVoyage(name string) {}

// Sail
func NewSail()  {}
func GetSails() {} // Load all sails at once, it's only for display

// BoatState
func UpdateBoatSate(b nmea.SABoat) {}
func GetBoatLatestState(id int)    {}
