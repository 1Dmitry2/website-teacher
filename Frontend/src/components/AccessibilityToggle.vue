<template>
  <button
    @click="toggleAccessibilityMode"
    class="accessibility-toggle"
    :class="{ 'active': isAccessibilityMode }"
    :aria-label="isAccessibilityMode ? 'Перейти в обычную версию' : 'Версия для слабовидящих'"
    title="Версия для слабовидящих"
  >
    <svg v-if="!isAccessibilityMode" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="icon">
      <path stroke-linecap="round" stroke-linejoin="round" d="M2.036 12.322a1.012 1.012 0 0 1 0-.639C3.423 7.51 7.36 4.5 12 4.5c4.638 0 8.573 3.007 9.963 7.178.07.207.07.431 0 .639C20.577 16.49 16.64 19.5 12 19.5c-4.638 0-8.573-3.007-9.963-7.178Z" />
      <path stroke-linecap="round" stroke-linejoin="round" d="M15 12a3 3 0 1 1-6 0 3 3 0 0 1 6 0Z" />
    </svg>
    <svg v-else xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="icon">
      <path stroke-linecap="round" stroke-linejoin="round" d="M3.98 8.223A10.477 10.477 0 0 0 1.934 12C3.226 16.338 7.244 19.5 12 19.5c.993 0 1.953-.138 2.863-.395M6.228 6.228A10.451 10.451 0 0 1 12 4.5c4.756 0 8.773 3.162 10.065 7.498a10.522 10.522 0 0 1-4.293 5.774M6.228 6.228 3 3m3.228 3.228 3.65 3.65m7.894 7.894L21 21m-3.228-3.228-3.65-3.65m0 0a3 3 0 1 0-4.243-4.243m4.242 4.242L9.88 9.88" />
    </svg>
    <span class="text">{{ isAccessibilityMode ? 'Обычная версия' : 'Версия для слабовидящих' }}</span>
  </button>
</template>

<script setup lang="ts">
import { useAccessibility } from '../composables/useAccessibility';

const { isAccessibilityMode, toggleAccessibilityMode } = useAccessibility();
</script>

<style scoped>
.accessibility-toggle {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  padding: 0.625rem 1rem;
  background: #ffffff;
  border: 2px solid #4f46e5;
  border-radius: 0.5rem;
  color: #4f46e5;
  font-size: 0.875rem;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
}

.accessibility-toggle:hover {
  background: #f3f4f6;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.15);
}

.accessibility-toggle.active {
  background: #4f46e5;
  color: #ffffff;
}

.accessibility-toggle.active:hover {
  background: #4338ca;
}

.accessibility-toggle .icon {
  width: 1.25rem;
  height: 1.25rem;
  flex-shrink: 0;
}

.accessibility-toggle .text {
  white-space: nowrap;
}

.accessibility-mode .accessibility-toggle {
  font-size: 1rem;
  padding: 0.875rem 1.25rem;
  border-width: 3px;
}

.accessibility-mode .accessibility-toggle .icon {
  width: 1.5rem;
  height: 1.5rem;
}

/* Мобильная версия: только иконка на экранах меньше 400px */
@media (max-width: 475px) {
  .accessibility-toggle {
    padding: 0.75rem;
    min-width: 3rem;
    min-height: 3rem;
    justify-content: center;
  }

  .accessibility-toggle .text {
    display: none;
  }

  .accessibility-toggle .icon {
    width: 1.75rem;
    height: 1.75rem;
  }

  .accessibility-mode .accessibility-toggle {
    padding: 0.875rem;
    min-width: 3.5rem;
    min-height: 3.5rem;
  }

  .accessibility-mode .accessibility-toggle .icon {
    width: 2rem;
    height: 2rem;
  }
}
</style>

