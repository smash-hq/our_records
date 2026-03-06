import request from "../utils/request"

// 获取未读通知数量
export function getUnreadNotifications() {
  return request.get("/notifications/unread")
}

// 获取通知列表
export function getNotifications(params) {
  return request.get("/notifications", { params })
}

// 标记通知为已读
export function markNotificationAsRead(id) {
  return request.put(`/notifications/${id}/read`)
}

// 标记所有通知为已读
export function markAllNotificationsAsRead() {
  return request.post("/notifications/read-all")
}
