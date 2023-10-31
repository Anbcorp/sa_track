package nmea

import (
	"fmt"
	"reflect"
)

// mokup
type SABoat struct {
	Heading uint16
	Stw     float64
	Cog     uint16
	Sog     float64
}

// Interface for a sentence
type nmea_sentence interface {
	FromBoat(SABoat)
	Values() string
}

// Common fields for a sentence
type nmea_data struct {
	id       string
	stype    string
	checksum uint8
}

var MessageTypes map[string]reflect.Type = make(map[string]reflect.Type)

func init() {
	MessageTypes["VHW"] = reflect.TypeOf((*VHW)(nil))
	MessageTypes["VTG"] = reflect.TypeOf((*VTG)(nil))
}

// VHW is composed by a nmea_sentence, plus heading and stw
type VHW struct {
	nmea_data
	heading uint16  // True and Mag heading are the same
	stw     float64 // TODO: which unit do we use internally ?
}

func (vhw *VHW) Values() string {
	// $IIVHW,279,T,279,M,8.2,N,15,K*75
	return fmt.Sprintf("%s%s,%d,T,%d,M,%.1f,N,%.1f,K", vhw.id, vhw.stype, vhw.heading, vhw.heading, vhw.stw, Knt2Km(vhw.stw))
}

func (vhw *VHW) FromBoat(b SABoat) {
	vhw.id = "II"
	vhw.stype = "VHW"
	vhw.heading = b.Heading
	vhw.stw = b.Stw
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
	v := n.Values()
	return fmt.Sprintf("$%s*%x", v, Checksum(v))
}

func Checksum(s string) byte {
	var result byte
	for i := 0; i < len(s); i++ {
		result ^= s[i]
	}
	return result
}

func WriteMessage(b SABoat, seq []string) {

	/*
		for type in seq
			m new(type)
			m.FromBoat(b)
			append(Sentence(m))
	*/
	n := new(VHW)
	fmt.Println(reflect.TypeOf(n))
	n.FromBoat(b)

	for i := 0; i < len(seq); i++ {
		v := reflect.New(MessageTypes[seq[i]].Elem()).Interface()
		fmt.Println(reflect.TypeOf(v))
		n := v.(nmea_sentence)
		fmt.Println(reflect.TypeOf(n))
		n.FromBoat(b)
		fmt.Println(Sentence(n))
	}
	/*n := new(VHW)
	n.FromBoat(b)
	fmt.Println(reflect.TypeOf(n))
	fmt.Println(Sentence(n))

	m := new(VTG)
	m.FromBoat(b)
	fmt.Println(reflect.TypeOf(m))
	fmt.Println(Sentence(m))*/
}
