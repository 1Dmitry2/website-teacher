<template>
  <div class="space-y-4">
    <div>
      <label class="block text-sm font-medium text-gray-700 mb-1">URL видео</label>
      <div class="flex gap-2">
        <ui-input 
          v-model="formData.url" 
          variant="primary" 
          type="text" 
          placeholder="https://youtube.com/watch?v=... или /upload/video.mp4"
          class="flex-1"
        />
        <div class="relative">
          <input
            type="file"
            ref="fileInput"
            accept="video/*"
            @change="handleFileUpload"
            class="hidden"
            id="video-file"
          />
          <ui-button 
            type="button" 
            variant="secondary"
            :disabled="uploading"
            @click="() => fileInput?.click()"
          >
            {{ uploading ? 'Загрузка...' : 'Загрузить' }}
          </ui-button>
        </div>
      </div>
      <p class="text-xs text-gray-500 mt-1">Поддерживаются ссылки на YouTube, прямые ссылки на видеофайлы или загрузка файла</p>
    </div>
    
    <div class="flex items-center">
      <input 
        type="checkbox" 
        v-model="formData.autoplay"
        id="autoplay"
        class="w-4 h-4 text-blue-600 border-gray-300 rounded focus:ring-blue-500"
      />
      <label for="autoplay" class="ml-2 text-sm font-medium text-gray-700">
        Автовоспроизведение
      </label>
    </div>
    
    <div>
      <label class="block text-sm font-medium text-gray-700 mb-1">Выравнивание</label>
      <select 
        v-model="formData.alignment" 
        class="w-full px-3 py-2 border border-gray-300 rounded-md focus:border-blue-500 focus:ring-2 focus:ring-blue-200 transition-all"
      >
        <option value="full-width">На всю ширину</option>
        <option value="left">Слева</option>
        <option value="center">По центру</option>
        <option value="right">Справа</option>
      </select>
    </div>
    
    <div>
      <label class="block text-sm font-medium text-gray-700 mb-1">Размер видео</label>
      <select 
        v-model="formData.size" 
        class="w-full px-3 py-2 border border-gray-300 rounded-md focus:border-blue-500 focus:ring-2 focus:ring-blue-200 transition-all"
      >
        <option value="small">Маленький</option>
        <option value="medium">Средний</option>
        <option value="large">Большой</option>
        <option value="xlarge">Крупный</option>
      </select>
    </div>
    
    <div v-if="formData.alignment && formData.alignment !== 'full-width'">
      <label class="block text-sm font-medium text-gray-700 mb-1">Максимальная ширина (пиксели)</label>
      <ui-input 
        :model-value="formData.maxWidth || 0"
        @update:model-value="(val: string | number) => formData.maxWidth = typeof val === 'number' ? (val || undefined) : (Number(val) || undefined)"
        variant="primary" 
        type="number" 
        placeholder="800"
      />
      <p class="text-xs text-gray-500 mt-1">Оставьте пустым для автоматической ширины</p>
    </div>

    <NotificationModal
      :is-open="notification.isOpen.value"
      :type="notification.options.value.type"
      :title="notification.options.value.title"
      :message="notification.options.value.message"
      @close="notification.close"
    />
  </div>
</template>

<script setup lang="ts">
import { nextTick, ref, watch } from 'vue';
import UiInput from '../ui/Ui-input.vue';
import UiButton from '../ui/Ui-button.vue';
import NotificationModal from '../ui/NotificationModal.vue';
import { adminApi } from '../../api/admin';
import { useNotification } from '../../composables/useModal';

import type { MediaSize } from '../../api/client';

export interface VideoBlockContent {
  url: string;
  autoplay: boolean;
  alignment?: 'left' | 'center' | 'right' | 'full-width';
  maxWidth?: number;
  size?: MediaSize;
}

const props = defineProps<{
  modelValue: VideoBlockContent;
}>();

const emit = defineEmits<{
  'update:modelValue': [value: VideoBlockContent];
}>();

const formData = ref<VideoBlockContent>({
  url: props.modelValue?.url || '',
  autoplay: props.modelValue?.autoplay ?? false,
  alignment: props.modelValue?.alignment || 'full-width',
  maxWidth: props.modelValue?.maxWidth,
  size: props.modelValue?.size || 'medium',
});

const fileInput = ref<HTMLInputElement | null>(null);
const uploading = ref(false);
const notification = useNotification();
const syncingFromModel = ref(false);

const handleFileUpload = async (event: Event) => {
  const target = event.target as HTMLInputElement;
  const file = target.files?.[0];
  if (!file) return;

  uploading.value = true;
  try {
    const url = await adminApi.uploadFile(file);
    formData.value.url = url;
    emit('update:modelValue', { ...formData.value });
  } catch (error) {
    notification.error(error instanceof Error ? error.message : 'Ошибка загрузки файла');
  } finally {
    uploading.value = false;
    if (target) {
      target.value = '';
    }
  }
};

watch(formData, (newVal) => {
  if (syncingFromModel.value) {
    return;
  }
  emit('update:modelValue', { ...newVal });
}, { deep: true });

watch(() => props.modelValue, (newVal) => {
  if (!newVal) {
    return;
  }
  syncingFromModel.value = true;
  formData.value = { ...newVal };
  nextTick(() => {
    syncingFromModel.value = false;
  });
}, { deep: true });
</script>

