<template>
    <InfoWindow v-if="activeWindow == device.device_id"
        :options="{ position: { lat: device.latitude, lng: device.longitude }, anchorPoint: 'CENTER' }">
        <div class="info-window">
            <div class="info-window__close" @click="closeWindow">&times;</div>

            <div class="info-window__header">
                <h3>{{ device.name }}</h3>
            </div>

            <div class="info-window__body mt-2 ">
                <div class="info-window__content">
                    <span>Device ID</span>
                    <p>{{ device.device_id }}</p>
                </div>
                <div class="info-window__content">
                    <span>Network</span>
                    <p>{{ device.network_type }}</p>
                </div>
                <div class="info-window__content">
                    <span>Firmware</span>
                    <p>{{ device.firmware_version }}</p>
                </div>
                <div class="info-window__content" v-if="formatToLocalDateTime(device.happened_at)">
                    <span>Status</span>
                    <p class="vacant" :class="{ 'occupied': device.is_occupied }">{{ device.is_occupied ? 'Occupied' :
                        'Vacant' }}</p>
                </div>
                <div class="info-window__content" v-else>
                    <span>Status</span>
                    <p class="unknown">Unknown</p>
                </div>
                <div class="info-window__content">
                    <span>Last Event</span>
                    <p v-html="timeSinceParked ?? 'Device has not reported <br />any parking activity yet.'"></p>
                </div>
                <div class="info-window__content mb-1" v-if="formatToLocalDateTime(device.happened_at)">
                    <span></span>
                    <p>{{ formatToLocalDateTime(device.happened_at) }}</p>
                </div>
            </div>

            <div class="info-window__footer mt-2 mb-1">
                <svg class="info-window__svg" @click="toggleDeviceInFavorites(device.device_id)">
                    <use xlink:href="@/assets/svg/sprite.svg#icon-star-2" v-if="isFavorite"></use>
                    <use xlink:href="@/assets/svg/sprite.svg#icon-star-1" v-else></use>
                </svg>
                <svg class="info-window__svg" @click="updatedMarkerLocation">
                    <use xlink:href="@/assets/svg/sprite.svg#icon-map-8"></use>
                </svg>
                <svg class="info-window__svg" @click="navigateToDevice()">
                    <use xlink:href="@/assets/svg/sprite.svg#icon-route-start"></use>
                </svg>
                <svg class="info-window__svg" style="padding: .2rem;" @click="editDevice(device.device_id)">
                    <use xlink:href="@/assets/svg/sprite.svg#icon-pencil-9"></use>
                </svg>
                <svg class="info-window__svg" style="padding: .4rem;">
                    <use xlink:href="@/assets/svg/sprite.svg#icon-bug"></use>
                </svg>
            </div>

        </div>
    </InfoWindow>
</template>

<!-- --------------------------------------------------------------- -->

<script setup>
import { useAppStore } from '@/stores/appStore';
import { formatToLocalDateTime, timeElapsed } from '@/utils/utils';
import { computed, ref, watch } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import { InfoWindow } from 'vue3-google-map';

// - Store -------------------------------------------------------------

const appStore = useAppStore();

// - Router ------------------------------------------------------------

const router = useRouter();

// - Emit -------------------------------------------------------------

const emit = defineEmits(['emitCloseWindow', 'emitUpdatedMarkerLocation'])

// - Props -------------------------------------------------------------
const props = defineProps({
    device: {
        type: Object,
        required: true,
    },
    activeWindow: {
        type: [Number, null, String],
        required: true,
    }
});

// - Data --------------------------------------------------------------

const timeSinceParked = ref('n/a');

// - Computed ----------------------------------------------------------

const isFavorite = computed(() => {
    return appStore.getUserFavorites.includes(props.device.device_id)
})



// - Methods -----------------------------------------------------------

function closeWindow() {
    emit('emitCloseWindow');
}

async function toggleDeviceInFavorites(deviceID) {
    try {
        appStore.toggleDeviceInFavorites(deviceID);
        const res = await appStore.updateUpdateFavorites();
    } catch (error) {
        console.error("! ParkingInfoWindow.toggleDeviceInFavorites !\n", error);
    }
}

function updatedMarkerLocation() {
    emit('emitUpdatedMarkerLocation', props.device.device_id)
}


function navigateToDevice() {

    const lat = props.device.latitude;
    const lng = props.device.longitude;
    const formattedCoordinates = `${lat},${lng}`;

    // Construct the Google Maps URL
    const googleMapsURL = `https://www.google.com/maps/place/${encodeURIComponent(formattedCoordinates)}/@${formattedCoordinates},17z/data=!3m1!4b1!4m4!3m3!8m2!3d${lat}!4d${lng}?entry=ttu`;

    window.open(googleMapsURL, '_blank');
}

function editDevice(deviceID) {
    closeWindow();
    router.push({
        name: 'deviceEditView',
        params: { deviceID },
    });
}

// - Hooks -------------------------------------------------------------

(() => {
    timeSinceParked.value = timeElapsed(props.device.happened_at);
})();

setInterval(() => {
    timeSinceParked.value = timeElapsed(props.device.happened_at);
}, 5000);




</script>

<!-- --------------------------------------------------------------- -->

<style lang="scss" scoped>

.info-window {

    border-radius: 3px !important;
    box-shadow: 0 4px 10px rgba(0, 0, 0, 0.1);
    overflow: hidden;
    min-width: 200px;

    &__close {
        cursor: pointer;
        position: absolute;
        font-size: 1.5rem;
        top: .0rem;
        right: .5rem;
        transition: all .1s ease;

        &:hover {
            color: $col-red-400;
        }
    }


    &__header {
        padding: 10px 10px;
        display: flex;
        justify-content: space-between;
        align-items: center;
        justify-content: center;
        border-bottom: 1px solid $col-zinc-700;

        h3 {
            margin: 0;
            font-size: 1rem;
            font-weight: bold;
        }

        .close-btn {
            background: none;
            border: none;
            color: $col-zinc-300;
            font-size: 1.4rem;
            font-weight: 100;
            cursor: pointer;

            &:hover {
                color: $col-red-400;
            }
        }
    }

    &__body {
        border-bottom: 1px solid $col-zinc-700;
    }

    &__content {
        padding: 3px 10px;
        font-size: .9rem;
        display: grid;
        grid-template-columns: 5.5rem 1fr;
        align-items: center;
        width: max-content !important;
        font-weight: 500;
        font-family: $font-action;

        span {
            // font-family: $font-display;;
            font-weight: 300;
        }
    }

    &__svg {
        fill: currentColor;
        width: 1.8rem;
        height: 1.8rem;
        padding: .3rem;
        border: .5px solid currentColor;
        border-radius: 2px;
        cursor: pointer;
        flex: 1;

        &:hover {
            color: $col-lime-400;
        }
    }

    &__footer {
        padding: 5px 10px;
        display: flex;
        flex-direction: row;
        justify-content: center;
        gap: .5rem;


    }
}

.vacant {
    color: $col-lime-500;
}

.occupied {
    color: $col-red-500;
}

.unknown {
    color: $col-indigo-500;
}
</style>


