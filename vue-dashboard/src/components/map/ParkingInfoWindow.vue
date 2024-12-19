<template>
    <InfoWindow v-if="activeWindow==device.device_id" :options="{ position: { lat: device.latitude, lng: device.longitude }, anchorPoint: 'CENTER' }">
        <div class="info-window">
            <div class="info-window__header">
                <h3>{{ device.name }}</h3>
           
            </div>
            <div class="info-window__content mt-2">
                <span>id:</span>
                <p>{{ device.device_id }}</p>
            </div>
            <div class="info-window__content">
                <span>nw:</span>
                <p>{{ device.network_type }}</p>
            </div>
            <div class="info-window__content">
                <span>fw:</span>
                <p>{{ device.firmware_version }}</p>
            </div>
            <div class="info-window-footer">
                <button class="action-btn" @click="onAction">More Info</button>
            </div>
        </div>
    </InfoWindow>
</template>

<!-- --------------------------------------------------------------- -->

<script setup>
import { ref } from 'vue';
import { InfoWindow } from 'vue3-google-map';

const props = defineProps({

    device: {
        type: Object,
        required: true,  
    },
    activeWindow: {
        type: [Number, null, String],
        required: true,  
    }
});
</script>

<!-- --------------------------------------------------------------- -->

<style lang="scss" scoped>
.info-window {

    border-radius: 3px !important;
    box-shadow: 0 4px 10px rgba(0, 0, 0, 0.1);
    overflow: hidden;
    min-width: 200px;
    


    &__header {
        padding: 10px 10px;
        display: flex;
        justify-content: space-between;
        align-items: center;
        justify-content: center;
        border-bottom: 1px solid $col-zinc-700;

        h3 {
            margin: 0;
            font-size: 1rem;
            font-weight: bold;         
        }

        .close-btn {
            background: none;
            border: none;
            color: $col-zinc-300;
            font-size: 1.4rem;
            font-weight: 100;
            cursor: pointer;
            &:hover {
                color: $col-red-400;
            }
        }
    }

    &__content {
        padding: 3px 10px;
        font-size: 1rem;
        display: grid;      
        grid-template-columns: 2rem 1fr; 
        align-items: center;
        width: max-content !important;

        span {
            font-family: $font-display;;
        }
    }

    &__icon {
        fill: currentColor;
        width: 1rem;
        height: 1rem;
    }

    &-footer {
        padding: 5px 10px;
        text-align: right;

        .action-btn {
            background-color: #007bff;
            color: white;
            border: none;
            padding: 5px 10px;
            border-radius: 4px;
            font-size: 0.9rem;
            cursor: pointer;

            &:hover {
                background-color: #0056b3;
            }
        }
    }
}
</style>