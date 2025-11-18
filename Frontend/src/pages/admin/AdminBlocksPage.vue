<template>
  <div class="max-w-7xl mx-auto px-4 py-10 space-y-6">
    <div class="flex flex-col gap-4 md:flex-row md:items-center md:justify-between">
      <div>
        <h1 class="text-3xl font-bold text-gray-900">Блоки страниц</h1>
        <p class="text-gray-500 mt-1">Управление контентными блоками</p>
      </div>
      <ui-button type="button" variant="primary" @click="showCreateModal = true">Создать блок</ui-button>
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
              <th class="px-3 sm:px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Страница</th>
              <th class="px-3 sm:px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider hidden sm:table-cell">Тип</th>
              <th class="px-3 sm:px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider hidden md:table-cell">Порядок</th>
              <th class="px-3 sm:px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Действия</th>
            </tr>
          </thead>
          <tbody class="bg-white divide-y divide-gray-200">
            <tr v-for="block in blocks" :key="block.id" class="hover:bg-gray-50">
              <td class="px-3 sm:px-6 py-4 whitespace-nowrap text-xs sm:text-sm text-gray-900 max-w-xs truncate sm:max-w-none sm:truncate-none">{{ block.page }}</td>
              <td class="px-3 sm:px-6 py-4 whitespace-nowrap text-xs sm:text-sm text-gray-500 hidden sm:table-cell">{{ block.type }}</td>
              <td class="px-3 sm:px-6 py-4 whitespace-nowrap text-xs sm:text-sm text-gray-500 hidden md:table-cell">{{ block.display_order }}</td>
              <td class="px-3 sm:px-6 py-4 whitespace-nowrap text-xs sm:text-sm font-medium">
                <div class="flex flex-col sm:flex-row gap-1 sm:gap-2">
                  <button @click="editBlock(block)" class="text-indigo-600 hover:text-indigo-900 text-left sm:text-center">Редактировать</button>
                  <button @click="deleteBlock(block.id)" class="text-red-600 hover:text-red-900 text-left sm:text-center">Удалить</button>
                </div>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <Transition name="modal">
      <div v-if="showCreateModal || editingBlock" class="fixed inset-0 bg-gray-900/60 backdrop-blur-sm flex items-center justify-center z-50 p-4" @click.self="closeModal">
        <div class="bg-white rounded-xl shadow-2xl max-w-2xl w-full max-h-[90vh] flex flex-col transform transition-all m-4">
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
            <h2 class="text-xl sm:text-2xl font-bold text-gray-900 pr-8">{{ editingBlock ? 'Редактировать блок' : 'Создать блок' }}</h2>
          </div>
          
          <div class="flex-1 overflow-y-auto px-4 sm:px-6 py-4">
            <div v-if="!editingBlock && !selectedBlockType" class="space-y-4">
              <p class="text-gray-600 mb-4 text-sm sm:text-base">Выберите тип блока:</p>
              <div class="grid grid-cols-1 sm:grid-cols-2 gap-3">
                <button
                  type="button"
                  @click="selectBlockType('text')"
                  class="p-4 border-2 border-gray-300 rounded-lg hover:border-blue-500 hover:bg-blue-50 transition-all text-left"
                >
                  <div class="text-2xl mb-2">📝</div>
                  <div class="font-medium text-gray-900">Текстовый блок</div>
                </button>
                <button
                  type="button"
                  @click="selectBlockType('slider')"
                  class="p-4 border-2 border-gray-300 rounded-lg hover:border-blue-500 hover:bg-blue-50 transition-all text-left"
                >
                  <div class="text-2xl mb-2">📸</div>
                  <div class="font-medium text-gray-900">Слайдер</div>
                </button>
                <button
                  type="button"
                  @click="selectBlockType('gallery')"
                  class="p-4 border-2 border-gray-300 rounded-lg hover:border-blue-500 hover:bg-blue-50 transition-all text-left"
                >
                  <div class="text-2xl mb-2">🖼</div>
                  <div class="font-medium text-gray-900">Галерея</div>
                </button>
                <button
                  type="button"
                  @click="selectBlockType('video')"
                  class="p-4 border-2 border-gray-300 rounded-lg hover:border-blue-500 hover:bg-blue-50 transition-all text-left"
                >
                  <div class="text-2xl mb-2">🎞</div>
                  <div class="font-medium text-gray-900">Видео</div>
                </button>
                <button
                  type="button"
                  @click="selectBlockType('text-with-image')"
                  class="p-4 border-2 border-gray-300 rounded-lg hover:border-blue-500 hover:bg-blue-50 transition-all text-left"
                >
                  <div class="text-2xl mb-2">📄</div>
                  <div class="font-medium text-gray-900">Текст с изображением</div>
                </button>
              </div>
            </div>

            <form v-else @submit.prevent="saveBlock" class="space-y-4">
              <div>
                <label class="block text-sm font-medium text-gray-700 mb-1">Страницы для отображения</label>
                <PageSelector
                  :selected-pages="formData.pages"
                  :available-routes="availableRoutes"
                  @update:selected-pages="formData.pages = $event"
                />
              </div>
              <div>
                <label class="block text-sm font-medium text-gray-700 mb-1">Порядок</label>
                <ui-input 
                  v-model.number="formData.display_order" 
                  variant="primary" 
                  type="number" 
                  placeholder="1" 
                  min="1"
                />
              </div>
              
              <div class="border-t border-gray-200 pt-4 mt-4">
                <h3 class="text-lg font-semibold text-gray-900 mb-4">Содержимое блока</h3>
                
                <BlockFormText 
                  v-if="formData.type === 'text' && blockContent"
                  v-model="blockContent as TextBlockContent"
                />
                
                <BlockFormSlider 
                  v-if="formData.type === 'slider' && blockContent"
                  v-model="blockContent as SliderBlockContent"
                />
                
                <BlockFormGallery 
                  v-if="formData.type === 'gallery' && blockContent"
                  v-model="blockContent as GalleryBlockContent"
                />
                
                <BlockFormVideo 
                  v-if="formData.type === 'video' && blockContent"
                  v-model="blockContent as VideoBlockContent"
                />
                
                <BlockFormTextWithImage 
                  v-if="formData.type === 'text-with-image' && blockContent"
                  v-model="blockContent as TextWithImageBlockContent"
                />
              </div>
            </form>
          </div>
          
          <div class="px-4 sm:px-6 py-4 border-t border-gray-200 flex flex-col sm:flex-row gap-2 justify-end flex-shrink-0">
            <ui-button v-if="!editingBlock && !selectedBlockType" type="button" variant="secondary" @click="closeModal">Отмена</ui-button>
            <template v-else>
              <ui-button type="button" variant="secondary" @click="closeModal">Отмена</ui-button>
              <ui-button type="button" variant="primary" @click="saveBlock">Сохранить</ui-button>
            </template>
          </div>
        </div>
      </div>
    </Transition>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, computed, watch } from 'vue';
