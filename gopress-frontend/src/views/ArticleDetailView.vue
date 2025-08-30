<script setup>
import { ref, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { fetchPostById } from '../api/post'

// 1. 获取当前路由信息
const route = useRoute()
// 从路由参数中获取文章 ID
const postId = route.params.id

// 2. 定义响应式变量来存储文章详情
const post = ref(null) // 初始值为 null，因为我们还没获取到数据
const isLoading = ref(true)
const error = ref(null)

// 3. 定义获取单篇文章数据的方法
async function getPostDetails() {
  try {
    isLoading.value = true
    const response = await fetchPostById(postId)
    // 后端直接返回文章对象，所以是 response.data.data
    post.value = response.data.data
  } catch (err) {
    error.value = '文章加载失败，可能已被删除或不存在。'
    console.error(err)
  } finally {
    isLoading.value = false
  }
}

// 4. 在组件挂载完成后，调用该方法
onMounted(() => {
  getPostDetails()
})
</script>

<template>
  <div class="article-detail-page">
    <router-link to="/" class="back-link">&larr; 返回列表</router-link>

    <div v-if="isLoading" class="loading">正在加载文章...</div>

    <div v-else-if="error" class="error">{{ error }}</div>

    <article v-else-if="post" class="post-content">
      <h1 class="post-title">{{ post.Title }}</h1>

      <div class="post-meta">
        <span>作者: {{ post.User.Username }}</span>
        <span>分类: {{ post.Category.Name }}</span>
        <span>发布于: {{ new Date(post.CreatedAt).toLocaleDateString() }}</span>
      </div>

      <div v-if="post.Tags && post.Tags.length" class="post-tags">
        <span v-for="tag in post.Tags" :key="tag.ID" class="tag">{{ tag.Name }}</span>
      </div>
      
      <div class="post-body" v-html="post.Content"></div>
    </article>
  </div>
</template>

<style scoped>
.article-detail-page {
  max-width: 800px;
  margin: 0 auto;
}
.back-link {
  display: inline-block;
  margin-bottom: 2rem;
  color: #007bff;
  text-decoration: none;
}
.back-link:hover {
  text-decoration: underline;
}
.loading, .error {
  text-align: center;
  padding: 3rem;
  color: #888;
}
.post-title {
  font-size: 2.5rem;
  margin-bottom: 1rem;
  line-height: 1.3;
}
.post-meta {
  color: #999;
  font-size: 0.9rem;
  margin-bottom: 1rem;
}
.post-meta span {
  margin-right: 1.5rem;
}
.post-tags {
  margin-bottom: 2rem;
}
.tag {
  display: inline-block;
  background-color: #f0f0f0;
  color: #555;
  padding: 0.25rem 0.75rem;
  border-radius: 1rem;
  font-size: 0.8rem;
  margin-right: 0.5rem;
  margin-bottom: 0.5rem;
}
.post-body {
  line-height: 1.8;
  font-size: 1.1rem;
  /* 确保 v-html 渲染的内容样式正常 */
  white-space: pre-wrap;
}
</style>