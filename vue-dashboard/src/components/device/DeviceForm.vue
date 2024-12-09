<template>
    <LocationModal 
        v-if="locationModalOpen"
        :modalIsOpen="locationModalOpen" 
        :markerPosition="{lat:latitude, lng:longitude}" 
        @emitCloseModal="locationModalOpen = false"
        @emitMarkerPosition="updateMarkerPosition"
    ></LocationModal>

    <form class="fform" autocomplete="off">
        <div class="fform__row mt-8 " @click="clearMessage" :class="{ 'fform__disabled': confirmOn }">
            <div class="fform__group ">
                <label class="fform__label" for="device_id">Device Id <span class="fform__required">*</span></label>
                <input class="fform__input" id="device_id" type="text" placeholder="Enter device ID"
                    v-model.trim="deviceId" :disabled="confirmOn">
            </div>
            <div class="fform__group ">
                <label class="fform__label" for="name">Device Name <span class="fform__required">*</span></label>
                <input class="fform__input" id="name" type="text" placeholder="Enter device name"
                    v-model.trim="name" :disabled="confirmOn">
            </div>
        </div>

        <div class="fform__row mt-8 " @click="clearMessage" :class="{ 'fform__disabled': confirmOn }">

            <div class="fform__row">
                
                <div class="fform__group ">
                    <TheSelector 
                        :options="returnNetworkType" 
                        :selectedOption="selectedOptions['networkType']"
                        fieldName="networkType" 
                        label="Network Type" 
                        :isDisabled="confirmOn" :isRequired="true"
                        @emitOption="selectedOptions['networkType'] = $event">
                    </TheSelector>
                </div>

                <div class="fform__group ">
                    <label class="fform__label" for="firmware_version">Firmware Version </label>
                    <input class="fform__input" 
                    id="firmware_version" 
                    type="text"  
                    placeholder="Enter the firmware version"
                    :value="formattedFirmwareVersion" 
                    @focusout="formatFirmwareVersion" 
                    :disabled="confirmOn">
                </div>

            </div>            

            <div class="fform__row">

                <div class="fform__row">
                    <div class="fform__group ">
                        <label class="fform__label" for="latitude">Latitude<span
                                class="fform__required">*</span></label>
                        <input class="fform__input" id="latitude" type="text" placeholder="Enter the latitude"
                            v-model.trim="latitude" :disabled="confirmOn">
                    </div>

                    <div class="fform__group ">
                        <label class="fform__label" for="longitude">Longitude<span
                                class="fform__required">*</span></label>
                        <input class="fform__input" id="longitude" type="text" placeholder="Enter the longitude"
                            v-model.trim="longitude" :disabled="confirmOn">

                        <svg @click="locationModalOpen = true" class="location-btn"><use xlink:href="@/assets/svg/sprite.svg#icon-google-maps-2"></use></svg>
                    </div>

                </div>
            </div>
            
        </div>
        <div class="checkbox__container mt-4">
            <TheCheckbox :is-checked="isWhiteListed" @emit-checkbox="isWhiteListed = !isWhiteListed">Is White Listed</TheCheckbox>
            <TheCheckbox :is-checked="isBlackListed" @emit-checkbox="isBlackListed = !isBlackListed">Is Black Listed</TheCheckbox>
            <TheCheckbox :is-checked="isHidden" @emit-checkbox="isHidden = !isHidden">Is Hidden</TheCheckbox>            
        </div>

        <transition name="fade" mode="out-in">
            <button v-if="!confirmOn" class="bbtn bbtn--blue mt-8"
                @click.prevent="initCreateDevice()" key="create-button">
                Create Device
            </button>

            <div v-else class="bbtn__row mt-8" key="confirm-buttons">
                <button class="bbtn bbtn--zinc-lt" @click.prevent="confirmOn = false">Cancel</button>
                <button class="bbtn bbtn--blue" @click.prevent="createDevice">Confirm</button>
            </div>
        </transition>
    </form>
</template>

<!-- --------------------------------------------------------------- -->

<script setup>
import { useMessageStore } from '@/stores/messageStore';
import TheSelector from '@/components/commen/TheSelector.vue'
import TheCheckbox from '../commen/TheCheckbox.vue';
import { computed, onMounted, reactive, ref, watch } from 'vue';
import LocationModal from '../commen/LocationModal.vue';
import { useAppStore } from '@/stores/appStore';
import { storeToRefs } from 'pinia';
import { useDeviceStore } from '@/stores/deviceStore';



// - Store -------------------------------------------------------------
const messageStore = useMessageStore();

const appStore = useAppStore();
const { getDefaultLatitude, getDefaultLongitude } = storeToRefs(appStore);

const deviceStore = useDeviceStore();

// - Data --------------------------------------------------------------
const confirmOn = ref(false);
const locationModalOpen = ref(false)

const deviceId = ref('');
const name = ref('');
const latitude = ref('');
const longitude = ref('');
const isBlackListed = ref(false);
const isWhiteListed = ref(false);
const isHidden = ref(false);

