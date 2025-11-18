import { ref } from 'vue';

interface ConfirmationOptions {
  title?: string;
  message: string;
  confirmText?: string;
  cancelText?: string;
  confirmVariant?: 'default' | 'danger';
}

interface NotificationOptions {
  type?: 'error' | 'success' | 'info';
  title?: string;
  message: string;
}

export function useConfirmation() {
  const isOpen = ref(false);
  const loading = ref(false);
  const options = ref<ConfirmationOptions>({ message: '' });
  let resolvePromise: ((value: boolean) => void) | null = null;

  const confirm = (opts: ConfirmationOptions): Promise<boolean> => {
    return new Promise((resolve) => {
      options.value = opts;
      isOpen.value = true;
      resolvePromise = resolve;
    });
  };

  const handleConfirm = () => {
    if (resolvePromise) {
      resolvePromise(true);
      resolvePromise = null;
    }
    isOpen.value = false;
  };

  const handleCancel = () => {
    if (resolvePromise) {
      resolvePromise(false);
      resolvePromise = null;
    }
    isOpen.value = false;
  };

  return {
    isOpen,
    loading,
    options,
    confirm,
    handleConfirm,
    handleCancel,
  };
}

export function useNotification() {
  const isOpen = ref(false);
  const options = ref<NotificationOptions>({ message: '' });
  let timeoutId: ReturnType<typeof setTimeout> | null = null;

  const show = (opts: NotificationOptions) => {
    options.value = opts;
    isOpen.value = true;

    // Автоматически закрываем через 5 секунд для info и success
    if (opts.type !== 'error') {
      if (timeoutId) clearTimeout(timeoutId);
      timeoutId = setTimeout(() => {
        isOpen.value = false;
      }, 5000);
    }
  };

  const close = () => {
    if (timeoutId) {
      clearTimeout(timeoutId);
      timeoutId = null;
    }
    isOpen.value = false;
  };

  // Удобные методы для разных типов уведомлений
  const error = (message: string, title?: string) => {
    show({ type: 'error', message, title: title || 'Ошибка' });
  };

  const success = (message: string, title?: string) => {
    show({ type: 'success', message, title: title || 'Успешно' });
  };

  const info = (message: string, title?: string) => {
    show({ type: 'info', message, title: title || 'Информация' });
  };

  return {
    isOpen,
    options,
    show,
    close,
    error,
    success,
    info,
  };
}

