<template>
    <main class="vview__main">
        <section class="vview__section">
            <div class="heading--2 ">Devises</div>
            <p class="flash-message my-4 p-2" :class="flashClasses">
                {{ flashMessage }}
            </p>
            <!-- <div class="heading--4">CREATE NEW ORGANISATION</div> -->

            <!-- <FormOrganisation></FormOrganisation> -->

            <div class="heading--4 mt-10">Devices List</div>      
            
            
            <TableDevice></TableDevice>
            
        </section>
    </main>
</template>

<!-- --------------------------------------------------------------- -->

<script setup>
import { ref } from 'vue';
// import FormOrganisation from '@/components/organisation/FormOrganisation.vue'
import TableDevice from '@/components/device/TableDevice.vue'



const flashMessage = ref('A device with this ID already exists in the system');
const flashClasses = ref('flash-message--orange');

</script>

<!-- --------------------------------------------------------------- -->

<style lang="scss" scoped>
// Placeholder comment to ensure global styles are imported correctly
</style>