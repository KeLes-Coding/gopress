<script setup>
import { ref, reactive } from 'vue'
// --- 新增代码 Start ---
import { useRouter } from 'vue-router'
import { useUserStore } from '../store/user'
import { login } from '../api/auth'
import { ElMessage } from 'element-plus'
// --- 新增代码 End ---

const loginFormRef = ref(null)
// --- 新增代码 Start ---
const router = useRouter() // 获取 router 实例，用于页面跳转
const userStore = useUserStore() // 获取 user store 实例
const loading = ref(false) // 添加一个 loading 状态，防止重复提交
// --- 新增代码 End ---

const loginForm = reactive({
  username: '',
  password: '',
})

const rules = reactive({
  username: [{ required: true, message: '请输入用户名', trigger: 'blur' }],
  password: [{ required: true, message: '请输入密码', trigger: 'blur' }],
})

const submitForm = () => {
  loginFormRef.value.validate(async (valid) => { // 将函数改为 async
    if (valid) {
      // --- 新增业务逻辑 ---
      loading.value = true
      try {
        const response = await login({
          username: loginForm.username,
          password: loginForm.password,
        })
        // 假设后端返回的数据结构是 { code: 200, data: { token: '...' } }
        const token = response.data.data.token
        // 调用 store 中的 action 来保存 token
        userStore.setToken(token)
        // 弹出成功提示
        ElMessage.success('登录成功！')
        // 跳转到后台仪表盘页面
        await router.push('/admin/dashboard')

      } catch (error) {
        // API 调用失败的错误已在响应拦截器中统一处理，这里可以留空
        console.error('登录失败:', error)
      } finally {
        loading.value = false // 无论成功失败，都结束 loading 状态
      }
    } else {
      console.log('表单校验失败')
      return false
    }
  })
}
</script>

<template>
  <div class="form-container">
    <el-card class="box-card">
      <template #header>
        <div class="card-header">
          <span>GoPress 后台登录</span>
        </div>
      </template>
      <el-form 
        :model="loginForm" 
        :rules="rules" 
        ref="loginFormRef" 
        label-width="80px"
        @submit.prevent="submitForm"
      >
        <el-form-item label="用户名" prop="username">
          <el-input v-model="loginForm.username"></el-input>
        </el-form-item>
        <el-form-item label="密码" prop="password">
          <el-input type="password" v-model="loginForm.password" show-password></el-input>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="submitForm" :loading="loading">登录</el-button>
        </el-form-item>
      </el-form>
    </el-card>
  </div>
</template>

<style scoped>
/* 样式保持不变 */
.form-container {
  display: flex;
  justify-content: center;
  align-items: center;
  height: 100vh;
  background-color: #f0f2f5;
}
.box-card {
  width: 400px;
}
.card-header {
  text-align: center;
  font-size: 20px;
}
</style>