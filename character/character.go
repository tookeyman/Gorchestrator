package character

import (
	"fmt"
	"strconv"
	"strings"
)

type Character struct {
	name, casting                                   string
	asyncChannel                                    chan string
	buffs, buffDur                                  []string
	hp, hpMax, mana, manaMax, end, endMax, id, zone int
	loc                                             *Location
}

type Location struct {
	x, y, z float64
	q       int64
}

func GetCharacterInstance(name string, netbotsPacket string) *Character {
	cha := Character{
		name:         name,
		asyncChannel: make(chan string),
	}
	cha.UpdateCharacter(netbotsPacket)
	return &cha
}

func (cha *Character) UpdateCharacter(packet string) {
	cha.processNetbotsPacket(packet)
	//@design: what else do we need to do on character update?
}

func (cha *Character) processNetbotsPacket(packet string) {
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

func (cha *Character) updatePart(id string, values string) {
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
		break
	case "P":
		//@research: what is this?
		break
	case "G":
		//@research: what is this?
		break
	case "S":
		//@research: what is this?
		break
	case "Y":
		//todo && @research:sit state
		break
	case "T":
		break
	case "Z":
		//@research find out what Z=<Zone>:???><CharID> is missing
		vals := strings.Split(values, ":")
		zoneLong, _ := strconv.ParseInt(vals[0], 10, 32)
		id := -1
		idSplit := strings.Split(vals[1], ">")
		if len(idSplit) > 1 { //No Zone change
			tid, _ := strconv.ParseInt(idSplit[1], 10, 32)
			id = int(tid)
		}
		cha.id = id
		cha.zone = int(zoneLong)
		break
	case "D":
		cha.buffDur = strings.Split(values, ":")
		break
	case "@":
		locAndQ := strings.Split(values, ":")
		x, _ := strconv.ParseFloat(locAndQ[1], 64)
		y, _ := strconv.ParseFloat(locAndQ[0], 64)
		z, _ := strconv.ParseFloat(locAndQ[2], 64)
		q, _ := strconv.ParseInt(locAndQ[3], 10, 32)
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
		//it's retarded that we have to swallow this every time
		break
	default:
		fmt.Println("Did not understand:\t", id)
	}

}
