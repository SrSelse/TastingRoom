<template>
  <div class="min-h-screen bg-gray-900 text-white">
    <header-component />
    <room :room-id="$route.params.roomId" />
    <footer-component />
  </div>
</template>

<script setup>
import HeaderComponent from '@components/Header.vue';
import FooterComponent from '@components/Footer.vue';
import Room from '@components/rooms/Room.vue';

</script>
