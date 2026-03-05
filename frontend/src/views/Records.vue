<template>
<div class="records-page">
<div class="records-header">
<h1><el-icon><Document/></el-icon><span>记录列表</span></h1>
<el-button type="primary" @click="$router.push('/upload')"><el-icon><Plus/></el-icon><span class="btn-text">新记录</span></el-button>
</div>
<div class="records-container" v-loading="loading">
<div class="search-bar">
<div class="search-wrapper">
<el-select v-model="queryForm.type" placeholder="全部类型" class="type-select" clearable>
<el-option label="文字" value="text">
<div class="option-content"><el-icon class="option-icon"><Document/></el-icon><span>文字</span></div>
</el-option>
<el-option label="图文" value="image">
<div class="option-content"><el-icon class="option-icon"><Picture/></el-icon><span>图文</span></div>
</el-option>
<el-option label="视频" value="video">
<div class="option-content"><el-icon class="option-icon"><VideoCamera/></el-icon><span>视频</span></div>
</el-option>
</el-select>
<div class="search-divider"></div>
<el-input v-model="queryForm.tags" placeholder="搜索标签..." class="search-input" clearable @keyup.enter="loadRecords">
<template #prefix><el-icon><Search/></el-icon></template>
</el-input>
<el-button type="primary" class="search-btn" @click="loadRecords">
<el-icon><Search/></el-icon>搜索
</el-button>
</div>
</div>
<div class="records-grid">
<el-card v-for="record in records" :key="record.id" class="record-item" shadow="hover" @click="viewRecord(record)">
<div class="record-cover" v-if="hasImage(record)">
<el-image :src="record.media_path" fit="cover" class="cover-image"/>
<div class="image-count" v-if="record.media_paths && record.media_paths.length > 1">
<el-icon><Picture/></el-icon> {{record.media_paths.length}}
</div>
</div>
<div class="record-cover record-cover-text" v-else-if="record.type==='text'">
<el-icon class="text-icon"><Document/></el-icon>
</div>
<div class="record-body">
<div class="record-meta">
<el-tag :type="getTypeTag(record.type)" size="small">{{getTypeLabel(record.type)}}</el-tag>
<span class="record-date">{{formatDate(record.created_at)}}</span>
</div>
<h3 class="record-title">{{record.title||'无标题'}}</h3>
<p class="record-excerpt">{{getExcerpt(record.content)}}</p>
<div class="record-tags" v-if="record.tags">
<el-tag v-for="tag in getTags(record.tags)" :key="tag" size="small" class="tag-item">#{{tag}}</el-tag>
</div>
</div>
<div class="record-actions" @click.stop>
<el-button class="btn-view" size="small" @click="viewRecord(record)"><el-icon><View/></el-icon>查看</el-button>
<el-button class="btn-delete" size="small" @click.stop="confirmDelete(record.id)"><el-icon><Delete/></el-icon>删除</el-button>
</div>
</el-card>
</div>
<el-empty v-if="!loading&&records.length===0" description="暂无记录"/>
</div>

<!-- 详情弹窗 -->
<el-dialog v-model="dialogVisible" title="记录详情" :width="dialogWidth">
<div v-if="currentRecord" class="detail-content">
<div class="detail-header">
<el-tag :type="getTypeTag(currentRecord.type)">{{getTypeLabel(currentRecord.type)}}</el-tag>
<span>{{formatDate(currentRecord.created_at)}}</span>
</div>
<h2>{{currentRecord.title}}</h2>
<p class="detail-text">{{currentRecord.content}}</p>
<div class="thumbnail-grid" v-if="currentRecord.media_paths && currentRecord.media_paths.length > 0">
<div v-for="(img, index) in currentRecord.media_paths" :key="index" class="thumb-item">
<el-image :src="img" fit="cover" class="thumb-img" preview-teleported :preview-src-list="currentRecord.media_paths" :initial-index="index"/>
</div>
</div>
<div class="detail-tags" v-if="currentRecord.tags">
<el-tag v-for="tag in getTags(currentRecord.tags)" :key="tag" size="small" style="margin-right:8px">#{{tag}}</el-tag>
</div>
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
import {ref,onMounted,onUnmounted} from "vue"
import {Document,Plus,Picture,View,Delete,Search,VideoCamera} from "@element-plus/icons-vue"
import {getRecords,deleteRecord} from "../api/record"
import {ElMessage,ElMessageBox} from "element-plus"

const records=ref([]),loading=ref(false),queryForm=ref({type:"",tags:""})
const dialogVisible=ref(false),currentRecord=ref(null),dialogWidth=ref("700px")

