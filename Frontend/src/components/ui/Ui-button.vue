<template>
  <button
      :type="props.type"
      class="cursor-pointer"
      :class="{
        'text-white bg-blue-600 hover:bg-blue-700 focus:ring-4 focus:ring-blue-400 shadow-sm font-medium rounded-md text-sm px-4 py-2.5 focus:outline-none transition-colors' : props.variant == 'primary',
        'text-gray-200 bg-gray-600 border border-gray-500 hover:bg-gray-700 hover:text-white focus:ring-4 focus:ring-gray-400 shadow-sm font-medium rounded-md text-sm px-4 py-2.5 focus:outline-none transition-colors' : props.variant == 'secondary',
        'text-white bg-green-600 hover:bg-green-700 focus:ring-4 focus:ring-green-400 shadow-sm font-medium rounded-md text-sm px-4 py-2.5 focus:outline-none transition-colors' : props.variant == 'success',
        'text-white bg-red-600 hover:bg-red-700 focus:ring-4 focus:ring-red-400 shadow-sm font-medium rounded-md text-sm px-4 py-2.5 focus:outline-none transition-colors' : props.variant == 'danger',
        'text-blue-600 bg-transparent hover:text-blue-700 hover:underline focus:ring-2 focus:ring-blue-400 font-medium text-sm px-2 py-1 focus:outline-none transition-colors' : props.variant == 'link',
      }"
      :disabled="props.disabled"
  >
    <slot></slot>
  </button>
</template>

<script setup lang="ts">
interface UIButton {
  variant: 'primary' | 'secondary' | 'success' | 'danger' | 'link',
  type: 'button' | 'submit' | 'reset',
  disabled?: boolean,
}
const props = defineProps<UIButton>();

</script>
