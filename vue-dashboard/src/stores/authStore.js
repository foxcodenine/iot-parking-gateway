import { ref, computed, watch, toRaw } from 'vue';
import { defineStore } from 'pinia';
import { useLocalStorage, useSessionStorage } from '@vueuse/core';
import { useJwtComposable } from '@/composables/useJwtComposable';
import { useRouter } from 'vue-router';

export const useAuthStore = defineStore("authStore", () => {

    // - State ---------------------------------------------------------
    
    const jwt = ref(null);  // Start with a non-persistent state
    const jwtLocalStorage = useLocalStorage("jwt", null); // Local storage for when remember me is true
    const jwtSessionStorage = useSessionStorage("jwt", null); // Session storage for when remember me is false

    const rememberMe = useLocalStorage("rememberMe", false);
    
    const redirectTo = ref("mapView");    
    
    // - Getters -------------------------------------------------------

    const isAuthenticated = computed(() => {
        return Boolean(jwt.value) || Boolean(jwtLocalStorage.value) || Boolean(jwtSessionStorage.value);           
    });

    const getJwt = computed(()=>{
        return jwt.value || jwtLocalStorage.value || jwtSessionStorage.value;
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
    function reset() {
        jwt.value = null;
        jwtLocalStorage.value = null;  
        jwtSessionStorage.value = null; 
    }

    function setJwt(token) {
        
        jwt.value = token;
        if (rememberMe.value) {
            jwtLocalStorage.value = token;  
        } else {
            jwtSessionStorage.value = token;
        }  
    }

    function clearJwt() {
        jwt.value = null;
        jwtLocalStorage.value = null;  
        jwtSessionStorage.value = null;  
    }

    function setRedirectTo(payload) {
        redirectTo.value = payload.name;
    }

    function toggleRememberMe() {
        rememberMe.value = !rememberMe.value;
    }   
    
    // - Expose --------------------------------------------------------

    return {
        getJwt,
        reset,
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
