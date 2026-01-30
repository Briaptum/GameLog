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
          <div class="flex-1 max-w-md mx-4 animate-fadeInUp search-bar-container">
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
            <button 
              @click="toggleMobileSearch"
              class="relative text-gray-700 hover:text-green-600 transition-all duration-300 group search-icon-mobile focus:outline-none"
            >
              <i class="fas fa-search text-xl transform transition-all duration-300 group-hover:scale-110"></i>
            </button>
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
    <!-- Mobile Search Bar -->
    <div 
      v-if="showMobileSearch"
      class="bg-white shadow-lg border-t border-gray-200 mobile-search-container"
    >
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-4">
        <div class="relative group">
          <input
            type="text"
            placeholder="Search..."
            v-model="searchQuery"
            @keyup.esc="closeMobileSearch"
            class="w-full px-2 py-3 pl-8 pr-10 text-gray-700 bg-transparent border-0 border-b-2 border-gray-300 focus:outline-none focus:border-gray-400 transition-all duration-300 hover:border-gray-400"
            ref="mobileSearchInput"
          />
          <i class="fas fa-search absolute left-0 top-1/2 transform -translate-y-1/2 text-gray-400 transition-all duration-300 group-focus-within:text-green-500"></i>
          <button
            @click="closeMobileSearch"
            class="absolute right-0 top-1/2 transform -translate-y-1/2 text-gray-400 hover:text-gray-600 transition-colors duration-300 focus:outline-none"
          >
            <i class="fas fa-times text-lg"></i>
          </button>
        </div>
      </div>
    </div>
    <main class="animate-fadeInUp">
      <div class="hero-section flex items-center justify-center min-h-[80vh] px-4">
        <div class="text-center max-w-6xl mx-auto">
          <h1 class="hero-title">
            <span class="hero-word">GRAB</span>
            <span class="hero-word">YOUR</span>
            <span class="hero-word">ORDER</span>
            <span class="hero-word-fast">FAST</span>
          </h1>
          
          <!-- Hero Buttons -->
          <div class="hero-buttons mt-12">
            <button class="btn-hero btn-orders">
              Orders
            </button>
            <button class="btn-hero btn-download">
              Download App
            </button>
          </div>
        </div>
      </div>
      
      <!-- Wave Divider -->
      <div class="wave-divider">
        <svg viewBox="0 0 1200 120" preserveAspectRatio="none" xmlns="http://www.w3.org/2000/svg">
          <defs>
            <!-- Gradient for wave -->
            <linearGradient id="waveGradient" x1="0%" y1="0%" x2="0%" y2="100%">
              <stop offset="0%" style="stop-color:#16a34a;stop-opacity:1" />
              <stop offset="50%" style="stop-color:#15803d;stop-opacity:0.9" />
              <stop offset="100%" style="stop-color:#16a34a;stop-opacity:0.8" />
            </linearGradient>
            <!-- Second gradient layer -->
            <linearGradient id="waveGradient2" x1="0%" y1="0%" x2="0%" y2="100%">
              <stop offset="0%" style="stop-color:#15803d;stop-opacity:0.7" />
              <stop offset="100%" style="stop-color:#16a34a;stop-opacity:0.6" />
            </linearGradient>
          </defs>
          <!-- Main wave with gradient -->
          <path d="M0,60 C300,100 400,20 600,60 C800,100 900,20 1200,60 L1200,120 L0,120 Z" fill="url(#waveGradient)"></path>
          <!-- Additional wave layer -->
          <path d="M0,80 C250,110 350,30 550,70 C750,110 850,30 1050,70 C1100,85 1150,45 1200,60 L1200,120 L0,120 Z" fill="url(#waveGradient2)"></path>
        </svg>
      </div>

      <!-- Hot Orders Section -->
      <section class="hot-orders-section pt-16 pb-12 px-4">
        <!-- Green Wave Background -->
        <div class="red-wave-background">
          <svg viewBox="0 0 2000 800" preserveAspectRatio="none" xmlns="http://www.w3.org/2000/svg">
            <defs>
              <!-- Gradient 1 - Top wave -->
              <linearGradient id="waveGradient1" x1="0%" y1="0%" x2="0%" y2="100%">
                <stop offset="0%" style="stop-color:#16a34a;stop-opacity:0.5" />
                <stop offset="50%" style="stop-color:#15803d;stop-opacity:0.4" />
                <stop offset="100%" style="stop-color:#16a34a;stop-opacity:0.3" />
              </linearGradient>
              <!-- Gradient 2 - Middle wave -->
              <linearGradient id="waveGradient2" x1="0%" y1="0%" x2="0%" y2="100%">
                <stop offset="0%" style="stop-color:#15803d;stop-opacity:0.4" />
                <stop offset="50%" style="stop-color:#16a34a;stop-opacity:0.35" />
                <stop offset="100%" style="stop-color:#15803d;stop-opacity:0.25" />
              </linearGradient>
              <!-- Gradient 3 - Bottom wave -->
              <linearGradient id="waveGradient3" x1="0%" y1="0%" x2="0%" y2="100%">
                <stop offset="0%" style="stop-color:#16a34a;stop-opacity:0.3" />
                <stop offset="50%" style="stop-color:#15803d;stop-opacity:0.25" />
                <stop offset="100%" style="stop-color:#16a34a;stop-opacity:0.2" />
              </linearGradient>
              <!-- Gradient 4 - Additional layer -->
              <linearGradient id="waveGradient4" x1="0%" y1="0%" x2="0%" y2="100%">
                <stop offset="0%" style="stop-color:#15803d;stop-opacity:0.25" />
                <stop offset="100%" style="stop-color:#16a34a;stop-opacity:0.15" />
              </linearGradient>
              <!-- Gradient 5 - Deep layer -->
              <linearGradient id="waveGradient5" x1="0%" y1="0%" x2="0%" y2="100%">
                <stop offset="0%" style="stop-color:#16a34a;stop-opacity:0.2" />
                <stop offset="100%" style="stop-color:#15803d;stop-opacity:0.1" />
              </linearGradient>
            </defs>
            <!-- Layer 1 - Top wave with gradient -->
            <path d="M-200,150 C200,50 600,250 1000,150 C1300,100 1500,200 1800,150 C1900,130 1950,170 2000,150 L2000,800 L-200,800 Z" fill="url(#waveGradient1)"></path>
            <!-- Layer 2 - Second wave -->
            <path d="M-200,250 C300,150 700,350 1100,250 C1400,200 1600,300 1900,250 C1950,230 1980,270 2000,250 L2000,800 L-200,800 Z" fill="url(#waveGradient2)"></path>
            <!-- Layer 3 - Third wave -->
            <path d="M-200,350 C250,250 750,450 1050,350 C1350,300 1550,400 1850,350 C1900,330 1950,370 2000,350 L2000,800 L-200,800 Z" fill="url(#waveGradient3)"></path>
            <!-- Layer 4 - Additional depth -->
            <path d="M-200,450 C400,350 800,550 1200,450 C1500,400 1700,500 2000,450 L2000,800 L-200,800 Z" fill="url(#waveGradient4)"></path>
            <!-- Layer 5 - Deep background layer -->
            <path d="M-200,550 C350,450 850,650 1150,550 C1450,500 1650,600 2000,550 L2000,800 L-200,800 Z" fill="url(#waveGradient5)"></path>
          </svg>
        </div>
        <div class="max-w-7xl mx-auto relative z-10 bg-white rounded-lg" style="background: rgba(255, 255, 255, 0.95);">
          <div class="text-center mb-12 pt-8">
            <h2 class="text-4xl md:text-5xl lg:text-6xl font-bold text-gray-900 mb-4">
              <span class="text-gray-900">HOT</span>
              <span class="orders-text ml-3">ORDERS</span>
            </h2>
            <p class="text-gray-600 text-lg md:text-xl max-w-2xl mx-auto mb-8">
              Trending items that everyone's ordering right now
            </p>
            
            <!-- Category Options -->
            <div class="flex flex-wrap justify-center gap-4 mt-8 mb-12">
              <button 
                @click="selectCategory('food')" 
                :class="['category-btn', { 'category-btn-active': selectedCategory === 'food' }]"
              >
                <i class="fas fa-utensils mr-2"></i>
                Food
              </button>
              <button 
                @click="selectCategory('drinks')" 
                :class="['category-btn', { 'category-btn-active': selectedCategory === 'drinks' }]"
              >
                <i class="fas fa-coffee mr-2"></i>
                Drinks
              </button>
              <button 
                @click="selectCategory('mart')" 
                :class="['category-btn', { 'category-btn-active': selectedCategory === 'mart' }]"
              >
                <i class="fas fa-shopping-basket mr-2"></i>
                Mart
              </button>
              <button 
                @click="selectCategory('express')" 
                :class="['category-btn', { 'category-btn-active': selectedCategory === 'express' }]"
              >
                <i class="fas fa-bolt mr-2"></i>
                Express
              </button>
            </div>
          </div>

          <!-- Items Grid -->
          <transition-group name="category-switch" tag="div" :class="['food-grid', { 'initial-load': isInitialLoad }]">
            <div 
              v-for="(item, index) in currentItems" 
              :key="`${selectedCategory}-${index}`" 
              class="food-item" 
              :style="{ '--item-index': index }"
            >
              <div class="food-image">
                <img :src="item.image" :alt="item.name" />
              </div>
              <div class="food-details">
                <h3 class="food-name">{{ item.name }}</h3>
                <p class="food-description">{{ item.description }}</p>
                <p class="food-price">{{ item.price }}</p>
                <div class="food-actions">
                  <button class="btn-order">Order</button>
                  <button class="btn-review">See Review</button>
                </div>
              </div>
            </div>
          </transition-group>
        </div>
      </section>
    </main>

    <!-- Footer -->
    <footer class="footer">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-12">
        <div class="footer-content">
          <!-- Company Info -->
          <div class="footer-section">
            <img 
              src="/Images/grabfast-logo.png" 
              alt="FastGrab" 
              class="footer-logo"
            />
            <p class="footer-description">
              Fast delivery service for food, drinks, groceries, and packages. Order fast, get it faster.
            </p>
          </div>

          <!-- Quick Links -->
          <div class="footer-section">
            <h3 class="footer-heading">Quick Links</h3>
            <ul class="footer-links">
              <li><a href="#">About Us</a></li>
              <li><a href="#">How It Works</a></li>
              <li><a href="#">Contact</a></li>
              <li><a href="#">Careers</a></li>
            </ul>
          </div>

          <!-- Services -->
          <div class="footer-section">
            <h3 class="footer-heading">Services</h3>
            <ul class="footer-links">
              <li><a href="#">Food Delivery</a></li>
              <li><a href="#">Drinks</a></li>
              <li><a href="#">Grocery Mart</a></li>
              <li><a href="#">Express Delivery</a></li>
            </ul>
          </div>

          <!-- Support -->
          <div class="footer-section">
            <h3 class="footer-heading">Support</h3>
            <ul class="footer-links">
              <li><a href="#">Help Center</a></li>
              <li><a href="#">Track Order</a></li>
              <li><a href="#">FAQs</a></li>
              <li><a href="#">Terms & Conditions</a></li>
            </ul>
          </div>
        </div>

        <!-- Footer Bottom -->
        <div class="footer-bottom">
          <p class="footer-copyright">
            &copy; {{ new Date().getFullYear() }} FastGrab. All rights reserved.
          </p>
          <div class="footer-social">
            <a href="#" class="social-link">
              <i class="fab fa-facebook"></i>
            </a>
            <a href="#" class="social-link">
              <i class="fab fa-twitter"></i>
            </a>
            <a href="#" class="social-link">
              <i class="fab fa-instagram"></i>
            </a>
            <a href="#" class="social-link">
              <i class="fab fa-linkedin"></i>
            </a>
          </div>
        </div>
      </div>
    </footer>
  </div>
