<template>
  <footer class="bg-zinc-900 text-[#E8E3D9] py-16 px-4">
    <div class="container mx-auto px-4">
      <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-12 mb-12">
        <div class="mb-4">
          <h4 class="font-medium mb-6 text-[#FAF7F2] tracking-wide" style="font-family: 'Cinzel', serif;">Contact</h4>
          <p class="text-[#E8E3D9] mb-3 text-[0.9rem] opacity-90">
            <a href="tel:+1-502-238-7400" class="text-[#E8E3D9] hover:text-[#FAF7F2] transition-colors duration-300">(502) 238-7400</a>
          </p>
          <p class="text-[#E8E3D9] mb-3 text-[0.9rem] opacity-90">8223 Shelbyville Rd.<br />Louisville, KY, 40222</p>
          <p class="text-[#E8E3D9] mb-4 text-[0.9rem] opacity-90">
            <a href="mailto:info@website-dummy.com" class="text-[#E8E3D9] hover:text-[#FAF7F2] transition-colors duration-300">info@website-dummy.com</a>
          </p>
          <div class="flex gap-4 mt-6">
            <a href="https://www.facebook.com/website-dummy-properties" target="_blank" rel="noopener noreferrer" class="text-[#E8E3D9] text-xl w-10 h-10 flex items-center justify-center rounded-full border border-[#A89F91] hover:text-[#FAF7F2] hover:border-[#164268] hover:bg-[#164268] transition-all duration-300" aria-label="Facebook">
              <font-awesome-icon :icon="['fab', 'facebook-f']" />
            </a>
            <a href="https://www.instagram.com/website-dummy-properties/" target="_blank" rel="noopener noreferrer" class="text-[#E8E3D9] text-xl w-10 h-10 flex items-center justify-center rounded-full border border-[#A89F91] hover:text-[#FAF7F2] hover:border-[#164268] hover:bg-[#164268] transition-all duration-300" aria-label="Instagram">
              <font-awesome-icon :icon="['fab', 'instagram']" />
            </a>
            <a href="https://www.linkedin.com/company/julie-pogue-properties/" target="_blank" rel="noopener noreferrer" class="text-[#E8E3D9] text-xl w-10 h-10 flex items-center justify-center rounded-full border border-[#A89F91] hover:text-[#FAF7F2] hover:border-[#164268] hover:bg-[#164268] transition-all duration-300" aria-label="LinkedIn">
              <font-awesome-icon :icon="['fab', 'linkedin']" />
            </a>
            <a href="https://share.google/XaJQzbw3k0q99A90Y" target="_blank" rel="noopener noreferrer" class="text-[#E8E3D9] text-xl w-10 h-10 flex items-center justify-center rounded-full border border-[#A89F91] hover:text-[#FAF7F2] hover:border-[#164268] hover:bg-[#164268] transition-all duration-300" aria-label="Google Places">
              <font-awesome-icon :icon="['fab', 'google']" />
            </a>
          </div>
        </div>
        
        <div class="mb-4 lg:col-span-2">
          <h4 class="font-medium mb-6 text-[#FAF7F2] tracking-wide" style="font-family: 'Cinzel', serif;">Communities</h4>
          <div class="grid grid-cols-1 md:grid-cols-2 gap-x-8 gap-y-3">
            <ul class="list-none p-0 m-0">
              <li 
                v-for="(neighborhood, index) in neighborhoods.slice(0, 6)" 
                :key="neighborhood.name"
                class="text-[#E8E3D9] mb-3 text-[0.9rem] opacity-90 hover:opacity-100 hover:text-[#FAF7F2] transition-all duration-300"
              >
                <a :href="`/${neighborhood.slug}`" class="text-[#E8E3D9] hover:text-[#FAF7F2] transition-colors duration-300">
                  {{ neighborhood.name }}
                </a>
              </li>
            </ul>
            <ul class="list-none p-0 m-0">
              <li 
                v-for="(neighborhood, index) in neighborhoods.slice(6)" 
                :key="neighborhood.name"
                class="text-[#E8E3D9] mb-3 text-[0.9rem] opacity-90 hover:opacity-100 hover:text-[#FAF7F2] transition-all duration-300"
              >
                <a :href="`/${neighborhood.slug}`" class="text-[#E8E3D9] hover:text-[#FAF7F2] transition-colors duration-300">
                  {{ neighborhood.name }}
                </a>
              </li>
            </ul>
          </div>
        </div>
        
        <div class="mb-4">
          <h4 class="font-medium mb-6 text-[#FAF7F2] tracking-wide" style="font-family: 'Cinzel', serif;">Quick Links</h4>
          <ul class="list-none p-0 m-0 space-y-3">
            <li>
              <a href="https://payments.earnnest.com/website-dummy/send/103043" class="inline-block bg-[#164268] text-white px-6 py-2.5 rounded text-sm font-medium no-underline hover:bg-[#0d2a42] hover:shadow-lg transition-all duration-300 tracking-wide" target="_blank" rel="noopener">
                Earnest Money Deposit
              </a>
            </li>
          </ul>
        </div>
      </div>
      <div class="border-t border-[#A89F91] pt-8 text-center text-[#E8E3D9] text-sm opacity-80">
        <p>&copy; 2025 Website Dummy Properties. All rights reserved.</p>
      </div>
    </div>
  </footer>
</template>

<script>
export default {
  name: 'Footer',
  data() {
    return {
      neighborhoods: [
        { name: 'Anchorage Homes', slug: 'anchorage-homes-for-sale' },
        { name: 'Audubon Park Homes', slug: 'audubon-park-homes-for-sale' },
        { name: 'Cherokee Triangle Homes', slug: 'cherokee-triangle-homes-for-sale' },
        { name: 'Crescent Hill Homes', slug: 'crescent-hill-homes-for-sale' },
        { name: 'Forest Springs Homes', slug: 'forest-springs-homes-for-sale' },
        { name: 'Highlands Homes', slug: 'highlands-homes-for-sale' },
        { name: 'Hunting Creek Homes', slug: 'hunting-creek-homes-for-sale' },
        { name: 'Indian Hills Homes', slug: 'indian-hills-homes-for-sale' },
        { name: 'Lake Forest Homes', slug: 'lake-forest-homes-for-sale' },
        { name: 'Mockingbird Valley Homes', slug: 'mockingbird-valley-homes-for-sale' },
        { name: 'Norton Commons Homes', slug: 'norton-commons-homes-for-sale' },
        { name: 'Owl Creek Homes', slug: 'owl-creek-homes-for-sale' },
        { name: 'Polo Fields Homes', slug: 'polo-fields-homes-for-sale' }
      ]
    }
  }
}
</script>

