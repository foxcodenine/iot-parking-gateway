<template>
    <div class="ttable__container">
        <input class="ttable__search mt-8" v-model="searchTerm" type="text" placeholder="Search...">

        <table class="ttable  mt-8">
            <thead>
                <tr>
                    <th class="cursor-pointer">
                        #
                        <svg class="t-sort-arrow t-sort-arrow--active">
                            <use xlink:href="@/assets/svg/sprite.svg#triangle-1"></use>
                        </svg>
                    </th>
                    <th class="cursor-pointer">
                        Email
                        <svg class="t-sort-arrow">
                            <use xlink:href="@/assets/svg/sprite.svg#triangle-1"></use>
                        </svg>
                    </th>
                    <th class="cursor-pointer">
                        Access Level
                        <svg class="t-sort-arrow">
                            <use xlink:href="@/assets/svg/sprite.svg#triangle-1"></use>
                        </svg>
                    </th>
                    <th class="cursor-pointer">
                        Enabled
                        <svg class="t-sort-arrow">
                            <use xlink:href="@/assets/svg/sprite.svg#triangle-1"></use>
                        </svg>
                    </th>

                    <th class="cursor-pointer">
                        Registered
                        <svg class="t-sort-arrow">
                            <use xlink:href="@/assets/svg/sprite.svg#triangle-1"></use>
                        </svg>
                    </th>
                    <th>
                    </th>
                </tr>
            </thead>
            <tbody>
                <tr v-for="user in getUsersList">
                    <td>{{ user.id }}</td>
                    <td>{{ user.email }}</td>
                    <td>{{ getAccessLevelDescriptionUtil(user.access_level) }}</td>
                    <td>
                        <svg class="t-icon-enabled ms-4" v-if="user.enabled">
                            <use xlink:href="@/assets/svg/sprite.svg#icon-c-enabled"></use>
                        </svg>
                        <svg class="t-icon-disabled ms-4" v-else>
                            <use xlink:href="@/assets/svg/sprite.svg#icon-c-disabled"></use>
                        </svg>
                    </td>
                    <td>{{ formatDateUtil(user.created_at) }}</td>
                    <td>
                        <div class="t-btns" v-if="user.access_level >= 1" @click="goToView('userEditView', user.id)">
                            <a class="t-btns__btn  ml-auto">
                                <svg class="t-btns__icon">
                                    <use xlink:href="@/assets/svg/sprite.svg#icon-pencil"></use>
                                </svg>
                            </a>
                        </div>
                        <div v-else></div>
                    </td>
                </tr>
            </tbody>
        </table>
    </div>
</template>

<!-- --------------------------------------------------------------- -->

<script setup>
import { useUserStore } from '@/stores/userStore';
import { formatDateUtil, getAccessLevelDescriptionUtil } from '@/utils/utils';
import { storeToRefs } from 'pinia';
import { ref } from 'vue';
import { useRouter } from 'vue-router';

const router = useRouter();
// - Store -------------------------------------------------------------

const userStore = useUserStore();
const { getUsersList } = storeToRefs(userStore);

// - Data --------------------------------------------------------------
const searchTerm = ref("");

// -- methods ----------------------------------------------------------

function goToView(view, id) {
    router.push({ name: view, params: { userID: id } });
}
</script>

<!-- --------------------------------------------------------------- -->

<style lang="scss" scoped>
// hello</style>