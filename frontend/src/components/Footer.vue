<template>
  <footer class="bg-white border-t border-gray-200">
    <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
      <div class="grid grid-cols-1 md:grid-cols-4 gap-8">
        <!-- BLOODSA Info -->
        <div class="col-span-1 md:col-span-2">
          <div class="flex items-center space-x-3 mb-4">
            <div class="w-8 h-8 bg-bloodsa-red rounded-full flex items-center justify-center">
              <svg class="w-5 h-5 text-white" fill="currentColor" viewBox="0 0 24 24">
                <path d="M12 2C6.48 2 2 6.48 2 12s4.48 10 10 10 10-4.48 10-10S17.52 2 12 2zm-2 15l-5-5 1.41-1.41L10 14.17l7.59-7.59L19 8l-9 9z"/>
              </svg>
            </div>
            <span class="text-lg font-bold text-bloodsa-red">BLOODSA</span>
          </div>
          <p class="text-gray-600 text-sm leading-relaxed">
            The Blood and Marrow Transplant Society of South Africa (BLOODSA) is dedicated to advancing 
            hematopoietic stem cell transplantation and cellular therapy in South Africa.
          </p>
        </div>

        <!-- Quick Links -->
        <div>
          <h3 class="text-sm font-semibold text-gray-900 uppercase tracking-wider mb-4">Quick Links</h3>
          <ul class="space-y-2">
            <li><a href="#" class="text-gray-600 hover:text-bloodsa-red text-sm transition-colors duration-200">About BLOODSA</a></li>
            <li><a href="#" class="text-gray-600 hover:text-bloodsa-red text-sm transition-colors duration-200">Membership</a></li>
            <li><a href="#" class="text-gray-600 hover:text-bloodsa-red text-sm transition-colors duration-200">Resources</a></li>
            <li><a href="#" class="text-gray-600 hover:text-bloodsa-red text-sm transition-colors duration-200">Events</a></li>
          </ul>
        </div>

        <!-- Contact -->
        <div>
          <h3 class="text-sm font-semibold text-gray-900 uppercase tracking-wider mb-4">Contact</h3>
          <ul class="space-y-2">
            <li class="text-gray-600 text-sm">info@bloodsa.org.za</li>
            <li class="text-gray-600 text-sm">+27 (0) 11 123 4567</li>
            <li class="text-gray-600 text-sm">Johannesburg, South Africa</li>
          </ul>
        </div>
      </div>

      <!-- Bottom Bar -->
      <div class="mt-8 pt-8 border-t border-gray-200">
        <div class="flex flex-col md:flex-row justify-between items-center">
          <p class="text-gray-500 text-sm">
            © {{ currentYear }} BLOODSA. All rights reserved. Website developed by <a href="https://tinyoptics.co.za/" target="_blank" class="text-bloodsa-red hover:text-bloodsa-light-red transition-colors duration-200">Tiny Optics</a>.
          </p>
          <div class="flex space-x-6 mt-4 md:mt-0">
            <a href="#" class="text-gray-400 hover:text-bloodsa-red transition-colors duration-200">
              <span class="sr-only">Privacy Policy</span>
              Privacy Policy
            </a>
            <a href="#" class="text-gray-400 hover:text-bloodsa-red transition-colors duration-200">
              <span class="sr-only">Terms of Service</span>
              Terms of Service
            </a>
          </div>
        </div>
      </div>
    </div>
  </footer>
</template>

<script setup lang="ts">
import { computed } from 'vue'

const currentYear = computed(() => new Date().getFullYear())
</script>
