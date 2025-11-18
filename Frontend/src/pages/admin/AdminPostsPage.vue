<template>
  <div class="max-w-7xl mx-auto px-4 py-10 space-y-6">
    <div class="flex flex-col gap-4 md:flex-row md:items-center md:justify-between">
      <div>
        <h1 class="text-3xl font-bold text-gray-900">Новости / Посты</h1>
        <p class="text-gray-500 mt-1">Управление постами</p>
      </div>
      <ui-button type="button" variant="primary" @click="showCreateModal = true">Создать пост</ui-button>
    </div>

    <div v-if="loading" class="bg-white rounded-xl shadow p-6 text-center text-gray-500">
      Загрузка...
    </div>

    <div v-else-if="error" class="bg-red-50 border border-red-200 rounded-xl p-6">
      <p class="text-red-700">{{ error }}</p>
    </div>

    <div v-else class="bg-white rounded-xl shadow overflow-hidden">
      <div class="overflow-x-auto -mx-4 sm:mx-0">
        <table class="min-w-full divide-y divide-gray-200">
          <thead class="bg-gray-50">
            <tr>
              <th class="px-3 sm:px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Заголовок</th>
              <th class="px-3 sm:px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider hidden sm:table-cell">Статус</th>
              <th class="px-3 sm:px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider hidden md:table-cell">Дата</th>
              <th class="px-3 sm:px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Действия</th>
            </tr>
          </thead>
          <tbody class="bg-white divide-y divide-gray-200">
            <tr v-for="post in posts" :key="post.id" class="hover:bg-gray-50">
              <td class="px-3 sm:px-6 py-4 text-xs sm:text-sm text-gray-900 max-w-xs truncate sm:max-w-none sm:truncate-none">{{ post.title }}</td>
              <td class="px-3 sm:px-6 py-4 whitespace-nowrap hidden sm:table-cell">
                <span :class="post.is_published ? 'bg-green-100 text-green-800' : 'bg-gray-100 text-gray-800'" class="px-2 py-1 text-xs font-semibold rounded-full">
                  {{ post.is_published ? 'Опубликован' : 'Черновик' }}
                </span>
              </td>
              <td class="px-3 sm:px-6 py-4 whitespace-nowrap text-xs sm:text-sm text-gray-500 hidden md:table-cell">{{ formatDate(post.created_at) }}</td>
              <td class="px-3 sm:px-6 py-4 whitespace-nowrap text-xs sm:text-sm font-medium">
                <div class="flex flex-col sm:flex-row gap-1 sm:gap-2">
                  <button @click="editPost(post)" class="text-indigo-600 hover:text-indigo-900 text-left sm:text-center">Редактировать</button>
                  <button @click="deletePost(post.id)" class="text-red-600 hover:text-red-900 text-left sm:text-center">Удалить</button>
                </div>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <Transition name="modal">
      <div v-if="showCreateModal || editingPost" class="fixed inset-0 bg-gray-900/60 backdrop-blur-sm flex items-center justify-center z-50 p-4" @click.self="closeModal">
        <div class="bg-white rounded-xl shadow-2xl max-w-4xl w-full max-h-[90vh] flex flex-col transform transition-all m-4">
          <div class="px-4 sm:px-6 pt-4 sm:pt-6 pb-4 border-b border-gray-200 relative flex-shrink-0">
            <button
              @click="closeModal"
              class="absolute top-4 right-4 sm:top-6 sm:right-6 text-gray-400 hover:text-gray-600 transition-colors focus:outline-none focus:ring-2 focus:ring-indigo-500 rounded-full p-1"
              type="button"
              aria-label="Закрыть"
            >
              <svg class="w-5 h-5 sm:w-6 sm:h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
              </svg>
            </button>
            <h2 class="text-xl sm:text-2xl font-bold text-gray-900 pr-8">{{ editingPost ? 'Редактировать пост' : 'Создать пост' }}</h2>
          </div>
          <form @submit.prevent="savePost" class="flex-1 overflow-y-auto px-4 sm:px-6 py-4 space-y-4">
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-1">Заголовок</label>
              <ui-input v-model="formData.title" variant="primary" type="text" required />
            </div>
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-1">Содержание</label>
              <textarea v-model="formData.content" rows="8" class="w-full px-3 py-2 border border-gray-300 rounded-md text-sm sm:text-base resize-y" required></textarea>
            </div>
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-1">Изображения</label>
              <div class="space-y-2">
                <div>
                  <ui-input v-model="imagesInput" variant="primary" type="text" placeholder="URL через запятую: https://example.com/img1.jpg, https://example.com/img2.jpg" />
                </div>
                <div class="flex items-center gap-2">
                  <span class="text-sm text-gray-600">или</span>
                </div>
                <div>
                  <input
                    type="file"
                    ref="fileInput"
                    @change="handleFileUpload"
                    multiple
                    accept="image/*"
                    class="hidden"
                  />
                  <button
                    type="button"
                    @click="fileInput?.click()"
                    class="px-4 py-2 bg-gray-100 hover:bg-gray-200 text-gray-700 rounded-md text-sm font-medium transition-colors"
                  >
                    Загрузить изображения
                  </button>
                  <div v-if="uploadingImages.length > 0" class="mt-2 space-y-1">
                    <div v-for="(status, index) in uploadingImages" :key="index" class="text-xs text-gray-600">
                      {{ status }}
                    </div>
                  </div>
                </div>
                <div v-if="formData.images.length > 0" class="mt-2 space-y-3">
                  <div class="text-xs text-gray-600 mb-1">Загруженные изображения:</div>
                  <div v-for="(img, index) in formData.images" :key="index" class="border border-gray-200 rounded p-2 space-y-2">
                    <div class="flex items-center gap-2">
                      <img :src="typeof img === 'string' ? img : img.url" :alt="`Изображение ${index + 1}`" class="w-20 h-20 object-cover rounded border border-gray-300" />
                      <div class="flex-1">
                        <label class="block text-xs font-medium text-gray-700 mb-1">Размер изображения</label>
                        <select
                          :value="(img as MediaItem).size || 'medium'"
                          @change="formData.images[index] = { ...(img as MediaItem), size: ($event.target as HTMLSelectElement).value as MediaSize }"
                          class="w-full px-2 py-1 text-xs border border-gray-300 rounded focus:border-blue-500 focus:ring-1 focus:ring-blue-200"
                        >
                          <option value="small">Маленький</option>
                          <option value="medium">Средний</option>
                          <option value="large">Большой</option>
                          <option value="xlarge">Крупный</option>
                        </select>
                      </div>
                      <button
                        type="button"
                        @click="removeImage(index)"
                        class="text-red-600 hover:text-red-800 text-sm font-medium"
                      >
                        Удалить
                      </button>
                    </div>
                  </div>
                </div>
              </div>
            </div>
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-1">Видео</label>
              <div class="space-y-2">
                <div>
                  <ui-input v-model="videosInput" variant="primary" type="text" placeholder="URL через запятую: https://example.com/video1.mp4, https://example.com/video2.mp4" />
                </div>
                <div class="flex items-center gap-2">
                  <span class="text-sm text-gray-600">или</span>
                </div>
                <div>
                  <input
                    type="file"
                    ref="videoInput"
                    @change="handleVideoUpload"
                    multiple
                    accept="video/*"
                    class="hidden"
                  />
                  <button
                    type="button"
                    @click="videoInput?.click()"
                    class="px-4 py-2 bg-gray-100 hover:bg-gray-200 text-gray-700 rounded-md text-sm font-medium transition-colors"
                  >
                    Загрузить видео
                  </button>
                  <div v-if="uploadingVideos.length > 0" class="mt-2 space-y-1">
                    <div v-for="(status, index) in uploadingVideos" :key="index" class="text-xs text-gray-600">
                      {{ status }}
                    </div>
                  </div>
                </div>
                <div v-if="formData.videos.length > 0" class="mt-2 space-y-3">
                  <div class="text-xs text-gray-600 mb-1">Загруженные видео:</div>
                  <div v-for="(video, index) in formData.videos" :key="index" class="border border-gray-200 rounded p-2 space-y-2">
                    <div class="flex items-center gap-2">
                      <video :src="typeof video === 'string' ? video : video.url" class="w-20 h-20 object-cover rounded border border-gray-300" muted></video>
                      <div class="flex-1">
                        <label class="block text-xs font-medium text-gray-700 mb-1">Размер видео</label>
                        <select
                          :value="(video as MediaItem).size || 'medium'"
                          @change="formData.videos[index] = { ...(video as MediaItem), size: ($event.target as HTMLSelectElement).value as MediaSize }"
                          class="w-full px-2 py-1 text-xs border border-gray-300 rounded focus:border-blue-500 focus:ring-1 focus:ring-blue-200"
                        >
                          <option value="small">Маленький</option>
                          <option value="medium">Средний</option>
                          <option value="large">Большой</option>
                          <option value="xlarge">Крупный</option>
                        </select>
                      </div>
                      <button
                        type="button"
                        @click="removeVideo(index)"
                        class="text-red-600 hover:text-red-800 text-sm font-medium"
                      >
                        Удалить
                      </button>
                    </div>
                  </div>
                </div>
              </div>
            </div>
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-1">Страницы для отображения</label>
              <PageSelector
                :selected-pages="formData.pages"
                :available-routes="availableRoutes"
                @update:selected-pages="formData.pages = $event"
              />
            </div>
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-1">Выравнивание поста</label>
              <select 
                v-model="formData.alignment" 
                class="w-full px-3 py-2 border border-gray-300 rounded-md focus:border-indigo-500 focus:ring-2 focus:ring-indigo-200 transition-all text-sm sm:text-base"
              >
                <option value="full-width">На всю ширину</option>
                <option value="left">Слева</option>
                <option value="center">По центру</option>
                <option value="right">Справа</option>
              </select>
            </div>
            <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
              <div>
                <label class="block text-sm font-medium text-gray-700 mb-1">Расположение заголовка</label>
                <select 
                  v-model="formData.title_position" 
                  class="w-full px-3 py-2 border border-gray-300 rounded-md focus:border-indigo-500 focus:ring-2 focus:ring-indigo-200 transition-all text-sm sm:text-base"
                >
                  <option value="top">Сверху</option>
                  <option value="bottom">Снизу</option>
                  <option value="left">Слева</option>
                  <option value="right">Справа</option>
                </select>
              </div>
              <div>
                <label class="block text-sm font-medium text-gray-700 mb-1">Расположение содержания</label>
                <select 
                  v-model="formData.content_position" 
                  class="w-full px-3 py-2 border border-gray-300 rounded-md focus:border-indigo-500 focus:ring-2 focus:ring-indigo-200 transition-all text-sm sm:text-base"
                >
                  <option value="top">Сверху</option>
                  <option value="bottom">Снизу</option>
                  <option value="left">Слева</option>
                  <option value="right">Справа</option>
                </select>
              </div>
            </div>
            <div class="flex items-center gap-2">
              <input type="checkbox" v-model="formData.is_published" id="published" class="w-4 h-4 text-indigo-600 border-gray-300 rounded" />
              <label for="published" class="text-sm font-medium text-gray-700">Опубликовать</label>
            </div>
          </form>
          <div class="px-4 sm:px-6 py-4 border-t border-gray-200 flex flex-col sm:flex-row gap-2 justify-end flex-shrink-0">
            <ui-button type="button" variant="secondary" @click="closeModal">Отмена</ui-button>
            <ui-button type="button" variant="primary" @click="savePost">Сохранить</ui-button>
          </div>
        </div>
      </div>
    </Transition>

    <ConfirmationModal
      :is-open="confirmation.isOpen.value"
      :title="confirmation.options.value.title"
      :message="confirmation.options.value.message"
      :confirm-text="confirmation.options.value.confirmText"
      :confirm-variant="confirmation.options.value.confirmVariant"
      @confirm="confirmation.handleConfirm"
      @cancel="confirmation.handleCancel"
    />

    <NotificationModal
      :is-open="notification.isOpen.value"
      :type="notification.options.value.type"
      :title="notification.options.value.title"
      :message="notification.options.value.message"
      @close="notification.close"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, computed } from 'vue';
