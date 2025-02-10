<template>
    <div>
        <form>
            <AssetDateSelector></AssetDateSelector>
            <button class="bbtn bbtn--sky mt-12" @click.prevent="fetchActivityLogs">Retrieve Activity Logs</button>
        </form>
    </div>
</template>

<!-- --------------------------------------------------------------- -->

<script setup>
import { useDebugStore } from '@/stores/debugStore';
import AssetDateSelector from '../AssetDateSelector.vue';

const debugStore = useDebugStore();

async function fetchActivityLogs() {

    const response = await debugStore.fetchActivityLogs();

    if (response?.status == 200) {
            const msg = response.data?.message ?? "Activity logs retrieved successfully.";
            messageStore.setFlashMessages([msg], "flash-message--green");
            console.log(response.data)

        }

}

</script>

<!-- --------------------------------------------------------------- -->

<style lang="scss" scoped>
// 
</style>