<template>
    <form class="ssign__form">
        <div>
            <div class="ssign__title mt-4 mb-8">
                Reset Your <b>Password</b> 
            </div>

            <div class="ssign__subtitle">
                "Enter your registered email address below, and we'll send you instructions to reset your password."
            </div>

            <div class="ssign__group mt-8">
                <label for="sign-email" class="ssign__label">Email</label>
                <input id="sign-email" type="text" class="ssign__input">
            </div>

            <div class="ssign__row mt-3">
                <div class="ssign__placeholder"></div>
                <div class="ssign__forgot" @click="goToView('loginView')">Back to Login</div>
            </div>
        </div>

        <button class="bbtn bbtn--sky">Send Reset Link</button>
    </form>
</template>

<!-- --------------------------------------------------------------- -->

<script setup>
import { ref } from 'vue';
import { useRouter } from 'vue-router';

// -- data -------------------------------------------------------------
const router = useRouter();

const message = ref("An email has been sent to your address with instructions to reset your password. Please check your inbox, and don't forget to check the spam folder if you don't see it.")
const warringMessage = ref("The email address you entered is not associated with an account. Please try again or contact support if you need further assistance.")

// -- methods ----------------------------------------------------------

function goToView(view ) {
    router.push({ name: view });
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
        b { font-weight: 650; }
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