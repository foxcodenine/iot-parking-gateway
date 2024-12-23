import { useLocalStorage, useSessionStorage } from '@vueuse/core';
import { ref, computed, reactive } from 'vue';
import { defineStore } from 'pinia';
import {  decryptAES } from '@/utils/cryptoUtils';
import { useAuthStore } from './authStore';
import { useDashboardStore } from './dashboardStore';
import axios from '@/axios';

export const useAppStore = defineStore("appStore", () => {

    const pageScrollDisabled = ref(false); 

    const appUrl = ref(import.meta.env.VITE_VUE_ENV === 'production' ? window.location.origin : import.meta.env.VITE_APP_URL);
    
    let appSettingsStorage = useAuthStore().getRemeberMe ? useLocalStorage('appSettings', null) : useSessionStorage('appSettings', null);

    const appSettings = computed({
        get: () => appSettingsStorage.value,
        set: (val) => appSettingsStorage.value = val
    });

    const googleApiKey = ref(null);

    const authUserStorage = useAuthStore().getRemeberMe ? useLocalStorage('authUser', null) : useSessionStorage('authUser', null);
    const authUser = computed({
        get: () => authUserStorage.value,
        set: (val) => authUserStorage.value = val
    });
    

    // - Getters -------------------------------------------------------

    const getAppUrl = computed(() => appUrl.value);

    const getPageScrollDisabled = computed(()=> pageScrollDisabled.value);

    const getAppSettings = computed(() => {
        if (appSettings.value) {
            return JSON.parse(appSettings.value)
        }
        return null;
    });    

    const getAuthUser = computed(() => {
        if (authUser.value) {
            return JSON.parse(authUser.value)
        }
        return null;
    });    

    // - Actions -------------------------------------------------------

    async function updateSettings(payload) {
        useDashboardStore().setIsLoading(true);

        console.log(`${getAppUrl.value}/api/setting`)
        try {            
            return await axios.put(`${getAppUrl.value}/api/setting`, payload);
            
        } catch (error){
            console.error('! appStore.updateSettings !');
            throw error;  
        } finally {
            useDashboardStore().setIsLoading(false);
        }
    }

    function resetAppStore() {
        appSettings.value = null;
        authUser.value = null;
        googleApiKey.value = null;
    }

    function setPageScrollDisabled(val) {
        pageScrollDisabled.value = val;
    } 

    function setAppSettings(newSettings) {
        appSettings.value = JSON.stringify(newSettings)
    }

    function setAuthUser(newAuthUser) {
        authUser.value = JSON.stringify(newAuthUser)
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
        getAppSettings, setAppSettings,
        updateSettings,
        getAuthUser,
        setAuthUser,
    };
});
