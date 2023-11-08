package nmea

import (
	"fmt"
	"reflect"
	"strings"

	"jason.go/model"
)

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

func WriteMessage(b model.Boat, seq []string) {
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
