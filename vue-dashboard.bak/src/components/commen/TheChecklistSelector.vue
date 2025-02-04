<template>
    <div class="checklist">

        <input class="checklist__search" v-if="addSearchBar"  v-model="searchTerm" type="text" placeholder="Search..." :disabled="isDisabled">

        <ul class="checklist__list" :class="{'checklist__list--flex': props.displayFlex}">
            <li  class="checklist__item checklist__item--all" @click="toggleCheckAll" v-if="!searchTerm && props.allCheckbox" >
                <svg class="checklist__icon">             
                    <use xlink:href="@/assets/svg/sprite.svg#icon-checkbox-3" v-if="checkAll"></use>          
                    <use xlink:href="@/assets/svg/sprite.svg#icon-checkbox-6"></use>          
                </svg>
                <div class="checklist__text">
                    All
                </div>
            </li>
            <li  class="checklist__item" v-for="item in  filteredChecklistItems" @click="updateCheckItems(item._key)">
                <svg class="checklist__icon">             
                    <use xlink:href="@/assets/svg/sprite.svg#icon-checkbox-3" v-if="props.checkedItems.includes(item._key)"></use>          
                    <use xlink:href="@/assets/svg/sprite.svg#icon-checkbox-6"></use>          
                </svg> 
                <div class="checklist__text">
                    {{ item._value }}
                </div>
            </li>     
        </ul>

    </div>
</template>

<!-- --------------------------------------------------------------- -->

<script setup>
import { computed, onMounted, ref, toRaw, watch } from 'vue';

// - Emit --------------------------------------------------------------

const emit = defineEmits(['emitNewCheckedItems'])

// - Props -------------------------------------------------------------
const props = defineProps({
    checklistItems: {
        type: Array,
        required: false,
        default: []
    },
    checkedItems: {
        type: Array,
        required: false,
        default: [],
    },
    isDisabled: {
        type: Boolean,
        default: false
    },
    allCheckbox: {
        type: Boolean,
        default: false
    },
    displayFlex: {
        type: Boolean,
        default: false
    },
    addSearchBar: {
        type: Boolean,
        default: false
    },
});
// - Data --------------------------------------------------------------

const checkAll = ref(false);
const allKeys = ref([]);
const searchTerm = ref("");

// - Computed -------------------------------------------------------------

const filteredChecklistItems = computed(() => {
    if (!searchTerm.value) {
        return props.checklistItems;
    }
    return props.checklistItems.filter(item => 
        item._value.toLowerCase().includes(searchTerm.value.toLowerCase())
    );
});

// - Watch -------------------------------------------------------------

watch(() => props.checkedItems, (newCheckedItems) => {
    checkAll.value = newCheckedItems.length === allKeys.value.length;
}, { deep: true });


// - Methods -----------------------------------------------------------

function updateCheckItems(itemKey) {
    if (props.isDisabled) return;

    const index = props.checkedItems.indexOf(itemKey);
    const newCheckedItems = [...props.checkedItems];

    if (index > -1) {
        newCheckedItems.splice(index, 1);
    } else {
        newCheckedItems.push(itemKey);
    }
    
    emit('emitNewCheckedItems', newCheckedItems);
}

function toggleCheckAll() {
    if (props.isDisabled) return;

    checkAll.value = !checkAll.value;
    if (checkAll.value) {
        emit('emitNewCheckedItems', allKeys.value);        
    } else {
        emit('emitNewCheckedItems', []);
    }
}

// - Hooks -------------------------------------------------------------

onMounted(()=>{
    allKeys.value = props.checklistItems.map(item => item._key);
});

</script>

<!-- --------------------------------------------------------------- -->

<style lang="scss" scoped>
.checklist {
    flex: 1;
    &__search {
        padding: 0.3rem 0.5rem;
        border: 1px solid $col-slate-300;
        border-radius: $border-radius;
        font-size: 0.9rem;
        width: 100%; 

        &:focus {
            border: 1px solid $col-blue-500;
        }
    }
    &__list {
        width: calc(100%);
        margin-top: 1rem;
        
        display: grid;
        grid-template-columns: repeat(auto-fit, minmax(15rem, 1fr));
        gap: 1rem;

        &--flex {
            display: flex;
            justify-content: space-between;
            gap: 0.5rem;
            flex-wrap: wrap
        }   
    }

    &__item {
        cursor: pointer;
        display: flex;
        align-items: center;
        gap: .5rem;
        color: $col-text-1;

        &--all {
            color: $col-red-600;
        }
    }

    &__icon {
        width: 1rem;
        height: 1rem;
        fill: currentColor;
    }

    &__text {
        font-size: 1rem;
        font-family: $font-primary;
    }
}
</style>