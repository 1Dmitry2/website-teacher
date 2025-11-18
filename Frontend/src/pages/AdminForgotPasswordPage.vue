<template>
  <div class="flex flex-col items-center justify-center min-h-screen bg-gray-100 px-4 py-12">
    <div class="bg-white shadow-lg rounded-xl p-8 w-full max-w-md space-y-5">
      <div class="text-center space-y-2">
        <h1 class="text-2xl font-bold text-gray-900">Восстановление пароля</h1>
        <p class="text-sm text-gray-500">
          Укажите email, на который отправим ссылку для сброса пароля
        </p>
      </div>

      <form class="space-y-4" @submit.prevent="handleSubmit">
        <ui-input
          type="email"
          variant="primary"
          v-model="email"
          placeholder="Админский email"
          :disabled="loading"
        />

        <div v-if="message" :class="messageClass">
          {{ message }}
        </div>

        <ui-button
          type="submit"
          variant="primary"
          class="w-full"
          :disabled="loading || !email"
        >
          {{ loading ? 'Отправляем...' : 'Отправить письмо' }}
        </ui-button>

        <router-link to="/admin-login" class="block text-center text-sm text-gray-500 hover:text-gray-700">
          Вернуться к входу
        </router-link>
      </form>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, ref } from "vue";
import UiButton from "../components/ui/Ui-button.vue";
import UiInput from "../components/ui/Ui-input.vue";
import { adminApi } from "../api/admin";

const email = ref("");
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

  if (!email.value) {
    message.value = "Введите email администратора";
    isError.value = true;
    return;
  }

  loading.value = true;
  try {
    await adminApi.forgotPassword(email.value);
    message.value = "Если email найден, ссылка для сброса пароля отправлена.";
  } catch (error) {
    isError.value = true;
    message.value =
      error instanceof Error ? error.message : "Не удалось отправить письмо.";
  } finally {
    loading.value = false;
  }
};
</script>

