<template>
    <div class="heading--4 mt-16 mb-4">Idle Voltage Over Time</div>
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

ChartJS.register(Title, Tooltip, Legend, BarElement, CategoryScale, LinearScale)

const props = defineProps({
    logs: {
        type: Array,
        default: () => []
    }
})

const { logs } = toRefs(props)

// Compute chart data for "idle_voltage"
const chartData = computed(() => {
    const labels = logs.value.map(log => formatToLocalDateTime(log.happened_at))
    const dataset = logs.value.map(log => log.idle_voltage)

    return {
        labels,
        datasets: [
            {
                label: 'Idle Voltage',
                data: dataset,
                backgroundColor: '#34d399'
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