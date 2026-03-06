<template>
<div id="app">
<!-- 主布局 -->
<el-container v-if="showLayout">
<el-header>
<div class="header-content">
<h1 class="logo"><img src="/logo.svg" class="logo-icon" alt="logo"/><span class="logo-text">记了么</span></h1>
<div class="nav-right">
<el-menu mode="horizontal" :ellipsis="false" :default-active="defaultRoute" router class="nav-menu">
<el-menu-item index="/upload" class="nav-item"><el-icon><Edit/></el-icon><span>写记录</span></el-menu-item>
<el-menu-item index="/timeline" class="nav-item"><el-icon><Clock/></el-icon><span>时光轴</span></el-menu-item>
<el-menu-item index="/records" class="nav-item"><el-icon><Document/></el-icon><span>记录列表</span></el-menu-item>
<el-menu-item index="/groups" class="nav-item"><el-icon><ChatDotRound/></el-icon><span>群组</span></el-menu-item>
</el-menu>
<div class="user-info">
<el-badge :value="unreadCount" :hidden="unreadCount===0" class="notification-badge">
<div class="notification-icon" @click="handleBellClick">
<el-icon :size="22"><Bell/></el-icon>
</div>
</el-badge>
<el-dropdown trigger="click" @command="handleUserCommand">
<div class="user-avatar">
<el-avatar :size="36" :src="user?.avatar" :key="user?.id || 'guest'">
<template #default>{{ user?.nickname?.charAt(0) || user?.username?.charAt(0) || 'U' }}</template>
</el-avatar>
<span class="user-name">{{ user?.nickname || user?.username || '用户' }}</span>
<el-icon class="el-icon--right"><ArrowDown/></el-icon>
</div>
<template #dropdown>
<el-dropdown-menu>
<el-dropdown-item command="profile"><el-icon><User/></el-icon>个人中心</el-dropdown-item>
<el-dropdown-item command="notifications" divided><el-icon><Bell/></el-icon>我的通知</el-dropdown-item>
<el-dropdown-item command="logout" divided><el-icon><SwitchButton/></el-icon>退出登录</el-dropdown-item>
</el-dropdown-menu>
</template>
</el-dropdown>
</div>
</div>
<el-button class="menu-btn" circle @click="drawerVisible=true" v-if="isMobile">
<el-icon><Menu/></el-icon>
</el-button>
</div>
</el-header>
<el-main>
<router-view/>
</el-main>
<el-drawer v-model="drawerVisible" direction="rtl" size="200px">
<el-menu mode="vertical" :default-active="defaultRoute" router @select="drawerVisible=false">
<el-menu-item index="/upload"><el-icon><Edit/></el-icon>写记录</el-menu-item>
<el-menu-item index="/timeline"><el-icon><Clock/></el-icon>时光轴</el-menu-item>
<el-menu-item index="/records"><el-icon><Document/></el-icon>记录列表</el-menu-item>
<el-menu-item index="/groups"><el-icon><ChatDotRound/></el-icon>群组</el-menu-item>
</el-menu>
</el-drawer>
</el-container>

<!-- 登录页面等 -->
<router-view v-if="!showLayout"/>

<!-- 通知弹窗 -->
<el-dialog v-model="notificationDialogVisible" title="我的通知" width="500px" class="notification-dialog">
<div class="notification-list" v-loading="notificationsLoading">
<div v-if="notifications.length > 0" class="notification-header">
<span class="notification-count">{{ unreadCount }} 条未读</span>
<el-button link type="primary" size="small" @click="markAllAsRead" v-if="unreadCount > 0">全部标记为已读</el-button>
</div>
<el-empty v-if="!notificationsLoading && notifications.length === 0" description="暂无通知"/>
<div v-for="notice in notifications" :key="notice.id" class="notification-item" :class="{'unread': !notice.is_read}" @click="goToNotification(notice)">
<div class="notification-avatar">
<el-avatar :size="40" :src="notice.from_avatar">
{{ notice.from_nickname?.charAt(0) || notice.from_username?.charAt(0) || 'U' }}
</el-avatar>
</div>
<div class="notification-content">
<div class="notification-info">
<span class="notification-from">{{ notice.from_nickname || notice.from_username }}</span>
<span class="notification-time">{{ formatTime(notice.created_at) }}</span>
</div>
<div class="notification-text">{{ notice.content }}</div>
<div class="notification-record">《{{ notice.record_title }}》</div>
</div>
<div class="notification-type-badge comment" v-if="notice.type === 'comment'">评论记录</div>
<div class="notification-type-badge reply" v-else-if="notice.type === 'reply'">回复评论</div>
</div>
</div>
<template #footer>
<el-button @click="notificationDialogVisible = false">关闭</el-button>
</template>
</el-dialog>
</div>
</template>

