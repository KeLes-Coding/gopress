// 1. 导入 axios
import axios from 'axios'

// 2. 导入 user store 和 Element Plus 的消息提示组件
import { useUserStore } from '../store/user'
import { ElMessage } from 'element-plus'

// 3. 创建一个新的 axios 实例
const service = axios.create({
    // baseURL: 'http://localhost:8080/api/v1', // 后端 API 的基础地址
    baseURL: import.meta.env.VITE_API_BASE_URL, // 使用环境变量
    timeout: 10000, // 请求超时时间 10 秒
})

// --- 请求拦截器 ---
// 在每次请求发送前，执行这个函数
service.interceptors.request.use(
    config => {
        // 获取 user store
        const userStore = useUserStore()

        // 如果 store 中有 token，则为请求头添加 Authorization 字段
        if (userStore.token) {
            config.headers['Authorization'] = 'Bearer ${useUserStore.token}'
        }
        return config
    },
    error => {
        // 对请求错误作处理
        console.log(error)
        return Promise.reject(error)
    }
)

// --- 响应拦截器 ---
// 在接收到响应后，执行这个函数
service.interceptors.response.use(
    response => {
        // 从响应中解构出业务状态码(code)、消息(message)和数据(data)
        const res = response.data
        // 我们的后端约定 code 为 200 时表示成功
        if (res.code !== 200) {
            // 如果 code 不为 200，则表示有业务错误
            ElMessage({
                message: res.message || 'Error',
                type: 'error',
                duration: 5 * 1000
            })
            // 在这里处理指定的错误码
            if (res.code === 401) {
                // TODO: 处理 token 失效，例如跳转到登录页
            }
            return Promise.reject(new Error(res.message || 'Error'))
        } else {
            // 如果成功，直接返回整个响应体，方便在 then 中解构
            return response
        }
    },
    error => {
        // 处理 HTTP 网络错误
        console.log('err' + error)
        ElMessage({
            message: error.message,
            type: 'error',
            duration: 5 * 1000
        })
        return Promise.reject(error)
    }
)

// 4. 导出实例
export default service