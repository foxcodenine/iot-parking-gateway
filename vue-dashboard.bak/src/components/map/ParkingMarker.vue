<template>
    <CustomMarker  :style="markerStyle()" 
        
        :options="{ position: { lat: device.latitude, lng: device.longitude }, anchorPoint: 'CENTER' }">
        <svg  :class="markersClass" :width="markerSize" :height="markerSize" viewBox="0 0 100 100" xmlns="http://www.w3.org/2000/svg">
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
import { formatToLocalDateTime, timeElapsed } from '@/utils/utils';

// - Props -------------------------------------------------------------

const props = defineProps({
    mapZoom: {
        type: Number,
        required: true,  
    },
    device: {
        type: Object,
        required: true,  
    },
    activeWindow: {
        type: [Number, null, String],
        required: true,
    }
});
// - Computed ----------------------------------------------------------

const markerSize = computed(() => {
    const zoom = Math.round(props.mapZoom);
    return Math.min(10 + (zoom - 10) * 3, 64); // Linearly scales marker size
});

const markersClass = computed(()=>{

    let c = '';

    if (props.activeWindow == props.device.device_id) {
        c += 'selected-device ';
    }
    

    if (formatToLocalDateTime(props.device.happened_at) == null) {
        c += 'unknown'
        return c;
    }
    if (props.device.is_occupied) {
        c += 'occupied'
        return c;
    } else {
        c += 'vacant'
        return c;
    }

    // also i would like to add 


})

// - Methods -----------------------------------------------------------
function markerStyle() {
    return "opacity: 1; "
}

</script>

<!-- --------------------------------------------------------------- -->
 
<style lang="scss" scoped>

.vacant {
    color: $col-green-600;
}
.occupied {
    color: $col-red-500;
}
.unknown {
    color: $col-indigo-500;
}



.selected-device {
  display: inline-block;
  background-color: inherit; // Choose a color that stands out
//   padding: 10px;
  border-radius: 50%;
  animation: pulse-animation 2s infinite;
  
  @keyframes pulse-animation {
    0% {
      box-shadow: 0 0 0 0 currentColor;
    }
    70% {
      box-shadow: 0 0 0 10px rgba(52, 152, 219, 0);
    }
    100% {
      box-shadow: 0 0 0 0 rgba(52, 152, 219, 0);
    }
  }
}


</style>