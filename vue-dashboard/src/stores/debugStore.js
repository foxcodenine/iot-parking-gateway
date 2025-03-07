import { ref, computed } from 'vue'
import { defineStore } from 'pinia'
import axios from '@/axios';
import { useAppStore } from './appStore';
import { useDashboardStore } from './dashboardStore';
import { useDeviceStore } from './deviceStore';
import { formatToLocalDateTime } from '@/utils/dateTimeUtils';


// ---------------------------------------------------------------------

export const useDebugStore = defineStore("debugStore", () => {


    // - State ---------------------------------------------------------

    const selectedDeviceID = ref(null);
    const selectedDeviceName = ref('');
    const fromDateTS = ref(0);
    const toDateTS = ref(0);

    const activityLogs = ref([]);   

    const activityLogsDetails = ref({
        fromDate: '',
        toDate: '',
        id: null,
        name: '',
    });

    const keepaliveLogs = ref([]);

    const keepaliveLogDetails = ref({
        fromDate: '',
        toDate: '',
        id: null,
        name: '',
    });

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

    const getKeepaliveLogs = computed(() => {
        return keepaliveLogs.value;
    })

    const getActivityLogsDetails = computed(() => {
        return activityLogsDetails.value;
    });

    const getKeepaliveLogDetails = computed(() => {
        return keepaliveLogDetails.value;
    });

    // - Actions -------------------------------------------------------

    function setSelectedDeviceID(payload) {
        selectedDeviceID.value = payload;
    }

    function setSelectedDeviceName() {
        selectedDeviceName.value = useDeviceStore().getDevicesList[selectedDeviceID.value]?.name;
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

        setSelectedDeviceName();
   
        try {          
   
            // Fetch data from the server
            const response = await axios.get(`${useAppStore().getAppUrl}/api/activity-logs/${selectedDeviceID.value}?from_date=${fromDate}&to_date=${toDate}`);  
            activityLogs.value = response.data.activity_logs;  

            keepaliveLogDetails.value = {
                fromDate: formatToLocalDateTime(fromDateTS.value),
                toDate: formatToLocalDateTime(toDateTS.value),
                id: selectedDeviceID.value,
                name: selectedDeviceName.value,
            };

            return response;
   
        } catch (error) {
            console.error('! Error in fetchActivityLogs !', error.message);
            throw error;  // Re-throw the error to be handled by the caller or error boundary

        } finally {
            useDashboardStore().setIsLoading(false);
        }
    }    


    async function fetchKeepaliveLogs() {
        useDashboardStore().setIsLoading(true);
    
        const fromDate = Math.floor(fromDateTS.value / 1000);
        const toDate = Math.ceil(toDateTS.value / 1000);

        setSelectedDeviceName();

        try {          
   
            // Fetch data from the server
            const response = await axios.get(`${useAppStore().getAppUrl}/api/keepalive-logs/${selectedDeviceID.value}?from_date=${fromDate}&to_date=${toDate}`);  
            keepaliveLogs.value = response.data.keepalive_logs;  

            keepaliveLogDetails.value = {
                fromDate: formatToLocalDateTime(fromDateTS.value),
                toDate: formatToLocalDateTime(toDateTS.value),
                id: selectedDeviceID.value,
                name: selectedDeviceName.value,
            };
            return response;
   
        } catch (error) {
            console.error('! Error in fetchKeepaliveLogs !', error.message);
            throw error;  // Re-throw the error to be handled by the caller or error boundary

        } finally {
            useDashboardStore().setIsLoading(false);
        }
    }


    // - Expose --------------------------------------------------------

    return {
        getSelectedDeviceID,
        getToDateTS,
        getFromDateTS,
        getKeepaliveLogs,
        getKeepaliveLogDetails,
        getActivityLogsDetails,
        
        setSelectedDeviceID, 
        setFromDateTS,
        setToDateTS,

        fetchActivityLogs,
        fetchKeepaliveLogs,
    }
});