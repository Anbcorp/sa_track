package nmea

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
