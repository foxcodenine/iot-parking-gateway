import { useLocalStorage, useSessionStorage } from '@vueuse/core';
import { ref, computed, reactive } from 'vue';
import { defineStore } from 'pinia';
import {  decryptAES } from '@/utils/cryptoUtils';
import { useAuthStore } from './authStore';

export const useAppStore = defineStore("appStore", () => {

    const pageScrollDisabled = ref(false); 

    const appUrl = ref(import.meta.env.VITE_VUE_ENV === 'production' ? window.location.origin : import.meta.env.VITE_APP_URL);
    
    let appSettingsStorage = useAuthStore().getRemeberMe ? useLocalStorage('appSettings', null) : useSessionStorage('appSettings', null);

    const appSettings = computed({
        get: () => appSettingsStorage.value,
        set: (val) => appSettingsStorage.value = val
    });

    const googleApiKey = ref(null);
    

    // - Getters -------------------------------------------------------

    const getAppUrl = computed(() => appUrl.value);

    const getPageScrollDisabled = computed(()=> pageScrollDisabled.value);

    const getAppSettings = computed(() => {
        if (appSettings.value) {
            return JSON.parse(appSettings.value)
        }
        return null;
    });    

    // - Actions -------------------------------------------------------

    function resetAppStore() {
        appSettings.value = null;
        googleApiKey.value = null;
    }

    function setPageScrollDisabled(val) {
        pageScrollDisabled.value = val;
    } 

    function setAppSettings(newSettings) {
        appSettings.value = JSON.stringify(newSettings)
    }

    async function getGoogleApiKey() {
        if (!googleApiKey.value) {
            const s = JSON.parse(appSettings.value); 
            googleApiKey.value = await decryptAES(s.google_api_key);
        }
        return googleApiKey.value;
    };

    // - Expose --------------------------------------------------------
    
    return {
        resetAppStore,
        getAppUrl,
        getGoogleApiKey,
        getPageScrollDisabled,
        setPageScrollDisabled,
        getAppSettings, setAppSettings 
    };
});
