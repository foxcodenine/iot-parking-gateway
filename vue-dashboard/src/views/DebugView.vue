<template>
    <main class="vview__main">
        <section class="vview__section">
            <div class="heading--2 ">Debug</div>
            <TheFlashMessage ></TheFlashMessage>


            <TheTabs
                class="mt-16"
                :tabsObjectData="tabsObjectData_1"
                :isDisabled="false"
                :layoutBreakpoint="500"
                @setActiveTab="tabsObjectData_1.activeTab = $event"
            ></TheTabs>

            <KeepAlive>
                <component :is="componentMap[tabsObjectData_1.activeTab]" class="mt-6"></component>
            </KeepAlive> 
        </section>
    </main>

</template>

<!-- --------------------------------------------------------------- -->

<script setup>
import TheTabs from '@/components/commen/TheTabs.vue'
import TheFlashMessage from '@/components/commen/TheFlashMessage.vue';
import ParkingHistory from '@/components/debug/parking_history/index.vue'
import KeepaliveHistory from '@/components/debug/keepalive_history/KeepaliveHistory.vue'
import { markRaw, reactive } from 'vue';
import { useDeviceStore } from '@/stores/deviceStore';


// - Store -------------------------------------------------------------
const deviceStore = useDeviceStore();

// - Data --------------------------------------------------------------
// Mapping of tab identifiers to component objects for dynamic loading
const componentMap = reactive({
    ParkingHistory: markRaw(ParkingHistory),
    KeepaliveHistory: markRaw(KeepaliveHistory),

});



// Reactive state to manage tab data and the active tab identifier
const tabsObjectData_1 = reactive({
    activeTab: 'ParkingHistory',
    tabs: {
        ParkingHistory: 'History',
        KeepaliveHistory: 'Keepalive',    
    },
});


// - Hooks -------------------------------------------------------------

try {
    deviceStore.fetchDevices()
} catch (error) {
    console.error('! DeviceView deviceStore.fetchDevices() !\n', error);
}

</script>

<!-- --------------------------------------------------------------- -->

<style lang="scss" scoped>
// 
</style>