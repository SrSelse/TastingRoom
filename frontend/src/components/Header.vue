<template>

  <header class="bg-gray-800 shadow-md text-white">
    <nav class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-4 flex justify-between items-center">
      <div class="flex items-center space-x-4">

        <router-link
          to="/"
          class="text-gray-300 hover:text-white px-3 py-2 rounded-md text-sm font-medium transition-colors"
        >
          <span class="text-xl font-bold">
            {{ appName }}
          </span>
        </router-link>
      </div>
      <div v-if="isAuthenticated">
        <router-link
          to="/rooms"
          class="text-gray-300 hover:text-white px-3 py-2 rounded-md text-sm font-medium transition-colors"
        >
          Rooms
        </router-link>
        <a
          href="#"
          class="text-gray-300 hover:text-white px-3 py-2 rounded-md text-sm font-medium transition-colors"
          @click="signOut"
        >
          Sign out
        </a>

      </div>
      <div v-else>
        <router-link
          to="/login"
          class="text-gray-300 hover:text-white px-3 py-2 rounded-md text-sm font-medium transition-colors"
        >
          Login
        </router-link>
        <router-link
          to="/register"
          class="text-gray-300 hover:text-white px-3 py-2 rounded-md text-sm font-medium transition-colors"
        >
          Register
        </router-link>
      </div>
    </nav>
  </header>
</template>

<script setup>
import { ref } from 'vue';
import { RouterLink, useRouter} from 'vue-router';
import { appName } from '@/classes/variables.js';
const router = useRouter();

const isAuthenticated = ref(localStorage.getItem('token'));
const signOut = () => {
  localStorage.removeItem('token');
  localStorage.removeItem('user');
  isAuthenticated.value = false;
  router.push('/');
}
</script>
