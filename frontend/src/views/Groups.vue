<template>
<div class="groups-page">
<div class="groups-header">
<h1><el-icon><User/></el-icon><span>群组管理</span></h1>
<el-button type="primary" @click="showCreateDialog=true"><el-icon><Plus/></el-icon><span>创建群组</span></el-button>
</div>

<div class="groups-container" v-loading="loading">
<div class="groups-grid">
<el-card v-for="group in groups" :key="group.id" class="group-item" shadow="hover">
<div class="group-header">
<div class="group-icon">
<el-avatar v-if="group.avatar" :size="56" :src="group.avatar" class="group-avatar"/>
<el-icon v-else><ChatDotRound/></el-icon>
</div>
<div class="group-info">
<h3 class="group-name">{{group.name}}</h3>
<p class="group-desc">{{group.description || '暂无描述'}}</p>
</div>
</div>
<div class="group-meta">
<div class="meta-item">
<el-icon><User/></el-icon>
<span>{{group.owner_name}}</span>
</div>
<div class="meta-item">
<el-icon><UserFilled/></el-icon>
<span>{{group.member_count}} 人</span>
</div>
</div>
<div class="group-actions">
<el-button size="small" @click="viewGroup(group)"><el-icon><View/></el-icon>查看</el-button>
<template v-if="group.owner_id===currentUserId">
<el-button size="small" @click="showEditDialog=true;editingGroup=group"><el-icon><Edit/></el-icon>编辑</el-button>
<el-button size="small" type="danger" @click="confirmDelete(group.id)"><el-icon><Delete/></el-icon>删除</el-button>
</template>
<template v-else>
<el-button size="small" type="warning" @click="confirmLeave(group.id)"><el-icon><SwitchButton/></el-icon>退出</el-button>
</template>
</div>
</el-card>
</div>
<el-empty v-if="!loading && groups.length===0" description="暂无群组，创建一个吧"/>
</div>

<!-- 创建群组对话框 -->
<el-dialog v-model="showCreateDialog" title="创建群组" width="500px">
<el-form :model="createForm" :rules="createRules" ref="createFormRef" label-width="80px">
<el-form-item label="群组名称" prop="name">
<el-input v-model="createForm.name" placeholder="请输入群组名称" maxlength="100"/>
</el-form-item>
<el-form-item label="群组描述" prop="description">
<el-input v-model="createForm.description" type="textarea" placeholder="请输入群组描述（选填）" maxlength="500" :rows="3"/>
</el-form-item>
</el-form>
<template #footer>
<div class="dialog-footer">
<el-button @click="showCreateDialog=false">取消</el-button>
<el-button type="primary" @click="handleCreate" :loading="creating">创建</el-button>
</div>
</template>
</el-dialog>

<!-- 编辑群组对话框 -->
<el-dialog v-model="showEditDialog" title="编辑群组" width="500px">
<el-form :model="editForm" :rules="createRules" ref="editFormRef" label-width="80px">
<el-form-item label="群组名称" prop="name">
<el-input v-model="editForm.name" placeholder="请输入群组名称" maxlength="100"/>
</el-form-item>
<el-form-item label="群组描述" prop="description">
<el-input v-model="editForm.description" type="textarea" placeholder="请输入群组描述（选填）" maxlength="500" :rows="3"/>
</el-form-item>
</el-form>
<template #footer>
<div class="dialog-footer">
<el-button @click="showEditDialog=false">取消</el-button>
<el-button type="primary" @click="handleUpdate" :loading="updating">保存</el-button>
</div>
</template>
</el-dialog>

<!-- 群组详情对话框 -->
<el-dialog v-model="showViewDialog" title="群组详情" width="700px" :close-on-click-modal="false">
<div v-if="currentGroup" class="view-content">
<div class="view-header">
<div class="view-icon-section">
<el-avatar v-if="currentGroup.avatar" :size="80" :src="currentGroup.avatar" class="group-detail-avatar"/>
<el-icon v-else class="view-icon"><ChatDotRound/></el-icon>
<el-upload
  v-if="isGroupOwner"
  class="group-avatar-uploader"
  :show-file-list="false"
  :before-upload="beforeGroupAvatarUpload"
  :http-request="(options)=>handleGroupAvatarUpload(options,currentGroup.id)"
  accept="image/png,image/jpeg,image/gif,image/webp"
