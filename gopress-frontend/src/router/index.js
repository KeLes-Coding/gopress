// 1. 从 vue-router 中导入所需函数
import { createRouter, createWebHistory } from 'vue-router'
// 1. 导入页面组件
import ArticleListView from '../views/ArticleListView.vue'
import ArticleDetailView from '../views/ArticleDetailView.vue'

// 2. 定义路由规则
// 每个路由规则都是一个对象，包含 path 和 component
const routes = [
    // 在这里添加具体的页面路由
    // 首页
    {
    path: '/', // 当用户访问网站根路径时
    name: 'ArticleList',
    component: ArticleListView // 显示 ArticleListView 组件
  },
  // 3. 添加第二条路由规则：文章详情页
  {
    path: '/posts/:id', // 当用户访问类似 /posts/1, /posts/2 这样的路径时
    name: 'ArticleDetail',
    component: ArticleDetailView
  },
]

// 3. 创建路由实例
const router = createRouter({
    // 使用 history 模式， URL 中不会出现
    history: createWebHistory(),
    routes, // routes: routes 的缩写
})

// 4. 导出路由实例，以便在 main.js 中使用
export default router