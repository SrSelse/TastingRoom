<template>
  <div class="min-h-screen bg-gray-900 text-white">
    <header-component />
    <beer
      :room-id="roomId"
      :beer-id="beerId"
    />
    <footer-component />
  </div>
</template>

<script setup>
import { ref, watch } from 'vue'
import { useRoute } from 'vue-router'
import HeaderComponent from '@components/Header.vue';
import FooterComponent from '@components/Footer.vue';
import Beer from '@components/beer/Beer.vue';
const route = useRoute()


const roomId = ref(route.params.roomId);
const beerId = ref(route.params.beerId);

watch(
  () => route.params.roomId,
  (newId, oldId) => {
    if (newId !== oldId) {
      roomId.value = newId;
    }
  }
)

watch(
  () => route.params.beerId,
  (newId, oldId) => {
    if (newId !== oldId) {
      beerId.value = newId;
    }
  }
)

</script>
