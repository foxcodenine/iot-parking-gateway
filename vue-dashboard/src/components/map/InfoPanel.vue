<template>
    <div class="info-panel">
        
        <div class="info-panel__selector">
        </div>

        <div class="info-panel__action">
        </div>

        <ParkingPanel></ParkingPanel>

    </div>
</template>
<!-- --------------------------------------------------------------- -->
<script setup>
import ParkingPanel from './info-panel/ParkingPanel.vue';
</script>
<!-- --------------------------------------------------------------- -->
<style lang="scss" scoped>
.info-panel {
    width: 18rem;

    position:  absolute;
    right: 0;
    top: 0.5rem;
    bottom: 0.5rem;
    border-top-left-radius: 3px;
    border-bottom-left-radius: 3px;

    background-color: rgba($col-zinc-900, 1);
    opacity: .9;


    // transform: translateX(17.5rem);
}


</style>