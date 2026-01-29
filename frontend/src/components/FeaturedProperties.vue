<template>
  <div class="bg-white rounded-lg border border-[#E8E3D9] p-6">
    <h3 class="text-2xl font-light text-[#3A3632] mb-6 tracking-[0.05em]" style="font-family: 'Noto Sans', sans-serif;">Featured Properties</h3>
    <div v-if="loading" class="text-center py-8">
      <p class="text-[#3A3632] text-sm" style="font-family: 'Noto Sans', sans-serif;">Loading properties...</p>
    </div>
    <div v-else-if="properties.length === 0" class="text-center py-8">
      <p class="text-[#3A3632] text-sm" style="font-family: 'Noto Sans', sans-serif;">No properties available at this time.</p>
    </div>
    <div v-else class="space-y-6">
      <PropertyCard
        v-for="property in properties.slice(0, limit)" 
        :key="property.address || property.id || property.listing_id"
        :property="property"
      />
    </div>
    <div class="mt-6 text-center">
      <a href="/search" class="inline-block text-brand-blue font-medium hover:text-brand-blue/80 transition-colors duration-300 text-sm" style="font-family: 'Noto Sans', sans-serif;">
        View All Properties →
      </a>
    </div>
  </div>
</template>

<script>
import PropertyCard from './PropertyCard.vue'

export default {
  name: 'FeaturedProperties',
  components: {
    PropertyCard
  },
  props: {
    properties: {
      type: Array,
      default: () => []
    },
    loading: {
      type: Boolean,
      default: false
    },
    limit: {
      type: Number,
      default: 3
    }
  }
}
</script>

