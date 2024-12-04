<template>

    <input class="ttable__search mt-8" v-model="searchTerm" type="text" placeholder="Search...">

    <div class="ttable__container">

        <table class="ttable  mt-8">
            <thead>
                <tr>

                    <th class="cursor-pointer w-48">
                        Device ID
                        <svg class="t-sort-arrow">
                            <use xlink:href="@/assets/svg/sprite.svg#triangle-1"></use>
                        </svg>
                    </th>
                    <th class="cursor-pointer w-52">
                        Name
                        <svg class="t-sort-arrow">
                            <use xlink:href="@/assets/svg/sprite.svg#triangle-1"></use>
                        </svg>
                    </th>
                    <th class="cursor-pointer w-24">
                        Network
                        <svg class="t-sort-arrow">
                            <use xlink:href="@/assets/svg/sprite.svg#triangle-1"></use>
                        </svg>
                    </th>
                    <th class="cursor-pointer w-16">
                        Firmware
                        <svg class="t-sort-arrow">
                            <use xlink:href="@/assets/svg/sprite.svg#triangle-1"></use>
                        </svg>
                    </th>
                    <th class="w-8"></th>
                    <th class="cursor-pointer w-24">
                        Latitude
                        <svg class="t-sort-arrow">
                            <use xlink:href="@/assets/svg/sprite.svg#triangle-1"></use>
                        </svg>
                    </th>
                    <th class="cursor-pointer w-24">
                        Longitude
                        <svg class="t-sort-arrow">
                            <use xlink:href="@/assets/svg/sprite.svg#triangle-1"></use>
                        </svg>
                    </th>
                    <th class="cursor-pointer" >
                        Is Occupied
                        <svg class="t-sort-arrow">
                            <use xlink:href="@/assets/svg/sprite.svg#triangle-1"></use>
                        </svg>
                    </th>
                    <th class="cursor-pointer" >
                        White List
                        <svg class="t-sort-arrow">
                            <use xlink:href="@/assets/svg/sprite.svg#triangle-1"></use>
                        </svg>
                    </th>
                    <th class="cursor-pointer">
                        Black List
                        <svg class="t-sort-arrow">
                            <use xlink:href="@/assets/svg/sprite.svg#triangle-1"></use>
                        </svg>
                    </th>
                    <th class="cursor-pointer">
                        Hide
                        <svg class="t-sort-arrow">
                            <use xlink:href="@/assets/svg/sprite.svg#triangle-1"></use>
                        </svg>
                    </th>
                    <th v-if="getUserAccessLevel <= 2">
                    </th>
                </tr>
            </thead>
            <tbody>



                <tr v-if="false">

                    <td>
                        <input class="ttable__input" type="text" value="Pam Beesly">
                    </td>
                    <td>860226067572735</td>
                    <td>NB-IoT</td>
                    <td>3.2</td>
                    <td>
                        <input class="ttable__input w-24" type="text" value="35.928676">
                    </td>
                    <td>
                        <input class="ttable__input w-24" type="text" value="14.418638">
                    </td>
                    <td>true</td>
                    <td>
                        <select>
                            <option value="active">false</option>
                            <option value="deactive">true</option>
                        </select>
                    </td>
                    <td>
                        <select>
                            <option value="active">false</option>
                            <option value="deactive">true</option>
                        </select>
                    </td>
                    <td>
                        <select>
                            <option value="active">false</option>
                            <option value="deactive">true</option>
                        </select>
                    </td>

                    <td>
                        <div class="t-btns">
                            <a class="t-btns__btn t-btns__btn--yes">
                                Yes
                            </a>
                            <a class="t-btns__btn t-btns__btn--no">
                                No
                            </a>
                        </div>
                    </td>
                </tr>



                <tr v-for="device in getDevicesList"
                    :class="{ 'bg-lime-200': selectedDevice.device_id == device.device_id }">

                    <td>{{ device.device_id }}</td>
                    
                    <td v-if="action=='edit'&&selectedDevice.device_id==device.device_id" >
                        <input class="ttable__input w-48" type="text" v-model="selectedDevice.name">
                    </td>
                    <td v-else>{{ device.name }}</td>
      
                    <td>{{ device.network_type }}</td>

                    <td v-if="action=='edit'&&selectedDevice.device_id==device.device_id" >
                        <input class="ttable__input w-12" type="text" v-model="selectedDevice.firmware_version">
                    </td>
                    <td v-else>{{ device.firmware_version }}</td>

                    <td class="w-8">                   
                        <svg v-if="action=='edit'&&selectedDevice.device_id==device.device_id" class="ttable__location-btn ">
                                <use xlink:href="@/assets/svg/sprite.svg#icon-google-maps-2"></use>
                        </svg>                     
                    </td>

                    <td v-if="action=='edit'&&selectedDevice.device_id==device.device_id" >
                        <input class="ttable__input w-20" type="text" v-model="selectedDevice.latitude">
                    </td>
                    <td v-else>{{ device.latitude }}</td>

                    <td v-if="action=='edit'&&selectedDevice.device_id==device.device_id" >
                        <input class="ttable__input w-20" type="text" v-model="selectedDevice.longitude">
                    </td>
                    <td v-else>{{ device.longitude }}</td>

                    <td @click="toggelOccupied(device)" class="ps-10">
                        <div v-if="action=='edit'&&selectedDevice.device_id==device.device_id" class="circle__outer circle__active" :class="{'circle__occupied': selectedDevice.is_occupied, 'circle__vacant': !selectedDevice.is_occupied}">
                            <div class="circle__inner"></div>
                        </div>
                        <div v-else class="circle__outer" :class="{'circle__occupied': device.is_occupied, 'circle__vacant': !device.is_occupied}">
                            <div class="circle__inner"></div>
                        </div>

                    </td>
                    <td @click="toggelAllowed(device)" class="ps-8">
                        <div v-if="action=='edit'&&selectedDevice.device_id==device.device_id" class="circle__outer circle__active" :class="{'circle__allowed': selectedDevice.is_allowed}">
                            <div class="circle__inner"></div>
                        </div>
                        <div v-else class="circle__outer" :class="{'circle__allowed': device.is_allowed}">
                            <div class="circle__inner"></div>
                        </div>
                    </td>
          
                    <td>{{ device.is_blocked }}</td>
                    <td>{{ device.is_hidden }}</td>
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
                            <a class="t-btns__btn t-btns__btn--yes" @click="editOrDelete" v-tooltip="{ content: action == 'edit' ? 'Edit Device' : 'Delete Device' }" >
                                Yes
                            </a>
                            <a class="t-btns__btn t-btns__btn--no" @click="cancel" v-tooltip="{ content: action == 'edit' ? 'Cancel Edit' : 'Cancel Delete' }">
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

