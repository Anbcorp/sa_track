package model

import "jason.go/nmea"

const dbfile string = "sa_track.db"

// Boat
func NewBoat(id int, name string) {

}
func GetBoats()      {}
func GetBoat(id int) {}
func FindBoat(name string) nmea.SABoat {
	return nmea.SABoat{}
}

// BoatType
func NewBoatType() {}

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
