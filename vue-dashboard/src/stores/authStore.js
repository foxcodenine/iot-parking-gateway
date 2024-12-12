import { ref, computed, watch, toRaw } from 'vue';
import { defineStore } from 'pinia';
import { useLocalStorage, useSessionStorage } from '@vueuse/core';
import { useJwtComposable } from '@/composables/useJwtComposable';
import { useRouter } from 'vue-router';

export const useAuthStore = defineStore("authStore", () => {

    // - State ---------------------------------------------------------
    const rememberMe = useLocalStorage("rememberMe", false);

    let jwtStorage = rememberMe.value ? useLocalStorage('jwt', null) : useSessionStorage('jwt', null);;

    const jwt = computed({
        get: () => jwtStorage.value,
        set: (val) => jwtStorage.value = val
    });   
    
    const redirectTo = ref("mapView");    
    
    // - Getters -------------------------------------------------------

    const isAuthenticated = computed(() => {
        return Boolean(jwt.value);           
    });

    const getJwt = computed(()=>{
        return jwt.value;
    });

    const getRemeberMe = computed(()=>{
        return rememberMe.value;
    });

    const getRedirectTo = computed(()=>{
        return redirectTo.value;
    });

    const getUserTokenData = computed(()=>{
        return useJwtComposable().parseJwt();
    });    

    const getUserAccessLevel = computed(()=>{
        return useJwtComposable().parseJwt().access_level;
    });    

    // - Actions -------------------------------------------------------

    function setJwt(token) {        
        jwt.value = token;
    }

    function clearJwt() {
        jwt.value = null;
    }

    function resetAuthStore() {
        jwt.value = null;
    }

    function setRedirectTo(payload) {
        redirectTo.value = payload.name;
    }

    function toggleRememberMe() {
        rememberMe.value = !rememberMe.value;
    }   
    
    // - Expose --------------------------------------------------------

    return {
        rememberMe,
        getJwt,
        resetAuthStore,
        isAuthenticated,
        setJwt,
        clearJwt,
        getRedirectTo,
        setRedirectTo,
        getRemeberMe, 
        toggleRememberMe, 
        getUserTokenData,
        getUserAccessLevel,
    }
});
