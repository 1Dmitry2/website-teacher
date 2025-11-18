<template>
  <div>
    <NavBar v-if="showNavBar" />
    <router-view/>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted } from 'vue';
import { useRoute } from 'vue-router';
import NavBar from './components/NavBar.vue';
import { useAccessibility } from './composables/useAccessibility';

const route = useRoute();

const showNavBar = computed(() => {
  const hideOnPages = ['/user-login', '/register', '/admin-login', '/admin/forgot-password', '/admin/reset'];
  return !hideOnPages.includes(route.path);
});

onMounted(() => {
  useAccessibility();
});
</script>
