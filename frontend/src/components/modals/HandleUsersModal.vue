<template>
  <div
    class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center p-4 z-50"
    @click.self="emit('close')"
  >
    <div class="bg-gray-800 rounded-lg p-6 max-w-2xl w-full max-h-[80vh] overflow-y-auto">
      <div class="flex justify-between items-center mb-6">
        <h2 class="text-xl font-bold">Manage Room Admins</h2>
        <button
          @click="emit('close')"
          class="text-gray-400 hover:text-white focus:outline-none"
        >
          <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
          </svg>
        </button>
      </div>

      <div class="mb-4">
        <p class="text-sm text-gray-400">Manage admin privileges for users in this room.</p>
      </div>

      <div class="overflow-x-auto">
        <table class="min-w-full divide-y divide-gray-700">
          <thead class="bg-gray-700">
            <tr>
              <th scope="col" class="px-6 py-3 text-left text-xs font-medium text-gray-300 uppercase tracking-wider">
                User
              </th>
              <th scope="col" class="px-6 py-3 text-left text-xs font-medium text-gray-300 uppercase tracking-wider">
                Admin
              </th>
              <th scope="col" class="px-6 py-3 text-right text-xs font-medium text-gray-300 uppercase tracking-wider">
                Actions
              </th>
            </tr>
          </thead>
          <tbody class="bg-gray-800 divide-y divide-gray-700">
            <tr v-for="user in users" :key="user.id">
              <td class="px-6 py-4 whitespace-nowrap">
                <div class="flex items-center">
                  <div class="flex-shrink-0 h-8 w-8 rounded-full bg-gray-700 flex items-center justify-center">
                    {{ user.name.charAt(0).toUpperCase() }}
                  </div>
                  <div class="ml-4">
                    <div class="text-sm font-medium text-white">{{ user.name }}</div>
                  </div>
                </div>
              </td>
              <td class="px-6 py-4 whitespace-nowrap">
                <input
                  type="checkbox"
                  :disabled="!isAdmin"
                  v-model="user.isAdmin"
                  @change="updateAdminStatus(user)"
                  class="rounded border-gray-600 text-indigo-600 focus:ring-indigo-500"
                />
              </td>
              <td class="px-6 py-4 whitespace-nowrap text-right text-sm font-medium">
                <button
                  :disabled="!isAdmin"
                  @click="removeUser(user.id)"
                  class="text-red-500 hover:text-red-400 focus:outline-none"
                  title="Remove user"
                >
                  <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
                  </svg>
                </button>
              </td>
            </tr>
          </tbody>
        </table>
      </div>

      <div class="mt-6 flex justify-between">
        <div>
          <button
            @click="leaveRoom"
            class="px-4 py-2 bg-red-700 hover:bg-gray-600 text-white rounded-md"
          >
            Leave room
          </button>
        </div>
        <div>
          <button
            @click="emit('close')"
            class="px-4 py-2 bg-gray-700 hover:bg-gray-600 text-white rounded-md"
          >
            Close
          </button>
        </div>
      </div>
      <template v-if="error">
        {{ error }}
      </template>
    </div>
  </div>

</template>

<script setup>
import { ref } from 'vue';
import { useRouter } from 'vue-router';

const router = useRouter();

const emit = defineEmits(['close', 'refresh']);

const props = defineProps({
  users: {
    type: Array,
    required: false,
    default: () => []
  },
  roomId: {
    type: String,
    required: true
  },
  isAdmin: {
    type: Boolean,
    required: true
  }
});
const error = ref('');

const leaveRoom = async (userId) => {
  if (!confirm('Are you sure you want to leave this room?')) {
    return;
  }

  try {
    const response = await fetch(`${import.meta.env.VITE_API_URL}/api/room/${props.roomId}/leave`, {
      method: 'POST',
      headers: {
        'Authorization': `Bearer ${localStorage.getItem('token')}`
      }
    });

    if (!response.ok) {
      const errorData = await response.json();
      throw new Error(errorData.message || 'Failed to remove user');
    }

    // Refresh users list
    router.push({name: 'roomlist'});
  } catch (err) {
    error.value = err.message;
    console.error('Error removing user:', err);
  }
};

const removeUser = async (userId) => {
  if (!confirm('Are you sure you want to remove this user from the room?')) {
    return;
  }

  try {
    const response = await fetch(`${import.meta.env.VITE_API_URL}/api/room/${props.roomId}/users/${userId}/remove`, {
      method: 'POST',
      headers: {
        'Authorization': `Bearer ${localStorage.getItem('token')}`
      }
    });

    if (!response.ok) {
      const errorData = await response.json();
      throw new Error(errorData.message || 'Failed to remove user');
    }

    emit('refresh');
    // Refresh users list
    // await fetchUsers();
  } catch (err) {
    error.value = err.message;
    console.error('Error removing user:', err);
  }
};

// Update admin status
const updateAdminStatus = async (user) => {
  try {
    const response = await fetch(`${import.meta.env.VITE_API_URL}/api/room/${props.roomId}/users/${user.id}/admin`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${localStorage.getItem('token')}`
      },
      body: JSON.stringify({
        isAdmin: user.isAdmin
      })
    });

    if (!response.ok) {
      const errorData = await response.json();
      throw new Error(errorData.message || 'Failed to update admin status');
    }

    emit('refresh');
    // Refresh users list
  } catch (err) {
    error.value = err.message;
    console.error('Error updating admin status:', err);
    // Revert the checkbox if the update failed
    user.isAdmin = !user.isAdmin;
  }
};
</script>
