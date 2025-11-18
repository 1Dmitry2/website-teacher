<template>
  <div class="block-gallery mb-8">
    <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-4">
      <div 
        v-for="(image, index) in content.images" 
        :key="index"
        class="relative group cursor-pointer overflow-hidden rounded-lg"
        @click="openLightbox(index)"
      >
        <img 
          :src="getImageUrl(image.src)" 
          :alt="image.caption || `Изображение ${index + 1}`"
          class="w-full h-64 object-cover transition-transform duration-300 group-hover:scale-110"
        />
        <div v-if="image.caption" class="absolute bottom-0 left-0 right-0 bg-black bg-opacity-75 text-white p-3 opacity-0 group-hover:opacity-100 transition-opacity">
          <p class="text-sm">{{ image.caption }}</p>
        </div>
      </div>
    </div>
    
    <!-- Lightbox -->
    <Transition name="lightbox-fade">
      <div 
        v-if="lightboxOpen"
        class="fixed inset-0 bg-black/95 backdrop-blur-sm z-50 flex items-center justify-center p-4"
        @click="closeLightbox"
      >
        <button 
          @click.stop="closeLightbox"
          class="absolute top-4 right-4 text-white hover:text-gray-300 text-5xl font-light transition-all duration-200 hover:scale-110 z-10"
          aria-label="Закрыть"
        >
          &times;
        </button>
        <button 
          v-if="content.images.length > 1"
          @click.stop="previousImage"
          class="absolute left-4 text-white hover:text-gray-300 text-5xl font-light transition-all duration-200 hover:scale-110 z-10 bg-black/30 hover:bg-black/50 rounded-full w-12 h-12 flex items-center justify-center"
          aria-label="Предыдущее изображение"
        >
          &#8249;
        </button>
        <button 
          v-if="content.images.length > 1"
          @click.stop="nextImage"
          class="absolute right-4 text-white hover:text-gray-300 text-5xl font-light transition-all duration-200 hover:scale-110 z-10 bg-black/30 hover:bg-black/50 rounded-full w-12 h-12 flex items-center justify-center"
          aria-label="Следующее изображение"
        >
          &#8250;
        </button>
        <Transition name="image-scale" mode="out-in">
          <div :key="lightboxIndex" class="max-w-4xl max-h-full">
            <img 
              :src="getImageUrl(content.images[lightboxIndex].src)" 
              :alt="content.images[lightboxIndex].caption || `Изображение ${lightboxIndex + 1}`"
              class="max-w-full max-h-[90vh] object-contain rounded-lg shadow-2xl"
            />
            <p v-if="content.images[lightboxIndex].caption" class="text-white text-center mt-4 text-lg">
              {{ content.images[lightboxIndex].caption }}
            </p>
          </div>
        </Transition>
      </div>
    </Transition>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue';
import { API_BASE_URL } from '../../api/client';

export interface GalleryImage {
  src: string;
  caption: string;
}

export interface GalleryBlockContent {
  images: GalleryImage[];
}

const props = defineProps<{
  content: GalleryBlockContent;
}>();

const lightboxOpen = ref(false);
const lightboxIndex = ref(0);

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
/* Анимация появления/исчезновения lightbox */
.lightbox-fade-enter-active,
.lightbox-fade-leave-active {
  transition: opacity 0.3s ease;
}

.lightbox-fade-enter-from,
.lightbox-fade-leave-to {
  opacity: 0;
}

/* Анимация масштабирования изображения */
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

