import { createApp } from 'vue'
import App from './App.vue'

// Import CSS files statically so Tailwind Vite plugin can process them
import './style.css'

const app = createApp(App)
app.mount('#app')
