<template>
  <div class="comment-item">
    <div 
      class="bg-gray-50 rounded-lg p-4"
      :class="{ 'border-l-4 border-indigo-500': comment.is_admin }"
    >
      <div class="flex items-start justify-between mb-2">
        <div class="flex items-center gap-2">
          <div 
            class="w-8 h-8 rounded-full flex items-center justify-center text-white text-sm font-semibold"
            :class="comment.is_admin ? 'bg-indigo-600' : 'bg-gray-400'"
          >
            {{ comment.is_admin ? 'А' : 'П' }}
          </div>
          <div>
            <div class="flex items-center gap-2">
              <span 
                class="font-semibold text-sm"
                :class="comment.is_admin ? 'text-indigo-600' : 'text-gray-900'"
              >
                {{ comment.is_admin ? 'Администратор' : 'Пользователь' }}
              </span>
              <span 
                v-if="comment.is_admin" 
                class="px-2 py-0.5 bg-indigo-100 text-indigo-700 text-xs font-medium rounded-full"
              >
                Ответ администратора
              </span>
            </div>
            <span class="text-gray-500 text-xs">{{ formatDate(comment.created_at) }}</span>
          </div>
        </div>
      </div>
      
      <p class="text-gray-800 mb-3 whitespace-pre-wrap">{{ comment.text }}</p>
      
      <div class="flex items-center justify-between">
        <div class="flex items-center gap-3">
          <button
            v-if="isAuthenticated && !showReplyForm"
            @click="showReplyForm = true"
            class="text-sm text-indigo-600 hover:text-indigo-700 font-medium"
          >
            Ответить
          </button>
          <button
            v-if="isAuthenticated && showReplyForm"
            @click="showReplyForm = false"
            class="text-sm text-gray-500 hover:text-gray-700"
          >
            Отмена
          </button>
        </div>
        <button
          v-if="isAuthenticated && !comment.is_admin && !isAdmin"
          @click="handleDeleteClick"
          class="text-sm text-red-600 hover:text-red-700 font-medium"
        >
          Удалить
        </button>
      </div>

      <!-- Форма ответа -->
      <div v-if="showReplyForm" class="mt-3 pt-3 border-t border-gray-200">
        <textarea
          v-model="replyText"
          placeholder="Напишите ответ..."
          rows="2"
          class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 resize-none text-sm"
        ></textarea>
        <div class="flex justify-end gap-2 mt-2">
          <button
            @click="showReplyForm = false"
            class="px-3 py-1 text-sm text-gray-600 hover:text-gray-800"
          >
            Отмена
          </button>
          <button
            @click="submitReply"
            :disabled="!replyText.trim()"
            class="px-3 py-1 text-sm bg-indigo-600 text-white rounded-lg hover:bg-indigo-700 disabled:opacity-50 disabled:cursor-not-allowed"
          >
            Отправить
          </button>
        </div>
      </div>
    </div>

    <!-- Ответы на комментарий -->
    <div v-if="replies.length > 0" class="ml-8 mt-3 space-y-3">
      <div
        v-for="reply in replies"
        :key="reply.id"
        class="relative"
      >
        <!-- Индикатор ответа администратора -->
        <div v-if="reply.is_admin" class="absolute -left-2 top-0 bottom-0 w-1 bg-indigo-500 rounded-l"></div>
        <CommentItem
          :comment="reply"
          :replies="[]"
          :is-authenticated="isAuthenticated"
          :is-admin="isAdmin"
          @reply="handleReply"
          @delete="handleDelete"
        />
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue';
import { type Comment } from '../api/client';

const props = defineProps<{
  comment: Comment;
  replies: Comment[];
  isAuthenticated: boolean;
  isAdmin?: boolean;
}>();

const emit = defineEmits<{
  reply: [commentId: string, text: string];
  delete: [commentId: string];
}>();

const showReplyForm = ref(false);
const replyText = ref('');

const formatDate = (date: string) => {
  return new Date(date).toLocaleDateString('ru-RU', {
    year: 'numeric',
    month: 'long',
    day: 'numeric',
    hour: '2-digit',
    minute: '2-digit'
  });
};

const submitReply = () => {
  if (!replyText.value.trim()) return;
  emit('reply', props.comment.id, replyText.value.trim());
  replyText.value = '';
  showReplyForm.value = false;
};

const handleReply = (commentId: string, text: string) => {
  emit('reply', commentId, text);
};

const handleDelete = (commentId: string) => {
  emit('delete', commentId);
};

const handleDeleteClick = () => {
  if (confirm('Удалить этот комментарий?')) {
    emit('delete', props.comment.id);
  }
};
</script>

<style scoped>
.comment-item {
  @apply w-full;
}
</style>

