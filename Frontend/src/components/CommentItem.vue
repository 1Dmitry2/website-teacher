<template>
  <div class="comment-item">
    <div 
      class="bg-gray-50 rounded-lg p-4"
      :class="{ 'border-l-4 border-indigo-500': comment.is_admin }"
    >
      <div class="flex items-start justify-between mb-2 gap-2">
        <div class="flex items-center gap-2 min-w-0 flex-1">
          <button
            @click="handleUserClick"
            class="w-7 h-7 sm:w-8 sm:h-8 rounded-full flex items-center justify-center text-white text-xs sm:text-sm font-semibold flex-shrink-0 cursor-pointer hover:opacity-80 transition-opacity"
            :class="comment.is_admin ? 'bg-indigo-600' : 'bg-gray-400'"
            :title="comment.user_email || 'Пользователь'"
          >
            {{ getInitial() }}
          </button>
          <div class="min-w-0 flex-1">
            <div class="flex items-center gap-1.5 sm:gap-2 flex-wrap">
              <span 
                class="font-semibold text-xs sm:text-sm truncate"
                :class="comment.is_admin ? 'text-indigo-600' : 'text-gray-900'"
              >
                {{ comment.is_admin ? 'Администратор' : 'Пользователь' }}
              </span>
              <span 
                v-if="comment.is_admin" 
                class="px-1.5 sm:px-2 py-0.5 bg-indigo-100 text-indigo-700 text-[10px] sm:text-xs font-medium rounded-full whitespace-nowrap"
              >
                Ответ администратора
              </span>
            </div>
            <span class="text-gray-500 text-[10px] sm:text-xs">{{ formatDate(comment.created_at) }}</span>
          </div>
        </div>
      </div>
      
      <p class="text-sm sm:text-base text-gray-800 mb-3 whitespace-pre-wrap break-words">{{ comment.text }}</p>
      
      <div v-if="hasReplies" class="mb-3 pb-2 border-b border-gray-200">
        <button
          @click="isCommentExpanded = !isCommentExpanded"
          class="text-xs sm:text-sm text-indigo-600 hover:text-indigo-700 font-medium flex items-center gap-1.5 transition-colors"
        >
          <svg 
            class="w-4 h-4 transition-transform duration-300 ease-in-out"
            :class="{ 'rotate-180': !isCommentExpanded }"
            fill="none" 
            stroke="currentColor" 
            viewBox="0 0 24 24"
          >
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7" />
          </svg>
          <Transition name="fade" mode="out-in">
            <span :key="isCommentExpanded">
              <template v-if="!isCommentExpanded">
                Показать {{ replies.length }} {{ replies.length === 1 ? 'ответ' : replies.length < 5 ? 'ответа' : 'ответов' }}
              </template>
              <template v-else>
                Скрыть ответы
              </template>
            </span>
          </Transition>
        </button>
      </div>
      
      <div class="flex flex-col sm:flex-row items-start sm:items-center justify-between gap-2">
        <div class="flex items-center gap-2 sm:gap-3">
          <button
            v-if="isAuthenticated && !showReplyForm"
            @click="showReplyForm = true"
            class="text-xs sm:text-sm text-indigo-600 hover:text-indigo-700 font-medium"
          >
            Ответить
          </button>
          <button
            v-if="isAuthenticated && showReplyForm"
            @click="showReplyForm = false"
            class="text-xs sm:text-sm text-gray-500 hover:text-gray-700"
          >
            Отмена
          </button>
        </div>
        <button
          v-if="isAuthenticated && !comment.is_admin && !isAdmin"
          @click="handleDeleteClick"
          class="text-xs sm:text-sm text-red-600 hover:text-red-700 font-medium"
        >
          Удалить
        </button>
      </div>

      <Transition name="slide-fade">
        <div v-if="showReplyForm" class="mt-3 pt-3 border-t border-gray-200">
          <textarea
            v-model="replyText"
            placeholder="Напишите ответ..."
            rows="2"
            class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 resize-none text-xs sm:text-sm"
          ></textarea>
          <div class="flex flex-col sm:flex-row justify-end gap-2 mt-2">
            <button
              @click="showReplyForm = false"
              class="px-3 py-1.5 text-xs sm:text-sm text-gray-600 hover:text-gray-800 rounded-lg"
            >
              Отмена
            </button>
            <button
              @click="submitReply"
              :disabled="!replyText.trim()"
              class="px-3 py-1.5 text-xs sm:text-sm bg-indigo-600 text-white rounded-lg hover:bg-indigo-700 disabled:opacity-50 disabled:cursor-not-allowed"
            >
              Отправить
            </button>
          </div>
        </div>
      </Transition>
    </div>

    <Transition name="slide-down">
      <div v-if="replies.length > 0 && isCommentExpanded" class="ml-4 sm:ml-8 mt-3">
        <div v-if="shouldShowCollapse && !isExpanded" class="mb-2">
          <button
            @click="isExpanded = true"
            class="text-xs sm:text-sm text-indigo-600 hover:text-indigo-700 font-medium flex items-center gap-1 transition-colors"
          >
            <svg class="w-4 h-4 transition-transform duration-200" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7" />
            </svg>
            Показать {{ replies.length }} {{ replies.length === 1 ? 'ответ' : replies.length < 5 ? 'ответа' : 'ответов' }}
          </button>
        </div>
        
        <div v-if="!shouldShowCollapse || isExpanded" class="space-y-3">
          <div v-if="shouldShowCollapse && isExpanded" class="mb-2">
            <button
              @click="isExpanded = false"
              class="text-xs sm:text-sm text-indigo-600 hover:text-indigo-700 font-medium flex items-center gap-1 transition-colors"
            >
              <svg class="w-4 h-4 rotate-180 transition-transform duration-200" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7" />
              </svg>
              Свернуть ответы
            </button>
          </div>
          
          <div
            v-for="reply in replies"
            :key="reply.id"
            class="relative"
          >
            <div v-if="reply.is_admin" class="absolute -left-2 top-0 bottom-0 w-1 bg-indigo-500 rounded-l"></div>
            <CommentItem
              :comment="reply"
              :all-comments="allComments"
              :is-authenticated="isAuthenticated"
              :is-admin="isAdmin"
              @reply="handleReply"
              @delete="handleDelete"
              @user-click="handleUserClickFromChild"
            />
          </div>
        </div>
      </div>
    </Transition>
  </div>

  <ConfirmationModal
    :is-open="confirmation.isOpen.value"
    :title="confirmation.options.value.title"
    :message="confirmation.options.value.message"
    :confirm-text="confirmation.options.value.confirmText"
    :confirm-variant="confirmation.options.value.confirmVariant"
    @confirm="confirmation.handleConfirm"
    @cancel="confirmation.handleCancel"
  />
