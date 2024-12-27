<template>
    <div class="parking-panel" ref="refParkinPanel">
        <div class="parking-panel__filter" @click="filtersStatusOn = !filtersStatusOn">
            <p>Filter by Status</p>
            <svg class="parking-panel__down-arrow" :class="{'parking-panel__down-arrow--active': filtersStatusOn}">
                <use xlink:href="@/assets/svg/sprite.svg#icon-interface-2" ></use>    
            </svg>
        </div>
        <div class="filters" v-show="filtersStatusOn">
            <div class="filters__item" :class="{ 'filters__item--selected': filter == 'all' }" @click="filter = 'all'">
                <svg class="filters__svg" style="padding: .4rem;">
                    <use xlink:href="@/assets/svg/sprite.svg#icon-radio-2" v-if="filter == 'all'"></use>
                    <use xlink:href="@/assets/svg/sprite.svg#icon-radio-1" v-else></use>
                </svg>
                <span class="filters__text">Devices on Map</span>
                <span class="filters__fraction">{{ qty.map }}/{{ qty.all }}</span>
            </div>
            <div class="filters__item" :class="{ 'filters__item--selected': filter == 'vacant' }"
                @click="filter = 'vacant'">
                <svg class="filters__svg" style="padding: .4rem;">
                    <use xlink:href="@/assets/svg/sprite.svg#icon-radio-2" v-if="filter == 'vacant'"></use>
                    <use xlink:href="@/assets/svg/sprite.svg#icon-radio-1" v-else></use>
                </svg>
                <span class="filters__text">Vacant</span>
                <span class="filters__fraction">{{ qty.vacant }}/{{ qty.all }}</span>
            </div>
            <div class="filters__item" :class="{ 'filters__item--selected': filter == 'occupied' }"
                @click="filter = 'occupied'">
                <svg class="filters__svg" style="padding: .4rem;">
                    <use xlink:href="@/assets/svg/sprite.svg#icon-radio-2" v-if="filter == 'occupied'"></use>
                    <use xlink:href="@/assets/svg/sprite.svg#icon-radio-1" v-else></use>
                </svg>
                <span class="filters__text">Occupied</span>
                <span class="filters__fraction">{{ qty.occupied }}/{{ qty.all }}</span>
            </div>
            <div class="filters__item" :class="{ 'filters__item--selected': filter == 'unknown' }"
                @click="filter = 'unknown'">
                <svg class="filters__svg" style="padding: .4rem;">
                    <use xlink:href="@/assets/svg/sprite.svg#icon-radio-2" v-if="filter == 'unknown'"></use>
                    <use xlink:href="@/assets/svg/sprite.svg#icon-radio-1" v-else></use>
                </svg>
                <span class="filters__text">Unknown</span>
                <span class="filters__fraction">{{ qty.unknown }}/{{ qty.all }}</span>
            </div>
        </div>

        <div class="parking-panel__filter" @click="filtersDateOn = !filtersDateOn">
            <p>Filter by Date</p>
            <svg class="parking-panel__down-arrow" :class="{'parking-panel__down-arrow--active': filtersDateOn}">
                <use xlink:href="@/assets/svg/sprite.svg#icon-interface-2" ></use>    
            </svg>
        </div>

        <div class="fdate" v-if="filtersDateOn">
            <div class="fdate__list">
                <div class="fdate__item" @click="filterDateBy = 'no_date_filter'">
                    <svg class="fdate__svg" style="padding: .4rem;">
                        <use xlink:href="@/assets/svg/sprite.svg#icon-radio-2" v-if="filterDateBy == 'no_date_filter'">
                        </use>
                        <use xlink:href="@/assets/svg/sprite.svg#icon-radio-1" v-else></use>
                    </svg>
                    <span class="fdate__text">No Date Filter (Show All Devices)</span>
                </div>
                <div class="fdate__item" @click="filterDateBy = 'not_communicating'">
                    <svg class="fdate__svg" style="padding: .4rem;">
                        <use xlink:href="@/assets/svg/sprite.svg#icon-radio-2" v-if="filterDateBy == 'not_communicating'">
                        </use>
                        <use xlink:href="@/assets/svg/sprite.svg#icon-radio-1" v-else></use>
                    </svg>
                    <span class="fdate__text">Devices Not Communication Since</span>
                </div>
                <div class="fdate__item" @click="filterDateBy = 'status_unchanged'">
                    <svg class="fdate__svg" style="padding: .4rem;">
                        <use xlink:href="@/assets/svg/sprite.svg#icon-radio-2" v-if="filterDateBy == 'status_unchanged'">
                        </use>
                        <use xlink:href="@/assets/svg/sprite.svg#icon-radio-1" v-else></use>
                    </svg>
                    <span class="fdate__text">Devices Status Unchanged Since</span>
                </div>
            </div>
            <DatePicker v-model="date" :attributes='attrs' expanded borderless transparent is-dark="true" />
        </div>



        <div class="parking-panel__search">
            <label for="search-input">Search Device</label>
            <input id="search-input" type="text" v-model="searchTerm" placeholder="Search By Name, ID or Group">
        </div>

        <div class="parking-panel__list">

            <a class="parking-device__item parking-device__item--favorite" v-for="device in getDevicesList.f"
                @click="setActiveWindow(device.device_id)">
                <svg :class="markersClass(device)" width="25px" height="25px" viewBox="0 0 100 100"
                    xmlns="http://www.w3.org/2000/svg">
                    <circle cx="50" cy="50" r="43" fill="#ffffff" />
                    <circle cx="50" cy="50" r="40" fill="currentColor" />
                    <text x="53" y="70" font-family="Arial, sans-serif" font-size="60" fill="#ffffff"
                        font-weight="normal" text-anchor="middle">P</text>
                </svg>
                <div class="parking-device__content">
                    <p class="parking-device__name">{{ device.name }}</p>
                    <p class="parking-device__id">{{ device.device_id }}</p>
                    <p class="parking-device__last-event">{{ timeElapsed(device.happened_at) }}</p>
                </div>
            </a>

            <div class="parking-device__line"></div>

            <div class="parking-device__item" v-for="device in getDevicesList.d"
                @click="setActiveWindow(device.device_id)">
                <svg :class="markersClass(device)" width="25px" height="25px" viewBox="0 0 100 100"
                    xmlns="http://www.w3.org/2000/svg">
                    <circle cx="50" cy="50" r="43" fill="#ffffff" />
                    <circle cx="50" cy="50" r="40" fill="currentColor" />
                    <text x="53" y="70" font-family="Arial, sans-serif" font-size="60" fill="#ffffff"
                        font-weight="normal" text-anchor="middle">P</text>
                </svg>
                <div class="parking-device__content">
                    <p class="parking-device__name">{{ device.name }}</p>
                    <p class="parking-device__id">{{ device.device_id }}</p>
                    <p class="parking-device__last-event">{{ timeElapsed(device.happened_at) }}</p>
                </div>
            </div>

        </div>


    </div>
