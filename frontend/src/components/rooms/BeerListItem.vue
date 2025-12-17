<template>
  <li
    class="py-6 flex flex-col sm:flex-row sm:items-center justify-between"
  >
    <div class="flex items-center space-x-4">
      <div class="flex-shrink-0 h-12 w-12 rounded-full bg-gray-800 flex items-center justify-center text-gray-300">
        <img
          v-if="beer.pictureUrl"
          :src="beer.pictureUrl"
          alt="Beer image"
          class="h-12 w-12 rounded-full object-cover"
        />
        <span v-else>🍺</span>
      </div>
      <div>
        <h2 class="text-xl font-semibold">{{ beer.name }}</h2>
        <p class="text-gray-300">{{ beer.style }}</p>
      </div>
    </div>
    <div class="flex mt-4 sm:mt-0">
      <router-link
        :to="`/rooms/${roomId}/beer/${beer.id}`"
      >
        <div
          v-if="showRatings"
          class="flex-col"
        >
          <div
            v-if="beer.average"
            class="flex items-center justify-center"
          >
            <span
              class="text-2xl font-bold text-yellow-400 mr-1"
              :class="{
                'blur-lg': !beer.published
              }"
            >
              {{ beer.average.toFixed(1) }}
            </span>

            <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5 text-yellow-400" viewBox="0 0 20 20" fill="currentColor">
              <path d="M9.049 2.927c.3-.921 1.603-.921 1.902 0l1.07 3.292a1 1 0 00.95.69h3.462c.969 0 1.371 1.24.588 1.81l-2.8 2.034a1 1 0 00-.364 1.118l1.07 3.292c.3.921-.755 1.688-1.54 1.118l-2.8-2.034a1 1 0 00-1.175 0l-2.8 2.034c-.784.57-1.838-.197-1.539-1.118l1.07-3.292a1 1 0 00-.364-1.118L2.98 8.72c-.783-.57-.38-1.81.588-1.81h3.461a1 1 0 00.951-.69l1.07-3.292z" />
            </svg>
          </div>
        </div>
        <span
          class="text-indigo-400 hover:text-indigo-300 text-sm font-medium"
        >
          View Details
        </span>
      </router-link>
      <button
        type="button"
        class="ml-4 cursor-pointer"
        @click="emit('edit')"
      >
        ✏️
      </button>
    </div>
  </li>
</template>

<script setup>
const props = defineProps({
  beer: {
    type: Object,
    required: true
  },
  roomId: {
    type: String,
    required: true
  },
  showRatings: {
    type: Boolean,
    required: true
  }
});
const emit = defineEmits([ 'edit' ]);
</script>
