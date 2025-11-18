<template>
  <div class="block-renderer w-full">
    <BlockDisplayText 
      v-if="block.type === 'text' && parsedContent"
      :content="parsedContent as any"
    />
    <BlockDisplaySlider 
      v-else-if="block.type === 'slider' && parsedContent"
      :content="parsedContent as any"
    />
    <BlockDisplayGallery 
      v-else-if="block.type === 'gallery' && parsedContent"
      :content="parsedContent as any"
    />
    <BlockDisplayVideo 
      v-else-if="block.type === 'video' && parsedContent"
      :content="parsedContent as any"
    />
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue';
import BlockDisplayText from './BlockDisplayText.vue';
import BlockDisplaySlider from './BlockDisplaySlider.vue';
import BlockDisplayGallery from './BlockDisplayGallery.vue';
import BlockDisplayVideo from './BlockDisplayVideo.vue';

export interface Block {
  id: string;
  page: string;
  pages: string[];
  type: 'text' | 'slider' | 'gallery' | 'video';
  content: Record<string, any> | string;
  display_order: number;
  created_at: string;
  updated_at: string;
}

const props = defineProps<{
  block: Block;
}>();

// Парсим content, если он пришел как строка JSON
const parsedContent = computed(() => {
  if (!props.block.content) return null;
  
  // Если content уже объект, возвращаем его
  if (typeof props.block.content === 'object' && !Array.isArray(props.block.content)) {
    return props.block.content;
  }
  
  // Если content - строка, пытаемся распарсить
  if (typeof props.block.content === 'string') {
    try {
      return JSON.parse(props.block.content);
    } catch (e) {
      console.error('Error parsing block content:', e);
      return null;
    }
  }
  
  return props.block.content;
});
</script>


