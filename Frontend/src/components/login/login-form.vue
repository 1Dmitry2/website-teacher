<template>
  <div class="flex flex-col items-center justify-center min-h-screen bg-gray-100 px-4 py-8">
    <div class="bg-white shadow-md rounded-xl p-6 sm:p-8 w-full max-w-sm">
      <slot name="admin-login">
      </slot>
      <slot name="login">
      </slot>
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
            placeholder="Пароль"
            :disabled="loading"
          />
        </div>
        <div v-if="error" class="text-red-600 text-xs sm:text-sm bg-red-50 border border-red-200 rounded-md p-2 sm:p-3">
          {{ error }}
        </div>
        <div class="text-center sm:text-left">
          <router-link to="/">
            <ui-button variant="link" type="button">Забыли пароль?</ui-button>
          </router-link>
        </div>
        <ui-button 
          type="submit" 
          variant="primary"
          :disabled="loading || !email || !password"
        >
          {{ loading ? 'Вход...' : 'Войти' }}
        </ui-button>
        <div class="text-center text-xs sm:text-sm text-gray-600 mt-4">
          Нет аккаунта? 
          <router-link to="/register" class="text-blue-600 hover:underline">
            Зарегистрироваться
          </router-link>
        </div>
      </form>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from "vue";
import { useRouter } from "vue-router";
import UiButton from "../ui/Ui-button.vue";
import UiInput from "../ui/Ui-input.vue";
import { apiClient } from "../../api/client";
import { authService } from "../../utils/auth";

const router = useRouter();
const email = ref('');
const password = ref('');
const error = ref('');
const loading = ref(false);

const onForm = async (e: Event) => {
  e.preventDefault();
  error.value = '';
  
  if (!email.value || !password.value) {
    error.value = 'Пожалуйста, заполните все поля';
    return;
  }

  loading.value = true;
  
  try {
    const response = await apiClient.login({
      email: email.value,
      password: password.value,
    });
    
    authService.setToken(response.token);
    router.push('/');
  } catch (err) {
    error.value = err instanceof Error ? err.message : 'Ошибка при входе. Проверьте email и пароль.';
  } finally {
    loading.value = false;
  }
}
</script>

<style lang="scss" scoped>

</style>
