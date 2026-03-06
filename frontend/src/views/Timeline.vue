<template>
  <div class="timeline-page">
    <div class="timeline-header">
      <h1 class="page-title">
        <el-icon>
          <Clock/>
        </el-icon>
        <span>时光轴</span></h1>
      <p class="page-subtitle">记录生活中的每一个美好瞬间</p>
    </div>
    <div class="timeline-container" v-loading="loading">
      <el-timeline class="timeline">
        <el-timeline-item v-for="record in records" :key="record.id" :timestamp="formatDate(record.created_at)"
                          placement="top" :color="getTypeColor(record.type)" size="large">
          <el-card class="record-card" shadow="hover" @click="viewRecord(record)">
            <div class="record-content-wrapper">
              <div class="record-thumb" v-if="hasImage(record)">
                <el-image :src="record.media_path" fit="cover" class="thumb-image"/>
                <div class="image-count" v-if="record.media_paths && record.media_paths.length > 1">
                  <el-icon>
                    <Picture/>
                  </el-icon>
                  {{ record.media_paths.length }}
                </div>
              </div>
              <div class="record-thumb record-thumb-text" v-else-if="record.type==='text'">
                <el-icon class="text-icon">
                  <Document/>
                </el-icon>
              </div>
              <div class="record-info">
                <div class="record-meta">
                  <el-tag :type="getTypeTag(record.type)" size="small" effect="plain">{{ getTypeLabel(record.type) }}
                  </el-tag>
                  <span class="record-time">{{ formatTime(record.created_at) }}</span>
                </div>
                <h3 class="record-title">{{ record.title || '无标题' }}</h3>
                <p class="record-brief">{{ getBrief(record.content) }}</p>
                <div class="record-footer">
                  <div class="record-tags" v-if="record.tags">
                    <el-tag v-for="tag in getTags(record.tags).slice(0,3)" :key="tag" size="small" class="tag-item">
                      #{{ tag }}
                    </el-tag>
                  </div>
                  <div class="record-comment-count" v-if="record.comment_count > 0" @click="viewRecord(record)">
                    <el-icon>
                      <ChatDotRound/>
                    </el-icon>
                    <span>{{ record.comment_count }}</span>
                  </div>
                </div>
              </div>
              <div class="card-author" v-if="record.username">
                <el-avatar :size="20" :src="record.user_avatar" class="author-avatar"/>
                <span class="author-name">{{ record.username }}</span>
              </div>
            </div>
            <div class="record-actions" @click.stop>
              <el-button class="btn-view" size="small" @click="viewRecord(record)">
                <el-icon>
                  <View/>
                </el-icon>
                查看
              </el-button>
              <el-button class="btn-delete" size="small" @click.stop="confirmDelete(record)">
                <el-icon>
                  <Delete/>
                </el-icon>
                删除
              </el-button>
            </div>
          </el-card>
        </el-timeline-item>
      </el-timeline>
      <el-empty v-if="!loading&&records.length===0" description="还没有记录，快来发布第一条吧~">
        <el-button type="primary" @click="goUpload">发布记录</el-button>
      </el-empty>
    </div>

    <!-- 详情弹窗 -->
    <el-dialog v-model="dialogVisible" title="记录详情" :width="dialogWidth" :close-on-click-modal="false">
      <div v-if="currentRecord" class="detail-content">
        <div class="detail-header">
          <el-tag :type="getTypeTag(currentRecord.type)">{{ getTypeLabel(currentRecord.type) }}</el-tag>
          <span>{{ formatDate(currentRecord.created_at) }}</span>
        </div>
        <h2>{{ currentRecord.title }}</h2>
        <p class="detail-text">{{ currentRecord.content }}</p>
        <div class="thumbnail-grid" v-if="currentRecord.media_paths && currentRecord.media_paths.length > 0">
          <div v-for="(img, index) in currentRecord.media_paths" :key="index" class="thumb-item">
            <el-image :src="img" fit="cover" class="thumb-img" preview-teleported
                      :preview-src-list="currentRecord.media_paths" :initial-index="index"/>
          </div>
        </div>
        <div class="detail-tags" v-if="currentRecord.tags">
          <el-tag v-for="tag in getTags(currentRecord.tags)" :key="tag" size="small" style="margin-right:8px">
            #{{ tag }}
          </el-tag>
        </div>
        <el-divider/>
        <comments-view :record-id="currentRecord.id"/>
      </div>
      <template #footer>
        <div class="dialog-footer">
          <el-button @click="dialogVisible=false">关闭</el-button>
        </div>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import {ref, onMounted, onUnmounted} from "vue"
import {Clock, Document, View, Delete, Picture, ChatDotRound} from "@element-plus/icons-vue"
import {getRecords, deleteRecord, updateRecord} from "../api/record"
import {ElMessage, ElMessageBox} from "element-plus"
import {useRouter} from "vue-router"
import CommentsView from "@/components/Comments.vue"

