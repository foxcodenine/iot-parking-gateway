<template>

    <LocationModal v-if="locationModalOpen" :modalIsOpen="locationModalOpen"
        :markerPosition="{ lat: selectedDevice.latitude, lng: selectedDevice.longitude }"
        @emitCloseModal="locationModalOpen = false" @emitMarkerPosition="updateMarkerPosition"></LocationModal>

    <input class="ttable__search mt-8" v-model="searchTerm" type="text" placeholder="Search...">

    <div class="ttable__container">

        <table class="ttable  mt-8" @click="clearMessage">
            <thead>
                <tr>
                    <th class="w-52">
                        <span class="cursor-pointer" @click="sortTable('device_id')">
                            Device ID
                        </span>
                        <svg class="t-sort-arrow" :class="sortArrow('device_id')">
                            <use xlink:href="@/assets/svg/sprite.svg#triangle-1"></use>
                        </svg>
                    </th>
                    <th class="w-52">
                        <span class="cursor-pointer" @click="sortTable('name')">
                            Name
                        </span>
                        <svg class="t-sort-arrow" :class="sortArrow('name')">
                            <use xlink:href="@/assets/svg/sprite.svg#triangle-1"></use>
                        </svg>
                    </th>
                    <th class="w-24">
                        <span class="cursor-pointer" @click="sortTable('network_type')">
                            Network
                        </span>
                        <svg class="t-sort-arrow" :class="sortArrow('network_type')">
                            <use xlink:href="@/assets/svg/sprite.svg#triangle-1"></use>
                        </svg>
                    </th>
                    <th class="w-16">
                        <span class="cursor-pointer" @click="sortTable('firmware_version')">
                            Firmware
                        </span>
                        <svg class="t-sort-arrow" :class="sortArrow('firmware_version')">
                            <use xlink:href="@/assets/svg/sprite.svg#triangle-1"></use>
                        </svg>
                    </th>
                    <th class="w-8"></th>
                    <th class="w-24">
                        <span class="cursor-pointer" @click="sortTable('latitude')">
                            Latitude
                        </span>
                        <svg class="t-sort-arrow" :class="sortArrow('latitude')">
                            <use xlink:href="@/assets/svg/sprite.svg#triangle-1"></use>
                        </svg>
                    </th>
                    <th class="w-24">
                        <span class="cursor-pointer" @click="sortTable('longitude')">
                            Longitude
                        </span>
                        <svg class="t-sort-arrow" :class="sortArrow('longitude')">
                            <use xlink:href="@/assets/svg/sprite.svg#triangle-1"></use>
                        </svg>
                    </th>
                    <th class="w-20">
                        <span class="cursor-pointer" @click="sortTable('is_occupied')">
                            Occupied
                        </span>
                        <svg class="t-sort-arrow" :class="sortArrow('is_occupied')">
                            <use xlink:href="@/assets/svg/sprite.svg#triangle-1"></use>
                        </svg>
                    </th>
                    <th class="w-20" v-if="getAppSettings?.device_access_mode != 'black_list'">
                        <span class="cursor-pointer" @click="sortTable('is_allowed')">
                            White List
                        </span>
                        <svg class="t-sort-arrow" :class="sortArrow('is_allowed')">
                            <use xlink:href="@/assets/svg/sprite.svg#triangle-1"></use>
                        </svg>
                    </th>
                    <th class="w-24" v-else>
                        <span class="cursor-pointer" @click="sortTable('is_blocked')">
                            Black List
                        </span>
                        <svg class="t-sort-arrow" :class="sortArrow('is_blocked')">
                            <use xlink:href="@/assets/svg/sprite.svg#triangle-1"></use>
                        </svg>
                    </th>
                    <th class="w-24">
                        <span class="cursor-pointer" @click="sortTable('is_hidden')">
                            Hidden
                        </span>
                        <svg class="t-sort-arrow" :class="sortArrow('is_hidden')">
                            <use xlink:href="@/assets/svg/sprite.svg#triangle-1"></use>
                        </svg>
                    </th>
                    <th class="w-24">
                        <span class="cursor-pointer" @click="sortTable('created_at')">
                            Registered
                        </span>
                        <svg class="t-sort-arrow" :class="sortArrow('created_at')">
                            <use xlink:href="@/assets/svg/sprite.svg#triangle-1"></use>
                        </svg>
                    </th>
                    <th v-if="getUserAccessLevel <= 2">
                    </th>
                </tr>
            </thead>

            <tbody>


                <tr  v-for="( device, _ ) in getDevicesList"
                    :class="{ 'bg-lime-200': selectedDevice.device_id == device.device_id }">

                    <td :id="`col-${device.device_id}`">{{ device.device_id }}</td>

                    <td v-if="action == 'edit' && selectedDevice.device_id == device.device_id">
                        <input class="ttable__input w-48" type="text" v-model="selectedDevice.name">
                    </td>
                    <td v-else>{{ device.name }}</td>

                    <td>{{ device.network_type }}</td>

                    <td v-if="action == 'edit' && selectedDevice.device_id == device.device_id">
                        <input class="ttable__input w-12" type="text" v-model="selectedDevice.firmware_version">
                    </td>
                    <td v-else>{{ device.firmware_version }}</td>

                    <td class="w-8">
                        <svg v-if="action == 'edit' && selectedDevice.device_id == device.device_id"
                            @click="locationModalOpen = true" class="ttable__location-btn ">
                            <use xlink:href="@/assets/svg/sprite.svg#icon-google-maps-2"></use>
                        </svg>
                    </td>

                    <td v-if="action == 'edit' && selectedDevice.device_id == device.device_id">
                        <input class="ttable__input w-20" type="text" v-model="selectedDevice.latitude">
                    </td>
                    <td v-else>{{ device.latitude }}</td>

                    <td v-if="action == 'edit' && selectedDevice.device_id == device.device_id">
                        <input class="ttable__input w-20" type="text" v-model="selectedDevice.longitude">
                    </td>
                    <td v-else>{{ device.longitude }}</td>

                    <td @click="toggelOccupied(device)" class="ps-7">
                        <div v-if="action == 'edit' && selectedDevice.device_id == device.device_id"
                            class="circle__outer circle__active"
                            :class="{ 'circle__occupied': selectedDevice.is_occupied, 'circle__vacant': !selectedDevice.is_occupied }">
                            <div class="circle__inner">
                                <p>P</p>
                            </div>
                        </div>
                        <div v-else class="circle__outer"
                            :class="{ 'circle__occupied': device.is_occupied, 'circle__vacant': !device.is_occupied }">
                            <div class="circle__inner">
                                <p>P</p>
                            </div>
                        </div>
                    </td>

                    <td @click="toggelAllowed(device)" class="ps-6"
                        v-if="getAppSettings?.device_access_mode != 'black_list'">
                        <div v-if="action == 'edit' && selectedDevice.device_id == device.device_id"
                            class="circle__outer circle__active"
                            :class="{ 'circle__allowed': selectedDevice.is_allowed }">
                            <div class="circle__inner">
                                <p v-if="selectedDevice.is_allowed">W</p>
                            </div>
                        </div>
                        <div v-else class="circle__outer" :class="{ 'circle__allowed': device.is_allowed }">
                            <div class="circle__inner">
                                <p v-if="device.is_allowed">W</p>
                            </div>
                        </div>
                    </td>

                    <td @click="toggelBlocked(device)" class="ps-6" v-else>
                        <div v-if="action == 'edit' && selectedDevice.device_id == device.device_id"
                            class="circle__outer circle__active"
                            :class="{ 'circle__blocked': selectedDevice.is_blocked }">
                            <div class="circle__inner">
                                <p v-if="selectedDevice.is_blocked">B</p>
                            </div>
                        </div>
                        <div v-else class="circle__outer" :class="{ 'circle__blocked': device.is_blocked }">
                            <div class="circle__inner">
                                <p v-if="device.is_blocked">B</p>
                            </div>
                        </div>
                    </td>

                    <td @click="toggelHidden(device)" class="ps-5">
                        <div v-if="action == 'edit' && selectedDevice.device_id == device.device_id"
                            class="circle__outer circle__active"
                            :class="{ 'circle__hidden': selectedDevice.is_hidden }">
                            <div class="circle__inner">
                                <p v-if="selectedDevice.is_hidden">H</p>
                            </div>
                        </div>
                        <div v-else class="circle__outer" :class="{ 'circle__hidden': device.is_hidden }">
                            <div class="circle__inner">
                                <p v-if="device.is_hidden">H</p>
                            </div>
                        </div>
                    </td>

                    <td>{{ formatDateUtil(device.created_at) }}</td>

                    <td v-if="getUserAccessLevel <= 2">
                        <div class="t-btns ml-auto" v-if="selectedDevice.device_id != device.device_id">
                            <a class="t-btns__btn " @click="initEditDevice(device)">
                                <svg class="t-btns__icon">
                                    <use xlink:href="@/assets/svg/sprite.svg#icon-pencil"></use>
                                </svg>
                            </a>
                            <a class="t-btns__btn" @click="initDeleteDevice(device)">
                                <svg class="t-btns__icon">
                                    <use xlink:href="@/assets/svg/sprite.svg#icon-trash-can-solid2"></use>
                                </svg>
                            </a>
                        </div>
                        <div class="t-btns ml-auto" v-else>
                            <a class="t-btns__btn t-btns__btn--yes" @click="editOrDelete"
                                v-tooltip="{ content: action == 'edit' ? 'Edit Device' : 'Delete Device' }">
                                Yes
                            </a>
                            <a class="t-btns__btn t-btns__btn--no" @click="resetSelectedDevice"
                                v-tooltip="{ content: action == 'edit' ? 'Cancel Edit' : 'Cancel Delete' }">
                                No
                            </a>
                        </div>
                    </td>
                </tr>

            </tbody>
        </table>

    </div>
