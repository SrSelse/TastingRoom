<template>
  <div class="min-h-screen bg-gray-900 text-white">
    <header-component />
    <new-room />
    <footer-component />
  </div>
</template>

<script setup>
import HeaderComponent from '@components/Header.vue';
import FooterComponent from '@components/Footer.vue';
import NewRoom from '@components/NewRoom.vue';

</script>
