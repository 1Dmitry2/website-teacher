import { createApp } from 'vue'
import "./styles/style.css"
import { router } from "./router/router.ts";
import App from './App.vue'
const app = createApp(App)
app.use(router)
app.mount('#app')
