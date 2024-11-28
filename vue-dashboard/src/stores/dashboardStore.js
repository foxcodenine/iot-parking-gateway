import { ref, computed } from 'vue'
import { defineStore } from 'pinia'

// ---------------------------------------------------------------------


export const useDashboardStore = defineStore("dashboardStore", () => {


    // ---- State ------------------------------------------------------

    const isUserMenuOpen = ref(false);
    const isLoading = ref(false)

    // ---- Getters ----------------------------------------------------

    const getIsUserMenuOpen = computed(()=>{
        return isUserMenuOpen.value;
    });

    const getIsLoading = computed(()=>{
        return isLoading.value;
    })

    // ---- Actions ----------------------------------------------------

    function toggleUserMenu() {
        isUserMenuOpen.value = !isUserMenuOpen.value
    }

    function updateUserMenu(val) {
        isUserMenuOpen.value = val;
    }

    function setIsLoading(val) {
        isLoading.value = val;
    }


    // - Expose --------------------------------------------------------

    return {
        getIsUserMenuOpen,
        toggleUserMenu,
        updateUserMenu,        
        getIsLoading,
        setIsLoading,
    }
});