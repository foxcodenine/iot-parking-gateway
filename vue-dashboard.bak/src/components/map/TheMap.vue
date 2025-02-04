<template>
    <div class="the-map">
        <GoogleMap v-if="apiKey"  @click="updatedMarkerLocation"
            :api-key="apiKey"  :map-id="getMapId" style="width: 100%; height: 100vh" 
            :center="getMapCenter"  :zoom="getMapZoom" @zoom_changed="zoomChanged"  ref="mapRef">

                <ParkingMarker v-for="device in getfilteredDevices" 
                    :activeWindow="getActiveWindow" 
                    :device="device" 
                    :mapZoom="mapZoom" 
                    @click="setActiveWindow(device.device_id)"                    
                ></ParkingMarker>

                <ParkingInfoWindow  v-for="device in getfilteredDevices" 
                    :activeWindow="getActiveWindow"    
                    :device="device"
                    @emitCloseWindow="setActiveWindow(null)"
                    @emitUpdatedMarkerLocation="initUpdatedMarkerLocation"
                ></ParkingInfoWindow>

                <InfoPanel></InfoPanel>

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
import { onMounted, ref, watch } from 'vue';
import { computed, toRaw } from 'vue';
import { GoogleMap, Marker, CustomMarker } from 'vue3-google-map';
import ParkingMarker from './ParkingMarker.vue';
import ParkingInfoWindow from './ParkingInfoWindow.vue';
import { useMapStore } from '@/stores/mapStore';
import InfoPanel from './InfoPanel.vue';

// - Store -------------------------------------------------------------

const appStore = useAppStore();

const deviceStore = useDeviceStore();
const { getDevicesList, getfilteredDevices } = storeToRefs(deviceStore);

const mapStore = useMapStore();
const { getMapCenter, getMapZoom, getActiveWindow } = storeToRefs(mapStore)

// - Data --------------------------------------------------------------

const apiKey = ref(null);
const mapZoom = ref(17)
const mapRef = ref(null);

const relocateMarker = ref(null);


// - Computed ----------------------------------------------------------



const getMapId = computed(() => { 
    try {
        return appStore.getAppSettings.google_map_id ?? "";
    } catch (error) {
        // Ignore the specific error you need to suppress
        if (error.message.includes("Cannot read properties of null (reading 'google_map_id')")) {
            return ""; // Return a default value if needed
        }
        // Log other errors
        console.error('! TheMap.getMapId !\n', error);
    }
});

// - Watchers ----------------------------------------------------------

watch(() => mapRef.value?.ready, (ready) => {
    if (!ready) return;
    
    disableDefaultInfoWindow();
    setTimeout(()=>{ adjustMapView(); }, 200);
})

// - Methods -----------------------------------------------------------

function zoomChanged() {
    mapZoom.value = mapRef.value.map.getZoom()
}


function adjustMapView() {

    let devicePositions = Object.values(getDevicesList.value).map(device => ({ lat: device.latitude, lng: device.longitude }));
    let uniquePositions = [...new Set(devicePositions.map(pos => JSON.stringify(pos)))].map(str => JSON.parse(str));

    if (uniquePositions.length === 0) return;
    if (uniquePositions.length === 1) {
        mapStore.setMapCenter(uniquePositions[0]);
        mapStore.setMapZoom(17);
    }
    if (uniquePositions.length > 1) {
        let bounds = new mapRef.value.api.LatLngBounds();
        uniquePositions.forEach(pos => {
            let point = new mapRef.value.api.LatLng(pos);
            bounds.extend(point);
        });

        mapRef.value.map.fitBounds(bounds);
    }
}

function disableDefaultInfoWindow() {
    mapRef.value.map.addListener("click", (event) => {
        if (event.placeId) {
            event.stop();
        }
    })
}

function initUpdatedMarkerLocation(deviceID) {    
    setActiveWindow(null);
    relocateMarker.value = deviceID;

    mapRef.value.map.setOptions({draggableCursor:'crosshair'});
}

async function updatedMarkerLocation(e) {
    try {
        if (!relocateMarker.value) return

        const markerID = relocateMarker.value;

        const location = e.latLng.toJSON();

        const device = getDevicesList.value[markerID];
        device.latitude = location.lat;
        device.longitude = location.lng;

        await deviceStore.updateDevice(device)

        mapRef.value.map.setOptions({ draggableCursor: 'url("https://maps.gstatic.com/mapfiles/openhand_8_8.cur"), default', });
        relocateMarker.value = null;        
        setActiveWindow(markerID);
    } catch (error) {
        console.error("! TheMap.updatedMarkerLocation !\n", error);
    }
}

function setActiveWindow(id) {
    mapStore.setActiveWindow(id)
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

    appStore.getUserFavorites
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