package model

type BoatType int

const (
	B_CRUISER38 BoatType = iota + 1
	B_MINI
	B_ROSE
	B_CATA
	B_PERF50
	B_TRAINING
	B_FOLKBOAT
	B_RACER32
	B_KETCH
	B_IMOCA
)

var BoatTypes map[BoatType]string = make(map[BoatType]string)

func init() {
	BoatTypes[B_CRUISER38] = "Sailaway Cruiser 38"
	BoatTypes[B_MINI] = "Mini Transat"
	BoatTypes[B_ROSE] = "52' Cruising Cat"
	BoatTypes[B_PERF50] = "50' Performance Cruiser"
	BoatTypes[B_TRAINING] = "Nordic Folkboat (Training)"
	BoatTypes[B_FOLKBOAT] = "Nordic Folkboat"
	BoatTypes[B_RACER32] = "32' Offshore Racer"
	BoatTypes[B_KETCH] = "45' Ketch"
	BoatTypes[B_IMOCA] = "Imoca 60"
}
