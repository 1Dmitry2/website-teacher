<template>
  <div class="max-w-7xl mx-auto px-4 py-10 space-y-6">
    <div>
      <h1 class="text-3xl font-bold text-gray-900">Пользователи</h1>
      <p class="text-gray-500 mt-1">Просмотр зарегистрированных пользователей</p>
    </div>

    <div v-if="loading" class="bg-white rounded-xl shadow p-6 text-center text-gray-500">Загрузка...</div>
    <div v-else-if="error" class="bg-red-50 border border-red-200 rounded-xl p-6"><p class="text-red-700">{{ error }}</p></div>
    <div v-else class="bg-white rounded-xl shadow overflow-hidden">
      <div class="overflow-x-auto">
        <table class="min-w-full divide-y divide-gray-200">
          <thead class="bg-gray-50">
            <tr>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">ID</th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Email</th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Статус</th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Дата регистрации</th>
            </tr>
          </thead>
          <tbody class="bg-white divide-y divide-gray-200">
            <tr v-for="user in users" :key="user.id">
              <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900">{{ user.id }}</td>
              <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900">{{ user.email }}</td>
              <td class="px-6 py-4 whitespace-nowrap">
                <span v-if="user.is_admin" class="px-2 py-1 text-xs font-semibold bg-indigo-100 text-indigo-800 rounded-full">Админ</span>
                <span v-else class="px-2 py-1 text-xs font-semibold bg-gray-100 text-gray-800 rounded-full">Пользователь</span>
              </td>
              <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500">{{ formatDate(user.created_at) }}</td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue';
import { useRouter } from 'vue-router';
import { adminApi, type User } from '../../api/admin';
import { adminAuthService } from '../../utils/adminAuth';

const router = useRouter();
const users = ref<User[]>([]);
const loading = ref(true);
const error = ref('');

const formatDate = (date?: string) => {
  if (!date) return '—';
  return new Date(date).toLocaleDateString('ru-RU');
};

const fetchUsers = async () => {
  loading.value = true;
  error.value = '';
  try {
    users.value = await adminApi.getUsers();
  } catch (err) {
    error.value = err instanceof Error ? err.message : 'Ошибка загрузки';
    if (error.value.toLowerCase().includes('unauthorized') || error.value.toLowerCase().includes('token')) {
      adminAuthService.removeToken();
      router.push('/admin-login');
    }
  } finally {
    loading.value = false;
  }
};

onMounted(fetchUsers);
</script>

