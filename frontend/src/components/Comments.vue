<template>
<div class="comments-section">
<div class="comments-header">
<h3>评论</h3>
<span class="comment-count">{{comments.length}} 条评论</span>
</div>

<!-- 发表评论 -->
<div class="comment-form">
<el-input
  v-model="newComment"
  type="textarea"
  :rows="3"
  placeholder="写下你的评论..."
  maxlength="500"
  show-word-limit
  class="comment-input"
/>
<div class="comment-actions">
<el-button @click="cancelReply" v-if="replyTo">取消回复</el-button>
<el-button type="primary" @click="submitComment" :loading="submitting">
{{replyTo ? '回复' : '发表'}}
</el-button>
</div>
</div>

<!-- 评论列表 -->
<div class="comments-list" v-loading="loading">
<div v-for="comment in comments" :key="comment.id" class="comment-item">
<div class="comment-avatar">
<el-avatar :size="40" :src="comment.avatar">
{{comment.nickname?.charAt(0) || comment.username?.charAt(0) || 'U'}}
</el-avatar>
</div>
<div class="comment-content">
<div class="comment-info">
<span class="comment-author">{{comment.nickname || comment.username}}</span>
<span class="comment-floor">{{comment.floor}}楼</span>
<span class="comment-time">{{formatTime(comment.created_at)}}</span>
</div>
<div class="comment-text">{{comment.content}}</div>
<div class="comment-meta">
<span class="comment-likes">
<el-icon @click="like(comment)" :class="{liked: comment.liked}"><Star/></el-icon>
{{comment.like_count}}
</span>
<el-button link type="primary" size="small" @click="showReplyInput(comment)">
<el-icon><ChatDotRound/></el-icon>回复
</el-button>
<el-button link type="danger" size="small" @click="deleteComment(comment)" v-if="canDelete(comment)">
<el-icon><Delete/></el-icon>删除
</el-button>
</div>

<!-- 回复输入框 -->
<div class="reply-form" v-if="replyTo && replyTo.id === comment.id">
<el-input
  v-model="replyContent"
  type="textarea"
  :rows="2"
  :placeholder="'回复 @' + (comment.nickname || comment.username)"
  maxlength="500"
  class="reply-input"
/>
<div class="reply-actions">
<el-button size="small" @click="cancelReply">取消</el-button>
<el-button size="small" type="primary" @click="submitReply(comment.id)" :loading="submitting">回复</el-button>
</div>
</div>

<!-- 回复列表 -->
<div v-if="comment.replies && comment.replies.length > 0" class="replies-list">
<div v-for="reply in comment.replies" :key="reply.id" class="reply-item">
<div class="reply-avatar">
<el-avatar :size="32" :src="reply.avatar">
{{reply.nickname?.charAt(0) || reply.username?.charAt(0) || 'U'}}
</el-avatar>
</div>
<div class="reply-content">
<div class="reply-info">
<span class="reply-author">{{reply.nickname || reply.username}}</span>
<span class="reply-time">{{formatTime(reply.created_at)}}</span>
</div>
<div class="reply-text">{{reply.content}}</div>
<div class="reply-meta">
<span class="reply-likes">
<el-icon><Star/></el-icon>
{{reply.like_count}}
</span>
<el-button link type="danger" size="small" @click="deleteComment(reply)" v-if="canDelete(reply)">
<el-icon><Delete/></el-icon>
</el-button>
</div>
</div>
</div>
<div v-if="comment.reply_count > (comment.replies?.length || 0)" class="view-more-replies">
<el-button link type="primary" @click="loadMoreReplies(comment)">
查看全部 {{comment.reply_count}} 条回复 <el-icon><ArrowRight/></el-icon>
</el-button>
</div>
</div>
</div>
</div>

<el-empty v-if="!loading && comments.length === 0" description="暂无评论，快来抢沙发吧"/>
</div>
</div>
</template>

<script setup>
import { ref, reactive, onMounted, watch } from "vue"
import { Star, Delete, ChatDotRound, ArrowRight } from "@element-plus/icons-vue"
import { ElMessage, ElMessageBox } from "element-plus"
import { getComments, createComment, deleteComment as apiDeleteComment, likeComment, getCommentReplies } from "@/api/comment"

const props = defineProps({
  recordId: {
    type: Number,
    required: true
  }
})

const comments = ref([])
const loading = ref(false)
const submitting = ref(false)
const newComment = ref("")
const replyContent = ref("")
const replyTo = ref(null)
const currentUserId = ref(null)

// 加载评论
const loadComments = async () => {
  loading.value = true
  try {
    const res = await getComments(props.recordId)
    comments.value = res.comments || []
  } catch (error) {
    console.error(error)
  } finally {
    loading.value = false
  }
}

// 获取当前用户 ID
onMounted(() => {
  const userStr = localStorage.getItem("user")
  if (userStr) {
    try {
      const user = JSON.parse(userStr)
      currentUserId.value = user.id
    } catch (e) {}
  }
  loadComments()
})

