<template>
    <div>
        <form>
            <AssetDateSelector></AssetDateSelector>
            <button class="bbtn bbtn--blue mt-12" @click.prevent="fetchKeepaliveLogs">Retrieve Keepalive Logs</button>
        </form>


        <paginate v-model="currentPage" :page-count="pageCount" :click-handler="handlePageClick" :prev-text="'Prev'"
            :next-text="'Next'" :container-class="'pagination mt-10'" />

        <Bar v-if="keepaliveLogs.length > 0" id="my-chart-id" :options="chartOptions" :data="chartData" />

        <div class="ttable__container" v-if="keepaliveLogs.length > 0">
            <table class="ttable  mt-8" @click="clearMessage">
                <thead>
                    <tr>
                        <th v-for="colname in Object.keys(keepaliveLogs[0])" @click="sortTable(colname)">
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
                    <tr v-for="keepaliveLog in returnKeepaliveLogs" :key="keepaliveLog.id">
                        <td v-for="recordValue in Object.values(keepaliveLog)">{{ recordValue }}</td>
                    </tr>
                </tbody>
            </table>
        </div>

    </div>
</template>

<!-- --------------------------------------------------------------- -->
<script setup>
import { ref, computed } from 'vue';
import { useDebugStore } from '@/stores/debugStore';
import { useMessageStore } from '@/stores/messageStore';
import AssetDateSelector from '../AssetDateSelector.vue';
import Paginate from 'vuejs-paginate-next';
import { formatToLocalDateTime } from '@/utils/dateTimeUtils';
import { Bar } from 'vue-chartjs';
import { Chart as ChartJS, Title, Tooltip, Legend, BarElement, CategoryScale, LinearScale } from 'chart.js';

// Register Chart.js components
ChartJS.register(Title, Tooltip, Legend, BarElement, CategoryScale, LinearScale);

const debugStore = useDebugStore();
const messageStore = useMessageStore();

// Pagination and data-related state
const perPage = 10;
const currentPage = ref(1);
const keepaliveLogs = ref([]); // full list of logs
const sortBy = ref('id');
const sortDesc = ref(true); // use boolean

// Compute total pages based on full data
const pageCount = computed(() => Math.ceil(keepaliveLogs.value.length / perPage));

// Compute items for the current page
const paginatedItems = computed(() => {
    const start = (currentPage.value - 1) * perPage;
    return keepaliveLogs.value.slice(start, start + perPage);
});

// Compute sorted and formatted logs for the table
const returnKeepaliveLogs = computed(() => {
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
    });
    return list;
});

// Computed chart data using the full list of logs
const chartData = computed(() => {
    return {
        labels: keepaliveLogs.value.map(log => formatToLocalDateTime(log.happened_at)),
        datasets: [{
            label: 'Current', // adjust label as needed
            data: keepaliveLogs.value.map(log => log.current),
            backgroundColor: '#3b82f6'
        }]
    };
});

// Chart options
const chartOptions = {
  responsive: true,
  scales: {
    x: {
      ticks: {
        // Rotate labels by 45 degrees
        maxRotation: 90,
        minRotation: 0,

        // Limit the total number of ticks displayed
        maxTicksLimit: 10,

        // Or shorten them using a callback
        callback: function(value, index, ticks) {
          // 'this' references the axis; get the full label
          const label = this.getLabelForValue(value);
          // e.g., show only the time or a shorter date
          return label.slice(0, 11); // "Feb 29, 20"
        },
      },
    },
  },
};


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
            console.log(keepaliveLogs.value);
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

// Toggle sort order and set sorting column
function sortTable(col) {
    if (sortBy.value === col) {
        sortDesc.value = !sortDesc.value;
    }
    sortBy.value = col;
}

// Return arrow classes for sorting indicators
function sortArrow(col) {
    if (col === sortBy.value) {
        return { 't-sort-arrow--active': true, 't-sort-arrow--desc': sortDesc.value };
    }
    return {};
}

// Convert snake_case to Title Case
function prettifyString(str) {
    return str
        .split('_')
        .map(word => word.charAt(0).toUpperCase() + word.slice(1).toLowerCase())
        .join(' ');
}
</script>


<!-- --------------------------------------------------------------- -->

<style lang="scss" scoped>
::v-deep .pagination {
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
</style>