import { useRouter } from 'vue-router';
import UiButton from '../../components/ui/Ui-button.vue';
import UiInput from '../../components/ui/Ui-input.vue';
import PageSelector from '../../components/ui/PageSelector.vue';
import ConfirmationModal from '../../components/ui/ConfirmationModal.vue';
import NotificationModal from '../../components/ui/NotificationModal.vue';
import { adminApi, type Post } from '../../api/admin';
import { adminAuthService } from '../../utils/adminAuth';
import { useConfirmation, useNotification } from '../../composables/useModal';
import type { MediaItem, MediaSize } from '../../api/client';

const router = useRouter();
const posts = ref<Post[]>([]);
const loading = ref(true);
const error = ref('');
const showCreateModal = ref(false);
const editingPost = ref<Post | null>(null);
const fileInput = ref<HTMLInputElement | null>(null);
const videoInput = ref<HTMLInputElement | null>(null);
const uploadingImages = ref<string[]>([]);
const uploadingVideos = ref<string[]>([]);

const routeDescriptions: Record<string, string> = {
  '/': 'Главная страница (Новости)',
  '/about': 'О себе',
  '/portfolio/about': 'Портфолио - О себе',
  '/portfolio/qualification': 'Портфолио - Повышение квалификации',
  '/portfolio/events': 'Портфолио - Участие в мероприятиях',
  '/portfolio/awards': 'Портфолио - Дипломы, сертификаты и благодарности',
  '/portfolio/publications': 'Портфолио - Публикации',
  '/methodology/programs': 'Методическая копилка - Программы',
  '/methodology/developments': 'Методическая копилка - Методические разработки',
  '/methodology/law': 'Методическая копилка - Нормативно-правовая база',
  '/methodology/videos': 'Методическая копилка - Видеозаписи',
  '/achievements/awards': 'Достижения - Дипломы и сертификаты',
  '/achievements/gallery': 'Достижения - Фотогалерея',
  '/parents/tips': 'Родителям - Советы',
  '/parents/day-schedule': 'Родителям - Режим дня школьника',
  '/parents/posture': 'Родителям - Нарушение осанки',
  '/parents/flat-feet': 'Родителям - Плоскостопие',
  '/parents/illness-prevention': 'Родителям - Профилактика заболеваний ОРЗ и ОРВИ',
  '/parents/hardening': 'Родителям - Закаливание в семье',
};

