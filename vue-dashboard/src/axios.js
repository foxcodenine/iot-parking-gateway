import axios from "axios";
import { useMessageStore } from "./stores/messageStore";
import { useAuthStore } from "./stores/authStore";

// ---------------------------------------------------------------------

// Create an axios instance
const apiClient = axios.create({

    // baseURL: import.meta.env.VITE_API_BASE_URL,

    headers: {
        'Content-Type': 'application/json' // Ensure all requests use JSON
    },

    timeout: 15000 // Timeout if the request takes longer than 15 seconds
});

// ---------------------------------------------------------------------

apiClient.interceptors.request.use(function (config) {

    const authStore = useAuthStore();

    const isAuthenticated = authStore.isAuthenticated;

    if (isAuthenticated) {

        const token = authStore.getJwt;
        config.headers.Authorization = `Bearer ${token}`;
    }

    return config;

}, function (error) {

    console.log('! apiClient.interceptors.request !');

    // Do something with request error
    return Promise.reject(error);
});

// ---------------------------------------------------------------------

apiClient.interceptors.response.use(

    response => {

        let data = response?.data;

        if (data?.messages && !data?.actions?.includes('hideMessage')) {
            useMessageStore().setFlashMessages(response.data.messages);
            useMessageStore().setFlashClass('flash-message--blue');
        }
        
        return response;
    },

    error => {
        console.error('! axios.interceptors.response !');

        if (error.status == 401 && useAuthStore().isAuthenticated) {
            window.location.assign('/logout');
            return Promise.reject(error);

        } else {
            return Promise.reject(error);
        }
       
    }
);

// ---------------------------------------------------------------------

export default apiClient;