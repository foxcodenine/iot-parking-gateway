<template>
    <form class="fform" autocomplete="off" :class="{'fform__disabled': confirmOn}">

        <div class="fform__row mt-8 ">
            <div class="fform__group ">
                <label class="fform__label" for="email">Email <span class="fform__required">*</span></label>
                <input class="fform__input" id="email" type="text" placeholder="Enter email" v-model="email">
            </div>
            
            <TheSelector
                :options="returnAccessLevelOptions"
                :selectedOption="selectedOptions['accessLevel']"
                fieldName="accessLevel"
                label="Access Level"
                :isDisabled="confirmOn"
                :isRequired="true"
                @emitOption="selectedOptions['accessLevel'] = $event"
            ></TheSelector>

        </div>
        <div class="fform__row mt-4">
            <div class="fform__group">
                <label class="fform__label" for="password1">Password<span class="fform__required">*</span></label>
                <input class="fform__input" id="password1" type="text" placeholder="Enter password" v-model="password1">
            </div>
            <div class="fform__group">
                <label class="fform__label" for="password2">Confirm Password<span class="fform__required">*</span></label>
                <input class="fform__input" id="password2" type="text" placeholder="Renter password" v-model="password2">
            </div>
            

        </div>
        <div class="fform__row">
            <button class="bbtn bbtn--blue mt-8" v-if="!confirmOn" @click="initCreateUser()">Create User</button>

            <button class="bbtn bbtn--zinc-lt mt-8" v-if="confirmOn" @click="confirmOn=false">Cancel</button>
            <button class="bbtn bbtn--blue mt-8" v-if="confirmOn">Confirm</button>

        </div>

    </form>
</template>

<!-- --------------------------------------------------------------- -->
<script setup>
import TheSelector from '@/components/commen/TheSelector.vue'
import { computed, reactive, ref } from 'vue';

const confirmOn = ref(false);

const accessLevelList = ref([
    // {id: 0, name: 'Root'},
    {id: 1, name: 'Administrator'},
    {id: 2, name: 'Editor'},
    {id: 3, name: 'Viewer'},
]);

const selectedOptions = reactive({
    'accessLevel': { _key: 1, _value: 'Administrator' }
});

const email = ref(null);
const password1 = ref(null);
const password2 = ref(null);

// - computed ----------------------------------------------------------

const returnAccessLevelOptions = computed(()=>{
    return accessLevelList.value.map((org => {
        // TODO: use utilStore.capitalizeFirstLetter() for value.
        return {...org, _key: org.id, _value: org.name}
    }))
})

// - methods -----------------------------------------------------------

function initCreateUser() {

    confirmOn.value = true;
}

</script>
<!-- --------------------------------------------------------------- -->

<style lang="scss" scoped>
// Placeholder comment to ensure global styles are imported correctly


</style>