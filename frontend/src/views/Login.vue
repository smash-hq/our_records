<template>
<div class="login-page">
<div class="login-container">
<div class="login-card">
<div class="login-header">
<img src="/logo.svg" class="logo-icon" alt="logo"/>
<h1 class="login-title">记了么</h1>
<p class="login-subtitle">记录生活中的每一个美好瞬间</p>
</div>

<el-form :model="loginForm" class="login-form" @submit.prevent="handleLogin">
<el-form-item>
<el-input
v-model="loginForm.username"
placeholder="请输入用户名"
size="large"
class="username-input"
maxlength="20"

>
<template #prefix>
<el-icon><User/></el-icon>
</template>
</el-input>
</el-form-item>

<el-form-item>
<el-button
type="primary"
@click="handleLogin"
:loading="loading"
size="large"
class="login-btn"
>
<el-icon><Unlock/></el-icon>
{{loading?"登录中...":"登录"}}
</el-button>
</el-form-item>
</el-form>
</div>
</div>
</div>
</template>

<script setup>
import {ref} from "vue"
import {User,Unlock} from "@element-plus/icons-vue"
import {useRouter} from "vue-router"

const router=useRouter()
const loginForm=ref({username:""})
const loading=ref(false)

const handleLogin=()=>{
  const username=loginForm.value.username.trim().toLowerCase()
  
  if(!username){
    return
  }
  
  if(username!=="huangqi"&&username!=="zhongyanling"){
    return
  }
  
  loading.value=true
  
  localStorage.setItem("lastLoginTime",new Date().toISOString())
  localStorage.setItem("lastLoginUser",username)
  localStorage.setItem("isLoggedIn","true")
  
  setTimeout(()=>{
    loading.value=false
    router.push("/timeline")
  },300)
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
.login-btn .el-icon{
  margin-right:6px;
  font-size:18px;
}
</style>