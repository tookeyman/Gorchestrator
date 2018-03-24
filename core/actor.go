package core

import (
	"fmt"
	"strconv"
	"strings"
)

type Actor struct {
	Name, casting                                   string
	asyncChannel                                    chan string
	buffs, buffDur, songs, memSpells, petBuffs      []string
	equipped, bagContents                           []string
	bagCapacity                                     []int
	hp, hpMax, mana, manaMax, end, endMax, id, zone int
	tID, tPctHP, petID, petPctHp, lvl, classID      int
	aaAssigned, aaSpent, aaAvailable                int
	loc                                             *Location
	heading                                         float64
	//unknown section
	f, x, n, o string
}

type Location struct {
	x, y, z float64
	q       int
}

func GetActorInstance(name string, netbotsPacket string) *Actor {
	cha := Actor{
		Name:         name,
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
	//fmt.Println(packet)
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
		cha.f = values
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
		cha.x = values
		break
	case "N":
		//todo @research: what is this?
		cha.n = values
		break
	case "L":
		//level:class
		lvlSplit := strings.Split(values, ":")
		lvl, _ := strconv.Atoi(lvlSplit[0])
		classID, _ := strconv.Atoi(lvlSplit[1])
		cha.lvl = lvl
		cha.classID = classID
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
		//pet buffs
		cha.petBuffs = strings.Split(values, ":")
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
		//memorized spells
		cha.memSpells = strings.Split(values, ":")
		break
	case "S":
		//songs
		cha.songs = strings.Split(values, ":")
		break
	case "Y":
		//movement/stace information i think this has info on snares and shit
		//todo && @research: sit state
		//parsedLong, _ := strconv.ParseInt(values, 10, 64)
		//fmt.Println("Sit State: ", strconv.FormatInt(parsedLong, 16))
		//standingRest := 0x10
		////standingMove := 02000
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
		//zone:???>myID
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
		cha.o = values
		break
	case "A":
		//AA: assigned:spent:available
		stringValues := strings.Split(values, ":")
		cha.aaAssigned, _ = strconv.Atoi(stringValues[0])
		cha.aaSpent, _ = strconv.Atoi(stringValues[1])
		cha.aaAvailable, _ = strconv.Atoi(stringValues[2])
		break
	case "I":
		//equipment items
		cha.equipped = strings.Split(values, ":")
		break
	case "R":
		//bag contents
		cha.bagContents = strings.Split(values, ":")
		break
	case "Q":
		//bag capacity
		stringArr := strings.Split(values, ":")
		bagArr := make([]int, len(stringArr))
		for i, item := range stringArr {
			capacity, _ := strconv.Atoi(item)
			bagArr[i] = capacity
		}
		cha.bagCapacity = bagArr
		break
	case "":
		//it's stupid that we have to swallow this every time
		break
	default:
		fmt.Println("Did not understand:\t", id)
	}
}
