<template>
  <div class="upload-page">
    <div class="upload-bg">
      <div class="bg-shape shape-1"></div>
      <div class="bg-shape shape-2"></div>
      <div class="bg-shape shape-3"></div>
    </div>
    <div class="upload-container">
      <div class="upload-header">
        <div class="header-icon">
          <el-icon class="icon-pen">
            <Edit/>
          </el-icon>
        </div>
        <h1>记录日常</h1>
        <p>分享你的生活点滴，记录每一个美好瞬间</p>
      </div>
      <el-card class="upload-card">
        <el-form :model="form" label-position="top">
          <el-form-item>
            <el-input v-model="form.title" placeholder="给今天的记录起个标题吧~" size="large" class="title-input"
                      maxlength="50" show-word-limit/>
          </el-form-item>
          <el-form-item>
            <el-input v-model="form.content" type="textarea" :rows="6"
                      placeholder="今天发生了什么有趣的事情？写下你的故事..." class="content-input" maxlength="2000"
                      show-word-limit/>
          </el-form-item>
          <el-form-item label="添加图片" class="image-form-item">
            <el-upload v-model:file-list="imageList" list-type="picture-card" :auto-upload="false"
                       :on-change="handleImageChange" :on-remove="handleRemove" :limit="9" multiple
                       class="image-uploader">
              <el-icon class="upload-icon">
                <Plus/>
              </el-icon>
            </el-upload>
            <div class="upload-tip">支持 jpg/png/gif/webp 格式，单张不超过 5MB，最多 9 张</div>
          </el-form-item>
          <el-form-item>
            <el-input v-model="form.tags" placeholder="添加标签，用逗号分隔，如：日常，美食，旅行" class="tags-input"
                      @focus="showHistory=true" ref="tagsInputRef">
              <template #prefix>
                <el-icon class="tag-icon">
                  <PriceTag/>
                </el-icon>
              </template>
            </el-input>
            <div class="tags-history" v-if="showHistory && historyTags.length>0" ref="historyRef">
              <div class="history-title">
                <el-icon>
                  <Clock/>
                </el-icon>
                <span>历史标签</span>
                <el-icon class="close-history" @click.stop="showHistory=false">
                  <Close/>
                </el-icon>
              </div>
              <div class="history-tags">
                <el-tag v-for="tag in historyTags" :key="tag"
                        class="history-tag"
                        size="default"
                        @click="selectHistoryTag(tag)">
                  {{ tag }}
                </el-tag>
              </div>
            </div>
          </el-form-item>
          <el-form-item label="可见性范围" class="visibility-form-item">
            <div class="visibility-options">
              <div class="visibility-option" :class="{active:form.visibility==='public'}"
                   @click="form.visibility='public';form.group_id=null">
                <div class="option-icon public-icon">
                  <el-icon>
                    <View/>
                  </el-icon>
                </div>
                <div class="option-content">
                  <div class="option-title">公开</div>
                  <div class="option-desc">所有人可见</div>
                </div>
                <div class="option-check" v-if="form.visibility==='public'">
                  <el-icon>
                    <CircleCheckFilled/>
                  </el-icon>
                </div>
              </div>
              <div class="visibility-option" :class="{active:form.visibility==='private'}"
                   @click="form.visibility='private'">
                <div class="option-icon private-icon">
                  <el-icon>
                    <Lock/>
                  </el-icon>
                </div>
                <div class="option-content">
                  <div class="option-title">仅组内可见</div>
                  <div class="option-desc">只有群组成员可见</div>
                </div>
                <div class="option-check" v-if="form.visibility==='private'">
                  <el-icon>
                    <CircleCheckFilled/>
                  </el-icon>
                </div>
              </div>
            </div>
          </el-form-item>
          <el-form-item label="选择可见群组" class="group-form-item" v-if="form.visibility==='private'">
            <el-select v-model="form.group_id" placeholder="请选择群组" class="group-select" clearable>
              <el-option v-for="g in myGroups" :key="g.id" :label="g.name" :value="g.id">
                <div class="group-option">
                  <span>{{ g.name }}</span>
                  <span class="group-members">{{ g.member_count }}人</span>
                </div>
              </el-option>
            </el-select>
            <div class="group-tip" v-if="myGroups.length===0">
              <el-icon>
                <InfoFilled/>
              </el-icon>
              <span>你还没有加入任何群组，<router-link to="/groups">去创建群组</router-link></span>
            </div>
          </el-form-item>
          <el-form-item class="submit-form-item">
            <div class="button-group">
              <el-button @click="reset" size="default" class="reset-btn">重置</el-button>
              <el-button type="primary" @click="submit" :loading="submitting" size="default" class="submit-btn">
                <el-icon class="submit-icon">
                  <UploadFilled/>
                </el-icon>
                <span class="submit-text">{{ submitting ? "发布中..." : "发布记录" }}</span>
              </el-button>
            </div>
          </el-form-item>
        </el-form>
      </el-card>
    </div>
  </div>
