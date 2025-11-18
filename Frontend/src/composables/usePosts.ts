import { ref, computed, onMounted, watch } from 'vue';
import { useRoute } from 'vue-router';
import { apiClient, type Post } from '../api/client';

export function usePosts() {
  const route = useRoute();
  const posts = ref<Post[]>([]);
  const loading = ref(true);
  const error = ref('');

  const filteredPosts = computed(() => {
    const currentPath = route.path;
    
    return posts.value.filter(post => {
      if (!post.pages) {
        return false;
      }
      
      const pagesArray = Array.isArray(post.pages) ? post.pages : [];
      if (pagesArray.length === 0) {
        return false;
      }
      
      const matches = pagesArray.some(page => {
        const normalizedPage = String(page).trim();
        const normalizedPath = currentPath.trim();
        return normalizedPage === normalizedPath;
      });
      
      return matches;
    });
  });

  const fetchPosts = async () => {
    loading.value = true;
    error.value = '';
    try {
      posts.value = await apiClient.getPosts();
    } catch (err) {
      error.value = err instanceof Error ? err.message : 'Ошибка загрузки постов';
    } finally {
      loading.value = false;
    }
  };

  onMounted(() => {
    fetchPosts();
  });

  watch(() => route.path, () => {
    fetchPosts();
  });

  return {
    posts,
    filteredPosts,
    loading,
    error,
    fetchPosts,
  };
}

