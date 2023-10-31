package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"slices"
	"strings"

	"github.com/carlmjohnson/requests"
	"jason.go/nmea"
)

type SAUserBoatStatus struct {
	Boats []SABoatStatus `json:"boats"`
}

type SABoatStatus struct {
	Awa             float64        `json:"awa"`             //-128.044
	Aws             float64        `json:"aws"`             //4.8431
	Backstay        float64        `json:"backstay"`        //0
	Boatname        string         `json:"boatname"`        //"Jade Erre"
	Boattype        string         `json:"boattype"`        //"45' Ketch"
	Cog             float64        `json:"cog"`             //200.456
	Divedegrees     float64        `json:"divedegrees"`     //-0.0434
	Drift           float64        `json:"drift"`           //0.1067
	Foilleft        float64        `json:"foilleft"`        //0
	Foilright       float64        `json:"foilright"`       //1
	Hdg             float64        `json:"hdg"`             //197.402
	Heeldegrees     float64        `json:"heeldegrees"`     //-26.3374
	Keelangle       float64        `json:"keelangle"`       //-0.764466
	Latitude        float64        `json:"latitude"`        //11.9833226409507
	Longitude       float64        `json:"longitude"`       //80.5654012958121
	Misnr           uint64         `json:"misnr"`           //0
	Raceorchallenge string         `json:"raceorchallenge"` //""
	Sails           []SASailStatus `json:"sails"`           //"array of SASailStatus
	Sog             float64        `json:"sog"`             //2.7364
	Spd             float64        `json:"spd"`             //3.297
	Twa             float64        `json:"twa"`             //-146.333
	Twd             float64        `json:"twd"`             //51.0694
	Tws             float64        `json:"tws"`             //6.8802
	Ubtnr           string         `json:"ubtnr"`           //"115636"
	Voyage          string         `json:"voyage"`          //"Tj\u00f8rv\u00e5g, M\u00f8re og Romsdal, Norway -> Unknown"
	Waterballast    float64        `json:"waterballast"`    //-0.982885
	Weatherhelm     float64        `json:"weatherhelm"`     //-0.1063
}

type SASailStatus struct {
	Barberhaulerdown  float64 `json:"barberhaulerdown"`  //0
	Barberhaulerin    float64 `json:"barberhaulerin"`    //0
	Downhaul          float64 `json:"downhaul"`          //0
	Furled            float64 `json:"furled"`            //0
	Halyard           float64 `json:"halyard"`           //1
	Outhaul           float64 `json:"outhaul"`           //0.175
	Reef1             float64 `json:"reef_1"`            //0
	Reef2             float64 `json:"reef_2"`            //0
	Reef3             float64 `json:"reef_3"`            //0
	Sail              string  `json:"sail"`              //"Mainsail"
	Sheet             float64 `json:"sheet"`             //0.240727
	Travelerorleadcar float64 `json:"travelerorleadcar"` //0
	Vang              float64 `json:"vang"`              //0
}

func getFreshData() []SABoatStatus {
	var str []SABoatStatus
	err := requests.URL("http://srv.sailaway.world/cgi-bin/sailaway/APIBoatInfo.pl?usrnr=<redacted>&key=<redacted>").ToJSON(&str).Fetch(context.Background())
	if err != nil {
		fmt.Println("error: ", err)
	}
	return str
}

func getSavedData(filename string) []SABoatStatus {
	// Read the JSON file
	data, err := os.ReadFile(filename)
	if err != nil {
		log.Fatal("Error reading JSON file:", err)
	}

	res := []SABoatStatus{}
	json.Unmarshal(data, &res)
	return res
}

func main() {

	//res := getFreshData()
	res := getSavedData("json/20231030_185101.json")
	//fmt.Println(res)
	/*file, _ := os.OpenFile("output.json", os.O_CREATE, os.ModePerm)
	defer file.Close()
	enc := json.NewEncoder(file)
	enc.Encode(res)*/

	// Look for my boats
	myBoats := []string{"Volovan", "Jade Erre"}
	for _, boat := range res {
		if slices.Contains(myBoats, boat.Boatname) {
			var sailName string
			var raisedSails []string
			for _, sail := range boat.Sails {
				if sail.Furled < 1 && sail.Halyard > 0.9 {
					if slices.Contains([]string{"Mainsail", "Mizzen"}, sail.Sail) {
						// Look for reefs
						reefs := sail.Reef1 + sail.Reef2 + sail.Reef3
						sailName = fmt.Sprintf("%s-%d", sail.Sail, int(reefs))
					} else {
						sailName = sail.Sail
					}
					raisedSails = append(raisedSails, sailName)
				}

				/*print(f"stw:{boat.route.spd:0.1f} cog:{int(boat.route.cog)}")
				print(f"tws:{boat.wind.tws:0.1f} twd:{boat.wind.twd:0.1f} twa:{boat.wind.twa:.1f}")
				print(f"heel:{boat.attitude.heel:0.1f} keel:{boat.heel.keelangle} foils:{boat.heel.foils}")*/
			}
			fmt.Println(boat.Boatname)
			fmt.Println()
			fmt.Println(strings.Join(raisedSails, " / "))
			fmt.Printf("stw:%.1f cog:%d\n", boat.Spd*3.6/1.852, int(boat.Cog))
			fmt.Printf("tws:%.1f twd:%d twa:%d\n", boat.Tws*3.6/1.852, int(boat.Twd), int(boat.Twa))
			fmt.Printf("heel:%.1f keel:%.1f\n", boat.Heeldegrees, boat.Keelangle)
			fmt.Println("-----------------------------------------------------------------------")
		}
	}

	b := nmea.SABoat{
		Heading: 279,
		Stw:     8.2,
		Cog:     281,
		Sog:     8.2,
		Tws:     13.7,
		Twd:     238,
		Twa:     -41,
		Awa:     -26,
		Aws:     20.7,
	}
	list := []string{"VHW", "VTG", "MWV", "MWV.R"}
	nmea.WriteMessage(b, list)
}
