<template>
  <Transition name="modal">
    <div v-if="isOpen" class="fixed inset-0 bg-gray-900/60 backdrop-blur-sm flex items-center justify-center z-50 p-4" @click.self="handleCancel">
      <div class="bg-white rounded-xl shadow-2xl max-w-md w-full transform transition-all">
        <div class="px-6 pt-6 pb-4">
          <h2 class="text-2xl font-bold text-gray-900 mb-2">{{ title }}</h2>
          <p class="text-gray-600">{{ message }}</p>
        </div>
        
        <div class="px-6 pb-6 flex gap-3 justify-end">
          <button
            @click="handleCancel"
            class="px-4 py-2 text-sm font-medium text-gray-700 bg-gray-100 hover:bg-gray-200 rounded-lg transition-colors focus:outline-none focus:ring-2 focus:ring-gray-400"
          >
            Отмена
          </button>
          <button
            @click="handleConfirm"
            :disabled="loading"
            class="px-4 py-2 text-sm font-medium text-white rounded-lg transition-colors focus:outline-none focus:ring-2 focus:ring-offset-2 disabled:opacity-50"
            :class="confirmVariant === 'danger' 
              ? 'bg-red-600 hover:bg-red-700 focus:ring-red-500' 
              : 'bg-blue-600 hover:bg-blue-700 focus:ring-blue-500'"
          >
            {{ loading ? 'Обработка...' : confirmText }}
          </button>
        </div>
      </div>
    </div>
  </Transition>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue';

interface Props {
  isOpen: boolean;
  title?: string;
  message: string;
  confirmText?: string;
  cancelText?: string;
  confirmVariant?: 'default' | 'danger';
  loading?: boolean;
}

const props = withDefaults(defineProps<Props>(), {
  title: 'Подтверждение',
  confirmText: 'Подтвердить',
  cancelText: 'Отмена',
  confirmVariant: 'default',
  loading: false,
});

const emit = defineEmits<{
  confirm: [];
  cancel: [];
}>();

const handleConfirm = () => {
  emit('confirm');
};

const handleCancel = () => {
  emit('cancel');
};
</script>

<style scoped>
.modal-enter-active {
  transition: opacity 0.3s ease;
}

.modal-leave-active {
  transition: opacity 0.2s ease;
}

.modal-enter-from,
.modal-leave-to {
  opacity: 0;
}

.modal-enter-active > div,
.modal-leave-active > div {
  transition: transform 0.3s ease, opacity 0.3s ease;
}

.modal-enter-from > div,
.modal-leave-to > div {
  transform: scale(0.95) translateY(-10px);
  opacity: 0;
}
</style>

