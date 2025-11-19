<template>
  <div class="space-y-4">
    <div>
      <label class="block text-sm font-medium text-gray-700 mb-1">Заголовок</label>
      <ui-input 
        v-model="formData.title" 
        variant="primary" 
        type="text" 
        placeholder="Введите заголовок"
      />
    </div>
    
    <div>
      <label class="block text-sm font-medium text-gray-700 mb-1">Текст</label>
      <textarea 
        v-model="formData.text" 
        rows="6"
        class="w-full px-3 py-2 border border-gray-300 rounded-md focus:border-blue-500 focus:ring-2 focus:ring-blue-200 transition-all placeholder:text-gray-400"
        placeholder="Введите текст блока"
      ></textarea>
    </div>
    
    <div>
      <label class="block text-sm font-medium text-gray-700 mb-1">Изображение</label>
      <div class="flex gap-2">
        <ui-input 
          v-model="formData.image" 
          variant="primary" 
          type="text" 
          placeholder="/upload/image.jpg или загрузите файл"
          class="flex-1"
        />
        <div class="relative">
          <input
            type="file"
            ref="fileInput"
            accept="image/*"
            @change="handleFileUpload"
            class="hidden"
            id="text-image-file"
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
      <p class="text-xs text-gray-500 mt-1">Поддерживаются прямые ссылки на изображения или загрузка файла</p>
    </div>
    
    <div>
      <label class="block text-sm font-medium text-gray-700 mb-1">Расположение текста относительно изображения</label>
      <select 
        v-model="formData.textPosition" 
        class="w-full px-3 py-2 border border-gray-300 rounded-md focus:border-blue-500 focus:ring-2 focus:ring-blue-200 transition-all"
      >
        <option value="right">Справа от изображения</option>
        <option value="left">Слева от изображения</option>
        <option value="top">Над изображением</option>
        <option value="bottom">Под изображением</option>
        <option value="overlay">Поверх изображения (overlay)</option>
      </select>
      <p class="text-xs text-gray-500 mt-1">Выберите, где будет расположен текст относительно изображения</p>
    </div>
    
    <div>
      <label class="block text-sm font-medium text-gray-700 mb-1">Выравнивание текста</label>
      <select 
        v-model="formData.alignment" 
        class="w-full px-3 py-2 border border-gray-300 rounded-md focus:border-blue-500 focus:ring-2 focus:ring-blue-200 transition-all"
      >
        <option value="left">Слева</option>
        <option value="center">По центру</option>
        <option value="right">Справа</option>
        <option value="justify">По ширине</option>
      </select>
    </div>
    
    <div>
      <label class="block text-sm font-medium text-gray-700 mb-1">Стиль</label>
      <select 
        v-model="formData.style" 
        class="w-full px-3 py-2 border border-gray-300 rounded-md focus:border-blue-500 focus:ring-2 focus:ring-blue-200 transition-all"
      >
        <option value="regular">Обычный</option>
        <option value="bold">Жирный</option>
        <option value="italic">Курсив</option>
        <option value="heading">Заголовок</option>
      </select>
    </div>
    
    <div>
      <label class="block text-sm font-medium text-gray-700 mb-1">Размер изображения</label>
      <select 
        v-model="formData.imageSize" 
        class="w-full px-3 py-2 border border-gray-300 rounded-md focus:border-blue-500 focus:ring-2 focus:ring-blue-200 transition-all"
      >
        <option value="small">Маленький</option>
        <option value="medium">Средний</option>
        <option value="large">Большой</option>
        <option value="xlarge">Крупный</option>
      </select>
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
  modelValue: TextWithImageBlockContent;
}>();

const emit = defineEmits<{
  'update:modelValue': [value: TextWithImageBlockContent];
}>();

const formData = ref<TextWithImageBlockContent>({
  title: props.modelValue?.title || '',
  text: props.modelValue?.text || '',
  image: props.modelValue?.image || '',
  textPosition: props.modelValue?.textPosition || 'right',
  alignment: props.modelValue?.alignment || 'left',
  style: props.modelValue?.style || 'regular',
  imageSize: props.modelValue?.imageSize || 'medium',
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
    formData.value.image = url;
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

