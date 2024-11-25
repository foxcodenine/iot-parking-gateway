<template>
    <div class="the-map" >
        <GoogleMap        
            :api-key="getGoogleApiKey" 
            style="width: 100%; height: 100vh" 
            :center="center" 
            :zoom="15">
            <Marker :options="{ position: center }" />
        </GoogleMap>
    </div>
</template>

<!-- --------------------------------------------------------------- -->
<script setup>

import { useAppStore } from '@/stores/appStore';
import { storeToRefs } from 'pinia';
import { computed, toRaw } from 'vue';
import { GoogleMap, Marker } from 'vue3-google-map';

const appStore = useAppStore();

const { getGoogleApiKey } = storeToRefs( appStore );

console.log('>>', getGoogleApiKey.value)


const center = { lat: 40.689247, lng: -74.044502 };

</script>

<!-- --------------------------------------------------------------- -->

<style lang="scss" scoped>
.the-map {
    background-color: $col-blue-300;
    // background-image: url("@/assets/images/map.jpg");
    // min-height: 100vh;
}
</style>