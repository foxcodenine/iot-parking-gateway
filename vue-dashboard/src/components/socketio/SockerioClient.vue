<template>
  <div style="display: none;">
      <h1>Socket.io with Vue</h1>
      <button @click="sendData">Send Data</button>
  </div>
</template>

<!-- --------------------------------------------------------------- -->

<script setup>
import { ref, onMounted, onBeforeUnmount } from 'vue';
import socketIOClient from 'socket.io-client';
import { useAppStore } from '@/stores/appStore';

const socket = ref(null);

// Function to send data to the server
const sendData = () => {
  if (socket.value) {
      socket.value.emit('update', 'Hello from Vue!');
  }
};

onMounted(() => {
  localStorage.setItem('debug', ''); // Disable debug logs

  socket.value = socketIOClient(useAppStore().getAppUrl, {
    //   path: '/socket.io/',
      transports: ['websocket'],
      reconnection: true,
  });

  // Handle connection
  socket.value.on('connect', () => {
      console.log('Connected to server');
  });

  // Handle custom events
  socket.value.on('received', (data) => {
      console.log(data);
  });

  socket.value.on('update', (data) => {
      console.log(data);
  });

  // Handle disconnection
  socket.value.on('disconnect', () => {
      console.log('Disconnected from server');
  });

  // Handle connection errors
  socket.value.on('connect_error', (error) => {
      console.error('Connection Error:', error);
  });
});

onBeforeUnmount(() => {
  // Disconnect socket when component is destroyed
  if (socket.value) {
      socket.value.disconnect();
  }
});
</script>

<!-- --------------------------------------------------------------- -->

<style scoped>
/* Add your component-specific styles here */
</style>
