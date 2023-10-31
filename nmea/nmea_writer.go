package nmea

import (
	"fmt"
	"reflect"
	"strings"
)

// mokup
type SABoat struct {
	Heading float64
	Stw     float64
	Cog     float64
	Sog     float64
	Tws     float64
	Twd     float64
	Twa     float64
	Awa     float64
	Aws     float64
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

var MessageTypes map[string]reflect.Type = make(map[string]reflect.Type)

func init() {
	//MessageTypes["DBT"] = reflect.TypeOf((*DBT)(nil))
	//MessageTypes["DPT"] = reflect.TypeOf((*DPT)(nil))
	//MessageTypes["GGA"] = reflect.TypeOf((*GGA)(nil))
	//MessageTypes["GLL"] = reflect.TypeOf((*GLL)(nil))
	//MessageTypes["HDT"] = reflect.TypeOf((*HDT)(nil))
	MessageTypes["MWV"] = reflect.TypeOf((*MWV)(nil))
	//MessageTypes["RMC"] = reflect.TypeOf((*RMC)(nil))
	MessageTypes["VHW"] = reflect.TypeOf((*VHW)(nil))
	MessageTypes["VTG"] = reflect.TypeOf((*VTG)(nil))
}

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
		mwv.wind_angle = Wind2Angle(b.Twa)
		mwv.wind_speed = b.Tws
		mwv.wind_ref = 'T'
	case 'R':
		mwv.wind_angle = Wind2Angle(b.Awa)
		mwv.wind_speed = b.Aws
		mwv.wind_ref = 'R'
	}
}

// Convert an relative wind direction (TWA,AWA) (+/- 0->180) to an absolute wind angle (0->359)
func Wind2Angle(dir float64) float64 {
	if dir > 0 {
		return dir
	} else {
		return 360 + dir
	}
}

func (mwv *MWV) Values() string {
	// $WIMWV,319,T,13.7,N,A*05
	return fmt.Sprintf("%d,%c,%.1f,N,A", int(mwv.wind_angle), mwv.wind_ref, mwv.wind_speed)
}

// VHW is composed by a nmea_sentence, plus heading and stw
type VHW struct {
	nmea_data
	heading float64 // True and Mag heading are the same
	stw     float64 // TODO: which unit do we use internally ?
}

func (vhw *VHW) FromBoat(b SABoat) {
	vhw.id = "II"
	vhw.stype = "VHW"
	vhw.heading = b.Heading
	vhw.stw = b.Stw
}

func (vhw *VHW) Values() string {
	// $IIVHW,279,T,279,M,8.2,N,15,K*75
	return fmt.Sprintf("%d,T,%d,M,%.1f,N,%.1f,K", int(vhw.heading), int(vhw.heading), vhw.stw, Knt2Km(vhw.stw))
}

// VTG is VHW with ground for reference
type VTG struct {
	VHW
}

func (vtg *VTG) FromBoat(b SABoat) {
	vtg.id = "II"
	vtg.stype = "VTG"
	vtg.heading = b.Cog
	vtg.stw = b.Sog
}

func Knt2Km(s float64) float64 {
	return s * 1.852
}

func Ms2Knt(ms float64) float64 {
	return ms * 3.6 / 1.852
}

func Sentence(n nmea_sentence) string {
	message := fmt.Sprintf("%s,%s", n.Header(), n.Values())
	return fmt.Sprintf("$%s*%02X", message, Checksum(message))
}

func Checksum(s string) byte {
	var result byte
	for i := 0; i < len(s); i++ {
		result ^= s[i]
	}
	return result
}

func WriteMessage(b SABoat, seq []string) {
	// Instanciate and print Sentences
	for i := 0; i < len(seq); i++ {
		s := strings.Split(seq[i], ".")
		mtype := s[0]

		v := reflect.New(MessageTypes[mtype].Elem()).Interface()
		n := v.(nmea_sentence)

		// Set option for message, if any
		if len(s) > 1 {
			n.SetOpt(s[1][0])
		}

		n.FromBoat(b)

		fmt.Println(Sentence(n))
	}
}
