<template>
  <div class="border border-gray-200 rounded-lg">
    <!-- Folder or File Header -->
    <div 
      class="flex items-center justify-between p-4 gap-4" 
      :class="{ 'cursor-pointer hover:bg-gray-50': item.isFolder }" 
      @click="item.isFolder ? toggleExpanded() : null"
    >
      <div class="flex items-center gap-3 flex-1 min-w-0">
        <!-- Folder/File Icon -->
        <svg v-if="item.isFolder" class="w-8 h-8 text-bloodsa-red flex-shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 7v10a2 2 0 002 2h14a2 2 0 002-2V9a2 2 0 00-2-2H5a2 2 0 00-2-2H3z" />
        </svg>
        
        <!-- PDF Icon -->
        <svg v-else-if="getFileIcon(item.name) === 'pdf'" class="w-8 h-8 text-red-600 flex-shrink-0" fill="currentColor" viewBox="0 0 24 24">
          <path d="M14,2H6A2,2 0 0,0 4,4V20A2,2 0 0,0 6,22H18A2,2 0 0,0 20,20V8L14,2M18,20H6V4H13V9H18V20Z" />
          <path d="M8,12H10V14H8V12M12,12H14V14H12V12M16,12H18V14H16V12M8,16H10V18H8V16M12,16H14V18H12V16M16,16H18V18H16V16Z" />
        </svg>
        
        <!-- Word Document Icon -->
        <svg v-else-if="getFileIcon(item.name) === 'word'" class="w-8 h-8 text-blue-600 flex-shrink-0" fill="currentColor" viewBox="0 0 24 24">
          <path d="M14,2H6A2,2 0 0,0 4,4V20A2,2 0 0,0 6,22H18A2,2 0 0,0 20,20V8L14,2M18,20H6V4H13V9H18V20Z" />
          <path d="M8,12H10V14H8V12M12,12H14V14H12V12M16,12H18V14H16V12M8,16H10V18H8V16M12,16H14V18H12V16M16,16H18V18H16V16Z" />
        </svg>
        
        <!-- Excel Icon -->
        <svg v-else-if="getFileIcon(item.name) === 'excel'" class="w-8 h-8 text-green-600 flex-shrink-0" fill="currentColor" viewBox="0 0 24 24">
          <path d="M14,2H6A2,2 0 0,0 4,4V20A2,2 0 0,0 6,22H18A2,2 0 0,0 20,20V8L14,2M18,20H6V4H13V9H18V20Z" />
          <path d="M8,12H10V14H8V12M12,12H14V14H12V12M16,12H18V14H16V12M8,16H10V18H8V16M12,16H14V18H12V16M16,16H18V18H16V16Z" />
        </svg>
        
        <!-- PowerPoint Icon -->
        <svg v-else-if="getFileIcon(item.name) === 'powerpoint'" class="w-8 h-8 text-orange-600 flex-shrink-0" fill="currentColor" viewBox="0 0 24 24">
          <path d="M14,2H6A2,2 0 0,0 4,4V20A2,2 0 0,0 6,22H18A2,2 0 0,0 20,20V8L14,2M18,20H6V4H13V9H18V20Z" />
          <path d="M8,12H10V14H8V12M12,12H14V14H12V12M16,12H18V14H16V12M8,16H10V18H8V16M12,16H14V18H12V16M16,16H18V18H16V16Z" />
        </svg>
        
        <!-- Image Icon -->
        <svg v-else-if="getFileIcon(item.name) === 'image'" class="w-8 h-8 text-purple-600 flex-shrink-0" fill="currentColor" viewBox="0 0 24 24">
          <path d="M8.5,13.5L11,16.5L14.5,12L19,18H5M21,19V5C21,3.89 20.1,3 19,3H5A2,2 0 0,0 3,5V19A2,2 0 0,0 5,21H19A2,2 0 0,0 21,19Z" />
        </svg>
        
        <!-- Text File Icon -->
        <svg v-else-if="getFileIcon(item.name) === 'text'" class="w-8 h-8 text-gray-600 flex-shrink-0" fill="currentColor" viewBox="0 0 24 24">
          <path d="M14,2H6A2,2 0 0,0 4,4V20A2,2 0 0,0 6,22H18A2,2 0 0,0 20,20V8L14,2M18,20H6V4H13V9H18V20Z" />
          <path d="M8,12H10V14H8V12M12,12H14V14H12V12M16,12H18V14H16V12M8,16H10V18H8V16M12,16H14V18H12V16M16,16H18V18H16V16Z" />
        </svg>
        
        <!-- Default File Icon -->
        <svg v-else class="w-8 h-8 text-bloodsa-red flex-shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M7 21h10a2 2 0 002-2V9.414a1 1 0 00-.293-.707l-5.414-5.414A1 1 0 0012.586 3H7a2 2 0 00-2 2v14a2 2 0 002 2z" />
        </svg>
        
        <div class="flex-1 min-w-0">
          <div class="flex items-center gap-2 min-w-0">
            <h3 class="text-xl font-semibold text-gray-700 truncate min-w-0">{{ item.name }}</h3>
            <!-- Expand/Collapse Icon for folders -->
            <svg v-if="item.isFolder" class="w-5 h-5 text-gray-400 transition-transform flex-shrink-0" :class="{ 'rotate-90': isExpanded }" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7" />
            </svg>
          </div>
          <p v-if="!item.isFolder" class="text-sm text-gray-500 mt-1">
            {{ formatFileSize(item.size) }} • {{ formatDate(item.modifiedTime) }}
          </p>
          <p v-else class="text-sm text-gray-500 mt-1">
            {{ item.children?.length || 0 }} items
          </p>
        </div>
      </div>
      
      <!-- Download Button for files -->
      <button 
        v-if="!item.isFolder"
        @click.stop="handleDownload(item)"
        :disabled="downloadingFile === item.path"
        class="px-4 py-2 bg-bloodsa-red text-white rounded-lg hover:bg-red-700 transition-colors disabled:opacity-50 disabled:cursor-wait flex items-center gap-2 flex-shrink-0"
      >
        <svg v-if="downloadingFile !== item.path" class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 10v6m0 0l-3-3m3 3l3-3m2 8H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
        </svg>
        <svg v-else class="animate-spin w-5 h-5" fill="none" viewBox="0 0 24 24">
          <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
          <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
        </svg>
        <span class="hidden sm:inline">{{ downloadingFile === item.path ? 'Downloading...' : 'Download' }}</span>
        <span class="sm:hidden">{{ downloadingFile === item.path ? '...' : '↓' }}</span>
      </button>
    </div>
    
    <!-- Folder Contents (Accordion) -->
    <div v-if="item.isFolder && isExpanded && item.children" class="border-t border-gray-200 bg-gray-50">
      <div class="p-4 space-y-2">
        <DocumentItem 
          v-for="child in item.children" 
          :key="child.path" 
          :item="child"
          :downloading-file="downloadingFile"
          @download="handleDownload"
        />
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'

