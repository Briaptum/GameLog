<template>
  <a 
    :href="propertyLink"
    class="block bg-white rounded-lg border border-[#E8E3D9] shadow-md hover:shadow-xl hover:-translate-y-1 transition-all duration-300 overflow-hidden no-underline"
  >
    <div class="h-80 overflow-hidden relative">
      <img 
        :src="propertyImage" 
        :alt="propertyAddress"
        class="w-full h-full object-cover transition-transform duration-300 hover:scale-105"
        @error="handleImageError"
      />
      <div v-if="propertyStatus" class="absolute top-3 right-3">
        <span class="bg-brand-dark-gold text-white px-3 py-1.5 rounded text-xs font-medium tracking-[0.1em] uppercase shadow-lg" style="font-family: 'Noto Sans', sans-serif;">{{ propertyStatus }}</span>
      </div>
      <!-- Price overlay on image -->
      <div class="absolute bottom-4 left-4 right-4">
        <div class="bg-black/70 backdrop-blur-sm rounded-lg px-3 py-1.5 inline-block">
          <p class="text-xl font-semibold text-white tracking-wide" style="font-family: 'Noto Sans', sans-serif;">{{ formattedPrice }}</p>
        </div>
      </div>
    </div>
    <div class="p-5">
      <p v-if="propertyLocation" class="text-xs text-[#A89F91] mb-1.5 tracking-[0.1em] uppercase" style="font-family: 'Noto Sans', sans-serif;">{{ propertyLocation }}</p>
      <h3 class="font-medium text-[#3A3632] mb-3 tracking-wide" style="font-family: 'Noto Sans', sans-serif;">{{ propertyAddress }}</h3>
      <div class="flex justify-between pt-3 border-t border-[#E8E3D9] text-sm text-[#A89F91] tracking-wide" style="font-family: 'Noto Sans', sans-serif;">
        <span v-if="propertyBeds !== null && propertyBeds !== undefined" class="font-medium flex items-center gap-1.5">
          <font-awesome-icon icon="bed" class="text-[#C9B8A8] text-xs" />
          <strong class="text-[#3A3632]">{{ propertyBeds }}</strong> Beds
        </span>
        <span v-if="propertyBaths !== null && propertyBaths !== undefined" class="font-medium flex items-center gap-1.5">
          <font-awesome-icon icon="bath" class="text-[#C9B8A8] text-xs" />
          <strong class="text-[#3A3632]">{{ formatBaths(propertyBaths) }}</strong> Baths
        </span>
        <span v-if="propertySqft !== null && propertySqft !== undefined" class="font-medium">
          <strong class="text-[#3A3632]">{{ formatSquareFeet(propertySqft) }}</strong> sq ft
        </span>
      </div>
    </div>
  </a>
</template>

<script>
export default {
  name: 'PropertyCard',
  props: {
    property: {
      type: Object,
      required: true
    },
    link: {
      type: String,
      default: '/search'
    }
  },
  computed: {
    propertyImage() {
      return this.property.image || this.property.primary_photo_url || ''
    },
    propertyAddress() {
      return this.property.address || 'Address not available'
    },
    propertyLocation() {
      return this.property.location || null
    },
    formattedPrice() {
      if (this.property.price) {
        // If it's already a formatted string, return it
        if (typeof this.property.price === 'string') {
          return this.property.price
        }
        // Otherwise format it
        const price = typeof this.property.price === 'number' ? this.property.price : parseFloat(this.property.price)
        if (!isNaN(price)) {
          return new Intl.NumberFormat('en-US', {
            style: 'currency',
            currency: 'USD',
            minimumFractionDigits: 0,
            maximumFractionDigits: 0
          }).format(price)
        }
      }
      return this.property.list_price ? this.formatPrice(this.property.list_price) : 'Price on request'
    },
    propertyStatus() {
      return this.property.status || null
    },
    propertyBeds() {
      return this.property.beds !== undefined ? this.property.beds : (this.property.bedrooms !== undefined ? this.property.bedrooms : null)
    },
    propertyBaths() {
      return this.property.baths !== undefined ? this.property.baths : (this.property.bathrooms !== undefined ? this.property.bathrooms : null)
    },
    propertySqft() {
      return this.property.sqft !== undefined ? this.property.sqft : (this.property.square_feet !== undefined ? this.property.square_feet : (this.property.living_area !== undefined ? this.property.living_area : null))
    },
    propertyLink() {
      if (this.property.listing_id) {
        return `/property/${this.property.listing_id}`
      }
      return this.link
    }
  },
  methods: {
    formatPrice(price) {
      const priceNum = typeof price === 'string' ? parseFloat(price) : price
      if (isNaN(priceNum)) return 'Price on request'
      return new Intl.NumberFormat('en-US', {
        style: 'currency',
        currency: 'USD',
        minimumFractionDigits: 0,
        maximumFractionDigits: 0
      }).format(priceNum)
    },
    formatBaths(baths) {
      if (baths === null || baths === undefined) return null
      const bathsNum = typeof baths === 'string' ? parseFloat(baths) : baths
      if (isNaN(bathsNum)) return null
      // If it's a whole number, show without decimal
      if (bathsNum % 1 === 0) {
        return bathsNum.toString()
      }
      return bathsNum.toFixed(1)
    },
    formatSquareFeet(sqft) {
      if (!sqft && sqft !== 0) return ''
      const num = typeof sqft === 'string' ? parseFloat(sqft) : sqft
      if (isNaN(num)) return ''
      return new Intl.NumberFormat('en-US').format(num)
    },
    handleImageError(event) {
      console.error('Property image failed to load:', event.target.src)
      event.target.style.display = 'none'
      event.target.parentElement.classList.add('bg-gradient-to-br', 'from-[#E8E3D9]', 'to-[#D4C9B8]', 'flex', 'items-center', 'justify-center')
      event.target.parentElement.innerHTML = '<span class="text-[#A89F91] text-sm tracking-[0.05em]">No Image</span>'
    }
  }
}
</script>

