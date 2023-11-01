package nmea

import (
	"fmt"

	"jason.go/util"
)

// MWV can be send for True wind or Aparent wind
type MWV struct {
	nmea_data
	wind_angle float64
	wind_speed float64
	wind_ref   byte
}

func (mwv *MWV) SetOpt(b byte) {
	mwv.wind_ref = b
}

/*
Wind angle is relative to the boat, from 0 to 359
*/
func (mwv *MWV) FromBoat(b SABoat) {
	mwv.id = "WI"
	mwv.stype = "MWV"
	switch mwv.wind_ref {
	default:
		mwv.wind_angle = util.Wind2Angle(b.Twa)
		mwv.wind_speed = b.Tws
		mwv.wind_ref = 'T'
	case 'R':
		mwv.wind_angle = util.Wind2Angle(b.Awa)
		mwv.wind_speed = b.Aws
		mwv.wind_ref = 'R'
	}
}

func (mwv *MWV) Values() string {
	// $WIMWV,319,T,13.7,N,A*05
	return fmt.Sprintf("%d,%c,%.1f,N,A", int(mwv.wind_angle), mwv.wind_ref, mwv.wind_speed)
}
