<template>
  <div class="flex flex-col items-center justify-center min-h-screen bg-gray-100 px-4 py-12">
    <div class="bg-white shadow-lg rounded-xl p-8 w-full max-w-md space-y-5">
      <div class="text-center space-y-2">
        <h1 class="text-2xl font-bold text-gray-900">Сброс пароля</h1>
        <p class="text-sm text-gray-500">
          Придумайте новый пароль для админского аккаунта
        </p>
      </div>

      <form class="space-y-4" @submit.prevent="handleSubmit">
        <ui-input
          type="password"
          variant="primary"
          v-model="newPassword"
          placeholder="Новый пароль"
          :disabled="loading"
        />
        <ui-input
          type="password"
          variant="primary"
          v-model="confirmPassword"
          placeholder="Повторите пароль"
          :disabled="loading"
        />

        <div v-if="message" :class="messageClass">
          {{ message }}
        </div>

        <ui-button
          type="submit"
          variant="primary"
          class="w-full"
          :disabled="loading || !newPassword || !confirmPassword"
        >
          {{ loading ? 'Сохраняем...' : 'Обновить пароль' }}
        </ui-button>

        <router-link to="/admin-login" class="block text-center text-sm text-gray-500 hover:text-gray-700">
          Вернуться к входу
        </router-link>
      </form>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, ref, watch } from "vue";
import { useRoute } from "vue-router";
import UiButton from "../components/ui/Ui-button.vue";
import UiInput from "../components/ui/Ui-input.vue";
import { adminApi } from "../api/admin";

const route = useRoute();

const token = ref<string>(typeof route.query.token === "string" ? route.query.token : "");
const newPassword = ref("");
const confirmPassword = ref("");
const loading = ref(false);
const message = ref("");
const isError = ref(false);

watch(
  () => route.query.token,
  (value) => {
    token.value = typeof value === "string" ? value : "";
  }
);

const messageClass = computed(() =>
  isError.value
    ? "text-xs text-red-600 bg-red-50 border border-red-200 rounded-md p-2"
    : "text-xs text-green-600 bg-green-50 border border-green-200 rounded-md p-2"
);

const handleSubmit = async () => {
  message.value = "";
  isError.value = false;

  if (!token.value) {
    message.value = "Неверная или устаревшая ссылка для сброса.";
    isError.value = true;
    return;
  }

  if (newPassword.value.length < 8) {
    message.value = "Пароль должен быть не короче 8 символов.";
    isError.value = true;
    return;
  }

  if (newPassword.value !== confirmPassword.value) {
    message.value = "Пароли не совпадают.";
    isError.value = true;
    return;
  }

  loading.value = true;
  try {
    await adminApi.resetPassword({
      token: token.value,
      newPassword: newPassword.value,
    });
    message.value = "Пароль успешно обновлён. Теперь можно войти.";
    newPassword.value = "";
    confirmPassword.value = "";
  } catch (error) {
    isError.value = true;
    message.value =
      error instanceof Error ? error.message : "Не удалось обновить пароль.";
  } finally {
    loading.value = false;
  }
};
</script>

