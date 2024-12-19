<template>
    <form class="fform" autocomplete="off">
        <transition name="modal" appear>
            <AdminConfirmationModal v-if="adminModalOn" 
                @emitCancel="adminModalOn = false" 
                @emitConfirm="updateSettings"
                appear
                >
            </AdminConfirmationModal>
        </transition>

        <div class="fform__row mt-20 " @click="clearMessage" :class="{ 'fform__disabled': edit_key != 'default_latitude' }">
            <div class="fform__description"><b>default_latitude</b> &nbsp Default latitude for map centering and initial device placement on the map.</div>
            <div class="fform__group ">
                <input class="fform__input " :class="{'fform__input--active':edit_key=='default_latitude'}" id="default_latitude" type="text" @blur="edit_key=null" 
                    v-model.trim="default_latitude" :disabled="edit_key != 'default_latitude'" placeholder="Enter new default latitude to change, else leave empty">
                <svg @click="edit_key='default_latitude'" class="fform__icon" :class="{'fform__icon--active':edit_key=='default_latitude'}"> <use xlink:href="@/assets/svg/sprite.svg#icon-pencil" ></use></svg>
            </div>
        </div>


        <div class="fform__row mt-10 " @click="clearMessage" :class="{ 'fform__disabled': edit_key != 'default_longitude' }">
            <div class="fform__description"><b>default_longitude</b> &nbsp Default longitude for map centering and initial device placement on the map.</div>
            <div class="fform__group ">
                <input class="fform__input " :class="{'fform__input--active':edit_key=='default_longitude'}" id="default_longitude" type="text" @blur="edit_key=null" 
                    v-model.trim="default_longitude" :disabled="edit_key != 'default_longitude'" placeholder="Enter new default longitude to change, else leave empty">
                <svg @click="edit_key='default_longitude'" class="fform__icon" :class="{'fform__icon--active':edit_key=='default_longitude'}"> <use xlink:href="@/assets/svg/sprite.svg#icon-pencil"></use></svg>
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


const adminModalOn = ref(false);

const edit_key = ref(null);
const default_latitude = ref(null);
const default_longitude = ref(null);
const google_map_id = ref(null);

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

function initUpdateSettings () {
    messageStore.clearFlashMessage();
    const message = []; // Clear previous messages
    let hasError = false;

    // Validate if default_latitude is a number
    if (isNaN(parseFloat(default_latitude.value)) || !isFinite(default_latitude.value)) {
        message.push("Default latitude must be a valid number.");
        hasError = true;
    }

    // Validate if default_longitude is a number
    if (isNaN(parseFloat(default_longitude.value)) || !isFinite(default_longitude.value)) {
        message.push("Default longitude must be a valid number.");
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


        payload.default_latitude = String(default_latitude.value);
        payload.default_longitude = String(default_longitude.value);
  

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
    default_latitude.value = settings.default_latitude;
    default_longitude.value = settings.default_longitude;
    google_map_id.value = settings.google_map_id;
});
(()=>{
    clearMessage();
})();

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

</style>