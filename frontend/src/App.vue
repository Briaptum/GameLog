<template>
  <div class="min-h-screen bg-brand-black flex flex-col items-center justify-center">
    <h1 class="text-brand-red text-6xl font-bold tracking-wide mb-8">GameLog</h1>
    <div class="w-full max-w-2xl px-4">
      <div class="relative">
        <input
          type="text"
          placeholder="Search..."
          class="w-full px-6 py-4 pr-24 bg-brand-black border-2 border-brand-red rounded-lg text-white placeholder-gray-500 focus:outline-none focus:border-brand-dark-red transition-colors"
        />
        <button
          type="button"
          class="absolute right-2 top-1/2 -translate-y-1/2 px-6 py-2 bg-brand-red text-white rounded-md hover:bg-brand-dark-red transition-colors font-medium"
        >
          Search
        </button>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  name: 'App'
}
</script>
