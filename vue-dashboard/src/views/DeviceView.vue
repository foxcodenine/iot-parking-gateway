<template>
    <main class="vview__main">
        <section class="vview__section">
            <div class="heading--2 ">Devises</div>
            <TheFlashMessage ></TheFlashMessage>
            <!-- <div class="heading--4">CREATE NEW ORGANISATION</div> -->

            <!-- <FormOrganisation></FormOrganisation> -->

            <div class="heading--4">Devices List</div>      
            
            
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

const flashMessage = ref('A device with this ID already exists in the system');
const flashClasses = ref('flash-message--orange');

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