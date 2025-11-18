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
  if (!props.content) return 'text-2xl font-normal mb-4';
  const styleMap = {
    regular: 'text-2xl font-normal mb-4',
    bold: 'text-2xl font-bold mb-4',
    italic: 'text-2xl font-normal italic mb-4',
    heading: 'text-3xl font-bold mb-6',
  };
  return styleMap[props.content.style] || 'text-2xl font-normal mb-4';
});

const textClass = computed(() => {
  if (!props.content) return 'text-base font-normal';
  const styleMap = {
    regular: 'text-base font-normal',
    bold: 'text-base font-bold',
    italic: 'text-base font-normal italic',
    heading: 'text-lg font-semibold',
  };
  return styleMap[props.content.style] || 'text-base font-normal';
});

const formattedText = computed(() => {
  if (!props.content?.text) return '';
  // Простое форматирование: заменяем переносы строк на <br>
  return props.content.text.replace(/\n/g, '<br>');
});
</script>