import { useRouter } from 'vue-router';
import UiButton from '../../components/ui/Ui-button.vue';
import UiInput from '../../components/ui/Ui-input.vue';
import BlockFormText from '../../components/blocks/BlockFormText.vue';
import BlockFormSlider from '../../components/blocks/BlockFormSlider.vue';
import BlockFormGallery from '../../components/blocks/BlockFormGallery.vue';
import BlockFormVideo from '../../components/blocks/BlockFormVideo.vue';
import BlockFormTextWithImage from '../../components/blocks/BlockFormTextWithImage.vue';
import PageSelector from '../../components/ui/PageSelector.vue';
import { adminApi, type Block } from '../../api/admin';
import { adminAuthService } from '../../utils/adminAuth';
import type { TextBlockContent } from '../../components/blocks/BlockFormText.vue';
import type { SliderBlockContent } from '../../components/blocks/BlockFormSlider.vue';
import type { GalleryBlockContent } from '../../components/blocks/BlockFormGallery.vue';
import type { VideoBlockContent } from '../../components/blocks/BlockFormVideo.vue';
import type { TextWithImageBlockContent } from '../../components/blocks/BlockFormTextWithImage.vue';

const router = useRouter();
const blocks = ref<Block[]>([]);
const loading = ref(true);
const error = ref('');
const showCreateModal = ref(false);
const editingBlock = ref<Block | null>(null);
const selectedBlockType = ref<Block['type'] | null>(null);
const blockContent = ref<TextBlockContent | SliderBlockContent | GalleryBlockContent | VideoBlockContent | TextWithImageBlockContent | null>(null);

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
  page: '',
  pages: [] as string[],
  type: 'text' as Block['type'],
  display_order: 1,
  content: {} as Record<string, any>,
});

