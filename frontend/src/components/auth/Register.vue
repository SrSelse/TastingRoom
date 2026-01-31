<template>
  <div class="min-h-screen flex items-center justify-center bg-gray-900 py-12 px-4 sm:px-6 lg:px-8">
    <div class="max-w-md w-full space-y-8">
      <div>
        <h2 class="mt-6 text-center text-3xl font-extrabold text-white">
          Create your account
        </h2>
      </div>
      <div class="rounded-md shadow-sm -space-y-px">
        <div>
          <label for="username" class="sr-only">Username</label>
          <input
            id="username"
            v-model="username"
            name="username"
            type="text"
            autocomplete="username"
            class="appearance-none rounded-none relative block w-full px-3 py-2 border border-gray-600 bg-gray-800 text-gray-100 placeholder-gray-400 rounded-t-md focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 focus:z-10 sm:text-sm"
            placeholder="Username"
            @keyup.enter="handleRegister"
          />
        </div>
        <div>
          <label for="password" class="sr-only">Password</label>
          <input
            id="password"
            v-model="password"
            name="password"
            type="password"
            autocomplete="new-password"
            class="appearance-none rounded-none relative block w-full px-3 py-2 border border-gray-600 bg-gray-800 text-gray-100 placeholder-gray-400 focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 focus:z-10 sm:text-sm"
            placeholder="Password"
            @keyup.enter="handleRegister"
          />
        </div>
        <div>
          <label for="confirm-password" class="sr-only">Confirm Password</label>
          <input
            id="confirm-password"
            v-model="confirmPassword"
            name="confirm-password"
            type="password"
            autocomplete="new-password"
            class="appearance-none rounded-none relative block w-full px-3 py-2 border border-gray-600 bg-gray-800 text-gray-100 placeholder-gray-400 rounded-b-md focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 focus:z-10 sm:text-sm"
            placeholder="Confirm Password"
            @keyup.enter="handleRegister"
          />
        </div>
      </div>

      <div>
        <button
          class="group relative w-full flex justify-center py-2 px-4 border border-transparent text-sm font-medium rounded-md text-white bg-indigo-600 hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500 focus:ring-offset-gray-900"
          @click="handleRegister"
        >
          Sign up
          <svg 
            v-if="registerInProgress"
            class="mr-3 ml-1 size-5 animate-spin text-white"
            xmlns="http://www.w3.org/2000/svg"
            fill="none"
            viewBox="0 0 24 24"
          >
            <circle 
              class="opacity-25"
              cx="12"
              cy="12"
              r="10"
              stroke="currentColor"
              stroke-width="4"
            >
            </circle>
            <path 
              class="opacity-75"
              fill="currentColor"
              d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"
            >
            </path>
          </svg>
        </button>
      </div>
      <div
        v-if="error"
        class="mt-6 text-center text-lg text-red-100"
      >
        {{ error }}
      </div>

    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue';
import { useRouter } from 'vue-router';

const username = ref('');
const password = ref('');
const confirmPassword = ref('');
const registerInProgress = ref(false);
const error = ref('');
const router = useRouter();

const handleRegister = async () => {
  if (registerInProgress.value) {
    return;
  }

  if (username.value === '' || password.value === '' || confirmPassword.value === '') {
    error.value = 'Please fill in all fields';
    return;
  }
  if (/^[a-zA-Z0-9_-]+$/.test(username.value) === false) {
    error.value = 'Username must only contain A-Z, a-z, 0-9, _ and -'
    return
  }
  if (password.value !== confirmPassword.value) {
    error.value = 'Passwords do not match!';
    return;
  }

  if (password.value.length < 8) {
    error.value = 'Password needs to be at least 8 characters';
    return;
  }

  try {
    registerInProgress.value = true;
    let resp = await fetch(`${import.meta.env.VITE_API_URL}/auth/register`, {
      method: 'POST',
      body: JSON.stringify({
        password: password.value,
        username: username.value
      }),
      headers: {
        'Content-Type': 'application/json'
      }
    });
    if (!resp.ok) {
      if (resp.status === 422) {
        throw new Error('Username already taken');
      }
      throw new Error('Internal server error');
    }

    resp = await fetch(`${import.meta.env.VITE_API_URL}/auth/login`, {
      method: 'POST',
      body: JSON.stringify({
        password: password.value,
        username: username.value
      }),
      headers: {
        'Content-Type': 'application/json'
      }
    });
    if (!resp.ok) {
      throw new Error('Internal server error');
    }
    const json = await resp.json();
    localStorage.setItem('token', json.token);
    localStorage.setItem('user', JSON.stringify(json.user));

    router.push({name: 'roomlist'});

  } catch (e) {
    error.value = e;
    console.error(e);
  } finally {
      registerInProgress.value = false;
  }
};
</script>
