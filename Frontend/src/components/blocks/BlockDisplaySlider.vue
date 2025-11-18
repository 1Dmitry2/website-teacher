<template>
  <div class="block-slider mb-8">
    <div class="relative w-full overflow-hidden rounded-lg" style="height: 400px;">
      <div 
        v-for="(slide, index) in content.slides" 
        :key="index"
        v-show="currentSlide === index"
        class="absolute inset-0 transition-opacity duration-500"
        :class="{ 'opacity-100': currentSlide === index, 'opacity-0': currentSlide !== index }"
      >
        <img 
          :src="getImageUrl(slide.image)" 
          :alt="slide.title || `Слайд ${index + 1}`"
          class="w-full h-full object-cover"
        />
        <div v-if="slide.title || slide.subtitle" class="absolute bottom-0 left-0 right-0 bg-black bg-opacity-50 text-white p-6">
          <h3 v-if="slide.title" class="text-2xl font-bold mb-2">{{ slide.title }}</h3>
          <p v-if="slide.subtitle" class="text-lg">{{ slide.subtitle }}</p>
        </div>
      </div>
      
      <!-- Навигация -->
      <button 
        v-if="content.slides.length > 1"
        @click="previousSlide"
        class="absolute left-4 top-1/2 transform -translate-y-1/2 bg-white bg-opacity-75 hover:bg-opacity-100 rounded-full p-2 transition-all"
        aria-label="Предыдущий слайд"
      >
        <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7" />
        </svg>
      </button>
      <button 
        v-if="content.slides.length > 1"
        @click="nextSlide"
        class="absolute right-4 top-1/2 transform -translate-y-1/2 bg-white bg-opacity-75 hover:bg-opacity-100 rounded-full p-2 transition-all"
        aria-label="Следующий слайд"
      >
        <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7" />
        </svg>
      </button>
      
      <!-- Индикаторы -->
      <div v-if="content.slides.length > 1" class="absolute bottom-4 left-1/2 transform -translate-x-1/2 flex gap-2">
        <button
          v-for="(slide, index) in content.slides"
          :key="index"
          @click="currentSlide = index"
          class="w-3 h-3 rounded-full transition-all"
          :class="currentSlide === index ? 'bg-white' : 'bg-white bg-opacity-50'"
          :aria-label="`Перейти к слайду ${index + 1}`"
        />
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue';
import { API_BASE_URL } from '../../api/client';

export interface SliderSlide {
  image: string;
  title?: string;
  subtitle?: string;
}

export interface SliderBlockContent {
  slides: SliderSlide[];
}

const props = defineProps<{
  content: SliderBlockContent;
}>();

const currentSlide = ref(0);
let autoSlideInterval: number | null = null;

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

const nextSlide = () => {
  if (props.content.slides.length === 0) return;
  currentSlide.value = (currentSlide.value + 1) % props.content.slides.length;
};

const previousSlide = () => {
  if (props.content.slides.length === 0) return;
  currentSlide.value = currentSlide.value === 0 ? props.content.slides.length - 1 : currentSlide.value - 1;
};

onMounted(() => {
  if (props.content.slides.length > 1) {
    autoSlideInterval = window.setInterval(() => {
      nextSlide();
    }, 5000);
  }
});

onUnmounted(() => {
  if (autoSlideInterval !== null) {
    clearInterval(autoSlideInterval);
  }
});
</script>