const parentPagesToExclude = ['/portfolio', '/methodology', '/achievements', '/parents'];

const availableRoutes = computed(() => {
  return router.getRoutes().filter(route => {
    if (route.path.startsWith('/admin') || 
        route.path === '/user-login' || 
        route.path === '/register') {
      return false;
    }
    if (parentPagesToExclude.includes(route.path)) {
      return false;
    }
    return true;
  }).map(route => {
    const description = routeDescriptions[route.path] || route.path;
    return {
      path: route.path,
      name: route.name || route.path,
      displayName: description
    };
  });
});

const formData = ref({
  title: '',
  content: '',
  images: [] as MediaItem[],
  videos: [] as MediaItem[],
  pages: [] as string[],
  is_published: false,
  alignment: 'full-width' as 'left' | 'center' | 'right' | 'full-width',
  title_position: 'top' as 'top' | 'bottom' | 'left' | 'right',
  content_position: 'bottom' as 'top' | 'bottom' | 'left' | 'right',
});

const imagesInput = computed({
  get: () => formData.value.images.map(img => typeof img === 'string' ? img : img.url).join(', '),
  set: (val) => {
    formData.value.images = val.split(',').map(s => s.trim()).filter(Boolean).map(url => ({ url, size: 'medium' as MediaSize }));
  },
});

