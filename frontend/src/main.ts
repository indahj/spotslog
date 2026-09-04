import './styles.css'

import { createApp } from 'vue'
import { createPinia } from 'pinia'

import App from './App.vue'
import router from './router'
import { useAuthStore } from './stores/auth'
import { FontAwesomeIcon } from '@fortawesome/vue-fontawesome'
import "./icons"

const app = createApp(App)

app.use(createPinia())
app.use(router)
app.component("FontAwesomeIcon", FontAwesomeIcon)

// Restore the session before mounting
// so guards and the nav bar see the real auth state on first paint instead of flashing a logged-out shell.
const auth = useAuthStore()
auth.restore().finally(() => app.mount("#app"))

