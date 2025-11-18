<template>
  <div class="space-y-4">
    <div class="space-y-3">
      <div class="flex items-center justify-between">
        <label class="block text-sm font-medium text-gray-700">Фотографии</label>
        <ui-button type="button" variant="primary" @click="addImage">
          + Добавить фотографию
        </ui-button>
      </div>
      
      <div v-for="(image, index) in formData.images" :key="index" class="border border-gray-300 rounded-lg p-4 space-y-3">
        <div class="flex items-center justify-between">
          <span class="text-sm font-medium text-gray-700">Фото {{ index + 1 }}</span>
          <div class="flex gap-2">
            <button 
              type="button"
              @click="moveUp(index)"
              :disabled="index === 0"
              class="text-blue-600 hover:text-blue-800 text-sm disabled:text-gray-400 disabled:cursor-not-allowed"
            >
              ↑
            </button>
            <button 
              type="button"
              @click="moveDown(index)"
              :disabled="index === formData.images.length - 1"
              class="text-blue-600 hover:text-blue-800 text-sm disabled:text-gray-400 disabled:cursor-not-allowed"
            >
              ↓
            </button>
            <button 
              type="button"
              @click="removeImage(index)"
              class="text-red-600 hover:text-red-800 text-sm"
            >
              Удалить
            </button>
          </div>
        </div>
        
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-1">URL изображения</label>
          <div class="flex gap-2">
            <ui-input 
              v-model="image.src" 
              variant="primary" 
              type="text" 
              placeholder="/upload/photo1.jpg или загрузите файл"
              class="flex-1"
            />
            <div class="relative">
              <input
                type="file"
                :ref="el => fileInputs[index] = el as HTMLInputElement"
                accept="image/*"
                @change="handleFileUpload(index, $event)"
                class="hidden"
                :id="`gallery-file-${index}`"
              />
              <ui-button 
                type="button" 
                variant="secondary"
                :disabled="uploading[index]"
                @click="() => fileInputs[index]?.click()"
              >
                {{ uploading[index] ? 'Загрузка...' : 'Загрузить' }}
              </ui-button>
            </div>
          </div>
        </div>
        
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-1">Подпись</label>
          <ui-input 
            v-model="image.caption" 
            variant="primary" 
            type="text" 
            placeholder="Подпись к фотографии"
          />
        </div>
      </div>
      
      <p v-if="formData.images.length === 0" class="text-sm text-gray-500 text-center py-4">
        Нет фотографий. Нажмите "+ Добавить фотографию" чтобы начать.
      </p>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue';
import UiInput from '../ui/Ui-input.vue';
import UiButton from '../ui/Ui-button.vue';
import { adminApi } from '../../api/admin';

export interface GalleryImage {
  src: string;
  caption: string;
}

export interface GalleryBlockContent {
  images: GalleryImage[];
}

const props = defineProps<{
  modelValue: GalleryBlockContent;
}>();

const emit = defineEmits<{
  'update:modelValue': [value: GalleryBlockContent];
}>();

const formData = ref<GalleryBlockContent>({
  images: props.modelValue?.images?.length ? [...props.modelValue.images] : [],
});

const fileInputs = ref<(HTMLInputElement | null)[]>(
  new Array(props.modelValue?.images?.length || 0).fill(null)
);
const uploading = ref<Record<number, boolean>>({});

const addImage = () => {
  formData.value.images.push({
    src: '',
    caption: '',
  });
  fileInputs.value.push(null);
};

const removeImage = (index: number) => {
  formData.value.images.splice(index, 1);
  fileInputs.value.splice(index, 1);
  emit('update:modelValue', { ...formData.value });
};

const handleFileUpload = async (index: number, event: Event) => {
  const target = event.target as HTMLInputElement;
  const file = target.files?.[0];
  if (!file) return;

  uploading.value[index] = true;
  try {
    const url = await adminApi.uploadFile(file);
    if (formData.value.images[index]) {
      formData.value.images[index].src = url;
      emit('update:modelValue', { ...formData.value });
    }
  } catch (error) {
    alert(error instanceof Error ? error.message : 'Ошибка загрузки файла');
  } finally {
    uploading.value[index] = false;
    // Сбрасываем input, чтобы можно было загрузить тот же файл снова
    if (target) {
      target.value = '';
    }
  }
};

const moveUp = (index: number) => {
  if (index > 0) {
    const current = formData.value.images[index];
    const prev = formData.value.images[index - 1];
    if (current && prev) {
      formData.value.images[index] = { ...prev };
      formData.value.images[index - 1] = { ...current };
      emit('update:modelValue', { ...formData.value });
    }
  }
};

const moveDown = (index: number) => {
  if (index < formData.value.images.length - 1) {
    const current = formData.value.images[index];
    const next = formData.value.images[index + 1];
    if (current && next) {
      formData.value.images[index] = { ...next };
      formData.value.images[index + 1] = { ...current };
      emit('update:modelValue', { ...formData.value });
    }
  }
};

watch(formData, (newVal) => {
  emit('update:modelValue', { ...newVal });
}, { deep: true });

watch(() => props.modelValue, (newVal) => {
  if (newVal && newVal.images) {
    formData.value = {
      images: [...newVal.images],
    };
    fileInputs.value = new Array(newVal.images.length).fill(null);
  }
}, { deep: true });
</script>

