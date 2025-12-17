<template>
  <div class="min-h-screen bg-gray-900 text-white py-12 px-4 sm:px-6 lg:px-8">
    <div class="max-w-2xl mx-auto">
      <div class="mb-8">
        <h1 class="text-3xl font-bold">Create New Tasting Room</h1>
        <p class="text-gray-300 mt-2">Set up a new room for your tasting event</p>
      </div>

      <div class="space-y-6 bg-gray-800 p-6 rounded-lg">
        <!-- Name Field -->
        <div>
          <label for="name" class="block text-sm font-medium text-gray-300 mb-1">Room Name</label>
          <input
            id="name"
            v-model="roomName"
            type="text"
            required
            class="w-full px-3 py-2 bg-gray-700 border border-gray-600 rounded-md text-white focus:outline-none focus:ring-indigo-500 focus:border-indigo-500"
            placeholder="e.g., Craft Beer Night"
          />
        </div>

        <!-- Description Field -->
        <div>
          <label for="description" class="block text-sm font-medium text-gray-300 mb-1">Description</label>
          <textarea
            id="description"
            v-model="roomDescription"
            rows="4"
            class="w-full px-3 py-2 bg-gray-700 border border-gray-600 rounded-md text-white focus:outline-none focus:ring-indigo-500 focus:border-indigo-500"
            placeholder="Describe your tasting event..."
          ></textarea>
        </div>

        <!-- Planned Date Field -->
        <div>
          <label for="planned_date" class="block text-sm font-medium text-gray-300 mb-1">Planned Date</label>
          <input
            id="planned_date"
            v-model="plannedDate"
            type="datetime-local"
            class="w-full px-3 py-2 bg-gray-700 border border-gray-600 rounded-md text-white focus:outline-none focus:ring-indigo-500 focus:border-indigo-500"
          />
        </div>

        <!-- Create Button -->
        <div class="flex justify-end">
          <button
            @click="createRoom"
            class="px-6 py-2 bg-indigo-600 hover:bg-indigo-700 rounded-md text-white font-medium transition-colors"
          >
            Create Room
          </button>
        </div>
      </div>

      <!-- Success Toast -->
      <div
        v-if="showSuccessToast"
        class="fixed bottom-4 right-4 px-4 py-2 bg-green-800 text-white rounded-md shadow-lg flex items-center space-x-2 z-50"
      >
        <span class="text-green-200">✓</span>
        <span>Tasting room created successfully!</span>
      </div>

      <!-- Error Toast -->
      <div
        v-if="showErrorToast"
        class="fixed bottom-4 right-4 px-4 py-2 bg-red-800 text-white rounded-md shadow-lg flex items-center space-x-2 z-50"
      >
        <span class="text-red-200">✗</span>
        <span>{{ errorMessage }}</span>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue';
import { useRouter } from 'vue-router';

const router = useRouter();

// Reactive data properties
const roomName = ref('');
const roomDescription = ref('');
const plannedDate = ref('');

// Toast notifications
const showSuccessToast = ref(false);
const showErrorToast = ref(false);
const errorMessage = ref('');

// AJAX function to create a new room
const createRoom = async () => {
  try {
    const response = await fetch(`${import.meta.env.VITE_API_URL}/api/room/create`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${localStorage.getItem('token')}`
      },
      body: JSON.stringify({
        name: roomName.value,
        description: roomDescription.value,
        plannedDate: plannedDate.value
      })
    });

    if (!response.ok) {
      const errorData = await response.json();
      throw new Error(errorData.message || 'Failed to create room');
    }

    // Show success toast
    showSuccessToast.value = true;

    // Get the newly created room data
    const newRoom = await response.json();

    // Redirect to the new room after a short delay
    setTimeout(() => {
      showSuccessToast.value = false;
      router.push(`/rooms/${newRoom.id}`);
    }, 2000);

  } catch (err) {
    errorMessage.value = err.message;
    showErrorToast.value = true;
    setTimeout(() => {
      showErrorToast.value = false;
    }, 3000);
  }
};
</script>
