<template>
  <nav class="bg-white shadow-md">
    <div class="max-w-7xl mx-auto px-3 sm:px-4 md:px-6 lg:px-8">
      <div class="flex justify-between items-center h-14 sm:h-16">
        <div class="flex items-center">
          <router-link to="/" class="text-lg sm:text-xl font-bold text-gray-800">
            Новости
          </router-link>
        </div>
        <div class="flex items-center gap-2 sm:gap-4">
          <span class="text-gray-700 text-xs sm:text-sm md:text-base truncate max-w-[120px] sm:max-w-none">
            {{ userEmail }}
          </span>
          <ui-button variant="primary" type="button" @click="handleLogout">
            Выйти
          </ui-button>
        </div>
      </div>
    </div>
  </nav>
</template>

<script setup lang="ts">
import { ref, onMounted, watch } from 'vue';
import { useRouter, useRoute } from 'vue-router';
import UiButton from './ui/Ui-button.vue';
import { authService } from '../utils/auth';
import { apiClient } from '../api/client';

const router = useRouter();
const route = useRoute();
const isAuthenticated = ref(false);
const userEmail = ref('');

const checkAuth = async () => {
  isAuthenticated.value = authService.isAuthenticated();
  if (isAuthenticated.value) {
    try {
      const user = await apiClient.getProfile();
      userEmail.value = user.email;
    } catch (error) {
      authService.removeToken();
      isAuthenticated.value = false;
      userEmail.value = '';
    }
  } else {
    userEmail.value = '';
  }
};

const handleLogout = () => {
  authService.removeToken();
  isAuthenticated.value = false;
  userEmail.value = '';
  router.push('/user-login');
};

onMounted(() => {
  checkAuth();
});
watch(() => route.path, () => {
  checkAuth();
});
</script>

<style lang="scss" scoped>

</style>

