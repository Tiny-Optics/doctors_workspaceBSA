<template>
  <div class="min-h-screen bg-gradient-to-br from-gray-50 to-gray-100 flex items-center justify-center px-4">
    <div class="max-w-md w-full text-center">
      <!-- Icon -->
      <div class="mb-8">
        <div class="inline-flex items-center justify-center w-24 h-24 bg-bloodsa-red bg-opacity-10 rounded-full">
          <svg class="w-12 h-12 text-bloodsa-red" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" />
          </svg>
        </div>
      </div>

      <!-- Title -->
      <h1 class="text-4xl font-bold text-gray-900 mb-4">
        Coming Soon
      </h1>

      <!-- Description -->
      <p class="text-lg text-gray-600 mb-8">
        This feature is currently under development. We're working hard to bring you the best experience.
      </p>

      <!-- Back Button -->
      <router-link
        to="/dashboard"
        class="inline-flex items-center px-6 py-3 bg-bloodsa-red text-white rounded-lg hover:bg-opacity-90 transition-all duration-200 font-medium"
      >
        <svg class="w-5 h-5 mr-2" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 19l-7-7m0 0l7-7m-7 7h18" />
        </svg>
        Back to Dashboard
      </router-link>

      <!-- Progress Info -->
      <div class="mt-12 p-6 bg-white rounded-lg shadow-md">
        <h3 class="font-semibold text-gray-900 mb-2">Development Status</h3>
        <p class="text-sm text-gray-600">
          This feature is part of our upcoming release. Stay tuned for updates!
        </p>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
// Coming soon page for features under development
</script>

