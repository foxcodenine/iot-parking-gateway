<template>
    <form class="ssign__form">
        <div>
            <div class="ssign__title mt-4 mb-12">
            Welcome to <b>IoTrack</b> Pro
        </div>

        <div class="ssign__group">
            <label for="sign-email" class="ssign__label">Email</label>
            <input id="sign-email" type="text" class="ssign__input">
        </div>

        <div class="ssign__group mt-6">
            <label for="sign-password" class="ssign__label">Password</label>
            <input id="sign-password" type="text" class="ssign__input">
        </div>

        <div class="ssign__row mt-3">
            <div class="ssign__remember" @click="rememberMe = !rememberMe">
                <svg class="ssign__icon" v-if="rememberMe">
                    <use xlink:href="@/assets/svg/sprite.svg#icon-checkbox-3" v-if="true"></use>
                </svg>
                <svg class="ssign__icon ssign__icon--empty" v-else></svg>
                Remeber me
            </div>
            <div class="ssign__forgot"  @click="goToView('forgotPasswordView')">Forgot your password?</div>
        </div>
        </div>

        <button class="bbtn bbtn--sky">Sign In</button>
    </form>
</template>

<!-- --------------------------------------------------------------- -->

<script setup>
import { ref } from 'vue';
import { useRouter } from 'vue-router';

// -- data -------------------------------------------------------------
const router = useRouter();

const rememberMe = ref(false)

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