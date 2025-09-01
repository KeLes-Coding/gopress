<script setup>
import { ref, reactive } from 'vue'
// --- 新增代码 ---
import { useRouter } from 'vue-router'
import { register } from '../api/auth'
import { ElMessage } from 'element-plus'

const registerFormRef = ref(null)
const router = useRouter()
const loading = ref(false)

const registerForm = reactive({
  username: '',
  email: '',
  password: '',
  confirmPassword: '',
})

const validatePass2 = (rule, value, callback) => {
  if (value === '') {
    callback(new Error('请再次输入密码'))
  } else if (value !== registerForm.password) {
    callback(new Error("两次输入的密码不一致!"))
  } else {
    callback()
  }
}

const rules = reactive({
  username: [
    { required: true, message: '请输入用户名', trigger: 'blur' },
    { min: 4, message: '用户名长度不能少于 4 位', trigger: 'blur' },
  ],
  email: [ // <--- 新增 email 规则
    { required: true, message: '请输入邮箱地址', trigger: 'blur' },
    { type: 'email', message: '请输入有效的邮箱地址', trigger: ['blur', 'change'] }
  ],
  password: [
    { required: true, message: '请输入密码', trigger: 'blur' },
    { min: 6, message: '密码长度不能少于 6 位', trigger: 'blur' },
  ],
  confirmPassword: [
    { validator: validatePass2, trigger: 'blur', required: true }
  ],
})

const submitForm = () => {
  registerFormRef.value.validate(async (valid) => {
    if (valid) {
      loading.value = true
      try {
        await register({
          username: registerForm.username,
          password: registerForm.password,
          email: registerForm.email,
        })
        ElMessage.success('注册成功！即将跳转到登录页。')
        // 注册成功后，延时一小会再跳转，让用户能看到提示
        setTimeout(() => {
          router.push('/login')
        }, 1500)
      } catch (error) {
        console.error('注册失败:', error)
      } finally {
        loading.value = false
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
              <span>注册 GoPress 管理员</span>
            </div>
          </template>
          <el-form 
            :model="registerForm" 
            :rules="rules" 
            ref="registerFormRef" 
            label-width="100px"
             @submit.prevent="submitForm"
          >
            <el-form-item label="用户名" prop="username">
              <el-input v-model="registerForm.username"></el-input>
            </el-form-item>
            <el-form-item label="邮箱" prop="email">
              <el-input v-model="registerForm.email"></el-input>
            </el-form-item>
            <el-form-item label="密码" prop="password">
              <el-input type="password" v-model="registerForm.password" show-password></el-input>
            </el-form-item>
            <el-form-item label="确认密码" prop="confirmPassword">
              <el-input type="password" v-model="registerForm.confirmPassword" show-password></el-input>
            </el-form-item>
            <el-form-item>
              <el-button type="primary" @click="submitForm" :loading="loading">注册</el-button>
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
  width: 450px;
}
.card-header {
  text-align: center;
  font-size: 20px;
}
</style>