<template>
  <nav ref="navRef" class="bg-white shadow-md">
    <div class="max-w-7xl mx-auto px-3 sm:px-4 md:px-6 lg:px-8">
      <div class="flex flex-wrap items-center justify-between gap-4 py-3">
        <div class="flex flex-1 items-center gap-4 min-w-0">
          <button
            class="md:hidden inline-flex items-center justify-center p-2 rounded-md text-gray-700 hover:bg-gray-100 focus:outline-none focus:ring-2 focus:ring-indigo-500"
            type="button"
            @click="toggleMobileMenu"
          >
            <span class="sr-only">Меню</span>
            <svg class="h-6 w-6" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
              <path stroke-linecap="round" stroke-linejoin="round" d="M4 6h16M4 12h16M4 18h16" />
            </svg>
          </button>
          <ul class="hidden md:flex items-center gap-4 flex-wrap">
            <li v-for="item in navigationItems" :key="item.title" class="relative">
              <router-link
                v-if="item.to && !item.children"
                :to="item.to"
                :class="[
                  'inline-flex items-center gap-2 text-sm font-medium transition-colors cursor-pointer px-2 py-1 rounded-md',
                  isActiveItem(item) ? 'bg-indigo-100 text-indigo-700' : 'text-gray-700 hover:text-indigo-600'
                ]"
                @click="handlePrimaryClick(item)"
              >
                <span
                  :class="[
                    'flex h-8 w-8 items-center justify-center rounded-full',
                    isActiveItem(item) ? 'bg-indigo-600 text-white' : 'bg-indigo-50 text-indigo-600'
                  ]"
                >
                  <svg class="h-4 w-4" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
                    <path
                      v-for="(path, index) in item.iconPaths"
                      :key="`${item.title}-icon-${index}`"
                      :d="path"
                      stroke-linecap="round"
                      stroke-linejoin="round"
                    />
                  </svg>
                </span>
                <span>{{ item.title }}</span>
              </router-link>
              <button
                v-else
                class="inline-flex items-center gap-2 text-sm font-medium text-gray-700 hover:text-indigo-600 transition-colors cursor-pointer"
                type="button"
                @click.stop="toggleDropdown(item.title)"
              >
                <span class="flex h-8 w-8 items-center justify-center rounded-full bg-indigo-50 text-indigo-600">
                  <svg class="h-4 w-4" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
                    <path
                      v-for="(path, index) in item.iconPaths"
                      :key="`${item.title}-icon-${index}`"
                      :d="path"
                      stroke-linecap="round"
                      stroke-linejoin="round"
                    />
                  </svg>
                </span>
                <span>{{ item.title }}</span>
                <svg class="h-4 w-4 text-gray-500" viewBox="0 0 20 20" fill="currentColor">
                  <path
                    fill-rule="evenodd"
                    d="M5.23 7.21a.75.75 0 011.06.02L10 11.186l3.71-3.955a.75.75 0 111.08 1.04l-4.24 4.52a.75.75 0 01-1.08 0L5.21 8.27a.75.75 0 01.02-1.06z"
                    clip-rule="evenodd"
                  />
                </svg>
              </button>

              <Transition name="dropdown">
                <div
                  v-if="item.children && activeDropdown === item.title"
                  class="absolute left-0 mt-3 w-64 rounded-md shadow-lg bg-white ring-1 ring-black ring-opacity-5 z-20"
                  @click.stop
                >
                  <ul class="py-3">
                    <li v-for="child in item.children" :key="child.title">
                      <button
                        class="w-full text-left px-4 py-2 text-sm text-gray-700 hover:bg-indigo-50 hover:text-indigo-600 transition-colors cursor-pointer"
                        type="button"
                        @click="handleChildSelect(child)"
                      >
                        {{ child.title }}
                      </button>
                    </li>
                  </ul>
                </div>
              </Transition>
            </li>
          </ul>
        </div>

        <div class="flex flex-wrap items-end gap-1.5 sm:gap-2 justify-end text-right">
          <template v-if="isAuthenticated">
            <div class="flex items-center gap-2">
              <div class="flex flex-col text-right leading-tight">
                <span class="text-gray-700 text-xs sm:text-sm md:text-base truncate max-w-[200px] sm:max-w-[240px] md:max-w-none">
                  {{ accountEmail }}
                  <template v-if="isAdminSession">
                    <span class="text-[11px] sm:text-xs uppercase tracking-wide text-indigo-600 font-semibold">
                      (ADMIN)
                    </span>
                  </template>
                </span>
              </div>
              <ui-button
                variant="primary"
                type="button"
                class="!w-auto px-3 py-1 text-xs sm:text-sm"
                @click="handleLogout"
              >
                Выйти
              </ui-button>
            </div>
          </template>
          <template v-else>
            <ui-button variant="primary" type="button" @click="handleLogin">
              Войти
            </ui-button>
          </template>
        </div>
      </div>

      <div v-if="isMobileMenuOpen" class="md:hidden border-t border-gray-200 pt-4 pb-6">
        <ul class="space-y-2">
          <li v-for="item in navigationItems" :key="`mobile-${item.title}`" class="border-b border-gray-100 pb-2">
            <router-link
              v-if="item.to && !item.children"
              :to="item.to"
                :class="[
                  'flex items-center gap-3 text-base font-medium py-2 cursor-pointer px-2 rounded-md',
                  isActiveItem(item) ? 'bg-indigo-50 text-indigo-700' : 'text-gray-800'
                ]"
              @click="handlePrimaryClick(item)"
            >
              <span
                :class="[
                  'flex h-8 w-8 items-center justify-center rounded-full',
                  isActiveItem(item) ? 'bg-indigo-600 text-white' : 'bg-indigo-50 text-indigo-600'
                ]"
              >
                <svg class="h-4 w-4" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
                  <path
                    v-for="(path, index) in item.iconPaths"
                    :key="`mobile-${item.title}-icon-${index}`"
                    :d="path"
                    stroke-linecap="round"
                    stroke-linejoin="round"
                  />
                </svg>
              </span>
              {{ item.title }}
            </router-link>
            <div v-else>
              <button
                class="flex w-full items-center justify-between text-base font-medium text-gray-800 py-2 cursor-pointer"
                type="button"
                @click="toggleDropdown(item.title)"
              >
                <span class="flex items-center gap-3">
                  <span class="flex h-8 w-8 items-center justify-center rounded-full bg-indigo-50 text-indigo-600">
                    <svg class="h-4 w-4" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
                      <path
                        v-for="(path, index) in item.iconPaths"
                        :key="`mobile-toggle-${item.title}-icon-${index}`"
                        :d="path"
                        stroke-linecap="round"
                        stroke-linejoin="round"
                      />
                    </svg>
                  </span>
                  <span>{{ item.title }}</span>
                </span>
                <svg class="h-5 w-5" viewBox="0 0 20 20" fill="currentColor">
                  <path
                    fill-rule="evenodd"
                    d="M5.23 7.21a.75.75 0 011.06.02L10 11.186l3.71-3.955a.75.75 0 111.08 1.04l-4.24 4.52a.75.75 0 01-1.08 0L5.21 8.27a.75.75 0 01.02-1.06z"
                    clip-rule="evenodd"
                  />
                </svg>
              </button>
              <Transition name="dropdown">
                <ul v-if="item.children && activeDropdown === item.title" class="mt-2 space-y-1 pl-4">
                <li v-for="child in item.children" :key="`mobile-child-${child.title}`">
                  <button
                    class="w-full text-left py-1 text-sm text-gray-600 hover:text-indigo-600 transition-colors cursor-pointer"
                    type="button"
                    @click="handleChildSelect(child)"
                  >
                    {{ child.title }}
                  </button>
                </li>
              </ul>
              </Transition>
            </div>
          </li>
        </ul>
      </div>
    </div>
  </nav>
