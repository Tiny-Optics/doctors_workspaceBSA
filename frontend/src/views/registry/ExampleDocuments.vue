<template>
  <div class="min-h-screen bg-bloodsa-red">
    <!-- Header Section -->
    <section class="relative overflow-hidden">
      <div class="relative max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-12 md:py-20">
        <!-- Back Button -->
        <button 
          @click="$router.push({ name: 'registry' })"
          class="mb-6 inline-flex items-center text-white hover:text-white/80 transition-colors"
        >
          <svg class="w-5 h-5 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7" />
          </svg>
          Back
        </button>
        
        <div class="text-center">
          <h1 class="text-4xl md:text-5xl font-bold text-white mb-4">Example Documents</h1>
          <p class="text-white max-w-2xl mx-auto">
            Download example documents and templates for the African HOPeR Registry submissions.
          </p>

          <!-- Search -->
          <div class="mt-6 max-w-xl mx-auto">
            <div class="relative">
              <input
                v-model="searchQuery"
                type="text"
                placeholder="Search documents..."
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
          <p class="text-gray-600 mt-4">Loading documents...</p>
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
            @click="loadDocuments"
            class="mt-4 px-4 py-2 bg-bloodsa-red text-white rounded-lg hover:bg-red-700 transition-colors"
          >
            Try Again
          </button>
        </div>
      </div>
    </section>

    <!-- Empty State -->
    <section v-else-if="!loading && filteredDocuments.length === 0" class="py-8">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div class="bg-white rounded-xl shadow-lg p-12 text-center">
          <svg class="w-16 h-16 text-gray-400 mx-auto mb-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M7 21h10a2 2 0 002-2V9.414a1 1 0 00-.293-.707l-5.414-5.414A1 1 0 0012.586 3H7a2 2 0 00-2 2v14a2 2 0 002 2z" />
          </svg>
          <h3 class="text-xl font-semibold text-gray-900 mb-2">
            {{ searchQuery ? 'No documents match your search' : 'No example documents available' }}
          </h3>
          <p class="text-gray-600 mb-4">
            {{ searchQuery ? 'Try a different search term' : 'Example documents will be available soon' }}
          </p>
          <p v-if="!searchQuery && documentsPath" class="text-sm text-gray-500">
            Configured path: <code class="bg-gray-100 px-2 py-1 rounded">{{ documentsPath }}</code>
          </p>
        </div>
      </div>
    </section>

    <!-- Documents List -->
    <section v-else class="py-8">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div class="bg-white rounded-xl shadow-lg p-6">
          <p class="italic text-gray-500 mb-4">Click to Download Document</p>
          <div class="space-y-2">
            <template v-for="item in filteredDocuments" :key="item.path">
              <DocumentItem 
                :item="item"
                :downloading-file="downloadingFile"
                @download="downloadFile"
              />
            </template>
          </div>
        </div>
      </div>
    </section>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { registryService } from '@/services/registryService'
import DocumentItem from '@/components/DocumentItem.vue'

interface DropboxFile {
  name: string
  path: string
  size: number
  modifiedTime: string
  isFolder: boolean
  children?: DropboxFile[]
}

const documents = ref<DropboxFile[]>([])
const documentsPath = ref<string>('')
const loading = ref(false)
const error = ref<string | null>(null)
const searchQuery = ref('')
const downloadingFile = ref<string | null>(null)

// Recursive search function for nested structure
function searchInDocuments(documents: DropboxFile[], query: string): DropboxFile[] {
  if (!query) return documents
  
  const results: DropboxFile[] = []
  const lowerQuery = query.toLowerCase()
  
  for (const doc of documents) {
    if (doc.isFolder && doc.children) {
      // Search in folder contents
      const matchingChildren = searchInDocuments(doc.children, query)
      if (matchingChildren.length > 0) {
        // Include folder if it has matching children
        results.push({
          ...doc,
          children: matchingChildren
        })
      }
    } else if (!doc.isFolder && doc.name.toLowerCase().includes(lowerQuery)) {
      // Include file if it matches search
      results.push(doc)
    }
  }
  
  return results
}

const filteredDocuments = computed(() => {
  if (!searchQuery.value) {
    return documents.value
  }
  return searchInDocuments(documents.value, searchQuery.value)
})

async function loadDocuments() {
  loading.value = true
  error.value = null
  
  try {
    // Get example documents from backend
    const response = await registryService.getExampleDocuments()
    
    // Handle null, undefined, or empty responses
    if (!response.files || !Array.isArray(response.files)) {
      documents.value = []
    } else {
      // Now we accept the full nested structure from the backend
      documents.value = response.files
    }
    
  } catch (err: any) {
    console.error('Failed to load documents:', err)
    error.value = err.message || 'Failed to load example documents'
    documents.value = []
  } finally {
    loading.value = false
  }
}

async function downloadFile(file: DropboxFile) {
  if (downloadingFile.value) return
  
  downloadingFile.value = file.path
  
  try {
    const response = await registryService.getDocumentDownloadLink(file.path)
    // Open in new tab for download
    window.open(response.link, '_blank')
  } catch (err: any) {
    console.error('Failed to download file:', err)
    error.value = err.message || 'Failed to download file'
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
  loadDocuments()
})
</script>

<style scoped>
.bg-bloodsa-red {
  background-color: #8B0000;
}
.text-bloodsa-red {
  color: #8B0000;
}
</style>

