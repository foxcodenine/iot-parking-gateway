<template>
    <div>
        <div class="fform__row mt-8 " @click="clearMessage" :class="{ 'fform__disabled': confirmOn }">
            <DeviceSelector :deviceID="deviceID" @emitDeviceId="deviceID=$event"></DeviceSelector>
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
import { useRoute, useRouter } from 'vue-router';

// - Router ------------------------------------------------------------

const route = useRoute();
const router = useRouter();

// - Store -------------------------------------------------------------
const messageStore = useMessageStore();

// - Data --------------------------------------------------------------
const confirmOn = ref(false);
const deviceID = ref(null);

const dateRange = ref({
    fromDate: new Date().setDate(new Date().getDate() + -7),
    toDate: new Date().setDate(new Date().getDate() + 0)
});


// - Watchers ----------------------------------------------------------

// Watch for changes in the route parameter (:id)
// If the route ID is different from the current `deviceID`, update `deviceID`
watch(() => route.params.id, (newId) => {
    if (newId != deviceID.value) {
        deviceID.value = newId;
    }
}, { immediate: true });

// Watch for changes in `deviceID` (selected device)
// If `deviceID` changes and is different from the current route, update the route without refreshing
watch(deviceID, (newId) => {
    if (newId != route.params.id) {
        router.replace(`/debug/${newId}`);        
    }
}, { immediate: true });

// - Methods -----------------------------------------------------------
function clearMessage() {
    messageStore.clearFlashMessage();
}


</script>

<!-- --------------------------------------------------------------- -->

<style lang="scss" scoped>
// 
</style>