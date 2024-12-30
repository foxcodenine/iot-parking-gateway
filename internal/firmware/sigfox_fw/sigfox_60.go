package sigfoxfw

import (
	"fmt"

	"github.com/foxcodenine/iot-parking-gateway/internal/helpers"
)

func Sigfox_60(hexStr string, timestamp int) (map[string]any, error) {
	myMap := map[string]any{}
	pkgAmount, parkingAmount, keepAliveAmount, settingsAmount := 0, 0, 0, 0

	var parkingPackages, keepAlivePackages, settingsPackages []map[string]any

	// Parse firmware version
	firmwareVersionTmp, nextOffset, err := helpers.ParseHexSubstring(hexStr, 0, 1)
	if err != nil {
		return nil, helpers.WrapError(err)
	}

	firmwareVersion := float64(firmwareVersionTmp) / 10.0

	myMap["firmware_version"] = firmwareVersion

	for nextOffset*2 < len(hexStr) {
		remain := hexStr[nextOffset:]

		eventID, nextOffset1, err := helpers.ParseHexSubstring(hexStr, nextOffset, 1)
		if err != nil {
			return nil, helpers.WrapError(err)
		}

		if !isValidEventID(eventID) {
			return nil, helpers.WrapError(fmt.Errorf("invalid event_id: %d, remain: %s", eventID, remain))
		}

		pkgAmount++

		switch eventID {
		case 26, 31:
			parkingAmount++
			pkg, err := parseParkingPackage60(hexStr, timestamp, nextOffset1)
			if err != nil {
				return nil, err
			}
			parkingPackages = append(parkingPackages, pkg)
			nextOffset1 = pkg["nextOffset"].(int)
		}

		nextOffset = nextOffset1
	}

	myMap["firmware_version"] = firmwareVersion

	myMap["pkg_amount"] = pkgAmount
	myMap["parking_amount"] = parkingAmount
	myMap["keep_alive_amount"] = keepAliveAmount
	myMap["settings_amount"] = settingsAmount
	myMap["parking_packages"] = parkingPackages
	myMap["keep_alive_packages"] = keepAlivePackages
	myMap["settings_packages"] = settingsPackages

	return myMap, nil
}

func parseParkingPackage60(hexStr string, timestamp, offset int) (map[string]any, error) {
	pkg := map[string]any{"timestamp": timestamp}
	var err error
	var nextOffset int

	if pkg["is_occupied"], nextOffset, err = helpers.ParseHexSubstring(hexStr, offset, 1); err != nil {
		return nil, helpers.WrapError(err)
	}

	if pkg["radar_cumulative"], nextOffset, err = helpers.ParseHexSubstring(hexStr, nextOffset, 1); err != nil {
		return nil, helpers.WrapError(err)
	}
	pkg["radar_cumulative"] = pkg["radar_cumulative"].(int) * 256

	// Parsing beacons array
	beacons := []map[string]any{}
	beaconsAmount := 0

	for {
		if nextOffset*2 >= len(hexStr) {
			break
		}

		beaconsAmount++
		var nextOffset1 int
		beacon := map[string]any{"beacon_number": beaconsAmount}

		if beacon["major"], nextOffset1, err = helpers.ParseHexSubstring(hexStr, nextOffset, 2); err != nil {
			return nil, helpers.WrapError(err)
		}

		if beacon["minor"], nextOffset1, err = helpers.ParseHexSubstring(hexStr, nextOffset1, 2); err != nil {
			return nil, helpers.WrapError(err)
		}

		beacons = append(beacons, beacon)

		// Ensure nextOffset is updated
		if nextOffset1 <= nextOffset {
			return nil, fmt.Errorf("nextOffset did not advance: nextOffset=%d, nextOffset1=%d", nextOffset, nextOffset1)
		}
		nextOffset = nextOffset1
	}
	pkg["beacons_amount"] = beaconsAmount
	pkg["beacons"] = beacons
	pkg["nextOffset"] = nextOffset
	return pkg, nil
}

// Validates event ID
func isValidEventID(eventID int) bool {
	validEventIDs := []int{6, 10, 26, 31}
	return helpers.Contains(validEventIDs, eventID)
}
