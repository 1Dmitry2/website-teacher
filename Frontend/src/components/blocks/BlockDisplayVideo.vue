<template>
  <div class="block-video mb-8">
    <div class="relative w-full" style="padding-bottom: 56.25%;">
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
import { API_BASE_URL } from '../../api/client';

export interface VideoBlockContent {
  url: string;
  autoplay: boolean;
}

const props = defineProps<{
  content: VideoBlockContent;
}>();

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


