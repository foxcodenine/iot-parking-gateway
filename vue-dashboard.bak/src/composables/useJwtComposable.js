import { useAuthStore } from '@/stores/authStore';
import { computed } from 'vue';

export function useJwtComposable() {
    const authStore = useAuthStore();

    const token = computed(() => authStore.getJwt);

    function parseJwt() {

        if (!token.value) { return false; }

        const base64Url = token.value.split('.')[1]; // get the payload part
        const base64 = base64Url.replace(/-/g, '+').replace(/_/g, '/'); // Convert Base64URL to Base64
        const jsonPayload = decodeURIComponent(atob(base64).split('').map(function(c) {
            return '%' + ('00' + c.charCodeAt(0).toString(16)).slice(-2);
        }).join(''));
    
        return JSON.parse(jsonPayload);
    }
    
    function checkJwtExpiration() {

        if (!token.value) { return false; }

        const payload = parseJwt();
        const currentTime = Math.floor(Date.now() / 1000); // current time in seconds since epoch
    
        if (payload?.exp && payload.exp < currentTime) {
            return false;
        } else {
            return true;
        }
    }

    return {
        parseJwt,
        checkJwtExpiration
    };
}
