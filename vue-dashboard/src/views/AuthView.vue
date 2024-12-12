<template>
    <main class="vview__main">

        <img class="background" src="@/assets/images/sign/photo-1724274876097-103bb600debb.avif" alt="">

        <div class="ssign">
            <div class="ssign__image">

            </div>
            <component :is="currentFormComponent" />
    
        </div>

    </main>
</template>

<!-- --------------------------------------------------------------- -->

<script setup>
import { computed, ref } from 'vue';
import LoginForm from '@/components/auth/LoginForm.vue';
import ForgotPasswordForm from '@/components/auth/ForgotPasswordForm.vue';
import { useRoute, useRouter } from 'vue-router';
import { watch } from 'vue';
import { useAuthStore } from '@/stores/authStore';
import { storeToRefs } from 'pinia';


// Use Vue Router to access the current route
const route = useRoute();
const router = useRouter();

const authStore = useAuthStore();
const { isAuthenticated, getRedirectTo } = storeToRefs(authStore);

// -- computed ----------------------------------------------------------

// Determine the correct component based on the route name
const currentFormComponent = computed(() => {
    if (route.name === 'loginView') {
        return LoginForm;
        
    } else if (route.name === 'forgotPasswordView') {
        return ForgotPasswordForm;
    }
    return null; // Or a default component if necessary
});

watch(isAuthenticated, (newVal) => {
    if (newVal) {
        router.push({ name: getRedirectTo.value });
    } else {
        router.push({ name: "loginView" });
    } 
}, {
    immediate: true,
});

</script>

<!-- --------------------------------------------------------------- -->

<style lang="scss" scoped>
.vview__main {
    position: relative; // Ensure .background and .ssign share the same stacking context

    @include respondDesktop(425) {
        padding: 01rem !important;
    }
}

.background {
    object-fit: cover;
    object-position: center;
    opacity: 0.5;
    position: absolute;
    top: 0;
    left: 0;
    width: 100%;
    height: 100%;
    z-index: 0; // Below .ssign
}

.ssign {
    position: relative;
    z-index: 1000 !important;
    margin: calc((80vh - 468px) / 2) auto;
    height: 468px;
    max-width: 700px;
    border: 1px solid $col-slate-300;
    display: flex;
    border-radius: 5px;
    background-color: $col-white;

    @include respondDesktop(425) {
        margin: calc((90vh - 468px) / 2) auto;
    }

    @include respondDesktop(700) {
        max-width: 430px;
    }


    &__image {
        height: 100%;
        flex: 1;
        background-image: url("@/assets/images/sign/photo-1445548671936-e1ff8a6a6b20.avif");
        background-size: cover;
        background-position: center;

        @include respondDesktop(700) {
            display: none;
        }
    }
}

</style>