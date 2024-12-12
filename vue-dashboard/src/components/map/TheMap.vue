<template>
    <div class="the-map">
        <GoogleMap
            v-if="apiKey"
            :api-key="apiKey"
            style="width: 100%; height: 100vh"
            :center="getMapCenter"
            :zoom="15"
        >
            <Marker :options="{ position: getMapCenter }" />
        </GoogleMap>
    </div>
</template>

<!-- --------------------------------------------------------------- -->
<script setup>

import { useAppStore } from '@/stores/appStore';
import { storeToRefs } from 'pinia';
import { onMounted, ref } from 'vue';
import { computed, toRaw } from 'vue';
import { GoogleMap, Marker } from 'vue3-google-map';

// - Store -------------------------------------------------------------

const appStore = useAppStore();

// - Data --------------------------------------------------------------

const apiKey = ref(null);

// - Computed ----------------------------------------------------------

const getMapCenter = computed(()=>{
    const lat = appStore.getAppSettings.default_latitude ? Number( appStore.getAppSettings.default_latitude) : 0;
    const lng = appStore.getAppSettings.default_longitude ? Number( appStore.getAppSettings.default_longitude) : 0;
    return { lat, lng };
});

// - Hooks -------------------------------------------------------------

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