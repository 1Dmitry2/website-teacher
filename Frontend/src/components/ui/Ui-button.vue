<template>
  <button
      :type="type"
      :class="{'primary': variant === 'primary'}"
      class="cursor-pointer">
    <slot></slot>
  </button>
</template>

<script setup lang="ts">
interface IButton {
  variant: {
    type: 'primary' | 'secondary' | 'success' | 'danger'
  }
}
const props = defineProps<({
  variant:{
    type: variant.type,
    default: 'primary'
  },
  type: {
    type: String,
    default: 'button'
  }
})
</script>
