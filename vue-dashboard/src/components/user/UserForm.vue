<template>
    <form class="fform" autocomplete="off">
        <transition name="modal" appear>
            <AdminConfirmationModal v-if="adminModalOn" 
                @emitCancel="adminModalOn = false" 
                @emitConfirm="updateUser"
                appear
                >
            </AdminConfirmationModal>
        </transition>

        <div class="fform__row mt-8 " @click="clearMessage" :class="{ 'fform__disabled': confirmOn }">
            <div class="fform__group ">
                <label class="fform__label" for="email">Email <span class="fform__required">*</span></label>
                <input class="fform__input" id="email" type="text" placeholder="Enter the user's email" v-model.trim="email"
                    :disabled="confirmOn">
            </div>

            <TheSelector :options="returnAccessLevelOptions" :selectedOption="selectedOptions['accessLevel']"
                fieldName="accessLevel" label="Access Level" :isDisabled="confirmOn" :isRequired="true"
                @emitOption="selectedOptions['accessLevel'] = $event"></TheSelector>

        </div>
        <div class="fform__row mt-4" @click="clearMessage" :class="{ 'fform__disabled': confirmOn }">
            <div class="fform__group">
                <label class="fform__label" for="password1">Password<span :class="{'fform__required': !editMode}">{{!editMode ? '*' : ''}}</span></label>
                <input class="fform__input" id="password1" type="password" :placeholder="!editMode ? 'Enter the user\'s password' : 'Leave blank to keep unchanged'"
                    v-model.trim="password1" :disabled="confirmOn">
            </div>
            <div class="fform__group">
                <label class="fform__label" for="password2">Confirm Password<span
                    :class="{'fform__required': !editMode}">{{!editMode ? '*' : ''}}</span></label>
                <input class="fform__input" id="password2" type="password" placeholder="Renter password"
                    v-model.trim="password2" :disabled="confirmOn">
            </div>
        </div>

        <TheCheckbox v-if="editMode" :is-checked="accoutEnabled" @emit-checkbox="accoutEnabled = !accoutEnabled">Account Enabled</TheCheckbox>

        <transition name="fade" mode="out-in">
            <button v-if="!confirmOn" class="bbtn mt-8" :class="{ 'bbtn--red': editMode, 'bbtn--blue': !editMode }"
                @click.prevent="initCreateOrUpdateUser()" key="create-button">
                {{ editMode ? "Update User" : "Create New User" }}
            </button>

            <div v-else class="bbtn__row mt-8" key="confirm-buttons">
                <button class="bbtn bbtn--zinc-lt" @click.prevent="confirmOn = false">Cancel</button>
                <button class="bbtn bbtn--blue" @click.prevent="createUser">Confirm</button>
            </div>
        </transition>



    </form>
</template>

<!-- --------------------------------------------------------------- -->
<script setup>
import TheSelector from '@/components/commen/TheSelector.vue'
import TheCheckbox from '../commen/TheCheckbox.vue';
import { useMessageStore } from '@/stores/messageStore';
import { useUserStore } from '@/stores/userStore';
import { computed, reactive, ref, watch } from 'vue';
import AdminConfirmationModal from '../commen/AdminConfirmationModal.vue';


// - Store -------------------------------------------------------------
const messageStore = useMessageStore();
const userStore = useUserStore();

// - Props -------------------------------------------------------------

const props = defineProps({
    userID: {
        type: String,
        required: false, // Because it won't be present in the user list view
    }
});

// - Data --------------------------------------------------------------
const confirmOn = ref(false);
const editMode = ref(false);
const adminModalOn = ref(false)


const accessLevelList = ref([
    // {id: 0, name: 'Root'},
    { id: 1, name: 'Administrator' },
    { id: 2, name: 'Editor' },
    { id: 3, name: 'Viewer' },
]);

const selectedOptions = reactive({
    'accessLevel': { _key: 1, _value: 'Administrator' }
});

const email = ref("");
const password1 = ref("");
const password2 = ref("");
const accessLevel = ref(1);
const accoutEnabled = ref(true)


// - computed ----------------------------------------------------------

