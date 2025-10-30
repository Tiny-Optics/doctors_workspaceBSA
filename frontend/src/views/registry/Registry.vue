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
              <!-- Logo Header -->
              <div class="flex items-center justify-center mb-6">
                <img src="/AfricanHoperRegistry.svg" alt="African HOPeR" class="max-w-full h-20 md:h-24" />
              </div>
              
            

              <!-- Embedded YouTube Video -->
              <div class="aspect-video w-full rounded-lg overflow-hidden shadow-lg relative">
                <!-- Loading skeleton -->
                <div 
                  v-if="!videoLoaded"
                  class="absolute inset-0 bg-gray-200 animate-pulse flex items-center justify-center"
                >
                  <div class="text-center">
                    <svg class="w-16 h-16 text-gray-400 mx-auto mb-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M14.752 11.168l-3.197-2.132A1 1 0 0010 9.87v4.263a1 1 0 001.555.832l3.197-2.132a1 1 0 000-1.664z" />
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
                    </svg>
                    <p class="text-gray-500 text-sm">Loading video...</p>
                  </div>
                </div>
                <!-- Video iframe -->
                <iframe
                  :src="getEmbedUrl(registryConfig.videoUrl)"
                  class="w-full h-full"
                  frameborder="0"
                  allow="accelerometer; autoplay; clipboard-write; encrypted-media; gyroscope; picture-in-picture"
                  allowfullscreen
                  @load="videoLoaded = true"
                ></iframe>
              </div>
            </div>
          </div>
          
          <!-- No Video State -->
          <div v-else class="max-w-4xl mx-auto">
            <div class="bg-white rounded-xl shadow-lg p-8">
              <div class="flex items-center justify-center mb-6">
                <img src="/AfricanHoperRegistry.svg" alt="African HOPeR" class="max-w-full h-20 md:h-24" />
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
          <div @click="showProcessModal = true" class="group bg-white rounded-xl shadow-lg hover:shadow-2xl transition-all duration-300 overflow-hidden border-2 border-gray-100 hover:border-bloodsa-red transform hover:-translate-y-1 cursor-pointer">
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
          <div 
            @click="$router.push({ name: 'registry-example-documents' })"
            class="group bg-white rounded-xl shadow-lg hover:shadow-2xl transition-all duration-300 overflow-hidden border-2 border-gray-100 hover:border-bloodsa-red transform hover:-translate-y-1 cursor-pointer"
          >
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
          <div 
            @click="$router.push({ name: 'registry-upload' })"
            class="group bg-white rounded-xl shadow-lg hover:shadow-2xl transition-all duration-300 overflow-hidden border-2 border-gray-100 hover:border-bloodsa-red transform hover:-translate-y-1 cursor-pointer"
          >
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
    
    <!-- Modal for Process Outline PDFs -->
    <div v-if="showProcessModal" class="fixed inset-0 z-50">
      <!-- Backdrop -->
      <div class="absolute inset-0 bg-black/50" @click="showProcessModal = false"></div>

      <!-- Modal container -->
      <div class="absolute inset-0 flex items-center justify-center p-4" @click.self="showProcessModal = false">
        <div class="w-full max-w-6xl bg-white rounded-xl shadow-2xl overflow-hidden">
          <!-- Header -->
          <div class="flex items-center justify-between px-6 py-4 border-b">
            <h2 class="text-lg font-semibold text-gray-800">Process Outline</h2>
            <button class="text-gray-500 hover:text-gray-700" @click="showProcessModal = false" aria-label="Close">
              <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"/></svg>
            </button>
          </div>

          <!-- Tabs -->
          <div class="px-6 pt-4">
            <div class="flex space-x-2">
              <button
                class="px-3 py-2 rounded-md text-sm font-medium border"
                :class="activePdf === 'summary' ? 'bg-bloodsa-red text-white border-bloodsa-red' : 'bg-white text-gray-700 border-gray-300'"
                @click="activePdf = 'summary'"
              >Setup Summary</button>
              <button
                class="px-3 py-2 rounded-md text-sm font-medium border"
                :class="activePdf === 'training' ? 'bg-bloodsa-red text-white border-bloodsa-red' : 'bg-white text-gray-700 border-gray-300'"
                @click="activePdf = 'training'"
              >REDCap User Training</button>
            </div>
          </div>

          <!-- Body -->
          <div class="px-6 pb-6 pt-4">
            <div class="grid grid-cols-1 gap-4">
              <!-- PDF Viewer -->
              <div class="w-full h-[70vh] border rounded-lg overflow-hidden bg-gray-50">
                <iframe
                  v-if="activePdf === 'summary'"
                  :src="pdfSummaryUrl"
                  class="w-full h-full"
                  frameborder="0"
                  title="Setup Summary PDF"
                ></iframe>
                <iframe
                  v-else
                  :src="pdfTrainingUrl"
                  class="w-full h-full"
                  frameborder="0"
                  title="REDCap User Training PDF"
                ></iframe>
              </div>

              <!-- Download links -->
              <div class="flex items-center justify-between">
                <div class="text-sm text-gray-600">PDF opens inline; you can also download it.</div>
                <div class="space-x-2">
                  <a :href="pdfSummaryUrl" download class="px-3 py-2 text-sm rounded-md border border-gray-300 hover:bg-gray-100">Download Setup Summary</a>
                  <a :href="pdfTrainingUrl" download class="px-3 py-2 text-sm rounded-md border border-gray-300 hover:bg-gray-100">Download REDCap Training</a>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { registryService } from '@/services/registryService'

const registryConfig = ref<any>(null)
const videoLoaded = ref(false)
const showProcessModal = ref(false)
const activePdf = ref<'summary' | 'training'>('summary')

// Public assets are served from the root in Vite
const pdfSummaryUrl = '/AHoperSetupSummary.pdf'
// Handle space in filename using %20
const pdfTrainingUrl = '/RedcapUserTraining%20guide.pdf'

async function loadRegistryConfig() {
  try {
    const config = await registryService.getPublicConfiguration()
    registryConfig.value = config
  } catch (error) {
    console.error('Failed to load registry configuration:', error)
  }
}

function getEmbedUrl(url: string): string {
  // Convert various YouTube URL formats to embed URL
  // Supports:
  // - https://www.youtube.com/watch?v=VIDEO_ID
  // - https://youtu.be/VIDEO_ID
  // - https://www.youtube.com/embed/VIDEO_ID (already embed format)
  
  if (!url) return ''
  
  // If already an embed URL, return as is
  if (url.includes('/embed/')) {
    return url
  }
  
  // Extract video ID from different formats
  let videoId = ''
  
  if (url.includes('youtube.com/watch?v=')) {
    videoId = url.split('v=')[1]?.split('&')[0] || ''
  } else if (url.includes('youtu.be/')) {
    videoId = url.split('youtu.be/')[1]?.split('?')[0] || ''
  }
  
  if (videoId) {
    return `https://www.youtube.com/embed/${videoId}`
  }
  
  return url
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