</template>
<!-- --------------------------------------------------------------- -->
<script setup>
import { useDeviceStore } from '@/stores/deviceStore';
import { useMapStore } from '@/stores/mapStore';
import { formatToLocalDateTime, timeElapsed } from '@/utils/utils';
import { storeToRefs } from 'pinia';
import { computed, onMounted, reactive, ref, watch } from 'vue';
import { Calendar, DatePicker } from 'v-calendar';
import 'v-calendar/style.css';
import { useAppStore } from '@/stores/appStore';

// - Store -------------------------------------------------------------

const deviceStore = useDeviceStore();

const mapStore = useMapStore();

const appStore = useAppStore();

// - Data --------------------------------------------------------------
const refParkinPanel = ref(null);
const searchTerm = ref("");
const filter = ref('all');

const filtersStatusOn = ref(false);
const filtersDateOn = ref(false);
const filterDateBy = ref('no_date_filter')


// Set the initial date to tomorrow
const date = ref(getDateOffset(-7));
const dateUTC = ref(convertToUTC(Date.now()));

const attrs = ref([
    {
        key: 'today',
        highlight: {
            //   color: 'purple',
            fillMode: 'outline',
            contentClass: 'italic',
        },

        dates: new Date(),
    },
]);

const qty = reactive({
    all: 0,
    map: 0,
    vacant: 0,
    occupied: 0,
    unknown: 0,
});


