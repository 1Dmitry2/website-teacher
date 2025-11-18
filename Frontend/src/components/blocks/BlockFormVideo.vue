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
  </div>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue';
import UiInput from '../ui/Ui-input.vue';
import UiButton from '../ui/Ui-button.vue';
import { adminApi } from '../../api/admin';

export interface VideoBlockContent {
  url: string;
  autoplay: boolean;
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
});

const fileInput = ref<HTMLInputElement | null>(null);
const uploading = ref(false);

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
    alert(error instanceof Error ? error.message : 'Ошибка загрузки файла');
  } finally {
    uploading.value = false;
    // Сбрасываем input, чтобы можно было загрузить тот же файл снова
    if (target) {
      target.value = '';
    }
  }
};

watch(formData, (newVal) => {
  emit('update:modelValue', { ...newVal });
}, { deep: true });

watch(() => props.modelValue, (newVal) => {
  if (newVal) {
    formData.value = { ...newVal };
  }
}, { deep: true });
</script>

