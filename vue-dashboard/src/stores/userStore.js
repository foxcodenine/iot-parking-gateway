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

    const getUsersList = computed(() => usersList.value );

    // Returns a computed property that gives a user by ID
     const getUserById = (id) => {
        return usersList.value.find(user => user.id === id);
    };

    // - Actions -------------------------------------------------------
    async function fetchUsers() {
        useDashboardStore().setIsLoading(true);
        try {
            const response =  await axios.get(useAppStore().getAppUrl + '/api/user');
            if (response.status == 200 && response.data?.users) {
                usersList.value = response.data.users;
                return usersList.value;
            }
        } catch (error) {
            console.error('! userStore.fetchUsers !');
            throw error; 
        } finally {
            useDashboardStore().setIsLoading(false);
        }
    }

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
        fetchUsers,
        getUsersList,
        getUserById,
    }
});