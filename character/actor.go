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
	tID, tPctHP, petID, petPctHp                    int
	loc                                             *Location
	heading                                         float64
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
		//buffs
		cha.buffs = strings.Split(values, ":")
		break
	case "F":
		//todo @research what is this?
		fmt.Println("F:", values)
		break
	case "C":
		//currently casting spell id
		cha.casting = values
		break
	case "E":
		//endurance
		vals := strings.Split(values, "/")
		end, _ := strconv.Atoi(vals[0])
		endMax, _ := strconv.Atoi(vals[1])
		cha.end = end
		cha.endMax = endMax
		break
	case "X":
		//todo @research: what is this?
		fmt.Println("X:", values)
		break
	case "N":
		//todo @research: what is this?
		fmt.Println("N:", values)
		break
	case "L":
		//todo @research: what is this?'
		fmt.Println("L:", values)
		break
	case "H":
		//health
		vals := strings.Split(values, "/")
		hp, _ := strconv.Atoi(vals[0])
		hpMax, _ := strconv.Atoi(vals[1])
		cha.hp = hp
		cha.hpMax = hpMax
		break
	case "M":
		//mana
		vals := strings.Split(values, "/")
		mana, _ := strconv.Atoi(vals[0])
		manaMax, _ := strconv.Atoi(vals[1])
		cha.mana = mana
		cha.manaMax = manaMax
		break
	case "W":
		//todo @research: what is this?
		fmt.Println("W:", values)
		break
	case "P":
		//petid:pcthp
		petSplit := strings.Split(values, ":")
		pID, _ := strconv.Atoi(petSplit[0])
		pPctHP, _ := strconv.Atoi(petSplit[1])
		cha.petID = pID
		cha.petPctHp = pPctHP
		break
	case "G":
		//todo memorized spells
		break
	case "S":
		//todo @research: what is this?
		//songs?
		fmt.Println("S:", values)
		break
	case "Y":
		//movement/stace information i think this has info on snares and shit
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
		//target info
		targetSplit := strings.Split(values, ":")
		tID, _ := strconv.Atoi(targetSplit[0])
		tPctHP, _ := strconv.Atoi(targetSplit[1])
		cha.tID = tID
		cha.tPctHP = tPctHP
		break
	case "Z":
		//zone/myID
		//todo @research find out what Z=[Zone]:???>[CharID] is missing
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
		//buff durations
		cha.buffDur = strings.Split(values, ":")
		break
	case "@":
		//location
		locAndQ := strings.Split(values, ":")
		x, _ := strconv.ParseFloat(locAndQ[1], 64)
		y, _ := strconv.ParseFloat(locAndQ[0], 64)
		z, _ := strconv.ParseFloat(locAndQ[2], 64)
		q, _ := strconv.Atoi(locAndQ[3])
		newLoc := Location{x: x, y: y, z: z, q: q}
		cha.loc = &newLoc
		break
	case "$":
		//heading
		heading, _ := strconv.ParseFloat(values, 64)
		cha.heading = heading
		break
	case "O":
		//todo @research: what is this?
		fmt.Println("O:", values)
		break
	case "A":
		//todo @research: what is this?
		fmt.Println("A:", values)
		break
	case "I":
		//todo equipment items
		break
	case "R":
		//todo bag contents
		break
	case "Q":
		//todo bag capacity
		break
	case "":
		//it's fucking stupid that we have to swallow this every time
		break
	default:
		fmt.Println("Did not understand:\t", id)
	}

}
