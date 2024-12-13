<template>
    <main class="vview__main">
        <section class="vview__section">
            <div class="heading--2 ">Application Settings</div>
            <TheFlashMessage ></TheFlashMessage>  
            <TheTabs
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
import { markRaw, reactive } from 'vue';

import AdminLevelSettings from '@/components/settings/AdminLevelSettings.vue';
import RootLevelSettings from '@/components/settings/RootLevelSettings.vue';

// Mapping of tab identifiers to component objects for dynamic loading
const componentMap = reactive({
    AdminLevelSettings: markRaw(AdminLevelSettings),
    RootLevelSettings: markRaw(RootLevelSettings),

});

// Reactive state to manage tab data and the active tab identifier
const tabsObjectData_1 = reactive({
    activeTab: 'AdminLevelSettings',
    tabs: {
        AdminLevelSettings: 'Admin Level Settings',
        RootLevelSettings: 'Root Level Settings',    
    },
});

</script>

<!-- --------------------------------------------------------------- -->

<style lang="scss" scoped>
// 
</style>