<template>
    <div>
        <form>
            <AssetDateSelector></AssetDateSelector>
            <button class="bbtn bbtn--blue mt-12" @click.prevent="fetchActivityLogs">Retrieve Activity Logs</button>
        </form>

        <!-- show here -->

        <paginate v-model="currentPage" :page-count="pageCount" :click-handler="handlePageClick" :prev-text="'Prev'" :next-text="'Next'"
            :container-class="'pagination mt-10'" />

            <div class="ttable__container" v-if="dataLogs.length > 0">
            <table class="ttable  mt-8" @click="clearMessage">
                <thead>
                    <tr>
                        <th v-for="colname in Object.keys(dataLogs[0])" @click="sortTable(colname)">
                            <span class="cursor-pointer">
                                {{ prettifyString(colname) }}
                            </span>
                            &nbsp;
                            <svg class="t-sort-arrow" :class="sortArrow(colname)">
                                <use xlink:href="@/assets/svg/sprite.svg#triangle-1"></use>
                            </svg>
                        </th>
                    </tr>
                </thead>
                <tbody>
                    <tr v-for="keepaliveLog in returnDataLogs" :key="keepaliveLog.id">
                        <td v-for="recordValue in Object.values(keepaliveLog)">{{ recordValue }}</td>
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
import { formatToLocalDateTime } from '@/utils/dateTimeUtils';
import Paginate from 'vuejs-paginate-next';


// Pagination settings
const perPage = 20;
const currentPage = ref(1);

// ---------------------------------------------------------------------

const debugStore = useDebugStore();
const messageStore = useMessageStore();

// - Data --------------------------------------------------------------

const sortBy = ref('id');
const sortDesc = ref('true');
const searchTerm = ref("");

const dataLogs = ref([]);

// - Computed ----------------------------------------------------------



// Compute total pages using the length of the data logs array
const pageCount = computed(() => Math.ceil(dataLogs.value.length / perPage));

// Compute the items for the current page
const paginatedItems = computed(() => {
    const start = (currentPage.value - 1) * perPage;
    return dataLogs.value.slice(start, start + perPage);
});

const returnDataLogs = computed(() => {

    let list = [...paginatedItems.value];


    list.forEach(log => {
        log.happened_at = formatToLocalDateTime(log.happened_at);
        delete log.created_at;
        delete log.raw_id;
    });

    list.sort((a, b) => {
        let modifier = sortDesc.value ? -1 : 1;


        if (a[sortBy.value] < b[sortBy.value]) return -1 * modifier;
        if (a[sortBy.value] > b[sortBy.value]) return 1 * modifier;

        return 0;

    })

    return list;
});


// -- Methods ----------------------------------------------------------

// Handle page click event
function handlePageClick(page) {
    currentPage.value = page;
}

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

function prettifyString(str) {
  return str
    .split('_')
    .map(word => word.charAt(0).toUpperCase() + word.slice(1).toLowerCase())
    .join(' ');
}

// ---------------------------------------------------------------------

</script>

<!-- --------------------------------------------------------------- -->

<style lang="scss" scoped>
// @import "https://cdn.jsdelivr.net/npm/bootstrap@5.1.3/dist/css/bootstrap.min.css";

:deep(){.pagination {
    display: flex;
    justify-content: center;
    list-style: none;
    gap: 0.5rem;

    font-family: $font-action;

    .page-item {
        .page-link {
            display: block;
            padding: 0.3rem 0.5rem;
            min-width: 2rem;
            display: flex;
            justify-content: center;
            border: 1px solid $col-gray-300;
            background-color: $col-white;
            color: $col-gray-700;
            text-decoration: none;
            border-radius: 0.375rem;
            transition: background-color 0.15s ease-in-out, border-color 0.15s ease-in-out, color 0.15s ease-in-out;
            cursor: pointer;
        }

        &:hover:not(.active):not(.disabled) .page-link {
            background-color: $col-gray-100;
            border-color: $col-gray-400;
            color: $col-gray-800;
        }

        &.disabled .page-link {
            cursor: not-allowed;
            opacity: 0.6;
            pointer-events: none;
            background-color: $col-gray-50;
            border-color: $col-gray-200;
        }

        &.active .page-link {
            z-index: 3;
            background-color: $col-zinc-200;
            border-color: $col-zinc-700;
            //   color: $col-white;
            box-shadow: 0 0.125rem 0.25rem rgba(0, 0, 0, 0.1);
        }
    }
}}



</style>