// 监听 recordId 变化，重新加载评论
watch(() => props.recordId, () => {
  loadComments()
})

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

const submitComment = async () => {
  if (!newComment.value.trim()) {
    ElMessage.warning("请输入评论内容")
    return
  }
  
  submitting.value = true
  try {
    await createComment(props.recordId, {
      content: newComment.value.trim(),
      parent_id: null
    })
    ElMessage.success("评论成功")
    newComment.value = ""
    loadComments()
  } catch (error) {
    console.error(error)
  } finally {
    submitting.value = false
  }
}

const showReplyInput = (comment) => {
  replyTo.value = comment
  replyContent.value = ""
}

const cancelReply = () => {
  replyTo.value = null
  replyContent.value = ""
}

const submitReply = async (parentId) => {
  if (!replyContent.value.trim()) {
    ElMessage.warning("请输入回复内容")
    return
  }
  
  submitting.value = true
  try {
    await createComment(props.recordId, {
      content: replyContent.value.trim(),
      parent_id: parentId
    })
    ElMessage.success("回复成功")
    cancelReply()
    loadComments()
  } catch (error) {
    console.error(error)
  } finally {
    submitting.value = false
  }
}

const deleteComment = async (comment) => {
  try {
    await ElMessageBox.confirm("确定要删除这条评论吗？", "删除确认", {
      confirmButtonText: "确定",
      cancelButtonText: "取消",
      type: "warning"
    })
    await apiDeleteComment(comment.id)
    ElMessage.success("删除成功")
    loadComments()
  } catch (error) {
    if (error !== "cancel") console.error(error)
  }
}

const canDelete = (comment) => {
  return currentUserId.value && (comment.user_id === currentUserId.value)
}

const like = async (comment) => {
  try {
    await likeComment(comment.id)
    comment.like_count++
    comment.liked = true
  } catch (error) {
    console.error(error)
  }
}

const loadMoreReplies = async (comment) => {
  try {
    const res = await getCommentReplies(props.recordId, comment.id)
    comment.replies = res.replies || []
  } catch (error) {
    console.error(error)
  }
}
</script>

<style scoped>
.comments-section{padding:20px 0}
.comments-header{display:flex;justify-content:space-between;align-items:center;margin-bottom:16px;padding-bottom:12px;border-bottom:1px solid #e4e7ed}
.comments-header h3{font-size:16px;color:#303133;margin:0}
.comment-count{font-size:13px;color:#909399}
.comment-form{margin-bottom:16px}
.comment-input{margin-bottom:10px}
.comment-input :deep(.el-textarea__inner){border-radius:8px;box-shadow:0 2px 8px rgba(0,0,0,0.05)}
.comment-actions{display:flex;justify-content:flex-end;gap:10px}
.comments-list{display:flex;flex-direction:column;gap:10px}
.comment-item{display:flex;gap:10px;padding:12px;background:#fff;border-radius:10px;border:1px solid #f0f0f0}
.comment-avatar{flex-shrink:0}
.comment-content{flex:1;min-width:0}
.comment-info{display:flex;align-items:center;gap:6px;margin-bottom:6px;flex-wrap:wrap}
.comment-author{font-size:14px;color:#303133;font-weight:600}
.comment-floor{font-size:12px;color:#909399;background:#f5f7fa;padding:2px 6px;border-radius:4px}
.comment-time{font-size:12px;color:#909399}
.comment-text{font-size:14px;color:#606266;line-height:1.5;margin-bottom:8px;white-space:pre-wrap}
.comment-meta{display:flex;align-items:center;gap:12px}
.comment-likes{display:flex;align-items:center;gap:4px;font-size:13px;color:#909399;cursor:pointer;transition:color 0.2s}
.comment-likes:hover{color:#667eea}
.comment-likes.liked{color:#667eea}
.reply-form{margin-top:10px;padding:10px;background:#f8f9ff;border-radius:8px}
.reply-input{margin-bottom:8px}
.reply-input :deep(.el-textarea__inner){border-radius:6px;font-size:13px}
.reply-actions{display:flex;justify-content:flex-end;gap:8px}
.replies-list{margin-top:10px;padding-left:10px;border-left:2px solid #e4e7ed;display:flex;flex-direction:column;gap:8px}
.reply-item{display:flex;gap:8px;padding:8px;background:#fafafa;border-radius:8px}
.reply-avatar{flex-shrink:0}
.reply-content{flex:1;min-width:0}
.reply-info{display:flex;align-items:center;gap:6px;margin-bottom:4px}
.reply-author{font-size:13px;color:#303133;font-weight:600}
.reply-time{font-size:12px;color:#909399}
.reply-text{font-size:13px;color:#606266;line-height:1.5;margin-bottom:6px;white-space:pre-wrap}
.reply-meta{display:flex;align-items:center;gap:12px}
.reply-likes{display:flex;align-items:center;gap:4px;font-size:12px;color:#909399}
.view-more-replies{margin-top:8px;padding-left:12px}
.view-more-replies .el-button{font-size:13px}
</style>