</template>

<script setup lang="ts">
import { ref, onMounted, watch, onBeforeUnmount, computed } from 'vue';
import { useRouter, useRoute } from 'vue-router';
import UiButton from './ui/Ui-button.vue';
import { authService } from '../utils/auth';
import { adminAuthService } from '../utils/adminAuth';
import { apiClient } from '../api/client';
import { adminApi } from '../api/admin';

type NavChild = {
  title: string;
  href?: string;
  to?: string;
};

type NavItem = {
  title: string;
  to?: string;
  children?: NavChild[];
  iconPaths: string[];
};

const router = useRouter();
const route = useRoute();
const authRole = ref<'guest' | 'user' | 'admin'>('guest');
const accountEmail = ref('');
const isMobileMenuOpen = ref(false);
const activeDropdown = ref<string | null>(null);
const navRef = ref<HTMLElement | null>(null);
const isAuthenticated = computed(() => authRole.value !== 'guest');
const isAdminSession = computed(() => authRole.value === 'admin');

const baseNavigation: NavItem[] = [
  {
    title: 'Новости',
    to: '/',
    iconPaths: ['M5 6H19V18H6a2 2 0 01-2-2V6z', 'M7 10h10', 'M7 13h8', 'M7 16h5'],
  },
  {
    title: 'Портфолио',
    to: '/portfolio',
    iconPaths: ['M3 7h18', 'M5 5h4l2 2h10v11H5a2 2 0 01-2-2V5z'],
    children: [
      { title: 'О себе', to: '/portfolio/about' },
      { title: 'Повышение квалификации', to: '/portfolio/qualification' },
      { title: 'Участие в мероприятиях', to: '/portfolio/events' },
      { title: 'Дипломы, сертификаты и благодарности', to: '/portfolio/awards' },
      { title: 'Публикации', to: '/portfolio/publications' },
    ],
  },
  {
    title: 'Методическая копилка',
    to: '/methodology',
    iconPaths: ['M6 5h7a2 2 0 012 2v12l-5-3-5 3V7a2 2 0 012-2z', 'M15 8h5v11l-5-3'],
    children: [
      { title: 'Программы', to: '/methodology/programs' },
      { title: 'Методические разработки', to: '/methodology/developments' },
      { title: 'Нормативно-правовая база', to: '/methodology/law' },
      { title: 'Видеозаписи', to: '/methodology/videos' },
    ],
  },
  {
    title: 'Достижения обучающихся',
    to: '/achievements',
    iconPaths: [
      'M8 5h8v4a4 4 0 01-4 4 4 4 0 01-4-4V5z',
      'M5 5h3v4a4 4 0 01-4 4V7a2 2 0 012-2z',
      'M16 5h3a2 2 0 012 2v6a4 4 0 01-4-4z',
      'M9 19h6',
      'M12 15v4',
    ],
    children: [
      { title: 'Дипломы и сертификаты', to: '/achievements/awards' },
      { title: 'Фотогалерея', to: '/achievements/gallery' },
    ],
  },
  {
    title: 'Родителям',
    to: '/parents',
    iconPaths: ['M12 12a4 4 0 100-8 4 4 0 000 8z', 'M4 20v-1a5 5 0 015-5h6a5 5 0 015 5v1'],
    children: [
      { title: 'Советы', to: '/parents/tips' },
      { title: 'Режим дня школьника', to: '/parents/day-schedule' },
      { title: 'Нарушение осанки', to: '/parents/posture' },
      { title: 'Плоскостопие', to: '/parents/flat-feet' },
      { title: 'Профилактика заболеваний ОРЗ и ОРВИ', to: '/parents/illness-prevention' },
      { title: 'Закаливание в семье', to: '/parents/hardening' },
    ],
  },
];