</template>

<script setup lang="ts">
import { ref, computed, watch } from 'vue';
import { type Comment } from '../api/client';
import ConfirmationModal from './ui/ConfirmationModal.vue';
import { useConfirmation } from '../composables/useModal';

const props = defineProps<{
  comment: Comment;
  replies?: Comment[];
  allComments?: Comment[];
  isAuthenticated: boolean;
  isAdmin?: boolean;
}>();

const emit = defineEmits<{
  reply: [commentId: string, text: string];
  delete: [commentId: string];
  'user-click': [userId: number];
}>();

const showReplyForm = ref(false);
const replyText = ref('');
const isExpanded = ref(true);

const replies = computed(() => {
  if (props.replies && props.replies.length > 0) {
    return props.replies;
  }
  if (props.allComments) {
    return props.allComments.filter(c => c.reply_to === props.comment.id);
  }
  return [];
});

const isCommentExpanded = ref(true);

watch(replies, (newReplies) => {
  if (newReplies.length > 3 && isCommentExpanded.value) {
    isCommentExpanded.value = false;
  }
}, { immediate: true });

const shouldShowCollapse = computed(() => {
  return replies.value.length > 5;
});

const hasReplies = computed(() => {
  return replies.value.length > 0;
});

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

const handleUserClickFromChild = (userId: number) => {
  emit('user-click', userId);
};

const confirmation = useConfirmation();

const handleDeleteClick = async () => {
  const confirmed = await confirmation.confirm({
    message: 'Удалить этот комментарий?',
    confirmVariant: 'danger',
    confirmText: 'Удалить',
  });
  if (confirmed) {
    emit('delete', props.comment.id);
  }
};

const getInitial = () => {
  if (props.comment.is_admin) {
    return 'А';
  }
  if (props.comment.user_email) {
    return props.comment.user_email.charAt(0).toUpperCase();
  }
  return 'П';
};

const handleUserClick = () => {
  emit('user-click', props.comment.user_id);
};
</script>

<style scoped>
.comment-item {
  @apply w-full;
}

/* Анимация для формы ответа */
.slide-fade-enter-active {
  transition: all 0.3s ease-out;
}

.slide-fade-leave-active {
  transition: all 0.25s ease-in;
}

.slide-fade-enter-from {
  transform: translateY(-10px);
  opacity: 0;
}

.slide-fade-leave-to {
  transform: translateY(-10px);
  opacity: 0;
}

/* Анимация для текста кнопки */
.fade-enter-active {
  transition: opacity 0.2s ease-in-out;
}

.fade-leave-active {
  transition: opacity 0.15s ease-in-out;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}

/* Анимация для раскрывающихся ответов */
.slide-down-enter-active {
  transition: max-height 0.5s cubic-bezier(0.4, 0, 0.2, 1),
              opacity 0.4s cubic-bezier(0.4, 0, 0.2, 1),
              transform 0.4s cubic-bezier(0.4, 0, 0.2, 1);
  overflow: hidden;
}

.slide-down-leave-active {
  transition: max-height 0.4s cubic-bezier(0.4, 0, 0.2, 1),
              opacity 0.3s cubic-bezier(0.4, 0, 0.2, 1),
              transform 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  overflow: hidden;
}

.slide-down-enter-from {
  max-height: 0;
  opacity: 0;
  transform: translateY(-10px);
}

.slide-down-enter-to {
  max-height: 5000px;
  opacity: 1;
  transform: translateY(0);
}

.slide-down-leave-from {
  max-height: 5000px;
  opacity: 1;
  transform: translateY(0);
}

.slide-down-leave-to {
  max-height: 0;
  opacity: 0;
  transform: translateY(-10px);
}
</style>

