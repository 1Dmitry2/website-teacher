<template>
  <div class="container mx-auto px-4 sm:px-6 lg:px-8 py-6 sm:py-8">
    <h1 class="text-2xl sm:text-3xl lg:text-4xl font-bold text-gray-800 mb-4 sm:mb-6">Новости</h1>
    
    <PageBlocks />
    
    <div v-if="loading" class="text-center py-8 text-gray-500">
      Загрузка новостей...
    </div>
    
    <div v-else-if="error" class="bg-red-50 border border-red-200 rounded-lg p-4 mb-6">
      <p class="text-red-700">{{ error }}</p>
    </div>
    
    <div v-else-if="shouldShowNoPostsMessage" class="text-gray-600 text-center py-8">
      <p>Пока нет новостей на этой странице.</p>
    </div>
    
    <div v-else class="space-y-6">
      <article 
        v-for="post in filteredPosts" 
        :key="post.id" 
        class="bg-white rounded-lg shadow-md overflow-hidden hover:shadow-lg transition-shadow"
      >
        <div v-if="post.images && post.images.length > 0" class="w-full h-64 overflow-hidden">
          <img 
            :src="post.images[0]" 
            :alt="post.title"
            class="w-full h-full object-cover"
          />
        </div>
        <div class="p-6">
          <h2 class="text-2xl font-bold text-gray-900 mb-2">{{ post.title }}</h2>
          <p class="text-gray-600 mb-4 whitespace-pre-line">{{ post.content }}</p>
          <div class="flex items-center justify-between text-sm text-gray-500">
            <span>{{ formatDate(post.created_at) }}</span>
            <div v-if="post.images && post.images.length > 1" class="flex gap-2">
              <img 
                v-for="(img, idx) in post.images.slice(1, 4)" 
                :key="idx"
                :src="img" 
                :alt="`${post.title} - изображение ${idx + 2}`"
                class="w-16 h-16 object-cover rounded"
              />
            </div>
          </div>
          
          <!-- Комментарии -->
          <PostComments :post-id="post.id" />
        </div>
      </article>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, computed } from 'vue';
import { useRoute } from 'vue-router';
import { apiClient, type Post } from '../api/client';
import PageBlocks from '../components/PageBlocks.vue';
import PostComments from '../components/PostComments.vue';

const route = useRoute();
const posts = ref<Post[]>([]);
const blocks = ref<any[]>([]);
const loading = ref(true);
const error = ref('');

const filteredPosts = computed(() => {
  const currentPath = route.path;
  return posts.value.filter(post => {
    if (!post.pages || post.pages.length === 0) {
      return false;
    }
    const pagesArray = Array.isArray(post.pages) ? post.pages : [];
    return pagesArray.includes(currentPath);
  });
});

const hasBlocks = computed(() => {
  return blocks.value && blocks.value.length > 0;
});

const shouldShowNoPostsMessage = computed(() => {
  return filteredPosts.value.length === 0 && !hasBlocks.value;
});

const formatDate = (date: string) => {
  return new Date(date).toLocaleDateString('ru-RU', {
    year: 'numeric',
    month: 'long',
    day: 'numeric'
  });
};

const fetchBlocks = async () => {
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
    console.error('Error fetching blocks:', err);
    blocks.value = [];
  }
};

const fetchPosts = async () => {
  loading.value = true;
  error.value = '';
  try {
    posts.value = await apiClient.getPosts();
  } catch (err) {
    error.value = err instanceof Error ? err.message : 'Ошибка загрузки новостей';
  } finally {
    loading.value = false;
  }
};

onMounted(async () => {
  await Promise.all([fetchPosts(), fetchBlocks()]);
});
</script>

<style lang="scss" scoped>

</style>
