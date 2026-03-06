import request from "../utils/request"

// 用户登录
export function login(data) {
  return request.post("/auth/login", data)
}

// 用户注册
export function register(data) {
  return request.post("/auth/register", data)
}

// 获取当前用户信息
export function getCurrentUser() {
  return request.get("/user")
}

// 修改密码
export function changePassword(data) {
  return request.put("/user/password", data)
}

// 上传头像
export function uploadAvatar(formData) {
  return request.post("/user/avatar", formData)
}
