<template>
        <TheSelector 
            :options="returnDeviceList" 
            :selectedOption="selectedOptions['deviceList']"
            fieldName="deviceList" 
            label="Devices" 
            :isDisabled="confirmOn" 
            :isRequired="true"
            @emitOption="updateSelectedOptions"
        ></TheSelector>
      
</template>

<!-- --------------------------------------------------------------- -->

<script setup>
import { useDeviceStore } from '@/stores/deviceStore';
import { useMessageStore } from '@/stores/messageStore';
import TheSelector from '@/components/commen/TheSelector.vue'
import { storeToRefs } from 'pinia';
import { onMounted, reactive, ref, watch } from 'vue';
import { computed } from 'vue';


const emit = defineEmits(['emitDeviceId'])

// - Store -------------------------------------------------------------
const messageStore = useMessageStore();
const deviceStore = useDeviceStore();


// - Data --------------------------------------------------------------
const confirmOn = ref(false);

const selectedOptions = reactive({
    'deviceList': { _key: '', _value: null }
});


// - Computed ----------------------------------------------------------
const returnDeviceList = computed(() => {    
 
    return Object.values(deviceStore.getDevicesList).map((device => {
    
        // TODO: use utilStore.capitalizeFirstLetter() for value.
        return { ...device, _key: device.device_id, _value: `<b>${device.device_id }</b><br/><i>${device.name }</i>`}
    }))
});



// - Methods -----------------------------------------------------------

function updateSelectedOptions(payload) {
    if (payload.fieldName == "deviceList") {    
        selectedOptions['deviceList'] = {...payload, _value: payload._key};
        emit('emitDeviceId', payload.device_id)
    }
}

onMounted(()=>{
    setTimeout(()=>{
        console.log();
    }, 2000)
})

</script>

<!-- --------------------------------------------------------------- -->

<style lang="scss" scoped>
// 
</style>