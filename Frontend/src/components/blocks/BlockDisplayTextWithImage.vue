<template>
  <div v-if="content" class="block-text-with-image mb-8" :class="containerClass">
    <div :class="layoutClass">
      <div :class="imageContainerClass">
        <img 
          :src="getImageUrl(content.image)" 
          :alt="content.title || 'Изображение'"
          :class="imageClass"
        />
      </div>
      
      <div :class="textContainerClass">
        <h2 v-if="content.title" :class="titleClass">{{ content.title }}</h2>
        <div v-if="content.text" :class="textClass" v-html="formattedText"></div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue';
import { API_BASE_URL } from '../../api/client';

import type { MediaSize } from '../../api/client';

export interface TextWithImageBlockContent {
  title: string;
  text: string;
  image: string;
  textPosition: 'left' | 'right' | 'top' | 'bottom' | 'overlay';
  alignment: 'left' | 'center' | 'right' | 'justify';
  style: 'regular' | 'bold' | 'italic' | 'heading';
  imageSize?: MediaSize;
}

const props = defineProps<{
  content: TextWithImageBlockContent;
}>();

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

const containerClass = computed(() => {
  if (props.content?.textPosition === 'overlay') {
    return 'relative';
  }
  return 'w-full';
});

const layoutClass = computed(() => {
  const position = props.content?.textPosition || 'right';
  
  if (position === 'overlay') {
    return 'relative';
  }
  
  if (position === 'top' || position === 'bottom') {
    return 'flex flex-col gap-4';
  }
  
  return 'flex flex-col md:flex-row gap-4 lg:gap-6';
});

const imageContainerClass = computed(() => {
  const position = props.content?.textPosition || 'right';
  
  if (position === 'overlay') {
    return 'absolute inset-0 z-0 w-full h-full';
  }
  
  if (position === 'top') {
    return 'w-full';
  }
  
  if (position === 'bottom') {
    return 'w-full order-2';
  }
  
  if (position === 'left') {
    return 'w-full md:w-1/2 order-2 md:order-1';
  }
  
  return 'w-full md:w-1/2 order-1 md:order-2';
});

const textContainerClass = computed(() => {
  const position = props.content?.textPosition || 'right';
  const alignment = props.content?.alignment || 'left';
  
  let classes = '';
  
  if (position === 'overlay') {
    classes = 'relative z-10 p-6 md:p-8';
    classes += ' bg-black bg-opacity-50 text-white rounded-lg';
  } else {
    classes = 'w-full';
    if (position === 'left') {
      classes += ' md:w-1/2 order-1 md:order-2';
    } else if (position === 'right') {
      classes += ' md:w-1/2 order-2 md:order-1';
    } else {
      classes += ' order-1';
    }
  }
  
  const alignmentMap = {
    left: 'text-left',
    center: 'text-center',
    right: 'text-right',
    justify: 'text-justify',
  };
  classes += ' ' + (alignmentMap[alignment] || 'text-left');
  
  return classes;
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

const getImageSizeClass = (): string => {
  const size = props.content?.imageSize || 'medium';
  const sizeMap: Record<MediaSize, string> = {
    small: 'min-h-[150px] sm:min-h-[200px] md:min-h-[250px]',
    medium: 'min-h-[200px] sm:min-h-[300px] md:min-h-[400px]',
    large: 'min-h-[300px] sm:min-h-[400px] md:min-h-[500px]',
    xlarge: 'min-h-[400px] sm:min-h-[500px] md:min-h-[600px]',
  };
  return sizeMap[size] || sizeMap.medium;
};

const imageClass = computed(() => {
  const position = props.content?.textPosition || 'right';
  
  if (position === 'overlay') {
    return 'w-full h-full object-cover rounded-lg';
  }
  
  return `w-full h-full object-cover rounded-lg ${getImageSizeClass()}`;
});

const formattedText = computed(() => {
  if (!props.content?.text) return '';
  return props.content.text.replace(/\n/g, '<br>');
});
</script>

<style scoped>
.block-text-with-image {
  @apply w-full;
}

.block-text-with-image .relative {
  min-height: 400px;
  position: relative;
}

.block-text-with-image .absolute.inset-0 {
  min-height: 400px;
}
</style>

