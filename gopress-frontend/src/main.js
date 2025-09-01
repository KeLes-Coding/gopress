import { createApp } from 'vue'
import './style.css'
import App from './App.vue'

// 1. 导入路由实例
import router from './router'

// 2. 导入 Pinia
import {createPinia} from 'pinia'

// 3. 导入 Element Plus
import ElementPlus from 'element-plus'
import 'element-plus/dist/index.css'

// 4. 创建 Pinia 实例
const pinia = createPinia()
const app = createApp(App)

// 5. 挂载应用前，使用 .use() 安装插件
app.use(router) // 安装路由插件
app.use(pinia) // 使用 Pinia 插件
app.use(ElementPlus)

// 6. 最后挂载应用
app.mount('#app')
