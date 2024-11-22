package firmware

import (
	"fmt"
	"math/big"

	"github.com/foxcodenine/iot-parking-gateway/internal/helpers"
)

func NB_53(hexStr string) (map[string]any, error) {

	myMap := map[string]any{}
	pkgAmount, parkingAmount, keepAliveAmount, settingsAmount := 0, 0, 0, 0

	var parkingPackages, keepAlivePackages, settingsPackages []map[string]any

	// Parse initial firmware version and device ID
	firmwareVersionTmp, nextOffset, err := helpers.ParseHexSubstring(hexStr, 0, 1)
	if err != nil {
		return nil, helpers.WrapError(err)
	}
	// Divide by 10 to convert to float64 and shift decimal place
	firmwareVersion := float64(firmwareVersionTmp) / 10.0

	// Parse device ID
	deviceID, nextOffset, err := helpers.ParseHexSubstring(hexStr, nextOffset, 7)
	if err != nil {
		return nil, helpers.WrapError(err)
	}

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
			pkg, err := parseParkingPackage53(hexStr, timestamp, nextOffset1)
			if err != nil {
				return nil, err
			}
			parkingPackages = append(parkingPackages, pkg)
			nextOffset1 = pkg["nextOffset"].(int)

		case 6:
			keepAliveAmount++
			pkg, err := parseKeepAlivePackage53(hexStr, timestamp, nextOffset1, keepAliveAmount)
			if err != nil {
				return nil, err
			}
			keepAlivePackages = append(keepAlivePackages, pkg)
			nextOffset1 = pkg["nextOffset"].(int)

		case 10:
			settingsAmount++
			pkg, err := parseSettingsPackage53(hexStr, timestamp, nextOffset1)
			if err != nil {
				return nil, err
			}
			settingsPackages = append(settingsPackages, pkg)
			nextOffset1 = pkg["nextOffset"].(int)
		}

		nextOffset = nextOffset1
	}

	myMap["firmware_version"] = firmwareVersion
	myMap["device_id"] = deviceID
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
func parseParkingPackage53(hexStr string, timestamp, offset int) (map[string]any, error) {
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
func parseKeepAlivePackage53(hexStr string, timestamp, offset, keepAliveAmount int) (map[string]any, error) {
	pkg := map[string]any{"timestamp": timestamp}
	var err error
	var nextOffset int

	// Parsing each field

	// Continue parsing additional fields...
	if pkg["idle_voltage"], nextOffset, err = helpers.ParseHexSubstring(hexStr, offset, 1); err != nil {
		return nil, helpers.WrapError(err)
	}
	pkg["idle_voltage"] = pkg["idle_voltage"].(int) / 16

	if pkg["battery_percentage"], nextOffset, err = helpers.ParseHexSubstring(hexStr, nextOffset, 1); err != nil {
		return nil, helpers.WrapError(err)
	}

	if pkg["current"], nextOffset, err = helpers.ParseHexSubstring(hexStr, nextOffset, 1); err != nil {
		return nil, helpers.WrapError(err)
	}

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

	if pkg["socket_error"], nextOffset, err = helpers.ParseHexSubstring(hexStr, nextOffset, 1); err != nil {
		return nil, helpers.WrapError(err)
	}

	if pkg["t3324"], nextOffset, err = helpers.ParseHexSubstring(hexStr, nextOffset, 1); err != nil {
		return nil, helpers.WrapError(err)
	}

	if pkg["t3412"], nextOffset, err = helpers.ParseHexSubstring(hexStr, nextOffset, 1); err != nil {
		return nil, helpers.WrapError(err)
	}

	if pkg["time_sync_rand_byte"], nextOffset, err = helpers.ParseHexSubstring(hexStr, nextOffset, 1); err != nil {
		return nil, helpers.WrapError(err)
	}

	if pkg["time_sync_rand_byte"].(int) != 0 || keepAliveAmount == 1 {
		var nextOffset3 int
		if pkg["time_sync_current_unix_time"], nextOffset3, err = helpers.ParseHexSubstring(hexStr, nextOffset, 4); err != nil {
			return nil, helpers.WrapError(err)
		}
		nextOffset = nextOffset3
	}

	pkg["nextOffset"] = nextOffset
	return pkg, nil
}

// Parses the Settings Package
func parseSettingsPackage53(hexStr string, timestamp, offset int) (map[string]any, error) {
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

	if pkg["radar_trail_cal_lo_th"], nextOffset, err = helpers.ParseHexSubstring(hexStr, nextOffset, 1); err != nil {
		return nil, helpers.WrapError(err)
	}
	pkg["radar_trail_cal_lo_th"] = pkg["radar_trail_cal_lo_th"].(int) * 256

	if pkg["radar_trail_cal_hi_th"], nextOffset, err = helpers.ParseHexSubstring(hexStr, nextOffset, 1); err != nil {
		return nil, helpers.WrapError(err)
	}
	pkg["radar_trail_cal_hi_th"] = pkg["radar_trail_cal_hi_th"].(int) * 256

	if pkg["radar_trail_uncal_lo_th"], nextOffset, err = helpers.ParseHexSubstring(hexStr, nextOffset, 1); err != nil {
		return nil, helpers.WrapError(err)
	}
	pkg["radar_trail_uncal_lo_th"] = pkg["radar_trail_uncal_lo_th"].(int) * 256

	if pkg["radar_trail_uncal_hi_th"], nextOffset, err = helpers.ParseHexSubstring(hexStr, nextOffset, 1); err != nil {
		return nil, helpers.WrapError(err)
	}
	pkg["radar_trail_uncal_hi_th"] = pkg["radar_trail_uncal_hi_th"].(int) * 256

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

	if pkg["ten_x_deep_sleep_time"], nextOffset, err = helpers.ParseHexSubstringBigInt(hexStr, nextOffset, 20); err != nil {
		return nil, helpers.WrapError(err)
	}
	// Divide the big.Int value by 4
	deepSleepTime := pkg["ten_x_deep_sleep_time"].(*big.Int)
	divisionFactor := big.NewInt(4)
	pkg["ten_x_deep_sleep_time"] = new(big.Int).Div(deepSleepTime, divisionFactor)

	if pkg["ten_x_action_before"], nextOffset, err = helpers.ParseHexSubstringBigInt(hexStr, nextOffset, 10); err != nil {
		return nil, helpers.WrapError(err)
	}

	if pkg["ten_x_action_after"], nextOffset, err = helpers.ParseHexSubstringBigInt(hexStr, nextOffset, 10); err != nil {
		return nil, helpers.WrapError(err)
	}

	if pkg["nb_iot_udp_ip_1"], nextOffset, err = helpers.ParseHexSubstring(hexStr, nextOffset, 1); err != nil {
		return nil, helpers.WrapError(err)
	}
	if pkg["nb_iot_udp_ip_2"], nextOffset, err = helpers.ParseHexSubstring(hexStr, nextOffset, 1); err != nil {
		return nil, helpers.WrapError(err)
	}
	if pkg["nb_iot_udp_ip_3"], nextOffset, err = helpers.ParseHexSubstring(hexStr, nextOffset, 1); err != nil {
		return nil, helpers.WrapError(err)
	}
	if pkg["nb_iot_udp_ip_4"], nextOffset, err = helpers.ParseHexSubstring(hexStr, nextOffset, 1); err != nil {
		return nil, helpers.WrapError(err)
	}

	// Combine the IP parts into a single string
	pkg["nb_iot_udp_ip"] = fmt.Sprintf("%d.%d.%d.%d",
		pkg["nb_iot_udp_ip_1"].(int),
		pkg["nb_iot_udp_ip_2"].(int),
		pkg["nb_iot_udp_ip_3"].(int),
		pkg["nb_iot_udp_ip_4"].(int))

	if pkg["nb_iot_udp_port"], nextOffset, err = helpers.ParseHexSubstring(hexStr, nextOffset, 2); err != nil {
		return nil, helpers.WrapError(err)
	}

	// Parse the APN length
	var apnLength int
	if apnLength, nextOffset, err = helpers.ParseHexSubstring(hexStr, nextOffset, 1); err != nil {
		return nil, helpers.WrapError(err)
	}
	pkg["nb_iot_apn_length"] = apnLength

	if pkg["nb_iot_apn"], nextOffset, err = helpers.ParseHexToASCIIString(hexStr, nextOffset, pkg["nb_iot_apn_length"].(int)); err != nil {
		return nil, helpers.WrapError(err)
	}

	if pkg["nb_iot_imsi"], nextOffset, err = helpers.ParseHexSubstringBigInt(hexStr, nextOffset, 7); err != nil {
		return nil, helpers.WrapError(err)
	}

	pkg["nextOffset"] = nextOffset
	return pkg, nil
}
