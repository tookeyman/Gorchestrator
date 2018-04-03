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
	contents, err := ioutil.ReadFile(filePath)
	if err != nil {
		createDefaultSettings(filePath)
		return &settings{
			buffsAvailable: make([]string, 0),
			buffsNeeded:    make([]string, 0),
			isDefault:      true,
		}
	}
	return parseSettingsFile(contents)
}

func parseSettingsFile(data []byte) *settings {
	if checkFileForDefault(data) {
		return &settings{
			buffsAvailable: make([]string, 0),
			buffsNeeded:    make([]string, 0),
			isDefault:      true,
		}
	}
	return nil
}

func checkFileForDefault(contents []byte) bool {
	header := string(contents[:10])
	if header == "isDefault" {
		return true
	}
	return false
}

func createDefaultSettings(filePath string) {
	contents := "isDefault\t\t//Delete this line!!!\n[Buffs Available]\n//List buffs here. One per line. " +
		"Either exact name or spell ID number.\n[Buffs Needed]\n//List buffs here. One per line. " +
		"Either exact name or spell ID number."
	err := ioutil.WriteFile(filePath, []byte(contents), 0777)
	if err != nil {
		panic(err)
	}
}
