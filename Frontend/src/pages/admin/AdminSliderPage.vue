<template>
  <div class="max-w-7xl mx-auto px-4 py-10 space-y-6">
    <div class="flex flex-col gap-4 md:flex-row md:items-center md:justify-between">
      <div>
        <h1 class="text-3xl font-bold text-gray-900">Слайдер</h1>
        <p class="text-gray-500 mt-1">Управление слайдером</p>
      </div>
      <ui-button type="button" variant="primary" @click="showCreateModal = true">Добавить элемент</ui-button>
    </div>

    <div v-if="loading" class="bg-white rounded-xl shadow p-6 text-center text-gray-500">Загрузка...</div>
    <div v-else-if="error" class="bg-red-50 border border-red-200 rounded-xl p-6"><p class="text-red-700">{{ error }}</p></div>
    <div v-else class="bg-white rounded-xl shadow overflow-hidden">
      <div class="overflow-x-auto">
        <table class="min-w-full divide-y divide-gray-200">
          <thead class="bg-gray-50"><tr>
            <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Порядок</th>
            <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Изображение</th>
            <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Название</th>
            <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Действия</th>
          </tr></thead>
          <tbody class="bg-white divide-y divide-gray-200">
            <tr v-for="item in slider" :key="item.id">
              <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500">{{ item.display_order }}</td>
              <td class="px-6 py-4"><img :src="item.image_url" :alt="item.title" class="h-16 w-24 object-cover rounded" /></td>
              <td class="px-6 py-4 text-sm text-gray-900">{{ item.title || 'Без названия' }}</td>
              <td class="px-6 py-4 whitespace-nowrap text-sm font-medium space-x-2">
                <button @click="editItem(item)" class="text-indigo-600 hover:text-indigo-900">Редактировать</button>
                <button @click="deleteItem(item.id)" class="text-red-600 hover:text-red-900">Удалить</button>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <Transition name="modal">
      <div v-if="showCreateModal || editingItem" class="fixed inset-0 bg-gray-900/60 backdrop-blur-sm flex items-center justify-center z-50 p-4" @click.self="closeModal">
        <div class="bg-white rounded-xl shadow-2xl max-w-md w-full max-h-[90vh] flex flex-col transform transition-all">
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
            <h2 class="text-2xl font-bold text-gray-900 pr-8">{{ editingItem ? 'Редактировать' : 'Добавить элемент' }}</h2>
          </div>
          <!-- Контент (прокручиваемый) -->
          <form @submit.prevent="saveItem" class="flex-1 overflow-y-auto px-6 py-4 space-y-4">
            <div><label class="block text-sm font-medium text-gray-700 mb-1">URL изображения</label><ui-input v-model="formData.image_url" variant="primary" type="text" required /></div>
            <div><label class="block text-sm font-medium text-gray-700 mb-1">Название</label><ui-input v-model="formData.title" variant="primary" type="text" /></div>
            <div><label class="block text-sm font-medium text-gray-700 mb-1">Описание</label><textarea v-model="formData.description" rows="3" class="w-full px-3 py-2 border border-gray-300 rounded-md"></textarea></div>
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
            <div><label class="block text-sm font-medium text-gray-700 mb-1">Порядок</label><ui-input v-model="formData.display_order" variant="primary" type="number" required /></div>
          </form>
          <!-- Кнопки (фиксированные) -->
          <div class="px-6 py-4 border-t border-gray-200 flex gap-2 justify-end flex-shrink-0">
            <ui-button type="button" variant="secondary" @click="closeModal">Отмена</ui-button>
            <ui-button type="button" variant="primary" @click="saveItem">Сохранить</ui-button>
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
import { adminApi, type SliderItem } from '../../api/admin';
import { adminAuthService } from '../../utils/adminAuth';

const router = useRouter();
const slider = ref<SliderItem[]>([]);
const loading = ref(true);
const error = ref('');
const showCreateModal = ref(false);
const editingItem = ref<SliderItem | null>(null);

// Маппинг роутов на описания
const routeDescriptions: Record<string, string> = {
  '/': 'Главная страница (Новости)',
  '/about': 'О себе',
  '/portfolio': 'Портфолио',
  '/portfolio/about': 'Портфолио - О себе',
  '/portfolio/qualification': 'Портфолио - Повышение квалификации',
  '/portfolio/events': 'Портфолио - Участие в мероприятиях',
  '/portfolio/awards': 'Портфолио - Дипломы, сертификаты и благодарности',
  '/portfolio/publications': 'Портфолио - Публикации',
  '/methodology': 'Методическая копилка',
  '/methodology/programs': 'Методическая копилка - Программы',
  '/methodology/developments': 'Методическая копилка - Методические разработки',
  '/methodology/law': 'Методическая копилка - Нормативно-правовая база',
  '/methodology/videos': 'Методическая копилка - Видеозаписи',
  '/achievements': 'Достижения обучающихся',
  '/achievements/awards': 'Достижения - Дипломы и сертификаты',
  '/achievements/gallery': 'Достижения - Фотогалерея',
  '/parents': 'Родителям',
  '/parents/tips': 'Родителям - Советы',
  '/parents/day-schedule': 'Родителям - Режим дня школьника',
  '/parents/posture': 'Родителям - Нарушение осанки',
  '/parents/flat-feet': 'Родителям - Плоскостопие',
  '/parents/illness-prevention': 'Родителям - Профилактика заболеваний ОРЗ и ОРВИ',
  '/parents/hardening': 'Родителям - Закаливание в семье',
};

// Получаем список доступных страниц из роутера (исключаем админские и служебные)
const availableRoutes = computed(() => {
  return router.getRoutes().filter(route => {
    // Исключаем админские страницы, страницы входа и регистрации
    return !route.path.startsWith('/admin') && 
           route.path !== '/user-login' && 
           route.path !== '/register';
  }).map(route => {
    const description = routeDescriptions[route.path] || route.path;
    return {
      path: route.path,
      name: route.name || route.path,
      displayName: description
    };
  });
});

const formData = ref({ image_url: '', title: '', description: '', pages: [] as string[], display_order: 0 });

const fetchSlider = async () => {
  loading.value = true;
  error.value = '';
  try {
    slider.value = await adminApi.getSlider();
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

const editItem = (item: SliderItem) => {
  editingItem.value = item;
  formData.value = { 
    image_url: item.image_url, 
    title: item.title || '', 
    description: item.description || '', 
    pages: item.pages || [],
    display_order: item.display_order 
  };
};

const deleteItem = async (id: string) => {
  if (!confirm('Удалить этот элемент?')) return;
  try {
    await adminApi.deleteSliderItem(id);
    await fetchSlider();
  } catch (err) {
    alert(err instanceof Error ? err.message : 'Ошибка удаления');
  }
};

const saveItem = async () => {
  try {
    if (editingItem.value) {
      await adminApi.updateSliderItem(editingItem.value.id, formData.value);
    } else {
      await adminApi.createSliderItem(formData.value);
    }
    closeModal();
    await fetchSlider();
  } catch (err) {
    alert(err instanceof Error ? err.message : 'Ошибка сохранения');
  }
};

const closeModal = () => {
  showCreateModal.value = false;
  editingItem.value = null;
  formData.value = { image_url: '', title: '', description: '', pages: [], display_order: 0 };
};

onMounted(fetchSlider);
</script>