const initializeBlockContent = (type: Block['type'], existingContent?: Record<string, any>) => {
  const content = existingContent && typeof existingContent === 'object' ? existingContent : {};
  
  switch (type) {
    case 'text':
      blockContent.value = {
        title: content.title || '',
        text: content.text || '',
        alignment: content.alignment || 'left',
        style: content.style || 'regular',
      } as TextBlockContent;
      break;
    case 'slider':
      blockContent.value = {
        slides: Array.isArray(content.slides) ? content.slides : [],
      } as SliderBlockContent;
      break;
    case 'gallery':
      blockContent.value = {
        images: Array.isArray(content.images) ? content.images : [],
      } as GalleryBlockContent;
      break;
    case 'video':
      blockContent.value = {
        url: content.url || '',
        autoplay: content.autoplay ?? false,
        alignment: content.alignment || 'full-width',
        maxWidth: content.maxWidth,
        size: content.size || 'medium',
      } as VideoBlockContent;
      break;
    case 'text-with-image':
      blockContent.value = {
        title: content.title || '',
        text: content.text || '',
        image: content.image || '',
        textPosition: content.textPosition || 'right',
        alignment: content.alignment || 'left',
        style: content.style || 'regular',
        imageSize: content.imageSize || 'medium',
      } as TextWithImageBlockContent;
      break;
  }
};

const selectBlockType = (type: Block['type']) => {
  selectedBlockType.value = type;
  formData.value.type = type;
  initializeBlockContent(type);
};

watch(() => formData.value.display_order, (newValue) => {
  if (typeof newValue === 'number' && (isNaN(newValue) || newValue < 1)) {
    formData.value.display_order = 1;
  }
});

const fetchBlocks = async () => {
  loading.value = true;
  error.value = '';
  try {
    blocks.value = await adminApi.getBlocks();
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

const editBlock = (block: Block) => {
  editingBlock.value = block;
  selectedBlockType.value = block.type;
  formData.value = {
    page: block.page,
    pages: block.pages || [],
    type: block.type,
    display_order: block.display_order || 1,
    content: block.content,
  };
  initializeBlockContent(block.type, block.content);
};

const deleteBlock = async (id: string) => {
  if (!confirm('Удалить этот блок?')) return;
  try {
    await adminApi.deleteBlock(id);
    await fetchBlocks();
  } catch (err) {
    alert(err instanceof Error ? err.message : 'Ошибка удаления');
  }
};

const saveBlock = async () => {
  try {
    if (blockContent.value) {
      formData.value.content = blockContent.value as Record<string, any>;
    }
    
    if (formData.value.pages && formData.value.pages.length > 0 && formData.value.pages[0]) {
      formData.value.page = formData.value.pages[0];
    } else {
      formData.value.page = '';
    }
    
    if (editingBlock.value) {
      await adminApi.updateBlock(editingBlock.value.id, formData.value);
    } else {
      await adminApi.createBlock(formData.value);
    }
    closeModal();
    await fetchBlocks();
  } catch (err) {
    alert(err instanceof Error ? err.message : 'Ошибка сохранения');
  }
};

const closeModal = () => {
  showCreateModal.value = false;
  editingBlock.value = null;
  selectedBlockType.value = null;
  blockContent.value = null;
  formData.value = {
    page: '',
    pages: [],
    type: 'text',
    display_order: 1,
    content: {},
  };
};

onMounted(fetchBlocks);
</script>

