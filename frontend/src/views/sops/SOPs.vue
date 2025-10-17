<template>
  <div class="min-h-screen bg-bloodsa-red">
    <!-- Header Section on Red Background -->
    <section class="relative overflow-hidden">
      <div class="relative max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-12 md:py-20">
        <!-- Back Button -->
        <button 
          @click="$router.go(-1)"
          class="mb-6 inline-flex items-center text-white hover:text-white/80 transition-colors"
        >
          <svg class="w-5 h-5 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7" />
          </svg>
          Back
        </button>
        
        <div class="text-center">
          <h1 class="text-4xl md:text-5xl font-bold text-white mb-4">SOPs</h1>
          <p class="text-white max-w-2xl mx-auto">
            Browse sample SOPs by disease area. Open a category to view and download ready-to-customise PDF templates.
          </p>
        </div>
      </div>
    </section>

    <!-- Cards Grid on Red Background with White Cards -->
    <section class="py-12">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-10">
          <!-- Card 1 -->
          <router-link
            :to="{ name: 'sops-list', params: { disease: 'anaemia' } }"
            class="group bg-white rounded-xl shadow-lg hover:shadow-2xl transition-all duration-300 overflow-hidden border-2 border-gray-100 hover:border-bloodsa-red transform hover:-translate-y-1"
          >
            <div class="p-3">
              <img class="w-full h-64 object-cover rounded-lg" src="https://images.unsplash.com/photo-1612538497338-6dff39a6064e?q=80&w=1400&auto=format&fit=crop" alt="Anaemia" />
            </div>
            <div class="px-6 pb-6">
              <h3 class="text-lg md:text-xl font-semibold text-gray-800 text-center">Anaemia</h3>
            </div>
          </router-link>

          <!-- Card 2 -->
          <router-link
            :to="{ name: 'sops-list', params: { disease: 'myeloma' } }"
            class="group bg-white rounded-xl shadow-lg hover:shadow-2xl transition-all duration-300 overflow-hidden border-2 border-gray-100 hover:border-bloodsa-red transform hover:-translate-y-1"
          >
            <div class="p-3">
              <img class="w-full h-64 object-cover rounded-lg" src="https://images.unsplash.com/photo-1559757175-08fda9a4a273?q=80&w=1400&auto=format&fit=crop" alt="Myeloma" />
            </div>
            <div class="px-6 pb-6">
              <h3 class="text-lg md:text-xl font-semibold text-gray-800 text-center">Myeloma</h3>
            </div>
          </router-link>

          <!-- Card 3 -->
          <router-link
            :to="{ name: 'sops-list', params: { disease: 'lymphoma' } }"
            class="group bg-white rounded-xl shadow-lg hover:shadow-2xl transition-all duration-300 overflow-hidden border-2 border-gray-100 hover:border-bloodsa-red transform hover:-translate-y-1"
          >
            <div class="p-3">
              <img class="w-full h-64 object-cover rounded-lg" src="https://images.unsplash.com/photo-1583912267550-d9f3e83d4a53?q=80&w=1400&auto=format&fit=crop" alt="Lymphoma" />
            </div>
            <div class="px-6 pb-6">
              <h3 class="text-lg md:text-xl font-semibold text-gray-800 text-center">Lymphoma</h3>
            </div>
          </router-link>
        </div>
      </div>
    </section>
  </div>
</template>

<script setup lang="ts">
</script>

<style scoped>
</style>


