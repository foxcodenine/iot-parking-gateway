<template>
    <main class="vview__main">
        <section class="vview__section">
            <div class="heading--2 ">Users</div>
            <TheFlashMessage ></TheFlashMessage>            
    
            <div class="heading--4" v-if="!props.userID && getUserAccessLevel <= 1">CREATE NEW USER</div>
            <div class="heading--4" v-if="props.userID && getUserAccessLevel <= 1">Edit USER</div>
            <KeepAlive>
                <UserForm v-if="getUserAccessLevel <= 1" :userID="props.userID"></UserForm>
            </KeepAlive>
            
            <div class="heading--4 mt-8">USER LIST</div>
            <UserTable :userID="props.userID"></UserTable>
            
        </section>
    </main>
</template>

<!-- --------------------------------------------------------------- -->

<script setup>
import { ref } from 'vue';
// import FormOrganisation from '@/components/organisation/FormOrganisation.vue'
import TheFlashMessage from '@/components/commen/TheFlashMessage.vue';
import UserForm from '@/components/user/UserForm.vue'
import UserTable from '@/components/user/UserTable.vue';
import { useAuthStore } from '@/stores/authStore';
import { storeToRefs } from 'pinia';
import { useUserStore } from '@/stores/userStore';

// - Store -------------------------------------------------------------
const userStore = useUserStore();

const authStore = useAuthStore();
const { getUserAccessLevel } = storeToRefs(authStore)

// - Props -------------------------------------------------------------

const props = defineProps({
    userID: {
        type: String,
        required: false, // Because it won't be present in the user list view
    }
});

// - Methods -----------------------------------------------------------


// - Hooks -------------------------------------------------------------

try {
    userStore.fetchUsers();
} catch (error) {
    console.error('! UserView userStore.fetchUsers() !\n', error);
}



</script>

<!-- --------------------------------------------------------------- -->

<style lang="scss" scoped>
// Placeholder comment to ensure global styles are imported correctly
</style>