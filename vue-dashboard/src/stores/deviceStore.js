import { ref, computed } from 'vue'
import { defineStore } from 'pinia'
import axios from '@/axios';
import { useAppStore } from './appStore';
import { useDashboardStore } from './dashboardStore';

// ---------------------------------------------------------------------


export const useDeviceStore = defineStore("deviceStore", () => {

    // - State ---------------------------------------------------------
    const devicesList = ref([]);

    // - Getters -------------------------------------------------------
    const getDevicesList = computed(() => devicesList.value );

    // - Actions -------------------------------------------------------

    function reset() {
    }    

    async function fetchDevices() {
        useDashboardStore().setIsLoading(true);
        try {
            const response =  await axios.get(useAppStore().getAppUrl + '/api/device');
            if (response.status == 200 && response.data?.devices) {
                devicesList.value = response.data.devices;
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
            devicesList.value.push(device);
        }  
    }

    function updateDeviceInList (device) {
        if (device && device.device_id) {
            const index = devicesList.value.findIndex(d => d.device_id == device.device_id)
            if (index == -1) { return }
            devicesList.value[index] = {...device};
        }
    }

    function removeDeviceFromList (deviceID) {
        if (deviceID) {
            const index = devicesList.value.findIndex(d => d.device_id == deviceID)
            if (index == -1) { return }

            devicesList.value.splice(index, 1)
        }
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
    }
});