<script setup>
import {ref,computed,onMounted,onUnmounted,watch} from "vue"
import {useRoute,useRouter} from "vue-router"
import {Star,Clock,Document,Edit,Menu,SwitchButton,ArrowDown,User,ChatDotRound,Bell} from "@element-plus/icons-vue"
import {ElMessageBox,ElMessage} from "element-plus"
import { getUnreadNotifications, getNotifications, markAllNotificationsAsRead } from "@/api/notification"

const route=useRoute()
const router=useRouter()
const drawerVisible=ref(false)
const isMobile=ref(false)
const user=ref(null)
const showLayout=computed(()=>route.path!=="/login")

// 通知相关
const unreadCount = ref(0)
const notifications = ref([])
const notificationsLoading = ref(false)
const notificationDialogVisible = ref(false)
const pollTimer = ref(null)

const defaultRoute=computed(()=>{
  const path=route.path
  if(path==="/upload")return"/upload"
  if(path==="/records")return"/records"
  return"/timeline"
})

// 加载用户信息
const loadUser=()=>{
  const userStr=localStorage.getItem("user")
  if(userStr){
    try{
      const userData=JSON.parse(userStr)
      user.value=userData
    }catch(e){
      console.error("解析用户信息失败")
      user.value=null
    }
  }else{
    user.value=null
  }
}

// 监听路由变化，在每次路由变化后重新加载用户信息
watch(()=>route.path, (newPath) => {
  // 在非登录页面时加载用户信息
  if (newPath !== '/login') {
    loadUser()
  }
})

// 监听用户变化，启动/停止轮询
watch(()=>user.value?.id, (newId, oldId) => {
  console.log('[通知] 用户 ID 变化:', newId, oldId)
  if (newId && newId !== oldId) {
    // 刚登录，启动轮询
    console.log('[通知] 启动轮询')
    startPolling()
    loadUnreadCount()
  } else if (!newId && oldId) {
    // 刚登出，停止轮询
    console.log('[通知] 停止轮询')
    stopPolling()
    unreadCount.value = 0
    notifications.value = []
  }
})

// 通知相关函数
const loadUnreadCount = async () => {
  try {
    const res = await getUnreadNotifications()
    unreadCount.value = res.unread_count || 0
    console.log('[通知] 未读数量:', unreadCount.value)
  } catch (error) {
    console.error("加载未读通知失败", error)
  }
}

const loadNotifications = async () => {
  notificationsLoading.value = true
  try {
    const res = await getNotifications({ page: 1, page_size: 50 })
    notifications.value = res.notifications || []
    console.log('[通知] 加载通知列表:', notifications.value.length, '条')
  } catch (error) {
    console.error("加载通知列表失败", error)
  } finally {
    notificationsLoading.value = false
  }
}

const startPolling = () => {
  console.log('[通知] 启动轮询')
  // 每 30 秒轮询一次
  pollTimer.value = setInterval(() => {
    loadUnreadCount()
  }, 30000)
  // 立即执行一次
  loadUnreadCount()
}

const stopPolling = () => {
  if (pollTimer.value) {
    clearInterval(pollTimer.value)
    pollTimer.value = null
  }
}

const handleNotificationCommand = (command) => {
  notificationDialogVisible.value = true
  loadNotifications()
  console.log('[通知] 打开通知列表')
}

// 点击铃铛直接打开通知列表
const handleBellClick = () => {
  notificationDialogVisible.value = true
  loadNotifications()
  console.log('[通知] 点击铃铛打开通知列表')
}

const markAllAsRead = async () => {
  try {
    await markAllNotificationsAsRead()
    await loadUnreadCount()
    await loadNotifications()
    ElMessage.success("已全部标记为已读")
  } catch (error) {
    console.error("标记已读失败", error)
  }
}

const formatTime = (date) => {
  if (!date) return ""
  const d = new Date(date)
  const now = new Date()
  const diff = now - d
  const minutes = Math.floor(diff / (1000 * 60))
  const hours = Math.floor(diff / (1000 * 60 * 60))
  const days = Math.floor(diff / (1000 * 60 * 60 * 24))

  if (minutes < 1) return "刚刚"
  if (minutes < 60) return `${minutes}分钟前`
  if (hours < 24) return `${hours}小时前`
  if (days < 30) return `${days}天前`
  return d.toLocaleDateString("zh-CN")
}

const goToNotification = (notice) => {
  console.log('[通知] 跳转到记录:', notice.record_id)
  // 跳转到记录列表页，并传递记录 ID
  router.push({ path: '/records', query: { view: notice.record_id } })
  notificationDialogVisible.value = false
}

