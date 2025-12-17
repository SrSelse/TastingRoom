<template>
    <div
      class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center p-4 z-50"
      @click.self="emit('close')"
    >
      <!-- Add Beer Modal Content -->
      <div class="bg-gray-800 rounded-lg p-6 max-w-md w-full">
        <h2 class="text-xl font-bold mb-4">Add New Beverage</h2>
        <div class="space-y-4">
          <div>
            <label for="beerName" class="block text-sm font-medium text-gray-300">Beverage Name</label>
            <input
              id="beerName"
              v-model="newBeerName"
              type="text"
              required
              class="mt-1 block w-full px-3 py-2 bg-gray-700 border border-gray-600 rounded-md text-white focus:outline-none focus:ring-indigo-500 focus:border-indigo-500"
            />
          </div>
          <div>
            <label for="beerStyle" class="block text-sm font-medium text-gray-300">Style</label>
            <input
              id="beerStyle"
              v-model="newBeerStyle"
              type="text"
              required
              class="mt-1 block w-full px-3 py-2 bg-gray-700 border border-gray-600 rounded-md text-white focus:outline-none focus:ring-indigo-500 focus:border-indigo-500"
            />
          </div>
          <div>
            <label for="beerPictureUrl" class="block text-sm font-medium text-gray-300">Picture URL</label>
            <input
              id="beerPictureUrl"
              v-model="newBeerUrl"
              type="text"
              class="mt-1 block w-full px-3 py-2 bg-gray-700 border border-gray-600 rounded-md text-white focus:outline-none focus:ring-indigo-500 focus:border-indigo-500"
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
              @click="addBeer"
              class="px-4 py-2 bg-indigo-600 hover:bg-indigo-700 rounded-md text-white"
            >
              Add
            </button>
          </div>
        </div>
      </div>
    </div>
</template>

<script setup>
import { ref } from 'vue';
const props = defineProps({
  roomId: {
    type: String,
    required: true
  }
});
const emit = defineEmits(['close', 'added']);
const newBeerName = ref('');
const newBeerStyle = ref('');
const newBeerUrl = ref('');

const addBeer = async () => {
  try {
    const response = await fetch(`${import.meta.env.VITE_API_URL}/api/room/${props.roomId}/beers/new`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${localStorage.getItem('token')}`
      },
      body: JSON.stringify({
        name: newBeerName.value,
        style: newBeerStyle.value,
        pictureUrl: newBeerUrl.value,
        roomId: Number(props.roomId),
      }),
    });
    if (!response.ok) throw new Error('Failed to add');
    const newBeer = await response.json();
    emit('added', newBeer);
    emit('close')
    newBeerName.value = '';
    newBeerStyle.value = '';
    newBeerUrl.value = '';
  } catch (err) {
    alert(`Error: ${err.message}`);
  }
};
</script>
