<template>
  <div class="page-blocks">
    <div v-if="loading" class="text-center py-4 text-gray-500">
      Загрузка контента...
    </div>
    <div v-else-if="error" class="bg-red-50 border border-red-200 rounded-lg p-4 mb-6">
      <p class="text-red-700">{{ error }}</p>
    </div>
    <div v-else-if="blocks.length === 0" class="text-center py-12 px-4">
      <div class="max-w-md mx-auto">
        <svg class="mx-auto h-16 w-16 text-gray-400 mb-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
        </svg>
        <p class="text-gray-500 text-lg">Контента на странице нет</p>
      </div>
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
    const currentPath = route.path;
    let pageParam: string;
    if (currentPath === '/') {
      pageParam = 'home';
    } else {
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

