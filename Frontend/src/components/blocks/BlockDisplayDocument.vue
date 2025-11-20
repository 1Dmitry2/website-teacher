<template>
  <div v-if="content" class="block-document w-full mb-8 bg-white rounded-xl shadow p-6">
    <div class="flex flex-col gap-4">
      <div>
        <h2 v-if="content.title" class="text-2xl font-semibold text-gray-900 mb-2">
          {{ content.title }}
        </h2>
        <p v-if="content.description" class="text-gray-600">
          {{ content.description }}
        </p>
      </div>

      <div v-if="content.mode === 'viewer' && viewerUrl" class="space-y-4">
        <div
          :class="content.viewerLayout === 'preview' ? 'border border-gray-200 rounded-lg overflow-hidden bg-gray-50' : 'rounded-lg overflow-hidden bg-gray-50'"
          :style="previewStyle"
        >
          <iframe
            v-if="viewerUrl"
            :src="viewerUrl"
            :title="content.title || 'Документ'"
            frameborder="0"
            class="w-full h-full"
            loading="lazy"
            allowfullscreen
          ></iframe>
        </div>

        <div v-if="content.showDownloadButton !== false" class="flex gap-3 items-center flex-wrap">
          <a
            v-if="resolvedDocumentUrl"
            :href="resolvedDocumentUrl"
            target="_blank"
            rel="noopener noreferrer"
            class="inline-flex items-center justify-center px-5 py-3 bg-indigo-600 text-white font-medium rounded-lg hover:bg-indigo-700 transition-colors"
            :download="downloadName"
          >
            Скачать документ
          </a>
          <span class="text-sm text-gray-500">
            {{ displayFileName }}
          </span>
        </div>
      </div>

      <div v-else class="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-3">
        <div class="flex items-center gap-3">
          <div class="flex-shrink-0 w-12 h-12 rounded-full bg-indigo-100 text-indigo-600 flex items-center justify-center text-2xl">
            📄
          </div>
          <div>
            <div class="font-medium text-gray-900">
              {{ displayFileName }}
            </div>
            <div class="text-sm text-gray-500">
              {{ content.description || 'Документ доступен для скачивания' }}
            </div>
          </div>
        </div>
        <div class="flex gap-2">
          <a
            v-if="resolvedDocumentUrl"
            :href="resolvedDocumentUrl"
            class="inline-flex items-center justify-center px-4 py-2 bg-indigo-600 text-white font-medium rounded-lg hover:bg-indigo-700 transition-colors"
            target="_blank"
            rel="noopener noreferrer"
            :download="downloadName"
          >
            Скачать
          </a>
          <a
            v-if="resolvedDocumentUrl"
            :href="resolvedDocumentUrl"
            class="inline-flex items-center justify-center px-4 py-2 border border-gray-300 text-gray-700 font-medium rounded-lg hover:bg-gray-50 transition-colors"
            target="_blank"
            rel="noopener noreferrer"
          >
            Открыть
          </a>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue';
import { API_BASE_URL } from '../../api/client';
import type { DocumentBlockContent } from './BlockFormDocument.vue';

const props = defineProps<{
  content: DocumentBlockContent;
}>();

const resolveDocumentUrl = (url: string) => {
  if (!url) return '';
  if (url.startsWith('http://') || url.startsWith('https://')) {
    return url;
  }
  if (url.startsWith('/')) {
    return `${API_BASE_URL}${url}`;
  }
  return `${API_BASE_URL}/uploads/${url}`;
};

const resolvedDocumentUrl = computed(() => resolveDocumentUrl(props.content?.fileUrl || ''));

const isPdf = computed(() => {
  const url = props.content?.fileUrl?.toLowerCase() || '';
  return url.endsWith('.pdf');
});

const viewerUrl = computed(() => {
  const url = resolvedDocumentUrl.value;
  if (!url) {
    return '';
  }
  if (isPdf.value) {
    return url + '#view=fitH';
  }
  return `https://docs.google.com/gview?embedded=1&url=${encodeURIComponent(url)}`;
});

const previewHeight = computed(() => props.content?.previewHeight || 480);

const previewStyle = computed(() => {
  if (props.content?.viewerLayout === 'preview') {
    return `height: ${previewHeight.value}px;`;
  }
  return 'min-height: 600px;';
});

const fallbackFileName = computed(() => {
  const url = resolvedDocumentUrl.value;
  if (!url) return 'Документ';
  try {
    const parts = url.split('/');
    return decodeURIComponent(parts[parts.length - 1] || 'Документ');
  } catch {
    return 'Документ';
  }
});

const displayFileName = computed(() => props.content?.fileName || props.content?.title || fallbackFileName.value);
const downloadName = computed(() => props.content?.fileName || fallbackFileName.value);
</script>

