import { createApp } from 'vue'
import './style.css'
import App from './App.vue'

// 1. 导入路由实例
import router from './router'

// 2. 导入 Pinia
import {createPinia} from 'pinia'

// 3. 创建 Pinia 实例
const pinia = createPinia()
const app = createApp(App)

// 4. 挂载应用前，使用 .use() 安装插件
app.use(router) // 安装路由插件
app.use(pinia) // 使用 Pinia 插件

// 5. 最后挂载应用
app.mount('#app')
