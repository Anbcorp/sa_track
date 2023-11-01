package model

type SailType int

const (
	S_MAIN SailType = iota + 1
	S_MIZZEN
	S_GENAKER
	S_GENOA
	S_JIB
	S_CODE0
	S_MIZZSTAY
	S_N1
	S_N2
	S_N3
	S_N4
	S_STORM
)

var SailTypes map[SailType]string = make(map[SailType]string)

func init() {
	SailTypes[S_MAIN] = "Mainsail"
	SailTypes[S_MIZZEN] = "Mizzen"
	SailTypes[S_GENAKER] = "Genaker"
	SailTypes[S_GENOA] = "Genoa"
	SailTypes[S_JIB] = "Jib"
	SailTypes[S_CODE0] = "Code 0"
	SailTypes[S_MIZZSTAY] = "Mizzen staysail"
	SailTypes[S_N1] = "Nr.1"
	SailTypes[S_N2] = "Nr.2"
	SailTypes[S_N3] = "Nr.3"
	SailTypes[S_N4] = "Nr.4"
	SailTypes[S_STORM] = "Stormjib"
}
