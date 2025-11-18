<template>
  <div class="space-y-4">
    <div class="space-y-3">
      <div class="flex items-center justify-between">
        <label class="block text-sm font-medium text-gray-700">Слайды</label>
        <ui-button type="button" variant="primary" @click="addSlide">
          + Добавить слайд
        </ui-button>
      </div>
      
      <div v-for="(slide, index) in formData.slides" :key="index" class="border border-gray-300 rounded-lg p-4 space-y-3">
        <div class="flex items-center justify-between">
          <span class="text-sm font-medium text-gray-700">Слайд {{ index + 1 }}</span>
          <button 
            type="button"
            @click="removeSlide(index)"
            class="text-red-600 hover:text-red-800 text-sm"
          >
            Удалить
          </button>
        </div>
        
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-1">Изображение</label>
          <div class="flex gap-2">
            <ui-input 
              v-model="slide.image" 
              variant="primary" 
              type="text" 
              placeholder="/upload/img1.jpg или загрузите файл"
              class="flex-1"
            />
            <div class="relative">
              <input
                type="file"
                :ref="el => fileInputs[index] = el as HTMLInputElement"
                accept="image/*"
                @change="handleFileUpload(index, $event)"
                class="hidden"
                :id="`slider-file-${index}`"
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
          <label class="block text-sm font-medium text-gray-700 mb-1">Заголовок (необязательно)</label>
          <ui-input 
            :model-value="slide.title || ''"
            @update:model-value="slide.title = String($event)"
            variant="primary" 
            type="text" 
            placeholder="Заголовок слайда"
          />
        </div>
        
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-1">Подзаголовок (необязательно)</label>
          <ui-input 
            :model-value="slide.subtitle || ''"
            @update:model-value="slide.subtitle = String($event)"
            variant="primary" 
            type="text" 
            placeholder="Подзаголовок слайда"
          />
        </div>
      </div>
      
      <p v-if="formData.slides.length === 0" class="text-sm text-gray-500 text-center py-4">
        Нет слайдов. Нажмите "+ Добавить слайд" чтобы начать.
      </p>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue';
import UiInput from '../ui/Ui-input.vue';
import UiButton from '../ui/Ui-button.vue';
import { adminApi } from '../../api/admin';

export interface SliderSlide {
  image: string;
  title?: string;
  subtitle?: string;
}

export interface SliderBlockContent {
  slides: SliderSlide[];
}

const props = defineProps<{
  modelValue: SliderBlockContent;
}>();

const emit = defineEmits<{
  'update:modelValue': [value: SliderBlockContent];
}>();

const formData = ref<SliderBlockContent>({
  slides: props.modelValue?.slides?.length ? [...props.modelValue.slides] : [],
});

const fileInputs = ref<(HTMLInputElement | null)[]>(
  new Array(props.modelValue?.slides?.length || 0).fill(null)
);
const uploading = ref<Record<number, boolean>>({});

const addSlide = () => {
  formData.value.slides.push({
    image: '',
    title: '',
    subtitle: '',
  });
  fileInputs.value.push(null);
};

const removeSlide = (index: number) => {
  formData.value.slides.splice(index, 1);
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
    if (formData.value.slides[index]) {
      formData.value.slides[index].image = url;
      emit('update:modelValue', { ...formData.value });
    }
  } catch (error) {
    alert(error instanceof Error ? error.message : 'Ошибка загрузки файла');
  } finally {
    uploading.value[index] = false;
    if (target) {
      target.value = '';
    }
  }
};

watch(formData, (newVal) => {
  emit('update:modelValue', { ...newVal });
}, { deep: true });

watch(() => props.modelValue, (newVal) => {
  if (newVal && newVal.slides) {
    formData.value = {
      slides: [...newVal.slides],
    };
    fileInputs.value = new Array(newVal.slides.length).fill(null);
  }
}, { deep: true });
</script>

