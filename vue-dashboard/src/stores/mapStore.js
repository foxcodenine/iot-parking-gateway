import { ref, computed } from 'vue'
import { defineStore } from 'pinia'
import { useAppStore } from './appStore';

// ---------------------------------------------------------------------


export const useMapStore = defineStore("mapStore", () => {


    // - State ---------------------------------------------------------
    const appStore = useAppStore();
    const mapZoom = ref(17);
    const mapCenter = ref({
        lat: Number(appStore.getAppSettings.default_latitude),
        lng: Number(appStore.getAppSettings.default_longitude)
    });


    // - Getters -------------------------------------------------------

    const getMapCenter = computed(()=>{

        return mapCenter.value;
    });

    const getMapZoom = computed(()=>{


        return mapZoom.value;
    })




    // - Actions -------------------------------------------------------

    function setMapCenter(payload) {
        mapCenter.value = payload;
    }

    function setMapZoom(payload) {
        mapZoom.value = payload;
    }
    

    // - Expose --------------------------------------------------------

    return {
        getMapCenter,
        getMapZoom,
        setMapCenter,
        setMapZoom,
    }
});