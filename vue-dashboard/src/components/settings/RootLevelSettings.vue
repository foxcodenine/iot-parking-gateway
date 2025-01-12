<template>
    <form class="fform" autocomplete="off">

        <p class="flash-message flash-message--zinc mt-10"><i>"Updating these settings will log out all users to apply the new configuration."</i></p>

        <transition name="modal" appear>
            <AdminConfirmationModal v-if="adminModalOn" 
                @emitCancel="adminModalOn = false" 
                @emitConfirm="updateSettings"
                appear
                >
            </AdminConfirmationModal>
        </transition>

        <div class="fform__row mt-10 " @click="clearMessage" :class="{ 'fform__disabled': edit_key != 'google_api_key' }">
            <div class="fform__description"><b>google_api_key</b> &nbsp API key used for accessing Google services like Maps and Places.</div>
            <div class="fform__group ">
                <input class="fform__input " :class="{'fform__input--active':edit_key=='google_api_key'}" id="google_api_key" type="text" @blur="edit_key=null" 
                    v-model.trim="google_api_key"  :disabled="edit_key != 'google_api_key'" placeholder="Enter new Google API key to change, else leave empty">
                <svg @click="edit_key='google_api_key'" class="fform__icon" :class="{'fform__icon--active':edit_key=='google_api_key'}"> <use xlink:href="@/assets/svg/sprite.svg#icon-pencil"></use></svg>
            </div>
        </div>

        <div class="fform__row mt-10 " @click="clearMessage" :class="{ 'fform__disabled': edit_key != 'google_map_id' }">
            <div class="fform__description"><b>google_map_id</b> &nbsp The Google Map ID used to customize and embed Google Maps in the application.</div>
            <div class="fform__group ">
                <input class="fform__input " :class="{'fform__input--active':edit_key=='google_map_id'}" id="google_map_id" type="text" @blur="edit_key=null" 
                    v-model.trim="google_map_id" :disabled="edit_key != 'google_map_id'" placeholder="Enter new default longitude to change, else leave empty">
                <svg @click="edit_key='google_map_id'" class="fform__icon" :class="{'fform__icon--active':edit_key=='google_map_id'}"> <use xlink:href="@/assets/svg/sprite.svg#icon-pencil"></use></svg>
            </div>
        </div>

        <div class="fform__row mt-10 " @click="clearMessage" :class="{ 'fform__disabled': edit_key != 'jwt_expiration_seconds' }">
            <div class="fform__description"><b>jwt_expiration_seconds</b> &nbsp Duration in seconds for which a user's JSON Web Token (JWT) remains valid after login.</div>
            <div class="fform__group ">
                <input class="fform__input " :class="{'fform__input--active':edit_key=='jwt_expiration_seconds'}" id="jwt_expiration_seconds" type="text" @blur="edit_key=null" 
                    v-model.trim="jwt_expiration_seconds" :disabled="edit_key != 'jwt_expiration_seconds'" placeholder="Enter new JWT duration to change, else leave empty">
                <svg @click="edit_key='jwt_expiration_seconds'" class="fform__icon" :class="{'fform__icon--active':edit_key=='jwt_expiration_seconds'}"> <use xlink:href="@/assets/svg/sprite.svg#icon-pencil"></use></svg>
            </div>
        </div>

        <div class="fform__row mt-10 " @click="clearMessage" :class="{ 'fform__disabled': edit_key != 'redis_ttl_seconds' }">
            <div class="fform__description"><b>redis_ttl_seconds</b> &nbsp Default time-to-live (TTL) in seconds for items stored in the Redis cache, impacting how long user and device data are cached.</div>
            <div class="fform__group ">
                <input class="fform__input " :class="{'fform__input--active':edit_key=='redis_ttl_seconds'}" id="redis_ttl_seconds" type="text" @blur="edit_key=null" 
                    v-model.trim="redis_ttl_seconds" :disabled="edit_key != 'redis_ttl_seconds'" placeholder="Enter new TTL to change, else leave empty">
                <svg @click="edit_key='redis_ttl_seconds'" class="fform__icon" :class="{'fform__icon--active':edit_key=='redis_ttl_seconds'}"> <use xlink:href="@/assets/svg/sprite.svg#icon-pencil"></use></svg>
            </div>
        </div>

        <div class="fform__row mt-10 " @click="clearMessage" :class="{ 'fform__disabled': edit_key != 'device_access_mode' }">
            <div class="fform__description"><b>device_access_mode</b> &nbsp Defines the access control mode for devices, determining whether they are managed via a blacklist or whitelist approach.</div>
            <div class="fform__group ">
                <select class="fform__input" :class="{'fform__input--active': edit_key == 'device_access_mode','fform__select': edit_key == 'device_access_mode'}" id="device_access_mode" 
                    v-model="device_access_mode" @blur="edit_key=null" :disabled="edit_key != 'device_access_mode'">
                    <option value="white_list">white_list</option>
                    <option value="black_list">black_list</option>
                </select>
                <svg @click="edit_key='device_access_mode'" class="fform__icon" :class="{'fform__icon--active': edit_key == 'device_access_mode'}">
                    <use xlink:href="@/assets/svg/sprite.svg#icon-pencil"></use>
                </svg>
            </div>
        </div>

        <div class="fform__row mt-10 " @click="clearMessage" :class="{ 'fform__disabled': edit_key != 'initial_parking_check_date' }">
            <div class="fform__description"><b>initial_parking_check_date</b> &nbsp The reference date for checking parking events. Devices with no events after this date are considered newly installed or inactive, and their status is marked as unknown.</div>
            <div class="fform__group ">
                <input class="fform__input " :class="{'fform__input--active':edit_key=='initial_parking_check_date'}" id="initial_parking_check_date" type="text" @blur="edit_key=null" 
                    v-model.trim="initial_parking_check_date" :disabled="edit_key != 'initial_parking_check_date'" placeholder="Enter new default longitude to change, else leave empty">
                <svg @click="edit_key='initial_parking_check_date'" class="fform__icon" :class="{'fform__icon--active':edit_key=='initial_parking_check_date'}"> <use xlink:href="@/assets/svg/sprite.svg#icon-pencil"></use></svg>
            </div>
        </div>

        <div class="fform__row mt-10 " @click="clearMessage" :class="{ 'fform__disabled': edit_key != 'cors_allowed_origins' }">
            <div class="fform__description"><b>cors_allowed_origins</b> &nbsp Specifies the domains that are permitted to access the API, including development hosts. Use '*' to allow all or specify domains individually, separated by a comma.</div>
            <div class="fform__group ">
                <input class="fform__input " :class="{'fform__input--active':edit_key=='cors_allowed_origins'}" id="cors_allowed_origins" type="text" @blur="edit_key=null" 
                    v-model.trim="cors_allowed_origins" :disabled="edit_key != 'cors_allowed_origins'" placeholder="Enter new default longitude to change, else leave empty">
                <svg @click="edit_key='cors_allowed_origins'" class="fform__icon" :class="{'fform__icon--active':edit_key=='cors_allowed_origins'}"> <use xlink:href="@/assets/svg/sprite.svg#icon-pencil"></use></svg>
            </div>
        </div>

        <div class="fform__row mt-10 " @click="clearMessage" :class="{ 'fform__disabled': edit_key != 'cors_allowed_origins' }">
            <div class="fform__description"><b>login_page_title</b> &nbsp Specifies the domains that are permitted to access the API, including development hosts. Use '*' to allow all or specify domains individually, separated by a comma.</div>
            <div class="fform__group ">
                <input class="fform__input " :class="{'fform__input--active':edit_key=='login_page_title'}" id="login_page_title" type="text" @blur="edit_key=null" 
                    v-model.trim="login_page_title" :disabled="edit_key != 'login_page_title'" placeholder="The HTML-formatted title text displayed on the login page of the IoTrack Pro application.">
                <svg @click="edit_key='login_page_title'" class="fform__icon" :class="{'fform__icon--active':edit_key=='login_page_title'}"> <use xlink:href="@/assets/svg/sprite.svg#icon-pencil"></use></svg>
            </div>
        </div>

        <button  class="bbtn bbtn--red mt-14"  @click.prevent="initUpdateSettings()">
            Update Settings
        </button>
    </form>

