<template>
    <div class="heading--4 mt-16 mb-4">Temperature Min/Max Trend</div>
    <div class="chart-container">
        <Line v-if="logs.length" :data="chartData" :options="chartOptions" />
    </div>    
</template>

<!-- --------------------------------------------------------------- -->

<script setup>
import { computed, toRefs } from 'vue'
import { Line } from 'vue-chartjs'
import {
    Chart as ChartJS,
    Title,
    Tooltip,
    Legend,
    PointElement,
    LineElement,
    CategoryScale,
    LinearScale
} from 'chart.js'
import { formatToLocalDateTime } from '@/utils/dateTimeUtils'

// Register line-related components
ChartJS.register(Title, Tooltip, Legend, PointElement, LineElement, CategoryScale, LinearScale)

const props = defineProps({
    logs: {
        type: Array,
        default: () => []
    }
})

const { logs } = toRefs(props)

// Compute chart data for temperature_min and temperature_max
const chartData = computed(() => {
    const labels = logs.value.map(log => formatToLocalDateTime(log.happened_at))
    const tempMinData = logs.value.map(log => log.temperature_min)
    const tempMaxData = logs.value.map(log => log.temperature_max)

    return {
        labels,
        datasets: [
            {
                label: 'Temperature Min',
                data: tempMinData,
                borderColor: '#34d399', 
                backgroundColor: '#ffffff',
                fill: false,
                tension: 0.1
            },
            {
                label: 'Temperature Max',
                data: tempMaxData,
                borderColor: '#ef4444', 
                backgroundColor: '#ffffff'
,
                fill: false,
                tension: 0.1
            }
        ]
    }
})

const chartOptions = {
    responsive: true,
    maintainAspectRatio: false,
    scales: {
        x: {
            ticks: {
                maxRotation: 90,
                minRotation: 0,
                maxTicksLimit: 10,
                callback(value) {
                    const label = this.getLabelForValue(value)
                    return label.slice(0, 11)
                }
            }
        }
    }
}
</script>

<!-- --------------------------------------------------------------- -->

<style scoped lang="scss">
.chart-container {
  /* Make the container taller on larger screens */
  min-height: 500px;
  position: relative; /* recommended for responsive charts */  


  @include respondDesktop(600) {
        min-height: 300px;
    }

}
</style>