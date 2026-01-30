<template>
  <div class="bg-gray-800 p-6 rounded-lg">
    <h2 class="text-xl font-semibold mb-4">User Ratings</h2>

    <!-- Publish Button and Average Rating -->
    <div class="flex justify-between items-center mb-6">
      <div>
        <button
          @click="publishRatings"
          class="px-4 py-2 bg-indigo-600 hover:bg-indigo-700 rounded-md text-white"
        >
          <span v-if="!localPublished">
            Publish Ratings
          </span>
          <span v-else>
            Unpublish Ratings
          </span>
        </button>
      </div>

      <!-- Average Rating Box -->
      <div
        v-if="users.length > 0 && averageRating !== null"
        class="bg-gray-700 p-4 rounded-lg text-center min-w-[120px]"
      >
        <p class="text-sm text-gray-400 mb-1">Average Rating</p>
        <div class="flex items-center justify-center">
          <span
            class="text-2xl font-bold text-yellow-400 mr-1"
            :class="{
              'blur-lg': !localPublished
            }"
          >
            {{ averageRating.toFixed(1) }}
          </span>

          <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5 text-yellow-400" viewBox="0 0 20 20" fill="currentColor">
            <path d="M9.049 2.927c.3-.921 1.603-.921 1.902 0l1.07 3.292a1 1 0 00.95.69h3.462c.969 0 1.371 1.24.588 1.81l-2.8 2.034a1 1 0 00-.364 1.118l1.07 3.292c.3.921-.755 1.688-1.54 1.118l-2.8-2.034a1 1 0 00-1.175 0l-2.8 2.034c-.784.57-1.838-.197-1.539-1.118l1.07-3.292a1 1 0 00-.364-1.118L2.98 8.72c-.783-.57-.38-1.81.588-1.81h3.461a1 1 0 00.951-.69l1.07-3.292z" />
          </svg>
        </div>
        <p class="text-xs text-gray-400 mt-1">{{ usersWithRatings}} ratings</p>
      </div>
    </div>

    <div class="space-y-4">
      <user-row
        v-for="user in users"
        :key="user.id"
        :user="user"
        :showRating="localPublished"
      />
    </div>
    <toast
      v-if="showToast"
      :text="publishedRatingsToast.text"
      :toast-type="publishedRatingsToast.type"
    />
  </div>
</template>

<script setup>
import { ref, onBeforeUnmount, watch, computed } from 'vue';
import UserRow from '@/components/beer/UserRow.vue';
import Toast from '@components/Toast.vue';
import { centrifuge } from '@/classes/centrifuge.js';
import { useRouter } from 'vue-router';
const router = useRouter();

const props = defineProps({
  roomId: {
    type: String,
    required: true,
  },
  beerId: {
    type: String,
    required: true,
  },
  published: {
    type: Boolean,
    required: true
  }
});

const users = ref([]);
const showToast = ref(false);
const localPublished = ref(props.published);

const publishedRatingsToast = computed(() => {
  if (localPublished.value) {
    return {
      text: "Rating published",
      type: "success"
    }
  }
  return {
    text: "Rating unpublished",
    type: "error"
  }
});

// Compute average rating
const averageRating = computed(() => {
  if (users.value.length === 0) return null;

  const sum = users.value.reduce((total, user) => {
    return user.rating ? total + user.rating : total;
  }, 0);

  const count = users.value.filter(user => user.rating).length;
  return count > 0 ? sum / count : 0;
});

const usersWithRatings = computed(() => {
  if (users.value.length === 0) return null;
return users.value.filter(user => user.rating).length

});

// Fetch ratings (users)
const fetchRatings = async () => {
  try {
    const response = await fetch(`${import.meta.env.VITE_API_URL}/api/room/${props.roomId}/beers/${props.beerId}/ratings`, {
      headers: {
        'Authorization': `Bearer ${localStorage.getItem('token')}`
      }
    });
    if (!response.ok) throw new Error('Failed to fetch ratings');
    users.value = await response.json();
  } catch (err) {
    console.error(err);
  }
};

// Publish ratings
const publishRatings = async () => {
  try {
    let url = `${import.meta.env.VITE_API_URL}/api/room/${props.roomId}/beers/${props.beerId}/publish`;
    if (localPublished.value) {
      url = `${import.meta.env.VITE_API_URL}/api/room/${props.roomId}/beers/${props.beerId}/unpublish`
    }
    const response = await fetch(url, {
      headers: {
        'Authorization': `Bearer ${localStorage.getItem('token')}`
      },
      method: 'POST',
    });
    if (!response.ok) throw new Error('Failed to publish ratings');
    localPublished.value = !localPublished.value;
    showToast.value = true;
    setTimeout(() => {
      showToast.value = false;
    }, 3000);
  } catch (err) {
    console.error(err);
  }
};


await centrifuge.init();
const sub = ref(null)
const createCentrifugoSub = () => {
  sub.value = centrifuge.newSubscription(`beers:beer-${props.beerId}`);
  sub.value.on('publication', (ctx) => {
    fetchRatings();
  }).subscribe();
}

watch(
  () => props.beerId,
  () => {
    localPublished.value = props.published;
    fetchRatings();
    if (sub.value) {
      centrifuge.removeSubscription(sub.value);
    }
    createCentrifugoSub();
  },
  { immediate: true },
);

watch(
  () => props.roomId,
  () => {
    router.push({name: 'room', params: {roomId: props.roomId}});
  },
);

onBeforeUnmount(() => {
  centrifuge.removeSubscription(sub.value)
})
</script>
