<template>
  <div class="flex flex-col items-center justify-center min-h-screen bg-gray-100">
    <div class="bg-white shadow-md rounded-xl p-8 w-full max-w-sm">
      <slot name="admin-login">
      </slot>
      <slot name="login">
      </slot>
      <form @submit="onForm" class="space-y-4 flex flex-col">
        <ui-input type="email" variant="primary" :model-value="email"/>
        <ui-input type="password" variant="primary" :model-value="password"/>
        <router-link to="/">
          <ui-button variant="link" type="button">Забыли пароль?</ui-button>
        </router-link>
        <ui-button type="submit" variant="primary">Войти</ui-button>
      </form>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from "vue";
import UiButton from "../ui/Ui-button.vue";
import UiInput from "../ui/Ui-input.vue";

const email = ref('');
const password = ref('');
const onForm = (e: Event) => {
  e.preventDefault()
}
</script>

<style lang="scss" scoped>

</style>
