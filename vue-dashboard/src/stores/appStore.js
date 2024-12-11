import { useLocalStorage } from '@vueuse/core';
import { ref, computed, reactive } from 'vue';
import { defineStore } from 'pinia';
import { decryptString, encryptString } from '@/utils/cryptoUtils';

export const useAppStore = defineStore("appStore", () => {

    const appUrl = ref(import.meta.env.VITE_VUE_ENV === 'production' ? window.location.origin : import.meta.env.VITE_APP_URL);
    const defaultLatitude = ref(import.meta.env.VITE_VUE_ENV === 'production' ? GO_DEFAULT_LATITUDE : import.meta.env.VITE_DEFAULT_LATITUDE);
    const defaultLongitude = ref(import.meta.env.VITE_VUE_ENV === 'production' ? GO_DEFAULT_LONGITUDE : import.meta.env.VITE_DEFAULT_LONGITUDE);

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


    const appSettings = useLocalStorage('appSettings', null)

    const whitelistBlacklistMode = ref("black"); 
    
    const pageScrollDisabled = useLocalStorage("googleApiKey", null); 

    // - Getters -------------------------------------------------------

    const getAppUrl = computed(() => appUrl.value);

    const getGoogleApiKey = computed(() => decryptString(googleApiKey.value));

    const getWhitelistBlacklistMode = computed(() =>{
        return whitelistBlacklistMode.value;
    });

    const getDefaultLatitude = computed(()=>{         
        return defaultLatitude.value;
     });
    const getDefaultLongitude = computed(()=>{ return defaultLongitude.value });

    const getPageScrollDisabled = computed(()=> pageScrollDisabled.value);

    const getAppSettings = computed(() => appSettings.value)

    

    // - Actions -------------------------------------------------------

    function resetAppState() {
        appUrlLocalStorage.value = null;
        googleApiKeyLocalStorage.value = null;
    }

    function setPageScrollDisabled(val) {
        pageScrollDisabled.value = val;
    } 

    function setAppSettings(newSettings) {
        appSettings.value = newSettings
      }

    // - Expose --------------------------------------------------------
    
    return {
        resetAppState,
        getAppUrl,
        getGoogleApiKey,
        getWhitelistBlacklistMode,
        getDefaultLatitude,
        getDefaultLongitude,
        getPageScrollDisabled,
        setPageScrollDisabled,
        getAppSettings, setAppSettings 
    };
});
