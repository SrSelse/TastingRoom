<template>
  <div class="min-h-screen bg-gray-900 text-white">
    <!-- Main Content -->
    <main class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-12">
      <div class="flex flex-col sm:flex-row sm:justify-between sm:items-center mb-8">
        <h1 class="text-3xl font-bold mb-4 sm:mb-0">Your Tasting Rooms</h1>

        <div class="flex flex-col sm:flex-row items-center space-y-4 sm:space-y-0 sm:space-x-4">
          <!-- Invitation Code Section -->
          <div class="flex items-center space-x-2">
            <input
              v-model="invitationCode"
              type="text"
              placeholder="Invitation code"
              class="px-3 py-2 bg-gray-700 border border-gray-600 rounded-md text-white focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 text-sm w-48"
            />
            <button
              @click="joinRoomWithCode"
              :disabled="isJoining || !invitationCode"
              class="px-3 py-2 bg-gray-700 hover:bg-gray-600 text-white rounded-md transition-colors text-sm flex items-center disabled:opacity-50"
            >
              <svg v-if="isJoining" class="animate-spin -ml-1 mr-2 h-4 w-4 text-white" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
                <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
              </svg>
              <span v-if="!isJoining">Join Room</span>
              <span v-else>Joining...</span>
            </button>
          </div>

          <!-- Create Room Button -->
          <router-link
            to="/rooms/new"
            class="inline-block bg-indigo-600 hover:bg-indigo-700 text-white font-medium py-2 px-4 rounded-md transition-colors"
          >
            Create New Room
          </router-link>
        </div>
      </div>

      <!-- Rooms List -->
      <div v-if="loading" class="text-center py-12">
        <p class="text-gray-300">Loading rooms...</p>
      </div>
      <div v-else-if="error" class="text-center py-12">
        <p class="text-red-400">Error loading rooms: {{ error }}</p>
      </div>
      <div v-else>
        <div v-if="rooms.length === 0" class="text-center py-12">
          <p class="text-gray-300">No rooms found. Create your first tasting room!</p>
        </div>
        <ul v-else class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
          <li
            v-for="room in rooms"
            :key="room.id"
            class="bg-gray-800 rounded-lg p-6 hover:shadow-lg transition-shadow"
          >
            <h2 class="text-xl font-semibold mb-2">{{ room.name }}</h2>
            <p class="text-gray-300 mb-4">{{ room.description }}</p>
            <div class="flex justify-between items-center">
              <span class="text-sm text-gray-400">{{ room.createdAt }}</span>
              <router-link
                :to="`/rooms/${room.id}`"
                class="text-indigo-400 hover:text-indigo-300 text-sm font-medium"
              >
                Enter Room
              </router-link>
            </div>
          </li>
        </ul>
      </div>
    </main>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue';
import { useRouter } from 'vue-router';

const router = useRouter();
const rooms = ref([]);
const loading = ref(true);
const error = ref(null);
const invitationCode = ref('');
const isJoining = ref(false);

// Simulate fetching rooms from an API
const fetchRooms = async () => {
  try {
    loading.value = true;
    const response = await fetch(`${import.meta.env.VITE_API_URL}/api/rooms`, {
      headers: {
        'Authorization': `Bearer ${localStorage.getItem('token')}`,
        'Content-Type': 'application/json'
      }
    });
    if (!response.ok) {
      throw new Error("Internal server error");
    }
    const data = await response.json();
    if (data.length === 0) {
      rooms.value = [];
      return;
    }
    rooms.value = data.map((room) => {
      const date = new Date(room.createdAt).toLocaleString('sv-SE').replace('T', ' ').substring(0, 16);
      return {
        ...room,
        createdAt: date
      }
    });
  } catch (err) {
    error.value = err.message;
  } finally {
    loading.value = false;
  }
};

// Join room with invitation code
const joinRoomWithCode = async () => {
  if (!invitationCode.value) {
    alert('Please enter an invitation code');
    return;
  }

  isJoining.value = true;

  try {
    const response = await fetch(`${import.meta.env.VITE_API_URL}/api/room/join`, {
      method: 'POST',
      headers: {
        'Authorization': `Bearer ${localStorage.getItem('token')}`,
        'Content-Type': 'application/json'
      },
      body: JSON.stringify({
        code: invitationCode.value
      })
    });

    if (!response.ok) {
      const errorData = await response.json();
      throw new Error(errorData.message || 'Failed to join room');
    }

    const room = await response.json();
    // Clear the invitation code field
    invitationCode.value = '';
    // Redirect to the joined room
    router.push(`/rooms/${room.id}`);
  } catch (err) {
    error.value = err.message;
    alert(`Error: ${err.message}`);
  } finally {
    isJoining.value = false;
  }
};

onMounted(fetchRooms);
</script>
