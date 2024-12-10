<template>
    <main class="vview__main">
        <section class="vview__section">

            <div class="heading--2 ">Devises</div>
            <TheFlashMessage ></TheFlashMessage>

            <div class="heading--4" v-if="getUserAccessLevel <= 1">Create new device</div>
            <DeviceForm  v-if="getUserAccessLevel <= 2"></DeviceForm>

            <div class="heading--4 mt-8">Devices list</div>
            <DeviceTable></DeviceTable>
            
        </section>
    </main>
</template>

<!-- --------------------------------------------------------------- -->

<script setup>
import { ref } from 'vue';
// import FormOrganisation from '@/components/organisation/FormOrganisation.vue'
import DeviceTable from '@/components/device/DeviceTable.vue';
import TheFlashMessage from '@/components/commen/TheFlashMessage.vue';
import { useAuthStore } from '@/stores/authStore';
import { useDeviceStore } from '@/stores/deviceStore';
import { storeToRefs } from 'pinia';
import DeviceForm from '@/components/device/DeviceForm.vue';


// - Store -------------------------------------------------------------
const authStore = useAuthStore();
const deviceStore = useDeviceStore();

const { getUserAccessLevel } = storeToRefs(authStore)

// - Hooks -------------------------------------------------------------

try {
    deviceStore.fetchDevices()
} catch (error) {
    console.error('! DeviceView deviceStore.fetchDevices() !\n', error);
}

</script>

<!-- --------------------------------------------------------------- -->

<style lang="scss" scoped>
// Placeholder comment to ensure global styles are imported correctly
</style>