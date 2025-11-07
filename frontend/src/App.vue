<script setup lang="ts">
import Header from './components/Header.vue'
import Footer from './components/Footer.vue'
import ToastNotification from './components/ToastNotification.vue'
import CookieConsent from './components/CookieConsent.vue'
</script>

<template>
  <div id="app" class="min-h-screen flex flex-col">
    <Header />
    <main class="flex-1">
      <router-view />
    </main>
    <Footer />
    
    <!-- Toast Notifications (Global) -->
    <ToastNotification />
    
    <!-- Cookie Consent Banner -->
    <CookieConsent />
  </div>
</template>

<style scoped></style>
