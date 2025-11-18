<template>
  <div class="w-full">
    <input
        class="w-full outline-none text-sm sm:text-base"
        :class="{
      'text-gray-900 bg-gray-50 border border-gray-300 focus:border-blue-500 focus:ring-2 focus:ring-blue-200 rounded-md px-3 py-2 sm:px-4 sm:py-2.5 transition-all placeholder:text-gray-400' : props.variant === 'primary',
      'text-gray-900 bg-white border border-gray-400 focus:border-blue-600 focus:ring-2 focus:ring-blue-300 rounded-md px-3 py-2 sm:px-4 sm:py-2.5 transition-all placeholder:text-gray-500' : props.variant === 'secondary',
      'border border-red-500 bg-red-50 text-red-700 placeholder-red-400 focus:ring-2 focus:ring-red-200 focus:border-red-600 rounded-md px-3 py-2 sm:px-4 sm:py-2.5 transition-all': props.variant === 'error',
      'border border-green-500 bg-green-50 text-green-700 placeholder-green-400 focus:ring-2 focus:ring-green-200 focus:border-green-600 rounded-md px-3 py-2 sm:px-4 sm:py-2.5 transition-all': props.variant === 'success',
        }"
        :value="modelValue"
        @input="onInput"
        :type="type"
        :placeholder="placeholder"
        :disabled="disabled"
        :min="$attrs.min"
        :required="$attrs.required"
    />
  </div>
</template>

<script setup lang="ts">
interface UIInput{
  type: 'text' | 'password' | 'email' | 'search' | 'number',
  placeholder?: string,
  disabled?: boolean,
  variant: 'primary' | 'secondary' | 'error' | 'success',
  modelValue: string | number

}

const props = defineProps<UIInput>();
const emit = defineEmits<{
  'update:modelValue': [value: string | number];
}>();
const onInput = (e: Event) => {
  const target = e.target as HTMLInputElement;
  const value = props.type === 'number' ? (target.valueAsNumber || 0) : target.value;
  emit('update:modelValue', value);
}
</script>
