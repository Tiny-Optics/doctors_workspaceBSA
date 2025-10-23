<template>
  <footer class="bg-white border-t border-gray-200">
    <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
      <div class="grid grid-cols-1 md:grid-cols-4 gap-8">
        <!-- BLOODSA Info -->
        <div class="col-span-1 md:col-span-2">
          <div class="flex items-center space-x-3 mb-4">
            <img 
              src="/BLOODSA-SVG-Logo.svg" 
              alt="BloodSA Logo" 
              class="w-32 h-32 object-contain"
            />
            <span class="text-lg font-bold text-bloodsa-red">Doctor's Workspace</span>
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
            <li>
              <a 
                href="mailto:info@bloodsa.org.za" 
                class="text-gray-600 hover:text-bloodsa-red text-sm transition-colors duration-200"
              >
                info@bloodsa.org.za
              </a>
            </li>
            <li>
              <a 
                href="tel:+27111234567" 
                class="text-gray-600 hover:text-bloodsa-red text-sm transition-colors duration-200"
              >
                +27 (0) 11 123 4567
              </a>
            </li>
            <li>
              <a 
                href="https://maps.app.goo.gl/xf4hEpThDYei8JNu8" 
                target="_blank" 
                rel="noopener noreferrer"
                class="text-gray-600 hover:text-bloodsa-red text-sm transition-colors duration-200"
              >
                Johannesburg, South Africa
              </a>
            </li>
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
