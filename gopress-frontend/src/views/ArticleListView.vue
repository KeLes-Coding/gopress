<script setup>
import { ref, onMounted } from 'vue'
import { fetchPostList } from '../api/post'

// 1. 定义响应式变量来存储文章列表
const posts = ref([])
// 定义加载状态和错误信息
const isLoading = ref(true)
const error = ref(null)

// 2. 定义一个方法来获取文章数据
async function getPosts() {
  try {
    isLoading.value = true
    const response = await fetchPostList({ page: 1, pageSize: 10 })
    // 从后端返回的数据结构中取出文章列表
    posts.value = response.data.data.post 
  } catch (err) {
    error.value = '加载文章失败，请稍后再试。'
    console.error(err)
  } finally {
    isLoading.value = false
  }
}

// 3. 在组件挂载完成后，调用该方法
onMounted(() => {
  getPosts()
})
</script>

<template>
  <div class="article-list-page">
    <h1>博客文章</h1>
    
    <div v-if="isLoading" class="loading">正在加载中...</div>
    
    <div v-else-if="error" class="error">{{ error }}</div>
    
    <ul v-else-if="posts.length > 0" class="post-list">
      <li v-for="post in posts" :key="post.ID" class="post-item">
        <h2>
          <router-link :to="`/posts/${post.ID}`">{{ post.Title }}</router-link>
        </h2>
        <p class="summary">{{ post.Summary }}</p>
        <div class="meta">
          <span>作者: {{ post.User.Username }}</span>
          <span>分类: {{ post.Category.Name }}</span>
          <span>发布于: {{ new Date(post.CreatedAt).toLocaleDateString() }}</span>
        </div>
      </li>
    </ul>
    
    <div v-else class="empty">暂无文章</div>
  </div>
</template>

<style scoped>
/* 样式只在本组件内生效 */
.article-list-page h1 {
  margin-bottom: 2rem;
}
.post-list {
  list-style: none;
  padding: 0;
}
.post-item {
  margin-bottom: 2.5rem;
  border-bottom: 1px solid #eee;
  padding-bottom: 1.5rem;
}
.post-item h2 {
  margin: 0 0 0.5rem 0;
}
.post-item h2 a {
  color: #333;
  text-decoration: none;
  transition: color 0.3s;
}
.post-item h2 a:hover {
  color: #007bff;
}
.summary {
  color: #666;
  line-height: 1.6;
}
.meta {
  font-size: 0.85rem;
  color: #999;
}
.meta span {
  margin-right: 1rem;
}
.loading, .error, .empty {
  text-align: center;
  padding: 2rem;
  color: #888;
}
</style>