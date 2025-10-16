<template>
  <div class="min-h-screen bg-gradient-to-br from-gray-50 to-gray-100">
    <!-- Hero Section -->
    <section class="relative overflow-hidden">
      <!-- Background Pattern -->
      <div class="absolute inset-0 opacity-5">
        <div class="absolute inset-0" style="background-image: repeating-linear-gradient(45deg, #8B0000 0, #8B0000 1px, transparent 0, transparent 50%); background-size: 10px 10px;"></div>
      </div>

      <div class="relative max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-20 md:py-32">
        <div class="text-center">
          <!-- Main Heading -->
          <h1 class="text-4xl md:text-6xl font-bold text-gray-900 mb-6">
            Welcome to the
            <span class="text-bloodsa-red block mt-2">Doctor's Workspace</span>
          </h1>
          
          <!-- Subheading -->
          <p class="text-xl md:text-2xl text-gray-600 max-w-3xl mx-auto mb-12 leading-relaxed">
            A comprehensive platform for haematologists, physicians, and data capturers to access clinical resources, manage workflows, and collaborate effectively.
          </p>

          <!-- CTA Button -->
          <div class="flex flex-col sm:flex-row gap-4 justify-center items-center">
            <button class="btn-primary text-lg px-8 py-4 shadow-lg hover:shadow-xl transform hover:-translate-y-0.5 transition-all duration-200">
              Login to Access Workspace
            </button>
            <a href="#features" class="text-bloodsa-red font-medium hover:text-bloodsa-light-red transition-colors duration-200">
              Learn More →
            </a>
          </div>
        </div>
      </div>
    </section>

    <!-- Features Section -->
    <section id="features" class="py-20 bg-white">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div class="text-center mb-16">
          <h2 class="text-3xl md:text-4xl font-bold text-gray-900 mb-4">
            What's Inside the Workspace
          </h2>
          <p class="text-lg text-gray-600 max-w-2xl mx-auto">
            Login to access a comprehensive suite of tools designed for healthcare professionals
          </p>
        </div>

        <!-- Feature Cards -->
        <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-8">
          <!-- Feature 1: SOPs -->
          <div class="group">
            <div class="bg-gradient-to-br from-gray-50 to-white p-8 rounded-xl border-2 border-gray-100 hover:border-bloodsa-red transition-all duration-300 hover:shadow-lg">
              <div class="w-12 h-12 bg-bloodsa-red rounded-lg flex items-center justify-center mb-4 group-hover:scale-110 transition-transform duration-300">
                <svg class="w-6 h-6 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
                </svg>
              </div>
              <h3 class="text-xl font-semibold text-gray-900 mb-2">Standard Operating Procedures</h3>
              <p class="text-gray-600">
                Access comprehensive SOPs for Anemia, Lymphoma, Myeloma, and General Business procedures.
              </p>
            </div>
          </div>

          <!-- Feature 2: Referrals -->
          <div class="group">
            <div class="bg-gradient-to-br from-gray-50 to-white p-8 rounded-xl border-2 border-gray-100 hover:border-bloodsa-red transition-all duration-300 hover:shadow-lg">
              <div class="w-12 h-12 bg-bloodsa-red rounded-lg flex items-center justify-center mb-4 group-hover:scale-110 transition-transform duration-300">
                <svg class="w-6 h-6 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 7h12m0 0l-4-4m4 4l-4 4m0 6H4m0 0l4 4m-4-4l4-4" />
                </svg>
              </div>
              <h3 class="text-xl font-semibold text-gray-900 mb-2">Transplant Referrals</h3>
              <p class="text-gray-600">
                Streamlined integration with REDCap for efficient transplant referral submissions and tracking.
              </p>
            </div>
          </div>

          <!-- Feature 3: Registry -->
          <div class="group">
            <div class="bg-gradient-to-br from-gray-50 to-white p-8 rounded-xl border-2 border-gray-100 hover:border-bloodsa-red transition-all duration-300 hover:shadow-lg">
              <div class="w-12 h-12 bg-bloodsa-red rounded-lg flex items-center justify-center mb-4 group-hover:scale-110 transition-transform duration-300">
                <svg class="w-6 h-6 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 11H5m14 0a2 2 0 012 2v6a2 2 0 01-2 2H5a2 2 0 01-2-2v-6a2 2 0 012-2m14 0V9a2 2 0 00-2-2M5 11V9a2 2 0 012-2m0 0V5a2 2 0 012-2h6a2 2 0 012 2v2M7 7h10" />
                </svg>
              </div>
              <h3 class="text-xl font-semibold text-gray-900 mb-2">African HOPeR Registry</h3>
              <p class="text-gray-600">
                Manage documentation, training materials, and ethics approvals for registry participation.
              </p>
            </div>
          </div>

          <!-- Feature 4: Document Management -->
          <div class="group">
            <div class="bg-gradient-to-br from-gray-50 to-white p-8 rounded-xl border-2 border-gray-100 hover:border-bloodsa-red transition-all duration-300 hover:shadow-lg">
              <div class="w-12 h-12 bg-bloodsa-red rounded-lg flex items-center justify-center mb-4 group-hover:scale-110 transition-transform duration-300">
                <svg class="w-6 h-6 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 7v10a2 2 0 002 2h14a2 2 0 002-2V9a2 2 0 00-2-2h-6l-2-2H5a2 2 0 00-2 2z" />
                </svg>
              </div>
              <h3 class="text-xl font-semibold text-gray-900 mb-2">Secure Document Storage</h3>
              <p class="text-gray-600">
                Upload and access de-identified documents with enterprise-grade security and backup.
              </p>
            </div>
          </div>

          <!-- Feature 5: Collaboration -->
          <div class="group">
            <div class="bg-gradient-to-br from-gray-50 to-white p-8 rounded-xl border-2 border-gray-100 hover:border-bloodsa-red transition-all duration-300 hover:shadow-lg">
              <div class="w-12 h-12 bg-bloodsa-red rounded-lg flex items-center justify-center mb-4 group-hover:scale-110 transition-transform duration-300">
                <svg class="w-6 h-6 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 20h5v-2a3 3 0 00-5.356-1.857M17 20H7m10 0v-2c0-.656-.126-1.283-.356-1.857M7 20H2v-2a3 3 0 015.356-1.857M7 20v-2c0-.656.126-1.283.356-1.857m0 0a5.002 5.002 0 019.288 0M15 7a3 3 0 11-6 0 3 3 0 016 0zm6 3a2 2 0 11-4 0 2 2 0 014 0zM7 10a2 2 0 11-4 0 2 2 0 014 0z" />
                </svg>
              </div>
              <h3 class="text-xl font-semibold text-gray-900 mb-2">Team Collaboration</h3>
              <p class="text-gray-600">
                Connect with haematologists and physicians across South Africa for better patient care.
              </p>
            </div>
          </div>

          <!-- Feature 6: Training -->
          <div class="group">
            <div class="bg-gradient-to-br from-gray-50 to-white p-8 rounded-xl border-2 border-gray-100 hover:border-bloodsa-red transition-all duration-300 hover:shadow-lg">
              <div class="w-12 h-12 bg-bloodsa-red rounded-lg flex items-center justify-center mb-4 group-hover:scale-110 transition-transform duration-300">
                <svg class="w-6 h-6 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 6.253v13m0-13C10.832 5.477 9.246 5 7.5 5S4.168 5.477 3 6.253v13C4.168 18.477 5.754 18 7.5 18s3.332.477 4.5 1.253m0-13C13.168 5.477 14.754 5 16.5 5c1.747 0 3.332.477 4.5 1.253v13C19.832 18.477 18.247 18 16.5 18c-1.746 0-3.332.477-4.5 1.253" />
                </svg>
              </div>
              <h3 class="text-xl font-semibold text-gray-900 mb-2">Training Resources</h3>
              <p class="text-gray-600">
                Access video tutorials, guidelines, and example documents to enhance your practice.
              </p>
            </div>
          </div>
        </div>
      </div>
    </section>

    <!-- CTA Section -->
    <section class="bg-bloodsa-red py-16">
      <div class="max-w-4xl mx-auto px-4 sm:px-6 lg:px-8 text-center">
        <h2 class="text-3xl md:text-4xl font-bold text-white mb-6">
          Ready to Get Started?
        </h2>
        <p class="text-xl text-white opacity-90 mb-8">
          Login to access your personalized workspace and start managing your clinical workflows more efficiently.
        </p>
        <button class="bg-white text-bloodsa-red px-8 py-4 rounded-md font-semibold text-lg hover:bg-gray-100 transition-colors duration-200 shadow-lg">
          Access Workspace Now
        </button>
      </div>
    </section>
  </div>
</template>

<script setup lang="ts">
// Home page component
</script>
