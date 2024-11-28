<template>
    <form class="fform" autocomplete="off" :class="{ 'fform__disabled': confirmOn }">

        <div class="fform__row mt-8 " @click="clearMessage">
            <div class="fform__group ">
                <label class="fform__label" for="email">Email <span class="fform__required">*</span></label>
                <input class="fform__input" id="email" type="text" placeholder="Enter email" v-model.trim="email">
            </div>

            <TheSelector :options="returnAccessLevelOptions" :selectedOption="selectedOptions['accessLevel']"
                fieldName="accessLevel" label="Access Level" :isDisabled="confirmOn" :isRequired="true"
                @emitOption="selectedOptions['accessLevel'] = $event"></TheSelector>

        </div>
        <div class="fform__row mt-4" @click="clearMessage">
            <div class="fform__group">
                <label class="fform__label" for="password1">Password<span class="fform__required">*</span></label>
                <input class="fform__input" id="password1" type="password" placeholder="Enter password"
                    v-model.trim="password1">
            </div>
            <div class="fform__group">
                <label class="fform__label" for="password2">Confirm Password<span
                        class="fform__required">*</span></label>
                <input class="fform__input" id="password2" type="password" placeholder="Renter password"
                    v-model.trim="password2">
            </div>


        </div>
     
            <button class="bbtn bbtn--blue mt-8" v-if="!confirmOn" @click.prevent="initCreateUser()">Create
                User</button>

            <div class="bbtn__row mt-8">

                <button class="bbtn bbtn--zinc-lt" v-if="confirmOn" @click.prevent="confirmOn = false">Cancel</button>
                <button class="bbtn bbtn--blue" v-if="confirmOn" @click.prevent="createUser">Confirm</button>
            </div>



    </form>
</template>

<!-- --------------------------------------------------------------- -->
<script setup>
import TheSelector from '@/components/commen/TheSelector.vue'
import { useMessageStore } from '@/stores/messageStore';
import { useUserStore } from '@/stores/userStore';
import { computed, reactive, ref, watch } from 'vue';

// - Store -------------------------------------------------------------
const messageStore = useMessageStore();
const userStore = useUserStore();

// - Data --------------------------------------------------------------
const confirmOn = ref(false);

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

// - computed ----------------------------------------------------------

const returnAccessLevelOptions = computed(() => {
    return accessLevelList.value.map((org => {
        // TODO: use utilStore.capitalizeFirstLetter() for value.
        return { ...org, _key: org.id, _value: org.name }
    }))
})

// - watchers -----------------------------------------------------------

watch(() => selectedOptions.accessLevel, (val, oldVal) => {
    accessLevel.value = val._key;
}, { deep: true });

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

function initCreateUser() {
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
    if (password1.value.length < 6) {
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
        messageStore.setFlashClass("flash-message--orange");
        return
    }

    confirmOn.value = true;
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



</script>
<!-- --------------------------------------------------------------- -->

<style lang="scss" scoped>
// Placeholder comment to ensure global styles are imported correctly</style>