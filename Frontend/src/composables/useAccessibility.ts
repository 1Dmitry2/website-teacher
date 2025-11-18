import { ref, watch, onMounted } from 'vue';

const ACCESSIBILITY_KEY = 'accessibility-mode';

export function useAccessibility() {
  const isAccessibilityMode = ref(false);

  const loadAccessibilityMode = () => {
    const saved = localStorage.getItem(ACCESSIBILITY_KEY);
    isAccessibilityMode.value = saved === 'true';
    applyAccessibilityMode(isAccessibilityMode.value);
  };

  const toggleAccessibilityMode = () => {
    isAccessibilityMode.value = !isAccessibilityMode.value;
    localStorage.setItem(ACCESSIBILITY_KEY, String(isAccessibilityMode.value));
    applyAccessibilityMode(isAccessibilityMode.value);
  };

  const applyAccessibilityMode = (enabled: boolean) => {
    const html = document.documentElement;
    if (enabled) {
      html.classList.add('accessibility-mode');
    } else {
      html.classList.remove('accessibility-mode');
    }
  };

  onMounted(() => {
    loadAccessibilityMode();
  });

  watch(isAccessibilityMode, (newValue) => {
    applyAccessibilityMode(newValue);
  });

  return {
    isAccessibilityMode,
    toggleAccessibilityMode,
  };
}

