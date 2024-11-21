<template>
    <form class="fform" autocomplete="off">

        <div class="fform__row mt-8">
            <div class="fform__group fform__disabled">
                <label class="fform__label" for="name">Name <span class="fform__required">*</span></label>
                <input class="fform__input" id="name" type="text" placeholder="Enter organisation name">
            </div>
            <TheSelector
                :options="returnOrganisationOptions"
                :selectedOption="selectedOptions['organisation']"
                fieldName="organisation"
                label="Organisation"
                :isDisabled="false"
                :isRequired="true"
                @emitOption="selectedOptions['organisation']=$event"
            ></TheSelector>
        </div>

        <button class="bbtn bbtn--blue mt-8">Create Organisation</button>

    </form>
</template>

<!-- --------------------------------------------------------------- -->

<script setup>
import { computed, reactive, ref } from 'vue';
import TheSelector from '@/components/commen/TheSelector.vue'


// -- organisation selector --------------------------------------------

const organisationList = ref([
    {id: 1, name: 'IoT Solutions Ltd'},
    {id: 2, name: 'Joseph Cachia & Son Ltd (JCS)'},
    {id: 3, name: 'Caruana Aluminium Systems Ltd'},
    {id: 4, name: 'Agenzija Sapport'},
    {id: 5, name: 'MZ Fuels Ltd'},
    {id: 6, name: 'IoT Solutions Ltd'},
    {id: 7, name: 'Joseph Cachia & Son Ltd (JCS)'},
    {id: 8, name: 'Caruana Aluminium Systems Ltd'},
    {id: 9, name: 'Agenzija Sapport'},
    {id: 10, name: 'MZ Fuels Ltd'},
]);

const selectedOptions = reactive({
    'organisation': { _key: 1, _value: 'IoT Solutions Ltd' }
})

const returnOrganisationOptions = computed(()=>{
    return organisationList.value.map((org => {
        // TODO: use utilStore.capitalizeFirstLetter() for value.
        return {...org, _key: org.id, _value: org.name}
    }))
})

// ---------------------------------------------------------------------



</script>

<!-- --------------------------------------------------------------- -->

<style lang="scss" scoped>
// Placeholder comment to ensure global styles are imported correctly


</style>