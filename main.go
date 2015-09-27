package main

import (
	"flag"
	"log"
	"strings"
	"encoding/hex"
	"encoding/json"
	"fmt"

	"github.com/2tvenom/golifx"
	"errors"
)

type (
	silent bool
)

var (
	lookupBulbs      bool
	showLabel        bool
	showHostInfo     bool
	showHostFirmware bool
	showWifiInfo     bool
	showWifiFirmware bool
	showPowerState   bool
	showBulbVersion  bool
	showInfo         bool
	showLocation     bool
	showGroup        bool
	showColor        bool

	mac string
	turnOn bool
	turnOff bool
	duration int

	hue        int
	saturation int
	brightness int
	kelvin     int

	jsonFormat bool
	isSilent   silent = silent(*flag.Bool("silent", false, "Silent"))
)

const (
	_MAX = 65535
)

func init() {
	flag.BoolVar(&lookupBulbs, "lookup", false, "Lookup bulbs")
	flag.BoolVar(&showLabel, "label", false, "Get bulb label. Use with --lookup")
	flag.BoolVar(&showHostInfo, "hostinfo", false, "Get bulb host info. Use with --lookup")
	flag.BoolVar(&showHostFirmware, "hostfirmware", false, "Get bulb host firmware. Use with --lookup")
	flag.BoolVar(&showWifiInfo, "wifiinfo", false, "Get bulb Wi-Fi info. Use with --lookup")
	flag.BoolVar(&showWifiFirmware, "wififirmware", false, "Get bulb Wi-Fi firmware. Use with --lookup")
	flag.BoolVar(&showPowerState, "powerstate", false, "Get bulb power state. Use with --lookup")
	flag.BoolVar(&showBulbVersion, "bulbversion", false, "Get bulb version. Use with --lookup")
	flag.BoolVar(&showInfo, "info", false, "Get bulb info. Use with --lookup")
	flag.BoolVar(&showLocation, "location", false, "Get bulb location. Use with --lookup")
	flag.BoolVar(&showGroup, "group", false, "Get bulb group. Use with --lookup")
	flag.BoolVar(&showColor, "color", false, "Get bulb color, label, power state. Use with --lookup")

	flag.StringVar(&mac, "bulb", "", "Set bulb MAC address for set color and power state. MAC format like 56:84:7a:fe:97:99")
	flag.BoolVar(&turnOn, "on", false, "Turn on bulb. Use with --bulb")
	flag.BoolVar(&turnOff, "off", false, "Turn off bulb. Use with --bulb")
	flag.IntVar(&duration, "duration", 0, "Turn off/Turn on duration. Use with --bulb and (--on or --off)")

	flag.IntVar(&hue, "hue", -1, "Set color hue. Range 0 to 65535. Use with --bulb")
	flag.IntVar(&saturation, "saturation", -1, "Set color saturation. Range 0 to 65535. Use with --bulb")
	flag.IntVar(&brightness, "brightness", -1, "Set color brightness. Range 0 to 65535. Use with --bulb")
	flag.IntVar(&kelvin, "kelvin", -1, "Set color kelvin. Range 2500° (warm) to 9000° (cool). Use with --bulb")

	flag.BoolVar(&jsonFormat, "json", false, "JSON format output")
}