>
<el-button type="primary" size="small" class="change-group-avatar-btn">
<el-icon><Camera/></el-icon>更换头像
</el-button>
</el-upload>
</div>
<div class="view-info">
<h3>{{currentGroup.name}}</h3>
<p>{{currentGroup.description || '暂无描述'}}</p>
</div>
</div>
<el-divider/>
<div class="view-meta">
<div class="meta-row">
<el-icon><User/></el-icon>
<span>群主：{{currentGroup.owner_name}}</span>
</div>
<div class="meta-row">
<el-icon><UserFilled/></el-icon>
<span>成员：{{currentGroup.member_count}} 人</span>
</div>
<div class="meta-row">
<el-icon><Clock/></el-icon>
<span>创建：{{formatDate(currentGroup.created_at)}}</span>
</div>
</div>
<el-divider/>
<div class="members-section">
<div class="members-header">
<div class="header-title">
<h4>群组成员列表</h4>
<span class="member-count">共 {{members.length}} 人</span>
</div>
<div class="header-actions">
<el-button v-if="isGroupOwner" type="primary" size="small" @click="showAddMemberDialog=true">
<el-icon><Plus/></el-icon>添加成员
</el-button>
</div>
</div>
<div class="members-list" v-loading="membersLoading">
<div v-for="(member, index) in members" :key="member.user_id" class="member-item">
<div class="member-index">{{index + 1}}</div>
<el-avatar :size="44" :src="member.avatar" class="member-avatar">
{{member.nickname?.charAt(0) || member.username?.charAt(0) || 'U'}}
</el-avatar>
<div class="member-info">
<div class="member-main">
<div class="member-name">{{member.nickname || member.username}}</div>
<el-tag :type="getRoleTagType(member.role)" size="small" class="role-tag">{{getRoleLabel(member.role)}}</el-tag>
</div>
<div class="member-meta">
<span class="username">@{{member.username}}</span>
<span class="divider">·</span>
<span class="join-time">加入于 {{formatJoinTime(member.joined_at)}}</span>
</div>
</div>
<el-button v-if="isGroupOwner && canRemoveMember(member)" 
  size="small" 
  type="danger" 
  plain
  @click="removeMember(member.user_id)"
  class="remove-btn">
<el-icon><Close/></el-icon>移除
</el-button>
</div>
<el-empty v-if="!membersLoading && members.length===0" description="暂无成员" :image-size="80"/>
</div>
</div>
</div>
<template #footer>
<div class="dialog-footer">
<el-button @click="showViewDialog=false">关闭</el-button>
</div>
</template>
</el-dialog>

<!-- 添加成员对话框 -->
<el-dialog v-model="showAddMemberDialog" title="添加成员" width="500px">
<div class="add-member-content">
<el-input v-model="searchKeyword" placeholder="搜索用户名或昵称" clearable @keyup.enter="handleSearch" :prefix-icon="Search">
<template #append>
<el-button @click="handleSearch"><el-icon><Search/></el-icon>搜索</el-button>
</template>
</el-input>
<div class="search-results" v-loading="searchLoading">
<el-checkbox-group v-model="selectedUsers">
<div v-for="user in searchResults" :key="user.id" class="user-item">
<el-checkbox :label="user.username">{{user.nickname || user.username}}</el-checkbox>
</div>
</el-checkbox-group>
<el-empty v-if="!searchLoading && searchResults.length===0" description="请输入关键词搜索用户" :image-size="60"/>
</div>
</div>
<template #footer>
<div class="dialog-footer">
<el-button @click="showAddMemberDialog=false">取消</el-button>
<el-button type="primary" @click="handleAddMembers" :loading="addingMembers">添加</el-button>
</div>
</template>
</el-dialog>
</div>
</template>

<script setup>
import { ref, reactive, onMounted, computed } from "vue"
import { User, Plus, ChatDotRound, Edit, Delete, View, UserFilled, Search, Close, SwitchButton, Camera } from "@element-plus/icons-vue"
import { ElMessage, ElMessageBox } from "element-plus"
import { getGroups, createGroup, updateGroup, deleteGroup, getGroupMembers, addGroupMember, removeGroupMember, searchUsers, leaveGroup, uploadGroupAvatar } from "@/api/group"

