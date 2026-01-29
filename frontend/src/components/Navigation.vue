<template>
  <div class="relative">
    <!-- Wrapper to maintain height when nav becomes fixed -->
      <div :class="isHomeAtTop ? 'h-24' : 'h-[72px]'">
        <!-- Navigation Header -->
        <nav 
          class="transition-[background-color,box-shadow] duration-300 ease-in-out z-[1000] box-border"
          :class="isHomeAtTop ? 'fixed top-0 left-0 right-0 bg-zinc-900/30' : (homePage === true && isScrolled) ? 'fixed top-0 left-0 right-0 bg-white shadow-md' : (isScrolled ? 'fixed top-0 left-0 right-0 bg-white shadow-md' : 'relative bg-white')"
        >
        <div class="container mx-auto px-4">
          <div class="flex items-stretch justify-between transition-[height] duration-300 ease-in-out" :class="isHomeAtTop ? 'h-24' : 'h-[72px]'">
          <!-- Navigation Links (Left) -->
          <div class="hidden lg:flex items-stretch flex-1">
            <template v-for="link in [
              { href: '/search', text: 'Search', active: isSearchPage },
              { href: '/buy', text: 'Buy', active: isBuyPage },
              { href: '/sell', text: 'Sell', active: isSellPage },
              { href: '/rentals', text: 'Rent', active: isRentalsPage },
              { href: '/about', text: 'About', active: isAboutPage },
              { href: '/blog', text: 'Blog', active: isBlogPage }
            ]" :key="link.href">
              <a 
                v-if="link.active"
                :href="link.href"
                :class="['px-4 h-full flex items-center font-normal no-underline bg-brand-blue text-white', isMounted ? 'transition-colors duration-300 ease-in-out' : '']"
                style="font-family: 'Noto Sans', sans-serif;"
              >{{ link.text }}</a>
              <a 
                v-else-if="isHomeAtTop"
                :href="link.href"
                :class="['px-4 h-full flex items-center font-normal no-underline text-white hover:bg-gray-200/20', isMounted ? 'transition-colors duration-300 ease-in-out' : '']"
                style="font-family: 'Noto Sans', sans-serif;"
              >{{ link.text }}</a>
              <a 
                v-else
                :href="link.href"
                :class="['px-4 h-full flex items-center font-normal no-underline text-[#3A3632] hover:bg-gray-200/50 hover:text-brand-blue', isMounted ? 'transition-colors duration-300 ease-in-out' : '']"
                style="font-family: 'Noto Sans', sans-serif;"
              >{{ link.text }}</a>
            </template>
          </div>
          
          <!-- Logo/Title (Center on desktop, Left on mobile) -->
          <div class="flex items-center lg:justify-center flex-1 lg:flex-1 h-full" :class="isHomeAtTop ? 'py-2' : ''">
            <a href="/" class="flex items-center gap-3 no-underline hover:opacity-80 transition-opacity duration-300">
              <!-- Logo image - single element with conditional classes for smooth transition -->
              <img 
                :src="isHomeAtTop ? '/jp-anchor-black.svg' : '/jp-anchor-blue.svg'"
                alt="Anchor" 
                :class="[
                  'w-auto',
                  isHomeAtTop ? 'brightness-0 invert h-20' : 'h-14',
                  isMounted ? 'transition-[filter,height] duration-300 ease-in-out' : ''
                ]"
              />
              <!-- Title text - single element with conditional classes to prevent re-renders -->
              <div class="flex flex-col leading-tight">
                <span 
                  :class="[
                    'font-semibold tracking-tight border-b-2 pb-1 -mt-[5px]',
                    isHomeAtTop ? 'text-white border-white text-2xl' : 'text-[#3A3632] border-brand-blue text-xl',
                    isMounted ? 'transition-[color,border-color,font-size] duration-300 ease-in-out' : ''
                  ]" 
                  style="font-family: 'Noto Serif', serif;"
                >
                  JULIE POGUE
                </span>
                <span 
                  :class="[
                    'font-light tracking-[0.55em]',
                    isHomeAtTop ? 'text-white text-sm' : 'text-[#3A3632] text-xs',
                    isMounted ? 'transition-[color,font-size] duration-300 ease-in-out' : ''
                  ]" 
                  style="font-family: 'Noto Sans', sans-serif;"
                >
                  PROPERTIES
                </span>
              </div>
            </a>
          </div>
          
          <!-- Contact Info (Right) -->
          <div class="hidden lg:flex items-stretch flex-1 justify-end">
            <template v-if="isHomeAtTop">
              <a href="tel:+1-502-238-7400" :class="['px-4 h-full flex items-center font-normal no-underline text-white hover:opacity-70', isMounted ? 'transition-[color,opacity] duration-300 ease-in-out' : '']" style="font-family: 'Noto Sans', sans-serif;">(502) 238-7400</a>
              <a 
                v-if="isContactActive"
                href="/contact"
                :class="['px-4 h-full flex items-center font-normal no-underline bg-brand-blue text-white', isMounted ? 'transition-colors duration-300 ease-in-out' : '']"
                style="font-family: 'Noto Sans', sans-serif;"
              >Contact</a>
              <a 
                v-else
                href="/contact"
                :class="['px-4 h-full flex items-center font-normal no-underline text-white hover:bg-gray-200/20', isMounted ? 'transition-colors duration-300 ease-in-out' : '']"
                style="font-family: 'Noto Sans', sans-serif;"
              >Contact</a>
              <div class="flex items-stretch">
                <a href="https://www.facebook.com/website-dummy-properties" target="_blank" rel="noopener noreferrer" :class="['px-4 h-full flex items-center font-normal no-underline text-white hover:opacity-70', isMounted ? 'transition-[color,opacity] duration-300 ease-in-out' : '']" aria-label="Facebook">
                  <font-awesome-icon :icon="['fab', 'facebook-f']" />
                </a>
                <a href="https://www.instagram.com/website-dummy-properties/" target="_blank" rel="noopener noreferrer" :class="['px-4 h-full flex items-center font-normal no-underline text-white hover:opacity-70', isMounted ? 'transition-[color,opacity] duration-300 ease-in-out' : '']" aria-label="Instagram">
                  <font-awesome-icon :icon="['fab', 'instagram']" />
                </a>
                <a href="https://www.linkedin.com/company/julie-pogue-properties/" target="_blank" rel="noopener noreferrer" :class="['px-4 h-full flex items-center font-normal no-underline text-white hover:opacity-70', isMounted ? 'transition-[color,opacity] duration-300 ease-in-out' : '']" aria-label="LinkedIn">
                  <font-awesome-icon :icon="['fab', 'linkedin']" />
                </a>
                <a href="https://share.google/XaJQzbw3k0q99A90Y" target="_blank" rel="noopener noreferrer" :class="['px-4 h-full flex items-center font-normal no-underline text-white hover:opacity-70', isMounted ? 'transition-[color,opacity] duration-300 ease-in-out' : '']" aria-label="Google Places">
                  <font-awesome-icon :icon="['fab', 'google']" />
                </a>
              </div>
            </template>
            <template v-else>
              <a href="tel:+1-502-238-7400" :class="['px-4 h-full flex items-center font-normal no-underline text-[#3A3632] hover:opacity-60', isMounted ? 'transition-[color,opacity] duration-300 ease-in-out' : '']" style="font-family: 'Noto Sans', sans-serif;">(502) 238-7400</a>
              <a 
                v-if="isContactActive"
                href="/contact"
                :class="['px-4 h-full flex items-center font-normal no-underline bg-brand-blue text-white', isMounted ? 'transition-colors duration-300 ease-in-out' : '']"
                style="font-family: 'Noto Sans', sans-serif;"
              >Contact</a>
              <a 
                v-else
                href="/contact"
                :class="['px-4 h-full flex items-center font-normal no-underline text-[#3A3632] hover:bg-gray-200/50 hover:text-brand-blue', isMounted ? 'transition-colors duration-300 ease-in-out' : '']"
                style="font-family: 'Noto Sans', sans-serif;"
              >Contact</a>
              <div class="flex items-stretch">
                <a href="https://www.facebook.com/website-dummy-properties" target="_blank" rel="noopener noreferrer" :class="['px-4 h-full flex items-center font-normal no-underline text-[#3A3632] hover:opacity-60', isMounted ? 'transition-[color,opacity] duration-300 ease-in-out' : '']" aria-label="Facebook">
                  <font-awesome-icon :icon="['fab', 'facebook-f']" />
                </a>
                <a href="https://www.instagram.com/website-dummy-properties/" target="_blank" rel="noopener noreferrer" :class="['px-4 h-full flex items-center font-normal no-underline text-[#3A3632] hover:opacity-60', isMounted ? 'transition-[color,opacity] duration-300 ease-in-out' : '']" aria-label="Instagram">
                  <font-awesome-icon :icon="['fab', 'instagram']" />
                </a>
                <a href="https://www.linkedin.com/company/julie-pogue-properties/" target="_blank" rel="noopener noreferrer" :class="['px-4 h-full flex items-center font-normal no-underline text-[#3A3632] hover:opacity-60', isMounted ? 'transition-[color,opacity] duration-300 ease-in-out' : '']" aria-label="LinkedIn">
                  <font-awesome-icon :icon="['fab', 'linkedin']" />
                </a>
                <a href="https://share.google/XaJQzbw3k0q99A90Y" target="_blank" rel="noopener noreferrer" :class="['px-4 h-full flex items-center font-normal no-underline text-[#3A3632] hover:opacity-60', isMounted ? 'transition-[color,opacity] duration-300 ease-in-out' : '']" aria-label="Google Places">
                  <font-awesome-icon :icon="['fab', 'google']" />
                </a>
              </div>
            </template>
          </div>

          <!-- Mobile Menu Button -->
          <button
            v-if="isHomeAtTop"
            @click="mobileMenuOpen = !mobileMenuOpen"
            class="lg:hidden transition-colors duration-300 p-2 text-white hover:text-brand-dark-gold"
            aria-label="Toggle menu"
          >
            <font-awesome-icon :icon="mobileMenuOpen ? 'times' : 'bars'" class="text-2xl" />
          </button>
          <button
            v-else
            @click="mobileMenuOpen = !mobileMenuOpen"
            class="lg:hidden transition-colors duration-300 p-2 text-[#3A3632] hover:text-brand-blue"
            aria-label="Toggle menu"
          >
            <font-awesome-icon :icon="mobileMenuOpen ? 'times' : 'bars'" class="text-2xl" />
          </button>
        </div>
      </div>
      </nav>
    </div>

    <!-- Mobile Menu -->
    <transition name="fade">
      <div 
        v-if="mobileMenuOpen"
        :class="isHomeAtTop ? 'lg:hidden border-b fixed left-0 right-0 shadow-lg bg-zinc-900 border-brand-dark-gold' : 'lg:hidden border-b fixed left-0 right-0 shadow-lg bg-[#FAF7F2] border-[#E8E3D9]'"
        :style="isHomeAtTop ? { top: '96px', zIndex: 9998, maxHeight: 'calc(100vh - 96px)', overflowY: 'auto', position: 'fixed' } : (isScrolled ? { top: '64px', zIndex: 9998, maxHeight: 'calc(100vh - 64px)', overflowY: 'auto', position: 'fixed' } : { top: '72px', zIndex: 9998, maxHeight: 'calc(100vh - 72px)', overflowY: 'auto', position: 'fixed' })"
      >
        <div class="px-4 py-4 space-y-3">
          <template v-for="link in [
            { href: '/search', text: 'Search', active: isSearchPage },
            { href: '/buy', text: 'Buy', active: isBuyPage },
            { href: '/sell', text: 'Sell', active: isSellPage },
            { href: '/rentals', text: 'Rent', active: isRentalsPage },
            { href: '/about', text: 'About', active: isAboutPage },
            { href: '/blog', text: 'Blog', active: isBlogPage }
          ]" :key="link.href">
            <a 
              v-if="link.active"
              :href="link.href"
              class="block transition-colors duration-200 py-2 px-3 rounded bg-brand-blue text-white"
              style="font-family: 'Noto Sans', sans-serif;"
              @click="mobileMenuOpen = false"
            >{{ link.text }}</a>
            <a 
              v-else-if="isHomeAtTop"
              :href="link.href"
              class="block transition-colors duration-200 py-2 px-3 rounded text-white hover:text-brand-dark-gold"
              style="font-family: 'Noto Sans', sans-serif;"
              @click="mobileMenuOpen = false"
            >{{ link.text }}</a>
            <a 
              v-else
              :href="link.href"
              class="block transition-colors duration-200 py-2 px-3 rounded text-[#3A3632] hover:text-brand-blue"
              style="font-family: 'Noto Sans', sans-serif;"
              @click="mobileMenuOpen = false"
            >{{ link.text }}</a>
          </template>
          <div :class="isHomeAtTop ? 'pt-2 mt-2 border-t border-brand-dark-gold' : 'pt-2 mt-2 border-t border-[#E8E3D9]'">
            <template v-if="isHomeAtTop">
              <a href="tel:+1-502-238-7400" class="block transition-colors duration-200 py-2 px-3 rounded text-white hover:text-brand-dark-gold" style="font-family: 'Noto Sans', sans-serif;" @click="mobileMenuOpen = false">(502) 238-7400</a>
              <a href="/contact" :class="isContactActive ? 'block transition-colors duration-200 py-2 px-3 rounded bg-brand-blue text-white' : 'block transition-colors duration-200 py-2 px-3 rounded text-white hover:text-brand-dark-gold'" style="font-family: 'Noto Sans', sans-serif;" @click="mobileMenuOpen = false">Contact</a>
            </template>
            <template v-else>
              <a href="tel:+1-502-238-7400" class="block transition-colors duration-200 py-2 px-3 rounded text-[#3A3632] hover:text-brand-blue" style="font-family: 'Noto Sans', sans-serif;" @click="mobileMenuOpen = false">(502) 238-7400</a>
              <a href="/contact" :class="isContactActive ? 'block transition-colors duration-200 py-2 px-3 rounded bg-brand-blue text-white' : 'block transition-colors duration-200 py-2 px-3 rounded text-[#3A3632] hover:text-brand-blue'" style="font-family: 'Noto Sans', sans-serif;" @click="mobileMenuOpen = false">Contact</a>
            </template>
            <div class="flex items-center gap-3 pt-2">
              <template v-if="isHomeAtTop">
                <a href="https://www.facebook.com/website-dummy-properties" target="_blank" rel="noopener noreferrer" class="block transition-colors duration-200 py-2 px-3 rounded text-white hover:text-brand-dark-gold" aria-label="Facebook" @click="mobileMenuOpen = false">
                  <font-awesome-icon :icon="['fab', 'facebook-f']" />
                </a>
                <a href="https://www.instagram.com/website-dummy-properties/" target="_blank" rel="noopener noreferrer" class="block transition-colors duration-200 py-2 px-3 rounded text-white hover:text-brand-dark-gold" aria-label="Instagram" @click="mobileMenuOpen = false">
                  <font-awesome-icon :icon="['fab', 'instagram']" />
                </a>
                <a href="https://www.linkedin.com/company/julie-pogue-properties/" target="_blank" rel="noopener noreferrer" class="block transition-colors duration-200 py-2 px-3 rounded text-white hover:text-brand-dark-gold" aria-label="LinkedIn" @click="mobileMenuOpen = false">
                  <font-awesome-icon :icon="['fab', 'linkedin']" />
                </a>
                <a href="https://share.google/XaJQzbw3k0q99A90Y" target="_blank" rel="noopener noreferrer" class="block transition-colors duration-200 py-2 px-3 rounded text-white hover:text-brand-dark-gold" aria-label="Google Places" @click="mobileMenuOpen = false">
                  <font-awesome-icon :icon="['fab', 'google']" />
                </a>
              </template>
              <template v-else>
                <a href="https://www.facebook.com/website-dummy-properties" target="_blank" rel="noopener noreferrer" class="block transition-colors duration-200 py-2 px-3 rounded text-[#3A3632] hover:text-brand-blue" aria-label="Facebook" @click="mobileMenuOpen = false">
                  <font-awesome-icon :icon="['fab', 'facebook-f']" />
                </a>
                <a href="https://www.instagram.com/website-dummy-properties/" target="_blank" rel="noopener noreferrer" class="block transition-colors duration-200 py-2 px-3 rounded text-[#3A3632] hover:text-brand-blue" aria-label="Instagram" @click="mobileMenuOpen = false">
                  <font-awesome-icon :icon="['fab', 'instagram']" />
                </a>
                <a href="https://www.linkedin.com/company/julie-pogue-properties/" target="_blank" rel="noopener noreferrer" class="block transition-colors duration-200 py-2 px-3 rounded text-[#3A3632] hover:text-brand-blue" aria-label="LinkedIn" @click="mobileMenuOpen = false">
                  <font-awesome-icon :icon="['fab', 'linkedin']" />
                </a>
                <a href="https://share.google/XaJQzbw3k0q99A90Y" target="_blank" rel="noopener noreferrer" class="block transition-colors duration-200 py-2 px-3 rounded text-[#3A3632] hover:text-brand-blue" aria-label="Google Places" @click="mobileMenuOpen = false">
                  <font-awesome-icon :icon="['fab', 'google']" />
                </a>
              </template>
            </div>
          </div>
        </div>
      </div>
    </transition>
  </div>
