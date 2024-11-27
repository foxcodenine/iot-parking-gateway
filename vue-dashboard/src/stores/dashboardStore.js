import { ref, computed } from 'vue'
import { defineStore } from 'pinia'

// ---------------------------------------------------------------------


export const useDashboardStore = defineStore("dashboardStore", () => {


    // ---- State ------------------------------------------------------

    const isUserMenuOpen = ref(false);
    const isFetching = ref(false)

    // ---- Getters ----------------------------------------------------

    const getIsUserMenuOpen = computed(()=>{
        return isUserMenuOpen.value;
    });

    const getIsFetching = computed(()=>{
        return isFetching.value;
    })

    // ---- Actions ----------------------------------------------------

    function toggleUserMenu() {
        isUserMenuOpen.value = !isUserMenuOpen.value
    }

    function updateUserMenu(val) {
        isUserMenuOpen.value = val;
    }

    function setIsFetching(val) {
        isFetching.value = val;
    }


    // - Expose --------------------------------------------------------

    return {
        getIsUserMenuOpen,
        toggleUserMenu,
        updateUserMenu,
        
        getIsFetching,
        setIsFetching,
    }
});