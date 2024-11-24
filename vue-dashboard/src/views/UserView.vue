<template>
    <main class="vview__main">
        <section class="vview__section">
            <div class="heading--2 ">Users</div>
            <TheFlashMessage></TheFlashMessage>
            <!-- <div class="heading--4">CREATE NEW ORGANISATION</div> -->

            <!-- <FormOrganisation></FormOrganisation> -->

                
            
            <UserForm></UserForm>
            
            
        </section>
    </main>
</template>

<!-- --------------------------------------------------------------- -->

<script setup>
import { ref } from 'vue';
// import FormOrganisation from '@/components/organisation/FormOrganisation.vue'
import UserForm from '@/components/user/UserForm.vue'
import TheFlashMessage from '@/components/commen/TheFlashMessage.vue';



</script>

<!-- --------------------------------------------------------------- -->

<style lang="scss" scoped>
// Placeholder comment to ensure global styles are imported correctly
</style>