const handleUserCommand=(command)=>{
  if(command==="logout"){
    ElMessageBox.confirm("确定要退出登录吗？","提示",{
      confirmButtonText:"确定",
      cancelButtonText:"取消",
      type:"warning"
    }).then(()=>{
      localStorage.removeItem("token")
      localStorage.removeItem("user")
      stopPolling()
      router.push("/login")
      ElMessage.success("已退出登录")
    }).catch(()=>{})
  }else if(command==="profile"){
    router.push("/profile")
  }else if(command==="notifications"){
    handleBellClick()
  }
}

const checkMobile=()=>{isMobile.value=window.innerWidth<=768}
checkMobile()
onMounted(()=>{
  window.addEventListener("resize",checkMobile)
  loadUser()
  // 如果用户已登录，启动轮询
  const userStr=localStorage.getItem("user")
  if(userStr){
    try{
      const userData=JSON.parse(userStr)
      if(userData && userData.id){
        console.log('[通知] 检测到已登录用户，启动轮询')
        startPolling()
        loadUnreadCount()
      }
    }catch(e){
      console.error("解析用户信息失败", e)
    }
  }
  // 监听 storage 变化，实现多标签页同步登出
  window.addEventListener("storage",(e)=>{
    if(e.key==="token"&&!e.newValue){
      user.value=null
      stopPolling()
      if(route.path!=="/login"){
        router.push("/login")
      }
    }
  })
})
onUnmounted(()=>{
  window.removeEventListener("resize",checkMobile)
  stopPolling()
})
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
.user-info{display:flex;align-items:center}
.user-avatar{display:flex;align-items:center;gap:8px;cursor:pointer;padding:8px 12px;border-radius:20px;transition:all 0.3s}
.user-avatar:hover{background:rgba(255,255,255,0.2)}
.user-name{color:#fff;font-size:14px;font-weight:500}
.user-avatar .el-icon--right{color:rgba(255,255,255,0.8);font-size:14px}
.notification-badge{margin-right:16px;cursor:pointer}
.notification-badge .notification-icon{display:flex;align-items:center;justify-content:center;padding:8px;border-radius:50%;transition:all 0.3s;color:rgba(255,255,255,0.8);cursor:pointer}
.notification-badge .notification-icon:hover{background:rgba(255,255,255,0.2);color:#fff}
.notification-dialog .notification-list{max-height:400px;overflow-y:auto;padding:8px}
.notification-dialog .notification-header{display:flex;justify-content:space-between;align-items:center;margin-bottom:16px;padding-bottom:12px;border-bottom:1px solid #f0f0f0}
.notification-dialog .notification-count{font-size:13px;color:#667eea;font-weight:500}
.notification-dialog .notification-item{display:flex;gap:14px;padding:18px;border-radius:12px;cursor:pointer;transition:all 0.3s;margin-bottom:12px;align-items:flex-start;background:#fff;border:1px solid #e8eaed;box-shadow:0 2px 8px rgba(0,0,0,0.04)}
.notification-dialog .notification-item:hover{transform:translateY(-2px);box-shadow:0 6px 16px rgba(102,126,234,0.15);border-color:#667eea}
.notification-dialog .notification-item.unread{background:linear-gradient(135deg, #f8faff 0%, #f0f4ff 100%);border-color:#c9d6ff;box-shadow:0 2px 10px rgba(102,126,234,0.1)}
.notification-dialog .notification-item.unread:hover{transform:translateY(-2px);box-shadow:0 8px 20px rgba(102,126,234,0.2)}
.notification-dialog .notification-avatar{flex-shrink:0}
.notification-dialog .notification-content{flex:1;min-width:0}
.notification-dialog .notification-info{display:flex;align-items:center;gap:8px;margin-bottom:6px;flex-wrap:wrap}
.notification-dialog .notification-from{font-size:14px;color:#303133;font-weight:600}
.notification-dialog .notification-type{font-size:12px;color:#667eea;background:#f0f4ff;padding:2px 8px;border-radius:4px}
.notification-dialog .notification-time{font-size:12px;color:#909399}
.notification-dialog .notification-text{font-size:13px;color:#606266;line-height:1.5;margin-bottom:6px;overflow:hidden;text-overflow:ellipsis;display:-webkit-box;-webkit-line-clamp:2;-webkit-box-orient:vertical}
.notification-dialog .notification-record{font-size:12px;color:#909399;font-style:italic}
.notification-dialog .notification-type-badge{flex-shrink:0;font-size:12px;padding:5px 12px;border-radius:6px;font-weight:500;transition:all 0.2s}
.notification-dialog .notification-type-badge.comment{color:#667eea;background:linear-gradient(135deg, #f0f4ff 0%, #e6eeff 100%);border:1px solid #c9d6ff}
.notification-dialog .notification-type-badge.reply{color:#f59e42;background:linear-gradient(135deg, #fff7ed 0%, #ffedd5 100%);border:1px solid #fed7aa}
.notification-dialog .notification-type-badge:hover{transform:scale(1.05)}
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