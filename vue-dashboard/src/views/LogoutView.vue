<template>
    <div style="display: none;"></div>
</template>

<script setup>
import { useAuthStore } from '@/stores/authStore';
import { useAppStore } from '@/stores/appStore';
import { useDashboardStore } from '@/stores/dashboardStore';
import axios from '@/axios';
import { onActivated, onMounted } from 'vue';
import { useRouter } from 'vue-router';

// Initialize composables outside of any functions
const authStore = useAuthStore();
const appStore = useAppStore();
const dashboardStore = useDashboardStore();
const router = useRouter();



onActivated(async () => {
    try {
        // Perform the logout API request
        await axios.post(`${appStore.getAppUrl}/api/auth/logout`);

        // Clear user-related states in the store
        authStore.resetAuthStore();
        appStore.resetAppStore();
        

    } catch (error) {
        console.error('Logout Error:', error);
        // Optionally handle errors specific to logout failure if needed
    } finally {
        dashboardStore.updateUserMenu(false);
        
        // Always clear local storage and session storage
        localStorage.clear();
        sessionStorage.clear();

        // Navigate to the login view regardless of the outcome of the above operations
        router.push({ name: 'loginView' });
        // window.location.assign('/login');
    }
});
</script>
