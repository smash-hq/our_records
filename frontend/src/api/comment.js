import request from "../utils/request"

// 获取记录评论列表
export function getComments(recordId) {
  return request.get("/records/" + recordId + "/comments")
}

// 创建评论
export function createComment(recordId, data) {
  return request.post("/records/" + recordId + "/comments", data)
}

// 获取评论回复
export function getCommentReplies(recordId, commentId) {
  return request.get("/records/" + recordId + "/comments/" + commentId + "/replies")
}

// 删除评论
export function deleteComment(commentId) {
  return request.delete("/comments/" + commentId)
}

// 点赞评论
export function likeComment(commentId) {
  return request.post("/comments/" + commentId + "/like")
}

// 获取我的评论
export function getMyComments(params) {
  return request.get("/comments/my", { params })
}
