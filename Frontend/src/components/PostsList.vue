<template>
  <div class="posts-list">
    <div v-if="loading" class="text-center py-8 text-gray-500">
      Загрузка постов...
    </div>
    
    <div v-else-if="error" class="bg-red-50 border border-red-200 rounded-lg p-4 mb-6">
      <p class="text-red-700">{{ error }}</p>
    </div>
    
    <div v-else-if="filteredPosts.length > 0" class="space-y-4 sm:space-y-6">
      <article 
        v-for="post in filteredPosts" 
        :key="post.id" 
        class="bg-white rounded-lg shadow-md overflow-hidden hover:shadow-lg transition-shadow mx-auto"
        :class="getPostAlignmentClass(post)"
      >
        <PostContent :post="post" />
      </article>
    </div>
  </div>
</template>

<script setup lang="ts">
import { watch } from 'vue';
import { usePosts } from '../composables/usePosts';
import { type Post } from '../api/client';
import PostContent from './PostContent.vue';

const emit = defineEmits<{
  (e: 'has-content-change', hasContent: boolean): void;
}>();

const { filteredPosts, loading, error } = usePosts();

const getPostAlignmentClass = (post: Post) => {
  const alignment = (post as any).alignment || 'full-width';
  if (alignment === 'full-width') {
    return 'w-full';
  }
  const alignmentMap = {
    left: 'max-w-3xl mx-auto',
    center: 'max-w-3xl mx-auto',
    right: 'max-w-3xl ml-auto mr-0',
  };
  return alignmentMap[alignment as keyof typeof alignmentMap] || 'w-full';
};

watch([filteredPosts, loading], () => {
  if (loading.value) {
    emit('has-content-change', false);
    return;
  }
  emit('has-content-change', filteredPosts.value.length > 0);
}, { immediate: true });
</script>

<style scoped>
.posts-list {
  @apply w-full;
}
</style>

