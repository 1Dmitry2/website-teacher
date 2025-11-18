<template>
  <div class="max-w-7xl mx-auto px-4 py-10 space-y-8">
    <div class="flex flex-col gap-4 md:flex-row md:items-center md:justify-between">
      <div>
        <p class="text-sm uppercase tracking-wide text-indigo-600 font-semibold">
          Админ-панель
        </p>
        <h1 class="text-3xl font-bold text-gray-900">Панель управления</h1>
        <p class="text-gray-500" v-if="profile">
          {{ profile.email }}
        </p>
      </div>
    </div>

    <div v-if="loading" class="bg-white rounded-xl shadow p-6 text-center text-gray-500">
      Загрузка данных...
    </div>

    <div v-else-if="error" class="bg-red-50 border border-red-200 rounded-xl p-6 space-y-3">
      <p class="text-red-700 text-sm">{{ error }}</p>
      <ui-button type="button" variant="primary" @click="fetchData">Повторить</ui-button>
    </div>

    <div v-else class="space-y-8">
      <div class="grid gap-6 md:grid-cols-2 lg:grid-cols-3">
        <router-link
          to="/admin/blocks"
          class="bg-white rounded-xl shadow p-6 hover:shadow-lg transition-shadow cursor-pointer group"
        >
          <div class="flex items-center gap-4">
            <div class="flex-shrink-0 w-12 h-12 bg-indigo-100 rounded-lg flex items-center justify-center group-hover:bg-indigo-200 transition-colors">
              <svg class="w-6 h-6 text-indigo-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 6h16M4 12h16M4 18h16" />
              </svg>
            </div>
            <div>
              <h3 class="text-lg font-semibold text-gray-900">Блоки страниц</h3>
              <p class="text-sm text-gray-500">Управление контентными блоками</p>
            </div>
          </div>
        </router-link>

        <router-link
          to="/admin/posts"
          class="bg-white rounded-xl shadow p-6 hover:shadow-lg transition-shadow cursor-pointer group"
        >
          <div class="flex items-center gap-4">
            <div class="flex-shrink-0 w-12 h-12 bg-blue-100 rounded-lg flex items-center justify-center group-hover:bg-blue-200 transition-colors">
              <svg class="w-6 h-6 text-blue-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 20H5a2 2 0 01-2-2V6a2 2 0 012-2h10a2 2 0 012 2v1m2 13a2 2 0 01-2-2V7m2 13a2 2 0 002-2V9a2 2 0 00-2-2h-2m-4-3H9M7 16h6M7 8h6v4H7V8z" />
              </svg>
            </div>
            <div>
              <h3 class="text-lg font-semibold text-gray-900">Новости / Посты</h3>
              <p class="text-sm text-gray-500">Создание и редактирование постов</p>
            </div>
          </div>
        </router-link>

        <router-link
          to="/admin/comments"
          class="bg-white rounded-xl shadow p-6 hover:shadow-lg transition-shadow cursor-pointer group"
        >
          <div class="flex items-center gap-4">
            <div class="flex-shrink-0 w-12 h-12 bg-yellow-100 rounded-lg flex items-center justify-center group-hover:bg-yellow-200 transition-colors">
              <svg class="w-6 h-6 text-yellow-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 12h.01M12 12h.01M16 12h.01M21 12c0 4.418-4.03 8-9 8a9.863 9.863 0 01-4.255-.949L3 20l1.395-3.72C3.512 15.042 3 13.574 3 12c0-4.418 4.03-8 9-8s9 3.582 9 8z" />
              </svg>
            </div>
            <div>
              <h3 class="text-lg font-semibold text-gray-900">Комментарии</h3>
              <p class="text-sm text-gray-500">Модерация комментариев</p>
            </div>
          </div>
        </router-link>

        <router-link
          to="/admin/users"
          class="bg-white rounded-xl shadow p-6 hover:shadow-lg transition-shadow cursor-pointer group"
        >
          <div class="flex items-center gap-4">
            <div class="flex-shrink-0 w-12 h-12 bg-gray-100 rounded-lg flex items-center justify-center group-hover:bg-gray-200 transition-colors">
              <svg class="w-6 h-6 text-gray-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4.354a4 4 0 110 5.292M15 21H3v-1a6 6 0 0112 0v1zm0 0h6v-1a6 6 0 00-9-5.197M13 7a4 4 0 11-8 0 4 4 0 018 0z" />
              </svg>
            </div>
            <div>
              <h3 class="text-lg font-semibold text-gray-900">Пользователи</h3>
              <p class="text-sm text-gray-500">Просмотр пользователей</p>
            </div>
          </div>
        </router-link>
      </div>

      <div class="grid gap-6 md:grid-cols-2">
        <div class="bg-white rounded-xl shadow p-6 space-y-4">
          <h2 class="text-xl font-semibold text-gray-900">Информация об администраторе</h2>
          <ul class="text-sm text-gray-600 space-y-2">
            <li>
              <span class="font-medium text-gray-800">Email:</span>
              <span class="ml-2">{{ profile?.email }}</span>
            </li>
            <li v-if="profile?.created_at">
              <span class="font-medium text-gray-800">Создан:</span>
              <span class="ml-2">{{ formatDate(profile?.created_at) }}</span>
            </li>
          </ul>
        </div>

        <div class="bg-white rounded-xl shadow p-6 space-y-3">
          <h2 class="text-xl font-semibold text-gray-900">Сводка</h2>
          <p class="text-gray-700">
            {{ dashboard?.message }}
          </p>
          <p class="text-xs text-gray-400" v-if="dashboard?.generated_at">
            Обновлено: {{ formatDate(dashboard.generated_at) }}
          </p>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { onMounted, ref } from "vue";
import { useRouter } from "vue-router";
import { adminApi, type AdminProfile, type DashboardSummary } from "../api/admin";
import { adminAuthService } from "../utils/adminAuth";

const router = useRouter();
const profile = ref<AdminProfile | null>(null);
const dashboard = ref<DashboardSummary | null>(null);
const loading = ref(true);
const error = ref("");

const formatDate = (value?: string) => {
  if (!value) return "";
  return new Date(value).toLocaleString();
};

const fetchData = async () => {
  loading.value = true;
  error.value = "";
  try {
    const [adminProfile, adminDashboard] = await Promise.all([
      adminApi.me(),
      adminApi.dashboard(),
    ]);
    profile.value = adminProfile;
    dashboard.value = adminDashboard;
  } catch (err) {
    const message =
      err instanceof Error
        ? err.message
        : "Не удалось загрузить данные админки.";
    if (message.toLowerCase().includes("token") || message.toLowerCase().includes("unauthorized")) {
      adminAuthService.removeToken();
      router.push("/admin-login");
      return;
    }
    error.value = message;
  } finally {
    loading.value = false;
  }
};

onMounted(fetchData);
</script>
