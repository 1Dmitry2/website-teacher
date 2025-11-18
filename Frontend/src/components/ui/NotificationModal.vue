<template>
  <Transition name="modal">
    <div v-if="isOpen" class="fixed inset-0 bg-gray-900/60 backdrop-blur-sm flex items-center justify-center z-50 p-4" @click.self="handleClose">
      <div class="bg-white rounded-xl shadow-2xl max-w-md w-full transform transition-all">
        <div class="px-6 pt-6 pb-4">
          <div class="flex items-start gap-4">
            <div 
              class="flex-shrink-0 w-10 h-10 rounded-full flex items-center justify-center"
              :class="iconClass"
            >
              <svg v-if="type === 'error'" class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
              </svg>
              <svg v-else-if="type === 'success'" class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
              </svg>
              <svg v-else class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
              </svg>
            </div>
            <div class="flex-1">
              <h2 class="text-xl font-bold mb-2" :class="titleClass">{{ title }}</h2>
              <p class="text-gray-600">{{ message }}</p>
            </div>
            <button
              @click="handleClose"
              class="flex-shrink-0 text-gray-400 hover:text-gray-600 transition-colors focus:outline-none focus:ring-2 focus:ring-gray-400 rounded-full p-1"
              type="button"
              aria-label="Закрыть"
            >
              <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
              </svg>
            </button>
          </div>
        </div>
        
        <div class="px-6 pb-6 flex justify-end">
          <button
            @click="handleClose"
            class="px-4 py-2 text-sm font-medium text-white rounded-lg transition-colors focus:outline-none focus:ring-2 focus:ring-offset-2"
            :class="buttonClass"
          >
            ОК
          </button>
        </div>
      </div>
    </div>
  </Transition>
</template>

<script setup lang="ts">
import { computed } from 'vue';

interface Props {
  isOpen: boolean;
  type?: 'error' | 'success' | 'info';
  title?: string;
  message: string;
}

const props = withDefaults(defineProps<Props>(), {
  type: 'info',
  title: '',
});

const emit = defineEmits<{
  close: [];
}>();

const iconClass = computed(() => {
  switch (props.type) {
    case 'error':
      return 'bg-red-100 text-red-600';
    case 'success':
      return 'bg-green-100 text-green-600';
    default:
      return 'bg-blue-100 text-blue-600';
  }
});

const titleClass = computed(() => {
  switch (props.type) {
    case 'error':
      return 'text-red-900';
    case 'success':
      return 'text-green-900';
    default:
      return 'text-gray-900';
  }
});

const buttonClass = computed(() => {
  switch (props.type) {
    case 'error':
      return 'bg-red-600 hover:bg-red-700 focus:ring-red-500';
    case 'success':
      return 'bg-green-600 hover:bg-green-700 focus:ring-green-500';
    default:
      return 'bg-blue-600 hover:bg-blue-700 focus:ring-blue-500';
  }
});

const handleClose = () => {
  emit('close');
};

// Автоматическое закрытие через 5 секунд для info и success
if (props.type !== 'error') {
  // Это будет обработано в родительском компоненте через composable
}
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

