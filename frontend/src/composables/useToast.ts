import { ref } from 'vue'

export interface Toast {
  id: number
  message: string
  type: 'success' | 'error' | 'warning' | 'info'
  duration: number
}

const toasts = ref<Toast[]>([])
let nextId = 0

export function useToast() {
  const addToast = (message: string, type: Toast['type'] = 'info', duration = 5000) => {
    const id = nextId++
    const toast: Toast = { id, message, type, duration }
    
    toasts.value.push(toast)
    
    // Auto-remove after duration
    if (duration > 0) {
      setTimeout(() => {
        removeToast(id)
      }, duration)
    }
    
    return id
  }
  
  const removeToast = (id: number) => {
    const index = toasts.value.findIndex(t => t.id === id)
    if (index > -1) {
      toasts.value.splice(index, 1)
    }
  }
  
  const success = (message: string, duration = 5000) => {
    return addToast(message, 'success', duration)
  }
  
  const error = (message: string, duration = 5000) => {
    return addToast(message, 'error', duration)
  }
  
  const warning = (message: string, duration = 5000) => {
    return addToast(message, 'warning', duration)
  }
  
  const info = (message: string, duration = 5000) => {
    return addToast(message, 'info', duration)
  }
  
  return {
    toasts,
    addToast,
    removeToast,
    success,
    error,
    warning,
    info
  }
}

