<template>
  <div>
    <NavBar v-if="showNavBar" />
    <router-view/>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue';
import { useRoute } from 'vue-router';
import NavBar from './components/NavBar.vue';

const route = useRoute();

const showNavBar = computed(() => {
  const hideOnPages = ['/user-login', '/register', '/admin-login', '/admin/forgot-password', '/admin/reset'];
  return !hideOnPages.includes(route.path);
});
</script>