// - Computed ----------------------------------------------------------

const getDevicesList = computed(() => {

    const initial_parking_check_date = appStore.getAppSettings['initial_parking_check_date'];

    
    if (!deviceStore.getDevicesList || deviceStore.getDevicesList.length === 0) {
        return [];
    }

    let devices = Object.values(deviceStore.getDevicesList);
    qty.all = devices.length;

    const vacant = devices.filter(item => {return item.is_occupied == false && item.happened_at > initial_parking_check_date});
    qty.vacant = vacant.length;

    const occupied = devices.filter(item => {return item.is_occupied == true && item.happened_at > initial_parking_check_date});
    qty.occupied = occupied.length;

    const unknown = devices.filter(item => {
        return item.happened_at <= initial_parking_check_date
    });
    qty.unknown = unknown.length;

    switch (filter.value) {
        case 'vacant':
            devices = vacant;
            break;
        case 'occupied':
            devices = occupied;
            break;
        case 'unknown':
            devices = unknown;
            break;

        default:
            break;
    }

    devices = devices.filter(item => {

        let date = "0";
        if (filterDateBy.value == 'not_communicating') {
            date = [item.happened_at, item.keepalive_at].reduce((a, b) => a > b ? a : b);
        }
        if (filterDateBy.value == 'status_unchanged') {
            date = item.happened_at
        }

        return date <= dateUTC.value;
    })

    devices = devices.filter(item => {
        return (
            item.device_id?.toLowerCase().includes(searchTerm.value.toLowerCase().trim()) ||
            item.name?.toLowerCase().includes(searchTerm.value.toLowerCase().trim()) ||
            item.network_type?.toLowerCase().includes(searchTerm.value.toLowerCase().trim()) ||
            String(item.firmware_version)?.toLowerCase().includes(searchTerm.value.toLowerCase().trim())
        )
    })

    // Sort by `happened_at` date
    devices.sort((a, b) => {
        const dateA = new Date(a.happened_at);
        const dateB = new Date(b.happened_at);
        return dateB - dateA; // For ascending order; use `dateB - dateA` for descending
    });

    qty.map = devices.length;

    const a = [...devices]
    const f = [];
    const d = [];

    devices.forEach(device => {
        if (appStore.getUserFavorites.includes(device.device_id)) {
            f.push(device);
        } else {
            d.push(device);
        }
    })

    return {a,f,d};
});

// - Watchers ----------------------------------------------------------

watch(date, (val) => {
    dateUTC.value = convertToUTC(val);
});



watch(() => getDevicesList.value.a, (val) => {
    deviceStore.setFilteredDevices(val);
}, {
    immediate: true,
})

// - Methods -----------------------------------------------------------

function setActiveWindow(id) {
    mapStore.setActiveWindow(id)
}

function markersClass(device) {

    if (formatToLocalDateTime(device.happened_at) == null) {
        return 'unknown'
    }

    if (device.is_occupied) {
        return 'occupied'
    } else {
        return 'vacant'
    }
};

function convertToUTC(dateString) {
    // Create a Date object from the input date string
    const date = new Date(dateString);

    // Convert the date to a UTC string in ISO 8601 format
    return date.toISOString();
}

function getDateOffset(d) {
    const today = new Date();
    const tomorrow = new Date(today);
    tomorrow.setDate(today.getDate() + d);
    return tomorrow;
}


</script>
<!-- --------------------------------------------------------------- -->
<style lang="scss" scoped>
.fdate {
    &__item {
        color: $col-zinc-50;
        display: flex;
        align-items: center;
        cursor: pointer;
    }



    &__svg {
        width: 1.6rem;
        height: 1.6rem;
    }

    &__text {
        font-family: $font-display;
        font-size: .8rem;
        text-transform: capitalize;
    }
}