</template>

<script setup>
import {ref, onMounted, onUnmounted} from "vue"
import {
  Edit,
  Plus,
  UploadFilled,
  PriceTag,
  View,
  Lock,
  CircleCheckFilled,
  InfoFilled,
  Clock,
  Close
} from "@element-plus/icons-vue"
import {createRecord, uploadFile} from "../api/record"
import {getGroups} from "@/api/group"
import {ElMessage} from "element-plus"
import {useRouter} from "vue-router"

const router = useRouter()
const form = ref({title: "", content: "", tags: "", visibility: "public", group_id: null})
const imageList = ref([])
const submitting = ref(false)
const myGroups = ref([])
const historyTags = ref([])
const showHistory = ref(false)
const tagsInputRef = ref(null)
const historyRef = ref(null)

// 最大历史标签数量
const MAX_HISTORY_TAGS = 20

// 加载历史标签
const loadHistoryTags = () => {
  const saved = localStorage.getItem("tagsHistory")
  if (saved) {
    try {
      historyTags.value = JSON.parse(saved)
    } catch (e) {
      historyTags.value = []
    }
  }
}

// 保存历史标签
const saveHistoryTags = () => {
  localStorage.setItem("tagsHistory", JSON.stringify(historyTags.value))
}

// 添加标签到历史
const addTagsToHistory = (tagsStr) => {
  if (!tagsStr.trim()) return
  const newTags = tagsStr.split(",").map(t => t.trim()).filter(t => t)
  const updated = [...newTags, ...historyTags.value]
  // 去重并保留最新的 MAX_HISTORY_TAGS 个
  const unique = [...new Set(updated)].slice(0, MAX_HISTORY_TAGS)
  historyTags.value = unique
  saveHistoryTags()
}

// 选择历史标签
const selectHistoryTag = (selectedTag) => {
  const currentTags = form.value.tags.split(",").map(t => t.trim()).filter(t => t)
  if (!currentTags.includes(selectedTag)) {
    currentTags.push(selectedTag)
    form.value.tags = currentTags.join(",")
  }
  // 保持焦点在输入框
  setTimeout(() => {
    tagsInputRef.value?.focus()
  }, 100)
}

// 点击外部关闭历史标签
const handleClickOutside = (e) => {
  if (historyRef.value && !historyRef.value.contains(e.target) && tagsInputRef.value?.$el && !tagsInputRef.value.$el.contains(e.target)) {
    showHistory.value = false
  }
}

onMounted(async () => {
  try {
    myGroups.value = await getGroups()
    loadHistoryTags()
    document.addEventListener("click", handleClickOutside)
  } catch (e) {
    console.error(e)
  }
})

onUnmounted(() => {
  document.removeEventListener("click", handleClickOutside)
})

