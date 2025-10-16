<template>
  <header class="bg-white shadow-sm border-b border-gray-200">
    <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
      <div class="flex justify-between items-center h-16">
        <!-- Logo and Title -->
        <div class="flex items-center space-x-3">
          <!-- Logo with red droplet and cross -->
          <div class="relative">
            <div class="w-10 h-10 bg-bloodsa-red rounded-full flex items-center justify-center">
              <svg class="w-6 h-6 text-white" fill="currentColor" viewBox="0 0 24 24">
                <path d="M12 2C6.48 2 2 6.48 2 12s4.48 10 10 10 10-4.48 10-10S17.52 2 12 2zm-2 15l-5-5 1.41-1.41L10 14.17l7.59-7.59L19 8l-9 9z"/>
              </svg>
            </div>
          </div>
          <h1 class="text-xl font-bold text-bloodsa-red">Doctor's Workspace</h1>
        </div>

        <!-- Navigation -->
        <nav class="flex items-center space-x-8">
          <router-link 
            to="/" 
            class="text-bloodsa-red font-medium hover:text-bloodsa-red-light-red transition-colors duration-200"
            :class="{ 'font-semibold': $route.path === '/' }"
          >
            HOME
          </router-link>
          
          <!-- Login Button -->
          <button class="btn-primary">
            Login
          </button>
        </nav>
      </div>
    </div>
  </header>
</template>

<script setup lang="ts">
// Component logic can be added here
</script>