</template>

<!-- --------------------------------------------------------------- -->

<script setup>
import { useAppStore } from '@/stores/appStore';
import { useMessageStore } from '@/stores/messageStore';
import { onMounted, ref, watch } from 'vue';
import AdminConfirmationModal from '@/components/commen/AdminConfirmationModal.vue';

const edit_key = ref(null);

const google_api_key = ref('');
const google_map_id = ref(null);
const jwt_expiration_seconds = ref(null);
const redis_ttl_seconds = ref(null);
const device_access_mode = ref(null);
const initial_parking_check_date = ref(null);
const cors_allowed_origins = ref(null);
const login_page_title = ref(null);


const adminModalOn = ref(false);

const appStore = useAppStore();
const messageStore = useMessageStore();

watch(edit_key, (val)=>{
    if (val == null) return;
    const inputEl = document.querySelector('#'+val);
    setTimeout(()=>{inputEl.focus();}, 200)       
})

function clearMessage() {
    messageStore.clearFlashMessage();
}

function initUpdateSettings() {
    messageStore.clearFlashMessage();
    const message = []; // Clear previous messages
    let hasError = false;

    // Validate if jwt_expiration_seconds is a valid positive integer
    if (!Number.isInteger(parseFloat(jwt_expiration_seconds.value)) || jwt_expiration_seconds.value <= 0) {
        message.push("JWT expiration seconds must be a valid positive integer.");
        hasError = true;
    }

    // Validate if redis_ttl_seconds is a valid positive integer
    if (!Number.isInteger(parseFloat(redis_ttl_seconds.value)) || redis_ttl_seconds.value <= 0) {
        message.push("Redis TTL seconds must be a valid positive integer.");
        hasError = true;
    }

    // Validate if device_access_mode is either 'black_list' or 'white_list'
    if (device_access_mode.value !== 'black_list' && device_access_mode.value !== 'white_list') {
        message.push("Device access mode must be either 'black_list' or 'white_list'.");
        hasError = true;
    }

    // Validate if google_map_id is null, empty, or a valid format
    if (google_map_id.value == null || google_map_id.value.trim() === '' || !/^[A-Za-z0-9_\-]{10,}$/.test(google_map_id.value)) {
        message.push("Google Map ID must be at least 10 characters long, containing letters, numbers, dashes, or underscores.");
        hasError = true;
    }

    // Validate if google_api_key is null, empty, or a valid format
    if (google_map_id.value !== null && google_api_key.value.trim() !== '' && !/^[A-Za-z0-9_\-]{10,}$/.test(google_api_key.value)) {
        message.push("Google MAP ID must be valid id with at least 20 characters, containing letters, numbers, dashes, or underscores.");
        hasError = true;
    }

    // Validate initial_parking_check_date for a valid ISO 8601 date format
    if (initial_parking_check_date.value && !/^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}Z$/.test(initial_parking_check_date.value)) {
        message.push("Initial parking check date must be in ISO 8601 format (YYYY-MM-DDTHH:MM:SSZ).");
        hasError = true;
    }

    // Validate CORS Allowed Origins
    if (!cors_allowed_origins.value || !cors_allowed_origins.value.split(',').every(origin => /^(\*|https?:\/\/[\w\-\.]+(:\d+)?(\/)?)$/.test(origin.trim()))) {
        message.push("CORS Allowed Origins must be a comma-separated list of valid URLs or '*'.");
        hasError = true;
    }

    // Validate login_page_title for non-empty and safe HTML content
    if (!login_page_title.value.trim() || /<script|<\/script>/i.test(login_page_title.value)) {
        message.push("Login page title must be provided and not include potentially dangerous content such as scripts.");
        hasError = true;
    }

    if (hasError) {
        messageStore.setFlashMessages(message);
        messageStore.setFlashClass("flash-message--yellow");
        return
    }

    adminModalOn.value = true;
}

