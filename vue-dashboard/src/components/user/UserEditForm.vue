<template>
    <form class="fform mb-12" autocomplete="off" :class="{ 'fform__disabled': confirmOn }">
        <div class="fform__row mt-8 " @click="clearMessage">
            <div class="fform__group ">
                <label class="fform__label" for="email">Email <span class="fform__required">*</span></label>
                <input class="fform__input" id="email" type="text" placeholder="Enter email" v-model.trim="email">
            </div>

            <TheSelector :options="returnAccessLevelOptions" :selectedOption="selectedOptions['accessLevel']"
                fieldName="accessLevel" label="Access Level" :isDisabled="confirmOn" :isRequired="true"
                @emitOption="selectedOptions['accessLevel'] = $event"></TheSelector>

        </div>
    </form>
</template>

<!-- --------------------------------------------------------------- -->
 
<script setup>
import TheSelector from '@/components/commen/TheSelector.vue'
import { useMessageStore } from '@/stores/messageStore';
import { useUserStore } from '@/stores/userStore';
import { computed, onMounted, reactive, ref, watch } from 'vue';

// - Store -------------------------------------------------------------
const messageStore = useMessageStore();
const userStore = useUserStore();

// - Data --------------------------------------------------------------
const confirmOn = ref(false);
const props = defineProps({
    userID: {
        type: String,
        required: false, // Because it won't be present in the user list view
    }
});


const accessLevelList = ref([
    // {id: 0, name: 'Root'},
    { id: 1, name: 'Administrator' },
    { id: 2, name: 'Editor' },
    { id: 3, name: 'Viewer' },
]);

const selectedOptions = reactive({
    'accessLevel': { _key: 1, _value: 'Administrator' }
});

const email = ref("");
const password1 = ref("");
const password2 = ref("");
const accessLevel = ref(1);

// - computed ----------------------------------------------------------

const returnAccessLevelOptions = computed(() => {
    return accessLevelList.value.map((accessLvl => {
        // TODO: use utilStore.capitalizeFirstLetter() for value.
        return { ...accessLvl, _key: accessLvl.id, _value: accessLvl.name }
    }))
})

const getUser = computed(()=>{
    return  userStore.getUserById(Number(props.userID));
})

// - watchers -----------------------------------------------------------

watch(() => selectedOptions.accessLevel, (val, oldVal) => {
    accessLevel.value = val._key;
}, { deep: true });

watch(() => getUser, (val, oldVal) => {
    if (val.value) {
        email.value = val.value.email;
        const accessLevel = accessLevelList.value.find(accesslvl => accesslvl.id === Number(val.value.access_level));
        selectedOptions['accessLevel'] = {...accessLevel, _key: accessLevel.id, _value: accessLevel.name}
    }
}, { deep: true, immediate: true });


</script>

<!-- --------------------------------------------------------------- -->
 
<style lang="scss" scoped>
// 
</style>