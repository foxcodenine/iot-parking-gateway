<template>
    <div id="create-asset-modal" class="my-modal fadeInTransition" :class="{'my-modal--on': modalIsOpen}">
        <div id="create-asset-modal-background" class="my-modal__background" @click="closeModal"></div>
        <div class="my-modal__wrapper">
            <div id="create-asset-modal-map" class="my-modal__map">
                <GoogleMap v-if="apiKey" :api-key="apiKey" :map-id="getMapId" class="map" :zoom="mapZoom" :center="mapCenter" @click="updatedMarkerLocation">
                    <Marker :options="{ position: markerPosition }"></Marker>
                </GoogleMap>
            </div>
            <div @click="closeModal" id="create-asset-modal-close-btn" class="my-modal__close">&times;</div>
        </div>
    </div>
</template>

<script setup>
import { ref, computed } from 'vue';
import { GoogleMap, Marker, CustomMarker } from "vue3-google-map";
import { defineProps, defineEmits } from 'vue';
import { useAppStore } from '@/stores/appStore';
import { storeToRefs } from 'pinia';
import { onMounted } from 'vue';

// - Store -------------------------------------------------------------

const appStore = useAppStore();


// - Props -------------------------------------------------------------

const props = defineProps({
    modalIsOpen: {
        type: Boolean,
        default: false
    },
    markerPosition: {
        type: Object,
        default: {lat:0, lng:0}
    }
});

// - Data --------------------------------------------------------------

const apiKey = ref(null);
const mapZoom = ref(15);

// - Computed ----------------------------------------------------------

const mapCenter = computed(() => props.markerPosition);
const getMapId = computed(() => { return appStore.getAppSettings.google_map_id ?? ""; });

// - Hooks -------------------------------------------------------------

onMounted(async () => {
    apiKey.value = await appStore.getGoogleApiKey();
});


const emits = defineEmits(['emitCloseModal', 'emitMarkerPosition']);



function closeModal() {
    emits('emitCloseModal');
}

function updatedMarkerLocation(e) {
    emits('emitMarkerPosition', e.latLng.toJSON());
}
</script>

<!-- --------------------------------------------------------------- -->

<style lang="scss" scoped>
/* Your existing styles can stay the same */
.map {
    height: 100%;
    width: 100%;
}
.my-modal {
    z-index: 500;
    position: fixed;
    top: 0; right: 0; bottom: 0; left: 0;
    grid-template-columns: 1fr;
    grid-template-rows: 1fr;
    justify-items: center;
    align-items: center;
    display: none;
    transition: all .2s ease-in;
    &--on { display: grid; }
    &__background {
        width: 100%;
        height: 100%;   
        grid-column: 1 / 2;
        grid-row: 1 / 2;
        background-color: $col-zinc-900;
        opacity: .3;
    }
    &__wrapper {
        grid-column: 1 / 2;
        grid-row: 1 / 2;
        z-index: 1000;
        padding: 1rem;
        padding-bottom: 4rem;
        width: 100%;
        height: 100%;
        max-width: 80rem;
        max-height: 46rem;     
        position: relative;
        @include respond(600) {
            max-height: 100%;
        }
    }
    &__map {
        border: 1px solid $col-zinc-50;
        border-radius: $border-radius;
        background-color: #fff;
        height: 100%;
        .gmnoprint { display: none !important; }
        .gm-fullscreen-control { display: none !important; }
    }
    &__close {
        cursor: pointer;
        z-index: 1000001;
        position: absolute;
        top: 0.1rem;
        right: 1.4rem;
        color: $col-white;
        font-size: 2.4rem;
        font-weight: 200;
       text-shadow: 
            0 0 2px hsl(216, 97%, 70%), 
            0 0 5px hsl(216, 97%, 70%);
        transition: all .1s ease-in;
        &:hover {
            color: $col-red-400;
            text-shadow: 
            0 0 2px rgba($col-red-400, .5), 
            0 0 5px rgba($col-red-400, .5);
        }
    }
}
</style>
