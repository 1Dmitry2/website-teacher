<template>
  <div class="page-selector">
    <div class="border border-gray-300 rounded-md max-h-64 overflow-y-auto">
      <div v-for="route in availableRoutes" :key="route.path" class="px-3 py-2 hover:bg-gray-50">
        <label class="flex items-start cursor-pointer gap-2">
          <input
            type="checkbox"
            :value="route.path"
            :checked="selectedPages.includes(route.path)"
            @change="handleChange"
            class="w-4 h-4 text-indigo-600 border-gray-300 rounded focus:ring-indigo-500"
          />
          <span class="flex flex-col">
            <span class="text-sm font-medium text-gray-900">
              {{ route.displayName }}
            </span>
            <span class="text-xs text-gray-500">
              {{ route.path }}
            </span>
          </span>
        </label>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
interface Route {
  path: string;
  displayName: string;
}

interface Props {
  selectedPages: string[];
  availableRoutes: Route[];
}

interface Emits {
  (e: 'update:selectedPages', value: string[]): void;
}

const props = defineProps<Props>();
const emit = defineEmits<Emits>();

const handleChange = (event: Event) => {
  const target = event.target as HTMLInputElement;
  const value = target.value;
  const checked = target.checked;
  
  let newSelected = [...props.selectedPages];
  if (checked) {
    if (!newSelected.includes(value)) {
      newSelected.push(value);
    }
  } else {
    newSelected = newSelected.filter(p => p !== value);
  }
  
  emit('update:selectedPages', newSelected);
};
</script>

<style scoped>
.page-selector {
  @apply w-full;
}
</style>

