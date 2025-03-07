<template>
    <div>
  
        <FirmwareVersion :logs="logs" />

        <CurrentChart :logs="logs" />

        <IdleVoltageChart :logs="logs" />

        <TemperatureChart :logs="logs" />
    </div>
</template>

<!-- --------------------------------------------------------------- -->

<script setup>
import { computed } from 'vue'
import { storeToRefs } from 'pinia'
import { useDebugStore } from '@/stores/debugStore'

// Import our sub-components from the charts directory
import FirmwareVersion from './charts/FirmwareVersion.vue'
import CurrentChart from './charts/CurrentChart.vue'
import IdleVoltageChart from './charts/IdleVoltageChart.vue'
import TemperatureChart from './charts/TemperatureChart.vue'

// Access logs from your debug store (Pinia)
const debugStore = useDebugStore()
const { getKeepaliveLogs } = storeToRefs(debugStore)

// Aliased for convenience
const logs = computed(() => getKeepaliveLogs.value)
</script>

<!-- --------------------------------------------------------------- -->


<style scoped>
/* Optional styling */
</style>