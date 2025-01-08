package sigfoxfw

import (
	"fmt"
	"math"

	"github.com/foxcodenine/iot-parking-gateway/internal/helpers"
)

func Sigfox_57(hexStr string, timestamp int) (map[string]any, error) {
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
			pkg, err := parseParkingPackag57(hexStr, timestamp, nextOffset1)
			if err != nil {
				return nil, err
			}
			parkingPackages = append(parkingPackages, pkg)
			nextOffset1 = pkg["nextOffset"].(int)
		case 6:
			keepAliveAmount++
			pkg, err := parseKeepAlivePackag57(hexStr, timestamp, nextOffset1)
			if err != nil {
				return nil, err
			}
			keepAlivePackages = append(keepAlivePackages, pkg)
			nextOffset1 = pkg["nextOffset"].(int)

		case 10:
			settingsAmount++
			pkg, err := parseSettingsPackag57(hexStr, timestamp, nextOffset1)
			if err != nil {
				return nil, err
			}
			settingsPackages = append(settingsPackages, pkg)
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

// Parses the Parking Package
func parseParkingPackag57(hexStr string, timestamp, offset int) (map[string]any, error) {
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

// Parses the Keep Alive Package
func parseKeepAlivePackag57(hexStr string, timestamp, offset int) (map[string]any, error) {
	pkg := map[string]any{"timestamp": timestamp}
	var err error
	var nextOffset int

	if pkg["idle_voltage"], nextOffset, err = helpers.ParseHexSubstring(hexStr, offset, 1); err != nil {
		return nil, helpers.WrapError(err)
	}
	pkg["idle_voltage"] = int(math.Floor(float64(pkg["idle_voltage"].(int)) * 16))

	if pkg["battery_percentage"], nextOffset, err = helpers.ParseHexSubstring(hexStr, nextOffset, 1); err != nil {
		return nil, helpers.WrapError(err)
	}

	if pkg["current"], nextOffset, err = helpers.ParseHexSubstring(hexStr, nextOffset, 1); err != nil {
		return nil, helpers.WrapError(err)
	}
	pkg["current"] = int(math.Floor(float64(pkg["current"].(int)) * 0.6104))

	if pkg["reset_count"], nextOffset, err = helpers.ParseHexSubstring(hexStr, nextOffset, 1); err != nil {
		return nil, helpers.WrapError(err)
	}

	if pkg["temperature_min"], nextOffset, err = helpers.ParseHexSubstring(hexStr, nextOffset, 1); err != nil {
		return nil, helpers.WrapError(err)
	}

	if pkg["temperature_max"], nextOffset, err = helpers.ParseHexSubstring(hexStr, nextOffset, 1); err != nil {
		return nil, helpers.WrapError(err)
	}

	if pkg["radar_error"], nextOffset, err = helpers.ParseHexSubstring(hexStr, nextOffset, 1); err != nil {
		return nil, helpers.WrapError(err)
	}

	if pkg["tcve_error"], nextOffset, err = helpers.ParseHexSubstring(hexStr, nextOffset, 1); err != nil {
		return nil, helpers.WrapError(err)
	}

	if pkg["radar_cumulative"], nextOffset, err = helpers.ParseHexSubstring(hexStr, nextOffset, 1); err != nil {
		return nil, helpers.WrapError(err)
	}
	pkg["radar_cumulative"] = int(math.Floor(float64(pkg["radar_cumulative"].(int)) * 256))

	if pkg["settings_checksum"], nextOffset, err = helpers.ParseHexSubstring(hexStr, nextOffset, 1); err != nil {
		return nil, helpers.WrapError(err)
	}

	pkg["nextOffset"] = nextOffset
	return pkg, nil
}

// Parses the Settings Package
func parseSettingsPackag57(hexStr string, timestamp, offset int) (map[string]any, error) {
	pkg := map[string]any{"timestamp": timestamp}
	var err error
	var nextOffset int

	if pkg["device_mode"], nextOffset, err = helpers.ParseHexSubstring(hexStr, offset, 1); err != nil {
		return nil, helpers.WrapError(err)
	}

	if pkg["device_enable"], nextOffset, err = helpers.ParseHexSubstring(hexStr, nextOffset, 1); err != nil {
		return nil, helpers.WrapError(err)
	}

	if pkg["radar_car_cal_lo_th"], nextOffset, err = helpers.ParseHexSubstring(hexStr, nextOffset, 1); err != nil {
		return nil, helpers.WrapError(err)
	}
	pkg["radar_car_cal_lo_th"] = pkg["radar_car_cal_lo_th"].(int) * 256

	if pkg["radar_car_cal_hi_th"], nextOffset, err = helpers.ParseHexSubstring(hexStr, nextOffset, 1); err != nil {
		return nil, helpers.WrapError(err)
	}
	pkg["radar_car_cal_hi_th"] = pkg["radar_car_cal_hi_th"].(int) * 256

	if pkg["radar_car_delta_th"], nextOffset, err = helpers.ParseHexSubstring(hexStr, nextOffset, 1); err != nil {
		return nil, helpers.WrapError(err)
	}
	pkg["radar_car_delta_th"] = pkg["radar_car_delta_th"].(int) * 256

	if pkg["downlink_en_7_bits_repeated_occupancy_period_mins"], nextOffset, err = helpers.ParseHexSubstring(hexStr, nextOffset, 1); err != nil {
		return nil, helpers.WrapError(err)
	}

	pkg["nextOffset"] = nextOffset
	return pkg, nil
}
