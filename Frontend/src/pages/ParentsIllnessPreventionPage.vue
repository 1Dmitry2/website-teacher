<template>
  <div class="container mx-auto px-4 sm:px-6 lg:px-8 py-6 sm:py-8">
    <h1 class="text-2xl sm:text-3xl lg:text-4xl font-bold text-gray-800 mb-4 sm:mb-6">Профилактика заболеваний ОРЗ и ОРВИ</h1>
    
    <PageBlocks :has-additional-content="hasPosts" />
    
    <PostsList @has-content-change="handlePostsContentChange" />
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue';
import PageBlocks from '../components/PageBlocks.vue';
import PostsList from '../components/PostsList.vue';

const hasPosts = ref(false);
const handlePostsContentChange = (value: boolean) => {
  hasPosts.value = value;
};
</script>

<style lang="scss" scoped>
</style>

