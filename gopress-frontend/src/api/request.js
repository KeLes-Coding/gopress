// 1. 导入 axios
import axios from 'axios'

// 2. 创建一个新的 axios 实例
const service = axios.create({
    // baseURL: 'http://localhost:8080/api/v1', // 后端 API 的基础地址
    baseURL: import.meta.env.VITE_API_BASE_URL, // 使用环境变量
    timeout: 10000, // 请求超时时间 10 秒
})

// 3. 导出实例
export default service