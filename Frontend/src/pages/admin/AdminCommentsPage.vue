<template>
  <div class="max-w-7xl mx-auto px-4 py-10 space-y-6">
    <div>
      <h1 class="text-3xl font-bold text-gray-900">Комментарии</h1>
      <p class="text-gray-500 mt-1">Модерация комментариев</p>
    </div>

    <div v-if="loading" class="bg-white rounded-xl shadow p-6 text-center text-gray-500">Загрузка...</div>
    <div v-else-if="error" class="bg-red-50 border border-red-200 rounded-xl p-6"><p class="text-red-700">{{ error }}</p></div>
    <div v-else class="space-y-4">
      <div v-for="comment in comments" :key="comment.id" class="bg-white rounded-xl shadow p-6">
        <div class="flex items-start justify-between">
          <div class="flex-1">
            <div class="flex items-center gap-2 mb-2">
              <span class="text-sm font-medium text-gray-900">ID поста: {{ comment.post_id }}</span>
              <span v-if="comment.is_admin" class="px-2 py-1 text-xs font-semibold bg-indigo-100 text-indigo-800 rounded-full">Админ</span>
              <span class="text-xs text-gray-500">{{ formatDate(comment.created_at) }}</span>
            </div>
            <p class="text-gray-700">{{ comment.text }}</p>
            <div v-if="comment.reply_to" class="mt-2 text-xs text-gray-500">Ответ на комментарий: {{ comment.reply_to }}</div>
          </div>
          <div class="flex gap-2 ml-4">
            <button @click="openReplyModal(comment)" class="text-sm text-indigo-600 hover:text-indigo-900">Ответить</button>
            <button @click="deleteComment(comment.id)" class="text-sm text-red-600 hover:text-red-900">Удалить</button>
          </div>
        </div>
      </div>
    </div>

    <Transition name="modal">
      <div v-if="isReplyModalOpen" class="fixed inset-0 bg-gray-900/60 backdrop-blur-sm flex items-center justify-center z-50 p-4" @click.self="closeReplyModal">
        <div class="bg-white rounded-xl shadow-2xl max-w-md w-full max-h-[90vh] flex flex-col transform transition-all">
          <!-- Заголовок (фиксированный) -->
          <div class="px-6 pt-6 pb-4 border-b border-gray-200 relative flex-shrink-0">
            <button
              @click="closeReplyModal"
              class="absolute top-6 right-6 text-gray-400 hover:text-gray-600 transition-colors focus:outline-none focus:ring-2 focus:ring-indigo-500 rounded-full p-1"
              type="button"
              aria-label="Закрыть"
            >
              <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
              </svg>
            </button>
            <h2 class="text-2xl font-bold text-gray-900 pr-8">Ответить на комментарий</h2>
          </div>
          <!-- Контент (прокручиваемый) -->
          <form @submit.prevent="submitReply" class="flex-1 overflow-y-auto px-6 py-4 space-y-4">
            <div><label class="block text-sm font-medium text-gray-700 mb-1">Ответ</label><textarea v-model="replyText" rows="4" class="w-full px-3 py-2 border border-gray-300 rounded-md" required></textarea></div>
          </form>
          <!-- Кнопки (фиксированные) -->
          <div class="px-6 py-4 border-t border-gray-200 flex gap-2 justify-end flex-shrink-0">
            <ui-button type="button" variant="secondary" @click="closeReplyModal">Отмена</ui-button>
            <ui-button type="button" variant="primary" @click="submitReply">Отправить</ui-button>
          </div>
        </div>
      </div>
    </Transition>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue';
import { useRouter } from 'vue-router';
import UiButton from '../../components/ui/Ui-button.vue';
import { adminApi, type Comment } from '../../api/admin';
import { adminAuthService } from '../../utils/adminAuth';

const router = useRouter();
const comments = ref<Comment[]>([]);
const loading = ref(true);
const error = ref('');
const isReplyModalOpen = ref(false);
const replyingTo = ref<Comment | null>(null);
const replyText = ref('');

const formatDate = (date: string) => {
  return new Date(date).toLocaleString('ru-RU');
};

const fetchComments = async () => {
  loading.value = true;
  error.value = '';
  try {
    comments.value = await adminApi.getComments();
  } catch (err) {
    error.value = err instanceof Error ? err.message : 'Ошибка загрузки';
    if (error.value.toLowerCase().includes('unauthorized') || error.value.toLowerCase().includes('token')) {
      adminAuthService.removeToken();
      router.push('/admin-login');
    }
  } finally {
    loading.value = false;
  }
};

const openReplyModal = (comment: Comment) => {
  replyingTo.value = comment;
  isReplyModalOpen.value = true;
};

const closeReplyModal = () => {
  isReplyModalOpen.value = false;
  replyingTo.value = null;
  replyText.value = '';
};

const submitReply = async () => {
  if (!replyingTo.value) return;
  try {
    await adminApi.replyToComment(replyingTo.value.id, replyText.value);
    closeReplyModal();
    await fetchComments();
  } catch (err) {
    alert(err instanceof Error ? err.message : 'Ошибка отправки ответа');
  }
};

const deleteComment = async (id: string) => {
  if (!confirm('Удалить этот комментарий?')) return;
  try {
    await adminApi.deleteComment(id);
    await fetchComments();
  } catch (err) {
    alert(err instanceof Error ? err.message : 'Ошибка удаления');
  }
};

onMounted(fetchComments);
</script>

