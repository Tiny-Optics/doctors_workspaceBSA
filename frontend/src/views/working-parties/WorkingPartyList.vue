<template>
  <div class="min-h-screen bg-bloodsa-red">
    <section class="relative overflow-hidden">
      <div class="absolute inset-0 opacity-5">
        <div class="absolute inset-0" style="background-image: repeating-linear-gradient(45deg, #8B0000 0, #8B0000 1px, transparent 0, transparent 50%); background-size: 10px 10px;"></div>
      </div>

      <div class="relative max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-12 md:py-20">
        <button
          @click="$router.push({ name: 'working-parties' })"
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
            {{ categoryDescription || `Browse and download documents for ${categoryName}. Click any item to download.` }}
          </p>

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

    <section v-if="loading" class="py-8">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div class="bg-white rounded-xl shadow-lg p-12 text-center">
          <div class="animate-spin rounded-full h-12 w-12 border-b-2 border-bloodsa-red mx-auto"></div>
          <p class="text-gray-600 mt-4">Loading files...</p>
        </div>
      </div>
    </section>

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

    <section v-else class="py-8">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div class="bg-white rounded-xl shadow-lg p-6">
          <p class="italic text-gray-500 mb-4">Click to download</p>
          <div class="space-y-2">
            <template v-for="item in filtered" :key="item.path">
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
import { computed, ref, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { workingPartyService } from '@/services/workingPartyService'
import type { WorkingPartyCategory, WorkingPartyFile } from '@/types/workingParty'
import DocumentItem from '@/components/DocumentItem.vue'

const route = useRoute()
const slug = computed(() => route.params.slug as string)

const category = ref<WorkingPartyCategory | null>(null)
const files = ref<WorkingPartyFile[]>([])
const loading = ref(false)
const error = ref<string | null>(null)
const query = ref('')
const downloadingFile = ref<string | null>(null)

const categoryName = computed(() => category.value?.name || 'Working Parties')
const categoryDescription = computed(() => category.value?.description)
const categoryDropboxPath = computed(() => {
  const path = category.value?.dropboxPath || ''
  try {
    return decodeURIComponent(path)
  } catch {
    return path
  }
})

function searchInFiles(fileList: WorkingPartyFile[], searchQuery: string): WorkingPartyFile[] {
  if (!searchQuery) return fileList

  const results: WorkingPartyFile[] = []
  const lowerQuery = searchQuery.toLowerCase()

  for (const file of fileList) {
    if (file.isFolder && file.children) {
      const matchingChildren = searchInFiles(file.children, searchQuery)
      if (matchingChildren.length > 0) {
        results.push({
          ...file,
          children: matchingChildren
        })
      }
    } else if (!file.isFolder && file.name.toLowerCase().includes(lowerQuery)) {
      results.push(file)
    }
  }

  return results
}

const filtered = computed(() => {
  if (!query.value) {
    return files.value
  }
  return searchInFiles(files.value, query.value)
})

async function loadData() {
  loading.value = true
  error.value = null

  try {
    const categoryData = await workingPartyService.getCategoryBySlug(slug.value)

    if (!categoryData) {
      error.value = 'Category not found'
      return
    }

    category.value = categoryData

    const fileData = await workingPartyService.getCategoryFiles(categoryData.id)

    if (!fileData || !Array.isArray(fileData)) {
      files.value = []
    } else {
      files.value = fileData
    }
  } catch (err: any) {
    console.error('Failed to load data:', err)
    error.value = err.message || 'Failed to load category or files'
  } finally {
    loading.value = false
  }
}

async function downloadFile(file: WorkingPartyFile) {
  if (!category.value || downloadingFile.value) return

  downloadingFile.value = file.path

  try {
    const downloadLink = await workingPartyService.getFileDownloadLink(category.value.id, file.path)
    window.open(downloadLink, '_blank')
  } catch (err: any) {
    console.error('Failed to download file:', err)
    alert(err.message || 'Failed to download file')
  } finally {
    downloadingFile.value = null
  }
}

onMounted(() => {
  loadData()
})
</script>
