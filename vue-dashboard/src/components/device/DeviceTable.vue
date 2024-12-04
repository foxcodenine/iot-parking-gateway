<template>

    <input class="ttable__search mt-8" v-model="searchTerm" type="text" placeholder="Search...">

    <div class="ttable__container">

        <table class="ttable  mt-8">
            <thead>
                <tr>

                    <th class="cursor-pointer">
                        Device ID
                        <svg class="t-sort-arrow">
                            <use xlink:href="@/assets/svg/sprite.svg#triangle-1"></use>
                        </svg>
                    </th>
                    <th class="cursor-pointer">
                        Name
                        <svg class="t-sort-arrow">
                            <use xlink:href="@/assets/svg/sprite.svg#triangle-1"></use>
                        </svg>
                    </th>
                    <th class="cursor-pointer">
                        Network
                        <svg class="t-sort-arrow">
                            <use xlink:href="@/assets/svg/sprite.svg#triangle-1"></use>
                        </svg>
                    </th>
                    <th class="cursor-pointer">
                        Firmware
                        <svg class="t-sort-arrow">
                            <use xlink:href="@/assets/svg/sprite.svg#triangle-1"></use>
                        </svg>
                    </th>
                    <th class="cursor-pointer ">
                        Latitude
                        <svg class="t-sort-arrow">
                            <use xlink:href="@/assets/svg/sprite.svg#triangle-1"></use>
                        </svg>
                    </th>
                    <th class="cursor-pointer">
                        Longitude
                        <svg class="t-sort-arrow">
                            <use xlink:href="@/assets/svg/sprite.svg#triangle-1"></use>
                        </svg>
                    </th>
                    <th class="cursor-pointer">
                        Is Occupied
                        <svg class="t-sort-arrow">
                            <use xlink:href="@/assets/svg/sprite.svg#triangle-1"></use>
                        </svg>
                    </th>
                    <th class="cursor-pointer">
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
                            <cite class="t-btns__cite">Update Device</cite>
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
                    <td>{{ device.name }}</td>
                    <td>{{ device.network_type }}</td>
                    <td>{{ device.firmware_version }}</td>
                    <td>{{ device.latitude }}</td>
                    <td>{{ device.longitude }}</td>
                    <td>{{ device.is_occupied }}</td>
                    <td>{{ device.is_allowed }}</td>
                    <td>{{ device.is_blocked }}</td>
                    <td>{{ device.is_hidden }}</td>
                    <td v-if="getUserAccessLevel <= 2">
                        <div class="t-btns ml-auto" v-if="true">
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
                        <div v-else></div>     
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

import { onMounted, ref, reactive } from 'vue';

// - Store -------------------------------------------------------------

const messageStore = useMessageStore();

const deviceStore = useDeviceStore();
const { getDevicesList } = storeToRefs(deviceStore);

const authStore = useAuthStore();
const { getUserAccessLevel } = storeToRefs(authStore)

// -- Data -------------------------------------------------------------
const searchTerm = ref("");
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

function initEditDevice() {
    
}
function initDeleteDevice() {

}



</script>

<!-- --------------------------------------------------------------- -->

<style lang="scss" scoped>
// Import global SCSS styles for consistent styling across the application</style>