</template>

<script>
export default {
  name: 'Navigation',
  props: {
    homePage: {
      type: Boolean,
      default: false,
      required: false
    }
  },
  data() {
    // Initialize state immediately to prevent flicker
    const initialPath = typeof window !== 'undefined' ? window.location.pathname : '/'
    const initialScroll = typeof window !== 'undefined' ? window.scrollY > 0 : false
    
    return {
      isScrolled: initialScroll,
      mobileMenuOpen: false,
      currentPath: initialPath,
      currentHash: typeof window !== 'undefined' ? window.location.hash : '',
      // Track if component has mounted to disable transitions on initial render
      isMounted: false
    }
  },
  watch: {
    currentPath() {
      // Close mobile menu when route changes
      this.mobileMenuOpen = false
    }
  },
  computed: {
    isHomeAtTop() {
      // Simple: home page prop is true AND not scrolled
      return this.homePage === true && !this.isScrolled
    },
    isSearchPage() {
      return this.currentPath === '/search' || this.currentPath.startsWith('/property/')
    },
    isBuyPage() {
      return this.currentPath === '/buy'
    },
    isSellPage() {
      return this.currentPath === '/sell'
    },
    isAboutPage() {
      return this.currentPath === '/about'
    },
    isRentalsPage() {
      return this.currentPath === '/rentals' || this.currentPath.startsWith('/rentals')
    },
    isBlogPage() {
      return this.currentPath === '/blog' || this.currentPath.startsWith('/blog/')
    },
    isContactActive() {
      return this.currentPath === '/contact'
    }
  },
  mounted() {
    // Update path immediately on mount
    this.updatePath()
    
    // Intercept link clicks to update path immediately
    this.interceptLinkClicks()
    
    window.addEventListener('scroll', this.handleScroll)
    window.addEventListener('popstate', this.updatePath)
    window.addEventListener('hashchange', this.handleHashChange)
    this.handleScroll() // Check initial scroll position
    
    // Mark as mounted after next tick to enable transitions for future state changes
    this.$nextTick(() => {
      this.isMounted = true
    })
    
    // Fallback: check pathname periodically in case navigation happens outside our intercept
    this.pathCheckInterval = setInterval(() => {
      const newPath = window.location.pathname
      if (newPath !== this.currentPath) {
        this.currentPath = newPath
      }
    }, 100)
  },
  beforeUnmount() {
    window.removeEventListener('scroll', this.handleScroll)
    window.removeEventListener('popstate', this.updatePath)
    window.removeEventListener('hashchange', this.handleHashChange)
    if (this.pathCheckInterval) {
      clearInterval(this.pathCheckInterval)
    }
    // Link clicks are handled by App.vue, not Navigation component
  },
  methods: {
    updatePath() {
      this.currentPath = window.location.pathname
      this.currentHash = window.location.hash
    },
    interceptLinkClicks() {
      // Navigation component doesn't handle link clicks - App.vue handles all routing
      // This method is kept for compatibility but does nothing
    },
    handleScroll() {
      this.isScrolled = window.scrollY > 0
    },
    handleHashChange() {
      this.updatePath()
    },
  }
}
</script>

<style scoped>
/* Mobile Menu Animation - Simple Fade */
.fade-enter-active {
  transition: opacity 0.3s ease;
}

.fade-leave-active {
  transition: opacity 0.3s ease;
}

.fade-enter-from {
  opacity: 0;
}

.fade-enter-to {
  opacity: 1;
}

.fade-leave-from {
  opacity: 1;
}

.fade-leave-to {
  opacity: 0;
}

/* Prevent visited link styling - always use our custom colors */
nav a:visited {
  color: inherit;
}

/* Ensure active state takes precedence */
nav a.bg-brand-blue {
  color: white !important;
}

/* Force white text on home page at top - override any other styles */
/* Target nav with dark background (home page at top) */
nav[class*="bg-zinc-900"] a {
  color: white !important;
}

/* Prevent image flicker on home page */
nav img {
  will-change: filter;
  backface-visibility: hidden;
  transform: translateZ(0);
}
</style>

