<template>
  <div class="min-h-screen bg-gray-900 text-white py-12 px-4 sm:px-6 lg:px-8">
    <div class="max-w-md mx-auto">
      <div class="bg-gray-800 rounded-lg shadow-lg p-8">
        <h1 class="text-2xl font-bold mb-6">Change Display Name</h1>

        <div class="mb-6">
          <label for="displayName" class="block text-sm font-medium text-gray-300 mb-2">
            Current Display Name
          </label>
          <div class="flex items-center space-x-2">
            <input
              id="displayName"
              v-model="currentDisplayName"
              type="text"
              readonly
              class="flex-grow px-3 py-2 bg-gray-700 border border-gray-600 rounded-md text-white focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 cursor-not-allowed"
            />
          </div>
        </div>

        <div class="mb-6">
          <label for="newDisplayName" class="block text-sm font-medium text-gray-300 mb-2">
            New Display Name
          </label>
          <input
            id="newDisplayName"
            v-model="newDisplayName"
            type="text"
            placeholder="Enter your new display name"
            class="w-full px-3 py-2 bg-gray-700 border border-gray-600 rounded-md text-white focus:outline-none focus:ring-indigo-500 focus:border-indigo-500"
            :class="{ 'border-red-500': error }"
            @keyup.enter="updateDisplayName"
          />
          <p v-if="error" class="mt-1 text-sm text-red-400">{{ error }}</p>
          <p class="mt-1 text-xs text-gray-400">
            Your display name will be visible to other users in tasting rooms.
          </p>
        </div>

        <div class="flex justify-end space-x-4">
          <button
            @click="updateDisplayName"
            :disabled="isUpdating || newDisplayName === currentDisplayName"
            class="px-4 py-2 bg-indigo-600 hover:bg-indigo-700 rounded-md text-white transition-colors disabled:opacity-50 disabled:cursor-not-allowed flex items-center"
          >
            <svg
              v-if="isUpdating"
              class="animate-spin -ml-1 mr-2 h-4 w-4 text-white"
              xmlns="http://www.w3.org/2000/svg"
              fill="none"
              viewBox="0 0 24 24"
            >
              <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
              <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
            </svg>
            <span>{{ isUpdating ? 'Updating...' : 'Update Name' }}</span>
          </button>
        </div>
      </div>

      <!-- Success Toast -->
      <toast
        v-if="showToast"
        :text="toastText"
        :toast-type="toastType"
      />
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue';
import { useRouter } from 'vue-router';
import Toast from '@/components/Toast.vue';

const router = useRouter();

// State
const currentDisplayName = ref('');
const newDisplayName = ref('');
const isUpdating = ref(false);
const error = ref('');
const showToast = ref(false);
const toastType = ref('');
const toastText = ref('');

// Fetch current display name
const fetchCurrentDisplayName = async () => {
  try {
    const response = await fetch(`${import.meta.env.VITE_API_URL}/api/user/profile`, {
      headers: {
        'Authorization': `Bearer ${localStorage.getItem('token')}`
      }
    });

    if (!response.ok) {
      throw new Error('Failed to fetch profile');
    }

    const data = await response.json();
    currentDisplayName.value = data.displayName || 'Not set';
    newDisplayName.value = data.displayName || '';
  } catch (err) {
    console.error('Error fetching display name:', err);
    error.value = 'Failed to load your current display name';
  }
};

// Update display name
const updateDisplayName = async () => {
  if (!newDisplayName.value) {
    error.value = 'Display name cannot be empty';
    return;
  }

  isUpdating.value = true;
  error.value = '';

  try {
    const response = await fetch(`${import.meta.env.VITE_API_URL}/api/user/updateProfile`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${localStorage.getItem('token')}`
      },
      body: JSON.stringify({
        displayName: newDisplayName.value
      })
    });

    if (!response.ok) {
      const errorData = await response.json();
      throw new Error(errorData.message || 'Failed to update display name');
    }

    // Update current display name
    currentDisplayName.value = newDisplayName.value;

    // Show success toast
    showToast.value = true;
    toastType.value = 'success';
    toastText.value = 'Display name updated';
    setTimeout(() => {
      showToast.value = false;
    }, 3000);

  } catch (err) {
    console.error('Error updating display name:', err);
    showToast.value = true;
    toastType.value = 'error';
    toastText.value = err.message
  } finally {
    isUpdating.value = false;
  }
};

// Reset form
const resetForm = () => {
  newDisplayName.value = currentDisplayName.value;
  error.value = '';
};

onMounted(() => {
  fetchCurrentDisplayName();
});
</script>

