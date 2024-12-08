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
                    v-model.trim="device_id" :disabled="confirmOn">
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
                    <label class="fform__label" for="firmware_version">Firmware Version <span class="fform__required">*</span></label>
                    <input class="fform__input" id="firmware_version" type="text" placeholder="Enter the firmware version"
                        v-model.trim="firmware_version" :disabled="confirmOn">
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



// - Store -------------------------------------------------------------
const messageStore = useMessageStore();
const appStore = useAppStore();
const { getDefaultLatitude, getDefaultLongitude } = storeToRefs(appStore);

// - Data --------------------------------------------------------------
const confirmOn = ref(false);
const locationModalOpen = ref(false)

const device_id = ref('');
const name = ref('');
const firmware_version = ref('');
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

function initCreateDevice() {
    confirmOn.value = true;
}

function createDevice() {

}

function updateMarkerPosition(latLng) {
    latitude.value = latLng.lat;
    longitude.value = latLng.lng;
}

onMounted(()=>{
    latitude.value = Number(getDefaultLatitude.value)
    longitude.value = Number(getDefaultLongitude.value)
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