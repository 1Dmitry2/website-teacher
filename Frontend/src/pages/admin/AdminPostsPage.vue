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
      <div class="overflow-x-auto">
        <table class="min-w-full divide-y divide-gray-200">
          <thead class="bg-gray-50">
            <tr>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Заголовок</th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Статус</th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Дата</th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Действия</th>
            </tr>
          </thead>
          <tbody class="bg-white divide-y divide-gray-200">
            <tr v-for="post in posts" :key="post.id">
              <td class="px-6 py-4 text-sm text-gray-900">{{ post.title }}</td>
              <td class="px-6 py-4 whitespace-nowrap">
                <span :class="post.is_published ? 'bg-green-100 text-green-800' : 'bg-gray-100 text-gray-800'" class="px-2 py-1 text-xs font-semibold rounded-full">
                  {{ post.is_published ? 'Опубликован' : 'Черновик' }}
                </span>
              </td>
              <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500">{{ formatDate(post.created_at) }}</td>
              <td class="px-6 py-4 whitespace-nowrap text-sm font-medium space-x-2">
                <button @click="editPost(post)" class="text-indigo-600 hover:text-indigo-900">Редактировать</button>
                <button @click="deletePost(post.id)" class="text-red-600 hover:text-red-900">Удалить</button>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <!-- Create/Edit Modal -->
    <Transition name="modal">
      <div v-if="showCreateModal || editingPost" class="fixed inset-0 bg-gray-900/60 backdrop-blur-sm flex items-center justify-center z-50 p-4" @click.self="closeModal">
        <div class="bg-white rounded-xl shadow-2xl max-w-4xl w-full max-h-[90vh] flex flex-col transform transition-all">
          <!-- Заголовок (фиксированный) -->
          <div class="px-6 pt-6 pb-4 border-b border-gray-200 relative flex-shrink-0">
            <button
              @click="closeModal"
              class="absolute top-6 right-6 text-gray-400 hover:text-gray-600 transition-colors focus:outline-none focus:ring-2 focus:ring-indigo-500 rounded-full p-1"
              type="button"
              aria-label="Закрыть"
            >
              <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
              </svg>
            </button>
            <h2 class="text-2xl font-bold text-gray-900 pr-8">{{ editingPost ? 'Редактировать пост' : 'Создать пост' }}</h2>
          </div>
          <!-- Контент (прокручиваемый) -->
          <form @submit.prevent="savePost" class="flex-1 overflow-y-auto px-6 py-4 space-y-4">
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-1">Заголовок</label>
              <ui-input v-model="formData.title" variant="primary" type="text" required />
            </div>
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-1">Содержание</label>
              <textarea v-model="formData.content" rows="10" class="w-full px-3 py-2 border border-gray-300 rounded-md" required></textarea>
            </div>
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-1">Изображения (URL через запятую)</label>
              <ui-input v-model="imagesInput" variant="primary" type="text" placeholder="https://example.com/img1.jpg, https://example.com/img2.jpg" />
            </div>
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-1">Страницы для отображения</label>
              <select 
                v-model="formData.pages" 
                multiple 
                class="w-full px-3 py-2 border border-gray-300 rounded-md min-h-[100px]"
              >
                <option v-for="route in availableRoutes" :key="route.path" :value="route.path">
                  {{ route.displayName }}
                </option>
              </select>
              <p class="text-xs text-gray-500 mt-1">Удерживайте Ctrl (Cmd на Mac) для выбора нескольких страниц</p>
            </div>
            <div class="flex items-center gap-2">
              <input type="checkbox" v-model="formData.is_published" id="published" class="w-4 h-4 text-indigo-600 border-gray-300 rounded" />
              <label for="published" class="text-sm font-medium text-gray-700">Опубликовать</label>
            </div>
          </form>
          <!-- Кнопки (фиксированные) -->
          <div class="px-6 py-4 border-t border-gray-200 flex gap-2 justify-end flex-shrink-0">
            <ui-button type="button" variant="secondary" @click="closeModal">Отмена</ui-button>
            <ui-button type="button" variant="primary" @click="savePost">Сохранить</ui-button>
          </div>
        </div>
      </div>
    </Transition>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, computed } from 'vue';
import { useRouter } from 'vue-router';
import UiButton from '../../components/ui/Ui-button.vue';
import UiInput from '../../components/ui/Ui-input.vue';
import { adminApi, type Post } from '../../api/admin';
import { adminAuthService } from '../../utils/adminAuth';

const router = useRouter();
const posts = ref<Post[]>([]);
const loading = ref(true);
const error = ref('');
const showCreateModal = ref(false);
const editingPost = ref<Post | null>(null);

// Маппинг роутов на описания
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

// Родительские страницы, которые нужно исключить (у них есть подстраницы)
const parentPagesToExclude = ['/portfolio', '/methodology', '/achievements', '/parents'];

// Получаем список доступных страниц из роутера (исключаем админские, служебные и родительские страницы)
const availableRoutes = computed(() => {
  return router.getRoutes().filter(route => {
    // Исключаем админские страницы, страницы входа и регистрации
    if (route.path.startsWith('/admin') || 
        route.path === '/user-login' || 
        route.path === '/register') {
      return false;
    }
    // Исключаем родительские страницы (кроме новостей "/")
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
  images: [] as string[],
  pages: [] as string[],
  is_published: false,
});

const imagesInput = computed({
  get: () => formData.value.images.join(', '),
  set: (val) => {
    formData.value.images = val.split(',').map(s => s.trim()).filter(Boolean);
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

const editPost = (post: Post) => {
  editingPost.value = post;
  formData.value = {
    title: post.title,
    content: post.content,
    images: post.images || [],
    pages: post.pages || [],
    is_published: post.is_published,
  };
};

const deletePost = async (id: string) => {
  if (!confirm('Удалить этот пост?')) return;
  try {
    await adminApi.deletePost(id);
    await fetchPosts();
  } catch (err) {
    alert(err instanceof Error ? err.message : 'Ошибка удаления');
  }
};

const savePost = async () => {
  try {
    if (editingPost.value) {
      await adminApi.updatePost(editingPost.value.id, formData.value);
    } else {
      await adminApi.createPost(formData.value);
    }
    closeModal();
    await fetchPosts();
  } catch (err) {
    alert(err instanceof Error ? err.message : 'Ошибка сохранения');
  }
};

const closeModal = () => {
  showCreateModal.value = false;
  editingPost.value = null;
  formData.value = {
    title: '',
    content: '',
    images: [],
    pages: [],
    is_published: false,
  };
};

onMounted(fetchPosts);
</script>

