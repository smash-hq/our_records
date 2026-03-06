<template>
<div class="login-page">
<div class="login-container">
<div class="login-card">
<div class="login-header">
<img src="/logo.svg" class="logo-icon" alt="logo"/>
<h1 class="login-title">记了么</h1>
<p class="login-subtitle">{{ isRegister ? '创建新账号' : '记录生活中的每一个美好瞬间' }}</p>
</div>

<el-form :model="form" :rules="rules" ref="formRef" class="login-form" @submit.prevent="handleSubmit">
<el-form-item prop="username">
<el-input
v-model="form.username"
placeholder="用户名"
size="large"
maxlength="50"
>
<template #prefix>
<el-icon><User/></el-icon>
</template>
</el-input>
</el-form-item>

<el-form-item prop="password" v-if="isRegister">
<el-input
v-model="form.password"
type="password"
placeholder="密码（至少 6 位）"
size="large"
maxlength="50"
show-password
>
<template #prefix>
<el-icon><Lock/></el-icon>
</template>
</el-input>
</el-form-item>

<el-form-item prop="password" v-else>
<el-input
v-model="form.password"
type="password"
placeholder="密码"
size="large"
maxlength="50"
show-password
>
<template #prefix>
<el-icon><Lock/></el-icon>
</template>
</el-input>
</el-form-item>

<el-form-item prop="email" v-if="isRegister">
<el-input
v-model="form.email"
placeholder="邮箱（选填）"
size="large"
maxlength="100"
>
<template #prefix>
<el-icon><Message/></el-icon>
</template>
</el-input>
</el-form-item>

<el-form-item>
<el-button
type="primary"
@click="handleSubmit"
:loading="loading"
size="large"
class="login-btn"
>
{{ loading ? (isRegister ? '注册中...' : '登录中...') : (isRegister ? '注册' : '登录') }}
</el-button>
</el-form-item>
</el-form>

<div class="switch-form">
<template v-if="isRegister">
已有账号？
<el-link type="primary" @click="isRegister = false">去登录</el-link>
</template>
<template v-else>
没有账号？
<el-link type="primary" @click="isRegister = true">去注册</el-link>
</template>
</div>
</div>
</div>
</div>
</template>

<script setup>
import { ref, reactive, computed } from "vue"
import { User, Lock, Message } from "@element-plus/icons-vue"
import { useRouter } from "vue-router"
import { ElMessage } from "element-plus"
import { login, register } from "@/api/auth"

const router = useRouter()
const formRef = ref(null)
const isRegister = ref(false)
const loading = ref(false)

const form = reactive({
  username: "",
  password: "",
  email: ""
})

const rules = computed(() => {
  const baseRules = {
    username: [
      { required: true, message: "请输入用户名", trigger: "blur" },
      { min: 3, max: 50, message: "用户名长度在 3-50 个字符", trigger: "blur" }
    ],
    password: [
      { required: true, message: "请输入密码", trigger: "blur" },
      { min: 6, max: 50, message: "密码长度在 6-50 个字符", trigger: "blur" }
    ]
  }
  if (isRegister.value) {
    baseRules.email = [
      { type: "email", message: "请输入有效的邮箱地址", trigger: "blur" }
    ]
  }
  return baseRules
})

const handleSubmit = async () => {
  if (!formRef.value) return
  
  await formRef.value.validate(async (valid) => {
    if (!valid) return
    
    loading.value = true
    try {
      if (isRegister.value) {
        // 注册
        const res = await register({
          username: form.username.trim(),
          password: form.password,
          email: form.email.trim() || undefined
        })
        ElMessage.success("注册成功，请登录")
        isRegister.value = false
        form.password = ""
      } else {
        // 登录
        const res = await login({
          username: form.username.trim(),
          password: form.password
        })
        // 保存 token 和用户信息
        localStorage.setItem("token", res.token)
        localStorage.setItem("user", JSON.stringify(res.user))
        ElMessage.success("登录成功")
        router.push("/timeline")
      }
    } catch (error) {
      console.error(error)
    } finally {
      loading.value = false
    }
  })
}
</script>

<style scoped>
.login-page{
  min-height:100vh;
  background:linear-gradient(135deg,#667eea 0%,#764ba2 100%);
  display:flex;
  align-items:center;
  justify-content:center;
  padding:20px;
}
.login-container{
  width:100%;
  max-width:420px;
}
.login-card{
  background:rgba(255,255,255,0.95);
  backdrop-filter:blur(10px);
  border-radius:24px;
  padding:48px 40px;
  box-shadow:0 20px 60px rgba(0,0,0,0.3);
}
.login-header{
  text-align:center;
  margin-bottom:40px;
}
.logo-icon{
  width:80px;
  height:80px;
  margin-bottom:16px;
}
.login-title{
  font-size:32px;
  color:#1a1a2e;
  margin-bottom:8px;
  font-weight:600;
}
.login-subtitle{
  font-size:14px;
  color:#666;
}
.login-form{
  margin-bottom:24px;
}
.login-form :deep(.el-form-item){
  margin-bottom:20px;
}
.login-form :deep(.el-input__wrapper){
  border-radius:12px;
  padding:12px 16px;
  box-shadow:0 2px 12px rgba(0,0,0,0.05);
  transition:all 0.3s;
}
.login-form :deep(.el-input__wrapper):hover{
  box-shadow:0 4px 16px rgba(0,0,0,0.1);
}
.login-form :deep(.el-input__wrapper.is-focus){
  box-shadow:0 4px 20px rgba(102,126,234,0.3);
}
.login-btn{
  width:100%;
  padding:14px;
  border-radius:12px;
  font-size:16px;
  font-weight:500;
  background:linear-gradient(135deg,#667eea,#764ba2);
  border:none;
  transition:all 0.3s;
}
.login-btn:hover{
  transform:translateY(-2px);
  box-shadow:0 8px 20px rgba(102,126,234,0.4);
}
.login-btn:active{
  transform:translateY(0);
}
.switch-form{
  text-align:center;
  font-size:14px;
  color:#666;
}
.switch-form .el-link{
  font-weight:500;
}
</style>