const updateDialogWidth=()=>{dialogWidth.value=window.innerWidth<=768?"90%":"700px"}
updateDialogWidth()
onMounted(()=>{window.addEventListener("resize",updateDialogWidth);loadRecords()})
onUnmounted(()=>window.removeEventListener("resize",updateDialogWidth))

const loadRecords=async()=>{loading.value=true;try{records.value=await getRecords(queryForm.value)}catch(e){}finally{loading.value=false}}
const formatDate=d=>new Date(d).toLocaleDateString("zh-CN",{year:"numeric",month:"long",day:"numeric"})
const getTypeTag=t=>({text:"",image:"success",video:"danger"})[t]||""
const getTypeLabel=t=>({text:"文字",image:"图文",video:"视频"})[t]||t
const hasImage=r=>r.type==="image"&&r.media_path
const getTags=t=>t?t.split(",").map(x=>x.trim()):[]
const getExcerpt=t=>{if(!t)return"";return t.substring(0,100)+(t.length>100?"...":"")}

const viewRecord=r=>{currentRecord.value=r;dialogVisible.value=true}

const confirmDelete=async(id)=>{
  try{
    await ElMessageBox.confirm("确定要删除这条记录吗？","删除确认",{confirmButtonText:"确定删除",cancelButtonText:"取消",type:"warning"})
    await deleteRecord(id)
    ElMessage.success("删除成功")
    loadRecords()
    if(currentRecord.value&&currentRecord.value.id===id){dialogVisible.value=false;currentRecord.value=null}
  }catch(e){if(e!=="cancel")console.error(e)}
}
</script>

