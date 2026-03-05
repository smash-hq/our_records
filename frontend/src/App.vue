<template>
<div id="app">
<el-container v-if="showLayout">
<el-header>
<div class="header-content">
<h1 class="logo"><img src="/logo.svg" class="logo-icon" alt="logo"/><span class="logo-text">记了么</span></h1>
<div class="nav-right">
<el-menu mode="horizontal" :ellipsis="false" :default-active="defaultRoute" router class="nav-menu">
<el-menu-item index="/upload" class="nav-item"><el-icon><Edit/></el-icon><span>写记录</span></el-menu-item>
<el-menu-item index="/timeline" class="nav-item"><el-icon><Clock/></el-icon><span>时光轴</span></el-menu-item>
<el-menu-item index="/records" class="nav-item"><el-icon><Document/></el-icon><span>记录列表</span></el-menu-item>
</el-menu>
<el-button class="logout-btn" @click="handleLogout">
<el-icon><SwitchButton/></el-icon>退出登录
</el-button>
</div>
<el-button class="menu-btn" circle @click="drawerVisible=true" v-if="isMobile">
<el-icon><Menu/></el-icon>
</el-button>
</div>
</el-header>
<el-main>
<router-view/>
</el-main>
</el-container>
<router-view v-if="!showLayout"/>
<el-drawer v-model="drawerVisible" direction="rtl" size="200px">
<el-menu mode="vertical" :default-active="defaultRoute" router @select="drawerVisible=false">
<el-menu-item index="/upload"><el-icon><Edit/></el-icon>写记录</el-menu-item>
<el-menu-item index="/timeline"><el-icon><Clock/></el-icon>时光轴</el-menu-item>
<el-menu-item index="/records"><el-icon><Document/></el-icon>记录列表</el-menu-item>
</el-menu>
</el-drawer>
</div>
</template>

<script setup>
import {ref,computed,onMounted,onUnmounted} from "vue"
import {useRoute,useRouter} from "vue-router"
import {Star,Clock,Document,Edit,Menu,SwitchButton} from "@element-plus/icons-vue"
import {ElMessageBox} from "element-plus"

const route=useRoute()
const router=useRouter()
const drawerVisible=ref(false)
const isMobile=ref(false)
const showLayout=computed(()=>route.path!=="/login")

const defaultRoute=computed(()=>{
  const path=route.path
  if(path==="/upload")return"/upload"
  if(path==="/records")return"/records"
  return"/timeline"
})

const handleLogout=()=>{
  ElMessageBox.confirm("确定要退出登录吗？","提示",{
    confirmButtonText:"确定",
    cancelButtonText:"取消",
    type:"warning"
  }).then(()=>{
    localStorage.removeItem("isLoggedIn")
    router.push("/login")
  }).catch(()=>{})
}

const checkMobile=()=>{isMobile.value=window.innerWidth<=768}
checkMobile()
onMounted(()=>window.addEventListener("resize",checkMobile))
onUnmounted(()=>window.removeEventListener("resize",checkMobile))
</script>

<style>
*{margin:0;padding:0;box-sizing:border-box}
html,body{height:100%;overflow-x:hidden}
#app{font-family:-apple-system,BlinkMacSystemFont,"Segoe UI",Roboto,"Helvetica Neue",Arial,sans-serif;min-height:100vh}
.el-header{background:linear-gradient(135deg,#667eea 0%,#764ba2 100%);color:#fff;display:flex;align-items:center;box-shadow:0 2px 12px rgba(0,0,0,0.1);height:56px;padding:0 16px}
.header-content{display:flex;align-items:center;width:100%;justify-content:space-between}
.logo{display:flex;align-items:center;gap:8px;font-size:20px}.logo-icon{width:28px;height:28px}
.logo-text{display:inline-block}
.nav-right{display:flex;align-items:center;flex:1;justify-content:flex-end;gap:20px}
.nav-menu{background:transparent;border:none}
.nav-menu .nav-item{color:rgba(255,255,255,0.9);border:none;padding:0 16px;border-radius:20px;margin:0 2px}
.nav-menu .nav-item:hover,.nav-menu .nav-item.is-active{background:rgba(255,255,255,0.2)!important;color:#fff}
.logout-btn{
  background:rgba(255,255,255,0.25);
  color:#fff;
  border:1px solid rgba(255,255,255,0.4);
  padding:10px 20px;
  border-radius:20px;
  font-size:14px;
  font-weight:500;
  transition:all 0.3s;
}
.logout-btn:hover{
  background:rgba(255,255,255,0.35);
  border-color:rgba(255,255,255,0.6);
  transform:translateY(-1px);
  box-shadow:0 4px 12px rgba(0,0,0,0.2);
}
.logout-btn .el-icon{margin-right:6px;font-size:16px}
.el-main{padding:0;background:#f5f7fa;min-height:calc(100vh - 56px)}
.menu-btn{background:rgba(255,255,255,0.2);border:none;color:#fff}
.menu-btn:hover{background:rgba(255,255,255,0.3)}
@media(max-width:768px){
.el-header{height:50px;padding:0 12px}
.logo{font-size:18px}
.logo-text{display:none}
.nav-right{display:none}
}
</style>