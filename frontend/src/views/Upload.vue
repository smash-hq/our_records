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
<el-icon class="icon-pen"><Edit/></el-icon>
</div>
<h1>记录日常</h1>
<p>分享你的生活点滴，记录每一个美好瞬间</p>
</div>
<el-card class="upload-card">
<el-form :model="form" label-position="top">
<el-form-item>
<el-input v-model="form.title" placeholder="给今天的记录起个标题吧~" size="large" class="title-input" maxlength="50" show-word-limit/>
</el-form-item>
<el-form-item>
<el-input v-model="form.content" type="textarea" :rows="6" placeholder="今天发生了什么有趣的事情？写下你的故事..." class="content-input" maxlength="2000" show-word-limit/>
</el-form-item>
<el-form-item label="添加图片" class="image-form-item">
<el-upload v-model:file-list="imageList" list-type="picture-card" :auto-upload="false" :on-change="handleImageChange" :on-remove="handleRemove" :limit="9" multiple class="image-uploader">
<el-icon class="upload-icon"><Plus/></el-icon>
</el-upload>
<div class="upload-tip">支持 jpg/png/gif/webp 格式，单张不超过 5MB，最多 9 张</div>
</el-form-item>
<el-form-item>
<el-input v-model="form.tags" placeholder="添加标签，用逗号分隔，如：日常，美食，旅行" class="tags-input">
<template #prefix><el-icon class="tag-icon"><PriceTag/></el-icon></template>
</el-input>
</el-form-item>
<el-form-item class="submit-form-item">
<div class="button-group">
<el-button @click="reset" size="default" class="reset-btn">重置</el-button>
<el-button type="primary" @click="submit" :loading="submitting" size="default" class="submit-btn">
<el-icon class="submit-icon"><UploadFilled/></el-icon>
<span class="submit-text">{{submitting?"发布中...":"发布记录"}}</span>
</el-button>
</div>
</el-form-item>
</el-form>
</el-card>
</div>
</div>
</template>

<script setup>
import {ref} from "vue"
import {Edit,Plus,UploadFilled,PriceTag} from "@element-plus/icons-vue"
import {createRecord,uploadFile} from "../api/record"
import {ElMessage} from "element-plus"
import {useRouter} from "vue-router"

const router=useRouter()
const form=ref({title:"",content:"",tags:""})
const imageList=ref([])
const submitting=ref(false)

const handleImageChange=(file)=>{
  const validTypes=["image/jpeg","image/png","image/gif","image/webp"]
  if(!validTypes.includes(file.raw.type)){ElMessage.error("只能上传 JPG/PNG/GIF/WEBP 格式的图片");return false}
  if(file.raw.size>5*1024*1024){ElMessage.error("图片大小不能超过 5MB");return false}
}
const handleRemove=(file,fileList)=>{imageList.value=fileList}

const submit=async()=>{
  if(!form.value.title.trim()){ElMessage.warning("请输入标题");return}
  if(!form.value.content.trim()&&imageList.value.length===0){ElMessage.warning("请输入内容或上传图片");return}
  submitting.value=true
  try{
    if(imageList.value.length>0){
      for(const img of imageList.value){
        const formData=new FormData()
        formData.append("file",img.raw)
        formData.append("title",form.value.title)
        formData.append("content",form.value.content)
        formData.append("tags",form.value.tags)
        formData.append("type","image")
        await uploadFile(formData)
      }
      ElMessage.success("发布成功")
    }else{
      await createRecord({type:"text",title:form.value.title,content:form.value.content,tags:form.value.tags})
      ElMessage.success("发布成功")
    }
    form.value={title:"",content:"",tags:""}
    imageList.value=[]
    setTimeout(()=>router.push("/"),500)
  }catch(e){console.error(e)}
  finally{submitting.value=false}
}
const reset=()=>{form.value={title:"",content:"",tags:""};imageList.value=[];ElMessage.info("已重置")}
</script>