const groups = ref([])
const loading = ref(false)
const showCreateDialog = ref(false)
const showEditDialog = ref(false)
const showViewDialog = ref(false)
const showAddMemberDialog = ref(false)
const creating = ref(false)
const updating = ref(false)
const membersLoading = ref(false)
const searchLoading = ref(false)
const addingMembers = ref(false)
const uploadingAvatar = ref(false)
const editingGroup = ref(null)
const currentGroup = ref(null)
const members = ref([])
const searchKeyword = ref("")
const searchResults = ref([])
const selectedUsers = ref([])

// 获取当前用户 ID
const currentUserId = computed(() => {
  const userStr = localStorage.getItem("user")
  if (userStr) {
    try {
      const user = JSON.parse(userStr)
      return user.id
    } catch (e) {
      return 0
    }
  }
  return 0
})

// 是否是群主
const isGroupOwner = computed(() => {
  if (!currentGroup.value || !currentUserId.value) return false
  return currentGroup.value.owner_id === currentUserId.value
})

const createForm = reactive({ name: "", description: "" })
const editForm = reactive({ name: "", description: "" })

const createRules = {
  name: [{ required: true, message: "请输入群组名称", trigger: "blur" }]
}

const loadGroups = async () => {
  loading.value = true
  try {
    groups.value = await getGroups()
  } catch (error) {
    console.error(error)
  } finally {
    loading.value = false
  }
}

const handleCreate = async () => {
  if (!createForm.name.trim()) {
    ElMessage.warning("请输入群组名称")
    return
  }
  creating.value = true
  try {
    await createGroup(createForm)
    ElMessage.success("创建成功")
    showCreateDialog.value = false
    createForm.name = ""
    createForm.description = ""
    loadGroups()
  } catch (error) {
    console.error(error)
  } finally {
    creating.value = false
  }
}

const handleUpdate = async () => {
  if (!editingGroup.value || !editForm.name.trim()) {
    ElMessage.warning("请输入群组名称")
    return
  }
  updating.value = true
  try {
    await updateGroup(editingGroup.value.id, editForm)
    ElMessage.success("更新成功")
    showEditDialog.value = false
    loadGroups()
  } catch (error) {
    console.error(error)
  } finally {
    updating.value = false
  }
}

const confirmDelete = async (id) => {
  try {
    await ElMessageBox.confirm("确定要删除此群组吗？删除后不可恢复。", "删除确认", {
      confirmButtonText: "确定删除",
      cancelButtonText: "取消",
      type: "warning"
    })
    await deleteGroup(id)
    ElMessage.success("删除成功")
    loadGroups()
  } catch (error) {
    if (error !== "cancel") console.error(error)
  }
}

const viewGroup = async (group) => {
  currentGroup.value = group
  showViewDialog.value = true
  await loadMembers(group.id)
}

const loadMembers = async (groupId) => {
  membersLoading.value = true
  try {
    const res = await getGroupMembers(groupId)
    members.value = res.members || []
    console.log('群组成员数据:', members.value)
  } catch (error) {
    console.error('加载成员失败:', error)
  } finally {
    membersLoading.value = false
  }
}

const handleSearch = async () => {
  if (!searchKeyword.value.trim()) return
  searchLoading.value = true
  try {
    searchResults.value = await searchUsers(searchKeyword.value.trim())
  } catch (error) {
    console.error(error)
  } finally {
    searchLoading.value = false
  }
}

const handleAddMembers = async () => {
  if (selectedUsers.value.length === 0) {
    ElMessage.warning("请选择要添加的用户")
    return
  }
  addingMembers.value = true
  try {
    await addGroupMember(currentGroup.value.id, { usernames: selectedUsers.value })
    ElMessage.success("添加成功")
    showAddMemberDialog.value = false
    selectedUsers.value = []
    searchResults.value = []
    searchKeyword.value = ""
    await loadMembers(currentGroup.value.id)
    currentGroup.value.member_count += selectedUsers.value.length
  } catch (error) {
    console.error(error)
  } finally {
    addingMembers.value = false
  }
}

const removeMember = async (userId) => {
  try {
    await ElMessageBox.confirm("确定要移除此成员吗？", "移除确认", {
      confirmButtonText: "确定",
      cancelButtonText: "取消",
      type: "warning"
    })
    await removeGroupMember(currentGroup.value.id, { user_id: userId })
    ElMessage.success("移除成功")
    await loadMembers(currentGroup.value.id)
    currentGroup.value.member_count = Math.max(0, currentGroup.value.member_count - 1)
  } catch (error) {
    if (error !== "cancel") console.error(error)
  }
}

