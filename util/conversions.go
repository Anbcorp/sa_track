package util

import (
	"fmt"
	"math"
	"time"
)

// knots to km/h
func Knt2Km(s float64) float64 {
	return s * 1.852
}

// m/s to knots
func Ms2Knt(ms float64) float64 {
	return ms * 3.6 / 1.852
}

// Convert an relative wind direction (TWA,AWA) (+/- 0->180) to an absolute wind angle (0->359)
func Wind2Angle(dir float64) float64 {
	if dir > 0 {
		return dir
	} else {
		return 360 + dir
	}
}

// Format a time.Time to NMEA GLL time
// HHmmss.ss
func NMEATimestamp(timestamp time.Time) string {
	// Golang way of formatting date is probably the worst I've ever seen, ugh...
	return timestamp.Format("150405.00")
}

// Format a time.Time to NMEA RMC date
// DDMMYY
func NMEADate(timestamp time.Time) string {
	return timestamp.Format("020106")
}

// Format a position in DDDMM.mmmmm format
func NMEALatLon(pos float64) string {
	dd, mins := math.Modf(math.Abs(pos))
	mm := mins * 60
	return fmt.Sprintf("%02d%02.04f", int(dd), mm)
}
