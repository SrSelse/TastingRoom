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
          @click="handleRegister"
          class="group relative w-full flex justify-center py-2 px-4 border border-transparent text-sm font-medium rounded-md text-white bg-indigo-600 hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500 focus:ring-offset-gray-900"
        >
          Sign up
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
const error = ref('');
const router = useRouter();

const handleRegister = async () => {
  if (username.value === '' || password.value === '' || confirmPassword.value === '') {
    error.value = 'Please fill in all fields';
    return;
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
  }
};
</script>
