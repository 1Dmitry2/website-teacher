<template>
  <div class="post-content">
    <template v-if="hasMedia">
      <div :class="getMediaLayoutClass()">
        <div v-if="titlePosition === 'top'" :class="getTitleClass('top')">
          <h2 class="text-xl sm:text-2xl font-bold text-gray-900 mb-2">{{ post.title }}</h2>
        </div>
        
        <div v-if="contentPosition === 'top'" :class="getContentClass('top')">
          <p class="text-sm sm:text-base text-gray-600 mb-4 whitespace-pre-line">{{ post.content }}</p>
        </div>
        
        <div :class="getMediaContainerClass()">
          <div v-if="post.videos && post.videos.length > 0 && post.videos[0]" class="w-full">
            <video 
              :src="getMediaUrl(post.videos[0])" 
              :alt="post.title"
              :class="`w-full h-auto ${getSizeClass(getMediaSize(post.videos[0]), true)} object-cover rounded-lg mx-auto`"
              controls
            ></video>
          </div>
          <div v-else-if="post.images && post.images.length > 0" class="w-full overflow-hidden flex justify-center">
            <img 
              :src="getMediaUrl(post.images[0])" 
              :alt="post.title"
              :class="`${getSizeClass(getMediaSize(post.images[0]))} object-cover rounded-lg`"
            />
          </div>
        </div>
        
        <div v-if="(titlePosition === 'left' || titlePosition === 'right') || (contentPosition === 'left' || contentPosition === 'right')" 
             :class="getTextContainerClass()">
          <div v-if="titlePosition === 'left' || titlePosition === 'right'">
            <h2 class="text-xl sm:text-2xl font-bold text-gray-900 mb-2">{{ post.title }}</h2>
          </div>
          <div v-if="contentPosition === 'left' || contentPosition === 'right'">
            <p class="text-sm sm:text-base text-gray-600 mb-4 whitespace-pre-line">{{ post.content }}</p>
          </div>
        </div>
        
        <div v-if="titlePosition === 'bottom'" :class="getTitleClass('bottom')">
          <h2 class="text-xl sm:text-2xl font-bold text-gray-900 mb-2">{{ post.title }}</h2>
        </div>
        
        <div v-if="contentPosition === 'bottom'" :class="getContentClass('bottom')">
          <p class="text-sm sm:text-base text-gray-600 mb-4 whitespace-pre-line">{{ post.content }}</p>
        </div>
      </div>
    </template>
    
    <template v-else>
      <div class="p-4 sm:p-6">
        <h2 class="text-xl sm:text-2xl font-bold text-gray-900 mb-2">{{ post.title }}</h2>
        <p class="text-sm sm:text-base text-gray-600 mb-4 whitespace-pre-line">{{ post.content }}</p>
      </div>
    </template>
    
    <div class="p-4 sm:p-6 pt-0">
      <div class="flex flex-col sm:flex-row items-start sm:items-center justify-between gap-2 sm:gap-0 text-xs sm:text-sm text-gray-500">
        <span>{{ formatDate(post.created_at) }}</span>
        <div v-if="post.images && post.images.length > 1" class="flex gap-1.5 sm:gap-2">
          <img 
            v-for="(img, idx) in post.images.slice(1, 4)" 
            :key="idx"
            :src="getMediaUrl(img)" 
            :alt="`${post.title} - изображение ${idx + 2}`"
            class="w-12 h-12 sm:w-16 sm:h-16 object-cover rounded"
          />
        </div>
      </div>
      
      <PostComments :post-id="post.id" />
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue';
import { type Post, API_BASE_URL, type MediaItem, type MediaSize } from '../api/client';
import PostComments from './PostComments.vue';

const props = defineProps<{
  post: Post;
}>();

const hasMedia = computed(() => {
  return (props.post.images && props.post.images.length > 0) || 
         (props.post.videos && props.post.videos.length > 0);
});

const titlePosition = computed(() => {
  return props.post.title_position || 'top';
});

const contentPosition = computed(() => {
  return props.post.content_position || 'bottom';
});

const getMediaUrl = (item: string | MediaItem): string => {
  const url = typeof item === 'string' ? item : item.url;
  if (!url) return '';
  if (url.startsWith('http://') || url.startsWith('https://')) {
    return url;
  }
  if (url.startsWith('/')) {
    return `${API_BASE_URL}${url}`;
  }
  return `${API_BASE_URL}/uploads/${url}`;
};

const getMediaSize = (item: string | MediaItem): MediaSize => {
  return typeof item === 'string' ? 'medium' : (item.size || 'medium');
};

const getSizeClass = (size: MediaSize, isVideo = false): string => {
  const sizeMap = {
    small: isVideo ? 'max-w-xs' : 'max-w-xs h-32 sm:h-40',
    medium: isVideo ? 'max-w-md' : 'max-w-md h-48 sm:h-64',
    large: isVideo ? 'max-w-2xl' : 'max-w-2xl h-64 sm:h-80',
    xlarge: isVideo ? 'max-w-4xl' : 'max-w-4xl h-80 sm:h-96',
  };
  return sizeMap[size] || sizeMap.medium;
};

const getMediaLayoutClass = () => {
  const hasHorizontalText = 
    titlePosition.value === 'left' || titlePosition.value === 'right' ||
    contentPosition.value === 'left' || contentPosition.value === 'right';
  
  if (hasHorizontalText) {
    return 'flex flex-col md:flex-row gap-4 p-4 sm:p-6';
  }
  return 'p-4 sm:p-6 space-y-4';
};

const getMediaContainerClass = () => {
  const hasHorizontalText = 
    titlePosition.value === 'left' || titlePosition.value === 'right' ||
    contentPosition.value === 'left' || contentPosition.value === 'right';
  
  if (hasHorizontalText) {
    const textOnLeft = titlePosition.value === 'left' || contentPosition.value === 'left';
    if (textOnLeft) {
      return 'w-full md:w-1/2 flex-shrink-0 order-2 md:order-2';
    } else {
      return 'w-full md:w-1/2 flex-shrink-0 order-2 md:order-2';
    }
  }
  return 'w-full order-2';
};

const getTitleClass = (position: string) => {
  if (position === 'top') {
    return 'w-full order-1';
  }
  if (position === 'bottom') {
    return 'w-full order-3';
  }
  return 'w-full';
};

const getContentClass = (position: string) => {
  if (position === 'top') {
    return 'w-full order-2';
  }
  if (position === 'bottom') {
    return 'w-full order-4';
  }
  return 'w-full';
};

const getTextContainerClass = () => {
  const hasLeft = titlePosition.value === 'left' || contentPosition.value === 'left';
  const hasRight = titlePosition.value === 'right' || contentPosition.value === 'right';
  
  if (hasLeft) {
    return 'w-full md:w-1/2 order-1 md:order-1';
  }
  if (hasRight) {
    return 'w-full md:w-1/2 order-1 md:order-1';
  }
  return 'w-full';
};

const formatDate = (date: string) => {
  return new Date(date).toLocaleDateString('ru-RU', {
    year: 'numeric',
    month: 'long',
    day: 'numeric'
  });
};
</script>

<style scoped>
.post-content {
  @apply w-full;
}
</style>