<style scoped>
.records-page{min-height:calc(100vh - 56px);padding:24px 16px}
.records-header{max-width:1200px;margin:0 auto 24px;display:flex;justify-content:space-between;align-items:center}
.records-header h1{font-size:24px;color:#303133;display:flex;align-items:center;gap:8px}
.records-container{max-width:1200px;margin:0 auto}
.search-bar{background:#fff;border-radius:16px;padding:12px 16px;margin-bottom:24px;box-shadow:0 2px 16px rgba(0,0,0,0.06)}
.search-wrapper{display:flex;align-items:center;gap:12px}
.type-select{width:140px}
.type-select :deep(.el-input__wrapper){border-radius:10px;background:#f5f7fa;box-shadow:none;padding:0 12px;height:40px}
.type-select :deep(.el-input__inner){font-size:14px;color:#606266}
.option-content{display:flex;align-items:center;gap:8px}
.option-icon{font-size:16px;color:#667eea}
.search-divider{width:1px;height:24px;background:#e4e7ed;margin:0 8px}
.search-input{flex:1;max-width:350px}
.search-input :deep(.el-input__wrapper){border-radius:10px;background:#f5f7fa;box-shadow:none;padding:0 12px 0 32px;height:40px}
.search-input :deep(.el-input__inner){font-size:14px;color:#606266}
.search-input :deep(.el-input__prefix){left:8px;color:#909399}
.search-btn{height:40px;padding:0 20px;border-radius:10px;background:linear-gradient(135deg,#667eea,#764ba2);border:none;font-weight:500}
.search-btn:hover{transform:translateY(-2px);box-shadow:0 6px 20px rgba(102,126,234,0.4)}
.search-btn .el-icon{margin-right:6px;font-size:16px}
.records-grid{display:grid;grid-template-columns:repeat(auto-fill,minmax(280px,1fr));gap:20px}
.record-item{border-radius:16px;transition:all 0.3s;overflow:hidden;display:flex;flex-direction:column;border:none;cursor:pointer}
.record-item:hover{transform:translateY(-6px);box-shadow:0 16px 40px rgba(0,0,0,0.12)}
.record-cover{height:180px;overflow:hidden;background:linear-gradient(135deg,#f5f7fa 0%,#e4e8ec 100%);position:relative;flex-shrink:0}
.record-cover-text{display:flex;align-items:center;justify-content:center}
.text-icon{font-size:64px;color:#c0c4cc}
.cover-image{width:100%;height:100%;object-fit:cover;transition:transform 0.4s}
.record-item:hover .cover-image{transform:scale(1.08)}
.image-count{position:absolute;bottom:8px;right:8px;background:rgba(0,0,0,0.7);color:#fff;border-radius:16px;padding:4px 12px;font-size:12px;display:flex;align-items:center;gap:4px;box-shadow:0 2px 8px rgba(0,0,0,0.2)}
.image-count .el-icon{font-size:14px}
.record-body{padding:18px;flex:1}
.record-meta{display:flex;justify-content:space-between;align-items:center;margin-bottom:12px}
.record-date{font-size:12px;color:#999}
.record-title{font-size:16px;color:#2c3e50;margin-bottom:10px;line-height:1.5;height:48px;overflow:hidden;display:-webkit-box;-webkit-line-clamp:2;-webkit-box-orient:vertical;font-weight:500}
.record-excerpt{font-size:13px;color:#7f8c8d;line-height:1.7;margin-bottom:14px;height:44px;overflow:hidden;display:-webkit-box;-webkit-line-clamp:2;-webkit-box-orient:vertical}
.record-tags{display:flex;flex-wrap:wrap;gap:6px;margin-bottom:14px;min-height:20px}
.tag-item{color:#667eea;border-color:#e0e6ed;background:#f8f9ff;font-size:12px;padding:3px 8px;border-radius:6px}
.record-actions{display:flex;justify-content:flex-end;gap:8px;padding:10px 16px;background:transparent;flex-shrink:0;border-top:1px solid #f5f5f5;margin-top:8px}
.btn-view{min-width:70px;padding:6px 12px;border-radius:6px;font-size:13px;font-weight:500;transition:all 0.2s;color:#409EFF;background:rgba(64,158,255,0.1);border:1px solid rgba(64,158,255,0.3)}
.btn-view:hover{background:rgba(64,158,255,0.2);border-color:rgba(64,158,255,0.5);color:#66b1ff}
.btn-delete{min-width:70px;padding:6px 12px;border-radius:6px;font-size:13px;font-weight:500;transition:all 0.2s;color:#F56C6C;background:rgba(245,108,108,0.1);border:1px solid rgba(245,108,108,0.3)}
.btn-delete:hover{background:rgba(245,108,108,0.2);border-color:rgba(245,108,108,0.5);color:#f78989}
.record-actions .el-button .el-icon{margin-right:3px;font-size:14px}
.detail-content{padding:10px}
.detail-header{display:flex;align-items:center;gap:15px;margin-bottom:20px}
.detail-text{font-size:15px;color:#606266;line-height:1.8;white-space:pre-wrap;margin:20px 0}
.thumbnail-grid{display:grid;grid-template-columns:repeat(auto-fill,minmax(100px,1fr));gap:10px;background:#fafafa;padding:12px;border-radius:10px;margin:20px 0}
.thumb-item{aspect-ratio:1;border-radius:8px;overflow:hidden;cursor:pointer;transition:transform 0.2s}
.thumb-item:hover{transform:scale(1.05)}
.thumb-img{width:100%;height:100%}
.detail-tags{display:flex;flex-wrap:wrap;gap:8px}
.detail-tags .tag-item{font-size:13px;padding:5px 10px;color:#2e7d32;border-color:#a5d6a7;background:#f1f8e9}
.dialog-footer{display:flex;justify-content:flex-end;gap:10px;padding:14px 20px;border-top:1px solid #f0f0f0;background:#fafafa}
@media(max-width:768px){
.records-page{padding:16px 12px}
.records-header{flex-direction:column;gap:12px;align-items:flex-start}
.records-header h1{font-size:20px}
.records-header .el-button{display:none}
.search-bar{border-radius:14px;padding:10px}
.search-wrapper{gap:8px;flex-wrap:nowrap}
.type-select{width:90px;flex-shrink:0}
.type-select :deep(.el-input__wrapper){height:36px;padding:0 8px}
.type-select :deep(.el-input__inner){font-size:13px}
.search-divider{width:1px;height:20px;background:#e4e7ed;margin:0 6px;flex-shrink:0}
.search-input{flex:1;min-width:0;max-width:none}
.search-input :deep(.el-input__wrapper){height:36px;padding:0 10px 0 30px}
.search-input :deep(.el-input__inner){font-size:13px}
.search-input :deep(.el-input__prefix){left:6px}
.search-input :deep(.el-input__prefix-inner>.el-icon){font-size:14px}
.search-btn{height:36px;padding:0 12px;border-radius:10px;flex-shrink:0}
.search-btn .el-icon{margin-right:4px;font-size:14px}
.search-btn span{font-size:13px}
.records-grid{grid-template-columns:1fr}
.record-cover{height:200px}
.record-title{font-size:15px}
.record-excerpt{font-size:12px}
.record-actions{gap:10px;padding:10px 12px}
.btn-view,.btn-delete{min-width:64px;padding:6px 10px;font-size:12px}
.thumbnail-grid{grid-template-columns:repeat(3,1fr)}
}
</style>