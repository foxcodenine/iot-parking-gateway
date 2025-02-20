<template>
    <div>
        <form>
            <AssetDateSelector></AssetDateSelector>
            <button class="bbtn bbtn--sky mt-12" @click.prevent="fetchActivityLogs">Retrieve Activity Logs</button>
        </form>

        <div class="ttable__container">
            <table class="ttable  mt-8" @click="clearMessage">

                <thead>
                    <tr>
                        <th>
                            <span class="cursor-pointer" @click="sortTable('id')">
                                #
                            </span>
                            <svg class="t-sort-arrow" :class="sortArrow('id')">
                                <use xlink:href="@/assets/svg/sprite.svg#triangle-1"></use>
                            </svg>
                        </th>

                        <th>
                            <span class="cursor-pointer" @click="sortTable('device_id')">
                                Device ID
                            </span>
                            <svg class="t-sort-arrow" :class="sortArrow('device_id')">
                                <use xlink:href="@/assets/svg/sprite.svg#triangle-1"></use>
                            </svg>
                        </th>

                        <th>
                            <span class="cursor-pointer" @click="sortTable('is_occupied')">
                                Is Occupied
                            </span>
                            <svg class="t-sort-arrow" :class="sortArrow('is_occupied')">
                                <use xlink:href="@/assets/svg/sprite.svg#triangle-1"></use>
                            </svg>
                        </th>

                        <th>
                            <span class="cursor-pointer" @click="sortTable('happened_at')">
                                Happened At
                            </span>
                            <svg class="t-sort-arrow" :class="sortArrow('happened_at')">
                                <use xlink:href="@/assets/svg/sprite.svg#triangle-1"></use>
                            </svg>
                        </th>

                        <th>
                            <span class="cursor-pointer" @click="sortTable('firmware_version')">
                                Firmware Version
                            </span>
                            <svg class="t-sort-arrow" :class="sortArrow('firmware_version')">
                                <use xlink:href="@/assets/svg/sprite.svg#triangle-1"></use>
                            </svg>
                        </th>

                        <th>
                            <span class="cursor-pointer" @click="sortTable('magnet_abs_total')">
                                Magnet ABS Total
                            </span>
                            <svg class="t-sort-arrow" :class="sortArrow('magnet_abs_total')">
                                <use xlink:href="@/assets/svg/sprite.svg#triangle-1"></use>
                            </svg>
                        </th>

                        <th>
                            <span class="cursor-pointer" @click="sortTable('peak_distance_cm')">
                                Peak Distance (cm)
                            </span>
                            <svg class="t-sort-arrow" :class="sortArrow('peak_distance_cm')">
                                <use xlink:href="@/assets/svg/sprite.svg#triangle-1"></use>
                            </svg>
                        </th>

                        <th>
                            <span class="cursor-pointer" @click="sortTable('radar_cumulative')">
                                Radar Cumulative
                            </span>
                            <svg class="t-sort-arrow" :class="sortArrow('radar_cumulative')">
                                <use xlink:href="@/assets/svg/sprite.svg#triangle-1"></use>
                            </svg>
                        </th>

                        <th>
                            <span class="cursor-pointer" @click="sortTable('beacons_amount')">
                                Beacons Qty
                            </span>
                            <svg class="t-sort-arrow" :class="sortArrow('beacons_amount')">
                                <use xlink:href="@/assets/svg/sprite.svg#triangle-1"></use>
                            </svg>
                        </th>

                        <th>
                            <span class="cursor-pointer" >
                                Beacons
                            </span>
                            <svg class="t-sort-arrow" :class="sortArrow('beacons')">
                                <use xlink:href="@/assets/svg/sprite.svg#triangle-1"></use>
                            </svg>
                        </th>
                    </tr>
                </thead>


                <tbody>
                    <tr v-for="dataLog in returnDataLogs" :key="dataLog.id">
                        <td>{{ dataLog.id }}</td> 
                        <td>{{ dataLog.device_id }}</td> 
                        <td>{{ dataLog.is_occupied ? 'Occupied' : 'Vacant' }}</td>                      
                        <td>{{ new Date(dataLog.happened_at).toLocaleString() }}</td>                       
                        <td>{{ dataLog.firmware_version }}</td> 
                        <td>{{ dataLog.magnet_abs_total }}</td> 
                        <td>{{ dataLog.peak_distance_cm }} cm</td> 
                        <td>{{ dataLog.radar_cumulative }}</td> 
                        <td>{{ dataLog.beacons_amount }}</td> 
                        <td>{{ dataLog.beacons }}</td> 
                    </tr>
                </tbody>


            </table>

        </div>

    </div>
</template>

<!-- --------------------------------------------------------------- -->

<script setup>
import { useDebugStore } from '@/stores/debugStore';
import { useMessageStore } from '@/stores/messageStore';

import AssetDateSelector from '../AssetDateSelector.vue';
import { computed, ref } from 'vue';

// ---------------------------------------------------------------------

const debugStore = useDebugStore();
const messageStore = useMessageStore();

// - Data --------------------------------------------------------------

const sortBy = ref('id');
const sortDesc = ref('true');
const searchTerm = ref("");

const dataLogs = ref([]);

// - Computed ----------------------------------------------------------

const returnDataLogs = computed(() => {

    let list =  [...dataLogs.value];

    list.sort((a, b) => {
        let modifier = sortDesc.value ? -1 : 1;

        if (a[sortBy.value] < b[sortBy.value]) return -1 * modifier;
        if (a[sortBy.value] > b[sortBy.value]) return 1 * modifier;

        return 0;

    })
    return list;
})

// -- Methods ----------------------------------------------------------

async function fetchActivityLogs() {
    // Ensure a device is selected before proceeding
    const selectedDeviceID = debugStore.getSelectedDeviceID;
    if (!selectedDeviceID || selectedDeviceID.length < 1) {
        messageStore.setFlashMessages(["Please select a device to fetch activity logs."], "flash-message--orange");
        return;
    }

    try {
        const response = await debugStore.fetchActivityLogs();

        if (response?.status == 200) {
            const msg = response.data?.message ?? "Activity logs retrieved successfully.";
            messageStore.setFlashMessages([msg], "flash-message--green");

            dataLogs.value = response.data.activity_logs;
        }
    } catch (error) {
        console.error("! ParkingHistory.fetchActivityLogs !\n", error);
        const errMsg = error.response?.data ?? "Failed to fetch Activity Logs";
        messageStore.setFlashMessages([errMsg], "flash-message--red");
    }
}

function clearMessage() {
    messageStore.clearFlashMessage();
}

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

// ---------------------------------------------------------------------

</script>

<!-- --------------------------------------------------------------- -->

<style lang="scss" scoped>
// </style>