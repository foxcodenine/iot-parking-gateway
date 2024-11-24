import { ref, computed } from 'vue'
import { defineStore } from 'pinia'
import axios from '@/axios';
import { useAppStore } from './appStore';

// ---------------------------------------------------------------------


export const useUserStore = defineStore("userStore", () => {


    // - State ---------------------------------------------------------


    // - Getters -------------------------------------------------------

    // - Actions -------------------------------------------------------
    async function createUser({email, password, accessLevel}) {
        try {
            const payload = {
                email: email,
                password: password,
                access_level: accessLevel,
            };
            const response = await axios.post(useAppStore().getAppUrl + '/api/user');
        } catch {
            console.error('! userStore.createUser !', error);
            throw error;  
        }
    }

    function reset() {
    }

    // - Expose --------------------------------------------------------

    return {
        reset,
        createUser,
   
    }
});