const videosInput = computed({
  get: () => formData.value.videos.map(vid => typeof vid === 'string' ? vid : vid.url).join(', '),
  set: (val) => {
    formData.value.videos = val.split(',').map(s => s.trim()).filter(Boolean).map(url => ({ url, size: 'medium' as MediaSize }));
  },
});

const formatDate = (date: string) => {
  return new Date(date).toLocaleDateString('ru-RU');
};

const fetchPosts = async () => {
  loading.value = true;
  error.value = '';
  try {
    posts.value = await adminApi.getPosts();
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

const normalizeMedia = (media: (string | MediaItem)[]): MediaItem[] => {
  return media.map(item => typeof item === 'string' ? { url: item, size: 'medium' as MediaSize } : item);
};

const editPost = (post: Post) => {
  editingPost.value = post;
  formData.value = {
    title: post.title,
    content: post.content,
    images: normalizeMedia(post.images || []),
    videos: normalizeMedia(post.videos || []),
    pages: post.pages || [],
    is_published: post.is_published,
    alignment: post.alignment || 'full-width',
    title_position: post.title_position || 'top',
    content_position: post.content_position || 'bottom',
  };
};

const confirmation = useConfirmation();
const notification = useNotification();

const deletePost = async (id: string) => {
  const confirmed = await confirmation.confirm({
    message: 'Удалить этот пост?',
    confirmVariant: 'danger',
    confirmText: 'Удалить',
  });
  if (!confirmed) return;
  try {
    await adminApi.deletePost(id);
    await fetchPosts();
  } catch (err) {
    notification.error(err instanceof Error ? err.message : 'Ошибка удаления');
  }
};

const savePost = async () => {
  try {
    const postData = {
      ...formData.value,
      images: formData.value.images.map(img => typeof img === 'string' ? img : img.url),
      videos: formData.value.videos.map(vid => typeof vid === 'string' ? vid : vid.url),
    };
    
    if (editingPost.value) {
      await adminApi.updatePost(editingPost.value.id, postData);
    } else {
      await adminApi.createPost(postData);
    }
    closeModal();
    await fetchPosts();
  } catch (err) {
    notification.error(err instanceof Error ? err.message : 'Ошибка сохранения');
  }
};

const handleFileUpload = async (event: Event) => {
  const target = event.target as HTMLInputElement;
  const files = target.files;
  if (!files || files.length === 0) return;

  uploadingImages.value = [];
  
  for (let i = 0; i < files.length; i++) {
    const file = files[i];
    if (!file) continue;
    uploadingImages.value.push(`Загрузка ${file.name}...`);
    
    try {
      const url = await adminApi.uploadFile(file);
      const exists = formData.value.images.some(img => (typeof img === 'string' ? img : img.url) === url);
      if (!exists) {
        formData.value.images.push({ url, size: 'medium' as MediaSize });
      }
      uploadingImages.value[i] = `✓ ${file.name} загружен`;
    } catch (error) {
      uploadingImages.value[i] = `✗ Ошибка загрузки ${file.name}`;
      console.error('Error uploading file:', error);
    }
  }
  
  setTimeout(() => {
    uploadingImages.value = [];
  }, 3000);
  
  if (target) {
    target.value = '';
  }
};

const handleVideoUpload = async (event: Event) => {
  const target = event.target as HTMLInputElement;
  const files = target.files;
  if (!files || files.length === 0) return;

  uploadingVideos.value = [];
  
  for (let i = 0; i < files.length; i++) {
    const file = files[i];
    if (!file) continue;
    uploadingVideos.value.push(`Загрузка ${file.name}...`);
    
    try {
      const url = await adminApi.uploadFile(file);
      const exists = formData.value.videos.some(vid => (typeof vid === 'string' ? vid : vid.url) === url);
      if (!exists) {
        formData.value.videos.push({ url, size: 'medium' as MediaSize });
      }
      uploadingVideos.value[i] = `✓ ${file.name} загружен`;
    } catch (error) {
      uploadingVideos.value[i] = `✗ Ошибка загрузки ${file.name}`;
      console.error('Error uploading file:', error);
    }
  }
  
  setTimeout(() => {
    uploadingVideos.value = [];
  }, 3000);
  
  if (target) {
    target.value = '';
  }
};

const removeImage = (index: number) => {
  formData.value.images.splice(index, 1);
};

const removeVideo = (index: number) => {
  formData.value.videos.splice(index, 1);
};

const closeModal = () => {
  showCreateModal.value = false;
  editingPost.value = null;
  uploadingImages.value = [];
  uploadingVideos.value = [];
  formData.value = {
    title: '',
    content: '',
    images: [],
    videos: [],
    pages: [],
    is_published: false,
    alignment: 'full-width',
    title_position: 'top',
    content_position: 'bottom',
  };
};

onMounted(fetchPosts);
</script>

