<template>
  <div class="flex flex-col items-center justify-center min-h-screen bg-gray-100">
    <div class="bg-white shadow-md rounded-xl p-8 w-full max-w-sm">
      <slot name="admin-login">
      </slot>
      <slot name="login">
      </slot>
      <form  @submit.prevent="handleLogin" class="space-y-4 flex flex-col">
        <router-link>Забыли пароль?</router-link>
        <ui-button>Войти</ui-button>
      </form>
    </div>
  </div>
</template>

<script setup lang="ts">
import UiButton from "../ui/Ui-button.vue";
</script>

<style lang="scss" scoped>

</style>
