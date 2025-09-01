import { defineStore } from 'pinia'
import { ref } from 'vue'

// 使用 defineStore 创建一个 store
// 第一个参数是 store 的唯一 ID，Pinia 会用它来连接 devtools
export const useUserStore = defineStore('user', () => {
    // ref() 创建了一个响应式的 state 属性，类似于组件中的 data
    // 我们从 localStorage 读取初始 token， 防止刷新页面后登录状态丢失
    const token = ref(localStorage.getItem('token') || null)

    // action: 定义一个方法来修改 state
    function setToken(newToken) {
        token.value = newToken
        // 将 token 同时存储到 localStorage
        if (newToken) {
            localStorage.setItem('token', newToken)
        } else {
            localStorage.removeItem('token')
        }
    }

    // action: 定义一个方法来清除 token （用于退出登陆）
    function clearToken() {
        setToken(null)
    }

    // 返回 state 和 actions，以便在组件中使用
    return { token, setToken, clearToken}
})