<template>
    <form class="ssign__form" @submit.prevent="submitForm">
        <div>
            <div class="ssign__title mt-4 mb-6 ">
                Welcome to <b>IoTrack</b> Pro
            </div>

            <div class="ssign__flash-message mb-4 h-5 text-red-600">
                {{ flashMessage }}
            </div>

            <div class="ssign__group" @click="resetFlashMassage()">
                <label for="sign-email" class="ssign__label">Email</label>
                <input id="sign-email" type="text" class="ssign__input" v-model.trim="email">
            </div>

            <div class="ssign__group mt-6" @click="resetFlashMassage()">
                <label for="sign-password" class="ssign__label">Password</label>
                <input id="sign-password" type="password" class="ssign__input" v-model.trim="password">
            </div>

            <div class="ssign__row mt-3" @click="resetFlashMassage()">
                <div class="ssign__remember" @click="toggleRememberMe()">
                    <svg class="ssign__icon" v-if="getRemeberMe">
                        <use xlink:href="@/assets/svg/sprite.svg#icon-checkbox-3"></use>
                    </svg>
                    <svg class="ssign__icon ssign__icon--empty" v-else></svg>
                    Remeber me
                </div>
                <div class="ssign__forgot" @click="goToView('forgotPasswordView')">Forgot your password?</div>
            </div>
        </div>

        <button class="bbtn bbtn--sky">Sign In</button>
    </form>
</template>

<!-- --------------------------------------------------------------- -->

<script setup>

import axios from '@/axios';
import validator from 'validator';
import { ref, watch } from 'vue';
import { useRouter } from 'vue-router';
import { useAppStore } from '@/stores/appStore';
import { useAuthStore } from '@/stores/authStore';
import { storeToRefs } from 'pinia';

// - Store -------------------------------------------------------------

const appStore = useAppStore();
const authStore = useAuthStore();
const { getRemeberMe } = storeToRefs(authStore);

// - Data --------------------------------------------------------------

const router = useRouter();
const email = ref("user@dev.com");
const password = ref("DevPass");
const flashMessage = ref("")



// - Methods -----------------------------------------------------------

function goToView(view) {
    router.push({ name: view });
}

function resetFlashMassage() {
    flashMessage.value = null;
}

async function submitForm() {
    flashMessage.value = '';
    try {
        if (validator.isEmpty(email.value) && validator.isEmpty(password.value)) {
            flashMessage.value = "Please enter your email and password to log in.";
            return
        } else if (validator.isEmpty(email.value)) {
            flashMessage.value = "Please enter your email to log in.";
            return
        } else if (validator.isEmpty(password.value)) {
            flashMessage.value = "Please enter your password to log in.";
            return
        } else if (!validator.isEmail(email.value)) {
            flashMessage.value = "Please enter a valid email address.";
            return
        }
        
        const response = await axios.post(`${appStore.getAppUrl}/api/auth/login`, {
            email: email.value,
            password: password.value
        });       

  

        if (response?.status == 200) {
            email.value = 'user@dev.com';
            password.value = 'DevPass';
            authStore.setJwt(response.data.token);
            appStore.setAppSettings(response.data.settings)
            
        }

        

    } catch (err) {
        const errorMessage = err.response?.data || err.message || 'Unable to log in. Please check your email and password.';
        flashMessage.value = errorMessage;

        console.error("! login form !", err);
    }
}

function toggleRememberMe() {
    authStore.toggleRememberMe();
}






</script>

<!-- --------------------------------------------------------------- -->

<style lang="scss" scoped>
.ssign {

    &__form {
        width: 430px;
        height: 100%;
        padding: 2rem;
        font-family: $font-action;
        font-size: .85rem;
        font-weight: 400;
        color: $col-slate-700;
        display: flex;
        flex-direction: column;
        justify-content: space-between;
    }

    &__title {
        font-size: 1.8rem;
        font-weight: 200;

        b {
            font-weight: 650;
        }
    }

    &__group {
        display: flex;
        flex-direction: column;
    }

    &__label {
        margin-bottom: .3rem;
    }

    &__input {
        width: 100%;
        height: 2rem;
        padding: .2rem .5rem;
        border: 1px solid $col-slate-300;
        border-radius: 3px;
        font-size: 1rem;
        color: $col-slate-600;
    }

    &__row {
        display: flex;
        justify-content: space-between;
        align-items: center;
    }

    &__remember {
        cursor: pointer;
    }

    &__icon {
        width: 1rem;
        height: 1rem;
        fill: $col-blue-600;
        transform: translateY(3px);
        border-radius: 2px;
        border: 1px solid $col-blue-600;

        &--empty {
            border: 1px solid $col-slate-500;
        }
    }

    &__forgot {
        cursor: pointer;
        color: $col-blue-600;

        &:hover {
            text-decoration: underline;
        }
    }
}
</style>