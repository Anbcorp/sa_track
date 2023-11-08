package model

import (
	"errors"
	"fmt"
)

type BoatType int

const (
	B_CRUISER38 BoatType = iota + 1
	B_MINI               // 2
	B_ROSE               // 3
	B_CATA               // 4
	B_PERF50             // 5
	B_TRAINING           // 6
	B_FOLKBOAT           // 7
	B_RACER32            // 8
	B_KETCH              // 9
	B_IMOCA              // 10
)

var BoatTypes map[BoatType]string = make(map[BoatType]string)

func init() {
	BoatTypes[B_CRUISER38] = "Sailaway Cruiser 38"
	BoatTypes[B_MINI] = "Mini Transat"
	BoatTypes[B_ROSE] = "Caribbean Rose"
	BoatTypes[B_CATA] = "52' Cruising Cat"
	BoatTypes[B_PERF50] = "50' Performance Cruiser"
	BoatTypes[B_TRAINING] = "Nordic Folkboat (Training)"
	BoatTypes[B_FOLKBOAT] = "Nordic Folkboat"
	BoatTypes[B_RACER32] = "32' Offshore Racer"
	BoatTypes[B_KETCH] = "45' Ketch"
	BoatTypes[B_IMOCA] = "Imoca 60"
}

func BoattypeFromName(name string) (btype BoatType, err error) {
	for t, v := range BoatTypes {
		if v == name {
			btype = t
			err = nil
			return
		}
	}
	return -1, errors.New(fmt.Sprintf("BoattypeFromName - Unknown boat type : '%s'", name))
}
