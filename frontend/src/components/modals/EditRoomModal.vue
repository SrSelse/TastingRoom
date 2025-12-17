<template>
  <div
    class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center p-4 z-50"
    @click.self="emit('close')"
  >
    <!-- Add Beer Modal Content -->
    <div class="bg-gray-800 rounded-lg p-6 max-w-md w-full">
      <h2 class="text-xl font-bold mb-4">Edit Room</h2>
      <div class="space-y-4">
        <div>
          <label for="roomName" class="block text-sm font-medium text-gray-300">Room Name</label>
          <input
            id="beerName"
            v-model="localRoom.name"
            type="text"
            required
            class="mt-1 block w-full px-3 py-2 bg-gray-700 border border-gray-600 rounded-md text-white focus:outline-none focus:ring-indigo-500 focus:border-indigo-500"
          />
        </div>
        <div>
          <label for="description" class="block text-sm font-medium text-gray-300 mb-1">Description</label>
          <textarea
            id="description"
            v-model="localRoom.description"
            rows="4"
            class="w-full px-3 py-2 bg-gray-700 border border-gray-600 rounded-md text-white focus:outline-none focus:ring-indigo-500 focus:border-indigo-500"
            placeholder="Describe your tasting event..."
          ></textarea>
        </div>
        <div>
          <label for="planned_date" class="block text-sm font-medium text-gray-300 mb-1">Planned Date</label>
          <input
            id="planned_date"
            v-model="localRoom.plannedDate"
            type="datetime-local"
            class="w-full px-3 py-2 bg-gray-700 border border-gray-600 rounded-md text-white focus:outline-none focus:ring-indigo-500 focus:border-indigo-500"
          />
        </div>
        <div class="flex justify-end space-x-4">
          <button
            type="button"
            @click="emit('close')"
            class="px-4 py-2 text-gray-300 hover:text-white"
          >
            Cancel
          </button>
          <button
            @click="editRoom"
            class="px-4 py-2 bg-indigo-600 hover:bg-indigo-700 rounded-md text-white"
          >
            Save
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue';
const props = defineProps({
  data: {
    type: Object,
    required: true
  }
});
const emit = defineEmits(['close', 'saved']);
const localRoom = ref(props.data);

const editRoom = async () => {
  try {
    const response = await fetch(`${import.meta.env.VITE_API_URL}/api/room/${props.data.id}/edit`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${localStorage.getItem('token')}`
      },
      body: JSON.stringify({
        name: localRoom.value.name,
        description: localRoom.value.description,
        plannedDate: localRoom.value.plannedDate,
      }),
    });
    if (!response.ok) throw new Error('Failed to add');
    const newBeer = await response.json();
    emit('saved');
    emit('close');
  } catch (err) {
    alert(`Error: ${err.message}`);
  }
};
</script>