interface DropboxFile {
  name: string
  path: string
  size: number
  modifiedTime: string
  isFolder: boolean
  children?: DropboxFile[]
}

const props = defineProps<{
  item: DropboxFile
  downloadingFile?: string | null
}>()

const emit = defineEmits<{
  download: [file: DropboxFile]
}>()

const isExpanded = ref(false)

const toggleExpanded = () => {
  isExpanded.value = !isExpanded.value
}

const handleDownload = (file: DropboxFile) => {
  emit('download', file)
}

const formatFileSize = (bytes: number): string => {
  if (bytes === 0) return '0 Bytes'
  const k = 1024
  const sizes = ['Bytes', 'KB', 'MB', 'GB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return Math.round(bytes / Math.pow(k, i) * 100) / 100 + ' ' + sizes[i]
}

const formatDate = (dateString: string): string => {
  const date = new Date(dateString)
  return date.toLocaleDateString('en-US', { year: 'numeric', month: 'short', day: 'numeric' })
}

const getFileIcon = (filename: string): string => {
  const extension = filename.split('.').pop()?.toLowerCase()
  
  // PDF files
  if (extension === 'pdf') return 'pdf'
  
  // Microsoft Word documents
  if (['doc', 'docx'].includes(extension || '')) return 'word'
  
  // Microsoft Excel documents
  if (['xls', 'xlsx', 'csv'].includes(extension || '')) return 'excel'
  
  // Microsoft PowerPoint documents
  if (['ppt', 'pptx'].includes(extension || '')) return 'powerpoint'
  
  // Image files
  if (['jpg', 'jpeg', 'png', 'gif', 'bmp', 'svg', 'webp', 'tiff', 'ico'].includes(extension || '')) return 'image'
  
  // Text files
  if (['txt', 'rtf', 'md', 'log'].includes(extension || '')) return 'text'
  
  // Default
  return 'default'
}
</script>

<style scoped>
.bg-bloodsa-red {
  background-color: #8B0000;
}
.text-bloodsa-red {
  color: #8B0000;
}
</style>
