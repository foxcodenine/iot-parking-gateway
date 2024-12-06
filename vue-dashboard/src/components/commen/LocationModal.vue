<template>
    <div id="create-asset-modal"  class="my-modal fadeInTransition" :class="{'my-modal--on': modalIsOpen}">
        <div id="create-asset-modal-background" class="my-modal__background" @click="closeModal"></div>

        <div class="my-modal__wrapper">
            <div id="create-asset-modal-map" class="my-modal__map">
                <GoogleMap :api-key="mapApiKey" class="map" :zoom="mapZoom" :center="mapCenter" @click="updatedMarkerLocation" >
                    <Marker :options="{ position: markerPosition }"></Marker>
                </GoogleMap>
            </div>
            <div @click="closeModal" id="create-asset-modal-close-btn" class="my-modal__close">&times;</div>
        </div>
        
    </div>
</template>

<!-- --------------------------------------------------------------- -->

<script>
import { GoogleMap, Marker, CustomMarker } from "vue3-google-map";
export default {
    setup() {       
        return {}
    },

    components: { GoogleMap, Marker, CustomMarker },

    emits: ['emitCloseModal', 'emitMarkerPostion'],

    props: {
        modalIsOpen: {
            type: Boolean,
            default: false
        },
        markerPosition: {
            type: Object,
            default: { lat: 35.94, lng: 14.4 }
        }
    },

    data() {
        return {
            mapApiKey: 'AIzaSyAIs3GtjOuWwwkDDSYfhe3E2q0qUrHPNFi',
            mapZoom: 11.2,
            mapCenter: { lat: 35.94, lng: 14.4 },
        }
    },
    
    methods: {
        closeModal() {
            console.log(123)
            this.$emit('emitCloseModal');
        },
        updatedMarkerLocation(e) {            
            this.$emit('emitMarkerPostion', e.latLng.toJSON());
        },
    }
}
</script>

<!-- --------------------------------------------------------------- -->

<style lang="scss" scoped>


.map {
    height: 100%;
    width: 100%;
}
.my-modal {
    z-index: 500;
    position: fixed;
    top: 0; right: 0; bottom: 0;  left: 0;

    grid-template-columns: 1fr;
    grid-template-rows: 1fr;

    justify-items: center;
    align-items: center;

    display: none;
    // opacity: 0;

    transition: all .2s ease-in;

    &--on { display: grid; }

    &__background {
        width: 100%;
        height: 100%;
        background-color:  #000;
        opacity: .5;

        grid-column: 1 / 2;
        grid-row: 1 / 2;
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

        border: 1px solid $col-white;
        border-radius: $border-radius;
        background-color:  #fff;

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
            0 0 5px hsl(216, 97%, 70%), 
            0 0 8px hsl(216, 97%, 70%), 
            0 0 12px hsl(216, 97%, 70%),
            0 0 16px hsl(216, 97%, 70%);

        transition: all .1s ease-in;

        &:hover {
            color: $col-red-600;

            text-shadow: 
            0 0 2px rgba($col-red-600, .5), 
            0 0 5px rgba($col-red-600, .5), 
            0 0 8px rgba($col-red-600, .5), 
            0 0 12px rgba($col-red-600, .5),
            0 0 16px rgba($col-red-600, .5);
        }
    }
}
</style>