<script setup lang="ts">
import Header from './components/Header.vue'
import Footer from './components/Footer.vue'
import ToastNotification from './components/ToastNotification.vue'
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
  </div>
</template>

<style scoped></style>
