<template>
  <div class="min-h-screen bg-gray-900 text-white py-12 px-4 sm:px-6 lg:px-8">
    <!-- Main Content -->
    <div class="max-w-7xl mx-auto">
      <div class="flex flex-col sm:flex-row sm:justify-between sm:items-center mb-8">
        <div class="flex-col mb-6 sm:mb-0">
          <h1 class="text-3xl font-bold mb-2">{{ roomInfo.name }}</h1>
          <small class="text-gray-400">{{ roomInfo.plannedDate }}</small>
          <div v-if="isAdmin">
            <button
              class="inline-block bg-indigo-600 hover:bg-indigo-700 text-white font-medium py-2 px-4 rounded-md transition-colors"
              @click="editRoom"
            >
              ✏️ Edit room info
            </button>
          </div>

          <!-- Invitation Code Section -->
        </div>
        <div class="mt-2 flex items-center space-x-2">
          <div class="bg-gray-700 px-2 py-1 rounded text-xs">
            Invitation code
          </div>
          <div class="bg-gray-700 px-2 py-1 rounded text-xs">
            <span class="text-indigo-300 font-mono">{{ roomInfo.code }}</span>
          </div>
          <button
            @click="copyInvitationCode"
            class="text-gray-400 hover:text-white p-1 focus:outline-none cursor-pointer"
            title="Copy invitation code"
          >
            <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" fill="currentColor" class="bi bi-copy" viewBox="0 0 16 16">
              <path fill-rule="evenodd" d="M4 2a2 2 0 0 1 2-2h8a2 2 0 0 1 2 2v8a2 2 0 0 1-2 2H6a2 2 0 0 1-2-2zm2-1a1 1 0 0 0-1 1v8a1 1 0 0 0 1 1h8a1 1 0 0 0 1-1V2a1 1 0 0 0-1-1zM2 5a1 1 0 0 0-1 1v8a1 1 0 0 0 1 1h8a1 1 0 0 0 1-1v-1h1v1a2 2 0 0 1-2 2H2a2 2 0 0 1-2-2V6a2 2 0 0 1 2-2h1v1z"/>
            </svg>
          </button>
          <!-- Admin Button -->
          <button
            @click="showAdminModal = true"
            class="text-xs bg-gray-700 hover:bg-gray-600 text-white px-2 py-1 rounded"
            title="Manage admins"
          >
            Manage Users
          </button>
        </div>
        <div class="flex">
          <label
            for="showRatings"
            class="block text-gray-500 font-bold md:text-right mb-1 md:mb-0 pr-4"
          >
            Show ratings
          </label>
          <input
            id="showRatings"
            v-model="showRatings"
            class="rounded-full h-4 w-4 cursor-pointer"
            type="checkbox"
          />
        </div>

        <button
          @click="openAddBeerModal"
          class="inline-block bg-indigo-600 hover:bg-indigo-700 text-white font-medium py-2 px-4 rounded-md transition-colors"
        >
          Add Beverage
        </button>
      </div>

      <!-- Beers List -->
      <div v-if="loading" class="text-center py-12">
        <p class="text-gray-300">Loading beverages...</p>
      </div>
      <div v-else-if="error" class="text-center py-12">
        <p class="text-red-400">Error loading beverages: {{ error }}</p>
      </div>
      <div v-else>
        <div v-if="beers.length === 0" class="text-center py-12">
          <p class="text-gray-300">No beverages added yet. Add your first beverage!</p>
        </div>
        <div v-else class="divide-y divide-gray-700">
          <beer-list-item
            v-for="beer in beers"
            :key="beer.id"
            :beer="beer"
            :room-id="roomId"
            :show-ratings="showRatings"
            @edit="editBeer(beer)"
          />
        </div>
      </div>
    </div>

    <add-beer-modal
      v-if="showAddBeerModal"
      :room-id="roomId"
      @close="showAddBeerModal = false"
      @added="beerAdded"
    />
    <!-- Add Beer Modal -->

    <!-- Admin Modal -->
    <handle-users-modal
      v-if="showAdminModal"
      :users="users"
      :room-id="roomId"
      :is-admin="isAdmin"
      @close="showAdminModal = false"
      @refresh="fetchUsers"
    />

    <edit-beer-modal
      v-if="showEditBeerModal"
      :data="editBeerModalInfo"
      @saved="fetchBeers"
      @close="closeEditBeerModal"
    />

    <edit-room-modal
      v-if="showEditRoomModal"
      :data="editRoomModalInfo"
      @saved="fetchRoomInfo"
      @close="closeEditRoomModal"
    />

  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue';
