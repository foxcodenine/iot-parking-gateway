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

            useAuthStore().clearJwt();
            useDashboardStore().updateUserMenu(false);

            // useAppStore().resetAppState();

            useRouter().push({ name: 'loginView' });

        } catch (error) {
            console.error('! logout !  \n', error);
        }
    })();

})


</script>

<!-- --------------------------------------------------------------- -->

<style lang="scss" scoped></style>