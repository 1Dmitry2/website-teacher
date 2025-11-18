<template>
  <div class="block-gallery mb-8">
    <div :class="galleryClass">
      <div 
        v-for="(image, index) in content.images" 
        :key="index"
        class="relative group cursor-pointer overflow-hidden rounded-lg"
        :class="imageContainerClass"
        :style="imageContainerStyle"
        @click="openLightbox(index)"
      >
        <img 
          :src="getImageUrl(image.src)" 
          :alt="image.caption || `Изображение ${index + 1}`"
          :class="`${imageClass} ${getImageSizeClass(image)}`"
        />
        <div v-if="image.caption" class="absolute bottom-0 left-0 right-0 bg-black bg-opacity-75 text-white p-3 opacity-0 group-hover:opacity-100 transition-opacity">
          <p class="text-sm">{{ image.caption }}</p>
        </div>
      </div>
    </div>
    
    <Transition name="lightbox-fade">
      <div 
        v-if="lightboxOpen"
        class="fixed inset-0 bg-black/95 backdrop-blur-sm z-50 flex items-center justify-center p-4"
        @click="closeLightbox"
      >
        <button 
          @click.stop="closeLightbox"
          class="absolute top-2 right-2 sm:top-4 sm:right-4 text-white hover:text-gray-300 text-4xl sm:text-5xl font-light transition-all duration-200 hover:scale-110 z-10 touch-manipulation"
          aria-label="Закрыть"
        >
          &times;
        </button>
        <button 
          v-if="content.images.length > 1"
          @click.stop="previousImage"
          class="absolute left-2 sm:left-4 text-white hover:text-gray-300 text-3xl sm:text-5xl font-light transition-all duration-200 hover:scale-110 z-10 bg-black/30 hover:bg-black/50 rounded-full w-10 h-10 sm:w-12 sm:h-12 flex items-center justify-center touch-manipulation"
          aria-label="Предыдущее изображение"
        >
          &#8249;
        </button>
        <button 
          v-if="content.images.length > 1"
          @click.stop="nextImage"
          class="absolute right-2 sm:right-4 text-white hover:text-gray-300 text-3xl sm:text-5xl font-light transition-all duration-200 hover:scale-110 z-10 bg-black/30 hover:bg-black/50 rounded-full w-10 h-10 sm:w-12 sm:h-12 flex items-center justify-center touch-manipulation"
          aria-label="Следующее изображение"
        >
          &#8250;
        </button>
        <Transition name="image-scale" mode="out-in">
          <div v-if="content.images && content.images[lightboxIndex]" :key="lightboxIndex" class="max-w-4xl max-h-full">
            <img 
              :src="getImageUrl(content.images[lightboxIndex]?.src || '')" 
              :alt="content.images[lightboxIndex]?.caption || `Изображение ${lightboxIndex + 1}`"
              class="max-w-full max-h-[90vh] object-contain rounded-lg shadow-2xl"
            />
            <p v-if="content.images[lightboxIndex]?.caption" class="text-white text-center mt-2 sm:mt-4 text-sm sm:text-lg px-4">
              {{ content.images[lightboxIndex]?.caption }}
            </p>
          </div>
        </Transition>
      </div>
    </Transition>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue';
import { API_BASE_URL } from '../../api/client';

import type { MediaSize } from '../../api/client';

export interface GalleryImage {
  src: string;
  caption: string;
  size?: MediaSize;
}

export interface GalleryBlockContent {
  images: GalleryImage[];
  layout?: 'grid' | 'flex' | 'single-row';
  columns?: number; // для grid: 1-6, для flex: количество элементов в строке
}

const props = defineProps<{
  content: GalleryBlockContent;
}>();

const lightboxOpen = ref(false);
const lightboxIndex = ref(0);

const galleryClass = computed(() => {
  const layout = props.content?.layout || 'grid';
  const columns = props.content?.columns || 3;
  
  if (layout === 'single-row') {
    return 'flex gap-4 overflow-x-auto';
  }
  
  if (layout === 'flex') {
    return `flex flex-wrap gap-4`;
  }
  
  const gridCols = {
    1: 'grid-cols-1',
    2: 'grid-cols-1 sm:grid-cols-2',
    3: 'grid-cols-1 sm:grid-cols-2 lg:grid-cols-3',
    4: 'grid-cols-1 sm:grid-cols-2 lg:grid-cols-4',
    5: 'grid-cols-1 sm:grid-cols-2 lg:grid-cols-5',
    6: 'grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-6',
  };
  
  return `grid ${gridCols[columns as keyof typeof gridCols] || gridCols[3]} gap-4`;
});

const imageContainerClass = computed(() => {
  const layout = props.content?.layout || 'grid';
  if (layout === 'flex') {
    return 'w-full';
  }
  return '';
});

const imageContainerStyle = computed(() => {
  const layout = props.content?.layout || 'grid';
  if (layout === 'flex') {
    const columns = props.content?.columns || 3;
    const gap = 16;
    return {
      flex: `0 0 calc(${100 / columns}% - ${gap * (columns - 1) / columns}px)`,
      maxWidth: `calc(${100 / columns}% - ${gap * (columns - 1) / columns}px)`,
    };
  }
  return {};
});

const getImageSizeClass = (image: GalleryImage): string => {
  const size = image.size || 'medium';
  const sizeMap: Record<MediaSize, string> = {
    small: 'h-32 sm:h-40',
    medium: 'h-48 sm:h-64',
    large: 'h-64 sm:h-80',
    xlarge: 'h-80 sm:h-96',
  };
  return sizeMap[size] || sizeMap.medium;
};

const imageClass = computed(() => {
  const layout = props.content?.layout || 'grid';
  const baseClass = 'object-cover transition-transform duration-300 group-hover:scale-110 w-full';
  
  if (layout === 'single-row') {
    return `flex-shrink-0 w-64 h-64 ${baseClass}`;
  }
  
  return baseClass;
});

const getImageUrl = (url: string) => {
  if (!url) return '';
  if (url.startsWith('http://') || url.startsWith('https://')) {
    return url;
  }
  if (url.startsWith('/')) {
    return `${API_BASE_URL}${url}`;
  }
  return `${API_BASE_URL}/uploads/${url}`;
};

const openLightbox = (index: number) => {
  lightboxIndex.value = index;
  lightboxOpen.value = true;
  document.body.style.overflow = 'hidden';
};

const closeLightbox = () => {
  lightboxOpen.value = false;
  document.body.style.overflow = '';
};

const nextImage = () => {
  if (props.content.images.length === 0) return;
  lightboxIndex.value = (lightboxIndex.value + 1) % props.content.images.length;
};

const previousImage = () => {
  if (props.content.images.length === 0) return;
  lightboxIndex.value = lightboxIndex.value === 0 ? props.content.images.length - 1 : lightboxIndex.value - 1;
};
</script>

<style scoped>
.lightbox-fade-enter-active,
.lightbox-fade-leave-active {
  transition: opacity 0.3s ease;
}

.lightbox-fade-enter-from,
.lightbox-fade-leave-to {
  opacity: 0;
}

.image-scale-enter-active {
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
}

.image-scale-leave-active {
  transition: all 0.25s cubic-bezier(0.4, 0, 0.2, 1);
}

.image-scale-enter-from {
  opacity: 0;
  transform: scale(0.9);
}

.image-scale-leave-to {
  opacity: 0;
  transform: scale(0.95);
}
</style>

