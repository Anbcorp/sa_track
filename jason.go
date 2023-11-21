package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"slices"
	"strings"
	"syscall"
	"time"

	"github.com/carlmjohnson/requests"
	"gopkg.in/yaml.v3"
	"jason.go/model"
	"jason.go/nmea"
	"jason.go/util"
)

type Config struct {
	Usrnr  int64  `yaml:"usrnr"`
	Apikey string `yaml:"apikey"`
}

var config Config

const SAAPI_URL string = "http://srv.sailaway.world/cgi-bin/sailaway/APIBoatInfo.pl"

func getFreshData() []model.SABoatStatus {
	var str []model.SABoatStatus
	url := fmt.Sprintf("%s?usrnr=%d&key=%s", SAAPI_URL, config.Usrnr, config.Apikey)
	err := requests.URL(url).ToJSON(&str).Fetch(context.Background())
	if err != nil {
		fmt.Println("error: ", err)
	}
	jsonname := fmt.Sprintf("json/%s.json", time.Now().Format("20060102_150405"))
	file, _ := os.OpenFile(jsonname, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, os.ModePerm)
	defer func() {
		file.Sync()
		file.Close()
	}()
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
	defer util.TimeMe(time.Now(), "testModel")
	//res := getFreshData()
	res := getSavedData("json/20231108_135120.json")
	now := time.Now()
	var boats []*model.Boat
	for _, rb := range res {
		b := model.Json2model(rb, now)
		//fmt.Println(rb.Boatname, rb.Boattype)
		//list := []string{"GLL", "GGA", "VHW", "HDT", "MWV", "MWV.R", "VTG", "RMC"}
		//nmea.WriteMessage(*b, list)
		boats = append(boats, b)
	}

	model.OpenDB()
	fmt.Println("here")
	for _, b := range boats {
		_, err := model.GetBoat(b.Ubtnr)
		if err != nil {
			fmt.Printf("Boat '%d: %s' not found, inserting...\n", b.Ubtnr, b.Name)
			model.NewBoat(b)
		} else {
			fmt.Printf("Boat %d: %s found!\n", b.Ubtnr, b.Name)
			if b.Spd > 0 {
				model.NewState(b)
			}
		}
	}
}

func getConfig() {
	confdata, err := os.ReadFile("sa_track.yaml")
	if err != nil {
		log.Fatal("Error reading config file:", err)
	}
	yaml.Unmarshal(confdata, &config)
}

func refreshBoats(t time.Time, boatnames []string) {
	defer util.TimeMe(t, "refreshBoats")
	log.Print("refreshing ", t.Format("2006/01/02 15:04:05"))
	//res := getSavedData("json/20231110_205644.json")
	res := getFreshData()
	var boats []*model.Boat
	for _, rb := range res {
		if true { //&& slices.Contains(boatnames, rb.Boatname) {
			b := model.Json2model(rb, t)
			boats = append(boats, b)
		}
	}
	for _, b := range boats {
		oldboat, err := model.GetBoat(b.Ubtnr)
		if err != nil {
			//fmt.Printf("Boat '%d: %s' not found, inserting...\n", b.Ubtnr, b.Name)
			model.NewBoat(b)
		}

		model.BoatRefresh(oldboat)
		if b.Latitude == oldboat.Latitude && b.Longitude == oldboat.Longitude {
			// Boat didn't move, skip
			log.Printf("%s didn't move, skipping...\n", b.Name)
		} else {
			log.Printf("%s at %.5f,%.5f heading %d at %.1f knt\n", b.Name, b.Latitude, b.Longitude, int(b.Hdg), b.Spd)
			model.NewState(b)
		}
	}

}

func PrintBoats(boats []string) {
	fmt.Println(boats)
}

func main() {
	getConfig()
	//fmt.Println(config)
	//getFreshData()
	//model.PopulateDB()
	//testModel()
	//printStatus()

	model.OpenDB()
	defer model.DbHandle.Close()

	done := make(chan bool)
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)

	go func() {
		for {
			select {
			case s := <-signalChan:
				switch s {
				case syscall.SIGINT, syscall.SIGTERM:
					log.Printf("Bye... o7")
					model.DbHandle.Close()
					done <- true
					os.Exit(1)
				case syscall.SIGHUP:
					PrintBoats([]string{"Volovan", "Jade erre", "Challenge Accepted"})
				}
			}
		}
	}()

	defer signal.Stop(signalChan)
	ticker := time.NewTicker(10 * time.Minute)
	refreshBoats(time.Now(), []string{"Volovan", "Jade Erre", "Challenge Accepted"})
	for {
		select {
		case <-done:
			ticker.Stop()
			break
		case t := <-ticker.C:
			refreshBoats(t, []string{"Volovan", "Jade Erre", "Challenge Accepted"})
			log.Print("\n==========> Waiting for next tick ... <==========\n")
		}
	}

}
