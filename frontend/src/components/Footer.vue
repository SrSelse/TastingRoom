<template>
  <footer class="bg-gray-800 text-white py-8">
    <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
      <div class="flex flex-col md:flex-row justify-between items-center">
        <!-- App Info -->
        <div class="mb-6 md:mb-0 text-center md:text-left">
          <h2 class="text-xl font-bold">{{ appName }}</h2>
          <p class="text-gray-300 mt-2 max-w-md">
            The ultimate platform for beverage enthusiasts to host and join tasting rooms, score beverages, and share experiences.
          </p>
        </div>

        <!-- Buy Me a Coffee Button -->
        <div class="flex flex-col items-center">
          <p class="text-gray-300 mb-2">Like this project?</p>
          <a
            href="https://www.buymeacoffee.com/beanlink"
            target="_blank"
            rel="noopener noreferrer"
            class="inline-flex items-center px-4 py-2 bg-indigo-500 hover:bg-yellow-600 rounded-md text-gray-900 font-medium transition-colors"
          >
            <span class="mr-2">🍺</span>
            Buy me a beer or a bottle of wine
            <span class="ml-2">🍾</span>
          </a>
        </div>
      </div>

      <!-- Copyright -->
      <div class="mt-8 pt-8 border-t border-gray-700 text-center text-gray-400 text-sm">
        <p>© {{ new Date().getFullYear() }} {{ appName }}. All rights reserved.</p>
      </div>
    </div>
  </footer>
</template>

<script setup>
import { appName } from '@/classes/variables.js';
</script>
