<template>
        <div class="my-date-picker__container mt-4">
            <div class="my-date-picker">
                <DatePicker 
                v-model.range="range" 
                borderless 
                expanded 
                :columns="columns" 
                :rows="rows" 
                mode="dateTime" is24hr
                hide-time-header
                />
            </div>
        </div>
</template>

<!-- --------------------------------------------------------------- -->

<script setup>
import { Calendar, DatePicker } from 'v-calendar';
import { useScreens } from 'vue-screen-utils';
import 'v-calendar/style.css';
import { ref, watch } from 'vue';

const emit = defineEmits(["emitDateRange"])

const props = defineProps({
    dateRange: {
        type: Object,
        default: {
            fromDate: new Date().setDate(new Date().getDate() + -3),
            toDate: new Date().getTime(),
        }
    },
});

const range = ref({
  start: props.dateRange.fromDate,
  end: props.dateRange.toDate,
});

const { mapCurrent } = useScreens({
  xs: '0',
  sm: '640px',
  md: '768px',
  lg: '1024px',
})

const columns = mapCurrent( {sm:2, md: 2, lg: 2 }, 1);
const rows = mapCurrent( {xs:2, }, 1);

watch(range, (val) => {

    const payload = {
        fromDate: new Date(val.start).getTime(),
        toDate: new Date(val.end).getTime(),
    }
    emit("emitDateRange", payload);    
}, {
    immediate:true
})


</script>

<!-- --------------------------------------------------------------- -->

<style lang="scss" scoped>
.my-date-picker {
    min-width: 100%;
    padding: 0;

    &__container {
        overflow: hidden;
        border: 1px solid $col-slate-300;
        border-radius: $border-radius;
    }
    
    & >* {
        @include respondDesktop(385) {
            margin: 0 -19px 0 -15px;
        }
    }
}
</style>

<style>
.vc-base-select select { border: none; }
</style>
