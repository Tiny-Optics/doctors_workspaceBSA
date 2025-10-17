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

    <!-- Loading State -->
    <section v-if="loading" class="py-12">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div class="text-center py-12">
          <div class="animate-spin rounded-full h-16 w-16 border-b-4 border-white mx-auto"></div>
          <p class="text-white mt-4 text-lg">Loading categories...</p>
        </div>
      </div>
    </section>

    <!-- Empty State -->
    <section v-else-if="!loading && categories.length === 0" class="py-12">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div class="bg-white rounded-xl shadow-lg p-12 text-center">
          <svg class="w-20 h-20 text-gray-400 mx-auto mb-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
          </svg>
          <h3 class="text-2xl font-semibold text-gray-900 mb-2">No SOP Categories Available</h3>
          <p class="text-gray-600 mb-4">Please check back later or contact your administrator.</p>
        </div>
      </div>
    </section>

    <!-- Cards Grid on Red Background with White Cards -->
    <section v-else class="py-12">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-10">
          <router-link
            v-for="category in sortedCategories"
            :key="category.id"
            :to="{ name: 'sops-list', params: { slug: category.slug } }"
            class="group bg-white rounded-xl shadow-lg hover:shadow-2xl transition-all duration-300 overflow-hidden border-2 border-gray-100 hover:border-bloodsa-red transform hover:-translate-y-1"
          >
            <div class="p-3">
              <img 
                v-if="category.imagePath" 
                class="w-full h-64 object-cover rounded-lg" 
                :src="`http://localhost:8080${category.imagePath}`" 
                :alt="category.name"
                @error="handleImageError"
              />
              <div 
                v-else
                class="w-full h-64 bg-gradient-to-br from-gray-200 to-gray-300 rounded-lg flex items-center justify-center"
              >
                <svg class="w-20 h-20 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
                </svg>
              </div>
            </div>
            <div class="px-6 pb-6">
              <h3 class="text-lg md:text-xl font-semibold text-gray-800 text-center mb-2">
                {{ category.name }}
              </h3>
              <p v-if="category.description" class="text-sm text-gray-600 text-center line-clamp-2">
                {{ category.description }}
              </p>
            </div>
          </router-link>
        </div>
      </div>
    </section>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { sopService } from '@/services/sopService'
import type { SOPCategory } from '@/types/sop'

const categories = ref<SOPCategory[]>([])
const loading = ref(false)

const sortedCategories = computed(() => {
  return [...categories.value].sort((a, b) => a.displayOrder - b.displayOrder)
})

async function loadCategories() {
  loading.value = true
  try {
    const response = await sopService.listCategories({ limit: 100 })
    // Only show active categories to regular users
    categories.value = response.categories.filter(c => c.isActive)
  } catch (error) {
    console.error('Failed to load SOP categories:', error)
  } finally {
    loading.value = false
  }
}

function handleImageError(event: Event) {
  const img = event.target as HTMLImageElement
  // Hide the image on error - the fallback div will show instead
  img.style.display = 'none'
}

onMounted(() => {
  loadCategories()
})
</script>

<style scoped>
.line-clamp-2 {
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}
</style>