const router = useRouter()
const records = ref([]), loading = ref(false), dialogVisible = ref(false), currentRecord = ref(null)
const dialogWidth = ref("750px"), isEdit = ref(false), saving = ref(false)
const editForm = ref({title: "", content: "", tags: ""})

const updateDialogWidth = () => {
  dialogWidth.value = window.innerWidth <= 768 ? "92%" : "750px"
}
updateDialogWidth()
onMounted(() => {
  window.addEventListener("resize", updateDialogWidth);
  loadRecords()
})
onUnmounted(() => window.removeEventListener("resize", updateDialogWidth))

const loadRecords = async () => {
  loading.value = true
  try {
    records.value = await getRecords()
  } catch (e) {
  } finally {
    loading.value = false
  }
}
const formatDate = d => new Date(d).toLocaleDateString("zh-CN", {year: "numeric", month: "long", day: "numeric"})
const formatDateTime = d => new Date(d).toLocaleString("zh-CN", {
  year: "numeric",
  month: "long",
  day: "numeric",
  hour: "2-digit",
  minute: "2-digit"
})
const formatTime = d => new Date(d).toLocaleTimeString("zh-CN", {hour: "2-digit", minute: "2-digit"})
const getTypeTag = t => ({text: "", image: "success", video: "danger"})[t] || ""
const getTypeLabel = t => ({text: "文字", image: "图文", video: "视频"})[t] || t
const getTypeColor = t => ({text: "#409EFF", image: "#67C23A", video: "#F56C6C"})[t] || "#909399"
const hasImage = r => r.type === "image" && r.media_path
const getTags = t => t ? t.split(",").map(x => x.trim()) : []
const getBrief = t => {
  if (!t) return "暂无内容";
  return t.length > 80 ? t.substring(0, 80) + "..." : t
}

const viewRecord = r => {
  currentRecord.value = r;
  isEdit.value = false;
  dialogVisible.value = true
}
const goUpload = () => router.push("/upload")

const saveEdit = async () => {
  if (!editForm.value.title.trim()) {
    ElMessage.warning("请输入标题");
    return
  }
  saving.value = true
  try {
    await updateRecord(currentRecord.value.id, {
      type: currentRecord.value.type,
      title: editForm.value.title,
      content: editForm.value.content,
      tags: editForm.value.tags
    })
    ElMessage.success("保存成功");
    dialogVisible.value = false;
    loadRecords()
  } catch (e) {
    console.error(e)
  } finally {
    saving.value = false
  }
}

const confirmDelete = async (record) => {
  try {
    await ElMessageBox.confirm("确定要删除这条记录吗？", "删除确认", {
      confirmButtonText: "确定删除",
      cancelButtonText: "取消",
      type: "warning"
    })
    await deleteRecord(record.id);
    ElMessage.success("删除成功");
    loadRecords()
  } catch (e) {
    if (e !== "cancel") console.error(e)
  }
}
</script>

<style scoped>
.timeline-page {
  min-height: calc(100vh - 56px);
  background: linear-gradient(180deg, #f5f7fa 0%, #e4e8ec 100%);
  padding: 24px 12px 40px
}

.timeline-header {
  text-align: center;
  margin-bottom: 30px
}

.page-title {
  font-size: 24px;
  color: #303133;
  margin-bottom: 6px;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px
}

.page-subtitle {
  font-size: 13px;
  color: #909399
}

.timeline-container {
  max-width: 750px;
  margin: 0 auto
}

.timeline {
  padding: 10px 0
}

.record-card {
  border-radius: 10px;
  transition: all 0.25s;
  margin-left: 16px;
  cursor: pointer;
  position: relative
}

.record-card:hover {
  transform: translateX(4px);
  box-shadow: 0 8px 24px rgba(0, 0, 0, 0.08)
}

.record-content-wrapper {
  display: flex;
  gap: 16px;
  position: relative
}

.card-author {
  position: absolute;
  top: -9px;
  right: 12px;
  display: flex;
  align-items: center;
  gap: 5px;
  z-index: 10;
  background: rgba(255, 255, 255, 0.95);
  padding: 4px 8px;
  border-radius: 12px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.15)
}

.card-author .author-avatar {
  border: 2px solid #fff;
  box-shadow: 0 2px 6px rgba(0, 0, 0, 0.1)
}

.card-author .author-name {
  font-size: 11px;
  color: #667eea;
  font-weight: 500
}

.record-thumb {
  width: 70px;
  height: 50px;
  flex-shrink: 0;
  border-radius: 8px;
  overflow: hidden;
  background: #fff;
  position: relative
}

.record-thumb-text {
  display: flex;
  align-items: center;
  justify-content: center;
  background: #e9ecef
}

.text-icon {
  font-size: 48px;
  color: #c0c4cc
}

.thumb-image {
  width: 100%;
  height: 100%;
  object-fit: cover
}

