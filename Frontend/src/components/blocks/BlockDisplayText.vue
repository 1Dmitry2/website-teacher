<template>
  <div v-if="content" class="block-text mb-8" :class="textAlignmentClass">
    <h2 v-if="content.title" :class="titleClass">{{ content.title }}</h2>
    <div v-if="content.text" :class="textClass" v-html="formattedText"></div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue';

export interface TextBlockContent {
  title: string;
  text: string;
  alignment: 'left' | 'center' | 'right' | 'justify';
  style: 'regular' | 'bold' | 'italic' | 'heading';
}

const props = defineProps<{
  content: TextBlockContent;
}>();

const textAlignmentClass = computed(() => {
  if (!props.content) return 'text-left';
  const alignmentMap = {
    left: 'text-left',
    center: 'text-center',
    right: 'text-right',
    justify: 'text-justify',
  };
  return alignmentMap[props.content.alignment] || 'text-left';
});

const titleClass = computed(() => {
  if (!props.content) return 'text-xl sm:text-2xl font-normal mb-3 sm:mb-4';
  const styleMap = {
    regular: 'text-xl sm:text-2xl font-normal mb-3 sm:mb-4',
    bold: 'text-xl sm:text-2xl font-bold mb-3 sm:mb-4',
    italic: 'text-xl sm:text-2xl font-normal italic mb-3 sm:mb-4',
    heading: 'text-2xl sm:text-3xl font-bold mb-4 sm:mb-6',
  };
  return styleMap[props.content.style] || 'text-xl sm:text-2xl font-normal mb-3 sm:mb-4';
});

const textClass = computed(() => {
  if (!props.content) return 'text-sm sm:text-base font-normal';
  const styleMap = {
    regular: 'text-sm sm:text-base font-normal',
    bold: 'text-sm sm:text-base font-bold',
    italic: 'text-sm sm:text-base font-normal italic',
    heading: 'text-base sm:text-lg font-semibold',
  };
  return styleMap[props.content.style] || 'text-sm sm:text-base font-normal';
});

const formattedText = computed(() => {
  if (!props.content?.text) return '';
  return props.content.text.replace(/\n/g, '<br>');
});
</script>


