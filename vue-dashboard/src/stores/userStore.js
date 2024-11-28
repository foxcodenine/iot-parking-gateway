import { ref, computed } from 'vue'
import { defineStore } from 'pinia'
import axios from '@/axios';
import { useAppStore } from './appStore';
import { useDashboardStore } from './dashboardStore';

// ---------------------------------------------------------------------


export const useUserStore = defineStore("userStore", () => {

    // - State ---------------------------------------------------------
    const usersList = ref([]);

    // - Getters -------------------------------------------------------

    // - Actions -------------------------------------------------------
    async function createUser({email, password1,  password2, accessLevel}) {
        useDashboardStore().setIsLoading(true);
        try {
            const payload = {
                email: email,
                password1: password1,
                password2: password2,
                access_level: accessLevel,
            };
            return await axios.post(useAppStore().getAppUrl + '/api/user', payload);
            
        } catch (error){
            console.error('! userStore.createUser !');
            throw error;  
        } finally {
            useDashboardStore().setIsLoading(false);
        }
    }

    function reset() {
    }

    function pushUserToList(user) {
        if (user && user.id) {
            usersList.value.push(user);
            console.log(usersList.value);
        }
        
    }

    // - Expose --------------------------------------------------------

    return {
        reset,
        createUser,
        pushUserToList,
    }
});