const canRemoveMember = (member) => {
  return member.role !== "owner" && member.role !== "admin"
}

const getRoleLabel = (role) => {
  const labels = { owner: "群主", admin: "管理员", member: "成员" }
  return labels[role] || "成员"
}

const getRoleTagType = (role) => {
  const types = { owner: "danger", admin: "warning", member: "" }
  return types[role] || ""
}

const formatDate = (date) => {
  if (!date) return "-"
  return new Date(date).toLocaleString("zh-CN")
}

const formatJoinTime = (date) => {
  if (!date) return "未知"
  const d = new Date(date)
  const now = new Date()
  
  // 获取日期部分（忽略时分秒）
  const dDate = new Date(d.getFullYear(), d.getMonth(), d.getDate())
  const nowDate = new Date(now.getFullYear(), now.getMonth(), now.getDate())
  const diff = nowDate - dDate
  const days = Math.floor(diff / (1000 * 60 * 60 * 24))
  
  if (days === 0) return "今天"
  if (days === 1) return "昨天"
  if (days < 30) return `${days}天前`
  const months = Math.floor(days / 30)
  if (months < 12) return `${months}个月前`
  return `${Math.floor(months / 12)}年前`
}

// 编辑时填充表单
const openEditDialog = (group) => {
  editingGroup.value = group
  editForm.name = group.name
  editForm.description = group.description
  showEditDialog.value = true
}

// 退出群组
const confirmLeave = async (id) => {
  try {
    await ElMessageBox.confirm("确定要退出此群组吗？", "退出确认", {
      confirmButtonText: "确定",
      cancelButtonText: "取消",
      type: "warning"
    })
    await leaveGroup(id)
    ElMessage.success("退出成功")
    loadGroups()
  } catch (error) {
    if (error !== "cancel") console.error(error)
  }
}

