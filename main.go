package main

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"

	"golang.org/x/sys/windows/registry"
)

// RegValueInfo - Holds information used to describe a field/value pair within a registry key
type RegValueInfo struct {
	Key     string
	Buf     []byte
	ValType uint32
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func getKey() registry.Key {
	iconPositions, err := registry.OpenKey(registry.CURRENT_USER, "SOFTWARE\\Microsoft\\Windows\\Shell\\Bags\\1\\Desktop", registry.ALL_ACCESS)

	checkErr(err)

	return iconPositions
}

func getKeyDataPath() (keyDataPath string) {
	//make sure file exists before we do anything
	exepath, _ := os.Executable()
	dir := filepath.Dir(exepath)
	keyDataPath = dir + string(os.PathSeparator) + "regKeyData"
	return
}

func export(iconPositions registry.Key) {

	KeyValPairs := []RegValueInfo{}

	valueNames, readErr := iconPositions.ReadValueNames(-1)
	checkErr(readErr)

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
			valType,
		})
	}

	fileContents, jsonErr := json.Marshal(KeyValPairs)
	checkErr(jsonErr)

	fileErr := ioutil.WriteFile(getKeyDataPath(), fileContents, 0777)
	checkErr(fileErr)
}

func clear(iconPositions registry.Key) {
	valueNames, _ := iconPositions.ReadValueNames(-1)

	for i := 0; i < len(valueNames); i++ {
		checkErr(iconPositions.DeleteValue(valueNames[i]))
	}
}

func importKey(iconPositions registry.Key) {
	var KeyValPairs []RegValueInfo

	fileData, _ := ioutil.ReadFile(getKeyDataPath())

	json.Unmarshal(fileData, &KeyValPairs)

	for i := 0; i < len(KeyValPairs); i++ {
		curData := KeyValPairs[i]
		switch curData.ValType {
		case registry.DWORD:
			value := binary.LittleEndian.Uint32(curData.Buf)
			checkErr(iconPositions.SetDWordValue(curData.Key, value))
		case registry.BINARY:
			checkErr(iconPositions.SetBinaryValue(curData.Key, curData.Buf))
		case registry.SZ:
			binData := bytes.NewBuffer(curData.Buf)
			strLength := binData.Len()
			strData := string(binData.Bytes()[:strLength])
			checkErr(iconPositions.SetStringValue(curData.Key, strData))
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
		keyDataPath := getKeyDataPath()
		if _, err := os.Stat(keyDataPath); os.IsNotExist(err) {
			panic("Prior Saved Reg Key Data doesn't exist!  Please Run 'IconRestore save' to create it first!")
		}

		clear(iconPositions)
		importKey(iconPositions)
		iconPositions.Close()

		restartExplorer()
	}

}
