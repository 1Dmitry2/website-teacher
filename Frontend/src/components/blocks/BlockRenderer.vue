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
    <BlockDisplayTextWithImage 
      v-else-if="block.type === 'text-with-image' && parsedContent"
      :content="parsedContent as any"
    />
    <BlockDisplayDocument
      v-else-if="block.type === 'document' && parsedContent"
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
import BlockDisplayTextWithImage from './BlockDisplayTextWithImage.vue';
import BlockDisplayDocument from './BlockDisplayDocument.vue';

export interface Block {
  id: string;
  page: string;
  pages: string[];
  type: 'text' | 'slider' | 'gallery' | 'video' | 'text-with-image' | 'document';
  content: Record<string, any> | string;
  display_order: number;
  created_at: string;
  updated_at: string;
}

const props = defineProps<{
  block: Block;
}>();

const parsedContent = computed(() => {
  if (!props.block.content) return null;
  
  if (typeof props.block.content === 'object' && !Array.isArray(props.block.content)) {
    return props.block.content;
  }
  
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