// 群组头像上传前验证
const beforeGroupAvatarUpload = (file) => {
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

// 处理群组头像上传
const handleGroupAvatarUpload = async (options, groupId) => {
  uploadingAvatar.value = true
  try {
    const formData = new FormData()
    formData.append('file', options.file)
    const res = await uploadGroupAvatar(groupId, formData)
    ElMessage.success('头像上传成功')
    // 更新当前群组信息
    if (currentGroup.value && currentGroup.value.id === groupId) {
      currentGroup.value.avatar = res.avatar
    }
    // 更新群组列表中的头像
    const group = groups.value.find(g => g.id === groupId)
    if (group) {
      group.avatar = res.avatar
    }
  } catch (error) {
    console.error(error)
  } finally {
    uploadingAvatar.value = false
  }
}

onMounted(() => {
  loadGroups()
})
</script>

<style scoped>
.groups-page{min-height:calc(100vh - 56px);padding:24px 16px;background:#f5f7fa}
.groups-header{max-width:1200px;margin:0 auto 24px;display:flex;justify-content:space-between;align-items:center}
.groups-header h1{font-size:24px;color:#303133;display:flex;align-items:center;gap:8px}
.groups-container{max-width:1200px;margin:0 auto}
.groups-grid{display:grid;grid-template-columns:repeat(auto-fill,minmax(320px,1fr));gap:20px}
.group-item{border-radius:16px;transition:all 0.3s;cursor:pointer}
.group-item:hover{transform:translateY(-4px);box-shadow:0 12px 32px rgba(0,0,0,0.1)}
.group-header{display:flex;align-items:center;gap:16px;margin-bottom:16px}
.group-icon{width:56px;height:56px;border-radius:12px;background:linear-gradient(135deg,#667eea,#764ba2);display:flex;align-items:center;justify-content:center;color:#fff;font-size:28px;flex-shrink:0;overflow:hidden}
.group-icon .group-avatar{width:56px;height:56px;border-radius:12px}
.group-info{flex:1;min-width:0}
.group-name{font-size:18px;color:#303133;margin-bottom:4px;font-weight:600;white-space:nowrap;overflow:hidden;text-overflow:ellipsis}
.group-desc{font-size:13px;color:#909399;white-space:nowrap;overflow:hidden;text-overflow:ellipsis}
.group-meta{display:flex;gap:16px;margin-bottom:16px;padding:12px;background:#f9f9f9;border-radius:8px}
.meta-item{display:flex;align-items:center;gap:6px;font-size:13px;color:#606266}
.meta-item .el-icon{color:#667eea}
.group-actions{display:flex;gap:8px}
.group-actions .el-button{flex:1}
.view-content{padding:10px}
.view-header{display:flex;align-items:center;gap:16px;margin-bottom:16px}
.view-icon-section{display:flex;flex-direction:column;align-items:center;gap:12px}
.view-icon{width:80px;height:80px;border-radius:16px;background:linear-gradient(135deg,#667eea,#764ba2);display:flex;align-items:center;justify-content:center;color:#fff;font-size:32px;flex-shrink:0}
.group-detail-avatar{width:80px;height:80px;border-radius:16px;box-shadow:0 4px 16px rgba(102,126,234,0.3)}
.change-group-avatar-btn{padding:6px 16px;border-radius:16px;font-size:13px}
.change-group-avatar-btn .el-icon{margin-right:4px;font-size:14px}
.view-info h3{font-size:20px;color:#303133;margin-bottom:8px}
.view-info p{font-size:14px;color:#909399}
.view-meta{padding:8px 0}
.meta-row{display:flex;align-items:center;gap:8px;padding:8px 0;font-size:14px;color:#606266}
.meta-row .el-icon{color:#667eea;width:20px}
.members-section{padding:8px 0}
.members-header{display:flex;justify-content:space-between;align-items:center;margin-bottom:16px;padding:12px;background:#f8f9ff;border-radius:12px}
.header-title{display:flex;align-items:center;gap:12px}
.header-title h4{font-size:15px;color:#303133;margin:0}
.member-count{font-size:13px;color:#909399;background:#e8eaed;padding:4px 10px;border-radius:12px}
.members-list{display:flex;flex-direction:column;gap:10px;max-height:400px;overflow-y:auto;padding:8px}
.members-list::-webkit-scrollbar{width:6px}
.members-list::-webkit-scrollbar-thumb{background:#c1c4c9;border-radius:3px}
.member-item{display:flex;align-items:center;gap:12px;padding:14px;background:#fff;border:1px solid #e4e7ed;border-radius:12px;transition:all 0.3s}
.member-item:hover{background:#f8f9ff;border-color:#667eea;box-shadow:0 2px 8px rgba(102,126,234,0.1)}
.member-index{width:24px;height:24px;background:#e8eaed;border-radius:50%;display:flex;align-items:center;justify-content:center;font-size:12px;color:#606266;font-weight:600;flex-shrink:0}
.member-avatar{flex-shrink:0;box-shadow:0 2px 8px rgba(0,0,0,0.1)}
.member-info{flex:1;min-width:0}
.member-main{display:flex;align-items:center;gap:8px;margin-bottom:6px}
.member-name{font-size:15px;color:#303133;font-weight:600;white-space:nowrap;overflow:hidden;text-overflow:ellipsis}
.role-tag{font-size:12px;padding:2px 8px}
.member-meta{display:flex;align-items:center;gap:6px;font-size:12px;color:#909399}
.username{font-family:monospace}
.divider{color:#c0c4cc}
.join-time{}
.remove-btn{flex-shrink:0;font-size:12px;padding:6px 10px}
.add-member-content{padding:10px}
.search-results{margin-top:16px;max-height:300px;overflow-y:auto;border:1px solid #e4e7ed;border-radius:8px;padding:12px}
.user-item{padding:8px 0;border-bottom:1px solid #f0f0f0}
.user-item:last-child{border-bottom:none}
.dialog-footer{display:flex;justify-content:flex-end;gap:10px;padding:14px 20px;border-top:1px solid #f0f0f0;background:#fafafa}
@media(max-width:768px){
.groups-page{padding:16px 12px}
.groups-header{flex-direction:column;gap:12px;align-items:flex-start}
.groups-header h1{font-size:20px}
.groups-header .el-button{display:none}
.groups-grid{grid-template-columns:1fr}
.view-content{padding:0}
.members-header{flex-direction:column;align-items:flex-start;gap:12px}
.header-actions{width:100%}
.header-actions .el-button{width:100%}
.member-item{padding:10px}
.member-index{display:none}
.member-main{flex-wrap:wrap}
.member-meta{flex-wrap:wrap}
}
</style>