</template>

<!-- --------------------------------------------------------------- -->

<script setup>
import { useAuthStore } from '@/stores/authStore';
import { useDeviceStore } from '@/stores/deviceStore';
import { useMessageStore } from '@/stores/messageStore';
import { storeToRefs } from 'pinia';
import { vTooltip } from 'floating-vue'
import 'floating-vue/dist/style.css';

import { onMounted, ref, reactive, watch, computed } from 'vue';
import { useAppStore } from '@/stores/appStore';
import { formatDateUtil } from '@/utils/dateTimeUtils';

import LocationModal from '../commen/LocationModal.vue';
import { useRouter } from 'vue-router';


const router = useRouter();
// - Store -------------------------------------------------------------
const messageStore = useMessageStore();

const deviceStore = useDeviceStore();

const authStore = useAuthStore();
const { getUserAccessLevel } = storeToRefs(authStore)

const appStore = useAppStore();
const { getAppSettings } = storeToRefs(appStore);

// - Props -------------------------------------------------------------

const props = defineProps({
    deviceID: {
        type: String,
        required: false,
    }
});



// -- Data -------------------------------------------------------------


const locationModalOpen = ref(false);

const action = ref(null);

const selectedDevice = reactive({
    device_id: null,
    name: null,
    network_type: null,
    firmware_version: null,
    latitude: null,
    longitude: null,
    is_occupied: null,
    is_allowed: null,
    is_blocked: null,
    is_hidden: null,
});

