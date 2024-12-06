<script setup>
import io from 'socket.io-client';

// Enable debug logs
localStorage.debug = '*';

const socket = io('http://localhost:6060/', {
    transports: ['websocket', 'polling']
});

socket.on('connect', () => {
  console.log('[Client] Connected to server:', socket.id);
});

socket.on('update', (data) => {
  console.log('[Client] Received update:', data);
});



socket.on('connect_error', (err) => {
  console.error('[Client] Connection error:', err.message);
});
</script>