async function updateSettings(payload) {
    adminModalOn.value = false;
    try {
        payload.admin_password = payload.adminPassword;
        delete payload.adminPassword;

        if (google_api_key.value !== null && google_api_key.value.trim() !== '') {
            payload.google_api_key = google_api_key.value;
        }

        payload.jwt_expiration_seconds = String(jwt_expiration_seconds.value);
        payload.redis_ttl_seconds = String(redis_ttl_seconds.value);
        payload.device_access_mode = device_access_mode.value;
        payload.google_map_id = google_map_id.value;
        payload.initial_parking_check_date = initial_parking_check_date.value;
        payload.cors_allowed_origins = cors_allowed_origins.value;
        payload.login_page_title = login_page_title.value;

        const response = await appStore.updateSettings(payload);

        if (response?.status == 200) {
            const msg = response.data?.message ?? "Settings updated successfully.";
            messageStore.setFlashMessages([msg], "flash-message--green");
     
            // update settings

        }

    } catch (error) {
        console.error("! RootLevelSettings.updateSettings !\n", error);
        const errMsg = error.response?.data ?? "Failed to update settings"
        messageStore.setFlashMessages([errMsg], "flash-message--red");

    } finally {

    }
}

onMounted(()=>{
    const settings = appStore.getAppSettings
    google_map_id.value = settings.google_map_id;
    jwt_expiration_seconds.value = settings.jwt_expiration_seconds;
    redis_ttl_seconds.value = settings.redis_ttl_seconds;
    device_access_mode.value = settings.device_access_mode;
    initial_parking_check_date.value = settings.initial_parking_check_date;
    cors_allowed_origins.value = settings.cors_allowed_origins;
    login_page_title.value = settings.login_page_title;
});

