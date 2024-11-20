import { ref, computed } from 'vue'
import { defineStore } from 'pinia'
import { decryptEnv } from '@/utils/cryptoUtils';

// ---------------------------------------------------------------------


export const useEnvStore = defineStore("envStore", () => {


    // - State ---------------------------------------------------------
    const env = ref({})

    // - Getters -------------------------------------------------------
    const getEnv = computed(() => {
        return env.value;
    })

    // - Actions -------------------------------------------------------
    function reset() {
        env.value = {};
    }

    async function loadEnvironmentVariables() {

        window.runtimeConfig = window.runtimeConfig || {}; // Initialize if not already set

        // Check if the app is running in a production environment
        if (import.meta.env.VITE_VUE_ENV === "production") {
            // Dynamically set runtime configuration for the production environment
            window.runtimeConfig.VITE_APP_URL = '${window.location.origin}'; // Set the app URL
        }

        // Fetch environment variables from the server
   
        const response = await fetch('http://localhost:9090/env');
        // `import.meta.env.VITE_APP_URL` is used as the base URL for the request
        // This allows the app to fetch environment variables from a backend API

        const data = await response.json(); 
        // Parse the response as JSON to get the environment variables object

        // Loop through each key-value pair in the fetched data
        for (let key in data) {         
            env.value[key] = data[key]
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