func main() {
	flag.Parse()

	if mac != "" {
		hwAddr, err := parseMac(mac)
		if err != nil {
			isSilent.error(err)
		}

		bulb := &golifx.Bulb{}
		bulb.SetHardwareAddress(hwAddr)

		if turnOn {
			if duration == 0 {
				bulb.SetPowerState(true)
			} else {
				bulb.SetPowerDurationState(true, uint32(duration))
			}
			return
		}

		if turnOff {
			if duration == 0 {
				bulb.SetPowerState(false)
			} else {
				bulb.SetPowerDurationState(false, uint32(duration))
			}
			return
		}

		if hue != -1 && (hue < 0 || hue > _MAX) {
			isSilent.error(errors.New("Incorrect hue"))
		}

		if saturation != -1 && (saturation < 0 || saturation > _MAX) {
			isSilent.error(errors.New("Incorrect saturation"))
		}

		if brightness != -1 && (brightness < 0 || brightness > _MAX) {
			isSilent.error(errors.New("Incorrect brightness"))
		}

		if kelvin != -1 && (kelvin < 2500 || kelvin > 9000) {
			isSilent.error(errors.New("Incorrect kelvin"))
		}

		if hue != -1 || saturation != -1 || brightness != -1 || kelvin != -1 {
			var err error
			var color *golifx.HSBK

			if hue != -1 && saturation != -1 && brightness != -1 && kelvin != -1 {
				color = &golifx.HSBK{
					Hue: uint16(hue),
					Saturation: uint16(saturation),
					Brightness: uint16(brightness),
					Kelvin: uint16(kelvin),
				}
			} else {
				bulbState, err := bulb.GetColorState()
				color = bulbState.Color

				if err != nil {
					isSilent.error(err)
				}

				if hue != -1 {
					color.Hue = uint16(hue)
				}

				if saturation != -1 {
					color.Saturation = uint16(saturation)
				}

				if brightness != -1 {
					color.Brightness = uint16(brightness)
				}

				if kelvin != -1 {
					color.Kelvin = uint16(kelvin)
				}
			}

			bulbState, err := bulb.SetColorStateWithResponse(color, uint32(duration))

			if err != nil {
				isSilent.error(err)
			}

			if !isSilent {
				if jsonFormat {
					b, _ := json.Marshal(bulbState)
					fmt.Println(string(b))
				} else {
					fmt.Println(bulbState)
				}
			}
		}
		return
	}

	if lookupBulbs {
		bulbs, err := golifx.LookupBulbs()

		if err != nil {
			isSilent.error(err)
		}

		for _, bulb := range bulbs {
			if showLabel {
				_, err := bulb.GetLabel()
				if err != nil {
					isSilent.error(err)
				}
			}

			if showHostInfo {
				_, err := bulb.GetStateHostInfo()
				if err != nil {
					isSilent.error(err)
				}
			}

			if showHostFirmware {
				_, err := bulb.GetHostFirmware()
				if err != nil {
					isSilent.error(err)
				}
			}

			if showWifiInfo {
				_, err := bulb.GetWifiInfo()
				if err != nil {
					isSilent.error(err)
				}
			}

			if showWifiFirmware {
				_, err := bulb.GetWifiFirmware()
				if err != nil {
					isSilent.error(err)
				}
			}

			if showPowerState {
				_, err := bulb.GetPowerState()
				if err != nil {
					isSilent.error(err)
				}
			}

			if showBulbVersion {
				_, err := bulb.GetVersion()
				if err != nil {
					isSilent.error(err)
				}
			}

			if showInfo {
				_, err := bulb.GetInfo()
				if err != nil {
					isSilent.error(err)
				}
			}

			if showLocation {
				_, err := bulb.GetLocation()
				if err != nil {
					isSilent.error(err)
				}
			}

			if showGroup {
				_, err := bulb.GetGroup()
				if err != nil {
					isSilent.error(err)
				}
			}

			if showColor {
				_, err := bulb.GetColorState()
				if err != nil {
					isSilent.error(err)
				}
			}
		}

		if jsonFormat {
			b, _ := json.Marshal(bulbs)
			fmt.Println(string(b))
		} else {
			for id, bulb := range bulbs {
				fmt.Print(bulb)
				if id < len(bulbs)-1 {
					fmt.Println("----------------")
				}
			}
		}

		return
	}
}

func (s silent) error(err error) {
	if !s {
		log.Fatalf("Error: %s\n", err)
	}
}

func parseMac(mac string) (uint64, error) {
	mac = strings.Replace(mac, ":", "", -1)
	hex, err := hex.DecodeString(mac)
	if err != nil {
		return 0, err
	}
	hex = append(hex, []byte{0,0}...)
	var digit uint64
	readUint64(hex, &digit)
	return digit, nil
}

func readUint64(buff []byte, dest *uint64) error {
	*dest = 0
	for i := 0; i < 8; i++ {
		*dest += uint64(buff[i]&0xFF) << uint(i*8)
	}

	return nil
}