const handleImageChange = (file) => {
  const validTypes = ["image/jpeg", "image/png", "image/gif", "image/webp"]
  if (!validTypes.includes(file.raw.type)) {
    ElMessage.error("只能上传 JPG/PNG/GIF/WEBP 格式的图片");
    return false
  }
  if (file.raw.size > 5 * 1024 * 1024) {
    ElMessage.error("图片大小不能超过 5MB");
    return false
  }
}
const handleRemove = (file, fileList) => {
  imageList.value = fileList
}

const submit = async () => {
  if (!form.value.title.trim()) {
    ElMessage.warning("请输入标题");
    return
  }
  if (!form.value.content.trim() && imageList.value.length === 0) {
    ElMessage.warning("请输入内容或上传图片");
    return
  }
  submitting.value = true
  const tagsToSave = form.value.tags  // 保存当前标签
  try {
    const recordData = {
      type: "text",
      title: form.value.title,
      content: form.value.content,
      tags: form.value.tags,
      visibility: form.value.visibility,
      group_id: form.value.group_id || null
    }
    if (imageList.value.length > 0) {
      for (const img of imageList.value) {
        const formData = new FormData()
        formData.append("file", img.raw)
        formData.append("title", form.value.title)
        formData.append("content", form.value.content)
        formData.append("tags", form.value.tags)
        formData.append("type", "image")
        formData.append("visibility", form.value.visibility)
        if (form.value.group_id) formData.append("group_id", form.value.group_id)
        await uploadFile(formData)
      }
      ElMessage.success("发布成功")
    } else {
      await createRecord(recordData)
      ElMessage.success("发布成功")
    }
    form.value = {title: "", content: "", tags: "", visibility: "public", group_id: null}
    imageList.value = []
    // 保存标签到历史
    if (tagsToSave) {
      addTagsToHistory(tagsToSave)
    }
    setTimeout(() => router.push("/timeline"), 500)
  } catch (e) {
    console.error(e)
  } finally {
    submitting.value = false
  }
}
const reset = () => {
  form.value = {title: "", content: "", tags: ""};
  imageList.value = [];
  ElMessage.info("已重置")
}
</script>

<style scoped>
.upload-page {
  min-height: calc(100vh - 56px);
  position: relative;
  overflow: hidden;
  padding: 40px 20px
}

.upload-bg {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  z-index: 0;
  overflow: hidden
}

.bg-shape {
  position: absolute;
  border-radius: 50%;
  opacity: 0.1;
  filter: blur(60px)
}

