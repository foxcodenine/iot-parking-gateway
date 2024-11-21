<template>
    <div class="the-map" >
        <GoogleMap        
            :api-key="googleAptKey" 
            style="width: 100%; height: 100vh" 
            :center="center" 
            :zoom="15">
            <Marker :options="{ position: center }" />
        </GoogleMap>
    </div>
</template>

<!-- --------------------------------------------------------------- -->
<script setup>

import { ref } from 'vue';
import { GoogleMap, Marker } from 'vue3-google-map';

const googleAptKey = ref(GO_GOOGLE_API_KEY);

if (import.meta.env.VITE_VUE_ENV == 'development') {
    googleAptKey.value = import.meta.env.VITE_GOOGLE_API_KEY
} else {
    googleAptKey.value = GO_GOOGLE_API_KEY;
    GO_GOOGLE_API_KEY = "";
}

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