<style scoped>
.upload-page{min-height:calc(100vh - 56px);position:relative;overflow:hidden;padding:40px 20px}
.upload-bg{position:absolute;top:0;left:0;width:100%;height:100%;z-index:0;overflow:hidden}
.bg-shape{position:absolute;border-radius:50%;opacity:0.1;filter:blur(60px)}
.shape-1{width:400px;height:400px;background:linear-gradient(135deg,#667eea,#764ba2);top:-100px;right:-100px;animation:float 8s ease-in-out infinite}
.shape-2{width:300px;height:300px;background:linear-gradient(135deg,#f093fb,#f5576c);bottom:-50px;left:-50px;animation:float 10s ease-in-out infinite reverse}
.shape-3{width:200px;height:200px;background:linear-gradient(135deg,#4facfe,#00f2fe);top:50%;left:50%;animation:float 12s ease-in-out infinite}
@keyframes float{0%,100%{transform:translate(0,0) scale(1)}50%{transform:translate(30px,-30px) scale(1.05)}}
.upload-container{position:relative;z-index:1;max-width:800px;margin:0 auto}
.upload-header{text-align:center;margin-bottom:32px}
.header-icon{width:80px;height:80px;margin:0 auto 16px;background:linear-gradient(135deg,#667eea,#764ba2);border-radius:50%;display:flex;align-items:center;justify-content:center;box-shadow:0 10px 30px rgba(102,126,234,0.4);animation:pulse 2s ease-in-out infinite}
.header-icon .icon-pen{font-size:40px;color:#fff}
@keyframes pulse{0%,100%{transform:scale(1);box-shadow:0 10px 30px rgba(102,126,234,0.4)}50%{transform:scale(1.05);box-shadow:0 15px 40px rgba(102,126,234,0.5)}}
.upload-header h1{font-size:32px;color:#1a1a2e;margin-bottom:8px;font-weight:600}
.upload-header p{font-size:15px;color:#666}
.upload-card{border-radius:20px;box-shadow:0 20px 60px rgba(0,0,0,0.1);background:rgba(255,255,255,0.95);backdrop-filter:blur(10px)}
.upload-card :deep(.el-card__body){padding:32px}
.title-input :deep(.el-input__wrapper),.content-input :deep(.el-input__wrapper),.tags-input :deep(.el-input__wrapper){border-radius:16px;box-shadow:0 2px 12px rgba(0,0,0,0.05);transition:all 0.3s}
.title-input :deep(.el-input__wrapper):hover,.content-input :deep(.el-input__wrapper):hover,.tags-input :deep(.el-input__wrapper):hover{box-shadow:0 4px 16px rgba(0,0,0,0.1)}
.title-input :deep(.el-input__wrapper.is-focus),.content-input :deep(.el-input__wrapper.is-focus),.tags-input :deep(.el-input__wrapper.is-focus){box-shadow:0 4px 20px rgba(102,126,234,0.3)}
.image-form-item{margin:24px 0}
.image-uploader{width:100%}
.image-uploader :deep(.el-upload--picture-card){width:100px;height:100px;border-radius:12px;border:2px dashed #667eea;background:#f8f9ff;transition:all 0.3s}
.image-uploader :deep(.el-upload--picture-card:hover){border-color:#764ba2;background:#f0f2ff}
.image-uploader :deep(.el-upload-list--picture-card .el-upload-list__item){width:100px;height:100px;border-radius:12px;transition:all 0.3s}
.image-uploader :deep(.el-upload-list--picture-card .el-upload-list__item:hover){transform:scale(1.05)}
.upload-icon{font-size:32px;color:#667eea;margin-bottom:8px}
.upload-tip{font-size:12px;color:#999;text-align:center;margin-top:8px}
.tag-icon{color:#667eea}
.submit-form-item{margin-top:32px}
.button-group{display:flex;justify-content:flex-end;gap:12px}
.reset-btn,.submit-btn{width:100px;height:40px;padding:0;border-radius:10px;font-size:14px;font-weight:500}
.reset-btn{background:#f5f7fa;color:#606266;border:1px solid #e4e7ed}
.reset-btn:hover{background:#e9ecef;color:#303133}
.submit-btn{background:linear-gradient(135deg,#667eea,#764ba2);border:none;color:#fff}
.submit-btn:hover{transform:translateY(-2px);box-shadow:0 6px 20px rgba(102,126,234,0.4)}
.submit-btn:active{transform:translateY(0)}
.submit-icon{margin-right:6px;font-size:16px}
.submit-text{font-weight:500}
@media(max-width:768px){
.upload-page{padding:24px 16px}
.header-icon{width:60px;height:60px}
.header-icon .icon-pen{font-size:28px}
.upload-header h1{font-size:24px}
.upload-header p{font-size:13px}
.upload-card :deep(.el-card__body){padding:20px}
.image-uploader :deep(.el-upload--picture-card),.image-uploader :deep(.el-upload-list--picture-card .el-upload-list__item){width:80px;height:80px}
.submit-form-item{text-align:center}
.button-group{display:flex;flex-direction:column;gap:12px;justify-content:center;align-items:center;width:100%}
.reset-btn,.submit-btn{width:280px;height:44px;margin:0 auto}
}
</style>