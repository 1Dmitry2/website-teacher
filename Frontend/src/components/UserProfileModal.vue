<template>
  <Transition name="modal">
    <div v-if="isOpen" class="fixed inset-0 bg-gray-900/60 backdrop-blur-sm flex items-center justify-center z-50 p-4" @click.self="close">
      <div class="bg-white rounded-xl shadow-2xl max-w-md w-full transform transition-all">
        <div class="px-6 pt-6 pb-4 border-b border-gray-200 relative">
          <button
            @click="close"
            class="absolute top-6 right-6 text-gray-400 hover:text-gray-600 transition-colors focus:outline-none focus:ring-2 focus:ring-indigo-500 rounded-full p-1"
            type="button"
            aria-label="Закрыть"
          >
            <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
            </svg>
          </button>
          <h2 class="text-2xl font-bold text-gray-900 pr-8">Профиль</h2>
        </div>
        
        <div v-if="loading" class="p-6 text-center text-gray-500">
          Загрузка...
        </div>
        
        <div v-else-if="error" class="p-6">
          <div class="bg-red-50 border border-red-200 rounded-lg p-4">
            <p class="text-red-700">{{ error }}</p>
          </div>
        </div>
        
        <div v-else-if="profile" class="p-6 space-y-4">
          <div class="flex items-center justify-center mb-4">
            <div class="w-20 h-20 rounded-full bg-indigo-100 flex items-center justify-center">
              <span class="text-3xl font-bold text-indigo-600 uppercase">
                {{ getInitial(profile.email) }}
              </span>
            </div>
          </div>
          
          <div class="space-y-3">
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-1">Email</label>
              <p class="text-base text-gray-900">{{ profile.email }}</p>
            </div>
            
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-1">Дата регистрации</label>
              <p class="text-base text-gray-900">{{ formatDate(profile.created_at) }}</p>
            </div>
            
            <div v-if="profile.comments_count !== undefined">
              <label class="block text-sm font-medium text-gray-700 mb-1">Комментариев</label>
              <p class="text-base text-gray-900">{{ profile.comments_count }}</p>
            </div>
            
            <div v-if="profile.banned" class="mt-2">
              <span class="inline-flex items-center px-3 py-1 rounded-full text-sm font-semibold bg-red-100 text-red-800">
                Заблокирован
              </span>
            </div>
            
            <div v-if="profile.is_admin" class="mt-2">
              <span class="inline-flex items-center px-3 py-1 rounded-full text-sm font-semibold bg-indigo-100 text-indigo-800">
                Администратор
              </span>
            </div>
            
            <div v-if="isAdmin && userId && !profile.is_admin" class="mt-4 pt-4 border-t border-gray-200">
              <button
                @click="handleBan"
                :disabled="banning"
                class="w-full px-4 py-2 text-sm font-medium rounded-lg transition-colors"
                :class="profile.banned 
                  ? 'bg-green-600 text-white hover:bg-green-700 disabled:opacity-50' 
                  : 'bg-red-600 text-white hover:bg-red-700 disabled:opacity-50'"
              >
                {{ banning ? 'Обработка...' : (profile.banned ? 'Разбанить пользователя' : 'Забанить пользователя') }}
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>
  </Transition>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue';
import { apiClient } from '../api/client';
import { adminApi } from '../api/admin';

interface Props {
  isOpen: boolean;
  userId?: number | null;
  email?: string;
  isAdmin?: boolean;
}

const props = defineProps<Props>();
const emit = defineEmits<{ (e: 'close'): void }>();

const profile = ref<any>(null);
const loading = ref(false);
const error = ref('');
const banning = ref(false);

const getInitial = (email?: string) => {
  if (!email) return '?';
  return email.charAt(0).toUpperCase();
};

const formatDate = (date?: string) => {
  if (!date) return 'Не указано';
  return new Date(date).toLocaleDateString('ru-RU', {
    year: 'numeric',
    month: 'long',
    day: 'numeric',
    hour: '2-digit',
    minute: '2-digit'
  });
};

const fetchProfile = async () => {
  if (!props.isOpen) return;
  
  loading.value = true;
  error.value = '';
  
  try {
    if (props.userId && props.isAdmin) {
      // Получаем информацию о пользователе по ID (для админа - с полной информацией)
      const user = await adminApi.getUserById(props.userId);
      profile.value = user;
    } else if (props.userId) {
      // Обычный пользователь может просматривать информацию о других пользователях (без бана)
      const user = await apiClient.getUserById(props.userId);
      profile.value = user;
    } else if (props.isAdmin) {
      // Текущий админ
      const admin = await adminApi.me();
      profile.value = {
        email: admin.email,
        created_at: admin.created_at || new Date().toISOString(),
        is_admin: true,
      };
    } else {
      // Текущий пользователь
      const user = await apiClient.getProfile();
      profile.value = user;
    }
  } catch (err) {
    error.value = err instanceof Error ? err.message : 'Ошибка загрузки профиля';
  } finally {
    loading.value = false;
  }
};

const handleBan = async () => {
  if (!props.userId || !props.isAdmin || !profile.value) return;
  
  const newBannedStatus = !profile.value.banned;
  const confirmMessage = newBannedStatus 
    ? `Забанить пользователя ${profile.value.email}?`
    : `Разбанить пользователя ${profile.value.email}?`;
  
  if (!confirm(confirmMessage)) return;
  
  banning.value = true;
  try {
    await adminApi.banUser(props.userId, newBannedStatus);
    profile.value.banned = newBannedStatus;
  } catch (err) {
    error.value = err instanceof Error ? err.message : 'Ошибка при изменении статуса бана';
  } finally {
    banning.value = false;
  }
};

const close = () => {
  emit('close');
};

watch(() => props.isOpen, (newValue) => {
  if (newValue) {
    fetchProfile();
  } else {
    profile.value = null;
    error.value = '';
  }
});
</script>

<style scoped>
.modal-enter-active {
  transition: opacity 0.3s ease;
}

.modal-leave-active {
  transition: opacity 0.2s ease;
}

.modal-enter-from,
.modal-leave-to {
  opacity: 0;
}

.modal-enter-active > div,
.modal-leave-active > div {
  transition: transform 0.3s ease, opacity 0.3s ease;
}

.modal-enter-from > div,
.modal-leave-to > div {
  transform: scale(0.95) translateY(-10px);
  opacity: 0;
}
</style>

