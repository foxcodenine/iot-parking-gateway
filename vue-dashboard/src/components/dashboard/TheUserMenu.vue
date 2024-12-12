<template>

    <div id="the-user-menu" class="menu m-1">
        <div class="menu__link">
            <div class="menu__text">Profile Settings</div>
            <svg class="menu__icon">
                <use xlink:href="@/assets/svg/sprite.svg#icon-user-settings"></use>
            </svg>
        </div>

        <div class="menu__link">
            <div class="menu__text">Application Settings</div>
            <svg class="menu__icon">
                <use xlink:href="@/assets/svg/sprite.svg#icon-settings-2"></use>
            </svg>
        </div>

        <div class="menu__link" @click="logout()">
            <div class="menu__text">Logout</div>
            <svg class="menu__icon">
                <use xlink:href="@/assets/svg/sprite.svg#icon-logout-1"></use>
            </svg>
        </div>
    </div>

</template>

<!-- --------------------------------------------------------------- -->

<script setup>
import { useRouter } from 'vue-router';


// - router ------------------------------------------------------------

const router = useRouter();

// - methods -----------------------------------------------------------

function logout() {
    router.push({ name: 'logoutView' });
}

</script>

<!-- --------------------------------------------------------------- -->
 
<style lang="scss" scoped>
.menu {
    background-color: $col-zinc-800;
    border-radius: $border-radius;
    box-shadow: $box-shadow-4;
    width: 14rem;


    &__link {
        cursor: pointer;
        display: flex;
        justify-content: space-between;
        font-family: $font-display;
        gap: 1rem;

        width: 100%;
        padding: .5rem .8rem;
        color: $col-zinc-100;
        font-size: .9rem;
        transition: all .1s ease;

        &:not(&:last-child) {
            border-bottom: 1px solid $col-zinc-200;
        }

        &:hover {
            background-color: $col-zinc-300;
            color: $col-zinc-800;
        }
    }

    &__icon {
        width: 1.2rem;
        height: 1.2rem;
        fill: currentColor;
        stroke: currentColor;
    }
}
</style>