package nmea

import (
	"fmt"
)

// HDT - Heading True - Actual vessel heading in degrees true
type HDT struct {
	nmea_data
	heading float64
}

func (hdt *HDT) FromBoat(b SABoat) {
	hdt.id = "II"
	hdt.stype = "HDT"
	hdt.heading = b.Heading
}

func (hdt *HDT) Values() string {
	// $IIHDT,265,T*3D
	return fmt.Sprintf("%d,T", int(hdt.heading))
}
