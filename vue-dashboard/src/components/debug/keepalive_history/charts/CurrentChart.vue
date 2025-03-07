<template>
    <div class="heading--4 mt-16 mb-4">Current (mA) Over Time</div>
    <div class="chart-container">
        <Bar v-if="logs.length" :data="chartData" :options="chartOptions" />
    </div>
</template>

<!-- --------------------------------------------------------------- -->

<script setup>
import { computed, toRefs } from 'vue'
import { Bar } from 'vue-chartjs'
import {
    Chart as ChartJS,
    Title,
    Tooltip,
    Legend,
    BarElement,
    CategoryScale,
    LinearScale
} from 'chart.js'
import { formatToLocalDateTime } from '@/utils/dateTimeUtils'

// Register only the Chart.js components we need
ChartJS.register(Title, Tooltip, Legend, BarElement, CategoryScale, LinearScale)

// Define props so the parent can pass in logs
const props = defineProps({
    logs: {
        type: Array,
        default: () => []
    }
})

const { logs } = toRefs(props)

// Compute chart data for "current"
const chartData = computed(() => {
    const labels = logs.value.map(log => formatToLocalDateTime(log.happened_at))
    const dataset = logs.value.map(log => log.current)

    return {
        labels,
        datasets: [
            {
                label: 'Current',
                data: dataset,
                backgroundColor: '#ef4444' 
            }
        ]
    }
})

// Basic chart options
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
                    return label.slice(0, 11) // short label
                }
            }
        }
    }
}
</script>

<!-- --------------------------------------------------------------- -->

<style lang="scss" scoped>
.chart-container {
  /* Make the container taller on larger screens */
  min-height: 500px;
  position: relative; /* recommended for responsive charts */  


  @include respondDesktop(600) {
        min-height: 300px;
    }

}
</style>