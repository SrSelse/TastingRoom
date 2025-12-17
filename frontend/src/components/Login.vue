<template>
  <div class="min-h-screen flex items-center justify-center bg-gray-900 py-12 px-4 sm:px-6 lg:px-8">
    <div class="max-w-md w-full space-y-8">
      <div>
        <h2 class="mt-6 text-center text-3xl font-extrabold text-white">
          Sign in to your account
        </h2>
      </div>
      <input type="hidden" name="remember" value="true" />
      <div class="rounded-md shadow-sm -space-y-px">
        <div>
          <label for="username" class="sr-only">Username</label>
          <input
            id="username"
            v-model="username"
            name="username"
            type="text"
            autocomplete="username"
            required
            class="appearance-none rounded-none relative block w-full px-3 py-2 border border-gray-600 bg-gray-800 text-gray-100 placeholder-gray-400 rounded-t-md focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 focus:z-10 sm:text-sm"
            placeholder="Username"
            @keyup.enter="handleLogin"
          />
        </div>
        <div>
          <label for="password" class="sr-only">Password</label>
          <input
            id="password"
            v-model="password"
            name="password"
            type="password"
            autocomplete="current-password"
            required
            class="appearance-none rounded-none relative block w-full px-3 py-2 border border-gray-600 bg-gray-800 text-gray-100 placeholder-gray-400 rounded-b-md focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 focus:z-10 sm:text-sm"
            placeholder="Password"
            @keyup.enter="handleLogin"
          />
        </div>
      </div>


      <div>
        <button
          @click="handleLogin"
          class="group relative w-full flex justify-center py-2 px-4 border border-transparent text-sm font-medium rounded-md text-white bg-indigo-600 hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500 focus:ring-offset-gray-900"
          @keyup.enter="handleLogin"
        >
          Sign in
        </button>
      </div>

      <div class="flex items-center justify-between">
        <div class="text-sm">
          <a href="/register" class="font-medium text-indigo-400 hover:text-indigo-300">
            Not a user? Register
          </a>
        </div>
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
const error = ref('');
const router = useRouter();

const handleLogin = async () => {
  try {
    error.value = '';
    const resp = await fetch(`${import.meta.env.VITE_API_URL}/auth/login`, {
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
      if ([401, 403].includes(resp.status)) {
        throw new Error('Username or password incorrect');
      }
      throw new Error('Internal server error');
    }
    const json = await resp.json();
    localStorage.setItem('token', json.token);
    localStorage.setItem('user', JSON.stringify(json.user));
    router.push('/');
  } catch (e) {
    error.value = e;
  }

};
</script>
