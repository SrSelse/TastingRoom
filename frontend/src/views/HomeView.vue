<template>
  <div class="min-h-screen bg-gray-900 text-white">
    <!-- Header -->
    <header-component />

    <!-- Hero Section -->
    <main class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-12">
      <section class="text-center">
        <h1 class="text-4xl font-extrabold mb-6">
          Host and Join Tasting Rooms
        </h1>
        <p class="text-xl text-gray-300 mb-8 max-w-3xl mx-auto">
          Welcome to <span class="font-semibold text-indigo-400">{{ appName }}</span>, the ultimate platform for beverage enthusiasts.
          Create your own tasting rooms, add your favorite beverages, and host tastings - either in person or online.
          Our platform lets you score and review beverages in real-time using your own devices, with results updating automatically for everyone in the room.
          Whether you're a casual drinker or a seasoned connoisseur, {{ appName }} makes it easy to share your passion and discover new flavors together.
        </p>
        <div class="mt-10">
          <router-link
            to="/rooms"
            class="inline-block bg-indigo-600 hover:bg-indigo-700 text-white font-medium py-3 px-6 rounded-md transition-colors"
          >
            Get Started
          </router-link>
        </div>
      </section>

      <!-- Features Section -->
      <section class="mt-20 grid grid-cols-1 md:grid-cols-3 gap-8 text-center">
        <div class="p-6 rounded-lg bg-gray-800">
          <h2 class="text-xl font-bold mb-4">Create Rooms</h2>
          <p class="text-gray-300">
            Set up private or public tasting rooms in seconds. Customize your room with the beverages you want to taste.
          </p>
        </div>
        <div class="p-6 rounded-lg bg-gray-800">
          <h2 class="text-xl font-bold mb-4">Real-Time Scoring</h2>
          <p class="text-gray-300">
            Score beverages on the fly using your phone, tablet, or computer. See results update instantly for everyone in the room.
          </p>
        </div>
        <div class="p-6 rounded-lg bg-gray-800">
          <h2 class="text-xl font-bold mb-4">Join Anywhere</h2>
          <p class="text-gray-300">
            Participate in tastings from anywhere. Connect with friends, family, or fellow enthusiasts - no matter the distance.
          </p>
        </div>
      </section>
    </main>
    <footer-component />
  </div>
</template>

<script setup>
import HeaderComponent from '@components/Header.vue';
import FooterComponent from '@components/Footer.vue';

import { appName } from '@/classes/variables.js';
</script>
