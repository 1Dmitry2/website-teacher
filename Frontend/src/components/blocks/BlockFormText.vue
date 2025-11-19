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
      <label class="block text-sm font-medium text-gray-700 mb-1">Выравнивание</label>
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
  </div>
</template>

<script setup lang="ts">
import { nextTick, ref, watch } from 'vue';
import UiInput from '../ui/Ui-input.vue';

export interface TextBlockContent {
  title: string;
  text: string;
  alignment: 'left' | 'center' | 'right' | 'justify';
  style: 'regular' | 'bold' | 'italic' | 'heading';
}

const props = defineProps<{
  modelValue: TextBlockContent;
}>();

const emit = defineEmits<{
  'update:modelValue': [value: TextBlockContent];
}>();

const formData = ref<TextBlockContent>({
  title: props.modelValue?.title || '',
  text: props.modelValue?.text || '',
  alignment: props.modelValue?.alignment || 'left',
  style: props.modelValue?.style || 'regular',
});

const syncingFromModel = ref(false);

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

