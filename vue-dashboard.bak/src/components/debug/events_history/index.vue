<template>
    <div>
        <div class="fform__row mt-8 " @click="clearMessage" :class="{ 'fform__disabled': confirmOn }">
            <DeviceSelector @emitDeviceId="deviceID=$event"></DeviceSelector>
        </div>

        <RangeDatePicker  :dateRange="dateRange" @emitDateRange="dateRange=$event" ></RangeDatePicker>
        
    </div>
</template>
<!-- --------------------------------------------------------------- -->
<script setup>
import { useMessageStore } from '@/stores/messageStore';
import RangeDatePicker from '@/components/commen/RangeDatePicker.vue';
import DeviceSelector from '@/components/commen/DeviceSelector.vue';
import { ref, watch } from 'vue';


// - Store -------------------------------------------------------------
const messageStore = useMessageStore();

// - Data --------------------------------------------------------------
const confirmOn = ref(false);

const deviceID = ref(null);

const dateRange = ref({
    fromDate: new Date().setDate(new Date().getDate() + -7),
    toDate: new Date().setDate(new Date().getDate() + 0)
});

// - Methods -----------------------------------------------------------
function clearMessage() {
    messageStore.clearFlashMessage();
}


</script>

<!-- --------------------------------------------------------------- -->

<style lang="scss" scoped>
// 
</style>