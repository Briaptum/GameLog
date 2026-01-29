<template>
  <div class="min-h-screen flex items-center justify-center bg-gray-50 py-12 px-4 sm:px-6 lg:px-8 relative">
    <!-- ANHELM Watermark Logo - Bottom Left -->
    <div class="absolute bottom-6 left-6 flex items-center">
      <div class="w-4 h-4 flex flex-col items-center justify-center mr-2">
        <!-- Top two dots -->
        <div class="flex gap-1 mb-0.5">
          <div class="w-1 h-1 bg-cyan-400 rounded-full"></div>
          <div class="w-1 h-1 bg-cyan-400 rounded-full"></div>
        </div>
        <!-- Bottom dot -->
        <div class="w-1 h-1 bg-cyan-400 rounded-full"></div>
      </div>
      <span class="text-xs font-light text-zinc-500 tracking-wider">ANHELM</span>
    </div>

    <div class="max-w-md w-full space-y-6">
      <!-- Customer Logo Section -->
      <div class="text-center">
        <div class="mx-auto flex items-center justify-center mb-8">
          <!-- ANHELM Logo (dots + text) -->
          <div class="flex items-center">
            <div class="w-8 h-8 flex flex-col items-center justify-center mr-3">
              <!-- Top two dots -->
              <div class="flex gap-1.5 mb-1">
                <div class="w-1.5 h-1.5 bg-cyan-400 rounded-full"></div>
                <div class="w-1.5 h-1.5 bg-cyan-400 rounded-full"></div>
              </div>
              <!-- Bottom dot -->
              <div class="w-1.5 h-1.5 bg-cyan-400 rounded-full"></div>
            </div>
            <h1 class="text-2xl font-light text-zinc-700 tracking-wider">ANHELM</h1>
          </div>
        </div>
        <div class="mb-2">
          <h2 class="text-lg font-medium text-zinc-600">Welcome back</h2>
          <p class="text-sm text-zinc-500 mt-1">Sign in to your account</p>
        </div>
      </div>

      <!-- Login Form -->
      <form @submit.prevent="$emit('login', loginForm)" class="bg-white py-8 px-6 shadow-lg rounded-2xl border border-gray-100">
        <div class="space-y-6">
          <!-- Error Message -->
          <div v-if="error" class="bg-red-50 border border-red-200 rounded-md p-4">
            <div class="flex">
              <div class="flex-shrink-0">
                <svg class="h-5 w-5 text-red-400" viewBox="0 0 20 20" fill="currentColor">
                  <path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zM8.707 7.293a1 1 0 00-1.414 1.414L8.586 10l-1.293 1.293a1 1 0 101.414 1.414L10 11.414l1.293 1.293a1 1 0 001.414-1.414L11.414 10l1.293-1.293a1 1 0 00-1.414-1.414L10 8.586 8.707 7.293z" clip-rule="evenodd" />
                </svg>
              </div>
              <div class="ml-3">
                <p class="text-sm text-red-800">{{ error }}</p>
              </div>
            </div>
          </div>

          <!-- Email Field -->
          <div>
            <label for="email" class="block text-sm font-medium text-zinc-700 mb-2">
              Email address
            </label>
            <div class="relative">
              <div class="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
                <font-awesome-icon icon="envelope" class="h-5 w-5 text-gray-400" />
              </div>
              <input
                id="email"
                v-model="loginForm.email"
                name="email"
                type="email"
                autocomplete="email"
                required
                class="block w-full pl-10 pr-3 py-3 border border-gray-300 rounded-md placeholder-gray-400 focus:outline-none focus:ring-1 focus:ring-cyan-500 focus:border-cyan-500 transition duration-150 ease-in-out text-sm"
                placeholder="Enter your email"
              />
            </div>
          </div>

          <!-- Password Field -->
          <div>
            <label for="password" class="block text-sm font-medium text-zinc-700 mb-2">
              Password
            </label>
            <div class="relative">
              <div class="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
                <font-awesome-icon icon="lock" class="h-5 w-5 text-gray-400" />
              </div>
              <input
                id="password"
                v-model="loginForm.password"
                name="password"
                type="password"
                autocomplete="current-password"
                required
                class="block w-full pl-10 pr-3 py-3 border border-gray-300 rounded-md placeholder-gray-400 focus:outline-none focus:ring-1 focus:ring-cyan-500 focus:border-cyan-500 transition duration-150 ease-in-out text-sm"
                placeholder="Enter your password"
              />
            </div>
          </div>

          <!-- Login Button -->
          <div>
            <button
              type="submit"
              :disabled="loading"
              class="group relative w-full flex justify-center py-3 px-4 border border-transparent text-sm font-medium rounded-md text-white bg-cyan-600 hover:bg-cyan-700 focus:outline-none focus:ring-1 focus:ring-cyan-500 disabled:opacity-50 disabled:cursor-not-allowed transition duration-150 ease-in-out"
            >
              <span v-if="loading" class="absolute left-0 inset-y-0 flex items-center pl-3">
                <font-awesome-icon icon="spinner" class="animate-spin h-5 w-5 text-cyan-500" />
              </span>
              <span v-else class="absolute left-0 inset-y-0 flex items-center pl-3">
                <font-awesome-icon icon="lock" class="h-5 w-5 text-cyan-500 group-hover:text-cyan-400" />
              </span>
              {{ loading ? 'Signing in...' : 'Sign in' }}
            </button>
          </div>
        </div>
      </form>
    </div>
  </div>
</template>

<script>
export default {
  name: 'LoginPage',
  props: {
    loginForm: {
      type: Object,
      required: true
    },
    error: {
      type: String,
      default: ''
    },
    loading: {
      type: Boolean,
      default: false
    }
  }
}
</script>

