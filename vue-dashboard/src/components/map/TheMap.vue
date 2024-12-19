<template>
    <div class="the-map">
        <GoogleMap v-if="apiKey" 
            :api-key="apiKey"  :map-id="getMapId" style="width: 100%; height: 100vh" 
            :center="getMapCenter"  :zoom="17" @zoom_changed="zoomChanged"  ref="mapRef">

                <ParkingMarker v-for="device in getDevicesList" :device="device" :mapZoom="mapZoom" @click="activeWindow=device.device_id"></ParkingMarker>
                <ParkingInfoWindow  v-for="device in getDevicesList" :activeWindow="activeWindow" :device="device"></ParkingInfoWindow>

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
import ParkingMarker from './ParkingMarker.vue';
import ParkingInfoWindow from './ParkingInfoWindow.vue';

// - Store -------------------------------------------------------------

const appStore = useAppStore();
const deviceStore = useDeviceStore();

const { getDevicesList } = storeToRefs(deviceStore);

// - Data --------------------------------------------------------------

const apiKey = ref(null);
const mapZoom = ref(17)
const mapRef = ref(null);

const activeWindow = ref(null);

// - Computed ----------------------------------------------------------

const getMapCenter = computed(() => {
    const lat = appStore.getAppSettings.default_latitude ? Number(appStore.getAppSettings.default_latitude) : 0;
    const lng = appStore.getAppSettings.default_longitude ? Number(appStore.getAppSettings.default_longitude) : 0;
    return { lat, lng };
});

const getMapId = computed(() => { return appStore.getAppSettings.google_map_id ?? ""; });



function zoomChanged() {
    mapZoom.value = mapRef.value.map.getZoom()
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