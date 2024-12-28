<template>
    <section id="the-sidebar" class="sidebar">
        <div class="sidebar__list ">
            <div id="menu-btn" class="sidebar__item" @click="toggleMenu" >
                <svg class="sidebar__svg ">
                    <use xlink:href="@/assets/svg/sprite.svg#icon-menu-1"></use>
                </svg>
                <div class="sidebar__text" >
                    <span>Menu</span>
                </div>
            </div>

            <div class="sidebar__item" @click="goToView('mapView')" >
                <svg class="sidebar__svg ">
                    <use xlink:href="@/assets/svg/sprite.svg#icon-map-8"></use>
                </svg>
                <div class="sidebar__text" >
                    <span>Map</span>
                </div>
            </div>
    
            <div class="sidebar__item" @click="goToView('userView')" >
                <svg class="sidebar__svg ">
                    <use xlink:href="@/assets/svg/sprite.svg#icon-user-31"></use>
                </svg>
                <div class="sidebar__text">
                    <span>Users</span>
                </div>
            </div>
       
            <div class="sidebar__item" @click="goToView('deviceView')" >
                <svg class="sidebar__svg ">
                    <use xlink:href="@/assets/svg/sprite.svg#icon-view-12"></use>
                </svg>
                <div class="sidebar__text">
                    <span>Devices</span>
                </div>
            </div>

            <div class="sidebar__item" @click="goToView('deviceView')" >
                <svg class="sidebar__svg ">
                    <use xlink:href="@/assets/svg/sprite.svg#icon-bug"></use>
                </svg>
                <div class="sidebar__text">
                    <span>Debug</span>
                </div>
            </div>

            <div class="sidebar__item" @click="goToView('deviceView')" >
                <svg class="sidebar__svg ">
                    <use xlink:href="@/assets/svg/sprite.svg#icon-docs"></use>
                </svg>
                <div class="sidebar__text">
                    <span>Doc's</span>
                </div>
            </div>

            <div class="sidebar__item" @click="goToView('deviceView')" >
                <svg class="sidebar__svg ">
                    <use xlink:href="@/assets/svg/sprite.svg#icon-about"></use>
                </svg>
                <div class="sidebar__text">
                    <span>About</span>
                </div>
            </div>

            <!-- <div class="sidebar__item" @click="goToView('viewHome')" >
                <svg class="sidebar__svg ">
                    <use xlink:href="@/assets/svg/sprite.svg#icon-gears-2"></use>
                </svg>
                <div class="sidebar__text">
                    <span>Settings</span>
                </div>
            </div> -->
         
        </div>

    </section>

</template>

<!-- --------------------------------------------------------------- -->

<script setup>

import { useDashboardStore } from '@/stores/dashboardStore';
import { useRouter } from 'vue-router';


// - Store -------------------------------------------------------------

const dashboardStore = useDashboardStore();


// -- data -------------------------------------------------------------
const router = useRouter();


// -- methods ----------------------------------------------------------

function goToView(view ) {
    router.push({ name: view });
}

function toggleMenu() {
    dashboardStore.toggleUserMenu();
}

</script>

<!-- --------------------------------------------------------------- -->

<style lang="scss" scoped>
.sidebar {
    background-color: $col-zinc-800; 
    position: fixed;
    top: 0;
    bottom: 0;
    left: 0;
    width: 4rem;

    &__line {
        width: 100%;
        border-top: 1px solid $col-zinc-600;
        // margin: 2px 0;
    }


    &__svg {
        width: 1.2rem;
        height: 1.2rem;
        fill: currentColor;

        &--expand {
            transform: translateX(-7px);

            // @include respondMobile($bp-medium) {
            //     transform: translateX(0px);
            // }
        }
    }

    &__list {
        display: flex;
        flex-direction: column;
        gap: 1px;
        // padding-top: .5rem;
    }

    &__item {
        cursor: pointer;
        display: flex;
        position: relative;
        gap: .3rem;
        flex-direction: column;
        align-items: center;
        justify-content: center;
        color: $col-zinc-100;
        height: 4.2rem;
        padding-top: 5px;
        // border: 1px solid crimson;
        transition: all .1s ease;

        &:hover {
            color: $col-zinc-800;
            background-color: $col-zinc-300;
        }
        

        // @include respondMobile($bp-medium) {
        //     flex-direction: row;
        //     justify-content: start;
        //     gap: 1rem;
        //     height: 3.8rem;
        //     padding: 0 1rem;
        // }

        &--hidden {
            background-color: $col-zinc-300;
            height: 3rem;

            &:not(&:first-child) {
                border-top: 1px solid $col-zinc-50;
            }
        }
    }


    &__text {
        font-family: $font-display;
        font-size: .8rem;
        font-weight: 400;

        // @include respondMobile($bp-medium) {
        //     font-size: 1rem;
        //     font-family: $font-display;
        //     flex: 1;
        // }

        &--mobile {
            display: block;

            // @include respondMobile($bp-medium) {
            //     display: none;
            // }
        }

        &--desktop {
            display: none;

            // @include respondMobile($bp-medium) {
            //     display: block;
            // }
        }
    }
}

</style>