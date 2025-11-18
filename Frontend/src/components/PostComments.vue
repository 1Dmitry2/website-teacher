<template>
  <div class="post-comments mt-6 border-t border-gray-200 pt-6">
    <h3 class="text-lg font-semibold text-gray-900 mb-4">
      Комментарии 
      <span v-if="comments.length > 0" class="text-gray-500 font-normal">({{ comments.length }})</span>
    </h3>

    <!-- Форма добавления комментария -->
    <div v-if="isAuthenticated" class="mb-6">
      <!-- Для администратора без авторизации пользователя - только ответы -->
      <div v-if="isAdmin && !isUserAuthenticated" class="p-4 bg-indigo-50 border border-indigo-200 rounded-lg">
        <p class="text-indigo-800 text-sm">
          <span class="font-semibold">Вы вошли как администратор.</span> Вы можете отвечать на комментарии, нажав кнопку "Ответить" под любым комментарием.
        </p>
      </div>
      <!-- Для обычных пользователей или администраторов с авторизацией пользователя -->
      <div v-else>
        <div class="flex gap-3">
          <div class="flex-1">
            <textarea
              v-model="newCommentText"
              placeholder="Напишите комментарий..."
              rows="3"
              class="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 resize-none"
              :disabled="submitting"
            ></textarea>
          </div>
        </div>
        <div class="flex justify-end mt-2">
          <button
            @click="submitComment"
            :disabled="!newCommentText.trim() || submitting"
            class="px-4 py-2 bg-indigo-600 text-white rounded-lg hover:bg-indigo-700 disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
          >
            {{ submitting ? 'Отправка...' : 'Отправить' }}
          </button>
        </div>
      </div>
    </div>
    <div v-else class="mb-6 p-4 bg-gray-50 rounded-lg text-center">
      <p class="text-gray-600 text-sm">
        <router-link to="/user-login" class="text-indigo-600 hover:text-indigo-700 font-medium">
          Войдите
        </router-link>
        , чтобы оставлять комментарии
      </p>
    </div>

    <!-- Список комментариев -->
    <div v-if="loading" class="text-center py-4 text-gray-500">
      Загрузка комментариев...
    </div>
    <div v-else-if="error" class="bg-red-50 border border-red-200 rounded-lg p-4 mb-4">
      <p class="text-red-700 text-sm">{{ error }}</p>
    </div>
    <div v-else-if="organizedComments.length === 0" class="text-center py-8 text-gray-500">
      <p>Пока нет комментариев. Будьте первым!</p>
    </div>
    <div v-else class="space-y-4">
      <div
        v-for="comment in organizedComments"
        :key="comment.id"
        class="comment-item"
      >
        <CommentItem
          :comment="comment"
          :replies="getReplies(comment.id)"
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
import { ref, computed, onMounted } from 'vue';
import { apiClient, type Comment } from '../api/client';
import { authService } from '../utils/auth';
import { adminAuthService } from '../utils/adminAuth';
import { adminApi } from '../api/admin';
import CommentItem from './CommentItem.vue';

const props = defineProps<{
  postId: string;
}>();

const comments = ref<Comment[]>([]);
const loading = ref(false);
const error = ref('');
const submitting = ref(false);
const newCommentText = ref('');

const isUserAuthenticated = computed(() => authService.isAuthenticated());
const isAdmin = computed(() => adminAuthService.isAuthenticated());
const isAuthenticated = computed(() => 
  isUserAuthenticated.value || isAdmin.value
);

// Организуем комментарии: родительские и ответы
const organizedComments = computed(() => {
  return comments.value.filter(comment => !comment.reply_to);
});

const getReplies = (parentId: string): Comment[] => {
  return comments.value.filter(comment => comment.reply_to === parentId);
};

const fetchComments = async () => {
  loading.value = true;
  error.value = '';
  try {
    comments.value = await apiClient.getPostComments(props.postId);
  } catch (err) {
    error.value = err instanceof Error ? err.message : 'Ошибка загрузки комментариев';
    console.error('Error fetching comments:', err);
  } finally {
    loading.value = false;
  }
};

const submitComment = async () => {
  if (!newCommentText.value.trim() || submitting.value) return;
  
  // Администраторы без авторизации пользователя не могут создавать комментарии
  // (форма не должна отображаться в этом случае, но на всякий случай проверяем)
  if (isAdmin.value && !isUserAuthenticated.value) {
    return;
  }
  
  submitting.value = true;
  error.value = '';
  try {
    const newComment = await apiClient.createComment(props.postId, newCommentText.value.trim());
    comments.value.push(newComment);
    newCommentText.value = '';
    await fetchComments(); // Обновляем для правильной сортировки
  } catch (err) {
    error.value = err instanceof Error ? err.message : 'Ошибка отправки комментария';
    if (error.value.toLowerCase().includes('unauthorized') || error.value.toLowerCase().includes('token')) {
      authService.removeToken();
      if (isAdmin.value) {
        adminAuthService.removeToken();
      }
      window.location.reload();
    }
  } finally {
    submitting.value = false;
  }
};

const handleReply = async (commentId: string, text: string) => {
  if (!text.trim() || submitting.value) return;
  
  submitting.value = true;
  error.value = '';
  try {
    // Если администратор, используем adminApi, иначе обычный apiClient
    const reply = isAdmin.value
      ? await adminApi.replyToComment(commentId, text.trim())
      : await apiClient.replyToComment(commentId, text.trim());
    comments.value.push(reply);
    await fetchComments(); // Обновляем для правильной сортировки
  } catch (err) {
    error.value = err instanceof Error ? err.message : 'Ошибка отправки ответа';
    if (error.value.toLowerCase().includes('unauthorized') || error.value.toLowerCase().includes('token')) {
      if (isAdmin.value) {
        adminAuthService.removeToken();
      } else {
        authService.removeToken();
      }
      window.location.reload();
    }
  } finally {
    submitting.value = false;
  }
};

const handleDelete = async (commentId: string) => {
  try {
    await apiClient.deleteComment(commentId);
    // Удаляем комментарий и все его ответы
    comments.value = comments.value.filter(c => c.id !== commentId && c.reply_to !== commentId);
  } catch (err) {
    const errorMsg = err instanceof Error ? err.message : 'Ошибка удаления комментария';
    error.value = errorMsg;
    // Если ошибка авторизации, перезагружаем страницу
    if (errorMsg.toLowerCase().includes('unauthorized') || errorMsg.toLowerCase().includes('forbidden') || errorMsg.toLowerCase().includes('token')) {
      authService.removeToken();
      setTimeout(() => window.location.reload(), 1000);
    }
  }
};

onMounted(() => {
  fetchComments();
});
</script>

<style scoped>
.post-comments {
  @apply w-full;
}
</style>

