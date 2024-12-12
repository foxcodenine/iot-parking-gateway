<template>
</template>

<!-- --------------------------------------------------------------- -->

<script setup>
import { useAppStore } from '@/stores/appStore';
import { useAuthStore } from '@/stores/authStore';
import axios from '@/axios';
import { onMounted } from 'vue';
import { useRouter } from 'vue-router';
import { useDashboardStore } from '@/stores/dashboardStore';

onMounted(async () => {

    (async () => {
        try {

            const user = useAuthStore().getUserTokenData;

            // if (user.id && user.email) {
            //     axios.post(`${import.meta.env.VITE_WEB_BACKEND_NODE_AUTH_URL}/auth/logout`, {user});
            // }

            useAuthStore().resetAuthStore();
            useAppStore().resetAppStore(); 
            useDashboardStore().updateUserMenu(false);

            useRouter().push({ name: 'loginView' });

            localStorage.clear();
            sessionStorage.clear();

        } catch (error) {
            console.error('! logout !  \n', error);
        }
    })();

})


</script>

<!-- --------------------------------------------------------------- -->

<style lang="scss" scoped></style>