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
	neededBuffs := make([]string, 0)
	availableBuffs := make([]string, 0)
	asString := removeComments(string(data))

	fmt.Printf("Parsed Settings File: %#v\n", removeComments(asString))
	return &settings{
		buffsAvailable: availableBuffs,
		buffsNeeded:    neededBuffs,
		isDefault:      asString[0:9] == "isDefault",
	}
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
	contents := "isDefault\t//Delete this line!!!\n[Buffs Available]\n//List buffs here. One per line. " +
		"Either exact name or spell ID number.\n[Buffs Needed]\n//List buffs here. One per line. " +
		"Either exact name or spell ID number."
	err := ioutil.WriteFile(filePath, []byte(contents), 0777)
	if err != nil {
		panic(err)
	}
}
