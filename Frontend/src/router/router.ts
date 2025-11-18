import { createRouter, createWebHistory } from 'vue-router';
import NewsPage from "../pages/NewsPage.vue";
import AboutPage from "../pages/AboutPage.vue";
import PortfolioPage from "../pages/PortfolioPage.vue";
import PortfolioAboutPage from "../pages/PortfolioAboutPage.vue";
import PortfolioQualificationPage from "../pages/PortfolioQualificationPage.vue";
import PortfolioEventsPage from "../pages/PortfolioEventsPage.vue";
import PortfolioAwardsPage from "../pages/PortfolioAwardsPage.vue";
import PortfolioPublicationsPage from "../pages/PortfolioPublicationsPage.vue";
import MethodologyPage from "../pages/MethodologyPage.vue";
import MethodologyProgramsPage from "../pages/MethodologyProgramsPage.vue";
import MethodologyDevelopmentsPage from "../pages/MethodologyDevelopmentsPage.vue";
import MethodologyLawPage from "../pages/MethodologyLawPage.vue";
import MethodologyVideosPage from "../pages/MethodologyVideosPage.vue";
import AchievementsPage from "../pages/AchievementsPage.vue";
import AchievementsAwardsPage from "../pages/AchievementsAwardsPage.vue";
import AchievementsGalleryPage from "../pages/AchievementsGalleryPage.vue";
import ParentsPage from "../pages/ParentsPage.vue";
import ParentsTipsPage from "../pages/ParentsTipsPage.vue";
import ParentsDaySchedulePage from "../pages/ParentsDaySchedulePage.vue";
import ParentsPosturePage from "../pages/ParentsPosturePage.vue";
import ParentsFlatFeetPage from "../pages/ParentsFlatFeetPage.vue";
import ParentsIllnessPreventionPage from "../pages/ParentsIllnessPreventionPage.vue";
import ParentsHardeningPage from "../pages/ParentsHardeningPage.vue";
import AdminPage from "../pages/AdminPage.vue";
import AdminLoginPage from "../pages/AdminLoginPage.vue";
import AdminForgotPasswordPage from "../pages/AdminForgotPasswordPage.vue";
import AdminResetPasswordPage from "../pages/AdminResetPasswordPage.vue";
import AdminBlocksPage from "../pages/admin/AdminBlocksPage.vue";
import AdminPostsPage from "../pages/admin/AdminPostsPage.vue";
import AdminCommentsPage from "../pages/admin/AdminCommentsPage.vue";
import AdminUsersPage from "../pages/admin/AdminUsersPage.vue";
import UserLogin from "../pages/UserLogin.vue";
import RegisterPage from "../pages/RegisterPage.vue";
import { authService } from "../utils/auth";
import { adminAuthService } from "../utils/adminAuth";


