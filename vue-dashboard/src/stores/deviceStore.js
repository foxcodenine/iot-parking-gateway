import { ref, computed } from 'vue'
import { defineStore } from 'pinia'
import axios from '@/axios';
import { useAppStore } from './appStore';
import { useDashboardStore } from './dashboardStore';

// ---------------------------------------------------------------------


export const useDeviceStore = defineStore("deviceStore", () => {

    // - State ---------------------------------------------------------
    const devicesList = ref({});
    const devicesFetched = ref(false);

    const filteredDevices = ref([]);

    // - Getters -------------------------------------------------------
    const getDevicesList = computed(() => {
         return devicesList.value;
    } );
    const getfilteredDevices = computed(() => {
         return filteredDevices.value;
    } );

    // - Actions -------------------------------------------------------

    function reset() {
    } 
    
    function setFilteredDevices(payload) {
        filteredDevices.value = payload;
    }

    async function fetchDevices() {
  
        useDashboardStore().setIsLoading(true);
        try {
            const response =  await axios.get(useAppStore().getAppUrl + '/api/device?map=true');
            if (response?.status == 200 && response.data?.devices) {
                devicesList.value = response.data.devices;
                devicesFetched.value = true;
                return devicesList.value;
            }
        } catch (error) {
            console.error('! deviceStore.fetchDevices !');
            throw error; 
        } finally {
            useDashboardStore().setIsLoading(false);
        }
    }

    async function createDevice({device_id, name, network_type, firmware_version, latitude, longitude, is_allowed, is_blocked, is_hidden}) {   

        useDashboardStore().setIsLoading(true);
     
        try {
            const payload = {
                device_id, 
                name, 
                network_type, 
                firmware_version, 
                latitude, 
                longitude, 
                is_allowed, 
                is_blocked, 
                is_hidden
            };  

            return await axios.post(useAppStore().getAppUrl + '/api/device', payload);
            
        } catch (error){
            console.error('! deviceStore.createDevice !');
            throw error;  
        } finally {
            useDashboardStore().setIsLoading(false);
        }
    }

    async function updateDevice({device_id, name, firmware_version,  latitude, longitude, is_occupied, is_allowed, is_blocked, is_hidden}) {
        useDashboardStore().setIsLoading(true);
        try {
            const payload = {                
                name,
                firmware_version,
                latitude,
                longitude,
                is_occupied,
                is_allowed,
                is_blocked,
                is_hidden,
            };
            
            return await axios.put(`${useAppStore().getAppUrl}/api/device/${device_id}`, payload);
            
        } catch (error){
            console.error('! deviceStore.updateDevice !');
            throw error;  
        } finally {
            useDashboardStore().setIsLoading(false);
        }
    }

    async function deleteDevice({device_id}) {
        useDashboardStore().setIsLoading(true);
        try {            
            return await axios.delete(`${useAppStore().getAppUrl}/api/device/${device_id}`);            
        } catch (error){
            console.error('! deviceStore.deleteDevice !');
            throw error;  
        } finally {
            useDashboardStore().setIsLoading(false);
        }
    }

    function pushDeviceToList(device) {
        
        if (device && device.device_id) {           
            devicesList.value[device.device_id] = device;
        }  
    }

    function updateDeviceInList (device) {
        if (device && device.device_id) {
            devicesList.value[device.device_id] = device;
        }
    }

    function removeDeviceFromList (deviceID) {
        if (deviceID) {
            delete devicesList.value[deviceID]
        }
    }

    function onParkingEvent(payload) {
  
        if (!devicesList.value[payload.device_id]) return
       devicesList.value[payload.device_id].is_occupied = payload.is_occupied;
       devicesList.value[payload.device_id].happened_at = payload.happened_at;
       devicesList.value[payload.device_id].firmware_version = payload.firmware_version;
       devicesList.value[payload.device_id].beacons = payload.beacons;
    }

    // - Expose --------------------------------------------------------

    return {
        reset,
        getDevicesList,
        fetchDevices,
        updateDevice,
        updateDeviceInList,
        deleteDevice,
        removeDeviceFromList,
        createDevice,
        pushDeviceToList,
        onParkingEvent,
        getfilteredDevices,
        setFilteredDevices
    }
});