import { useRouter } from 'vue-router';
import BeerListItem from './BeerListItem.vue';
import AddBeerModal from '@components/modals/AddBeerModal.vue';
import EditBeerModal from '@components/modals/EditBeerModal.vue';
import EditRoomModal from '@components/modals/EditRoomModal.vue';
import HandleUsersModal from '@components/modals/HandleUsersModal.vue';

const props = defineProps({
  roomId: {
    type: String,
    required: true
  }
});
const beers = ref([]);
const users = ref([]);
const roomInfo = ref({name: 'Tasting room'});
const loading = ref(true);
const error = ref(null);
const showAddBeerModal = ref(false);
const showAdminModal = ref(false);
const isAdmin = ref(false);
const router = useRouter();
const showRatings = ref(false);

// Copy invitation code to clipboard
const copyInvitationCode = () => {
  if (!roomInfo.value.code) return;

  navigator.clipboard.writeText(roomInfo.value.code)
    .then(() => {
      console.log('Invitation code copied to clipboard');
    })
    .catch(err => {
      console.error('Failed to copy: ', err);
    });
};

const checkAdminStatus = async () => {
  try {
    const response = await fetch(`${import.meta.env.VITE_API_URL}/api/room/${props.roomId}/is-admin`, {
      headers: {
        'Authorization': `Bearer ${localStorage.getItem('token')}`
      }
    });
    if (!response.ok) throw new Error('Failed to check admin status');
    isAdmin.value = await response.json();
  } catch (err) {
    console.error(err);
  }
};



// Remove user from room

// Simulate fetching beers for the room
const fetchBeers = async () => {
  try {
    loading.value = true;
    const response = await fetch(`${import.meta.env.VITE_API_URL}/api/room/${props.roomId}/beers`, {
      headers: {
        'Authorization': `Bearer ${localStorage.getItem('token')}`
      }
    });
    beers.value = await response.json();
  } catch (err) {
    error.value = err.message;
  } finally {
    loading.value = false;
  }
};

const fetchRoomInfo = async () => {
  try {
    loading.value = true;
    const response = await fetch(`${import.meta.env.VITE_API_URL}/api/room/${props.roomId}`, {
      headers: {
        'Authorization': `Bearer ${localStorage.getItem('token')}`
      }
    });
    if (!response.ok) {
      throw new Error('Something went wrong while fetching room info');
    }
    const data = await response.json();
    if (!data) {
      return;
    }

    const date = new Date(data.plannedDate).toLocaleString('sv-SE').replace('T', ' ').substring(0, 16);
    roomInfo.value = {
      ...data,
      plannedDate: date
    };
  } catch (err) {
    error.value = err.message;
    return;
  } finally {
    loading.value = false;
  }
}

const fetchUsers = async () => {
  try {
    loading.value = true;
    const response = await fetch(`${import.meta.env.VITE_API_URL}/api/room/${props.roomId}/users`, {
      headers: {
        'Authorization': `Bearer ${localStorage.getItem('token')}`
      }
    });
    if (!response.ok) {
      throw new Error('Something went wrong while fetching participants');
    }
    users.value = await response.json();
  } catch (err) {
    error.value = err.message;
    return;
  } finally {
    loading.value = false;
  }
}

const fetchInfo = async () => {
  await fetchRoomInfo();
  await fetchBeers();
  await fetchUsers();
}

// Simulate adding a new beer via AJAX

const openAddBeerModal = () => {
  showAddBeerModal.value = true;
};

const beerAdded = (beer) => {
  beers.value.push(beer);
};


const showEditBeerModal = ref(false);
const editBeerModalInfo = ref({});
const editBeer = (beer) => {
  editBeerModalInfo.value = {
    id: beer.id,
    roomId: props.roomId,
    name: beer.name,
    style: beer.style,
    pictureUrl: beer.pictureUrl
  };
  showEditBeerModal.value = true;
}

const closeEditBeerModal = () => {
  editBeerModalInfo.value = {};
  showEditBeerModal.value = false;
}

const showEditRoomModal = ref(false);
const editRoomModalInfo = ref({});
const editRoom = () => {
  editRoomModalInfo.value = {
    id: props.roomId,
    name: roomInfo.value.name,
    description: roomInfo.value.description,
    plannedDate: roomInfo.value.plannedDate
  };
  showEditRoomModal.value = true;
}

const closeEditRoomModal = () => {
  editRoomModalInfo.value = {};
  showEditRoomModal.value = false;
}

onMounted(() => {
  fetchInfo();
  checkAdminStatus();
});
</script>
