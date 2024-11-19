import { ref, computed } from 'vue'
import { defineStore } from 'pinia'
import { decryptEnv } from '@/utils/cryptoUtils';

// ---------------------------------------------------------------------


export const useEnvStore = defineStore("envStore", () => {


    // - State ---------------------------------------------------------
    const env = ref({})

    // - Getters -------------------------------------------------------
    const getEnv = computed(()=>{
        return env.value;
    })

    // - Actions -------------------------------------------------------
    function reset() {
        env.value = {};
    }

    async function loadEnvironmentVariables() {
        //   ${window.location.origin}
        const response = await fetch(`${window.location.origin}` + '/env');
  
        const data = await response.json();
    
        for (let key in data) {
            env.value[key] = decryptEnv(data[key]);
        }
    
        return env.value;
    }   


    // - Expose --------------------------------------------------------

    return {
        reset,
        getEnv,
        loadEnvironmentVariables,
    }
});