const networkType = ref([
    // {id: 0, name: 'Root'},
    { id: 'NB-IoT', name: 'NB-IoT' },
    { id: 'LoRa', name: 'LoRa' },
    { id: 'SigFox', name: 'SigFox' },
]);

const selectedOptions = reactive({
    'networkType': { _key: 'NB-IoT', _value: 'NB-IoT' }
});

// - firmwareVersion ----------------------------------------------------

const firmwareVersion = ref('0');

// Computed property to format the firmware version for display
const formattedFirmwareVersion = computed(() => {
    return parseFloat(firmwareVersion.value).toFixed(1);
});

// Method to handle input and format to one decimal place
function formatFirmwareVersion(event) {
    let value = event.target.value;
    if (value === '' || isNaN(value)) {
        firmwareVersion.value = '';
    } else {
        // Restrict to one decimal point
        let floatVal = parseFloat(value);
        firmwareVersion.value = floatVal.toFixed(1);
    }
}

// - Computed ----------------------------------------------------------

const returnNetworkType = computed(() => {
    return networkType.value.map((net => {
        // TODO: use utilStore.capitalizeFirstLetter() for value.
        return { ...net, _key: net.id, _value: net.name }
    }))
});



// - Watchers ----------------------------------------------------------

watch( locationModalOpen, (val)=>{
    appStore.setPageScrollDisabled(val);
});


// - Methods -----------------------------------------------------------





function clearMessage() {
    messageStore.clearFlashMessage();
}

function resetForm() {
    // Reset all form fields
    deviceId.value = '';
    name.value = '';
    firmwareVersion.value = '0';
    latitude.value = getDefaultLatitude.value ? Number(getDefaultLatitude.value) : 0; 
    longitude.value = getDefaultLongitude.value ? Number(getDefaultLongitude.value) : 0; 
    isWhiteListed.value = false;
    isBlackListed.value = false;
    isHidden.value = false;

    selectedOptions.networkType = { _key: 'NB-IoT', _value: 'NB-IoT' };
}


function initCreateDevice() {
    // Clear previous messages
    messageStore.clearFlashMessage();

    // Prepare an array to store error messages
    const errors = [];

    // Validate that the deviceId is not empty
    if (!deviceId.value.trim()) {
        errors.push("Device ID is required.");
    }

    // Validate that the name is not empty
    if (!name.value.trim()) {
        errors.push("Device Name is required.");
    }

    // Validate that a valid network type is selected
    if (!selectedOptions.networkType || !networkType.value.some(nt => nt.id === selectedOptions.networkType._key)) {
        errors.push("A valid Network Type must be selected.");
    }

    // Validate the firmware version to be a valid number
    if (!firmwareVersion.value || isNaN(firmwareVersion.value) || firmwareVersion.value.trim() === '') {
        errors.push("Firmware Version must be a valid number.");
    }

    // If there are any errors, display them and do not proceed
    if (errors.length > 0) {
        messageStore.setFlashMessages(errors, 'flash-message--yellow'); // Update to use your actual method for displaying errors
        return; // Stop the function if there are errors
    }

    // If no errors, proceed
    confirmOn.value = true;
}

async function createDevice() {
    try {
        // Prepare the payload ensuring all fields are correctly referenced
        const payload = {
            device_id: deviceId.value,  
            name: name.value,           
            network_type: selectedOptions.networkType._key, 
            firmware_version: Number(firmwareVersion.value), 
            latitude: latitude.value,   
            longitude: longitude.value, 
            is_allowed: isWhiteListed.value, 
            is_blocked: isBlackListed.value, 
            is_hidden: isHidden.value   
        };

        console.log(payload)

        // Make the API call to create the device
        const response = await deviceStore.createDevice(payload);
        console.log(response)

        if (response.status == 200) {
            const msg = response.data?.message ?? "Device created successfully.";
            messageStore.setFlashMessages([msg], "flash-message--green");
            resetForm();
            // deviceStore.pushUserToList(response.data?.user);

        }


    } catch (error) {
        console.error("! DeviceForm.createDevice !\n", error);
        const errMsg = error.response?.data ?? "Failed to create device";
        messageStore.setFlashMessages([errMsg], "flash-message--red");

    } finally {
        confirmOn.value = false;
    }
}

function updateMarkerPosition(latLng) {
    latitude.value = latLng.lat;
    longitude.value = latLng.lng;
}


onMounted(()=>{
    latitude.value = getDefaultLatitude.value ? Number(getDefaultLatitude.value) : 0; 
    longitude.value = getDefaultLongitude.value ? Number(getDefaultLongitude.value) : 0; 
})

</script>

<!-- --------------------------------------------------------------- -->

<style lang="scss" scoped>
.checkbox__container {
    display: flex;
    flex-wrap: wrap;
    column-gap: 2rem;
    
}

.location-btn {
    cursor: pointer;
    width: 2.5rem;
    height: 2.5rem;
    transition: all .1s ease;
    position: absolute;
    top: 0.8rem;
    right: 0.2rem;
  

    &:hover {
        transform: scale(1.1);
    }
}
.no-scrolling {
    overflow-x:hidden;
}
</style>