</script>

<!-- --------------------------------------------------------------- -->

<style lang="scss" scoped>
.fform__row {
    flex-direction: column;
}

.fform__description {
    font-family: $font-primary;
    line-height: 1.2rem;
    margin-bottom: -0.5rem;
}

.fform__input {
    height: 2.5rem !important;
    padding: 0.5rem 0.5rem 0.5rem 0.5rem;

    &--active {
        border-color: $col-blue-600 !important;;
    }
}

.fform__icon {
    width: 2.5rem;
    height: 2.5rem;
    color: $col-zinc-300 !important;
    fill: $col-white;
    background-color: currentColor;
    padding: .5rem;
    border: 1px solid currentColor;
    border-radius: 5px;
    position: absolute;
    top: 0;
    right: 0;
    cursor: pointer;

    &:hover, &--active {
        color: $col-blue-600 !important;
        fill: currentColor;
        background-color: $col-white;
    }

    &:active {    
        fill: $col-white;
        background-color: $col-blue-600 !important;
    }
}

.fform__select {
    cursor: pointer;
}



#google_api_key::placeholder {
    font-size: .8rem;  /* Example font size */
    text-wrap: wrap;
    width: 200px;
    transform: translate(0,-9px);
    color: $col-red-600;
    line-height: .8rem !important;


    @include respondMobile(445) {
        font-size: .8rem;
        width: 100%;
        transform: translate(0,0);
    }
    @include respondMobile(462) { font-size: .85rem; }
    @include respondMobile(480) { font-size: .9rem; }
    @include respondMobile(520) { font-size: 1rem; }
}

</style>