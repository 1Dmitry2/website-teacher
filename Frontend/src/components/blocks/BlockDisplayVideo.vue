<template>
  <div class="block-video mb-8" :class="alignmentClass">
    <div class="relative" :class="containerClass" :style="containerStyle">
      <iframe
        v-if="isYouTubeUrl(content.url)"
        :src="getYouTubeEmbedUrl(content.url)"
        class="absolute top-0 left-0 w-full h-full rounded-lg"
        frameborder="0"
        allow="accelerometer; autoplay; clipboard-write; encrypted-media; gyroscope; picture-in-picture"
        allowfullscreen
      ></iframe>
      <video
        v-else
        :src="getVideoUrl(content.url)"
        :autoplay="content.autoplay"
        :controls="true"
        class="absolute top-0 left-0 w-full h-full rounded-lg"
      ></video>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue';
import { API_BASE_URL } from '../../api/client';

import type { MediaSize } from '../../api/client';

export interface VideoBlockContent {
  url: string;
  autoplay: boolean;
  alignment?: 'left' | 'center' | 'right' | 'full-width';
  maxWidth?: number;
  size?: MediaSize;
}

const props = defineProps<{
  content: VideoBlockContent;
}>();

const alignmentClass = computed(() => {
  if (!props.content?.alignment || props.content.alignment === 'full-width') {
    return 'w-full';
  }
  const alignmentMap = {
    left: 'flex justify-start',
    center: 'flex justify-center',
    right: 'flex justify-end',
  };
  return alignmentMap[props.content.alignment] || 'w-full';
});

const containerClass = computed(() => {
  if (props.content?.alignment === 'full-width') {
    return 'w-full';
  }
  return `${getSizeClass()} w-full`;
});

const getSizeClass = (): string => {
  const size = props.content?.size || 'medium';
  const sizeMap: Record<MediaSize, string> = {
    small: 'max-w-xs',
    medium: 'max-w-md',
    large: 'max-w-2xl',
    xlarge: 'max-w-4xl',
  };
  return sizeMap[size] || sizeMap.medium;
};

const containerStyle = computed(() => {
  const maxWidth = props.content?.maxWidth;
  const size = props.content?.size || 'medium';
  
  let width = maxWidth;
  if (!width && props.content?.alignment && props.content.alignment !== 'full-width') {
    const sizeWidthMap: Record<MediaSize, number> = {
      small: 320,
      medium: 448,
      large: 672,
      xlarge: 896,
    };
    width = sizeWidthMap[size] || sizeWidthMap.medium;
  }
  
  if (width && props.content?.alignment && props.content.alignment !== 'full-width') {
    return {
      maxWidth: `${width}px`,
      paddingBottom: '56.25%',
    };
  }
  return {
    paddingBottom: '56.25%',
  };
});

const isYouTubeUrl = (url: string): boolean => {
  if (!url) return false;
  const youtubeRegex = /(?:youtube\.com\/(?:[^\/]+\/.+\/|(?:v|e(?:mbed)?)\/|.*[?&]v=)|youtu\.be\/)([^"&?\/\s]{11})/;
  return youtubeRegex.test(url);
};

const getYouTubeEmbedUrl = (url: string): string => {
  if (!url) return '';
  const youtubeRegex = /(?:youtube\.com\/(?:[^\/]+\/.+\/|(?:v|e(?:mbed)?)\/|.*[?&]v=)|youtu\.be\/)([^"&?\/\s]{11})/;
  const match = url.match(youtubeRegex);
  if (match && match[1]) {
    return `https://www.youtube.com/embed/${match[1]}${props.content.autoplay ? '?autoplay=1' : ''}`;
  }
  return url;
};

const getVideoUrl = (url: string): string => {
  if (!url) return '';
  if (url.startsWith('http://') || url.startsWith('https://')) {
    return url;
  }
  if (url.startsWith('/')) {
    return `${API_BASE_URL}${url}`;
  }
  return `${API_BASE_URL}/uploads/${url}`;
};
</script>


