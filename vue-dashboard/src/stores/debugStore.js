import { ref, computed } from 'vue'
import { defineStore } from 'pinia'
import axios from '@/axios';
import { useAppStore } from './appStore';
import { useDashboardStore } from './dashboardStore';


// ---------------------------------------------------------------------


export const useDebugStore = defineStore("debugStore", () => {


    // - State ---------------------------------------------------------


    const selectedDeviceID = ref('')
    const fromDateTS = ref(0)
    const toDateTS = ref(0)

    const activityLogs = ref([]);


    // - Getters -------------------------------------------------------

    const getSelectedDeviceID = computed(() => {
        return selectedDeviceID.value;
    });
    const getFromDateTS = computed(() => {
        return fromDateTS.value;
    });
    const getToDateTS = computed(() => {
        return toDateTS.value;
    });


    // - Actions -------------------------------------------------------

    function setSelectedDeviceID(payload) {
        selectedDeviceID.value = payload;
    }
    function setFromDateTS(payload) {
        fromDateTS.value = payload;      
    }
    function setToDateTS(payload) {
        toDateTS.value = payload;
    }

    async function fetchActivityLogs() {  
        useDashboardStore().setIsLoading(true);

        const fromDate = Math.floor(fromDateTS.value / 1000);
        const toDate = Math.ceil(toDateTS.value / 1000);


        try {
            const response =  await axios.get(useAppStore().getAppUrl + `/api/activity-logs/${selectedDeviceID.value}?from_date=${fromDate}&to_date=${toDate}`);
            console.log(response)
            if (response?.status == 200 && response.data?.activity_logs) {
                activityLogs.value = response.data.activity_logs;
             
                return activityLogs.value;
            }

        } catch (error) {
            console.error('! debugStore.fetchActivityLogs !');
            throw error; 

        } finally {
            useDashboardStore().setIsLoading(false);
        }
    }


    // - Expose --------------------------------------------------------

    return {
        getSelectedDeviceID,
        setSelectedDeviceID,
        getToDateTS,
        getFromDateTS,
        setFromDateTS,
        setToDateTS,
        fetchActivityLogs,
    }
});