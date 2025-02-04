<template>
    <div class="ccheckbox mt-4" @click="clickCheckbox">
        <svg class="ccheckbox__icon" v-if="props.isChecked">
            <use xlink:href="@/assets/svg/sprite.svg#icon-checked"></use>
        </svg>
        <svg class="ccheckbox__icon" v-else>
            <use xlink:href="@/assets/svg/sprite.svg#icon-unchecked"></use>
        </svg>
        <p class="ccheckbox__text"><slot></slot></p>
    </div>
</template>

<!-- --------------------------------------------------------------- -->

<script setup>

const props = defineProps({
    isChecked: {
        type: Boolean,
        required: true,        
    },
});

const emit = defineEmits(['emitCheckbox'])

function clickCheckbox() { emit('emitCheckbox') }

</script>

<!-- --------------------------------------------------------------- -->

<style lang="scss" scoped>

.ccheckbox {
    display: flex;
    gap: .2rem;
    cursor: pointer;
    fill: $col-blue-600;
    color: $col-blue-600;
    width: max-content;

    &__icon {
        width: 1.6rem;
        height: 1.6rem;
    }

    &__text {
        font-family: $font-primary;
        font-size: 1rem;
        transform: translateY(0.5px);
    }
}

</style>