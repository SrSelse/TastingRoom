<template>
  <div
    class="fixed bottom-4 right-4 px-4 py-2 bg-green-800 text-white rounded-md shadow-lg flex items-center space-x-2 z-50"
  >
    <span class="text-green-200">✓</span>
    <span>{{ text }}</span>
  </div>
</template>

<script setup>
const props = defineProps({
  text: {
    type: String,
    required: true
  }
});
</script>
