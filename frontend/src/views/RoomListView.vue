<template>
  <div class="min-h-screen bg-gray-900 text-white">
    <header-component />
    <rooms-list />
    <footer-component />
  </div>
</template>

<script setup>
import HeaderComponent from '@components/Header.vue';
import FooterComponent from '@components/Footer.vue';
import RoomsList from '@components/rooms/RoomsList.vue';
</script>
