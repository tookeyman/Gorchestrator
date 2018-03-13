package netcore

type packet int

const (
	Login packet = iota
	Disconnect
	Pong
	MsgAll
	BCI
	Localecho
	Tell
	Names
	Channels
)

var packets = [...]string{
	"LOGIN=Orchestrator;\tLOCALECHO 1\tNBMSGECHO 1\n",
	"\tPONG\n",
	"\tMSGALL\n",
	"\tBCI\n",
	"\tLOCALECHO ",
	"\tTELL\n",
	"\tNAMES\n",
	"\tCHANNELS\n",
}

func (pack packet) String() string {
	return packets[pack]
}
