<template>
  <div class="bg-gray-800 p-6 rounded-lg">
    <h2 class="text-xl font-semibold mb-4">Rate This</h2>
    <div class="mb-4">
      <label for="rating" class="block text-sm font-medium text-gray-300 mb-2">Your Rating (0-5)</label>
      <input
        id="rating"
        v-model="myRating.rating"
        type="number"
        min="0"
        max="5"
        step="0.1"
        class="w-full px-3 py-2 bg-gray-700 border border-gray-600 rounded-md text-white focus:outline-none focus:ring-indigo-500 focus:border-indigo-500"
      />
    </div>

    <!-- Expandable Notes Section -->
    <div class="mb-4">
      <button
        @click="showNotes = !showNotes"
        class="flex items-center text-sm font-medium text-gray-300 mb-2 focus:outline-none"
      >
        <span>{{ showNotes ? '▼' : '▶' }} Add Notes</span>
      </button>
      <div v-if="showNotes" class="mt-2">
        <textarea
          v-model="myRating.note"
          rows="4"
          placeholder="Write your tasting notes here..."
          class="w-full px-3 py-2 bg-gray-700 border border-gray-600 rounded-md text-white focus:outline-none focus:ring-indigo-500 focus:border-indigo-500"
        ></textarea>
      </div>
    </div>

    <button
      @click="submitRating"
      class="px-4 py-2 bg-indigo-600 hover:bg-indigo-700 rounded-md text-white"
    >
      Submit Rating
    </button>

    <!-- Toast Notification -->
    <toast
      v-if="showToast"
      text="Rating submitted successfully!"
    />
  </div>
</template>

<script setup>
import { ref, onMounted, onBeforeUnmount, watch } from 'vue';
import Toast from '@components/Toast.vue';
import { centrifuge } from '@/classes/centrifuge.js';
import { useRouter } from 'vue-router';

const router = useRouter();



const props = defineProps({
  beerId: String,
  roomId: String,
});



const emit = defineEmits(['rating-submitted']);
const showNotes = ref(false);
const showToast = ref(false);
const myRating = ref({
  rating: 0,
  note: ""
});

watch(
  () => props.beerId,
  () => fetchMyRating(),
  // { immediate: true },
);

watch(
  () => props.roomId,
  () => {
    router.push({name: 'room', params: {roomId: props.roomId}});

  },
  // { immediate: true },
  );

const submitRating = async () => {
  try {
    const response = await fetch(`${import.meta.env.VITE_API_URL}/api/room/${props.roomId}/beers/${props.beerId}/rate`, {
      method: 'POST',
      headers: {
        'Authorization': `Bearer ${localStorage.getItem('token')}`,
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({
        ...myRating.value
      }),
    });
    if (!response.ok) throw new Error('Failed to submit rating');

    // Show toast notification
    showToast.value = true;
    setTimeout(() => {
      showToast.value = false;
    }, 3000);

    emit('rating-submitted');
  } catch (err) {
    console.error(err);
  }
};

const fetchMyRating = async () => {
  try {
    // Fetch beer details
    const response = await fetch(`${import.meta.env.VITE_API_URL}/api/room/${props.roomId}/beers/${props.beerId}/my-rating`, {
      headers: {
        'Authorization': `Bearer ${localStorage.getItem('token')}`
      }
    });
    if (!response.ok) throw new Error('Failed to fetch beverage');
    const data = await response.json();
    if (!!data) {
      myRating.value = data;
    } else {
      myRating.value = {
        rating: 0,
        note: ""
      };
    }
  } catch (err) {
    console.error(err);
  }
};

onMounted(() => {
  fetchMyRating();
});

await centrifuge.init();
const sub = centrifuge.newSubscription(`rooms:${props.roomId}-next-beer`);
sub.on('publication', (ctx) => {
  console.log(ctx);
  router.push({name: 'beer', params: {roomId: props.roomId, beerId: ctx.data.beerId}});
}).subscribe();

onBeforeUnmount(() => {
  centrifuge.removeSubscription(sub)
})
</script>
