<template>
  <div class="min-h-screen bg-gray-900 text-white py-12 px-4 sm:px-6 lg:px-8">
    <div class="max-w-4xl mx-auto">
      <!-- Navigation Buttons -->
      <div class="flex justify-between mb-8">
        <router-link
          :to="`/rooms/${roomId}`"
          class="px-4 py-2 bg-gray-700 hover:bg-gray-600 rounded-md text-white text-sm font-medium transition-colors"
        >
          ← Back to Room
        </router-link>
        <button
          v-if="isAdmin"
          @click="goToRandomBeer"
          class="px-4 py-2 bg-indigo-700 hover:bg-indigo-600 rounded-md text-white text-sm font-medium transition-colors"
        >
          Random Beverage →
        </button>
      </div>

      <!-- Beer Info -->
      <div class="flex flex-col md:flex-row md:items-start mb-8">
        <div class="w-full md:w-1/3 mb-6 md:mb-0 md:mr-8 flex justify-center md:justify-start">
          <img
            v-if="beer.pictureUrl"
            :src="beer.pictureUrl"
            alt="Beer image"
            class="h-48 w-48 object-cover rounded-lg"
          />
          <div v-else class="h-48 w-48 bg-gray-800 rounded-lg flex items-center justify-center text-gray-400 text-4xl">
            🍺
          </div>
        </div>
        <div class="w-full md:w-2/3">
          <h1 class="text-3xl font-bold mb-2">{{ beer.name }}</h1>
          <p class="text-xl text-gray-300 mb-4">{{ beer.style }}</p>
        </div>
      </div>

      <!-- Mode Switcher -->
      <div class="mb-6">
        <div class="flex space-x-4">
          <button
            @click="mode = 'satellite'"
            :class="['px-4 py-2 rounded-md', mode === 'satellite' ? 'bg-indigo-600' : 'bg-gray-700']"
          >
            Satellite Mode
          </button>
          <button
            v-if="isAdmin"
            @click="mode = 'admin'"
            :class="['px-4 py-2 rounded-md', mode === 'admin' ? 'bg-indigo-600' : 'bg-gray-700']"
          >
            Admin Mode
          </button>
        </div>
      </div>

      <!-- Satellite Mode -->
      <Suspense>
        <SatelliteMode
          v-if="mode === 'satellite'"
          :beer-id="beerId"
          :room-id="roomId"
        />

      </Suspense>
      <Suspense>
        <!-- Admin Mode -->
        <template v-if="isAdmin && beer.id">
          <AdminMode
            v-show="mode === 'admin'"
            :room-id="roomId"
            :beer-id="beerId"
            :published="beer.published"
          />
        </template>
      </Suspense>
    </div>
  </div>
</template>

<script setup>
import { ref, watch, onMounted } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import SatelliteMode from '@/components/beer/SatelliteMode.vue';
import AdminMode from '@/components/beer/AdminMode.vue';

const props = defineProps({
  roomId: {
    type: String,
    required: true
  },
  beerId: {
    type: String,
    required: true
  },
});

const beer = ref({});
const isAdmin = ref(false);
const mode = ref('satellite');
const loading = ref(true);
const error = ref(null);
const router = useRouter();

const fetchBeer = async () => {
  try {
    beer.value = {};
    loading.value = true;
    // Fetch beer details
    const response = await fetch(`${import.meta.env.VITE_API_URL}/api/room/${props.roomId}/beers/${props.beerId}`, {
      headers: {
        'Authorization': `Bearer ${localStorage.getItem('token')}`
      }
    });
    if (!response.ok) throw new Error('Failed to fetch beverage');
    beer.value = await response.json();
  } catch (err) {
    error.value = err.message;
  } finally {
    loading.value = false;
  }
};

const checkAdminStatus = async () => {
  try {
    const response = await fetch(`${import.meta.env.VITE_API_URL}/api/room/${props.roomId}/is-admin`, {
      headers: {
        'Authorization': `Bearer ${localStorage.getItem('token')}`
      }
    });
    if (!response.ok) throw new Error('Failed to check admin status');
    isAdmin.value = await response.json();
  } catch (err) {
    console.error(err);
  }
};

const goToRandomBeer = async () => {
  try {
    const response = await fetch(`${import.meta.env.VITE_API_URL}/api/room/${props.roomId}/beers/random`, {
      headers: {
        'Authorization': `Bearer ${localStorage.getItem('token')}`
      }
    });
    if (!response.ok) throw new Error('Failed to check admin status');
    const data = await response.json();
    if (!!data) {
      router.push({name: 'beer', params: {roomId: props.roomId, beerId: data.id}});
    }
  } catch (err) {
    console.error(err)
  }
};

onMounted(() => {
  fetchBeer();
  checkAdminStatus();
});

watch(
  () => props.beerId,
  () => fetchBeer(),
  // { immediate: true },
);

watch(
  () => props.roomId,
  () => {
    router.push({name: 'room', params: {roomId: props.roomId}});

  },
  // { immediate: true },
  );
</script>
