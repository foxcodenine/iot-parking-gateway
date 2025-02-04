<template>
    <div class="the-selector " :class="{'the-selector__disabled': props.isDisabled}" @click="toggleOptions" :id="fieldId">
        <label class="the-selector__label" for="type">{{ label }}<span v-if="isRequired" class="the-selector__required">*</span></label>
        <div class="the-selector__text" :class="{'the-selector__text--open': showOptions}">{{ selectedOption._value }}</div>
        <ul class="the-selector__options" v-show="showOptions">

            <li @click="showOptions = false">
                <input v-model="searchTerm" class="the-selector__input" type="text" placeholder="Filter options...">
            </li>

            <li v-for="option, index in filterOptions" :key="`${index}-${option._key}`" @click="selectOption(option)">
                {{ option._value }}
            </li>
        </ul>
        <svg class="the-selector__downarrow">
            <use xlink:href="@/assets/svg/sprite.svg#icon-arrow-9"></use>
        </svg>
    </div>
</template>

<!-- --------------------------------------------------------------- -->

<script setup>
import { computed, onMounted, ref, watch } from 'vue';

// -- Emits ------------------------------------------------------------

const emit = defineEmits(['emitOption']);

// - Props -------------------------------------------------------------
const props = defineProps({
    options: {
        type: Array,
        default: []
    },
    selectedOption: {
        type: Object,
        required: true
    },
    fieldName: {
        type: String,
        required: true
    },
    label: {
        type: String,
        required: true
    },
    isRequired: {
        type: Boolean,
        required: false
    },
    isDisabled: {
        type: Boolean,
        default: false
    }

});

// -- Data -------------------------------------------------------------

const fieldId = ref(`_${Math.floor(Math.random() * 1000) + 1}`);
const searchTerm = ref('');
const showOptions = ref(false);

// -- Computed ---------------------------------------------------------

// Creates a computed property that filters options 
// based on the search term, matching case-insensitively.
const filterOptions = computed(()=>{
    return props.options.filter(option => {
        return option._value.toLowerCase().includes(searchTerm.value.toLowerCase().trim()) ;
    });
});

// -- Watchers ---------------------------------------------------------

// Watch for changes to isDisabled prop and close options if true
watch(() => props.isDisabled, (newValue, oldValue) => {
    
    if (Boolean(newValue) === true) {
        showOptions.value = false;
    }
});

// -- Method -----------------------------------------------------------

// Toggles the visibility of the options list, 
// preventing toggle if the field is disabled
function toggleOptions() {
    if (props.isDisabled === true) {
        showOptions.value = false;
        return;
    }
    showOptions.value = !showOptions.value;
}

// Assigns the fieldName to the option, 
// and emits the 'emitOption' event with the selected option
function selectOption(option) { 
    option.fieldName = props.fieldName;
    emit('emitOption', option);
}

// -- Hooks ------------------------------------------------------------

onMounted(()=>{

    // Disable options if field is initially disabled
    if (props.isDisabled === true) {
        showOptions.value = false;
    }

    // Close options on clicks outside the selector
    document.querySelector('body').addEventListener('click', (e)=>{
        if (e.target.closest(`#${fieldId.value}`) === null) {
            showOptions.value = false;
        }
    })
});

</script>

<!-- --------------------------------------------------------------- -->
 
<style lang="scss" scoped>
.the-selector {
  position: relative;
  flex: 1;
  height: 4rem;
  overflow: visible !important;

  &__label {
    position: absolute;
    top: 0.5rem;
    left: 0.5rem;
    font-family: $font-display;
    font-size: 0.8rem;
    font-weight: 500;
    text-transform: uppercase;
    color: $col-text-2;
    z-index: 10;
    display: flex;
    align-items: center;
    gap: 0.3rem;
    line-height: 1rem;
  }

  &__text {
    cursor: text;
    width: 100%;
    height: 4rem;
    margin: 0;
    border: none;
    padding: 2rem 0.5rem 0.5rem 0.5rem;
    background-color: rgba($col-white, .5);   
    color: $col-text-1;
    font-family: $font-action;
    font-size: 1rem;
    font-weight: 300;
    border: 1px solid $col-slate-300;
    border-radius: $border-radius;

    &--open,
    &:focus {
      border: 1px solid $col-blue-500;
    }
  }

  &__downarrow {
    width: 0.7rem;
    height: 0.7rem;
    position: absolute;
    bottom: 1.7rem;
    right: 0.7rem;
    fill: currentColor;
  }

  &__input {
    width: 100%;
    padding: 0.5rem;
    border: none;
    background-color: $col-white;
    color: $col-text-1;
    font-size: 0.9rem;
    border: 1px solid transparent;

    &:focus {
      outline: none !important;
      border: 1px solid $col-blue-500;
    }
  }

  &__options {
    position: absolute;
    top: 100%;
    left: 1px;
    right: 0;
    list-style: none;
    padding: 0;
    margin: 0;
    transform: translateY(0.3rem);
    z-index: 1000;
    max-height: 300px;
    overflow: auto;
    border: 1px solid $col-slate-300;
    border-radius: $border-radius;

    @extend %custom-scrollbar;

    li {
      cursor: pointer;
      padding: 0.2rem 0.5rem;
      border-bottom: 1px solid $col-slate-300;
      background-color: $col-blue-50;
      font-size: 0.9rem;
      font-family: $font-primary;

      &:first-child {
        padding: 0;
      }

      &:last-child {
        border-bottom: none;
      }

      &:hover {
        background-color: $col-sky-200;
      }
    }
  }

  &__required {
    color: $col-red-400;
    font-size: 1rem;
    transform: translateY(1px);
  }

  &__disabled {
    .the-selector__text {
      cursor: not-allowed;
      opacity: 0.7;
      background-color: $col-slate-100;
      border-color: $col-slate-300 !important;
    }

    .the-selector__label {
      opacity: 0.5;
    }

    .the-selector__downarrow {
      fill: rgba($col-text-1, 0.4);
    }
  }
}
</style>