.parking-panel {
    display: flex;
    flex-direction: column;
    align-items: stretch;
    gap: 0.5rem;
    padding: 0.5rem;
    width: 100%;
    height: 100%;
    /* Ensures the panel takes up the full height of its parent */
    font-size: 1rem;
    font-family: $font-action;
    text-transform: uppercase;
    color: $col-zinc-200;

    &__filter {
        display: flex;
        align-items: center;
        justify-content: space-between;
        padding: 0.1rem 0.5rem;
        border: 1px solid $col-zinc-700;
        background-color: rgba($col-slate-600, 0.1);
        cursor: pointer;

        &:hover {
            color: $col-zinc-50;
            border: 1px solid $col-zinc-500;
        }
    }

    &__search {
        display: flex;
        flex-direction: column;
        padding: 0.1rem 0.5rem;
        border: 1px solid $col-zinc-700;
        background-color: rgba($col-slate-600, 0.1);

        &:hover {
            color: $col-zinc-50;
            border: 1px solid $col-zinc-500;
        }

        label {
            cursor: pointer;
        }

        input {
            width: 100%;
            height: 1.5rem;
            background-color: transparent;
            color: $col-zinc-50;
            font-size: 1rem;
            font-family: $font-primary;
            border: none;
        }
    }

    &__list {
        flex: 1;
        /* Take up the remaining space in the container */
        display: flex;
        flex-direction: column;
        gap: 0.5rem;
        overflow-y: scroll;
        margin-top: 0.5rem;
        /* Adds spacing below the search bar */
        @extend %custom-scrollbar;

    }

    &__down-arrow {
        width: 0.35rem;
        height: 0.35rem;
        fill: $col-zinc-300;
        transform: rotate(-90deg);

        &--active {
            transform: rotate(180deg);
        }
    }
}

.parking-device {
    &__item {
        display: flex;
        align-items: center;
        gap: 0.5rem;
        min-height: 4rem;
        padding-bottom: 0.5rem;
        cursor: pointer;
        overflow: hidden;


        &:not(&:last-child) {
            border-bottom: 1px solid $col-zinc-600;
        }

        &--favorite {
            color: $col-yellow-200;
            .parking-device__last-event {
                color: rgba($col-yellow-200, .55) !important;
            }

        }
    }



    &__content {
        display: flex;
        flex-direction: column;
        gap: 0.2rem;
        font-size: 1rem;
        line-height: 1.1rem;
        text-transform: capitalize;
    }

    &__name {
        font-size: 1rem;
    }

    &__id {
        font-size: 0.8rem;
    }

    &__last-event {
        font-family: $font-display;
        color: $col-zinc-400;
    }
}

.filters {
    &__item {
        display: flex;
        align-items: center;
        text-transform: capitalize;
        font-family: $font-display;
        color: $col-zinc-300;
        border: 1px solid $col-zinc-700;
        margin-bottom: .5rem;
        padding: .2rem .4rem;
        border-radius: 2rem;
        cursor: pointer;
        transition: all .1s ease;

        &:hover,
        &--selected {
            border-color: $col-zinc-200;
        }
    }

    &__svg {
        width: 1.6rem;
        height: 1.6rem;
        margin-left: -5px;
    }

    &__fraction {
        flex: 1;
        text-align: end;
    }
}

.vacant {
    color: $col-green-600;
}

.occupied {
    color: $col-red-500;
}

.unknown {
    color: $col-indigo-500;
}
</style>

<style>
.fdate {

    .vc-title,
    .vc-arrow,
    .vc-nav-item,
    .vc-nav-title,
    .vc-nav-arrow {
        background-color: transparent !important;
    }

    .vc-nav-item {
        color: #d1d5db;
    }

    .vc-popover-content.direction-bottom {
        background-color: #18181b;
        border-color: #52525b;
    }
}
</style>