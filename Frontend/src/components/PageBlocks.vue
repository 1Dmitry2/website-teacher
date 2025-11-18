<template>
  <div class="page-blocks">
    <div v-if="loading" class="text-center py-4 text-gray-500">
      Загрузка контента...
    </div>
    <div v-else-if="error" class="bg-red-50 border border-red-200 rounded-lg p-4 mb-6">
      <p class="text-red-700">{{ error }}</p>
    </div>
    <div v-else>
      <BlockRenderer 
        v-for="block in blocks" 
        :key="block.id"
        :block="block"
      />
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, watch } from 'vue';
import { useRoute } from 'vue-router';
import { apiClient } from '../api/client';
import BlockRenderer, { type Block } from './blocks/BlockRenderer.vue';

const route = useRoute();
const blocks = ref<Block[]>([]);
const loading = ref(true);
const error = ref('');

const fetchBlocks = async () => {
  loading.value = true;
  error.value = '';
  try {
    // Получаем текущий путь страницы
    const currentPath = route.path;
    // Бэкенд ожидает параметр page, где "/" передается как "home"
    // Для путей с несколькими сегментами (например /portfolio/about) нужно передать как "portfolio/about"
    // Но бэкенд конвертирует "home" в "/", а остальные пути должны быть с начальным слешем в базе
    let pageParam: string;
    if (currentPath === '/') {
      pageParam = 'home';
    } else {
      // Убираем первый слеш для параметра URL
      pageParam = currentPath.slice(1);
    }
    
    const fetchedBlocks = await apiClient.getPageBlocks(pageParam);
    blocks.value = fetchedBlocks || [];
  } catch (err) {
    error.value = err instanceof Error ? err.message : 'Ошибка загрузки блоков';
    console.error('Error fetching blocks:', err);
  } finally {
    loading.value = false;
  }
};

onMounted(() => {
  fetchBlocks();
});

watch(() => route.path, () => {
  fetchBlocks();
});
</script>

<style scoped>
.page-blocks {
  @apply w-full;
}
</style>

