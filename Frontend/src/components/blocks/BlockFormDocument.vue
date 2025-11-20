<template>
  <div class="space-y-4">
    <div>
      <label class="block text-sm font-medium text-gray-700 mb-1">Заголовок</label>
      <ui-input
        v-model="formData.title"
        variant="primary"
        type="text"
        placeholder="Например, Рабочая программа"
      />
    </div>

    <div>
      <label class="block text-sm font-medium text-gray-700 mb-1">Описание</label>
      <textarea
        v-model="formData.description"
        rows="3"
        class="w-full px-3 py-2 border border-gray-300 rounded-md focus:border-blue-500 focus:ring-2 focus:ring-blue-200 transition-all placeholder:text-gray-400"
        placeholder="Добавьте краткое описание содержимого документа"
      ></textarea>
    </div>

    <div>
      <label class="block text-sm font-medium text-gray-700 mb-1">Файл документа</label>
      <div class="flex flex-col gap-2">
        <div class="flex gap-2">
          <ui-input
            v-model="formData.fileUrl"
            variant="primary"
            type="text"
            placeholder="/uploads/document.pdf или загрузите файл"
            class="flex-1"
          />
          <div class="relative">
            <input
              type="file"
              ref="fileInput"
              accept=".pdf,.doc,.docx,.ppt,.pptx,.xls,.xlsx,.rtf,.odt,.txt,application/pdf,application/msword,application/vnd.openxmlformats-officedocument.wordprocessingml.document"
              @change="handleFileUpload"
              class="hidden"
              id="document-file-input"
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
        <p class="text-xs text-gray-500">
          Поддерживаются PDF, DOC/DOCX, PPT/PPTX, XLS/XLSX, RTF, ODT, TXT. Можно указать прямую ссылку или загрузить файл.
        </p>
      </div>
    </div>

    <div>
      <label class="block text-sm font-medium text-gray-700 mb-1">Название файла (для скачивания)</label>
      <ui-input
        v-model="formData.fileName"
        variant="primary"
        type="text"
        placeholder="Например, Рабочая программа 2025.pdf"
      />
    </div>

    <div>
      <label class="block text-sm font-medium text-gray-700 mb-1">Режим отображения</label>
      <select
        v-model="formData.mode"
        class="w-full px-3 py-2 border border-gray-300 rounded-md focus:border-blue-500 focus:ring-2 focus:ring-blue-200 transition-all"
      >
        <option value="download">Только кнопка скачивания</option>
        <option value="viewer">Встроенный просмотр</option>
      </select>
      <p class="text-xs text-gray-500 mt-1">
        Выберите, как будет отображаться документ на странице.
      </p>
    </div>

    <div v-if="formData.mode === 'viewer'" class="space-y-4 border border-gray-200 rounded-lg p-4 bg-gray-50">
      <div>
        <label class="block text-sm font-medium text-gray-700 mb-1">Вариант отображения</label>
        <select
          v-model="formData.viewerLayout"
          class="w-full px-3 py-2 border border-gray-300 rounded-md focus:border-blue-500 focus:ring-2 focus:ring-blue-200 transition-all"
        >
          <option value="full">На всю ширину страницы</option>
          <option value="preview">Превью + кнопка скачивания</option>
        </select>
      </div>

      <div v-if="formData.viewerLayout === 'preview'">
        <label class="block text-sm font-medium text-gray-700 mb-1">Высота превью (px)</label>
        <ui-input
          v-model.number="formData.previewHeight"
          variant="primary"
          type="number"
          min="200"
          max="1200"
          placeholder="500"
        />
        <p class="text-xs text-gray-500 mt-1">
          Высота области предпросмотра. По умолчанию 480px.
        </p>
      </div>

      <label class="inline-flex items-center gap-2 text-sm text-gray-700">
        <input
          type="checkbox"
          v-model="formData.showDownloadButton"
          class="rounded text-indigo-600 focus:ring-indigo-500 border-gray-300"
        />
        Показывать кнопку скачивания
      </label>
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

export interface DocumentBlockContent {
  title: string;
  description?: string;
  fileUrl: string;
  fileName?: string;
  mode: 'download' | 'viewer';
  viewerLayout?: 'full' | 'preview';
  previewHeight?: number;
  showDownloadButton?: boolean;
}

const props = defineProps<{
  modelValue: DocumentBlockContent;
}>();

const emit = defineEmits<{
  'update:modelValue': [value: DocumentBlockContent];
}>();

const fileInput = ref<HTMLInputElement | null>(null);
const uploading = ref(false);
const notification = useNotification();
const syncingFromModel = ref(false);

const formData = ref<DocumentBlockContent>({
  title: props.modelValue?.title || '',
  description: props.modelValue?.description || '',
  fileUrl: props.modelValue?.fileUrl || '',
  fileName: props.modelValue?.fileName || '',
  mode: props.modelValue?.mode || 'download',
  viewerLayout: props.modelValue?.viewerLayout || 'full',
  previewHeight: props.modelValue?.previewHeight || 480,
  showDownloadButton: props.modelValue?.showDownloadButton ?? true,
});

const handleFileUpload = async (event: Event) => {
  const target = event.target as HTMLInputElement;
  const file = target.files?.[0];
  if (!file) return;

  uploading.value = true;
  try {
    const url = await adminApi.uploadFile(file);
    formData.value.fileUrl = url;
    if (!formData.value.fileName) {
      formData.value.fileName = file.name;
    }
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
  formData.value = {
    title: newVal.title || '',
    description: newVal.description || '',
    fileUrl: newVal.fileUrl || '',
    fileName: newVal.fileName || '',
    mode: newVal.mode || 'download',
    viewerLayout: newVal.viewerLayout || 'full',
    previewHeight: newVal.previewHeight || 480,
    showDownloadButton: newVal.showDownloadButton ?? true,
  };
  nextTick(() => {
    syncingFromModel.value = false;
  });
}, { deep: true });
</script>

