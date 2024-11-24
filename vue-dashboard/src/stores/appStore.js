import { ref, computed } from 'vue'
import { defineStore } from 'pinia'

// ---------------------------------------------------------------------


export const useAppStore = defineStore("appStore", () => {


    // - State ---------------------------------------------------------
    const appUrl = ref( import.meta.env.VITE_APP_URL)
    if (import.meta.env.VITE_VUE_ENV == 'production') {
        appUrl.value = `${window.location.origin}`;
    } 

    // - Getters -------------------------------------------------------
    const getAppUrl = computed(()=>{
        return appUrl.value
    });


    // - Actions -------------------------------------------------------
    function resetAppState() {
    }

    // - Expose --------------------------------------------------------

    return {
        resetAppState,
        getAppUrl,
    }
});