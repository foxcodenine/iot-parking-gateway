<template>

    <div class="ttable__container" v-if="getKeepaliveLogs.length > 0">
        <table class="ttable  mt-8" @click="clearMessage">
            <thead>
                <tr>
                    <th v-for="colname in Object.keys(getKeepaliveLogs[0])" @click="sortTable(colname)">
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
</template>
<!-- --------------------------------------------------------------- -->
<script setup>
import { useDebugStore } from '@/stores/debugStore';
import { storeToRefs } from 'pinia';
import { computed, ref } from 'vue';
import { formatToLocalDateTime } from '@/utils/dateTimeUtils';


const debugStore = useDebugStore();
const { getKeepaliveLogs } = storeToRefs(debugStore)

const props = defineProps({
    paginatedItems: {
        type: Array,
        default: [],
    }
});

const sortBy = ref('id');
const sortDesc = ref(true); // use boolean

// Compute sorted and formatted logs for the table
const returnKeepaliveLogs = computed(() => {
    let list = [...props.paginatedItems];
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
// 
</style>