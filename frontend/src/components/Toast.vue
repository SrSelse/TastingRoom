<template>
  <div
    class="fixed bottom-4 right-4 px-4 py-2 text-white rounded-md shadow-lg flex items-center space-x-2 z-50"
    :class="backgroundClass"
  >
    <span 
      v-if="toastType === 'success'"
      class="text-green-200"
    >
      ✓
    </span>
    <span>{{ text }}</span>
  </div>
</template>

<script setup>
import { computed } from 'vue';
const props = defineProps({
  text: {
    type: String,
    required: true
  },
  toastType: {
    type: String,
    required: false,
    default: 'success'
  }
});

const backgroundClass = computed(() => {
  if (props.toastType === 'success') {
    return 'bg-green-800'
  }
  return 'bg-red-800';
});
</script>
