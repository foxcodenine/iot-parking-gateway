<template>
    <p class="flash-message my-4 p-2" :class="flashClasses">
        {{ flashMessage }}
    </p>
</template>

<!-- --------------------------------------------------------------- -->

<script setup>
import { ref } from 'vue';


const flashMessage = ref('A device with this ID already exists in the system');
const flashClasses = ref('flash-message--orange');

</script>

<!-- --------------------------------------------------------------- -->

<style lang="scss" scoped>
// Placeholder comment to ensure global styles are imported correctly
</style>