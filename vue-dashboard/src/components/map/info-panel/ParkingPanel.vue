<template>
    <div class="parking-panel">
        <div class="parking-panel__filter">
            <p>Filters</p>
        </div>

        <div class="parking-panel__search">
            <label for="search-input">Search Device</label>
            <input id="search-input" type="text" placeholder="Search By Name, ID or Group">
        </div>

        <div class="parking-panel__list">

            <div class="parking-device__item" v-for="device in getDevicesList">
                <svg :class="markersClass(device)" width="25px" height="25px" viewBox="0 0 100 100" xmlns="http://www.w3.org/2000/svg">
                    <circle cx="50" cy="50" r="43" fill="#ffffff" />
                    <circle cx="50" cy="50" r="40" fill="currentColor" />
                    <text x="53" y="70" font-family="Arial, sans-serif" font-size="60" fill="#ffffff"
                        font-weight="normal" text-anchor="middle">P</text>
                </svg>
                <div class="parking-device__content">
                    <p class="parking-device__name">{{ device.name }}</p>                    
                    <p class="parking-device__last-event">{{ timeElapsed(device.happened_at) }}</p>
                </div>
            </div>

        </div>


    </div>
</template>
<!-- --------------------------------------------------------------- -->
<script setup>
import { useDeviceStore } from '@/stores/deviceStore';
import { formatToLocalDateTime, timeElapsed } from '@/utils/utils';
import { storeToRefs } from 'pinia';
import { computed } from 'vue';


const deviceStore = useDeviceStore();


const getDevicesList = computed(()=>{
    return deviceStore.getDevicesList;
});

function markersClass(device) {

    if (formatToLocalDateTime(device.happened_at) == null) {
        return 'inactive'
    }

    if (device.is_occupied) {
        return 'occupied'
    } else {
        return 'vacant'
    }
};


</script>
<!-- --------------------------------------------------------------- -->
<style lang="scss" scoped>
.parking-panel {
    display: flex;
    flex-direction: column;
    align-items: stretch;
    gap: 0.5rem;
    padding: 0.5rem;
    width: 100%;
    height: 100%; /* Ensures the panel takes up the full height of its parent */
    font-size: 1rem;
    font-family: $font-action;
    text-transform: uppercase;
    color: $col-zinc-200;

    &__filter {
        display: flex;
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
        flex: 1; /* Take up the remaining space in the container */
        display: flex;
        flex-direction: column;
        gap: 0.5rem;
        overflow-y: scroll;
        margin-top: 0.5rem; /* Adds spacing below the search bar */
        @extend %custom-scrollbar;
      
    }
}

.parking-device {
    &__item {
        display: flex;
        align-items: center;
        gap: 0.5rem;
        min-height: 3rem;        
        padding-bottom: 0.5rem;
        cursor: pointer;
        overflow: hidden;
       

        &:not(&:last-child) {
            border-bottom: 1px solid $col-zinc-600;
        }
    }

    &__content {
        display: flex;
        flex-direction: column;
        gap: 0.2rem;
        font-size: 1rem;
        line-height: 1.2rem;
        text-transform: capitalize;
    }

    &__name {
        font-size: 1rem;
    }

    &__last-event {
        font-family: $font-display;
        color: $col-zinc-400;
    }
}

.vacant {
    color: $col-green-600;
}

.occupied {
    color: $col-red-500;
}

.inactive {
    color: $col-indigo-500;
}


</style>