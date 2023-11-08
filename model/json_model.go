package model

import (
	"fmt"
	"log"
	"slices"
	"strconv"
	"time"
)

type SAUserBoats struct {
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

func (s *SASailStatus) IsActive() bool {
	if s.Furled < 1 && s.Halyard >= 0.9 {
		return true
	}
	return false
}

func (s *SASailStatus) IsReefable() bool {
	if slices.Contains([]string{"Mainsail", "Mizzen"}, s.Sail) {
		return true
	}
	return false
}

// Prep model from json structure for internal use and DB save
func Json2model(b SABoatStatus, t time.Time) *Boat {
	d := new(Boat)

	d.Timestamp = t
	d.Ubtnr, _ = strconv.Atoi(b.Ubtnr)
	d.Name = b.Boatname
	// Convert from boat name to internal type
	btype, err := BoattypeFromName(b.Boattype)
	if err != nil {
		log.Fatal(err)
	}
	d.Type = btype
	// TODO: use reflect to fill this from a list of field names
	d.Latitude = b.Latitude
	d.Longitude = b.Longitude
	d.Cog = b.Cog
	d.Sog = b.Sog
	d.Spd = b.Spd
	d.Hdg = b.Hdg
	d.Awa = b.Awa
	d.Aws = b.Aws
	d.Twa = b.Twa
	d.Tws = b.Tws
	// Twd is not stored
	d.Divedegrees = b.Divedegrees
	d.Drift = b.Drift
	d.Foilleft = b.Foilleft
	d.Foilright = b.Foilright
	d.Heeldegrees = b.Heeldegrees
	d.Keelangle = b.Keelangle
	d.Waterballast = b.Waterballast
	d.Weatherhelm = b.Weatherhelm

	// Active sails
	for _, sail := range b.Sails {
		fmt.Printf("testing %s\n", sail.Sail)
		if sail.IsActive() {
			fmt.Printf("Active\n")
			dsail := new(Sail)

			sailn, err := SailtypeFromName(sail.Sail)
			if err != nil {
				log.Fatal(err)
			}
			dsail.Type = sailn

			if sail.IsReefable() {
				reefs := sail.Reef1 + sail.Reef2 + sail.Reef3
				dsail.Reefs = int(reefs) * 33
			} else {
				dsail.FurledPct = int(sail.Furled * 100)
			}
			d.ActiveSails = append(d.ActiveSails, *dsail)
		}
	}

	// Voyage, either misnr or voyage id. Voyage id is not known until we ask the DB
	if b.Misnr != 0 {
		d.Voyage.Id = int(b.Misnr)
		d.Voyage.Name = b.Raceorchallenge
	} else {
		//TODO: Can't fill voyage id and name for a user voyage unless we hit the db
	}

	return d
}
