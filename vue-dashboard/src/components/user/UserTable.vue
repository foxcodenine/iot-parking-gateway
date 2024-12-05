<template>
    <input class="ttable__search mt-8" v-model="searchTerm" type="text" placeholder="Search...">
    <AdminConfirmationModal v-if="adminModalOn" 
            @emitCancel="adminModalOn = false" 
            @emitConfirm="deleteUser"
            appear
            >
    </AdminConfirmationModal>
    <div class="ttable__container">

        <table class="ttable  mt-8" @click="clearMessage">
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
                    <th v-if="getUserAccessLevel <= 1">
                    </th>
                </tr>
            </thead>
            <tbody>
                <tr v-for="user in getUsersList" :class="{'bg-lime-200': props.userID == user.id}">
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
                    <td v-if="getUserAccessLevel <= 1" >
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
import { ref } from 'vue';
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
const { getUsersList } = storeToRefs(userStore);

const authStore = useAuthStore();
const { getUserAccessLevel } = storeToRefs(authStore)


// - Data --------------------------------------------------------------
const searchTerm = ref("");
const adminModalOn = ref(false);

// -- methods ----------------------------------------------------------
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