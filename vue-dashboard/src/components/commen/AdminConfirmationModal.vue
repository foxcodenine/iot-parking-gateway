<template>
    <section class="modal">
        <div class="modal__background"></div>
        <div class="modal__content">
            <div class="modal__title">
                Admin Password Confirmation Required
            </div>
            <div class="modal__body">
                <p class="modal__text">
                    For security purposes, please re-enter your admin password to confirm this action.
                </p>
                <input type="password" class="modal__input mt-4" placeholder="Enter your password" v-model.trim="adminPassword">
            </div>
            <div class="modal__footer">
                <div class="bbtn bbtn--zinc-lt" @click="cancel">Cancel</div>
                <div class="bbtn bbtn--emerald" @click="confirm">Confirm</div>
            </div>
        </div>
    </section>
</template>

<!-- --------------------------------------------------------------- -->

<script setup>
import { onMounted, ref } from 'vue';

const emit = defineEmits(['emitCancel', 'emitConfirm'])

const adminPassword = ref("")

function cancel() {
    emit('emitCancel');
}

function confirm() {
    if (adminPassword.value.length < 6) { return }
    emit('emitConfirm', {adminPassword: adminPassword.value});
}

onMounted(()=>{
    setTimeout(()=>{
        document.querySelector('.modal__content').classList.add('modal__show');
    }, 50);

    setTimeout(()=>{
        adminPassword.value = '';
        emit('emitCancel');
    }, 20000)
})

</script>

<!-- --------------------------------------------------------------- -->

<style lang="scss" scoped>

.modal {
    position: fixed;
    top: 0;
    left: 0;
    right: 0;
    bottom: 0;
    z-index: 2000;

    display:  grid;
    grid-template-columns: 1fr;
    grid-template-rows: 1fr;
    justify-items: center;
    align-items: center;

    &__background {
        background-color: $col-zinc-900;
        width: 100%;
        height: 100%;
        opacity: .3;
        grid-column: 1 / -1;
        grid-row: 1 / -1;
        z-index: 1;
    }

    &__content {
        min-width: 10rem;
        min-height: 5rem;
        margin: 1rem;
        background-color: $col-zinc-50;
        grid-column: 1 / -1;
        grid-row: 1 / -1;
        z-index: 10;
        opacity: 1;
        border: 1px solid $col-zinc-400;
        border-radius: 5px;        
        font-family: $font-primary;
        opacity: 0;
        transition: all .4s ease;
    }

    &__show {
        opacity: 1;
    }

    &__title {
        padding: 1rem;
        border-bottom: 1px solid $col-zinc-200;
        font-size: 1.2rem;
    }
    &__body {
        padding: 1rem;
        border-bottom: 1px solid $col-zinc-200;
    }

    &__input {
        cursor: text;
        width: 100%;
        height: 100%;    
        padding: 0.3rem 0.5rem;
        background-color: rgba(white, .5); 
        color: $col-text-1;
        font-family: $font-action;
        font-size: 1rem;
        font-weight: 300;
        border: 1px solid $col-slate-300;
        border-radius: $border-radius;

        &:focus {
            border: 1px solid $col-blue-500;
        }
    }

    &__footer {
        padding: 1rem;
        display: flex;
        justify-content: space-between;
    }
}
.bbtn {
    width: 8rem;
}

</style>