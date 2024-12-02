import { useLocalStorage } from '@vueuse/core';
import { ref, computed } from 'vue';
import { defineStore } from 'pinia';
import { decryptString, encryptString } from '@/utils/cryptoUtils';

export const useAppStore = defineStore("appStore", () => {

    const appUrl = ref(import.meta.env.VITE_VUE_ENV === 'production' ? window.location.origin : import.meta.env.VITE_APP_URL);
   

    // Using local storage to manage the Google API key
    const googleApiKey = ref(null);
    const googleApiKeyLocalStorage = useLocalStorage("googleApiKey", null); 

    // Initialize Google API Key if not already set in local storage
    if (googleApiKeyLocalStorage.value === null) {
        const apiKey = encryptString(import.meta.env.VITE_VUE_ENV === 'production' ? decryptString(GO_GOOGLE_API_KEY): import.meta.env.VITE_GOOGLE_API_KEY); 
        googleApiKeyLocalStorage.value = apiKey;  // Store the API key in local storage
        googleApiKey.value = apiKey;  // Set the reactive reference
        GO_GOOGLE_API_KEY = null;
    } else {
        googleApiKey.value = googleApiKeyLocalStorage.value;  // Use the stored key
        GO_GOOGLE_API_KEY = null;
    }

    // - Getters -------------------------------------------------------

    const getAppUrl = computed(() => appUrl.value);
    const getGoogleApiKey = computed(() => decryptString(googleApiKey.value));

    // - Actions -------------------------------------------------------

    function resetAppState() {
        appUrlLocalStorage.value = null;
        googleApiKeyLocalStorage.value = null;
    }

    // - Expose --------------------------------------------------------
    
    return {
        resetAppState,
        getAppUrl,
        getGoogleApiKey,
    };
});