</template>

<script>
export default {
  name: 'App',
  data() {
    return {
      showMobileSearch: false,
      searchQuery: '',
      selectedCategory: 'food',
      isInitialLoad: true,
      categories: {
        food: [
          { name: 'Classic Burger', description: 'Juicy beef patty with fresh toppings', price: '$12.99', image: 'https://images.unsplash.com/photo-1568901346375-23c9450c58cd?w=400&h=400&fit=crop' },
          { name: 'Margherita Pizza', description: 'Fresh mozzarella and basil', price: '$15.99', image: 'https://images.unsplash.com/photo-1565299624946-b28f40a0ae38?w=400&h=400&fit=crop' },
          { name: 'Caesar Salad', description: 'Crisp romaine with parmesan', price: '$9.99', image: 'https://images.unsplash.com/photo-1546069901-ba9599a7e63c?w=400&h=400&fit=crop' },
          { name: 'Sushi Platter', description: 'Assorted fresh sushi rolls', price: '$18.99', image: 'https://images.unsplash.com/photo-1563379926898-05f4575a45d8?w=400&h=400&fit=crop' },
          { name: 'Pasta Carbonara', description: 'Creamy Italian classic', price: '$14.99', image: 'https://images.unsplash.com/photo-1555939594-58d7cb561ad1?w=400&h=400&fit=crop' },
          { name: 'Street Tacos', description: 'Authentic Mexican street food', price: '$11.99', image: 'https://images.unsplash.com/photo-1565958011703-44f9829ba187?w=400&h=400&fit=crop' },
          { name: 'Tonkotsu Ramen', description: 'Rich pork broth noodles', price: '$13.99', image: 'https://images.unsplash.com/photo-1567620905732-2d1ec7ab7445?w=400&h=400&fit=crop' },
          { name: 'Club Sandwich', description: 'Triple-decker delight', price: '$10.99', image: 'https://images.unsplash.com/photo-1606502281004-f86cf1282af5?w=400&h=400&fit=crop' },
          { name: 'Fried Chicken', description: 'Crispy golden perfection', price: '$12.99', image: 'https://images.unsplash.com/photo-1603133872878-684f208fb84b?w=400&h=400&fit=crop' }
        ],
        drinks: [
          { name: 'Fresh Orange Juice', description: 'Freshly squeezed oranges', price: '$5.99', image: 'https://images.unsplash.com/photo-1600271886742-f049cd451bba?w=400&h=400&fit=crop' },
          { name: 'Iced Coffee', description: 'Cold brew perfection', price: '$4.99', image: 'https://images.unsplash.com/photo-1517487881594-2787fef5ebf7?w=400&h=400&fit=crop' },
          { name: 'Smoothie Bowl', description: 'Mixed berry smoothie', price: '$7.99', image: 'https://images.unsplash.com/photo-1590301157890-4810ed352733?w=400&h=400&fit=crop' },
          { name: 'Matcha Latte', description: 'Japanese green tea latte', price: '$6.49', image: 'https://images.unsplash.com/photo-1536013453441-42467a254d66?w=400&h=400&fit=crop' },
          { name: 'Bubble Tea', description: 'Classic milk tea with tapioca', price: '$5.49', image: 'https://images.unsplash.com/photo-1525385133512-2f3bdd039054?w=400&h=400&fit=crop' },
          { name: 'Lemonade', description: 'Homemade fresh lemonade', price: '$3.99', image: 'https://images.unsplash.com/photo-1523677011781-c91d1bbe1f33?w=400&h=400&fit=crop' },
          { name: 'Cappuccino', description: 'Espresso with steamed milk', price: '$4.49', image: 'https://images.unsplash.com/photo-1572442388796-11668a67e53d?w=400&h=400&fit=crop' },
          { name: 'Protein Shake', description: 'Chocolate protein blend', price: '$6.99', image: 'https://images.unsplash.com/photo-1622483767028-3f66f32aef97?w=400&h=400&fit=crop' },
          { name: 'Green Juice', description: 'Detox veggie blend', price: '$7.49', image: 'https://images.unsplash.com/photo-1610970881699-44a5587cabec?w=400&h=400&fit=crop' }
        ],
        mart: [
          { name: 'Fresh Milk', description: 'Organic whole milk 1L', price: '$3.99', image: 'https://images.unsplash.com/photo-1563636619-e9143da7973b?w=400&h=400&fit=crop' },
          { name: 'Bread Loaf', description: 'Artisan whole wheat bread', price: '$4.49', image: 'https://images.unsplash.com/photo-1509440159596-0249088772ff?w=400&h=400&fit=crop' },
          { name: 'Fresh Eggs', description: 'Free-range eggs dozen', price: '$5.99', image: 'https://images.unsplash.com/photo-1518569656558-1f25e69d93d7?w=400&h=400&fit=crop' },
          { name: 'Bananas', description: 'Organic ripe bananas', price: '$2.99', image: 'https://images.unsplash.com/photo-1571771894821-ce9b6c11b08e?w=400&h=400&fit=crop' },
          { name: 'Tomatoes', description: 'Vine-ripened tomatoes', price: '$4.29', image: 'https://images.unsplash.com/photo-1546094096-0df4bcaaa337?w=400&h=400&fit=crop' },
          { name: 'Cheese Block', description: 'Aged cheddar cheese', price: '$6.99', image: 'https://images.unsplash.com/photo-1486297678162-eb2a19b0a32d?w=400&h=400&fit=crop' },
          { name: 'Pasta Pack', description: 'Italian spaghetti 500g', price: '$3.49', image: 'https://images.unsplash.com/photo-1621996346565-e3dbc646d9a9?w=400&h=400&fit=crop' },
          { name: 'Olive Oil', description: 'Extra virgin olive oil', price: '$8.99', image: 'https://images.unsplash.com/photo-1474979266404-7eaacbcd87c5?w=400&h=400&fit=crop' },
          { name: 'Coffee Beans', description: 'Premium arabica beans', price: '$12.99', image: 'https://images.unsplash.com/photo-1559056199-641a0ac8b55e?w=400&h=400&fit=crop' }
        ],
        express: [
          { name: 'Same Day Delivery', description: 'Express package within 24hrs', price: '$8.99', image: 'https://images.unsplash.com/photo-1566576721346-d4a3b4eaeb55?w=400&h=400&fit=crop' },
          { name: 'Next Day Shipping', description: 'Standard overnight delivery', price: '$12.99', image: 'https://images.unsplash.com/photo-1586528116311-ad8dd3c8310d?w=400&h=400&fit=crop' },
          { name: 'Priority Express', description: '2-3 hour rush delivery', price: '$24.99', image: 'https://images.unsplash.com/photo-1605902394069-ff2ae2430e62?w=400&h=400&fit=crop' },
          { name: 'International Shipping', description: 'Global package delivery', price: '$45.99', image: 'https://images.unsplash.com/photo-1578575437130-527eed3abbec?w=400&h=400&fit=crop' },
          { name: 'Bulk Packages', description: 'Multiple items delivery', price: '$19.99', image: 'https://images.unsplash.com/photo-1566576912321-d58ddd7a6088?w=400&h=400&fit=crop' },
          { name: 'Document Courier', description: 'Secure document delivery', price: '$6.99', image: 'https://images.unsplash.com/photo-1586339949916-3e9457bef6d3?w=400&h=400&fit=crop' },
          { name: 'Fragile Item Care', description: 'Special handling package', price: '$15.99', image: 'https://images.unsplash.com/photo-1607166452427-7e4477079cb9?w=400&h=400&fit=crop' },
          { name: 'Large Item Delivery', description: 'Oversized package service', price: '$29.99', image: 'https://images.unsplash.com/photo-1494412574643-ff11b0a5c1c3?w=400&h=400&fit=crop' },
          { name: 'Scheduled Pickup', description: 'Book a pickup time slot', price: '$4.99', image: 'https://images.unsplash.com/photo-1600880292203-757bb62b4baf?w=400&h=400&fit=crop' }
        ]
      }
    }
  },
  computed: {
    currentItems() {
      return this.categories[this.selectedCategory]
    }
  },
  methods: {
    toggleMobileSearch() {
      this.showMobileSearch = !this.showMobileSearch
      if (this.showMobileSearch) {
        this.$nextTick(() => {
          this.$refs.mobileSearchInput?.focus()
        })
      }
    },
    closeMobileSearch() {
      this.showMobileSearch = false
      this.searchQuery = ''
    },
    selectCategory(category) {
      if (this.isInitialLoad) {
        this.isInitialLoad = false
      }
      this.selectedCategory = category
    }
  }
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

@keyframes slideDownSearch {
  from {
    transform: translateY(-100%);
    opacity: 0;
    max-height: 0;
  }
  to {
    transform: translateY(0);
    opacity: 1;
    max-height: 200px;
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

/* Hide search bar and show magnifying glass at 800px and below */
@media (max-width: 800px) {
  .search-bar-container {
    display: none;
  }
  
  .search-icon-mobile {
    display: block;
  }
  
  .hero-title {
    gap: 1rem 1.25rem;
  }
  
  .hero-word,
  .hero-word-fast {
    font-size: clamp(2rem, 6vw, 3rem);
  }
}

/* Show search bar and hide magnifying glass above 800px */
@media (min-width: 801px) {
  .search-bar-container {
    display: block;
  }
  
  .search-icon-mobile {
    display: none;
  }
  
  .mobile-search-container {
    display: none;
  }
}

.mobile-search-container {
  animation: slideDownSearch 0.4s ease-out;
  position: sticky;
  top: 64px;
  z-index: 40;
}

/* Wave Divider */
.wave-divider {
  width: 100%;
  overflow: hidden;
  line-height: 0;
  position: relative;
  margin-top: -1px;
}

.wave-divider svg {
  position: relative;
  display: block;
  width: calc(100% + 1.3px);
  height: 80px;
}

@media (min-width: 768px) {
  .wave-divider svg {
    height: 120px;
  }
}

@keyframes waveFloat {
  0%, 100% {
    transform: translateX(0) translateY(0);
  }
  33% {
    transform: translateX(-20px) translateY(10px);
  }
  66% {
    transform: translateX(10px) translateY(-10px);
  }
}

@keyframes wavePulse {
  0% {
    opacity: 0.2;
  }
  100% {
    opacity: 0.4;
  }
}

/* Hero Section Animations */
@keyframes speedIn {
  0% {
    transform: translateX(-100px);
    opacity: 0;
    filter: blur(10px);
  }
  60% {
    transform: translateX(10px);
  }
  100% {
    transform: translateX(0);
    opacity: 1;
    filter: blur(0);
  }
}

@keyframes fastZoom {
  0% {
    transform: scale(0.5) translateX(-200px);
    opacity: 0;
    filter: blur(15px);
  }
  50% {
    transform: scale(1.1) translateX(5px);
    filter: blur(2px);
  }
  100% {
    transform: scale(1) translateX(0);
    opacity: 1;
    filter: blur(0);
  }
}

@keyframes shake {
  0%, 100% {
    transform: translateX(0);
  }
  10%, 30%, 50%, 70%, 90% {
    transform: translateX(-2px);
  }
  20%, 40%, 60%, 80% {
    transform: translateX(2px);
  }
}

.hero-title {
  font-size: clamp(2.5rem, 8vw, 7rem);
  font-weight: 900;
  line-height: 1.1;
  letter-spacing: 0.02em;
  display: flex;
  flex-wrap: wrap;
  justify-content: center;
  align-items: center;
  gap: 1.5rem 2rem;
  padding: 2rem 0;
}

.hero-word {
  display: inline-block;
  color: #1f2937;
  position: relative;
  animation: speedIn 0.8s ease-out forwards;
  animation-delay: calc(var(--word-index) * 0.15s);
  opacity: 0;
  transition: transform 0.3s ease;
}

.hero-word:hover {
  transform: translateY(-4px);
}

.hero-word:nth-child(1) {
  --word-index: 0;
}

.hero-word:nth-child(2) {
  --word-index: 1;
}

.hero-word:nth-child(3) {
  --word-index: 2;
}

.hero-word-fast {
  display: inline-block;
  color: #16a34a;
  font-weight: 900;
  font-style: italic;
  text-shadow: 2px 2px 4px rgba(22, 163, 74, 0.2);
  text-decoration: underline;
  text-decoration-color: #16a34a;
  text-decoration-thickness: 4px;
  text-underline-offset: 8px;
  animation: fastZoom 1s ease-out forwards, shake 0.5s ease-in-out 1.2s;
  animation-delay: 0.45s, 1.2s;
  opacity: 0;
  position: relative;
  transition: transform 0.3s ease;
}

.hero-word-fast:hover {
  transform: translateY(-4px) scale(1.05);
}

/* Hero Buttons */
.hero-buttons {
  display: flex;
  gap: 1.5rem;
  justify-content: center;
  flex-wrap: wrap;
  opacity: 0;
  animation: fadeInUp 0.8s ease-out 1.5s forwards;
}

.btn-hero {
  padding: 1rem 2.5rem;
  font-size: 1.125rem;
  font-weight: 600;
  border: none;
  border-radius: 12px;
  cursor: pointer;
  transition: all 0.3s ease;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
  letter-spacing: 0.025em;
}

.btn-orders {
  background: #16a34a;
  color: white;
}

.btn-orders:hover {
  background: #15803d;
  transform: translateY(-2px);
  box-shadow: 0 6px 16px rgba(22, 163, 74, 0.4);
}

.btn-download {
  background: #D02752;
  color: white;
}

.btn-download:hover {
  background: #b02145;
  transform: translateY(-2px);
  box-shadow: 0 6px 16px rgba(208, 39, 82, 0.4);
}

@keyframes speedLine {
  0% {
    transform: translateX(-100%) translateY(-50%);
    opacity: 1;
  }
  100% {
    transform: translateX(100%) translateY(-50%);
    opacity: 0;
  }
}

/* Hot Orders Section */
.hot-orders-section {
  position: relative;
  background: transparent;
  overflow: hidden;
}

.red-wave-background {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  z-index: 0;
  pointer-events: none;
  overflow: hidden;
}

.red-wave-background svg {
  width: 120%;
  height: 100%;
  margin-left: -10%;
  animation: waveFloat 10s ease-in-out infinite;
}

.red-wave-background svg path:nth-child(1) {
  animation: wavePulse 6s ease-in-out infinite alternate;
}

.red-wave-background svg path:nth-child(2) {
  animation: wavePulse 6s ease-in-out infinite alternate;
  animation-delay: 0.8s;
}

.red-wave-background svg path:nth-child(3) {
  animation: wavePulse 6s ease-in-out infinite alternate;
  animation-delay: 1.6s;
}

.red-wave-background svg path:nth-child(4) {
  animation: wavePulse 6s ease-in-out infinite alternate;
  animation-delay: 2.4s;
}

.red-wave-background svg path:nth-child(5) {
  animation: wavePulse 6s ease-in-out infinite alternate;
  animation-delay: 3.2s;
}

.hot-orders-section h2 {
  position: relative;
  display: inline-block;
  opacity: 0;
  animation: titleBounceIn 0.8s ease-out 0.6s forwards;
}

.hot-orders-section p {
  opacity: 0;
  animation: fadeInUp 0.6s ease-out 0.9s forwards;
}

.hot-orders-section .flex {
  opacity: 0;
  animation: fadeInUp 0.6s ease-out 1.1s forwards;
}

@keyframes titleBounceIn {
  0% {
    opacity: 0;
    transform: scale(0.3) translateY(-30px);
  }
  50% {
    transform: scale(1.05);
  }
  70% {
    transform: scale(0.97);
  }
  100% {
    opacity: 1;
    transform: scale(1) translateY(0);
  }
}

.orders-text {
  color: #D02752;
  font-style: italic;
  text-decoration: underline;
  text-decoration-color: #D02752;
  text-decoration-thickness: 4px;
  text-underline-offset: 4px;
}

.category-btn {
  padding: 0.875rem 2.5rem;
  background: transparent;
  border: none;
  border-bottom: 3px solid transparent;
  font-size: 1.125rem;
  font-weight: 500;
  color: #6b7280;
  cursor: pointer;
  transition: all 0.3s ease;
  letter-spacing: 0.025em;
}

.category-btn:hover {
  color: #16a34a;
  border-bottom-color: #16a34a;
}

.category-btn:active {
  color: #15803d;
}

.category-btn-active {
  color: #16a34a;
  border-bottom-color: #16a34a;
  font-weight: 600;
}

.category-btn i {
  color: #9ca3af;
  transition: color 0.3s ease;
}

.category-btn:hover i,
.category-btn-active i {
  color: #16a34a;
}

/* Food Grid */
.food-grid {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 1.5rem;
  max-width: 1200px;
  margin: 0 auto;
  position: relative;
}

/* Category Switch Animations */
.category-switch-enter-active {
  transition: all 0.5s ease-out;
}

.category-switch-leave-active {
  transition: all 0.4s ease-in;
  position: absolute;
  width: calc(33.333% - 1rem);
}

.category-switch-enter-from {
  opacity: 0;
  transform: translateY(30px) scale(0.9);
}

.category-switch-enter-to {
  opacity: 1;
  transform: translateY(0) scale(1);
}

.category-switch-leave-from {
  opacity: 1;
  transform: translateY(0) scale(1);
}

.category-switch-leave-to {
  opacity: 0;
  transform: translateY(-20px) scale(0.95);
}

@media (max-width: 1024px) {
  .category-switch-leave-active {
    width: calc(50% - 0.75rem);
  }
}

@media (max-width: 640px) {
  .category-switch-leave-active {
    width: 100%;
  }
}

@media (max-width: 1024px) {
  .food-grid {
    grid-template-columns: repeat(2, 1fr);
  }
}

@media (max-width: 640px) {
  .food-grid {
    grid-template-columns: 1fr;
  }
}

.food-item {
  background: white;
  border-radius: 12px;
  overflow: hidden;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.08);
  transition: all 0.3s ease;
  cursor: pointer;
}

/* Initial load animation */
.food-grid.initial-load .food-item {
  opacity: 0;
  animation: foodItemAppear 0.6s ease-out forwards;
  animation-delay: calc(var(--item-index) * 0.1s + 1s);
}

.food-item:hover {
  transform: translateY(-4px);
  box-shadow: 0 8px 16px rgba(0, 0, 0, 0.12);
}

@keyframes foodItemAppear {
  0% {
    opacity: 0;
    transform: translateY(30px) scale(0.95);
  }
  100% {
    opacity: 1;
    transform: translateY(0) scale(1);
  }
}

.food-image {
  width: 100%;
  height: 180px;
  overflow: hidden;
  background: #f3f4f6;
}

.food-image img {
  width: 100%;
  height: 100%;
  object-fit: cover;
  transition: transform 0.3s ease;
}

.food-item:hover .food-image img {
  transform: scale(1.05);
}

.food-details {
  padding: 1rem;
}

.food-name {
  font-size: 1.125rem;
  font-weight: 600;
  color: #1f2937;
  margin-bottom: 0.375rem;
}

.food-description {
  font-size: 0.8125rem;
  color: #6b7280;
  margin-bottom: 0.625rem;
  line-height: 1.5;
}

.food-price {
  font-size: 1rem;
  font-weight: 700;
  color: #16a34a;
  margin-bottom: 0.75rem;
}

.food-actions {
  display: flex;
  gap: 0.5rem;
}

.btn-order,
.btn-review {
  flex: 1;
  padding: 0.5rem 0.75rem;
  border: none;
  border-radius: 6px;
  font-size: 0.8125rem;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.3s ease;
}

.btn-order {
  background: #16a34a;
  color: white;
}

.btn-order:hover {
  background: #15803d;
  transform: translateY(-2px);
  box-shadow: 0 4px 8px rgba(22, 163, 74, 0.3);
}

.btn-review {
  background: transparent;
  color: #6b7280;
  border: 1px solid #d1d5db;
}

.btn-review:hover {
  color: #374151;
  border-color: #9ca3af;
  background: #f9fafb;
}

/* Footer */
.footer {
  background: #1f2937;
  color: #e5e7eb;
  margin-top: 0;
}

.footer-content {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
  gap: 3rem;
  margin-bottom: 3rem;
}

.footer-section {
  display: flex;
  flex-direction: column;
}

.footer-logo {
  height: 40px;
  margin-bottom: 1rem;
  filter: brightness(0) invert(1);
}

.footer-description {
  color: #9ca3af;
  font-size: 0.875rem;
  line-height: 1.6;
  max-width: 280px;
}

.footer-heading {
  font-size: 1.125rem;
  font-weight: 600;
  color: white;
  margin-bottom: 1rem;
}

.footer-links {
  list-style: none;
  padding: 0;
  margin: 0;
}

.footer-links li {
  margin-bottom: 0.75rem;
}

.footer-links a {
  color: #9ca3af;
  text-decoration: none;
  font-size: 0.875rem;
  transition: color 0.3s ease;
}

.footer-links a:hover {
  color: #16a34a;
}

.footer-bottom {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding-top: 2rem;
  border-top: 1px solid #374151;
  flex-wrap: wrap;
  gap: 1rem;
}

.footer-copyright {
  color: #9ca3af;
  font-size: 0.875rem;
  margin: 0;
}

.footer-social {
  display: flex;
  gap: 1rem;
}

.social-link {
  width: 36px;
  height: 36px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: #374151;
  border-radius: 50%;
  color: #9ca3af;
  text-decoration: none;
  transition: all 0.3s ease;
}

.social-link:hover {
  background: #16a34a;
  color: white;
  transform: translateY(-2px);
}

@media (max-width: 768px) {
  .footer-content {
    grid-template-columns: 1fr;
    gap: 2rem;
  }

  .footer-bottom {
    flex-direction: column;
    text-align: center;
  }
}
</style>