const sortBy = ref('created_at');
const sortDesc = ref('true');
const searchTerm = ref("");






// -- Computed ---------------------------------------------------------

const getDevicesList = computed(() => {
    let list = [...Object.values(deviceStore.getDevicesList)];

    list = list.filter(item => {
        return (
            item.device_id?.toLowerCase().includes(searchTerm.value.toLowerCase().trim()) ||
            item.name?.toLowerCase().includes(searchTerm.value.toLowerCase().trim()) ||
            item.network_type?.toLowerCase().includes(searchTerm.value.toLowerCase().trim()) ||
            String(item.firmware_version)?.toLowerCase().includes(searchTerm.value.toLowerCase().trim())
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


// - Watchers ----------------------------------------------------------

watch(locationModalOpen, (val) => {
    appStore.setPageScrollDisabled(val);
});


/**
 * Watches the `deviceID` prop for changes and performs the following actions:
 * 
 * - If `deviceID` is valid:
 *   1. Waits for a brief delay (250ms) to ensure the DOM is updated.
 *   2. Fetches the device data from the store and initializes the edit process (`initEditDevice`).
 *   3. Scrolls smoothly to the element with the corresponding `id` (`col-${deviceID}`).
 * 
 * - Navigates to the 'deviceView' route using Vue Router, ensuring the user is redirected to the device management view.
 * 
 * This implementation serves as a workaround for scenarios where clicking "Edit" on a device's InfoWindow
 * in the MapView does not correctly select the device to edit, but updates the URL to reflect the device being edited.
 * By selecting the device to edit from the `deviceID` in the URL, updating the view from `deviceEditView` to `deviceView`,
 * and scrolling to the target element, this ensures a seamless user experience.
 * 
 * Note: The 250ms timeout ensures that DOM changes (e.g., list rendering) are complete before scrolling.
 */
watch(() => props.deviceID, async (deviceID) => {
    if (!deviceID) return

    setTimeout(() => {
        const device = deviceStore.getDevicesList[deviceID];
        // Select Device to edit
        initEditDevice(device);
        // Scroll to the element 
        const element = document.getElementById(`col-${deviceID}`);
        if (element) {
            element.scrollIntoView({ behavior: 'smooth' });
        }
    }, 250)

    router.push({ name: 'deviceView' });

}, {
    immediate: true
});

// - Methods -----------------------------------------------------------

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

function resetSelectedDevice() {
    action.value = null;
    selectedDevice.device_id = null;
    selectedDevice.name = null;
    selectedDevice.network_type = null;
    selectedDevice.firmware_version = null;
    selectedDevice.latitude = 0;
    selectedDevice.longitude = 0;
    selectedDevice.is_occupied = null;
    selectedDevice.is_allowed = null;
    selectedDevice.is_blocked = null;
    selectedDevice.is_hidden = null;
}

function initEditDevice(d) {
    action.value = "edit";
    selectedDevice.device_id = d.device_id;
    selectedDevice.name = d.name;
    selectedDevice.network_type = d.network_type;
    selectedDevice.firmware_version = d.firmware_version;
    selectedDevice.latitude = d.latitude;
    selectedDevice.longitude = d.longitude;
    selectedDevice.is_occupied = d.is_occupied;
    selectedDevice.is_allowed = d.is_allowed;
    selectedDevice.is_blocked = d.is_blocked;
    selectedDevice.is_hidden = d.is_hidden;
}

function initDeleteDevice(d) {
    action.value = "delete";
    selectedDevice.device_id = d.device_id;
}

function editOrDelete() {
    action.value == "edit" ? editDevice() : deleteDevice();
}

async function editDevice() {
    try {
        const response = await deviceStore.updateDevice(selectedDevice);

        if (response?.status == 200) {
            const msg = response.data?.message ?? "Device updated successfully.";
            messageStore.setFlashMessages([msg], "flash-message--green");

            if (response.data?.device) {
                console.log('A')
                deviceStore.updateDeviceInList(response.data.device);
            }
        }

    } catch (error) {
        console.error("! UserForm.updateUser !\n", error);
        const errMsg = error.response?.data ?? "Failed to update device"
        messageStore.setFlashMessages([errMsg], "flash-message--red");
    } finally {
        resetSelectedDevice();
    }
}

async function deleteDevice() {
    try {

        const response = await deviceStore.deleteDevice({
            device_id: selectedDevice.device_id,
        });

        if (response?.status == 200) {
            const msg = response.data?.message ?? "Device deleted successfully.";
            deviceStore.removeDeviceFromList(selectedDevice.device_id);
            messageStore.setFlashMessages([msg], "flash-message--green");
        }

    } catch (error) {
        console.error("! DeviceForm.deleteDevice !\n", error);
        const errMsg = error.response?.data ?? "Failed to delete device"
        messageStore.setFlashMessages([errMsg], "flash-message--red");
    } finally {
    }
}

function toggelOccupied(device) {
    if (action.value == 'edit' && device.device_id == selectedDevice.device_id) {

        selectedDevice.is_occupied = !selectedDevice.is_occupied;
    }
}

function toggelAllowed(device) {
    if (action.value == 'edit' && device.device_id == selectedDevice.device_id) {

        selectedDevice.is_allowed = !selectedDevice.is_allowed;
    }
}

function toggelBlocked(device) {
    if (action.value == 'edit' && device.device_id == selectedDevice.device_id) {

        selectedDevice.is_blocked = !selectedDevice.is_blocked;
    }
}

function toggelHidden(device) {
    if (action.value == 'edit' && device.device_id == selectedDevice.device_id) {

        selectedDevice.is_hidden = !selectedDevice.is_hidden;
    }
}

function updateMarkerPosition(latLng) {
    selectedDevice.latitude = latLng.lat;
    selectedDevice.longitude = latLng.lng;
}

</script>

<!-- --------------------------------------------------------------- -->

<style lang="scss" scoped>
.t-btns {
    justify-content: end !important;
}

.ttable {
    min-width: 79rem;
}

.ttable__input,
.ttable__select {
    background-color: rgba($col-white, .7) !important;
    padding: .2rem .3rem !important;
    height: 1.8rem !important;
}

.ttable__location-btn {
    cursor: pointer;
    width: 1.3rem;
    height: 1.3rem;
    transition: all .1s ease;
    margin-bottom: -4px;
    margin-right: -6px;

    &:hover {
        transform: scale(1.1);
    }
}

.circle {
    &__outer {
        opacity: .9;
        width: 1.5rem;
        height: 1.5rem;
        border-radius: 50%;
        color: $col-blue-600;
        border: 1px solid currentColor;
        padding: 1px;
    }

    &__inner {
        width: 100%;
        height: 100%;
        border-radius: 50%;
        border: 1px solid currentColor;
        background-color: currentColor;
        display: flex;
        align-items: center;
        justify-content: center;

        p {
            color: $col-zinc-50;
            font-family: $font-action;
            font-size: .8rem;
            font-weight: 700;
            text-align: center;
            // transform: translate(.1px, 1px)
        }
    }

    &__occupied {
        color: $col-red-700 !important;
    }

    &__vacant {
        color: $col-green-700 !important;
    }

    &__allowed {
        color: $col-zinc-600 !important;

        &>* {
            // border: none;
            border-color: $col-zinc-500 !important;
            background-color: $col-zinc-50 !important;

            p {
                color: $col-zinc-700;
            }
        }
    }

    &__blocked {
        color: $col-zinc-800 !important;

        &>* {
            // border: none;
            border-color: $col-zinc-800 !important;
        }
    }

    &__hidden {
        color: $col-zinc-600 !important;

        &>* {
            // border: none;
            border-color: $col-zinc-600 !important;
            background-color: $col-yellow-300 !important;

            p {
                color: $col-zinc-700;
            }
        }
    }

    &__active {
        cursor: pointer;
        opacity: 1;
    }
}
</style>
