package lorafw

import (
	"fmt"
	"math"

	"github.com/foxcodenine/iot-parking-gateway/internal/helpers"
)

func Lora_58(hexStr string) (map[string]any, error) {
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

		timestamp, nextOffset1, err := helpers.ParseHexSubstring(hexStr, nextOffset, 4)
		if err != nil {
			return nil, helpers.WrapError(err)
		}

		eventID, nextOffset1, err := helpers.ParseHexSubstring(hexStr, nextOffset1, 1)
		if err != nil {
			return nil, helpers.WrapError(err)
		}

		if !isValidEventID(eventID) {
			return nil, helpers.WrapError(fmt.Errorf("invalid event_id: %d, remain: %s", eventID, remain))
		}

		pkgAmount++

		switch eventID {
		case 26:
			parkingAmount++
			pkg, err := parseParkingPackage58(hexStr, timestamp, nextOffset1)
			if err != nil {
				return nil, err
			}
			parkingPackages = append(parkingPackages, pkg)
			nextOffset1 = pkg["nextOffset"].(int)
		case 6:
			keepAliveAmount++
			pkg, err := parseKeepAlivePackage58(hexStr, timestamp, nextOffset1, keepAliveAmount)
			if err != nil {
				return nil, err
			}
			keepAlivePackages = append(keepAlivePackages, pkg)
			nextOffset1 = pkg["nextOffset"].(int)
		case 10:
			settingsAmount++
			settingsAmount++
			pkg, err := parseSettingsPackage58(hexStr, timestamp, nextOffset1)
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
func parseParkingPackage58(hexStr string, timestamp, offset int) (map[string]any, error) {
	pkg := map[string]any{"timestamp": timestamp}
	var err error
	var nextOffset int

	if pkg["peak_distance_cm"], nextOffset, err = helpers.ParseHexSubstring(hexStr, offset, 1); err != nil {
		return nil, helpers.WrapError(err)
	}

	if pkg["is_occupied"], nextOffset, err = helpers.ParseHexSubstring(hexStr, nextOffset, 1); err != nil {
		return nil, helpers.WrapError(err)
	}

	if pkg["radar_cumulative"], nextOffset, err = helpers.ParseHexSubstring(hexStr, nextOffset, 1); err != nil {
		return nil, helpers.WrapError(err)
	}
	pkg["radar_cumulative"] = pkg["radar_cumulative"].(int) * 256

	if pkg["magnet_abs_total"], nextOffset, err = helpers.ParseHexSubstring(hexStr, nextOffset, 2); err != nil {
		return nil, helpers.WrapError(err)
	}

	if pkg["beacons_amount"], nextOffset, err = helpers.ParseHexSubstring(hexStr, nextOffset, 1); err != nil {
		return nil, helpers.WrapError(err)
	}

	// Parsing beacons array
	beacons := []map[string]any{}
	beaconsAmount := pkg["beacons_amount"].(int)

	for i := 1; i <= beaconsAmount; i++ {
		var nextOffset1 int
		beacon := map[string]any{"beacon_number": i}

		if beacon["major"], nextOffset1, err = helpers.ParseHexSubstring(hexStr, nextOffset, 2); err != nil {
			return nil, helpers.WrapError(err)
		}
		if beacon["minor"], nextOffset1, err = helpers.ParseHexSubstring(hexStr, nextOffset1, 2); err != nil {
			return nil, helpers.WrapError(err)
		}
		if beacon["rssi"], nextOffset1, err = helpers.ParseHexSubstring(hexStr, nextOffset1, 1); err != nil {
			return nil, helpers.WrapError(err)
		}
		beacons = append(beacons, beacon)

		nextOffset = nextOffset1
	}

	pkg["beacons"] = beacons
	pkg["nextOffset"] = nextOffset
	return pkg, nil
}

// Parses the Keep Alive Package
func parseKeepAlivePackage58(hexStr string, timestamp, offset, keepAliveAmount int) (map[string]any, error) {
	pkg := map[string]any{"timestamp": timestamp}
	var err error
	var nextOffset int

	if pkg["idle_voltage"], nextOffset, err = helpers.ParseHexSubstring(hexStr, offset, 2); err != nil {
		return nil, helpers.WrapError(err)
	}
	pkg["idle_voltage"] = int(math.Floor(float64(pkg["idle_voltage"].(int)) * 0.2197))

	if pkg["battery_percentage"], nextOffset, err = helpers.ParseHexSubstring(hexStr, nextOffset, 1); err != nil {
		return nil, helpers.WrapError(err)
	}

	if pkg["current"], nextOffset, err = helpers.ParseHexSubstring(hexStr, nextOffset, 2); err != nil {
		return nil, helpers.WrapError(err)
	}
	pkg["current"] = int(math.Floor(float64(pkg["current"].(int)) * 0.6104))

	if pkg["reset_count"], nextOffset, err = helpers.ParseHexSubstring(hexStr, nextOffset, 1); err != nil {
		return nil, helpers.WrapError(err)
	}

	if pkg["manual_calibration"], nextOffset, err = helpers.ParseHexSubstring(hexStr, nextOffset, 1); err != nil {
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

	if pkg["mag_error"], nextOffset, err = helpers.ParseHexSubstring(hexStr, nextOffset, 1); err != nil {
		return nil, helpers.WrapError(err)
	}

	if pkg["tcve_error"], nextOffset, err = helpers.ParseHexSubstring(hexStr, nextOffset, 1); err != nil {
		return nil, helpers.WrapError(err)
	}

	if pkg["ble_security_issues"], nextOffset, err = helpers.ParseHexSubstring(hexStr, nextOffset, 1); err != nil {
		return nil, helpers.WrapError(err)
	}

	if pkg["radar_cumulative_total"], nextOffset, err = helpers.ParseHexSubstring(hexStr, nextOffset, 1); err != nil {
		return nil, helpers.WrapError(err)
	}
	pkg["radar_cumulative_total"] = pkg["radar_cumulative_total"].(int) * 256

	if pkg["mag_total"], nextOffset, err = helpers.ParseHexSubstring(hexStr, nextOffset, 2); err != nil {
		return nil, helpers.WrapError(err)
	}

	if pkg["network_registration_ok"], nextOffset, err = helpers.ParseHexSubstring(hexStr, nextOffset, 1); err != nil {
		return nil, helpers.WrapError(err)
	}

	if pkg["network_registration_nok"], nextOffset, err = helpers.ParseHexSubstring(hexStr, nextOffset, 1); err != nil {
		return nil, helpers.WrapError(err)
	}

	if pkg["rssi_average"], nextOffset, err = helpers.ParseHexSubstring(hexStr, nextOffset, 1); err != nil {
		return nil, helpers.WrapError(err)
	}

	if pkg["network_message_attempts"], nextOffset, err = helpers.ParseHexSubstring(hexStr, nextOffset, 1); err != nil {
		return nil, helpers.WrapError(err)
	}

	if pkg["network_ack_1ds"], nextOffset, err = helpers.ParseHexSubstring(hexStr, nextOffset, 1); err != nil {
		return nil, helpers.WrapError(err)
	}

	if pkg["network_1ack_ds"], nextOffset, err = helpers.ParseHexSubstring(hexStr, nextOffset, 1); err != nil {
		return nil, helpers.WrapError(err)
	}

	if pkg["network_1ack_1ds"], nextOffset, err = helpers.ParseHexSubstring(hexStr, nextOffset, 1); err != nil {
		return nil, helpers.WrapError(err)
	}

	if pkg["tcvr_deep_sleep_min"], nextOffset, err = helpers.ParseHexSubstring(hexStr, nextOffset, 2); err != nil {
		return nil, helpers.WrapError(err)
	}

	if pkg["tcvr_deep_sleep_max"], nextOffset, err = helpers.ParseHexSubstring(hexStr, nextOffset, 2); err != nil {
		return nil, helpers.WrapError(err)
	}

	if pkg["tcvr_deep_sleep_average"], nextOffset, err = helpers.ParseHexSubstring(hexStr, nextOffset, 2); err != nil {
		return nil, helpers.WrapError(err)
	}

	if pkg["settings_checksum"], nextOffset, err = helpers.ParseHexSubstring(hexStr, nextOffset, 1); err != nil {
		return nil, helpers.WrapError(err)
	}

	if pkg["time_sync_rand_byte"], nextOffset, err = helpers.ParseHexSubstring(hexStr, nextOffset, 1); err != nil {
		return nil, helpers.WrapError(err)
	}

	if keepAliveAmount == 1 || pkg["time_sync_rand_byte"].(int) > 0 {
		if pkg["time_sync_current_unix_time"], nextOffset, err = helpers.ParseHexSubstring(hexStr, nextOffset, 4); err != nil {
			return nil, helpers.WrapError(err)
		}
	}

	pkg["nextOffset"] = nextOffset
	return pkg, nil
}

// Parses the Settings Package
func parseSettingsPackage58(hexStr string, timestamp, offset int) (map[string]any, error) {
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

	if pkg["radar_car_uncal_lo_th"], nextOffset, err = helpers.ParseHexSubstring(hexStr, nextOffset, 1); err != nil {
		return nil, helpers.WrapError(err)
	}
	pkg["radar_car_uncal_lo_th"] = pkg["radar_car_uncal_lo_th"].(int) * 256

	if pkg["radar_car_uncal_hi_th"], nextOffset, err = helpers.ParseHexSubstring(hexStr, nextOffset, 1); err != nil {
		return nil, helpers.WrapError(err)
	}
	pkg["radar_car_uncal_hi_th"] = pkg["radar_car_uncal_hi_th"].(int) * 256

	if pkg["radar_car_delta_th"], nextOffset, err = helpers.ParseHexSubstring(hexStr, nextOffset, 1); err != nil {
		return nil, helpers.WrapError(err)
	}
	pkg["radar_car_delta_th"] = pkg["radar_car_delta_th"].(int) * 256

	if pkg["mag_car_lo"], nextOffset, err = helpers.ParseHexSubstring(hexStr, nextOffset, 2); err != nil {
		return nil, helpers.WrapError(err)
	}
	if pkg["mag_car_hi"], nextOffset, err = helpers.ParseHexSubstring(hexStr, nextOffset, 2); err != nil {
		return nil, helpers.WrapError(err)
	}

	if pkg["debug_period"], nextOffset, err = helpers.ParseHexSubstring(hexStr, nextOffset, 1); err != nil {
		return nil, helpers.WrapError(err)
	}
	if pkg["debug_mode"], nextOffset, err = helpers.ParseHexSubstring(hexStr, nextOffset, 1); err != nil {
		return nil, helpers.WrapError(err)
	}

	if pkg["logs_mode"], nextOffset, err = helpers.ParseHexSubstring(hexStr, nextOffset, 1); err != nil {
		return nil, helpers.WrapError(err)
	}
	if pkg["logs_amount"], nextOffset, err = helpers.ParseHexSubstring(hexStr, nextOffset, 1); err != nil {
		return nil, helpers.WrapError(err)
	}

	if pkg["maximum_registration_time"], nextOffset, err = helpers.ParseHexSubstring(hexStr, nextOffset, 1); err != nil {
		return nil, helpers.WrapError(err)
	}
	pkg["maximum_registration_time"] = pkg["maximum_registration_time"].(int) / 4

	if pkg["maximum_registration_attempts"], nextOffset, err = helpers.ParseHexSubstring(hexStr, nextOffset, 1); err != nil {
		return nil, helpers.WrapError(err)
	}
	if pkg["maximum_deep_sleep_time"], nextOffset, err = helpers.ParseHexSubstring(hexStr, nextOffset, 1); err != nil {
		return nil, helpers.WrapError(err)
	}
	pkg["maximum_deep_sleep_time"] = pkg["maximum_deep_sleep_time"].(int) / 2

	if pkg["deep_sleep_time_1"], nextOffset, err = helpers.ParseHexSubstring(hexStr, nextOffset, 2); err != nil {
		return nil, helpers.WrapError(err)
	}
	pkg["deep_sleep_time_1"] = pkg["deep_sleep_time_1"].(int) * 4

	if pkg["action_before_1"], nextOffset, err = helpers.ParseHexSubstring(hexStr, nextOffset, 1); err != nil {
		return nil, helpers.WrapError(err)
	}

	if pkg["action_after_1"], nextOffset, err = helpers.ParseHexSubstring(hexStr, nextOffset, 1); err != nil {
		return nil, helpers.WrapError(err)
	}

	if pkg["deep_sleep_time_2"], nextOffset, err = helpers.ParseHexSubstring(hexStr, nextOffset, 2); err != nil {
		return nil, helpers.WrapError(err)
	}
	pkg["deep_sleep_time_2"] = pkg["deep_sleep_time_2"].(int) * 4

	if pkg["action_before_2"], nextOffset, err = helpers.ParseHexSubstring(hexStr, nextOffset, 1); err != nil {
		return nil, helpers.WrapError(err)
	}

	if pkg["action_after_2"], nextOffset, err = helpers.ParseHexSubstring(hexStr, nextOffset, 1); err != nil {
		return nil, helpers.WrapError(err)
	}

	if pkg["deep_sleep_time_3"], nextOffset, err = helpers.ParseHexSubstring(hexStr, nextOffset, 2); err != nil {
		return nil, helpers.WrapError(err)
	}
	pkg["deep_sleep_time_3"] = pkg["deep_sleep_time_3"].(int) * 4

	if pkg["action_before_3"], nextOffset, err = helpers.ParseHexSubstring(hexStr, nextOffset, 1); err != nil {
		return nil, helpers.WrapError(err)
	}

	if pkg["action_after_3"], nextOffset, err = helpers.ParseHexSubstring(hexStr, nextOffset, 1); err != nil {
		return nil, helpers.WrapError(err)
	}

	if pkg["deep_sleep_time_4"], nextOffset, err = helpers.ParseHexSubstring(hexStr, nextOffset, 2); err != nil {
		return nil, helpers.WrapError(err)
	}
	pkg["deep_sleep_time_4"] = pkg["deep_sleep_time_4"].(int) * 4

	if pkg["action_before_4"], nextOffset, err = helpers.ParseHexSubstring(hexStr, nextOffset, 1); err != nil {
		return nil, helpers.WrapError(err)
	}

	if pkg["action_after_4"], nextOffset, err = helpers.ParseHexSubstring(hexStr, nextOffset, 1); err != nil {
		return nil, helpers.WrapError(err)
	}

	if pkg["deep_sleep_time_5"], nextOffset, err = helpers.ParseHexSubstring(hexStr, nextOffset, 2); err != nil {
		return nil, helpers.WrapError(err)
	}
	pkg["deep_sleep_time_5"] = pkg["deep_sleep_time_5"].(int) * 4

	if pkg["action_before_5"], nextOffset, err = helpers.ParseHexSubstring(hexStr, nextOffset, 1); err != nil {
		return nil, helpers.WrapError(err)
	}

	if pkg["action_after_5"], nextOffset, err = helpers.ParseHexSubstring(hexStr, nextOffset, 1); err != nil {
		return nil, helpers.WrapError(err)
	}

	if pkg["deep_sleep_time_6"], nextOffset, err = helpers.ParseHexSubstring(hexStr, nextOffset, 2); err != nil {
		return nil, helpers.WrapError(err)
	}
	pkg["deep_sleep_time_6"] = pkg["deep_sleep_time_6"].(int) * 4

	if pkg["action_before_6"], nextOffset, err = helpers.ParseHexSubstring(hexStr, nextOffset, 1); err != nil {
		return nil, helpers.WrapError(err)
	}

	if pkg["action_after_6"], nextOffset, err = helpers.ParseHexSubstring(hexStr, nextOffset, 1); err != nil {
		return nil, helpers.WrapError(err)
	}

	if pkg["deep_sleep_time_7"], nextOffset, err = helpers.ParseHexSubstring(hexStr, nextOffset, 2); err != nil {
		return nil, helpers.WrapError(err)
	}
	pkg["deep_sleep_time_7"] = pkg["deep_sleep_time_7"].(int) * 4

	if pkg["action_before_7"], nextOffset, err = helpers.ParseHexSubstring(hexStr, nextOffset, 1); err != nil {
		return nil, helpers.WrapError(err)
	}

	if pkg["action_after_7"], nextOffset, err = helpers.ParseHexSubstring(hexStr, nextOffset, 1); err != nil {
		return nil, helpers.WrapError(err)
	}

	if pkg["deep_sleep_time_8"], nextOffset, err = helpers.ParseHexSubstring(hexStr, nextOffset, 2); err != nil {
		return nil, helpers.WrapError(err)
	}
	pkg["deep_sleep_time_8"] = pkg["deep_sleep_time_8"].(int) * 4

	if pkg["action_before_8"], nextOffset, err = helpers.ParseHexSubstring(hexStr, nextOffset, 1); err != nil {
		return nil, helpers.WrapError(err)
	}

	if pkg["action_after_8"], nextOffset, err = helpers.ParseHexSubstring(hexStr, nextOffset, 1); err != nil {
		return nil, helpers.WrapError(err)
	}

	if pkg["deep_sleep_time_9"], nextOffset, err = helpers.ParseHexSubstring(hexStr, nextOffset, 2); err != nil {
		return nil, helpers.WrapError(err)
	}
	pkg["deep_sleep_time_9"] = pkg["deep_sleep_time_9"].(int) * 4

	if pkg["action_before_9"], nextOffset, err = helpers.ParseHexSubstring(hexStr, nextOffset, 1); err != nil {
		return nil, helpers.WrapError(err)
	}

	if pkg["action_after_9"], nextOffset, err = helpers.ParseHexSubstring(hexStr, nextOffset, 1); err != nil {
		return nil, helpers.WrapError(err)
	}

	if pkg["deep_sleep_time_10"], nextOffset, err = helpers.ParseHexSubstring(hexStr, nextOffset, 2); err != nil {
		return nil, helpers.WrapError(err)
	}
	pkg["deep_sleep_time_10"] = pkg["deep_sleep_time_10"].(int) * 4

	if pkg["action_before_10"], nextOffset, err = helpers.ParseHexSubstring(hexStr, nextOffset, 1); err != nil {
		return nil, helpers.WrapError(err)
	}

	if pkg["action_after_10"], nextOffset, err = helpers.ParseHexSubstring(hexStr, nextOffset, 1); err != nil {
		return nil, helpers.WrapError(err)
	}

	if pkg["lora_data_rate"], nextOffset, err = helpers.ParseHexSubstring(hexStr, nextOffset, 1); err != nil {
		return nil, helpers.WrapError(err)
	}

	if pkg["retries"], nextOffset, err = helpers.ParseHexSubstring(hexStr, nextOffset, 1); err != nil {
		return nil, helpers.WrapError(err)
	}

	pkg["nextOffset"] = nextOffset
	return pkg, nil
}

// Validates event ID
func isValidEventID(eventID int) bool {
	validEventIDs := []int{6, 10, 26}
	return helpers.Contains(validEventIDs, eventID)
}
