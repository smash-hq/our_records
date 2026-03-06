<template>
<div class="profile-page">
<div class="profile-container">
<el-card class="profile-card">
<div class="profile-header">
<div class="avatar-section">
<el-avatar :size="100" :src="user?.avatar" class="profile-avatar">
<template #default>{{ user?.nickname?.charAt(0) || user?.username?.charAt(0) || 'U' }}</template>
</el-avatar>
<el-upload
  class="avatar-uploader"
  :show-file-list="false"
  :before-upload="beforeAvatarUpload"
  :http-request="handleAvatarUpload"
  accept="image/png,image/jpeg,image/gif,image/webp"
>
<el-button type="primary" size="small" class="change-avatar-btn">
<el-icon><Camera/></el-icon>更换头像
</el-button>
</el-upload>
</div>
<h2 class="profile-name">{{ user?.nickname || user?.username }}</h2>
<p class="profile-username">@{{ user?.username }}</p>
</div>

<el-divider/>

<div class="profile-info">
<div class="info-item">
<el-icon class="info-icon"><User/></el-icon>
<div class="info-content">
<span class="info-label">用户名</span>
<span class="info-value">{{ user?.username }}</span>
</div>
</div>
<div class="info-item">
<el-icon class="info-icon"><Star/></el-icon>
<div class="info-content">
<span class="info-label">昵称</span>
<span class="info-value">{{ user?.nickname || '-' }}</span>
</div>
</div>
<div class="info-item">
<el-icon class="info-icon"><Message/></el-icon>
<div class="info-content">
<span class="info-label">邮箱</span>
<span class="info-value">{{ user?.email || '-' }}</span>
</div>
</div>
<div class="info-item">
<el-icon class="info-icon"><Clock/></el-icon>
<div class="info-content">
<span class="info-label">注册时间</span>
<span class="info-value">{{ formatDate(user?.created_at) }}</span>
</div>
</div>
</div>

<el-divider/>

<div class="profile-actions">
<el-button type="primary" @click="showPasswordDialog=true" class="action-btn">
<el-icon><Lock/></el-icon>修改密码
</el-button>
</div>
</el-card>

<!-- 我的评论 -->
<el-card class="profile-card comments-card">
<div class="card-header">
<h3><el-icon><ChatDotRound/></el-icon>我的评论</h3>
</div>
<div class="my-comments" v-loading="commentsLoading">
<div v-for="comment in myComments" :key="comment.id" class="comment-item">
<div class="comment-content">
<p class="comment-text">{{comment.content}}</p>
<div class="comment-meta">
<span class="comment-record">发布于《{{comment.record_title}}》</span>
<span class="comment-date">{{formatTime(comment.created_at)}}</span>
<span class="comment-likes"><el-icon><Star/></el-icon>{{comment.like_count}}</span>
</div>
</div>
</div>
<el-empty v-if="!commentsLoading && myComments.length===0" description="还没有评论"/>
<div class="pagination" v-if="myComments.length > 0">
<el-pagination
  layout="prev, pager, next"
  :total="commentsTotal"
  :page-size="commentsPageSize"
  :current-page="commentsPage"
  @current-change="loadMyComments"
/>
</div>
</div>
</el-card>
</div>

<!-- 修改密码对话框 -->
<el-dialog v-model="showPasswordDialog" title="修改密码" width="400px">
<el-form :model="passwordForm" :rules="passwordRules" ref="passwordFormRef" label-width="80px">
<el-form-item label="原密码" prop="old_password">
<el-input
v-model="passwordForm.old_password"
type="password"
placeholder="请输入原密码"
show-password
/>
</el-form-item>
<el-form-item label="新密码" prop="new_password">
<el-input
v-model="passwordForm.new_password"
type="password"
placeholder="请输入新密码（至少 6 位）"
show-password
/>
</el-form-item>
</el-form>
<template #footer>
<div class="dialog-footer">
<el-button @click="showPasswordDialog=false">取消</el-button>
<el-button type="primary" @click="handleChangePassword" :loading="changingPassword">确定</el-button>
</div>
</template>
</el-dialog>
</div>
</template>

<script setup>
import { ref, reactive, onMounted } from "vue"
import { User, Star, Message, Clock, Lock, Camera, ChatDotRound } from "@element-plus/icons-vue"
import { useRouter } from "vue-router"
import { ElMessage } from "element-plus"
import { getCurrentUser, changePassword, uploadAvatar } from "@/api/auth"
import { getMyComments } from "@/api/comment"

const router = useRouter()
const user = ref(null)
const showPasswordDialog = ref(false)
const changingPassword = ref(false)
const passwordFormRef = ref(null)
const uploadingAvatar = ref(false)

const passwordForm = reactive({
  old_password: "",
  new_password: ""
})

const passwordRules = {
  old_password: [
    { required: true, message: "请输入原密码", trigger: "blur" }
  ],
  new_password: [
    { required: true, message: "请输入新密码", trigger: "blur" },
    { min: 6, max: 50, message: "密码长度在 6-50 个字符", trigger: "blur" }
  ]
}

// 我的评论相关
const myComments = ref([])
const commentsLoading = ref(false)
const commentsTotal = ref(0)
const commentsPage = ref(1)
const commentsPageSize = ref(10)

