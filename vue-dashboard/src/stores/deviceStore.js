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

    // - Expose --------------------------------------------------------

    return {
        reset,
        getDevicesList,
        fetchDevices,
    }
});