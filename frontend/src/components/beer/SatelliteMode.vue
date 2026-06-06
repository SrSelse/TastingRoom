<template>
  <div class="bg-gray-800 p-6 rounded-lg">
    <h2 class="text-xl font-semibold mb-4">Rate This</h2>
    <div class="mb-4">
      <label for="rating" class="block text-sm font-medium text-gray-300 mb-2">Your Rating</label>
      <input
        id="rating"
        v-model="myRating.rating"
        type="number"
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

    <!-- Published Ratings -->
    <div v-if="isPublished && users.length > 0" class="mt-8">
      <h2 class="text-xl font-semibold mb-4">Ratings</h2>
      <div class="space-y-4">
        <user-row
          v-for="user in users"
          :key="user.id"
          :user="user"
          :showRating="true"
        />
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, onBeforeUnmount, watch } from 'vue';
import Toast from '@components/Toast.vue';
import UserRow from '@/components/beer/UserRow.vue';
import { centrifuge } from '@/classes/centrifuge.js';
import { useRouter } from 'vue-router';

const router = useRouter();

const props = defineProps({
  beerId: String,
  roomId: String,
  published: {
    type: Boolean,
    default: false,
  },
});

const emit = defineEmits(['rating-submitted']);
const showNotes = ref(false);
const showToast = ref(false);
const isPublished = ref(props.published);
const users = ref([]);
const myRating = ref({
  rating: 0,
  note: ""
});

watch(
  () => props.published,
  (val) => {
    isPublished.value = val;
    if (val) fetchRatings();
  },
);

watch(
  () => props.beerId,
  () => {
    isPublished.value = props.published;
    fetchMyRating();
    if (isPublished.value) fetchRatings();
  },
);

watch(
  () => props.roomId,
  () => {
    router.push({name: 'room', params: {roomId: props.roomId}});
  },
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

    showToast.value = true;
    setTimeout(() => {
      showToast.value = false;
    }, 3000);

    emit('rating-submitted');
    fetchMyRating()
  } catch (err) {
    console.error(err);
  }
};

const fetchMyRating = async () => {
  try {
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

onMounted(() => {
  fetchMyRating();
  if (isPublished.value) fetchRatings();
});

await centrifuge.init();

const onNextBeerPublication = (ctx) => {
  router.push({name: 'beer', params: {roomId: props.roomId, beerId: ctx.data.beerId}});
};

const onBeerPublication = (ctx) => {
  if (ctx.data.reason === 'ratings-published') {
    isPublished.value = true;
    fetchRatings();
  } else if (ctx.data.reason === 'ratings-unpublished') {
    isPublished.value = false;
    users.value = [];
  }
};

const subNextBeer = centrifuge.newSubscription(`rooms:${props.roomId}-next-beer`);
subNextBeer.on('publication', onNextBeerPublication).subscribe();

let subBeer = null;
const subscribeToBeer = () => {
  subBeer = centrifuge.newSubscription(`beers:beer-${props.beerId}`);
  subBeer.on('publication', onBeerPublication).subscribe();
};

const unsubscribeFromBeer = () => {
  if (subBeer) {
    subBeer.off('publication', onBeerPublication);
    centrifuge.removeSubscription(subBeer);
    subBeer = null;
  }
};

subscribeToBeer();

watch(
  () => props.beerId,
  () => {
    unsubscribeFromBeer();
    subscribeToBeer();
  },
);

onBeforeUnmount(() => {
  subNextBeer.off('publication', onNextBeerPublication);
  centrifuge.removeSubscription(subNextBeer);
  unsubscribeFromBeer();
});
</script>
