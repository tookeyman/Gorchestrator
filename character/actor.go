package character

import (
	"fmt"
	"strconv"
	"strings"
)

type Actor struct {
	name, casting                                   string
	asyncChannel                                    chan string
	buffs, buffDur                                  []string
	hp, hpMax, mana, manaMax, end, endMax, id, zone int
	tID, tPctHP                                     int
	loc                                             *Location
}

type Location struct {
	x, y, z float64
	q       int
}

func GetActorInstance(name string, netbotsPacket string) *Actor {
	cha := Actor{
		name:         name,
		asyncChannel: make(chan string),
	}
	cha.UpdateActor(netbotsPacket)
	return &cha
}

func (cha *Actor) UpdateActor(packet string) {
	cha.processNetbotsPacket(packet)
	//@design: what else do we need to do on character update?
}

func (cha *Actor) processNetbotsPacket(packet string) {
	parts := strings.Split(packet, "|")
	for _, part := range parts {
		pack := strings.Split(part, "=")
		id := pack[0]
		var values = ""
		if len(pack) > 1 {
			values = pack[1]
		}
		cha.updatePart(id, values)
	}
}

func (cha *Actor) updatePart(id string, values string) {
	switch id {
	case "B":
		cha.buffs = strings.Split(values, ":")
		break
	case "F":
		break
	case "C":
		cha.casting = values
		break
	case "E":
		vals := strings.Split(values, "/")
		end, _ := strconv.Atoi(vals[0])
		endMax, _ := strconv.Atoi(vals[1])
		cha.end = end
		cha.endMax = endMax
		break
	case "X":
		//@research: what is this?
		break
	case "N":
		//@research: what is this?
		break
	case "L":
		//@research: what is this?
		break
	case "H":
		vals := strings.Split(values, "/")
		hp, _ := strconv.Atoi(vals[0])
		hpMax, _ := strconv.Atoi(vals[1])
		cha.hp = hp
		cha.hpMax = hpMax
		break
	case "M":
		vals := strings.Split(values, "/")
		mana, _ := strconv.Atoi(vals[0])
		manaMax, _ := strconv.Atoi(vals[1])
		cha.mana = mana
		cha.manaMax = manaMax
		break
	case "W":
		//@research: what is this?
		fmt.Println("W:", values)
		break
	case "P":
		//@research: what is this?
		//pet?
		break
	case "G":
		//@research: what is this?
		fmt.Println("G:", values)
		break
	case "S":
		//@research: what is this?
		fmt.Println("S:", values)
		break
	case "Y":
		//todo && @research: sit state
		//parsedLong, _ := strconv.ParseInt(values, 10, 64)
		//fmt.Println("Sit State: ", strconv.FormatInt(parsedLong, 16))
		//standingRest := 0x10
		//standingMove := 02000
		//if int(parsedLong) == standingRest{
		//	fmt.Println("Standing, rest")
		//}else if int(parsedLong) == (standingRest << 1){
		//	fmt.Println("Sitting, rest")
		//}else if int(parsedLong) == (standingRest >> 2){
		//	fmt.Println("Crouching, rest")
		//}
		break
	case "T":
		targetSplit := strings.Split(values, ":")
		tID, _ := strconv.Atoi(targetSplit[0])
		tPctHP, _ := strconv.Atoi(targetSplit[1])
		cha.tID = tID
		cha.tPctHP = tPctHP
		break
	case "Z":
		//@research find out what Z=[Zone]:???>[CharID] is missing
		vals := strings.Split(values, ":")
		zoneID, _ := strconv.Atoi(vals[0])
		id := -1
		idSplit := strings.Split(vals[1], ">")
		if len(idSplit) > 1 { //No Zone change
			tid, _ := strconv.Atoi(idSplit[1])
			id = tid
		}
		cha.id = id
		cha.zone = zoneID
		break
	case "D":
		cha.buffDur = strings.Split(values, ":")
		break
	case "@":
		locAndQ := strings.Split(values, ":")
		x, _ := strconv.ParseFloat(locAndQ[1], 64)
		y, _ := strconv.ParseFloat(locAndQ[0], 64)
		z, _ := strconv.ParseFloat(locAndQ[2], 64)
		q, _ := strconv.Atoi(locAndQ[3])
		newLoc := Location{x: x, y: y, z: z, q: q}
		cha.loc = &newLoc
		break
	case "$":
		//@research: what is this?
		break
	case "O":
		//@research: what is this?
		break
	case "A":
		//@research: what is this?
		break
	case "I":
		//@research: what is this?
		break
	case "R":
		//@research: what is this?
		break
	case "Q":
		//@research: what is this?
		break
	case "":
		//it's fucking stupid that we have to swallow this every time
		break
	default:
		fmt.Println("Did not understand:\t", id)
	}

}
