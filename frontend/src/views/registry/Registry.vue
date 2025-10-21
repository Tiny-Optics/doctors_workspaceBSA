<template>
  <div class="min-h-screen bg-bloodsa-red">
    <!-- Header Section with Video -->
    <section class="relative overflow-hidden">
      <div class="relative max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-12 md:py-20">
        <!-- Back Button -->
        <button 
          @click="$router.push({ name: 'dashboard' })"
          class="mb-6 inline-flex items-center text-white hover:text-white/80 transition-colors"
        >
          <svg class="w-5 h-5 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7" />
          </svg>
          Back
        </button>
        
        <div class="text-center">
          <h1 class="text-4xl md:text-5xl font-bold text-white mb-4">African HOPeR</h1>
          <p class="text-white max-w-2xl mx-auto mb-8">
            A space for our registry.
          </p>
          
          <!-- Video Section -->
          <div v-if="registryConfig?.videoUrl" class="max-w-4xl mx-auto">
            <div class="bg-white rounded-xl shadow-lg p-8">
              <div class="flex items-center justify-center mb-6">
                <div class="flex items-center">
                  <div class="relative">
                    <!-- African H PeR Logo with Continent Icon -->
                    <div class="flex items-center">
                      <span class="text-6xl font-bold text-bloodsa-red mr-2">A</span>
                      <div class="relative">
                        <span class="text-6xl font-bold text-bloodsa-red">H</span>
                        <div class="absolute -top-2 -right-2 w-8 h-8 bg-gray-300 rounded-full flex items-center justify-center">
                          <svg class="w-6 h-6 text-bloodsa-red" viewBox="0 0 24 24" fill="currentColor">
                            <path d="M12 2C6.48 2 2 6.48 2 12s4.48 10 10 10 10-4.48 10-10S17.52 2 12 2zm-1 17.93c-3.94-.49-7-3.85-7-7.93 0-.62.08-1.21.21-1.79L9 15v1c0 1.1.9 2 2 2v1.93zm6.9-2.54c-.26-.81-1-1.39-1.9-1.39h-1v-3c0-.55-.45-1-1-1H8v-2h2c.55 0 1-.45 1-1V7h2c1.1 0 2-.9 2-2v-.41c2.93 1.19 5 4.06 5 7.41 0 2.08-.8 3.97-2.1 5.39z"/>
                          </svg>
                        </div>
                      </div>
                      <span class="text-6xl font-bold text-bloodsa-red ml-2">PeR</span>
                    </div>
                    
                    <!-- YouTube Play Button Overlay -->
                    <div class="absolute inset-0 flex items-center justify-center">
                      <button 
                        @click="playVideo"
                        class="w-20 h-20 bg-red-600 hover:bg-red-700 rounded-full flex items-center justify-center shadow-lg transform hover:scale-110 transition-all duration-200"
                      >
                        <svg class="w-8 h-8 text-white ml-1" fill="currentColor" viewBox="0 0 24 24">
                          <path d="M8 5v14l11-7z"/>
                        </svg>
                      </button>
                    </div>
                  </div>
                </div>
              </div>
              
              <p class="text-bloodsa-red text-lg font-medium text-center">
                African Haematology Oncology Patient Electronic Registry
              </p>
            </div>
          </div>
          
          <!-- No Video State -->
          <div v-else class="max-w-4xl mx-auto">
            <div class="bg-white rounded-xl shadow-lg p-8">
              <div class="flex items-center justify-center mb-6">
                <div class="flex items-center">
                  <span class="text-6xl font-bold text-bloodsa-red mr-2">A</span>
                  <div class="relative">
                    <span class="text-6xl font-bold text-bloodsa-red">H</span>
                    <div class="absolute -top-2 -right-2 w-8 h-8 bg-gray-300 rounded-full flex items-center justify-center">
                      <svg class="w-6 h-6 text-bloodsa-red" viewBox="0 0 24 24" fill="currentColor">
                        <path d="M12 2C6.48 2 2 6.48 2 12s4.48 10 10 10 10-4.48 10-10S17.52 2 12 2zm-1 17.93c-3.94-.49-7-3.85-7-7.93 0-.62.08-1.21.21-1.79L9 15v1c0 1.1.9 2 2 2v1.93zm6.9-2.54c-.26-.81-1-1.39-1.9-1.39h-1v-3c0-.55-.45-1-1-1H8v-2h2c.55 0 1-.45 1-1V7h2c1.1 0 2-.9 2-2v-.41c2.93 1.19 5 4.06 5 7.41 0 2.08-.8 3.97-2.1 5.39z"/>
                      </svg>
                    </div>
                  </div>
                  <span class="text-6xl font-bold text-bloodsa-red ml-2">PeR</span>
                </div>
              </div>
              
              <p class="text-bloodsa-red text-lg font-medium text-center">
                African Haematology Oncology Patient Electronic Registry
              </p>
            </div>
          </div>
        </div>
      </div>
    </section>

    <!-- Three Cards Section -->
    <section class="py-12">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div class="grid grid-cols-1 md:grid-cols-3 gap-10">
          <!-- Process Outline Card -->
          <div class="group bg-white rounded-xl shadow-lg hover:shadow-2xl transition-all duration-300 overflow-hidden border-2 border-gray-100 hover:border-bloodsa-red transform hover:-translate-y-1 cursor-pointer">
            <div class="p-6">
              <div class="w-full h-48 bg-gradient-to-br from-yellow-200 to-yellow-300 rounded-lg flex items-center justify-center mb-6">
                <div class="text-center">
                  <svg class="w-16 h-16 text-yellow-600 mx-auto mb-2" fill="currentColor" viewBox="0 0 24 24">
                    <path d="M12 2C6.48 2 2 6.48 2 12s4.48 10 10 10 10-4.48 10-10S17.52 2 12 2zm-1 17.93c-3.94-.49-7-3.85-7-7.93 0-.62.08-1.21.21-1.79L9 15v1c0 1.1.9 2 2 2v1.93zm6.9-2.54c-.26-.81-1-1.39-1.9-1.39h-1v-3c0-.55-.45-1-1-1H8v-2h2c.55 0 1-.45 1-1V7h2c1.1 0 2-.9 2-2v-.41c2.93 1.19 5 4.06 5 7.41 0 2.08-.8 3.97-2.1 5.39z"/>
                  </svg>
                  <div class="flex justify-center space-x-1">
                    <div class="w-2 h-2 bg-yellow-600 rounded-full"></div>
                    <div class="w-2 h-2 bg-yellow-500 rounded-full"></div>
                    <div class="w-2 h-2 bg-yellow-400 rounded-full"></div>
                  </div>
                </div>
              </div>
              <h3 class="text-xl font-semibold text-gray-800 text-center">
                Process Outline
              </h3>
            </div>
          </div>

          <!-- Example Documents Card -->
          <div class="group bg-white rounded-xl shadow-lg hover:shadow-2xl transition-all duration-300 overflow-hidden border-2 border-gray-100 hover:border-bloodsa-red transform hover:-translate-y-1 cursor-pointer">
            <div class="p-6">
              <div class="w-full h-48 bg-gradient-to-br from-blue-200 to-blue-300 rounded-lg flex items-center justify-center mb-6">
                <div class="text-center">
                  <svg class="w-16 h-16 text-blue-600 mx-auto mb-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
                  </svg>
                  <div class="flex justify-center space-x-1">
                    <div class="w-2 h-2 bg-blue-600 rounded-full"></div>
                    <div class="w-2 h-2 bg-blue-500 rounded-full"></div>
                    <div class="w-2 h-2 bg-blue-400 rounded-full"></div>
                  </div>
                </div>
              </div>
              <h3 class="text-xl font-semibold text-gray-800 text-center">
                Example Documents
              </h3>
            </div>
          </div>

          <!-- Upload Final Approvals Card -->
          <div class="group bg-white rounded-xl shadow-lg hover:shadow-2xl transition-all duration-300 overflow-hidden border-2 border-gray-100 hover:border-bloodsa-red transform hover:-translate-y-1 cursor-pointer">
            <div class="p-6">
              <div class="w-full h-48 bg-gradient-to-br from-green-200 to-green-300 rounded-lg flex items-center justify-center mb-6">
                <div class="text-center">
                  <svg class="w-16 h-16 text-green-600 mx-auto mb-2" fill="currentColor" viewBox="0 0 24 24">
                    <path d="M19.35 10.04C18.67 6.59 15.64 4 12 4 9.11 4 6.6 5.64 5.35 8.04 2.34 8.36 0 10.91 0 14c0 3.31 2.69 6 6 6h13c2.76 0 5-2.24 5-5 0-2.64-2.05-4.78-4.65-4.96zM14 13v4h-4v-4H7l5-5 5 5h-3z"/>
                  </svg>
                  <div class="flex justify-center space-x-1">
                    <div class="w-2 h-2 bg-green-600 rounded-full"></div>
                    <div class="w-2 h-2 bg-green-500 rounded-full"></div>
                    <div class="w-2 h-2 bg-green-400 rounded-full"></div>
                  </div>
                </div>
              </div>
              <h3 class="text-xl font-semibold text-gray-800 text-center">
                Upload Final Approvals
              </h3>
            </div>
          </div>
        </div>
      </div>
    </section>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { registryService } from '@/services/registryService'

const registryConfig = ref<any>(null)

async function loadRegistryConfig() {
  try {
    const config = await registryService.getPublicConfiguration()
    registryConfig.value = config
  } catch (error) {
    console.error('Failed to load registry configuration:', error)
  }
}

function playVideo() {
  if (registryConfig.value?.videoUrl) {
    window.open(registryConfig.value.videoUrl, '_blank')
  }
}

onMounted(() => {
  loadRegistryConfig()
})
</script>

<style scoped>
.bg-bloodsa-red {
  background-color: #8B0000;
}
.text-bloodsa-red {
  color: #8B0000;
}
.border-bloodsa-red {
  border-color: #8B0000;
}
</style>
