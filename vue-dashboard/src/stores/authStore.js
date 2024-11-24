import { ref, computed, watch, toRaw } from 'vue';
import { defineStore } from 'pinia';
import { useLocalStorage } from '@vueuse/core';
import { useJwtComposable } from '@/composables/useJwtComposable';
import { useRouter } from 'vue-router';

export const useAuthStore = defineStore("authStore", () => {

    // - State ---------------------------------------------------------
    
    const jwt = ref(null);  // Start with a non-persistent state

    const jwtLocalStorage = useLocalStorage("jwt", null);  // Local storage management

    const rememberMe = useLocalStorage("remember-me", false);   
    
    const redirectTo = ref("viewHome");    
    
    // - Getters -------------------------------------------------------

    const isAuthenticated = computed(() => {
        if (Boolean(jwt.value) || Boolean(jwtLocalStorage.value)) {
            return true;
        }
        return false;
    });

    const getJwt = computed(()=>{
        return jwt.value || jwtLocalStorage.value;
    });

    const getRedirectTo = computed(()=>{
        return redirectTo.value;
    })

    const getRemeberMe = computed(()=>{
        return rememberMe.value;
    });

    const getUserTokenData = computed(()=>{
        return useJwtComposable().parseJwt();
    });
    

    // - Actions -------------------------------------------------------
    function reset() {
    }

    function setJwt(token) {
        
        jwt.value = token;
        if (rememberMe.value) {
            jwtLocalStorage.value = token;  // Save to local storage if "remember me" is enabled
        }  
    }

    function clearJwt() {
        jwt.value = null;
        jwtLocalStorage.value = null;  // Clear from local storage

        // useRouter().push({ name: 'viewLogin' });
    }

    function setRedirectTo(payload) {
        redirectTo.value = payload.name;
    }

    function toggleRememberMe(val) {
        rememberMe.value = !rememberMe.value;
    }

    

   
    
    // - Expose --------------------------------------------------------

    return {
        getJwt,
        isAuthenticated,
        setJwt,
        clearJwt,
        getRedirectTo,
        setRedirectTo,
        getRemeberMe, 
        toggleRememberMe, 
        getUserTokenData,
    }
});
