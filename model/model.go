package model

import (
	"math"
	"time"
)

type Boat struct {
	Ubtnr int
	Name  string
	Type  BoatType
	BoatState
}

type BoatState struct {
	Timestamp    time.Time
	Latitude     float64
	Longitude    float64
	Cog          float64
	Sog          float64
	Stw          float64
	Hdg          float64
	Awa          float64
	Aws          float64
	Twa          float64
	Tws          float64
	Divedegrees  float64
	Drift        float64
	Foilleft     float64
	Foilright    float64
	Heeldegrees  float64
	Keelangle    float64
	Waterballast float64
	Weatherhelm  float64
	ActiveSails  []Sail
	Voyage
}

// DB rep: +/- sailtype.reduction
// positive is reefed
// negative is furled
type Sail struct {
	Type      SailType
	Reefs     int
	FurledPct int
}

func (s *Sail) ToDB() float64 {
	reduction := 0
	if s.FurledPct > 0 || s.FurledPct < 1 {
		reduction = s.FurledPct
	} else if s.Reefs > 0 {
		reduction = s.Reefs
	}
	return float64(s.Type) + float64(reduction)/100
}

func SailFromDB(s float64) *Sail {
	sail := new(Sail)
	stype, reduction := math.Modf(math.Abs(s))
	sail.Type = SailType(stype)
	if s >= 0 {
		sail.Reefs = int(reduction * 100)
	} else {
		sail.FurledPct = int(reduction * 100)
	}
	return sail
}

type Voyage struct {
	Id   int
	Name string
}
