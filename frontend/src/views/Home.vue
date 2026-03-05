<template>
  <div class="home">
    <el-card>
      <h2>欢迎使用 Our Records</h2>
      <p>这是一个生活记录应用，支持记录图片、文字、语音、视频等信息。</p>
      <el-row :gutter="20" style="margin-top: 30px">
        <el-col :span="8">
          <el-statistic title="记录类型" :value="4" />
        </el-col>
        <el-col :span="8">
          <el-statistic title="存储方式" :value="2" />
        </el-col>
        <el-col :span="8">
          <el-statistic title="支持格式" value="10+" />
        </el-col>
      </el-row>
      <div style="margin-top: 30px">
        <el-button type="primary" @click="goUpload">开始记录</el-button>
        <el-button @click="goRecords">查看记录</el-button>
      </div>
    </el-card>
  </div>
</template>

<script setup>
import { useRouter } from "vue-router"
const router = useRouter()
const goUpload = () => router.push("/upload")
const goRecords = () => router.push("/records")
</script>

<style scoped>
.home { max-width: 800px; margin: 0 auto; }
.home h2 { color: #409EFF; }
</style>