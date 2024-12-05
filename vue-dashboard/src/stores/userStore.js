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

    async function updateUser({user_id, email, password1,  password2, access_level, enabled, admin_password}) {
        useDashboardStore().setIsLoading(true);
        try {
            const payload = {                
                email,
                password1,
                password2,
                access_level,
                enabled,
                admin_password
            };
            
            return await axios.put(`${useAppStore().getAppUrl}/api/user/${user_id}`, payload);
            
        } catch (error){
            console.error('! userStore.updateUser !');
            throw error;  
        } finally {
            useDashboardStore().setIsLoading(false);
        }
    }

    async function deleteUser({user_id, admin_password}) {
        useDashboardStore().setIsLoading(true);
        try {
            const payload = { 
                admin_password
            };
            
            return await axios.delete(`${useAppStore().getAppUrl}/api/user/${user_id}`, {data: {...payload}});
            
        } catch (error){
            console.error('! userStore.deleteUser !');
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

    function removeUserFromList(userID) {
        if (userID) {
            const index = usersList.value.findIndex(u => u.id == userID)
            if (index == -1) { return }

            usersList.value.splice(index, 1)
        }
    }


    function updateUserInList(user) {
        if (user && user.id) {
            const index = usersList.value.findIndex(u => u.id == user.id)
            if (index == -1) { return }
            usersList.value[index].email = user.email;
            usersList.value[index].access_level = user.access_level;
            usersList.value[index].enabled = user.enabled;
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
        updateUser,
        updateUserInList,
        deleteUser,
        removeUserFromList,
    }
});