const routes = [
    { 
        path: '/', 
        component: NewsPage,
        meta: { requiresAuth: false }
    },
    { 
        path: '/about', 
        component: AboutPage,
        meta: { requiresAuth: false }
    },
    { 
        path: '/portfolio', 
        component: PortfolioPage,
        meta: { requiresAuth: false }
    },
    { 
        path: '/portfolio/about', 
        component: PortfolioAboutPage,
        meta: { requiresAuth: false }
    },
    { 
        path: '/portfolio/qualification', 
        component: PortfolioQualificationPage,
        meta: { requiresAuth: false }
    },
    { 
        path: '/portfolio/events', 
        component: PortfolioEventsPage,
        meta: { requiresAuth: false }
    },
    { 
        path: '/portfolio/awards', 
        component: PortfolioAwardsPage,
        meta: { requiresAuth: false }
    },
    { 
        path: '/portfolio/publications', 
        component: PortfolioPublicationsPage,
        meta: { requiresAuth: false }
    },
    { 
        path: '/methodology', 
        component: MethodologyPage,
        meta: { requiresAuth: false }
    },
    { 
        path: '/methodology/programs', 
        component: MethodologyProgramsPage,
        meta: { requiresAuth: false }
    },
    { 
        path: '/methodology/developments', 
        component: MethodologyDevelopmentsPage,
        meta: { requiresAuth: false }
    },
    { 
        path: '/methodology/law', 
        component: MethodologyLawPage,
        meta: { requiresAuth: false }
    },
    { 
        path: '/methodology/videos', 
        component: MethodologyVideosPage,
        meta: { requiresAuth: false }
    },
    { 
        path: '/achievements', 
        component: AchievementsPage,
        meta: { requiresAuth: false }
    },
    { 
        path: '/achievements/awards', 
        component: AchievementsAwardsPage,
        meta: { requiresAuth: false }
    },
    { 
        path: '/achievements/gallery', 
        component: AchievementsGalleryPage,
        meta: { requiresAuth: false }
    },
    { 
        path: '/parents', 
        component: ParentsPage,
        meta: { requiresAuth: false }
    },
    { 
        path: '/parents/tips', 
        component: ParentsTipsPage,
        meta: { requiresAuth: false }
    },
    { 
        path: '/parents/day-schedule', 
        component: ParentsDaySchedulePage,
        meta: { requiresAuth: false }
    },
    { 
        path: '/parents/posture', 
        component: ParentsPosturePage,
        meta: { requiresAuth: false }
    },
    { 
        path: '/parents/flat-feet', 
        component: ParentsFlatFeetPage,
        meta: { requiresAuth: false }
    },
    { 
        path: '/parents/illness-prevention', 
        component: ParentsIllnessPreventionPage,
        meta: { requiresAuth: false }
    },
    { 
        path: '/parents/hardening', 
        component: ParentsHardeningPage,
        meta: { requiresAuth: false }
    },
    { 
        path: '/admin', 
        component: AdminPage,
        meta: { requiresAdmin: true }
    },
    {
        path: '/admin/blocks',
        component: AdminBlocksPage,
        meta: { requiresAdmin: true }
    },
    {
        path: '/admin/posts',
        component: AdminPostsPage,
        meta: { requiresAdmin: true }
    },
    {
        path: '/admin/comments',
        component: AdminCommentsPage,
        meta: { requiresAdmin: true }
    },
    {
        path: '/admin/users',
        component: AdminUsersPage,
        meta: { requiresAdmin: true }
    },
    { 
        path: '/admin-login', 
        component: AdminLoginPage,
        meta: { adminGuestOnly: true }
    },
    {
        path: '/admin/forgot-password',
        component: AdminForgotPasswordPage,
        meta: { adminGuestOnly: true }
    },
    {
        path: '/admin/reset',
        component: AdminResetPasswordPage,
        meta: { adminGuestOnly: true }
    },
    { 
        path: '/user-login', 
        component: UserLogin,
        meta: { requiresAuth: false, userGuestOnly: true }
    },
    { 
        path: '/register', 
        component: RegisterPage,
        meta: { requiresAuth: false, userGuestOnly: true }
    },
]

export const router = createRouter({
    history: createWebHistory(),
    routes,
})

router.beforeEach((to, _from, next) => {
    const requiresUserAuth = to.matched.some(record => record.meta.requiresAuth);
    const requiresAdminAuth = to.matched.some(record => record.meta.requiresAdmin);
    const adminGuestOnly = to.matched.some(record => record.meta.adminGuestOnly);
    const userGuestOnly = to.matched.some(record => record.meta.userGuestOnly);

    const isUserAuthenticated = authService.isAuthenticated();
    const isAdminAuthenticated = adminAuthService.isAuthenticated();

    if (requiresAdminAuth && !isAdminAuthenticated) {
        next('/admin-login');
        return;
    }

    if (adminGuestOnly && isAdminAuthenticated) {
        next('/admin');
        return;
    }

    if (requiresUserAuth && !isUserAuthenticated) {
        next('/user-login');
        return;
    }

    if (userGuestOnly && isUserAuthenticated) {
        next('/');
        return;
    }

    next();
})
