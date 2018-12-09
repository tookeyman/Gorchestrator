package core

import (
	"fmt"
	"math/rand"
	"sync"
)

type CharacterGraphManager struct {
	spellIdToName map[string]string
	rng           rand.Rand

	lock sync.Mutex
}

func CreateCharacterManager() CharacterGraphManager {
	man := CharacterGraphManager{}
	m := make(map[string]string)
	lock := sync.Mutex{}
	man.spellIdToName = m
	man.lock = lock
	return man
}

func (manager *CharacterGraphManager) CheckBuffs(characters *map[string]*Character) {
	if len(*characters) == 0 {
		return
	}
	//m := *characters
	buffers := buildBuffAvailableMap(characters)
	neededBuffs := buildNeededBuffsMap(characters)
	fmt.Printf("%#v\n%#v\n", buffers, neededBuffs)
}

func buildBuffAvailableMap(characters *map[string]*Character) map[string][]string {
	m := make(map[string][]string)
	for toon := range *characters {
		fmt.Println("Checking available buffs from ", toon, " ", (*characters)[toon].Props.buffsAvailable)
		for _, buff := range (*characters)[toon].Props.buffsAvailable {
			fmt.Println("Checking for ", buff)
			tList, ok := m[buff]
			if !ok {
				tList = []string{toon}
			} else {
				tList = append(tList, toon)
			}
			m[buff] = tList
		}
	}
	return m
}

func buildNeededBuffsMap(characters *map[string]*Character) map[string][]string {
	m := make(map[string][]string)
	for toon := range *characters {
		fmt.Println("Checking needed buffs from ", toon, " ", (*characters)[toon].Props.buffsNeeded)
		for _, buff := range (*characters)[toon].Props.buffsNeeded {
			fmt.Println("Checking for ", buff)
			tList, ok := m[buff]
			if !ok {
				tList = []string{toon}
			} else {
				tList = append(tList, toon)
			}
			m[buff] = tList
		}
	}
	return m
}
