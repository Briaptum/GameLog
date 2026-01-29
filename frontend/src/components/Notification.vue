<template>
  <Teleport to="body">
    <Transition name="notification">
      <div
        v-if="visible"
        :class="[
          'fixed top-24 right-4 z-[1100] max-w-md w-full shadow-lg rounded-lg p-4 flex items-start gap-3',
          type === 'success' ? 'bg-green-50 border border-green-200' : 'bg-red-50 border border-red-200'
        ]"
      >
        <div :class="type === 'success' ? 'text-green-600' : 'text-red-600'">
          <svg v-if="type === 'success'" class="w-5 h-5" fill="currentColor" viewBox="0 0 20 20">
            <path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z" clip-rule="evenodd" />
          </svg>
          <svg v-else class="w-5 h-5" fill="currentColor" viewBox="0 0 20 20">
            <path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zM8.707 7.293a1 1 0 00-1.414 1.414L8.586 10l-1.293 1.293a1 1 0 101.414 1.414L10 11.414l1.293 1.293a1 1 0 001.414-1.414L11.414 10l1.293-1.293a1 1 0 00-1.414-1.414L10 8.586 8.707 7.293z" clip-rule="evenodd" />
          </svg>
        </div>
        <div class="flex-1">
          <p :class="type === 'success' ? 'text-green-800' : 'text-red-800'" class="text-sm font-medium">
            {{ message }}
          </p>
        </div>
        <button
          @click="close"
          :class="type === 'success' ? 'text-green-600 hover:text-green-800' : 'text-red-600 hover:text-red-800'"
          class="flex-shrink-0"
        >
          <svg class="w-4 h-4" fill="currentColor" viewBox="0 0 20 20">
            <path fill-rule="evenodd" d="M4.293 4.293a1 1 0 011.414 0L10 8.586l4.293-4.293a1 1 0 111.414 1.414L11.414 10l4.293 4.293a1 1 0 01-1.414 1.414L10 11.414l-4.293 4.293a1 1 0 01-1.414-1.414L8.586 10 4.293 5.707a1 1 0 010-1.414z" clip-rule="evenodd" />
          </svg>
        </button>
      </div>
    </Transition>
  </Teleport>
</template>

<script>
export default {
  name: 'Notification',
  emits: ['close'],
  props: {
    message: {
      type: String,
      required: true
    },
    type: {
      type: String,
      default: 'success',
      validator: (value) => ['success', 'error'].includes(value)
    },
    duration: {
      type: Number,
      default: 3000
    }
  },
  data() {
    return {
      visible: false,
      timeout: null
    }
  },
  mounted() {
    this.visible = true
    if (this.duration > 0) {
      this.timeout = setTimeout(() => {
        this.close()
      }, this.duration)
    }
  },
  beforeUnmount() {
    if (this.timeout) {
      clearTimeout(this.timeout)
    }
  },
  methods: {
    close() {
      this.visible = false
      if (this.timeout) {
        clearTimeout(this.timeout)
      }
      this.$emit('close')
    }
  }
}
</script>

<style scoped>
.notification-enter-active,
.notification-leave-active {
  transition: all 0.3s ease;
}

.notification-enter-from {
  opacity: 0;
  transform: translateX(100%);
}

.notification-leave-to {
  opacity: 0;
  transform: translateX(100%);
}
</style>