const returnAccessLevelOptions = computed(() => {
    return accessLevelList.value.map((accessLvl => {
        // TODO: use utilStore.capitalizeFirstLetter() for value.
        return { ...accessLvl, _key: accessLvl.id, _value: accessLvl.name }
    }))
});

const getUser = computed(() => {
    return userStore.getUserById(Number(props.userID));
});

// - watchers -----------------------------------------------------------

watch(() => selectedOptions.accessLevel, (val, oldVal) => {
    accessLevel.value = val._key;
}, { deep: true });

watch(() => getUser, (val, oldVal) => {
    confirmOn.value = false;
    password1.value = '';
    password2.value = '';
    if (val.value) {
        editMode.value = true;
        email.value = val.value.email;
        accoutEnabled.value = val.value.enabled
        const accessLevel = accessLevelList.value.find(accesslvl => accesslvl.id === Number(val.value.access_level));
        selectedOptions['accessLevel'] = { ...accessLevel, _key: accessLevel.id, _value: accessLevel.name }
    } else {
        editMode.value = false;
    }
}, { deep: true, immediate: true });

watch(editMode, (val) => {
    if (!val) {
        email.value = '';
        accessLevel.value = 1;
        selectedOptions.accessLevel = { _key: 1, _value: 'Administrator' };
    }
})

// - methods -----------------------------------------------------------

function clearMessage() {
    messageStore.clearFlashMessage();
}

function resetForm() {
    email.value = "";
    password1.value = "";
    password2.value = "";
    accessLevel.value = 1;
}

function initCreateOrUpdateUser() {
    messageStore.clearFlashMessage();
    const message = []; // Clear previous messages
    let hasError = false;

    // Check if the email is in the correct format
    const emailPattern = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
    if (!emailPattern.test(email.value)) {
        message.push("Invalid email format.");
        hasError = true;
    }

    // Check if access level is within the valid range
    if (accessLevel.value < 0 || accessLevel.value > 3) {
        message.push("Access level must be between 0 and 3.");
        hasError = true;
    }

    // Check if password is more than 6 characters
    if (!editMode.value && password1.value.length < 6) {
        message.push("Password must be longer than 6 characters.");
        hasError = true;
    }

    // Check if passwords match
    if (password1.value !== password2.value) {
        message.push("Passwords do not match.");
        hasError = true;
    }

    if (hasError) {
        messageStore.setFlashMessages(message);
        messageStore.setFlashClass("flash-message--yellow");
        return
    }

    if (!editMode.value) {
        confirmOn.value = true;
    } else {
        adminModalOn.value = true;
    }
}



async function createUser() {
    try {
        const response = await userStore.createUser({
            email: email.value,
            password1: password1.value,
            password2: password2.value,
            accessLevel: accessLevel.value
        });

        if (response.status == 201) {
            const msg = response.data?.message ?? "User created successfully.";
            messageStore.setFlashMessages([msg], "flash-message--green");
            resetForm();
            userStore.pushUserToList(response.data?.user);

        }


    } catch (error) {
        console.error("! UserForm.createUser !\n", error);
        const errMsg = error.response?.data ?? "Failed to create user"
        messageStore.setFlashMessages([errMsg], "flash-message--red");

    } finally {
        confirmOn.value = false;
    }
}

async function updateUser(payload) {
    adminModalOn.value = false;
    try {
        const response = await userStore.updateUser({
            user_id: props.userID,
            email: email.value,
            password1: password1.value,
            password2: password2.value,
            access_level: accessLevel.value,
            enabled: accoutEnabled.value,
            admin_password: payload.adminPassword,
        });

        if (response.status == 200) {
            const msg = response.data?.message ?? "User updated successfully.";
            messageStore.setFlashMessages([msg], "flash-message--green");
     
            if (response.data?.user) {
                userStore.updateUserInList(response.data.user);
            }

     
        }
    } catch (error) {
        console.error("! UserForm.updateUser !\n", error);
        const errMsg = error.response?.data ?? "Failed to update user"
        messageStore.setFlashMessages([errMsg], "flash-message--red");
    } finally {
        confirmOn.value = false;
    }
}




</script>
<!-- --------------------------------------------------------------- -->

<style lang="scss" scoped>
.display-none {
    display: none !important;
}
</style>