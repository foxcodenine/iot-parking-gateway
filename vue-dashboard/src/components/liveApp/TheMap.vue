<template>
    <div class="the-map" >
        <GoogleMap        
            :api-key="getAppVar.googleAptKey" 
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
import { computed, toRaw } from 'vue';
import { GoogleMap, Marker } from 'vue3-google-map';

const getAppVar = computed(()=>{
    console.log(toRaw(useAppStore().appVariables))

    console.log(toRaw(useAppStore().appVariables2))
  
    return null;
});


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