const formatDate = (date) => {
  if (!date) return "-"
  return new Date(date).toLocaleString("zh-CN")
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

const loadUser = async () => {
  try {
    const res = await getCurrentUser()
    user.value = res
    localStorage.setItem("user", JSON.stringify(res))
  } catch (error) {
    console.error("加载用户信息失败", error)
  }
}

const loadMyComments = async (page = 1) => {
  commentsLoading.value = true
  try {
    const res = await getMyComments({ page, page_size: commentsPageSize.value })
    myComments.value = res.comments || []
    commentsTotal.value = res.total || 0
    commentsPage.value = page
  } catch (error) {
    console.error(error)
  } finally {
    commentsLoading.value = false
  }
}

const beforeAvatarUpload = (file) => {
  const validTypes = ['image/jpeg', 'image/png', 'image/gif', 'image/webp']
  if (!validTypes.includes(file.type)) {
    ElMessage.error('只支持 JPG/PNG/GIF/WEBP 格式的图片')
    return false
  }
  if (file.size > 5 * 1024 * 1024) {
    ElMessage.error('图片大小不能超过 5MB')
    return false
  }
  return true
}

const handleAvatarUpload = async (options) => {
  uploadingAvatar.value = true
  try {
    const formData = new FormData()
    formData.append('file', options.file)
    const res = await uploadAvatar(formData)
    ElMessage.success('头像上传成功')
    user.value = res.user
    localStorage.setItem("user", JSON.stringify(res.user))
  } catch (error) {
    console.error('头像上传失败:', error)
    ElMessage.error('头像上传失败：' + (error.message || '未知错误'))
  } finally {
    uploadingAvatar.value = false
  }
}

const handleChangePassword = async () => {
  if (!passwordFormRef.value) return

  await passwordFormRef.value.validate(async (valid) => {
    if (!valid) return

    changingPassword.value = true
    try {
      await changePassword({
        old_password: passwordForm.old_password,
        new_password: passwordForm.new_password
      })
      ElMessage.success("密码修改成功，请重新登录")
      showPasswordDialog.value = false
      passwordForm.old_password = ""
      passwordForm.new_password = ""
      localStorage.removeItem("token")
      localStorage.removeItem("user")
      setTimeout(() => {
        router.push("/login")
      }, 1000)
    } catch (error) {
      console.error(error)
    } finally {
      changingPassword.value = false
    }
  })
}

onMounted(() => {
  loadUser()
  loadMyComments()
})
</script>

<style scoped>
.profile-page{
  min-height:calc(100vh - 56px);
  padding:24px 16px;
  background:#f5f7fa;
}
.profile-container{
  max-width:800px;
  margin:0 auto;
}
.profile-card{
  border-radius:16px;
  margin-bottom:20px;
}
.profile-header{
  text-align:center;
  padding:20px 0;
}
.avatar-section{
  display:flex;
  flex-direction:column;
  align-items:center;
  gap:12px;
  margin-bottom:16px;
}
.profile-avatar{
  box-shadow:0 4px 16px rgba(102,126,234,0.3);
}
.change-avatar-btn{
  padding:6px 16px;
  border-radius:16px;
  font-size:13px;
}
.change-avatar-btn .el-icon{
  margin-right:4px;
  font-size:14px;
}
.profile-name{
  font-size:24px;
  color:#303133;
  margin-bottom:8px;
  font-weight:600;
}
.profile-username{
  font-size:14px;
  color:#909399;
}
.profile-info{
  padding:8px 0;
}
.info-item{
  display:flex;
  align-items:center;
  padding:16px 0;
  border-bottom:1px solid #f5f5f5;
}
.info-item:last-child{
  border-bottom:none;
}
.info-icon{
  font-size:20px;
  color:#667eea;
  margin-right:16px;
  flex-shrink:0;
}
.info-content{
  flex:1;
  display:flex;
  justify-content:space-between;
  align-items:center;
}
.info-label{
  font-size:14px;
  color:#909399;
}
.info-value{
  font-size:14px;
  color:#303133;
  font-weight:500;
}
.profile-actions{
  padding:20px 0;
  text-align:center;
}
.action-btn{
  padding:12px 32px;
  border-radius:12px;
  font-size:15px;
}
.action-btn .el-icon{
  margin-right:6px;
  font-size:18px;
}
.card-header{
  display:flex;
  align-items:center;
  justify-content:space-between;
  padding:16px 20px;
  border-bottom:1px solid #e4e7ed;
}
.card-header h3{
  font-size:16px;
  color:#303133;
  margin:0;
  display:flex;
  align-items:center;
  gap:8px;
}
.my-comments{
  padding:0;
}
.comment-item{
  padding:16px 20px;
  border-bottom:1px solid #f5f5f5;
  transition:background 0.2s;
}
.comment-item:hover{
  background:#f8f9ff;
}
.comment-item:last-child{
  border-bottom:none;
}
.comment-text{
  font-size:14px;
  color:#303133;
  line-height:1.6;
  margin-bottom:10px;
  white-space:pre-wrap;
}
.comment-meta{
  display:flex;
  align-items:center;
  gap:12px;
  font-size:12px;
  color:#909399;
  flex-wrap:wrap;
}
.comment-record{
  color:#667eea;
}
.comment-likes{
  display:flex;
  align-items:center;
  gap:4px;
}
.comment-likes .el-icon{
  font-size:14px;
}
.pagination{
  padding:16px 20px;
  display:flex;
  justify-content:center;
}
.dialog-footer{
  display:flex;
  justify-content:flex-end;
  gap:10px;
}
</style>
