<template>
  <div class="min-h-screen bg-bloodsa-red">
    <!-- Header Section -->
    <section class="relative overflow-hidden">
      <div class="absolute inset-0 opacity-5">
        <div class="absolute inset-0" style="background-image: repeating-linear-gradient(45deg, #8B0000 0, #8B0000 1px, transparent 0, transparent 50%); background-size: 10px 10px;"></div>
      </div>

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
          <h1 class="text-4xl md:text-5xl font-bold text-white mb-4">{{ pageTitle }}</h1>
          <p class="text-white max-w-2xl mx-auto">
            Browse and download SOPs for {{ diseaseTitle }}. Click any item to download the PDF template.
          </p>

          <!-- Search -->
          <div class="mt-6 max-w-xl mx-auto">
            <div class="relative">
              <input
                v-model="query"
                type="text"
                placeholder="Search..."
                class="w-full pl-10 pr-4 py-2 border border-gray-300 rounded-full focus:ring-2 focus:ring-bloodsa-red focus:border-transparent"
              />
              <svg class="absolute left-3 top-2.5 w-5 h-5 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
              </svg>
            </div>
          </div>
        </div>
      </div>
    </section>

    <!-- SOP List -->
    <section class="py-8">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div class="bg-white rounded-xl shadow-lg p-6">
          <p class="italic text-gray-500 mb-4">Click to Download PDF</p>
          <ul class="divide-y divide-gray-200">
            <li v-for="item in filtered" :key="item" class="py-4">
              <a href="#" class="text-xl font-semibold text-gray-700 hover:text-bloodsa-red transition-colors">{{ item }}</a>
            </li>
          </ul>
        </div>
      </div>
    </section>
  </div>
</template>

<script setup lang="ts">
import { computed, ref } from 'vue'
import { useRoute } from 'vue-router'

const route = useRoute()
const disease = computed(() => (route.params.disease as string) || 'anaemia')
const category = computed(() => (route.params.category as string) || 'unit-management')

const diseaseTitle = computed(() => {
  const d = disease.value
  if (d === 'anaemia') return 'Anaemia'
  if (d === 'myeloma') return 'Myeloma'
  if (d === 'lymphoma') return 'Lymphoma'
  return d.charAt(0).toUpperCase() + d.slice(1)
})

const pageTitle = computed(() => {
  const c = category.value
  if (c === 'unit-management') return 'Anaemia'
  if (c === 'treatment') return diseaseTitle.value
  return diseaseTitle.value
})

const query = ref('')
const items = ref<string[]>([
  'Aplastic Anaemia',
  'Iron deficiency Anaemia',
  'B12 Deficiency',
  'Something Else',
  'Something Else',
  'Something Else',
  'Something Else',
  'Something Else',
  'Something Else',
  'Something Else',
  'Something Else',
  'Something Else'
])

const filtered = computed(() =>
  items.value.filter(i => i.toLowerCase().includes(query.value.toLowerCase()))
)
</script>

<style scoped>
</style>


