<template>
  <div class="container mx-auto px-4 sm:px-6 lg:px-8 py-6 sm:py-8">
    <h1 class="text-2xl sm:text-3xl lg:text-4xl font-bold text-gray-800 mb-4 sm:mb-6">О себе</h1>
    
    <PageBlocks />
    
    <div v-if="!hasBlocks" class="prose max-w-none">
      <p class="text-gray-600">Здесь будет информация о преподавателе.</p>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue';
import PageBlocks from '../components/PageBlocks.vue';

const hasBlocks = ref(false);
</script>

<style lang="scss" scoped>
</style>

