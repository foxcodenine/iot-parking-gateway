<template>
    <CustomMarker  :style="markerStyle()" 
        
        @click="aaa(device)"
        :options="{ position: { lat: device.latitude, lng: device.longitude }, anchorPoint: 'CENTER' }">
        <svg class="vacant" :class="{'occupied': device.is_occupied}" :width="markerSize" :height="markerSize" viewBox="0 0 100 100" xmlns="http://www.w3.org/2000/svg">
            <circle cx="50" cy="50" r="48" fill="#ffffff" stroke="currentColor" stroke-width="4" />
            <circle cx="50" cy="50" r="40" fill="currentColor" />
            <text x="53" y="70" font-family="Arial, sans-serif" font-size="60" fill="#ffffff" font-weight="bold" text-anchor="middle">P</text>
        </svg>
    </CustomMarker>
</template>

<!-- --------------------------------------------------------------- -->

<script setup>
import { computed, onMounted } from 'vue';
import { GoogleMap, Marker, CustomMarker } from 'vue3-google-map';

const props = defineProps({
    mapZoom: {
        type: Number,
        required: true,  
    },
    device: {
        type: Object,
        required: true,  
    }
});

function markerStyle() {
    return "opacity: 1;"
}

const markerSize = computed(() => {
    const zoom = Math.round(props.mapZoom);
    return Math.min(10 + (zoom - 10) * 3, 64); // Linearly scales marker size
});

function aaa(d){
    console.log(d)
}

onMounted(()=>{
  
})

</script>

<!-- --------------------------------------------------------------- -->
 
<style lang="scss" scoped>

.vacant {
    color: $col-blue-600;
}
.occupied {
    color: $col-red-500 !important;
}


</style>