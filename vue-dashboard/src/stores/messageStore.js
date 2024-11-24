import { ref, computed } from 'vue'
import { defineStore } from 'pinia'

// ---------------------------------------------------------------------


export const useMessageStore = defineStore("messageStore", () => {


    // - State ---------------------------------------------------------


    const flashMessages =ref(["Email and password are required."]);
    const flashClass = ref("flash-message--red");

    const persistFlashMessage = ref(0);


    // - Getters -------------------------------------------------------



    const getFlashMessages = computed(()=>{
        return flashMessages.value;
    });

    const getFlashClass = computed(()=>{
        return flashClass.value;
    });

    const getPersistFlashMessage = computed(()=>{
        return persistFlashMessage.value;
    });


    // - Actions -------------------------------------------------------


    function setFlashMessages(msg, msgClass=false) {       
        flashMessages.value = msg;
        if (msgClass) {
            flashClass.value = msgClass;
        }
    }
    
    function setFlashClass(msgClass) {
         flashClass.value = msgClass;
    }

    function clearFlashMessage() {
           flashMessages.value = [];
        flashClass.value = '';
    }

    function setPersistFlashMessage(val) {
        persistFlashMessage.value = val
    }

    function decreasePersistFlashMessage() {
        
        if (persistFlashMessage.value != 0) {
            persistFlashMessage.value -= 1;
        }
        console.log(persistFlashMessage.value)
    }
    

    // - Expose --------------------------------------------------------

    return {
        
        getFlashMessages,
        setFlashMessages,
        
        setFlashClass,
        getFlashClass,

        clearFlashMessage,

        getPersistFlashMessage,
        setPersistFlashMessage,
        decreasePersistFlashMessage,
    }
});