import { onMounted, ref, reactive } from 'vue';

// - Store -------------------------------------------------------------

const messageStore = useMessageStore();

const deviceStore = useDeviceStore();
const { getDevicesList } = storeToRefs(deviceStore);

const authStore = useAuthStore();
const { getUserAccessLevel } = storeToRefs(authStore)

// -- Data -------------------------------------------------------------
const searchTerm = ref("");

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

function cancel() {
    action.value = null;
    selectedDevice.device_id = null;
    selectedDevice.name = null;
    selectedDevice.network_type = null;
    selectedDevice.firmware_version = null;
    selectedDevice.latitude = null;
    selectedDevice.longitude = null;
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



</script>

<!-- --------------------------------------------------------------- -->

<style lang="scss" scoped>
.t-btns {
    justify-content: end !important;
}
.ttable__input, .ttable__select {
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
        opacity: .6;
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
    }
    &__occupied {
        color: $col-red-700 !important;
    }
    &__vacant {
        color: $col-green-700 !important;
    }
    &__allowed {
        color: $col-zinc-600 !important;
        
        & > * {
            // border: none;
            border-color: $col-zinc-500 !important;
            background-color: $col-zinc-50 !important;
        }
    }
    &__active {
        cursor: pointer;
        opacity: 1;
    }
}
</style>
