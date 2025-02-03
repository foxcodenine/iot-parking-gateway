<template>
    <main class="vview__main">
        <section class="vview__section">
            <div class="heading--2 ">Debug</div>
            <TheFlashMessage ></TheFlashMessage>
            <EventsHistory></EventsHistory>
            <KeepAliveLogs></KeepAliveLogs>
        </section>
    </main>

</template>

<!-- --------------------------------------------------------------- -->

<script setup>
import TheFlashMessage from '@/components/commen/TheFlashMessage.vue';
import EventsHistory from '@/components/debug/events_history/index.vue'
import KeepAliveLogs from '@/components/debug/keepalive/index.vue'
</script>

<!-- --------------------------------------------------------------- -->

<style lang="scss" scoped>
// 
</style>