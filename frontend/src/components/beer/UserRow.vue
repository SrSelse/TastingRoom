<template>
  <div class="flex flex-col p-4 bg-gray-700 rounded-lg">
    <div class="flex justify-between items-center ">
      <div class="flex items-center space-x-4">
        <div
          class="h-10 w-10 rounded-full flex items-center justify-center text-gray-300"
          :class="user.rating ? 'ring-2 ring-green-500 bg-gray-800' : 'bg-gray-600'"
        >
          {{ user.name.charAt(0).toUpperCase() }}
        </div>
        <div>
          <p class="font-medium">{{ user.name }}</p>
        </div>
      </div>
      <div v-if="showRating" class="text-gray-300">
        <span>
          {{ user.rating }}
        </span>

      </div>
      <div v-else class="text-gray-500">
        Hidden
      </div>
    </div>
    <div
      v-if="showRating && user.note"
      class="flex w-full"
    >
      <div class="flex-col w-full">
        <button
          @click="showNotes = !showNotes"
          class="flex items-center text-sm font-medium text-gray-300 mb-2 focus:outline-none mt-2"
        >
        <span>{{ showNotes ? '▼' : '▶' }} Show Notes</span>
        </button>
        <div v-if="showNotes" class="mt-2">
          <div
            v-text="user.note"
            class="w-full px-3 py-2 bg-gray-700 border border-gray-600 rounded-md text-white focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 whitespace-pre-line"
          >
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue';
const props = defineProps({
  user: Object,
  showRating: Boolean,
});

const showNotes = ref(false);
</script>
