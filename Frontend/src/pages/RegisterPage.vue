<template>
  <div class="flex flex-col items-center justify-center min-h-screen bg-gray-100 px-4 py-8">
    <div class="bg-white shadow-md rounded-xl p-5 sm:p-6 md:p-8 w-full max-w-sm">
      <h1 class="text-xl sm:text-2xl font-bold mb-4 sm:mb-6 text-center">Регистрация</h1>
      <form @submit="onForm" class="space-y-4 flex flex-col">
        <div>
          <ui-input 
            type="email" 
            variant="primary" 
            v-model="email"
            placeholder="Email"
            :disabled="loading"
          />
        </div>
        <div>
          <ui-input 
            type="password" 
            variant="primary" 
            v-model="password"
            placeholder="Пароль (минимум 8 символов)"
            :disabled="loading"
          />
        </div>
        <div>
          <ui-input 
            type="password" 
            variant="primary" 
            v-model="confirmPassword"
            placeholder="Подтвердите пароль"
            :disabled="loading"
          />
        </div>
        <div v-if="error" class="text-red-600 text-xs sm:text-sm bg-red-50 border border-red-200 rounded-md p-2 sm:p-3 break-words">
          {{ error }}
        </div>
        <div v-if="success" class="text-green-600 text-xs sm:text-sm bg-green-50 border border-green-200 rounded-md p-2 sm:p-3 break-words">
          {{ success }}
        </div>
        <ui-button 
          type="submit" 
          variant="primary"
          :disabled="loading || !email || !password || !confirmPassword"
        >
          {{ loading ? 'Регистрация...' : 'Зарегистрироваться' }}
        </ui-button>
        <div class="text-center text-xs sm:text-sm text-gray-600 mt-4">
          Уже есть аккаунт? 
          <router-link to="/user-login" class="text-blue-600 hover:underline">
            Войти
          </router-link>
        </div>
      </form>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from "vue";
import { useRouter } from "vue-router";
import UiButton from "../components/ui/Ui-button.vue";
import UiInput from "../components/ui/Ui-input.vue";
import { apiClient } from "../api/client";
import { authService } from "../utils/auth";

const router = useRouter();
const email = ref('');
const password = ref('');
const confirmPassword = ref('');
const error = ref('');
const success = ref('');
const loading = ref(false);

const onForm = async (e: Event) => {
  e.preventDefault();
  error.value = '';
  success.value = '';
  
  if (!email.value || !password.value || !confirmPassword.value) {
    error.value = 'Пожалуйста, заполните все поля';
    return;
  }

  if (password.value.length < 8) {
    error.value = 'Пароль должен содержать минимум 8 символов';
    return;
  }

  if (password.value !== confirmPassword.value) {
    error.value = 'Пароли не совпадают';
    return;
  }

  loading.value = true;
  
  try {
    const response = await apiClient.register({
      email: email.value,
      password: password.value,
    });
    
    authService.setToken(response.token);
    success.value = 'Регистрация успешна! Перенаправление...';
    
    setTimeout(() => {
      router.push('/');
    }, 1000);
  } catch (err) {
    error.value = err instanceof Error ? err.message : 'Ошибка при регистрации. Возможно, такой email уже зарегистрирован.';
  } finally {
    loading.value = false;
  }
}
</script>

<style lang="scss" scoped>

</style>
