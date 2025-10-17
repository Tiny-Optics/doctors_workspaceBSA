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
          @click="$router.push({ name: 'sops' })"
          class="mb-6 inline-flex items-center text-white hover:text-white/80 transition-colors"
        >
          <svg class="w-5 h-5 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7" />
          </svg>
          Back
        </button>
        
        <div class="text-center">
          <h1 class="text-4xl md:text-5xl font-bold text-white mb-4">{{ categoryName }}</h1>
          <p class="text-white max-w-2xl mx-auto">
            {{ categoryDescription || `Browse and download SOPs for ${categoryName}. Click any item to download the PDF template.` }}
          </p>

          <!-- Search -->
          <div class="mt-6 max-w-xl mx-auto">
            <div class="relative">
              <input
                v-model="query"
                type="text"
                placeholder="Search files..."
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

    <!-- Loading State -->
    <section v-if="loading" class="py-8">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div class="bg-white rounded-xl shadow-lg p-12 text-center">
          <div class="animate-spin rounded-full h-12 w-12 border-b-2 border-bloodsa-red mx-auto"></div>
          <p class="text-gray-600 mt-4">Loading files...</p>
        </div>
      </div>
    </section>

    <!-- Error State -->
    <section v-else-if="error" class="py-8">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div class="bg-white rounded-xl shadow-lg p-12 text-center">
          <svg class="w-16 h-16 text-red-500 mx-auto mb-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
          </svg>
          <h3 class="text-xl font-semibold text-gray-900 mb-2">{{ error }}</h3>
          <button
            @click="loadData"
            class="mt-4 px-4 py-2 bg-bloodsa-red text-white rounded-lg hover:bg-red-700 transition-colors"
          >
            Try Again
          </button>
        </div>
      </div>
    </section>

    <!-- Empty State -->
    <section v-else-if="!loading && filtered.length === 0" class="py-8">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div class="bg-white rounded-xl shadow-lg p-12 text-center">
          <svg class="w-16 h-16 text-gray-400 mx-auto mb-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M7 21h10a2 2 0 002-2V9.414a1 1 0 00-.293-.707l-5.414-5.414A1 1 0 0012.586 3H7a2 2 0 00-2 2v14a2 2 0 002 2z" />
          </svg>
          <h3 class="text-xl font-semibold text-gray-900 mb-2">
            {{ query ? 'No files match your search' : 'No files in this category' }}
          </h3>
          <p class="text-gray-600 mb-4">
            {{ query ? 'Try a different search term' : 'Files need to be uploaded manually to Dropbox' }}
          </p>
          <p v-if="!query" class="text-sm text-gray-500">
            Upload files to: <code class="bg-gray-100 px-2 py-1 rounded">{{ categoryDropboxPath }}</code>
          </p>
        </div>
      </div>
    </section>

    <!-- SOP List -->
    <section v-else class="py-8">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div class="bg-white rounded-xl shadow-lg p-6">
          <p class="italic text-gray-500 mb-4">Click to Download PDF</p>
          <ul class="divide-y divide-gray-200">
            <li v-for="file in filtered" :key="file.path" class="py-4">
              <div class="flex items-center justify-between">
                <div class="flex items-center gap-3 flex-1">
                  <svg class="w-8 h-8 text-red-500 flex-shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M7 21h10a2 2 0 002-2V9.414a1 1 0 00-.293-.707l-5.414-5.414A1 1 0 0012.586 3H7a2 2 0 00-2 2v14a2 2 0 002 2z" />
                  </svg>
                  <div class="flex-1 min-w-0">
                    <button 
                      @click="downloadFile(file)" 
                      class="text-xl font-semibold text-gray-700 hover:text-bloodsa-red transition-colors text-left truncate block w-full"
                      :class="{ 'opacity-50 cursor-wait': downloadingFile === file.path }"
                      :disabled="downloadingFile === file.path"
                    >
                      {{ file.name }}
                    </button>
                    <p class="text-sm text-gray-500 mt-1">
                      {{ formatFileSize(file.size) }} • {{ formatDate(file.modifiedTime) }}
                    </p>
                  </div>
                </div>
                <button 
                  @click="downloadFile(file)"
                  :disabled="downloadingFile === file.path"
                  class="ml-4 px-4 py-2 bg-bloodsa-red text-white rounded-lg hover:bg-red-700 transition-colors disabled:opacity-50 disabled:cursor-wait flex items-center gap-2"
                >
                  <svg v-if="downloadingFile !== file.path" class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 10v6m0 0l-3-3m3 3l3-3m2 8H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
                  </svg>
                  <svg v-else class="animate-spin w-5 h-5" fill="none" viewBox="0 0 24 24">
                    <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                    <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                  </svg>
                  {{ downloadingFile === file.path ? 'Downloading...' : 'Download' }}
                </button>
              </div>
            </li>
          </ul>
        </div>
      </div>
    </section>
  </div>
</template>

<script setup lang="ts">
import { computed, ref, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { sopService } from '@/services/sopService'
import type { SOPCategory, SOPFile } from '@/types/sop'

const route = useRoute()
const slug = computed(() => route.params.slug as string)

const category = ref<SOPCategory | null>(null)
const files = ref<SOPFile[]>([])
const loading = ref(false)
const error = ref<string | null>(null)
const query = ref('')
const downloadingFile = ref<string | null>(null)

const categoryName = computed(() => category.value?.name || 'SOPs')
const categoryDescription = computed(() => category.value?.description)
const categoryDropboxPath = computed(() => category.value?.dropboxPath || '')

const filtered = computed(() =>
  files.value.filter(file => 
    file.name.toLowerCase().includes(query.value.toLowerCase())
  )
)

async function loadData() {
  loading.value = true
  error.value = null
  
  try {
    // Get category by slug
    const categoryData = await sopService.getCategoryBySlug(slug.value)
    
    if (!categoryData) {
      error.value = 'Category not found'
      return
    }
    
    category.value = categoryData
    
    // Get files for this category
    const fileData = await sopService.getCategoryFiles(categoryData.id)
    
    // Handle null, undefined, or empty responses
    if (!fileData || !Array.isArray(fileData)) {
      files.value = []
    } else {
      files.value = fileData.filter(f => !f.isFolder) // Only show files, not folders
    }
    
  } catch (err: any) {
    console.error('Failed to load data:', err)
    error.value = err.message || 'Failed to load category or files'
  } finally {
    loading.value = false
  }
}

async function downloadFile(file: SOPFile) {
  if (!category.value || downloadingFile.value) return
  
  downloadingFile.value = file.path
  
  try {
    const downloadLink = await sopService.getFileDownloadLink(category.value.id, file.name)
    // Open in new tab for download
    window.open(downloadLink, '_blank')
  } catch (err: any) {
    console.error('Failed to download file:', err)
    alert(err.message || 'Failed to download file')
  } finally {
    downloadingFile.value = null
  }
}

function formatFileSize(bytes: number): string {
  if (bytes === 0) return '0 Bytes'
  const k = 1024
  const sizes = ['Bytes', 'KB', 'MB', 'GB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return Math.round(bytes / Math.pow(k, i) * 100) / 100 + ' ' + sizes[i]
}

function formatDate(dateString: string): string {
  const date = new Date(dateString)
  return date.toLocaleDateString('en-US', { year: 'numeric', month: 'short', day: 'numeric' })
}

onMounted(() => {
  loadData()
})
</script>

<style scoped>
</style>
