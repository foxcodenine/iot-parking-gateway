import { useLocalStorage } from '@vueuse/core';
import { ref, computed, reactive } from 'vue';
import { defineStore } from 'pinia';

// ---------------------------------------------------------------------

export const useAppStore = defineStore("appStore", () => {

    // - State ---------------------------------------------------------
    const appVariables = useLocalStorage("appVariables", {
        appUrl: import.meta.env.VITE_VUE_ENV === 'production' ? window.location.origin : import.meta.env.VITE_APP_URL,      
    });

    // Adjust googleApiKey based on environment and availability
    if (import.meta.env.VITE_VUE_ENV === 'production' && GO_GOOGLE_API_KEY) {
        appVariables.value.googleApiKey = GO_GOOGLE_API_KEY; // Use value to properly reference the reactive object
        

    } else if (import.meta.env.VITE_VUE_ENV === 'development') {
        appVariables.value.googleApiKey = import.meta.env.VITE_GOOGLE_API_KEY;
    }



    // - Getters -------------------------------------------------------
    const getAppVariables = () => appVariables;

    // - Actions -------------------------------------------------------



    function resetAppState() {

    }

    // - Expose --------------------------------------------------------
    return {
        resetAppState,
        getAppVariables,
        appVariables: appVariables.value, 
    };
});