.shape-1 {
  width: 400px;
  height: 400px;
  background: linear-gradient(135deg, #667eea, #764ba2);
  top: -100px;
  right: -100px;
  animation: float 8s ease-in-out infinite
}

.shape-2 {
  width: 300px;
  height: 300px;
  background: linear-gradient(135deg, #f093fb, #f5576c);
  bottom: -50px;
  left: -50px;
  animation: float 10s ease-in-out infinite reverse
}

.shape-3 {
  width: 200px;
  height: 200px;
  background: linear-gradient(135deg, #4facfe, #00f2fe);
  top: 50%;
  left: 50%;
  animation: float 12s ease-in-out infinite
}

@keyframes float {
  0%, 100% {
    transform: translate(0, 0) scale(1)
  }
  50% {
    transform: translate(30px, -30px) scale(1.05)
  }
}

.upload-container {
  position: relative;
  z-index: 1;
  max-width: 800px;
  margin: 0 auto
}

.upload-header {
  text-align: center;
  margin-bottom: 32px
}

.header-icon {
  width: 80px;
  height: 80px;
  margin: 0 auto 16px;
  background: linear-gradient(135deg, #667eea, #764ba2);
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  box-shadow: 0 10px 30px rgba(102, 126, 234, 0.4);
  animation: pulse 2s ease-in-out infinite
}

.header-icon .icon-pen {
  font-size: 40px;
  color: #fff
}

@keyframes pulse {
  0%, 100% {
    transform: scale(1);
    box-shadow: 0 10px 30px rgba(102, 126, 234, 0.4)
  }
  50% {
    transform: scale(1.05);
    box-shadow: 0 15px 40px rgba(102, 126, 234, 0.5)
  }
}

.upload-header h1 {
  font-size: 32px;
  color: #1a1a2e;
  margin-bottom: 8px;
  font-weight: 600
}

.upload-header p {
  font-size: 15px;
  color: #666
}

.upload-card {
  border-radius: 20px;
  box-shadow: 0 20px 60px rgba(0, 0, 0, 0.1);
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(10px)
}

.upload-card :deep(.el-card__body) {
  padding: 32px
}

.title-input :deep(.el-input__wrapper), .content-input :deep(.el-input__wrapper), .tags-input :deep(.el-input__wrapper) {
  border-radius: 16px;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.05);
  transition: all 0.3s
}

.title-input :deep(.el-input__wrapper):hover, .content-input :deep(.el-input__wrapper):hover, .tags-input :deep(.el-input__wrapper):hover {
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.1)
}

.title-input :deep(.el-input__wrapper.is-focus), .content-input :deep(.el-input__wrapper.is-focus), .tags-input :deep(.el-input__wrapper.is-focus) {
  box-shadow: 0 4px 20px rgba(102, 126, 234, 0.3)
}

.image-form-item {
  margin: 24px 0
}

.image-uploader {
  width: 100%
}

.image-uploader :deep(.el-upload--picture-card) {
  width: 100px;
  height: 100px;
  border-radius: 12px;
  border: 2px dashed #667eea;
  background: #f8f9ff;
  transition: all 0.3s
}

.image-uploader :deep(.el-upload--picture-card:hover) {
  border-color: #764ba2;
  background: #f0f2ff
}

.image-uploader :deep(.el-upload-list--picture-card .el-upload-list__item) {
  width: 100px;
  height: 100px;
  border-radius: 12px;
  transition: all 0.3s
}

.image-uploader :deep(.el-upload-list--picture-card .el-upload-list__item:hover) {
  transform: scale(1.05)
}

.upload-icon {
  font-size: 32px;
  color: #667eea;
  margin-bottom: 8px
}

.upload-tip {
  font-size: 12px;
  color: #999;
  text-align: center;
  margin-top: 8px
}

.tag-icon {
  color: #667eea
}

.tags-input {
  position: relative
}

.tags-history {
  position: absolute;
  top: 100%;
  left: 0;
  right: 0;
  margin-top: 8px;
  background: #fff;
  border: 1px solid #e4e7ed;
  border-radius: 12px;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.1);
  z-index: 100;
  max-height: 200px;
  overflow-y: auto
}

.tags-history::-webkit-scrollbar {
  width: 6px
}

.tags-history::-webkit-scrollbar-thumb {
  background: #c1c4c9;
  border-radius: 3px
}

.history-title {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 10px 12px;
  border-bottom: 1px solid #f0f0f0;
  font-size: 13px;
  color: #606266;
  font-weight: 600
}

.history-title .el-icon {
  margin-right: 6px;
  font-size: 14px;
  color: #667eea
}

.close-history {
  cursor: pointer;
  font-size: 14px;
  color: #909399;
  transition: color 0.2s
}

.close-history:hover {
  color: #667eea
}

.history-tags {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
  padding: 12px
}

.history-tag {
  cursor: pointer;
  padding: 6px 12px;
  border-radius: 6px;
  font-size: 13px;
  background: #f5f7fa;
  border: 1px solid #e4e7ed;
  color: #606266;
  transition: all 0.2s
}

.history-tag:hover {
  background: #667eea;
  border-color: #667eea;
  color: #fff;
  transform: translateY(-1px);
  box-shadow: 0 2px 8px rgba(102, 126, 234, 0.3)
}

.visibility-form-item {
  margin: 20px 0
}

.visibility-options {
  display: flex;
  gap: 12px;
  width: 100%
}

.visibility-option {
  flex: 1;
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 8px 10px;
  border: 2px solid #e4e7ed;
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.3s;
  background: #fff
}

.visibility-option:hover {
  border-color: #667eea;
  background: #f8f9ff;
  transform: translateY(-1px);
  box-shadow: 0 2px 8px rgba(102, 126, 234, 0.1)
}

.visibility-option.active {
  border-color: #667eea;
  background: linear-gradient(135deg, rgba(102, 126, 234, 0.05), rgba(118, 75, 162, 0.05));
  box-shadow: 0 2px 10px rgba(102, 126, 234, 0.15)
}

.option-icon {
  width: 28px;
  height: 28px;
  border-radius: 6px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 14px;
  flex-shrink: 0
}

.public-icon {
  background: linear-gradient(135deg, #667eea, #764ba2);
  color: #fff
}

.private-icon {
  background: linear-gradient(135deg, #f093fb, #f5576c);
  color: #fff
}

.option-content {
  flex: 1;
  min-width: 0
}

.option-title {
  font-size: 13px;
  color: #303133;
  font-weight: 600;
  margin-bottom: -8px
}

.option-desc {
  font-size: 12px;
  color: #909399
}

.option-check {
  width: 18px;
  height: 18px;
  color: #667eea;
  font-size: 18px;
  flex-shrink: 0
}

.group-form-item {
  margin: 20px 0
}

.group-select {
  width: 100%
}

.group-select :deep(.el-input__wrapper) {
  border-radius: 12px;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.05)
}

.group-option {
  display: flex;
  justify-content: space-between;
  align-items: center
}

.group-members {
  font-size: 12px;
  color: #909399
}

.group-tip {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 12px;
  background: #f8f9ff;
  border-radius: 8px;
  font-size: 13px;
  color: #667eea
}

.group-tip .el-icon {
  font-size: 16px
}

.group-tip a {
  color: #667eea;
  text-decoration: none;
  font-weight: 500
}

.group-tip a:hover {
  text-decoration: underline
}

.submit-form-item {
  margin-top: 32px
}

.button-group {
  display: flex;
  justify-content: flex-end;
  gap: 12px
}

.reset-btn, .submit-btn {
  width: 100px;
  height: 40px;
  padding: 0;
  border-radius: 10px;
  font-size: 14px;
  font-weight: 500
}

.reset-btn {
  background: #f5f7fa;
  color: #606266;
  border: 1px solid #e4e7ed
}

.reset-btn:hover {
  background: #e9ecef;
  color: #303133
}

.submit-btn {
  background: linear-gradient(135deg, #667eea, #764ba2);
  border: none;
  color: #fff
}

.submit-btn:hover {
  transform: translateY(-2px);
  box-shadow: 0 6px 20px rgba(102, 126, 234, 0.4)
}

.submit-btn:active {
  transform: translateY(0)
}

.submit-icon {
  margin-right: 6px;
  font-size: 16px
}

.submit-text {
  font-weight: 500
}

@media (max-width: 768px) {
  .upload-page {
    padding: 24px 16px
  }

  .header-icon {
    width: 60px;
    height: 60px
  }

  .header-icon .icon-pen {
    font-size: 28px
  }

  .upload-header h1 {
    font-size: 24px
  }

  .upload-header p {
    font-size: 13px
  }

  .upload-card :deep(.el-card__body) {
    padding: 20px
  }

  .image-uploader :deep(.el-upload--picture-card), .image-uploader :deep(.el-upload-list--picture-card .el-upload-list__item) {
    width: 80px;
    height: 80px
  }

  .submit-form-item {
    text-align: center
  }

  .button-group {
    display: flex;
    flex-direction: column;
    gap: 12px;
    justify-content: center;
    align-items: center;
    width: 100%
  }

  .reset-btn, .submit-btn {
    width: 280px;
    height: 44px;
    margin: 0 auto
  }

  .visibility-options {
    flex-direction: column
  }
}
</style>