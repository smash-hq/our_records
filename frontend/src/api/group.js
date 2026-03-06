import request from "../utils/request"

// 获取我的群组列表
export function getGroups() {
  return request.get("/groups")
}

// 创建群组
export function createGroup(data) {
  return request.post("/groups", data)
}

// 获取群组详情
export function getGroup(id) {
  return request.get("/groups/" + id)
}

// 更新群组
export function updateGroup(id, data) {
  return request.put("/groups/" + id, data)
}

// 删除群组
export function deleteGroup(id) {
  return request.delete("/groups/" + id)
}

// 上传群组头像
export function uploadGroupAvatar(id, formData) {
  return request.post("/groups/" + id + "/avatar", formData)
}

// 添加群组成员
export function addGroupMember(id, data) {
  return request.post("/groups/" + id + "/members", data)
}

// 移除群组成员
export function removeGroupMember(id, data) {
  return request.delete("/groups/" + id + "/members", data)
}

// 获取群组成员列表
export function getGroupMembers(id) {
  return request.get("/groups/" + id + "/members")
}

// 退出群组
export function leaveGroup(id) {
  return request.post("/groups/" + id + "/leave")
}

// 搜索用户
export function searchUsers(keyword) {
  return request.get("/users/search", { params: { keyword } })
}
