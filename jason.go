package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"slices"
	"strings"
	"time"

	"github.com/carlmjohnson/requests"
	"jason.go/model"
	"jason.go/nmea"
)

func getFreshData() []model.SABoatStatus {
	var str []model.SABoatStatus
	err := requests.URL("http://srv.sailaway.world/cgi-bin/sailaway/APIBoatInfo.pl?usrnr=<redacted>&key=<redacted>").ToJSON(&str).Fetch(context.Background())
	if err != nil {
		fmt.Println("error: ", err)
	}
	jsonname := fmt.Sprintf("json/%s.json", time.Now().Format("20060102_150405"))
	file, _ := os.OpenFile(jsonname, os.O_CREATE|os.O_TRUNC, os.ModePerm)
	defer file.Close()
	enc := json.NewEncoder(file)
	enc.Encode(str)
	return str
}

func getSavedData(filename string) []model.SABoatStatus {
	// Read the JSON file
	data, err := os.ReadFile(filename)
	if err != nil {
		log.Fatal("Error reading JSON file:", err)
	}

	res := []model.SABoatStatus{}
	json.Unmarshal(data, &res)
	return res
}

func printStatus() {
	//res := getFreshData()
	res := getSavedData("json/20231030_185101.json")
	//fmt.Println(res)

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

	b := model.Boat{
		BoatState: model.BoatState{
			Hdg: 279,
			Spd: 8.2,
			Cog: 281,
			Sog: 8.2,
			Tws: 13.7,
			//Twd:          238,
			Twa:       -41,
			Awa:       -26,
			Aws:       20.7,
			Latitude:  46.5086619986689,
			Longitude: -9.57924805707325,
			Timestamp: time.Now().UTC(),
		},
	}
	list := []string{"GLL", "GGA", "VHW", "HDT", "MWV", "MWV.R", "VTG", "RMC"}
	nmea.WriteMessage(b, list)
}

func dumpBoats() {
	res := getSavedData("json/20231030_185101.json")
	// Look for my boats
	myBoats := []string{"Volovan", "Jade Erre"}
	for _, boat := range res {
		if slices.Contains(myBoats, boat.Boatname) {
			var sailName string
			var raisedSails []string
			for _, sail := range boat.Sails {
				if sail.IsActive() {
					if sail.IsReefable() {
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
}

func testModel() {
	//res := getFreshData()
	res := getSavedData("json/20231108_125556.json")
	now := time.Now()
	for _, rb := range res {
		b := model.Json2model(rb, now)
		fmt.Println(b)
		list := []string{"GLL", "GGA", "VHW", "HDT", "MWV", "MWV.R", "VTG", "RMC"}
		nmea.WriteMessage(*b, list)
	}

}

func main() {
	testModel()
	//printStatus()
}
