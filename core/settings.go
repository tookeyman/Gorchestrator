package core

import (
	"fmt"
	"io/ioutil"
)

type settings struct {
	buffsNeeded    []string
	buffsAvailable []string
	isDefault      bool
}

const (
	buffsAvailableHeader = "Buffs Available"
	buffsNeededHeader    = "Buffs Needed"
)

func ReadSettingsFile(name string) *settings {
	filePath := fmt.Sprintf("d:/MQ2/Orchestrator/%s.props", name)
	content, err := ioutil.ReadFile(filePath)
	blankSettings := &settings{
		buffsAvailable: make([]string, 0),
		buffsNeeded:    make([]string, 0),
		isDefault:      true,
	}
	if err != nil {
		createDefaultSettings(filePath)
		return blankSettings
	}
	if checkFileForDefault(content) {
		return blankSettings
	}
	return parseSettingsFile(content)
}

func parseSettingsFile(data []byte) *settings {

	asString := removeComments(string(data))
	cleanRawSettings := removeComments(asString)
	settingsMap := getSettingsMap(cleanRawSettings)
	fmt.Printf("SettingsMap:\n%#v\n", settingsMap)
	neededBuffs := settingsMap[buffsNeededHeader]
	availableBuffs := settingsMap[buffsAvailableHeader]
	return &settings{
		buffsAvailable: availableBuffs,
		buffsNeeded:    neededBuffs,
		isDefault:      asString[0:9] == "isDefault",
	}
}

func getSettingsMap(clean string) map[string][]string {
	fmt.Println(clean)
	m := make(map[string][]string)
	interestPointer := -1
	inHeader := false
	currentHeader := ""
	currentOptionsBuffer := make([]string, 100)
	currentOptionsLen := 0
	for idx, c := range clean {
		if c == ']' {
			if !inHeader {
				panic("Illegal ']' character outside of setting header")
			}
			currentHeader = string(clean[interestPointer+1 : idx])
			inHeader = false
		}
		if inHeader {
			continue
		}
		if c == '[' {
			interestPointer = idx
			inHeader = true
			if len(currentHeader) > 0 {
				m[currentHeader] = currentOptionsBuffer[0:currentOptionsLen]
				currentHeader = ""
				currentOptionsBuffer = make([]string, 100)
				currentOptionsLen = 0
			}
		}
		if c == '\n' {
			if clean[idx-1] == ']' {
				interestPointer = idx
				continue
			}
			line := string(clean[interestPointer+1 : idx])
			currentOptionsBuffer[currentOptionsLen] = line
			currentOptionsLen++
			fmt.Println("found ", line, "in ", currentHeader)
		}
	}
	if !inHeader {
		line := string(clean[interestPointer+1:])
		fmt.Println("Final line: ", line)
		currentOptionsBuffer[currentOptionsLen] = line
		m[currentHeader] = currentOptionsBuffer[0:currentOptionsLen]
	}
	return m
}

func removeComments(lines string) string {
	comment := false
	contents := ""
	for i := 0; i < len(lines)-1; i++ {
		char := byte(lines[i])
		nextChar := byte(lines[i+1])
		if char == byte('/') && nextChar == byte('/') {
			comment = true
		}

		if !comment {
			contents += string(char)
			if i == len(lines)-1 {
				contents += string(nextChar)
			}
		}
		if comment && char == byte('\n') {
			comment = false
		}
	}
	return contents
}

func checkFileForDefault(contents []byte) bool {
	header := string(contents[:10])
	if header == "isDefault" {
		return true
	}
	return false
}

func createDefaultSettings(filePath string) {
	contents := "isDefault\t//Delete this line!!!\n[" + buffsAvailableHeader + "]\n//List buffs here. One per line. " +
		"Either exact name or spell ID number.\n[" + buffsNeededHeader + "]\n//List buffs here. One per line. " +
		"Either exact name or spell ID number."
	err := ioutil.WriteFile(filePath, []byte(contents), 0777)
	if err != nil {
		panic(err)
	}
}
