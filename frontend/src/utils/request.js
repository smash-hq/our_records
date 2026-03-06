import axios from "axios"
import { ElMessage } from "element-plus"

const request = axios.create({
  baseURL: "/api",
  timeout: 10000
})

// 请求拦截器 - 添加 token
request.interceptors.request.use(
  config => {
    const token = localStorage.getItem("token")
    if (token) {
      config.headers.Authorization = `Bearer ${token}`
    }
    return config
  },
  error => {
    return Promise.reject(error)
  }
)

// 登出处理函数
const handleLogout = () => {
  localStorage.removeItem("token")
  localStorage.removeItem("user")
  // 触发 storage 事件通知其他组件
  window.dispatchEvent(new Event("storage"))
  // 跳转到登录页
  window.location.href = "/login"
}

// 响应拦截器
request.interceptors.response.use(
  response => response.data,
  error => {
    if (error.response?.status === 401) {
      // 只在非通知接口返回 401 时登出
      if (!error.config?.url?.includes('/notifications')) {
        handleLogout()
        return Promise.reject(new Error("未登录或登录已过期"))
      }
    }
    // 通知接口 401 时不显示错误提示
    if (!error.config?.url?.includes('/notifications')) {
      ElMessage.error(error.response?.data?.error || "请求失败")
    }
    return Promise.reject(error)
  }
)

export default request