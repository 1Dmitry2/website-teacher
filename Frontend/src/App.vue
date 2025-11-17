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
import { authService } from './utils/auth';

const route = useRoute();

const showNavBar = computed(() => {
  const hideOnPages = ['/user-login', '/register', '/admin-login'];
  if (hideOnPages.includes(route.path)) {
    return false;
  }
  return authService.isAuthenticated();
});
</script>
