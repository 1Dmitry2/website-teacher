import { createRouter, createWebHistory } from 'vue-router';
import NewsPage from "../pages/NewsPage.vue";
import AdminPage from "../pages/AdminPage.vue";
import AdminLoginPage from "../pages/AdminLoginPage.vue";


const routes = [
    { path: '/', component: NewsPage },
    { path: '/admin', component: AdminPage },
    { path: '/admin-login', component: AdminLoginPage },
]

export const router = createRouter({
    history: createWebHistory(),
    routes,
})