.image-count {
  position: absolute;
  bottom: 4px;
  right: 4px;
  background: rgba(0, 0, 0, 0.6);
  color: #fff;
  border-radius: 12px;
  padding: 2px 8px;
  font-size: 11px;
  display: flex;
  align-items: center;
  gap: 3px
}

.image-count .el-icon {
  font-size: 12px
}

.record-info {
  flex: 1;
  min-width: 0
}

.record-meta {
  display: flex;
  align-items: center;
  gap: 10px;
  margin-bottom: 8px
}

.record-time {
  font-size: 12px;
  color: #909399
}

.record-title {
  font-size: 15px;
  color: #303133;
  margin-bottom: 6px;
  line-height: 1.4;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis
}

.record-brief {
  font-size: 13px;
  color: #606266;
  line-height: 1.6;
  margin-bottom: 10px;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden
}

.record-footer {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 10px;
  margin-bottom: 12px
}

.record-tags {
  display: flex;
  flex-wrap: wrap;
  gap: 6px;
  flex: 1
}

.tag-item {
  background: #f5f7fa;
  color: #667eea;
  border-color: #e4e7ed;
  font-size: 12px;
  padding: 4px 10px
}

.record-comment-count {
  display: flex;
  align-items: center;
  gap: 4px;
  font-size: 12px;
  color: #909399;
  cursor: pointer;
  padding: 4px 8px;
  border-radius: 4px;
  transition: all 0.2s
}

.record-comment-count:hover {
  background: #f5f7fa;
  color: #667eea
}

.record-comment-count .el-icon {
  font-size: 14px
}

.record-actions {
  display: flex;
  justify-content: flex-end;
  gap: 8px;
  padding: 12px 16px;
  background: #fff;
  flex-shrink: 0;
  border-top: 1px solid #f5f5f5;
  align-items: center;
  margin-top: auto
}

/* 查看按钮 - 蓝色系 */
.btn-view {
  min-width: 70px;
  padding: 6px 12px;
  border-radius: 6px;
  font-size: 13px;
  font-weight: 500;
  transition: all 0.2s;
  color: #409EFF;
  background: rgba(64, 158, 255, 0.1);
  border: 1px solid rgba(64, 158, 255, 0.3);
}

.btn-view:hover {
  background: rgba(64, 158, 255, 0.2);
  border-color: rgba(64, 158, 255, 0.5);
  color: #66b1ff;
}

/* 删除按钮 - 红色系 */
.btn-delete {
  min-width: 70px;
  padding: 6px 12px;
  border-radius: 6px;
  font-size: 13px;
  font-weight: 500;
  transition: all 0.2s;
  color: #F56C6C;
  background: rgba(245, 108, 108, 0.1);
  border: 1px solid rgba(245, 108, 108, 0.3);
}

.btn-delete:hover {
  background: rgba(245, 108, 108, 0.2);
  border-color: rgba(245, 108, 108, 0.5);
  color: #f78989;
}

.record-actions .el-button .el-icon {
  margin-right: 3px;
  font-size: 14px
}

.detail-content {
  padding: 10px
}

.detail-header {
  display: flex;
  align-items: center;
  gap: 15px;
  margin-bottom: 20px
}

.detail-text {
  font-size: 15px;
  color: #606266;
  line-height: 1.8;
  white-space: pre-wrap;
  margin: 20px 0
}

.thumbnail-grid {
  display: grid;
  grid-template-columns:repeat(auto-fill, minmax(100px, 1fr));
  gap: 10px;
  background: #ffffff;
  padding: 12px;
  border-radius: 10px;
  margin: 20px 0
}

.thumb-item {
  aspect-ratio: 1;
  border-radius: 8px;
  overflow: hidden;
  cursor: pointer;
  transition: transform 0.2s
}

.thumb-item:hover {
  transform: scale(1.05)
}

.thumb-img {
  width: 100%;
  height: 100%
}

.detail-tags {
  display: flex;
  flex-wrap: wrap;
  gap: 8px
}

.detail-tags .tag-item {
  font-size: 13px;
  padding: 5px 10px;
  color: #2e7d32;
  border-color: #a5d6a7;
  background: #f1f8e9
}

.dialog-footer {
  display: flex;
  justify-content: flex-end;
  gap: 10px;
  padding: 14px 20px;
  border-top: 1px solid #f0f0f0;
  background: #fafafa
}

@media (max-width: 768px) {
  .timeline-page {
    padding: 16px 8px 24px;
    min-height: calc(100vh - 50px)
  }

  .page-title {
    font-size: 20px
  }

  .page-subtitle {
    font-size: 12px
  }

  .record-content-wrapper {
    flex-direction: column
  }

  .record-thumb {
    width: 100%;
    height: 90px
  }

  .record-card {
    margin-left: 8px
  }

  .thumbnail-grid {
    grid-template-columns:repeat(3, 1fr)
  }

  .record-actions {
    gap: 10px;
    padding: 12px 12px
  }

  .btn-view, .btn-delete {
    min-width: 64px;
    padding: 6px 10px;
    font-size: 12px
  }
}
</style>