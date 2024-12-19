<template>
    <div class="the-map">
        <GoogleMap v-if="apiKey" 
        :api-key="apiKey" 
        map-id="9e8ca8994cbac798" 
        style="width: 100%; height: 100vh" :center="getMapCenter" :zoom="17" @zoom_changed="zoomChanged" ref="mapRef">



            <CustomMarker v-for="device in getDevicesList"  @zoom_changed="zoomChanged" :style="markerStyle()"
                @click="aaa(device)"
                :options="{ position: { lat: device.latitude, lng: device.longitude }, anchorPoint: 'CENTER' }">
                <svg :width="markerSize" :height="markerSize" viewBox="0 0 100 100" xmlns="http://www.w3.org/2000/svg">
                    <circle cx="50" cy="50" r="48" fill="#ffffff" stroke="#0000ff" stroke-width="4" />
                    <circle cx="50" cy="50" r="40" fill="#0000ff" />
                    <text x="53" y="70" font-family="Arial, sans-serif" font-size="60" fill="#ffffff" font-weight="bold" text-anchor="middle">P</text>
                </svg>
            </CustomMarker>

            <!-- <Marker v-for="device in getDevicesList"
            :options="{ position: { lat: device.latitude, lng: device.longitude } }" /> -->

        </GoogleMap>
    </div>
</template>

<!-- --------------------------------------------------------------- -->
<script setup>

import { useAppStore } from '@/stores/appStore';
import { useDeviceStore } from '@/stores/deviceStore';
import { storeToRefs } from 'pinia';
import { onMounted, ref } from 'vue';
import { computed, toRaw } from 'vue';
import { GoogleMap, Marker, CustomMarker } from 'vue3-google-map';

// - Store -------------------------------------------------------------

const appStore = useAppStore();
const deviceStore = useDeviceStore();

const { getDevicesList } = storeToRefs(deviceStore);

// - Data --------------------------------------------------------------

const apiKey = ref(null);
const mapZoom = ref(17)
const mapRef = ref(null);

// - Computed ----------------------------------------------------------

const getMapCenter = computed(() => {
    const lat = appStore.getAppSettings.default_latitude ? Number(appStore.getAppSettings.default_latitude) : 0;
    const lng = appStore.getAppSettings.default_longitude ? Number(appStore.getAppSettings.default_longitude) : 0;
    return { lat, lng };
});


function zoomChanged() {
    mapZoom.value = mapRef.value.map.getZoom()
}

function markerStyle() {
    return "opacity: 1;"
}

const markerSize = computed(() => {
    const zoom = Math.round(mapZoom.value);
    return Math.min(10 + (zoom - 10) * 3, 64); // Linearly scales marker size
});
function aaa(d){
    console.log(d)
}

// - Hooks -------------------------------------------------------------
(async () => {
    try {
        await deviceStore.fetchDevices();
    } catch (error) {
        console.error('! DeviceView deviceStore.fetchDevices() !\n', error);
    }
})()

onMounted(async () => {
    apiKey.value = await appStore.getGoogleApiKey();
});

// ---------------------------------------------------------------------

</script>

<!-- --------------------------------------------------------------- -->

<style lang="scss" scoped>
.the-map {
    background-color: $col-blue-300;
    // background-image: url("@/assets/images/map.jpg");
    // min-height: 100vh;
}
</style>