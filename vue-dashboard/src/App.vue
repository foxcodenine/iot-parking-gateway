<template>
    <main class="dashboard" :class="{ 'no-side-bar': !showSideBar }" >

        <SockerioClient></SockerioClient>

        <VueLoadingOverlay :active="getIsLoading" :is-full-page="true" :lock-scroll="true" :width="128" :height="128"
            transition="fade" :opacity="0.0" />


        <section class="modal" v-if="getIsUserMenuOpen">
            <TheUserMenu></TheUserMenu>
        </section>

        <section class="sidebar" v-if="showSideBar">
            <TheSidebar></TheSidebar>
        </section>

        <section class="page" >
            <router-view v-slot="{ Component }">
                <component :is="Component" />
            </router-view>
        </section>

    </main>

</template>

<!-- --------------------------------------------------------------- -->

<script setup>
import VueLoadingOverlay from 'vue-loading-overlay';
import 'vue-loading-overlay/dist/css/index.css';
import { RouterLink, RouterView, useRoute } from 'vue-router'
import TheSidebar from './components/dashboard/TheSidebar.vue';
import TheUserMenu from './components/dashboard/TheUserMenu.vue';
import SockerioClient from './components/socketio/SockerioClient.vue';
import { computed, onMounted, ref, watch } from 'vue';
import { useDashboardStore } from './stores/dashboardStore';
import { storeToRefs } from 'pinia';
import { useScrollLock } from '@vueuse/core'
import { useAppStore } from './stores/appStore';

// - Store -------------------------------------------------------------

const dashboardStore = useDashboardStore();
const appStore = useAppStore();
const { getPageScrollDisabled } = storeToRefs(appStore);

// - Routes ------------------------------------------------------------

const route = useRoute();
const { getIsUserMenuOpen, getIsLoading } = storeToRefs(dashboardStore);


// - Page scrolling ----------------------------------------------------

const htmlEL = document.querySelector('html');
const isLocked = useScrollLock(htmlEL);
isLocked.value = getPageScrollDisabled.value;
watch(getPageScrollDisabled, (val) => {
    isLocked.value = val
})

// - Computed ----------------------------------------------------------


const showSideBar = computed(() => {
    return !['loginView', 'forgotPasswordView'].includes(route.name)
});

// - Methods -----------------------------------------------------------

// Close the user menu when clicks outside the user menu and top bar image
function closeUserMenuOnClickOutside() {
    document.querySelector('body').addEventListener('click', (e) => {
        const userMenuOrBtn = e.target.closest('#the-user-menu') || e.target.closest('#menu-btn');
        if (!userMenuOrBtn) {
            dashboardStore.updateUserMenu(false);
        }
    });
}

onMounted(() => {
    closeUserMenuOnClickOutside();
});

</script>

<!-- --------------------------------------------------------------- -->

<style lang="scss" scoped>
.dashboard {
    background-color: $col-slate-50;
    display: grid;
    grid-template-columns: 4rem 1fr;
    min-height: 100vh;
    min-width: 100vw;

    // @include respondMobile($bp-medium) {
    //     grid-template-columns: 17rem 1fr;
    // }
}

.modal {
    position: fixed;
    left: 4.25rem;
    top: .25rem;
    z-index: 500;
}

.no-side-bar {
    grid-template-columns: 1fr
}
</style>
