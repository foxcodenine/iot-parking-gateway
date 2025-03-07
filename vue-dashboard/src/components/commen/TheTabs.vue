<template>
    <div ref="theTabs" class="tabs-mb" :class="{ 'tabs-dt': isDesktop }">
        <div v-for="(value, key) in props.tabsObjectData.tabs" :key="key" @click="emitActiveTab(key)"
        class="tabs-mb__item" :class="{
            'tabs-dt__item': isDesktop,
            'tabs-mb__item--active': tabsObjectData.activeTab === key && !isDesktop,
            'tabs-dt__item--active': tabsObjectData.activeTab === key && isDesktop
        }">
            {{ value }}
        </div>
        <p class="tabs-mb__empty" :class="{ 'tabs-dt__empty': isDesktop }"></p>
    </div>
</template>

<!-- --------------------------------------------------------------- -->

<script setup>
import { computed, onMounted, onUnmounted, ref } from 'vue';


// -- Emits ------------------------------------------------------------

const emit = defineEmits(['setActiveTab']);

// - Props -------------------------------------------------------------
const props = defineProps({

    tabsObjectData: {
        type: Object,
        required: true
    },
    isDisabled: {
        type: Boolean,
        default: false
    },
    layoutBreakpoint: {
        type: Number,
        default: 600
    },

});

// -- Data -------------------------------------------------------------

const theTabs = ref(null);
const tabsWidth = ref(700);

// -- Computed ---------------------------------------------------------

const isDesktop = computed(() => tabsWidth.value >= props.layoutBreakpoint);

// -- Methods for Responsiveness and Resizing ----------------------------

function updateWidth () {
    if (theTabs.value) {
        tabsWidth.value = theTabs.value.offsetWidth;
    }
};

function setupResizeObserver() {
    const resizeObserver = new ResizeObserver(() => {
        requestAnimationFrame(() => {
            updateWidth();
        });
    });

    if (theTabs.value) {
        resizeObserver.observe(theTabs.value);
    }

    onUnmounted(() => {
        resizeObserver.disconnect();
    });
};


// -- Other Methods ----------------------------------------------------

function emitActiveTab(tab) {
    if (props.isDisabled) return;
    emit('setActiveTab', tab);
}

// -- Hooks ------------------------------------------------------------

onMounted(() => {
    updateWidth();
    setupResizeObserver();
});



</script>

<!-- --------------------------------------------------------------- -->


<style lang="scss" scoped>
.tabs-mb,
.tabs-dt {
    display: flex;
    flex-direction: column;
    font-size: 1.1rem;
}

.tabs-mb {
    &__empty {
        display: none;
    }

    &__item {
        cursor: pointer;
        font-family: $font-action;
        color: $col-zinc-400;
        font-weight: 400;
        padding: 2px 1rem;
        border: 1px solid $col-zinc-400;
        text-wrap: nowrap;

        &:not(:nth-last-child(2)) {
            border-bottom: none;
        }

        &:hover {
            color: $col-text-2;
            background-color: $col-blue-200;
            border-color: $col-text-2;

            &+.tabs-mb__item {
                border-top: 1px solid $col-text-2;
            }
        }

        &--active {
            color: $col-white;
            background-color: $col-zinc-400;
            border-color: $col-zinc-400;
        }
    }
}

.tabs-dt {
    flex-direction: row;

    &__empty {
        display: block;
        border-bottom: 1px solid $col-zinc-400 !important;
        flex: 1;
    }

    &__item {
        cursor: pointer;
        font-family: $font-action;
        color: $col-zinc-400;
        font-weight: 400;
        padding: 2px 0.5rem;
        min-width: 11rem;
        border: 1px solid $col-zinc-400;
        display: flex;
        justify-content: center;
        font-size: .9rem;
        text-wrap: nowrap;

        &:not(:nth-last-child(2)) {
            border-bottom: none;
        }

        &:not(:nth-last-child(2)) {
            border-bottom: 1px solid $col-zinc-400;
            border-right: none;
        }

        &:hover {
            color: $col-text-2;
            background-color: $col-blue-200;
            border-color: $col-text-2;
            border-bottom: 1px solid $col-text-2 !important;

            &+.tabs-dt__item {
                border-top: 1px solid $col-zinc-400;
                border-left: 1px solid $col-text-2;
            }
        }

        &--active {
            color: $col-white;
            background-color: $col-zinc-400;
            border-color: $col-zinc-400;
        }
    }
}
</style>
