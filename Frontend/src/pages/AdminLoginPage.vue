<template>
  <div class="flex flex-col items-center justify-center min-h-screen bg-gray-100 px-4 py-12">
    <div class="bg-white shadow-lg rounded-xl p-6 sm:p-8 w-full max-w-md space-y-6">
      <div class="text-center space-y-2">
        <h1 class="text-2xl sm:text-3xl font-bold text-gray-900">Вход в админ-панель</h1>
        <p class="text-gray-500 text-xs sm:text-sm">Используйте учетные данные администратора</p>
      </div>

      <form class="space-y-4" @submit.prevent="handleSubmit">
        <ui-input
          type="email"
          variant="primary"
          v-model="email"
          placeholder="Email"
          :disabled="loading"
        />
        <ui-input
          type="password"
          variant="primary"
          v-model="password"
          placeholder="Пароль"
          :disabled="loading"
        />

        <div v-if="message" :class="messageClass">
          {{ message }}
        </div>

        <div class="flex flex-wrap items-center justify-between text-sm">
          <router-link to="/admin/forgot-password" class="text-indigo-600 hover:underline">
            Забыли пароль?
          </router-link>
          <router-link to="/" class="text-gray-500 hover:text-gray-700">
            На главную
          </router-link>
        </div>

        <ui-button
          type="submit"
          variant="primary"
          class="w-full"
          :disabled="loading || !email || !password"
        >
          {{ loading ? 'Вход...' : 'Войти' }}
        </ui-button>
      </form>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, ref } from "vue";
import { useRouter } from "vue-router";
import UiButton from "../components/ui/Ui-button.vue";
import UiInput from "../components/ui/Ui-input.vue";
import { adminApi } from "../api/admin";
import { adminAuthService } from "../utils/adminAuth";

const router = useRouter();
const email = ref("");
const password = ref("");
const loading = ref(false);
const message = ref("");
const isError = ref(false);

const messageClass = computed(() =>
  isError.value
    ? "text-xs text-red-600 bg-red-50 border border-red-200 rounded-md p-2"
    : "text-xs text-green-600 bg-green-50 border border-green-200 rounded-md p-2"
);

const handleSubmit = async () => {
  message.value = "";
  isError.value = false;

  if (!email.value || !password.value) {
    message.value = "Пожалуйста, заполните все поля";
    isError.value = true;
    return;
  }

  loading.value = true;
  try {
    const response = await adminApi.login({
      email: email.value,
      password: password.value,
    });
    adminAuthService.setToken(response.token);
    router.push("/admin");
  } catch (error) {
    isError.value = true;
    message.value =
      error instanceof Error ? error.message : "Не удалось выполнить вход";
  } finally {
    loading.value = false;
  }
};
</script>
