import { createRouter, createWebHistory } from 'vue-router';
import NewsPage from "../pages/NewsPage.vue";
import AdminPage from "../pages/AdminPage.vue";
import AdminLoginPage from "../pages/AdminLoginPage.vue";
import UserLogin from "../pages/UserLogin.vue";
import RegisterPage from "../pages/RegisterPage.vue";
import { authService } from "../utils/auth";


const routes = [
    { 
        path: '/', 
        component: NewsPage,
        meta: { requiresAuth: false }
    },
    { 
        path: '/admin', 
        component: AdminPage,
        meta: { requiresAuth: true }
    },
    { 
        path: '/admin-login', 
        component: AdminLoginPage,
        meta: { requiresAuth: false }
    },
    { 
        path: '/user-login', 
        component: UserLogin,
        meta: { requiresAuth: false }
    },
    { 
        path: '/register', 
        component: RegisterPage,
        meta: { requiresAuth: false }
    },
]

export const router = createRouter({
    history: createWebHistory(),
    routes,
})

router.beforeEach((to, _from, next) => {
    const requiresAuth = to.matched.some(record => record.meta.requiresAuth);
    const isAuthenticated = authService.isAuthenticated();

    if (requiresAuth && !isAuthenticated) {
        next('/user-login');
    } else if ((to.path === '/user-login' || to.path === '/register' || to.path === '/admin-login') && isAuthenticated) {
        next('/');
    } else {
        next();
    }
})
