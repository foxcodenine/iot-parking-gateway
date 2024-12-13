<template>
    <input class="ttable__search mt-8" v-model="searchTerm" type="text" placeholder="Search...">
    <AdminConfirmationModal v-if="adminModalOn" @emitCancel="adminModalOn = false" @emitConfirm="deleteUser" appear>
    </AdminConfirmationModal>
    <div class="ttable__container">

        <table class="ttable  mt-8" @click="clearMessage">
            <thead>
                <tr>
                    <th>
                        <span class="cursor-pointer" @click="sortTable('id')">
                            #
                        </span>
                        <svg class="t-sort-arrow" :class="sortArrow('id')">
                            <use xlink:href="@/assets/svg/sprite.svg#triangle-1"></use>
                        </svg>
                    </th>
                    <th>
                        <span class="cursor-pointer" @click="sortTable('email')">
                            Email
                        </span>
                        <svg class="t-sort-arrow" :class="sortArrow('email')">
                            <use xlink:href="@/assets/svg/sprite.svg#triangle-1"></use>
                        </svg>
                    </th>
                    <th>
                        <span class="cursor-pointer" @click="sortTable('access_level')">
                            Access Level
                        </span>
                        <svg class="t-sort-arrow" :class="sortArrow('access_level')">
                            <use xlink:href="@/assets/svg/sprite.svg#triangle-1"></use>
                        </svg>
                    </th>
                    <th>
                        <span class="cursor-pointer" @click="sortTable('enabled')">
                            Enabled
                        </span>
                        <svg class="t-sort-arrow" :class="sortArrow('enabled')">
                            <use xlink:href="@/assets/svg/sprite.svg#triangle-1"></use>
                        </svg>
                    </th>

                    <th>
                        <span class="cursor-pointer" @click="sortTable('created_at')">
                            Registered
                        </span>
                        <svg class="t-sort-arrow" :class="sortArrow('created_at')">
                            <use xlink:href="@/assets/svg/sprite.svg#triangle-1"></use>
                        </svg>
                    </th>
                    <th v-if="getUserAccessLevel <= 1">
                    </th>
                </tr>
            </thead>
            <tbody>
                <tr v-for="user in getUsersList" :class="{ 'bg-lime-200': props.userID == user.id }">
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
                    <td v-if="getUserAccessLevel <= 1">
                        <div class="t-btns ml-auto" v-if="user.access_level >= 1">
                            <a class="t-btns__btn " @click="goToView('userEditView', user.id)">
                                <svg class="t-btns__icon">
                                    <use xlink:href="@/assets/svg/sprite.svg#icon-pencil"></use>
                                </svg>
                            </a>
                            <a class="t-btns__btn" @click="initDeleteUser(user.id)">
                                <svg class="t-btns__icon">
                                    <use xlink:href="@/assets/svg/sprite.svg#icon-trash-can-solid2"></use>
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
import { computed, ref } from 'vue';
import { useRouter } from 'vue-router';
import AdminConfirmationModal from '../commen/AdminConfirmationModal.vue';
import { useMessageStore } from '@/stores/messageStore';
import { useAuthStore } from '@/stores/authStore';


// - Props -------------------------------------------------------------

const props = defineProps({
    userID: {
        type: String,
        required: false, // Because it won't be present in the user list view
    }
});

// - Route -------------------------------------------------------------

const router = useRouter();

// - Store -------------------------------------------------------------

const  messageStore = useMessageStore();

const userStore = useUserStore();

const authStore = useAuthStore();
const { getUserAccessLevel } = storeToRefs(authStore)


// - Data --------------------------------------------------------------

const sortBy = ref('created_at');
const sortDesc = ref('true');
const searchTerm = ref("");

const adminModalOn = ref(false);

// -- Computed ---------------------------------------------------------

const getUsersList = computed(() => {
    
    let list = [...userStore.getUsersList];

    list = list.filter(item => {
        return (
            String(item.id)?.toLowerCase().includes(searchTerm.value.toLowerCase().trim()) ||
            item.email?.toLowerCase().includes(searchTerm.value.toLowerCase().trim()) ||
            getAccessLevelDescriptionUtil(item.access_level).toLowerCase().includes(searchTerm.value.toLowerCase())
        )
    })

    list.sort((a, b) => {
        let modifier = sortDesc.value ? -1 : 1;

        if (a[sortBy.value] < b[sortBy.value]) return -1 * modifier;
        if (a[sortBy.value] > b[sortBy.value]) return 1 * modifier;

        return 0;

    })
    return list;
});

// -- Methods ----------------------------------------------------------

function sortTable(col) {
    if (sortBy.value == col) {
        sortDesc.value = !sortDesc.value
    }    
    sortBy.value = col;
}

function sortArrow(col) {
    if (col === sortBy.value) {
        return { 't-sort-arrow--active': true, 't-sort-arrow--desc': sortDesc.value }
    }
}

function clearMessage() {
    messageStore.clearFlashMessage();
}

function goToView(view, id) {    
    router.push({ name: view, params: { userID: id } });
}

function initDeleteUser(id) {    
    router.push({ name: 'userEditView', params: { userID: id } });    
    adminModalOn.value = true;
}

async function deleteUser(payload) {

    adminModalOn.value = false;
    
    try {
        const response = await userStore.deleteUser({
            user_id: props.userID,
            admin_password: payload.adminPassword,
        });

        if (response.status == 200) {
            const msg = response.data?.message ?? "User deleted successfully.";
            userStore.removeUserFromList(props.userID);
            router.push({ name: 'userView' });
            setTimeout(()=> {
                messageStore.setFlashMessages([msg], "flash-message--green");
            }, 500);
        }
        
    } catch (error) {
        console.error("! UserForm.deleteUser !\n", error);
        const errMsg = error.response?.data ?? "Failed to delete user"
        messageStore.setFlashMessages([errMsg], "flash-message--red");
    } finally {
    }      
}
</script>

<!-- --------------------------------------------------------------- -->

<style lang="scss" scoped>
.t-btns {
    justify-content: end;
}


</style>