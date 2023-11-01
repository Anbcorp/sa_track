package nmea

import (
	"fmt"
	"reflect"
)

var MessageTypes map[string]reflect.Type = make(map[string]reflect.Type)

func init() {
	/* Depth is not available in APIBoatInfo.pl
	MessageTypes["DBT"] = reflect.TypeOf((*DBT)(nil))
	MessageTypes["DPT"] = reflect.TypeOf((*DPT)(nil)) */
	MessageTypes["GGA"] = reflect.TypeOf((*GGA)(nil))
	MessageTypes["GLL"] = reflect.TypeOf((*GLL)(nil))
	MessageTypes["HDT"] = reflect.TypeOf((*HDT)(nil))
	MessageTypes["MWV"] = reflect.TypeOf((*MWV)(nil))
	MessageTypes["RMC"] = reflect.TypeOf((*RMC)(nil))
	MessageTypes["VHW"] = reflect.TypeOf((*VHW)(nil))
	MessageTypes["VTG"] = reflect.TypeOf((*VTG)(nil))
}

// Interface for a sentence
type nmea_sentence interface {
	FromBoat(SABoat)
	Values() string
	Header() string
	SetOpt(byte)
}

// Common fields for a sentence
type nmea_data struct {
	id    string
	stype string
}

func (n *nmea_data) SetOpt(b byte) {}
func (n *nmea_data) Header() string {
	return fmt.Sprintf("%s%s", n.id, n.stype)
}
