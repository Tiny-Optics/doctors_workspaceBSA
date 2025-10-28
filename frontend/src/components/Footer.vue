<template>
  <footer class="bg-white border-t border-gray-200">
    <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
      <div class="grid grid-cols-1 md:grid-cols-4 gap-8">
        <!-- BLOODSA Info -->
        <div class="col-span-1 md:col-span-2">
          <div class="flex items-center space-x-3 mb-4">
            <img 
              src="/BLOODSA-Logo.svg" 
              alt="BloodSA Logo" 
              class="w-12 h-12 object-contain"
            />
            <span class="text-lg font-bold text-bloodsa-red">Doctor's Workspace</span>
          </div>
          <p class="text-gray-600 text-sm leading-relaxed">
            BLOODSA Foundation NPC is registered in South Africa with Company Reg. No: 2021/990954/08, Non-Profit Organisation Reg. No: 268-421 NPO  and Public Benefit Organisation Ref. No: 930076265.
          </p>
        </div>

        <!-- Quick Links -->
        <div>
          <h3 class="text-sm font-semibold text-gray-900 uppercase tracking-wider mb-4">Quick Links</h3>
          <ul class="space-y-2">
            <li><a href="https://bloodsa.org.za/about/" target="_blank" class="text-gray-600 hover:text-bloodsa-red text-sm transition-colors duration-200">About BLOODSA</a></li>
            <li><a href="https://bloodsa.org.za/contact/" target="_blank" class="text-gray-600 hover:text-bloodsa-red text-sm transition-colors duration-200">Contact Us</a></li>
            <li><a href="https://bloodsa.org.za/resources/" target="_blank" class="text-gray-600 hover:text-bloodsa-red text-sm transition-colors duration-200">Resources</a></li>
            <li><a href="https://bloodsa.org.za/upcoming-events/" target="_blank" class="text-gray-600 hover:text-bloodsa-red text-sm transition-colors duration-200">Events</a></li>
          </ul>
        </div>

        <!-- Contact -->
        <div>
          <h3 class="text-sm font-semibold text-gray-900 uppercase tracking-wider mb-4">Contact</h3>
          <ul class="space-y-2">
            <li>
              <a 
                href="mailto:admin@bloodsa.org.za" 
                class="text-gray-600 hover:text-bloodsa-red text-sm transition-colors duration-200"
              >
                admin@bloodsa.org.za 
              </a>
            </li>
            <li>
              <a 
                href="tel:+27214043070" 
                class="text-gray-600 hover:text-bloodsa-red text-sm transition-colors duration-200"
              >
              +27 (0) 21 404 3070
              </a>
            </li>
            <li>
              <a 
                href="https://maps.app.goo.gl/xf4hEpThDYei8JNu8" 
                target="_blank" 
                rel="noopener noreferrer"
                class="text-gray-600 hover:text-bloodsa-red text-sm transition-colors duration-200"
              >
              Room E4-58, E5 Haematology Clinic E-floor, New Main Building, Groote Schuur Hospital, Observatory, Cape Town, 7935
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
