package main

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"

	"golang.org/x/sys/windows/registry"
)

// RegValueInfo - Holds information used to describe a field/value pair within a registry key
type RegValueInfo struct {
	Key     string
	Buf     []byte
	ValType int
}

func getKey() registry.Key {
	iconPositions, err := registry.OpenKey(registry.CURRENT_USER, "SOFTWARE\\Microsoft\\Windows\\Shell\\Bags\\1\\Desktop", registry.QUERY_VALUE)

	if err != nil {
		panic(err)
	}

	return iconPositions
}

func export(iconPositions registry.Key) {

	KeyValPairs := []RegValueInfo{}

	valueNames, _ := iconPositions.ReadValueNames(-1)

	for i := 0; i < len(valueNames); i++ {
		curKey := valueNames[i]
		buf, valType, err := iconPositions.GetBinaryValue(curKey)

		if err == nil {
			KeyValPairs = append(KeyValPairs, RegValueInfo{
				curKey,
				buf,
				registry.BINARY,
			})
			continue
		}

		tempBuffer := new(bytes.Buffer)
		switch valType {
		case registry.DWORD:
			fallthrough
		case registry.QWORD:

			intVal, _, _ := iconPositions.GetIntegerValue(curKey)
			binary.Write(tempBuffer, binary.LittleEndian, intVal)
		case registry.SZ:
			fallthrough
		case registry.EXPAND_SZ:
			stringVal, _, _ := iconPositions.GetStringValue(curKey)
			binary.Write(tempBuffer, binary.LittleEndian, stringVal)
		default:
			panic(valType)
		}

		KeyValPairs = append(KeyValPairs, RegValueInfo{
			curKey,
			tempBuffer.Bytes(),
			int(valType),
		})
	}

	fileContents, jsonErr := json.Marshal(KeyValPairs)

	if jsonErr != nil {
		panic(jsonErr)
	}

	fileErr := ioutil.WriteFile("regKeyData", fileContents, 0777)

	if fileErr != nil {
		panic(fileErr)
	}
}

func clear(iconPositions registry.Key) {
	valueNames, _ := iconPositions.ReadValueNames(-1)

	for i := 0; i < len(valueNames); i++ {
		iconPositions.DeleteValue(valueNames[i])
	}
}

func importKey(iconPositions registry.Key) {
	var KeyValPairs []RegValueInfo

	fileData, _ := ioutil.ReadFile("regKeyData")

	json.Unmarshal(fileData, &KeyValPairs)

	for i := 0; i < len(KeyValPairs); i++ {
		curData := KeyValPairs[i]
		reader := bytes.NewReader(curData.Buf)
		switch curData.ValType {
		case registry.DWORD:
			var value uint32
			binary.Read(reader, binary.LittleEndian, &value)
			iconPositions.SetDWordValue(curData.Key, value)
		case registry.BINARY:
			var binData []byte
			binary.Read(reader, binary.LittleEndian, &binData)
			iconPositions.SetBinaryValue(curData.Key, binData)
		case registry.SZ:
			var strData string
			binary.Read(reader, binary.LittleEndian, &strData)
			iconPositions.SetStringValue(curData.Key, strData)
		}
	}
}

func restartExplorer() {
	stopCommand := exec.Command("Taskkill", "/IM", "explorer.exe", "/F")
	startCommand := exec.Command("explorer.exe")

	stopCommand.Run()
	startCommand.Start()
}

func main() {
	args := os.Args[1:]
	iconPositions := getKey()

	if len(args) > 0 && args[0] == "save" {
		fmt.Printf("Exporting Key to file 'regKeyData'")
		export(iconPositions)
	} else {
		//make sure file exists before we do anything
		if _, err := os.Stat("regKeyData"); os.IsNotExist(err) {
			panic("Prior Saved Reg Key Data doesn't exist!  Please Run 'IconRestore save' to create it first!")
		}

		clear(iconPositions)
		importKey(iconPositions)
		restartExplorer()
	}

}
