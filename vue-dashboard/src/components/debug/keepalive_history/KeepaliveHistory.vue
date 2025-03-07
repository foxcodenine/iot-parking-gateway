<template>
    <div>
        <form>
            <AssetDateSelector></AssetDateSelector>
            <button class="bbtn bbtn--blue mt-12" @click.prevent="fetchKeepaliveLogs">Retrieve Keepalive Logs</button>
        </form>

        <div v-if="keepaliveLogs.length" class="report-title mt-8">
            <span>{{ getKeepaliveLogDetails.name }}</span>
            <span>{{ getKeepaliveLogDetails.id }}</span>
            <span>{{ getKeepaliveLogDetails.fromDate }}</span>
            <span>{{ getKeepaliveLogDetails.toDate }}</span>
        </div>

        <TheTabs class="mt-8" v-if="keepaliveLogs.length" :tabsObjectData="tabsObjectData_1" :isDisabled="false"
            :layoutBreakpoint="500" @setActiveTab="tabsObjectData_1.activeTab = $event"></TheTabs>

        <paginate v-if="tabsObjectData_1.activeTab == 'KeepaliveHistoryTable'" v-model="currentPage"
            :page-count="pageCount" :click-handler="handlePageClick" :prev-text="'Prev'" :next-text="'Next'"
            :container-class="'pagination mt-10'" />

        <KeepAlive>
            <component :paginatedItems="paginatedItems" :is="componentMap[tabsObjectData_1.activeTab]" class="mt-6">
            </component>
        </KeepAlive>


    </div>
</template>

<!-- --------------------------------------------------------------- -->
<script setup>

import { ref, computed, markRaw, reactive } from 'vue';
import { useDebugStore } from '@/stores/debugStore';
import { useMessageStore } from '@/stores/messageStore';
import AssetDateSelector from '../AssetDateSelector.vue';
import Paginate from 'vuejs-paginate-next';
import KeepaliveHistoryCharts from './KeepaliveHistoryCharts.vue';
import KeepaliveHistoryTable from './KeepaliveHistoryTable.vue';
import TheTabs from '@/components/commen/TheTabs.vue'
import { storeToRefs } from 'pinia';

// -- Store ------------------------------------------------------------

const debugStore = useDebugStore();
const messageStore = useMessageStore();

const { getKeepaliveLogDetails } = storeToRefs(debugStore);

// -- Data -------------------------------------------------------------

// Pagination and data-related state
const perPage = 10;
const currentPage = ref(1);
const keepaliveLogs = ref([]); // full list of logs


// Reactive state to manage tab data and the active tab identifier
const tabsObjectData_1 = reactive({
    activeTab: 'KeepaliveHistoryTable',
    tabs: {
        KeepaliveHistoryTable: 'Table',
        KeepaliveHistoryCharts: 'Charts',
    },
});

const componentMap = reactive({
    KeepaliveHistoryCharts: markRaw(KeepaliveHistoryCharts),
    KeepaliveHistoryTable: markRaw(KeepaliveHistoryTable),

});

// -- Computed ---------------------------------------------------------

// Compute total pages based on full data
const pageCount = computed(() => Math.ceil(keepaliveLogs.value.length / perPage));

// Compute items for the current page
const paginatedItems = computed(() => {
    const start = (currentPage.value - 1) * perPage;
    return keepaliveLogs.value.slice(start, start + perPage);
});

// -- Methods ----------------------------------------------------------

// Fetch keepalive logs from the API
async function fetchKeepaliveLogs() {
    const selectedDeviceID = debugStore.getSelectedDeviceID;
    if (!selectedDeviceID || selectedDeviceID.length < 1) {
        messageStore.setFlashMessages(["Please select a device to fetch keepalive logs."], "flash-message--orange");
        return;
    }
    try {
        const response = await debugStore.fetchKeepaliveLogs();
        if (response?.status === 200) {
            const msg = response.data?.message ?? "Keepalive logs retrieved successfully.";
            messageStore.setFlashMessages([msg], "flash-message--green");
            keepaliveLogs.value = response.data.keepalive_logs;
        }
    } catch (error) {
        console.error("! KeepaliveLog.fetchKeepaliveLogs !\n", error);
        const errMsg = error.response?.data ?? "Failed to fetch Keepalive Logs";
        messageStore.setFlashMessages([errMsg], "flash-message--red");
    }
}

// Handle pagination click event
function handlePageClick(page) {
    currentPage.value = page;
}

</script>


<!-- --------------------------------------------------------------- -->



<style lang="scss" scoped>
:deep() {
    .pagination {
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
    }
}

.report-title {
    border: 0.5px solid $col-red-400;
    border-radius: 3px;
    color: $col-red-600 !important;
    padding: .2rem.5rem;
    gap: 2rem;
    display: flex;
    justify-content: center;
    align-items: center;
    font-size: .9rem;  
    font-weight: 400;
    text-transform: uppercase;
    font-family: $font-display;

    @include respondDesktop(600) {
        flex-direction: column;
        gap: 0.5rem;
    }
}

</style>
