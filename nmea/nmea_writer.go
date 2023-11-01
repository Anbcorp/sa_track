package nmea

import (
	"fmt"
	"reflect"
	"strings"
	"time"
)

// mokup
type SABoat struct {
	Heading      float64
	Stw          float64
	Cog          float64
	Sog          float64
	Tws          float64
	Twd          float64
	Twa          float64
	Awa          float64
	Aws          float64
	Latitude     float64
	Longitude    float64
	PositionDate time.Time
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
	// Always mark data as from a virtual origin
	fmt.Println("$SOL")

	// Instanciate and print Sentences
	for i := 0; i < len(seq); i++ {
		s := strings.Split(seq[i], ".")

		// Ignore unsupported types
		mtype, ok := MessageTypes[s[0]]
		if ok {

			v := reflect.New(mtype.Elem()).Interface()
			n := v.(nmea_sentence)

			// Set option for message, if any
			if len(s) > 1 {
				n.SetOpt(s[1][0])
			}

			n.FromBoat(b)

			fmt.Println(Sentence(n))
		}
	}
}