const adminNavItem: NavItem = {
  title: 'Админ-панель',
  to: '/admin',
  iconPaths: ['M6 6h12v12H6z', 'M9 3v4', 'M15 3v4', 'M6 11h12', 'M6 15h12'],
};

const navigationItems = computed(() =>
  isAdminSession.value ? [...baseNavigation, adminNavItem] : baseNavigation
);

const resetAuthState = () => {
  authRole.value = 'guest';
  accountEmail.value = '';
};

const checkAuth = async () => {
  const adminLoggedIn = adminAuthService.isAuthenticated();
  if (adminLoggedIn) {
    try {
      const admin = await adminApi.me();
      authRole.value = 'admin';
      accountEmail.value = admin.email;
      return;
    } catch (_error) {
      adminAuthService.removeToken();
    }
  }

  const userLoggedIn = authService.isAuthenticated();
  if (userLoggedIn) {
    try {
      const user = await apiClient.getProfile();
      authRole.value = 'user';
      accountEmail.value = user.email;
      return;
    } catch (_error) {
      authService.removeToken();
    }
  }

  resetAuthState();
};

const handleLogout = () => {
  closeMobileMenu();
  if (authRole.value === 'admin') {
    adminAuthService.removeToken();
    resetAuthState();
    router.push('/admin-login');
    return;
  }
  authService.removeToken();
  resetAuthState();
  router.push('/user-login');
};

const handleLogin = () => {
  closeMobileMenu();
  router.push('/user-login');
};

const toggleMobileMenu = () => {
  isMobileMenuOpen.value = !isMobileMenuOpen.value;
};

const closeMobileMenu = () => {
  isMobileMenuOpen.value = false;
  activeDropdown.value = null;
};

const toggleDropdown = (title: string) => {
  activeDropdown.value = activeDropdown.value === title ? null : title;
};

const handlePrimaryClick = (_item: NavItem) => {
  activeDropdown.value = null;
  closeMobileMenu();
};

const scrollToHash = (hash: string) => {
  if (!hash.startsWith('#')) {
    return;
  }
  const target = document.querySelector(hash);
  if (target) {
    target.scrollIntoView({ behavior: 'smooth', block: 'start' });
  }
};

const handleChildSelect = (child: NavChild) => {
  activeDropdown.value = null;
  closeMobileMenu();
  if (child.to) {
    router.push(child.to);
    return;
  }
  if (child.href) {
    if (child.href.startsWith('#')) {
      scrollToHash(child.href);
      router.replace({ hash: child.href });
    } else {
      window.open(child.href, '_blank', 'noopener');
    }
  }
};

const handleDocumentClick = (event: MouseEvent) => {
  if (navRef.value?.contains(event.target as Node)) {
    return;
  }
  activeDropdown.value = null;
};

const isActiveItem = (item: NavItem) => {
  if (!item.to) {
    return false;
  }
  return router.resolve(item.to).path === route.path;
};

onMounted(() => {
  checkAuth();
  document.addEventListener('click', handleDocumentClick);
});

onBeforeUnmount(() => {
  document.removeEventListener('click', handleDocumentClick);
});

watch(() => route.path, () => {
  checkAuth();
});
</script>

<style lang="scss" scoped>

</style>

