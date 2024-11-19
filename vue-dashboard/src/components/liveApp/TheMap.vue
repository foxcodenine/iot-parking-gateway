<template>
    <div class="the-map" >
        <GoogleMap 
            v-if="loadMap"
            :api-key="getAptKey.GOOGLE_API_KEY" 
            style="width: 100%; height: 100vh" 
            :center="center" 
            :zoom="15">
            <Marker :options="{ position: center }" />
        </GoogleMap>
    </div>
</template>

<!-- --------------------------------------------------------------- -->
<script setup>
import { useEnvStore } from '@/stores/envStore';
import { computed, onMounted, ref } from 'vue';
import { GoogleMap, Marker } from 'vue3-google-map';



const loadMap = ref(false);
const center = { lat: 40.689247, lng: -74.044502 };

const getAptKey = computed(()=>{
    return useEnvStore().getEnv
});



(async()=>{
	await useEnvStore().loadEnvironmentVariables();
	loadMap.value = true;
})()

</script>

<!-- --------------------------------------------------------------- -->

<style lang="scss" scoped>
.the-map {
    background-color: $col-blue-300;
    // background-image: url("@/assets/images/map.jpg");
    // min-height: 100vh;
}
</style>