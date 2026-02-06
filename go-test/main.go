package main

import (
	"fmt"
	"github.com/holoplot/go-evdev"
	"strconv"
)

func GetDeviceInfo(device *evdev.InputDevice) (string, error) {
	var name, errName = device.Name()
	if errName != nil {
		fmt.Printf("Error: %s\n", errName)
		return "", errName
	}

	var inputId, errId = device.InputID()
	if errId != nil {
		fmt.Printf("Error: %s\n", errId)
		return name, errId
	}

	var infoString string = fmt.Sprintf("Device Name: %s\n", name)
	infoString += fmt.Sprintf("Bus Type: %d\n", inputId.BusType)
	infoString += fmt.Sprintf("Vendor ID: %d\n", inputId.Vendor)
	infoString += fmt.Sprintf("Product ID: %d\n", inputId.Product)
	infoString += fmt.Sprintf("Version: %d\n", inputId.Version)
	return infoString, nil
}

func GetDeviceProperties(device *evdev.InputDevice) (string, error) {
	var properties []evdev.EvProp = device.Properties()
	if len(properties) == 0 {
		return "No special properties found for this device.", nil
	}

	var propertiesString string = "Device Properties:\n"
	for _, prop := range properties {
		fmt.Printf("Property: %d\n", prop)
		switch prop {
		case evdev.INPUT_PROP_POINTER:
			propertiesString += "- Pointer\n"
		case evdev.INPUT_PROP_DIRECT:
			propertiesString += "- Direct\n"
		case evdev.INPUT_PROP_BUTTONPAD:
			propertiesString += "- Button Pad\n"
		case evdev.INPUT_PROP_SEMI_MT:
			propertiesString += "- Semi-Multitouch\n"
		case evdev.INPUT_PROP_TOPBUTTONPAD:
			propertiesString += "- Top Button Pad\n"
		case evdev.INPUT_PROP_POINTING_STICK:
			propertiesString += "- Pointing Stick\n"
		case evdev.INPUT_PROP_ACCELEROMETER:
			propertiesString += "- Accelerometer\n"
		default:
			propertiesString += fmt.Sprintf("- Unknown Property: %d\n", prop)
		}
	}

	return propertiesString, nil
}

func GetDeviceAbsInfos(device *evdev.InputDevice) (string, error) {
	var absInfos map[evdev.EvCode]evdev.AbsInfo = nil
	var err error = nil

	absInfos, err = device.AbsInfos()
	if err != nil {
		return "", err
	}

	if len(absInfos) == 0 {
		return "No absolute axis information found for this device.", nil
	}

	var absInfosString string = "Device Absolute Axis Information:\n"
	for code, info := range absInfos {
		absInfosString += fmt.Sprintf("- Code: %d, Fuzz: %d, Flat: %d, Resolution: %d\n",
			code, info.Fuzz, info.Flat, info.Resolution)
	}

	return absInfosString, nil
}

func GetDeviceCapableTypes(device *evdev.InputDevice) (string, error) {
	var evTypes []evdev.EvType = device.CapableTypes()
	if len(evTypes) == 0 {
		return "No event types found for this device.", nil
	}

	var evTypesString string = "Device Capable Event Types:\n"
	for _, t := range evTypes {
		evTypesString += fmt.Sprintf("- Event Type: 0x%02x [%d]\n", t, t)
	}

	return evTypesString, nil
}

func GetDeviceCapableEvents(device *evdev.InputDevice, evType evdev.EvType) (string, error) {
	var evCodes []evdev.EvCode = device.CapableEvents(evType)
	if len(evCodes) == 0 {
		return "No codes found for this device.", nil
	}

	var evCodesString string = fmt.Sprintf("Device Capable Codes for Event Type %d:\n", evType)
	for _, code := range evCodes {
		evCodesString += fmt.Sprintf("- Code: 0x%02x [%d]\n", code, code)
	}

	return evCodesString, nil
}

func main() {
	var eventPath string = "/dev/input/event"
	var eventId int = 7
	var fullEventPath string = eventPath + strconv.Itoa(eventId)

	var device, errPath = evdev.Open(fullEventPath)
	if errPath != nil {
		fmt.Printf("Error: %s\n", errPath)
		return
	}

	var info, errInfo = GetDeviceInfo(device)
	if errInfo != nil {
		fmt.Printf("Error: %s\n", errInfo)
		device.Close()
		return
	}

	var properties, errProperties = GetDeviceProperties(device)
	if errProperties != nil {
		fmt.Printf("Error: %s\n", errProperties)
		device.Close()
		return
	}

	var absInfos, errAbsInfos = GetDeviceAbsInfos(device)
	if errAbsInfos != nil {
		fmt.Printf("Error: %s\n", errAbsInfos)
		device.Close()
		return
	}

	var capableTypes, errCapableTypes = GetDeviceCapableTypes(device)
	if errCapableTypes != nil {
		fmt.Printf("Error: %s\n", errCapableTypes)
		device.Close()
		return
	}

	var capableEvents, errCapableEvents = GetDeviceCapableEvents(device, evdev.EV_KEY)
	if errCapableEvents != nil {
		fmt.Printf("Error: %s\n", errCapableEvents)
		device.Close()
		return
	}

	fmt.Println(info)
	fmt.Println(properties)
	fmt.Println(absInfos)
	fmt.Println(capableTypes)
	fmt.Println(capableEvents)

	fmt.Print("Reading events (press 'q' to quit, not on keyboard? ctrl +c...)\n")
	for true {
		var event, errRead = device.ReadOne()
		if errRead != nil {
			fmt.Printf("Error: %s\n", errRead)
		} else {
			fmt.Printf("Read Event: %s\n", event)
			if event.Type == evdev.EV_KEY && event.Value == 1 {
				if event.Code == evdev.KEY_Q {
					break
				}

				if event.Code == evdev.KEY_C {
					fmt.Print("\033[H\033[2J") // Clear console
				}
			}

		}
	}
	device.Close()
}
