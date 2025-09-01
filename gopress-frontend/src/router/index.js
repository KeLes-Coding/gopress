// 1. 从 vue-router 中导入所需函数
import { createRouter, createWebHistory } from 'vue-router'
import { useUserStore } from '../store/user'
// 2. 导入页面组件
// 前台博客页面
import ArticleListView from '../views/ArticleListView.vue'
import ArticleDetailView from '../views/ArticleDetailView.vue'

// 后台管理页面
import AdminLayout from '../layout/AdminLayout.vue'
import DashboardView from '../views/admin/DashboardView.vue'
import PostManagementView from '../views/admin/PostManagementView.vue'
import CategoryManagementView from '../views/admin/CategoryManagementView.vue'
import LoginView from '../views/LoginView.vue'
import RegisterView from '../views/RegisterView.vue'

// 3. 定义路由规则
// 每个路由规则都是一个对象，包含 path 和 component
const routes = [
    // 在这里添加具体的页面路由
    // 首页
    {
    path: '/', // 当用户访问网站根路径时
    name: 'ArticleList',
    component: ArticleListView // 显示 ArticleListView 组件
  },
  // 文章详情页
  {
    path: '/posts/:id', // 当用户访问类似 /posts/1, /posts/2 这样的路径时
    name: 'ArticleDetail',
    component: ArticleDetailView
  },

  // 登陆和注册路由
  {
    path: '/login',
    name: 'Login',
    component: LoginView,
  },
  {
    path: '/register',
    name: 'Register',
    component: RegisterView,
  },

  // 后台管理路由
  {
    path: '/admin', // 所有以 /admin 开头的路径都会先加载 AdminLayout 组件
    name: 'Admin',
    component: AdminLayout, // AdminLayout 作为这组路由的父级，提供了整体布局
    meta: { requiresAuth: true }, // 添加 meta 字段
    children: [
      {
        path: 'dashboard',
        name: 'AdminDashboard', 
        component: DashboardView,
      },
      {
        path: 'posts',
        name: 'AdminPostManagerment',
        component: PostManagementView,
      },
      {
        path: 'categories',
        name: 'AdminCategoryManagerment',
        component: CategoryManagementView,
      }
    ]
  }
]

// 4. 创建路由实例
const router = createRouter({
    // 使用 history 模式， URL 中不会出现
    history: createWebHistory(),
    routes, // routes: routes 的缩写
})

// 全局前置守卫
router.beforeEach((to, from, next) => {
  const userStore = useUserStore()

  // 检查目标路由是否需要认证
  if (to.matched.some(record => record.meta.requiresAuth)) {
    // 如果需要认证，但用户没有 token
    if (!userStore.token) {
      // 重定向到登陆页面
      next({ name: 'Login' })
    } else {
      // 如果有 token，则允许访问
      next()
    }
  } else {
    // 如果目标路由不需要认证，则直接允许访问
    next()
  }
})

// 5. 导出路由实例，以便在 main.js 中使用
export default router