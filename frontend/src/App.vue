<template>
  <div>
    <nav class="bg-white shadow-md sticky top-0 z-50 backdrop-blur-sm bg-white/95 animate-slideDown">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div class="flex justify-between items-center h-16 gap-4">
          <div class="flex-shrink-0 animate-fadeInLeft">
            <img 
              src="/Images/grabfast-logo.png" 
              alt="FastGrab" 
              class="h-10 cursor-pointer transition-all duration-300 hover:scale-105 transform"
            />
          </div>
          <div class="flex-1 max-w-md mx-4 animate-fadeInUp">
            <div class="relative group">
              <input
                type="text"
                placeholder="Search..."
                class="w-full px-2 py-2 pl-8 pr-4 text-gray-700 bg-transparent border-0 border-b-2 border-gray-300 focus:outline-none focus:border-green-500 transition-all duration-300 hover:border-gray-400"
              />
              <i class="fas fa-search absolute left-0 top-1/2 transform -translate-y-1/2 text-gray-400 transition-all duration-300 group-hover:text-green-500 group-focus-within:text-green-500"></i>
            </div>
          </div>
          <div class="flex-shrink-0 flex items-center gap-6 animate-fadeInRight">
            <a 
              href="#" 
              class="relative text-gray-700 hover:text-green-600 font-medium transition-all duration-300 group"
            >
              Order
              <span class="absolute bottom-0 left-0 w-0 h-0.5 bg-green-600 transition-all duration-300 group-hover:w-full"></span>
            </a>
            <a 
              href="#" 
              class="relative text-gray-700 hover:text-green-600 font-medium transition-all duration-300 group"
            >
              Partnership
              <span class="absolute bottom-0 left-0 w-0 h-0.5 bg-green-600 transition-all duration-300 group-hover:w-full"></span>
            </a>
            <a 
              href="#" 
              class="relative text-gray-700 hover:text-green-600 font-medium transition-all duration-300 group"
            >
              Sign Up
              <span class="absolute bottom-0 left-0 w-0 h-0.5 bg-green-600 transition-all duration-300 group-hover:w-full"></span>
            </a>
            <a 
              href="#" 
              class="relative text-gray-700 hover:text-green-600 transition-all duration-300 group"
            >
              <i class="fas fa-shopping-cart text-xl transform transition-all duration-300 group-hover:scale-110 group-hover:rotate-[-5deg]"></i>
              <span class="absolute -top-1 -right-1 w-3 h-3 bg-red-500 rounded-full opacity-0 group-hover:opacity-100 transition-opacity duration-300"></span>
            </a>
          </div>
        </div>
      </div>
    </nav>
    <main class="animate-fadeInUp">
    </main>
  </div>
</template>

<script>
export default {
  name: 'App'
}
</script>

<style scoped>
@keyframes slideDown {
  from {
    transform: translateY(-100%);
    opacity: 0;
  }
  to {
    transform: translateY(0);
    opacity: 1;
  }
}

@keyframes fadeInLeft {
  from {
    transform: translateX(-30px);
    opacity: 0;
  }
  to {
    transform: translateX(0);
    opacity: 1;
  }
}

@keyframes fadeInRight {
  from {
    transform: translateX(30px);
    opacity: 0;
  }
  to {
    transform: translateX(0);
    opacity: 1;
  }
}

@keyframes fadeInUp {
  from {
    transform: translateY(20px);
    opacity: 0;
  }
  to {
    transform: translateY(0);
    opacity: 1;
  }
}

.animate-slideDown {
  animation: slideDown 0.6s ease-out;
}

.animate-fadeInLeft {
  animation: fadeInLeft 0.8s ease-out 0.2s both;
}

.animate-fadeInRight {
  animation: fadeInRight 0.8s ease-out 0.4s both;
}

.animate-fadeInUp {
  animation: fadeInUp 0.8s ease-out 0.3s both;
}
</style>
