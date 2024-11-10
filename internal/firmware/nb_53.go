package firmware

import (
	"fmt"
	"math/big"

	"github.com/foxcodenine/iot-parking-gateway/internal/helpers"
)

func NB_53(hexStr string) (map[string]any, error) {
	myMap := map[string]any{}

	// Parse firmware version
	firmwareVersion, nextOffset, err := helpers.ParseHexSubstring(hexStr, 0, 1)
	if err != nil {
		return nil, helpers.WrapError(err)
	}

	// Parse device ID
	deviceID, nextOffset, err := helpers.ParseHexSubstring(hexStr, nextOffset, 7)
	if err != nil {
		return nil, helpers.WrapError(err)
	}

	pkgAmount := 0
	parkingAmount := 0
	keepAliveAmount := 0
	settingsAmount := 0

	var parkingPackages []map[string]any
	var keepAlivePackages []map[string]any
	var settingsPackages []map[string]any

	for {

		if nextOffset*2 >= len(hexStr) {
			break
		}

		timestamp, nextOffset1, err := helpers.ParseHexSubstring(hexStr, nextOffset, 4)
		if err != nil {
			return nil, helpers.WrapError(err)
		}

		eventID, nextOffset1, err := helpers.ParseHexSubstring(hexStr, nextOffset1, 1)
		if err != nil {
			return nil, helpers.WrapError(err)
		}

		validEventIDs := []int{6, 10, 26}
		if !helpers.Contains(validEventIDs, eventID) {

			//  TODO:  add remainder
			return nil, fmt.Errorf("invalid event_id: %d remainder: ", eventID)
		}

		pkgAmount++

		// -- Parking --------------------------------------------------
		if eventID == 26 {

			pkg := map[string]any{"timestamp": timestamp}
			var nextOffset2 int

			parkingAmount++

			if pkg["peak_distance_cm"], nextOffset2, err = helpers.ParseHexSubstring(hexStr, nextOffset1, 1); err != nil {
				return nil, helpers.WrapError(err)
			}

			if pkg["vehicle_occupancy"], nextOffset2, err = helpers.ParseHexSubstring(hexStr, nextOffset2, 1); err != nil {
				return nil, helpers.WrapError(err)
			}

			if pkg["radar_cumulative"], nextOffset2, err = helpers.ParseHexSubstring(hexStr, nextOffset2, 1); err != nil {
				return nil, helpers.WrapError(err)
			}
			pkg["radar_cumulative"] = pkg["radar_cumulative"].(int) * 256

			if pkg["magnet_abs_total"], nextOffset2, err = helpers.ParseHexSubstring(hexStr, nextOffset2, 2); err != nil {
				return nil, helpers.WrapError(err)
			}

			if pkg["beacons_amount"], nextOffset2, err = helpers.ParseHexSubstring(hexStr, nextOffset2, 1); err != nil {
				return nil, helpers.WrapError(err)
			}

			// Parsing beacons array
			beacons := []map[string]any{}
			beaconsAmount := pkg["beacons_amount"].(int)

			for i := 1; i <= beaconsAmount; i++ {
				var nextOffset3 int
				beacon := map[string]any{"no": i}

				if beacon["major"], nextOffset3, err = helpers.ParseHexSubstring(hexStr, nextOffset2, 2); err != nil {
					return nil, helpers.WrapError(err)
				}
				if beacon["minor"], nextOffset3, err = helpers.ParseHexSubstring(hexStr, nextOffset3, 2); err != nil {
					return nil, helpers.WrapError(err)
				}
				if beacon["rssi"], nextOffset3, err = helpers.ParseHexSubstring(hexStr, nextOffset3, 1); err != nil {
					return nil, helpers.WrapError(err)
				}
				beacons = append(beacons, beacon)

				nextOffset2 = nextOffset3
			}

			pkg["beacons"] = beacons
			nextOffset1 = nextOffset2
			parkingPackages = append(parkingPackages, pkg)
		}

		// -- Keep Alive` ----------------------------------------------

		if eventID == 6 {
			pkg := map[string]any{"timestamp": timestamp}
			var nextOffset2 int

			keepAliveAmount++

			if pkg["idle_voltage"], nextOffset2, err = helpers.ParseHexSubstring(hexStr, nextOffset1, 1); err != nil {
				return nil, helpers.WrapError(err)
			}
			pkg["idle_voltage"] = pkg["idle_voltage"].(int) / 16

			if pkg["battery_percentage"], nextOffset2, err = helpers.ParseHexSubstring(hexStr, nextOffset2, 1); err != nil {
				return nil, helpers.WrapError(err)
			}

			if pkg["current"], nextOffset2, err = helpers.ParseHexSubstring(hexStr, nextOffset2, 1); err != nil {
				return nil, helpers.WrapError(err)
			}

			if pkg["reset_count"], nextOffset2, err = helpers.ParseHexSubstring(hexStr, nextOffset2, 1); err != nil {
				return nil, helpers.WrapError(err)
			}

			if pkg["manual_calibration"], nextOffset2, err = helpers.ParseHexSubstring(hexStr, nextOffset2, 1); err != nil {
				return nil, helpers.WrapError(err)
			}

			if pkg["temperature_min"], nextOffset2, err = helpers.ParseHexSubstring(hexStr, nextOffset2, 1); err != nil {
				return nil, helpers.WrapError(err)
			}

			if pkg["temperature_max"], nextOffset2, err = helpers.ParseHexSubstring(hexStr, nextOffset2, 1); err != nil {
				return nil, helpers.WrapError(err)
			}

			if pkg["radar_error"], nextOffset2, err = helpers.ParseHexSubstring(hexStr, nextOffset2, 1); err != nil {
				return nil, helpers.WrapError(err)
			}

			if pkg["mag_error"], nextOffset2, err = helpers.ParseHexSubstring(hexStr, nextOffset2, 1); err != nil {
				return nil, helpers.WrapError(err)
			}

			if pkg["tcve_error"], nextOffset2, err = helpers.ParseHexSubstring(hexStr, nextOffset2, 1); err != nil {
				return nil, helpers.WrapError(err)
			}

			if pkg["ble_security_issues"], nextOffset2, err = helpers.ParseHexSubstring(hexStr, nextOffset2, 1); err != nil {
				return nil, helpers.WrapError(err)
			}

			if pkg["radar_cumulative_total"], nextOffset2, err = helpers.ParseHexSubstring(hexStr, nextOffset2, 1); err != nil {
				return nil, helpers.WrapError(err)
			}
			pkg["radar_cumulative_total"] = pkg["radar_cumulative_total"].(int) * 256

			if pkg["mag_total"], nextOffset2, err = helpers.ParseHexSubstring(hexStr, nextOffset2, 2); err != nil {
				return nil, helpers.WrapError(err)
			}

			if pkg["network_registration_ok"], nextOffset2, err = helpers.ParseHexSubstring(hexStr, nextOffset2, 1); err != nil {
				return nil, helpers.WrapError(err)
			}

			if pkg["network_registration_nok"], nextOffset2, err = helpers.ParseHexSubstring(hexStr, nextOffset2, 1); err != nil {
				return nil, helpers.WrapError(err)
			}

			if pkg["rssi_average"], nextOffset2, err = helpers.ParseHexSubstring(hexStr, nextOffset2, 1); err != nil {
				return nil, helpers.WrapError(err)
			}

			if pkg["network_message_attempts"], nextOffset2, err = helpers.ParseHexSubstring(hexStr, nextOffset2, 1); err != nil {
				return nil, helpers.WrapError(err)
			}

			if pkg["network_ack_1ds"], nextOffset2, err = helpers.ParseHexSubstring(hexStr, nextOffset2, 1); err != nil {
				return nil, helpers.WrapError(err)
			}

			if pkg["network_1ack_ds"], nextOffset2, err = helpers.ParseHexSubstring(hexStr, nextOffset2, 1); err != nil {
				return nil, helpers.WrapError(err)
			}

			if pkg["network_1ack_1ds"], nextOffset2, err = helpers.ParseHexSubstring(hexStr, nextOffset2, 1); err != nil {
				return nil, helpers.WrapError(err)
			}

			if pkg["tcvr_deep_sleep_min"], nextOffset2, err = helpers.ParseHexSubstring(hexStr, nextOffset2, 2); err != nil {
				return nil, helpers.WrapError(err)
			}

			if pkg["tcvr_deep_sleep_max"], nextOffset2, err = helpers.ParseHexSubstring(hexStr, nextOffset2, 2); err != nil {
				return nil, helpers.WrapError(err)
			}

			if pkg["tcvr_deep_sleep_average"], nextOffset2, err = helpers.ParseHexSubstring(hexStr, nextOffset2, 2); err != nil {
				return nil, helpers.WrapError(err)
			}

			if pkg["settings_checksum"], nextOffset2, err = helpers.ParseHexSubstring(hexStr, nextOffset2, 1); err != nil {
				return nil, helpers.WrapError(err)
			}

			if pkg["socket_error"], nextOffset2, err = helpers.ParseHexSubstring(hexStr, nextOffset2, 1); err != nil {
				return nil, helpers.WrapError(err)
			}

			if pkg["t3324"], nextOffset2, err = helpers.ParseHexSubstring(hexStr, nextOffset2, 1); err != nil {
				return nil, helpers.WrapError(err)
			}

			if pkg["t3412"], nextOffset2, err = helpers.ParseHexSubstring(hexStr, nextOffset2, 1); err != nil {
				return nil, helpers.WrapError(err)
			}

			if pkg["time_sync_rand_byte"], nextOffset2, err = helpers.ParseHexSubstring(hexStr, nextOffset2, 1); err != nil {
				return nil, helpers.WrapError(err)
			}

			if pkg["time_sync_rand_byte"].(int) != 0 || keepAliveAmount == 1 {
				var nextOffset3 int
				if pkg["time_sync_rand_byte"], nextOffset3, err = helpers.ParseHexSubstring(hexStr, nextOffset2, 4); err != nil {
					return nil, helpers.WrapError(err)
				}
				nextOffset2 = nextOffset3
			}

			keepAlivePackages = append(keepAlivePackages, pkg)

			nextOffset1 = nextOffset2
		}

		// ---------------------------------------------------------------

		if eventID == 10 {

			pkg := map[string]any{"timestamp": timestamp}
			var nextOffset2 int
			settingsAmount++

			// Parsing settings fields
			if pkg["device_mode"], nextOffset2, err = helpers.ParseHexSubstring(hexStr, nextOffset1, 1); err != nil {
				return nil, helpers.WrapError(err)
			}
			if pkg["device_enable"], nextOffset2, err = helpers.ParseHexSubstring(hexStr, nextOffset2, 1); err != nil {
				return nil, helpers.WrapError(err)
			}

			if pkg["radar_car_cal_lo_th"], nextOffset2, err = helpers.ParseHexSubstring(hexStr, nextOffset2, 1); err != nil {
				return nil, helpers.WrapError(err)
			}
			pkg["radar_car_cal_lo_th"] = pkg["radar_car_cal_lo_th"].(int) * 256

			if pkg["radar_car_cal_hi_th"], nextOffset2, err = helpers.ParseHexSubstring(hexStr, nextOffset2, 1); err != nil {
				return nil, helpers.WrapError(err)
			}
			pkg["radar_car_cal_hi_th"] = pkg["radar_car_cal_hi_th"].(int) * 256

			if pkg["radar_car_uncal_lo_th"], nextOffset2, err = helpers.ParseHexSubstring(hexStr, nextOffset2, 1); err != nil {
				return nil, helpers.WrapError(err)
			}
			pkg["radar_car_uncal_lo_th"] = pkg["radar_car_uncal_lo_th"].(int) * 256

			if pkg["radar_car_uncal_hi_th"], nextOffset2, err = helpers.ParseHexSubstring(hexStr, nextOffset2, 1); err != nil {
				return nil, helpers.WrapError(err)
			}
			pkg["radar_car_uncal_hi_th"] = pkg["radar_car_uncal_hi_th"].(int) * 256

			if pkg["radar_car_delta_th"], nextOffset2, err = helpers.ParseHexSubstring(hexStr, nextOffset2, 1); err != nil {
				return nil, helpers.WrapError(err)
			}
			pkg["radar_car_delta_th"] = pkg["radar_car_delta_th"].(int) * 256

			if pkg["mag_car_lo"], nextOffset2, err = helpers.ParseHexSubstring(hexStr, nextOffset2, 2); err != nil {
				return nil, helpers.WrapError(err)
			}
			if pkg["mag_car_hi"], nextOffset2, err = helpers.ParseHexSubstring(hexStr, nextOffset2, 2); err != nil {
				return nil, helpers.WrapError(err)
			}

			if pkg["radar_trail_cal_lo_th"], nextOffset2, err = helpers.ParseHexSubstring(hexStr, nextOffset2, 1); err != nil {
				return nil, helpers.WrapError(err)
			}
			pkg["radar_trail_cal_lo_th"] = pkg["radar_trail_cal_lo_th"].(int) * 256

			if pkg["radar_trail_cal_hi_th"], nextOffset2, err = helpers.ParseHexSubstring(hexStr, nextOffset2, 1); err != nil {
				return nil, helpers.WrapError(err)
			}
			pkg["radar_trail_cal_hi_th"] = pkg["radar_trail_cal_hi_th"].(int) * 256

			if pkg["radar_trail_uncal_lo_th"], nextOffset2, err = helpers.ParseHexSubstring(hexStr, nextOffset2, 1); err != nil {
				return nil, helpers.WrapError(err)
			}
			pkg["radar_trail_uncal_lo_th"] = pkg["radar_trail_uncal_lo_th"].(int) * 256

			if pkg["radar_trail_uncal_hi_th"], nextOffset2, err = helpers.ParseHexSubstring(hexStr, nextOffset2, 1); err != nil {
				return nil, helpers.WrapError(err)
			}
			pkg["radar_trail_uncal_hi_th"] = pkg["radar_trail_uncal_hi_th"].(int) * 256

			if pkg["debug_period"], nextOffset2, err = helpers.ParseHexSubstring(hexStr, nextOffset2, 1); err != nil {
				return nil, helpers.WrapError(err)
			}
			if pkg["debug_mode"], nextOffset2, err = helpers.ParseHexSubstring(hexStr, nextOffset2, 1); err != nil {
				return nil, helpers.WrapError(err)
			}

			if pkg["logs_mode"], nextOffset2, err = helpers.ParseHexSubstring(hexStr, nextOffset2, 1); err != nil {
				return nil, helpers.WrapError(err)
			}
			if pkg["logs_amount"], nextOffset2, err = helpers.ParseHexSubstring(hexStr, nextOffset2, 1); err != nil {
				return nil, helpers.WrapError(err)
			}

			if pkg["maximum_registration_time"], nextOffset2, err = helpers.ParseHexSubstring(hexStr, nextOffset2, 1); err != nil {
				return nil, helpers.WrapError(err)
			}
			pkg["maximum_registration_time"] = pkg["maximum_registration_time"].(int) / 4

			if pkg["maximum_registration_attempts"], nextOffset2, err = helpers.ParseHexSubstring(hexStr, nextOffset2, 1); err != nil {
				return nil, helpers.WrapError(err)
			}
			if pkg["maximum_deep_sleep_time"], nextOffset2, err = helpers.ParseHexSubstring(hexStr, nextOffset2, 1); err != nil {
				return nil, helpers.WrapError(err)
			}
			pkg["maximum_deep_sleep_time"] = pkg["maximum_deep_sleep_time"].(int) / 2

			if pkg["ten_x_deep_sleep_time"], nextOffset2, err = helpers.ParseHexSubstringBigInt(hexStr, nextOffset2, 20); err != nil {
				return nil, helpers.WrapError(err)
			}
			// Divide the big.Int value by 4
			deepSleepTime := pkg["ten_x_deep_sleep_time"].(*big.Int)
			divisionFactor := big.NewInt(4)
			pkg["ten_x_deep_sleep_time"] = new(big.Int).Div(deepSleepTime, divisionFactor)

			if pkg["ten_x_action_before"], nextOffset2, err = helpers.ParseHexSubstringBigInt(hexStr, nextOffset2, 10); err != nil {
				return nil, helpers.WrapError(err)
			}

			if pkg["ten_x_action_after"], nextOffset2, err = helpers.ParseHexSubstringBigInt(hexStr, nextOffset2, 10); err != nil {
				return nil, helpers.WrapError(err)
			}

			if pkg["nb_iot_udp_ip"], nextOffset2, err = helpers.ParseHexSubstring(hexStr, nextOffset2, 4); err != nil {
				return nil, helpers.WrapError(err)
			}

			if pkg["nb_iot_udp_port"], nextOffset2, err = helpers.ParseHexSubstring(hexStr, nextOffset2, 2); err != nil {
				return nil, helpers.WrapError(err)
			}

			if pkg["nb_iot_apn_length"], nextOffset2, err = helpers.ParseHexSubstring(hexStr, nextOffset2, 1); err != nil {
				return nil, helpers.WrapError(err)
			}

			if pkg["nb_iot_apn"], nextOffset2, err = helpers.ParseHexSubstringBigInt(hexStr, nextOffset2, pkg["nb_iot_apn_length"].(int)); err != nil {
				return nil, helpers.WrapError(err)
			}

			if pkg["nb_iot_imsi"], nextOffset2, err = helpers.ParseHexSubstring(hexStr, nextOffset2, 7); err != nil {
				return nil, helpers.WrapError(err)
			}

			settingsPackages = append(settingsPackages, pkg)

			nextOffset1 = nextOffset2
		}

		nextOffset = nextOffset1
	}

	myMap["firmware_version"] = firmwareVersion
	myMap["device_id"] = deviceID
	myMap["parking_packages"] = parkingPackages
	myMap["keep_alive_packages"] = keepAlivePackages
	myMap["settings_packages"] = settingsPackages

	return myMap, nil
}
