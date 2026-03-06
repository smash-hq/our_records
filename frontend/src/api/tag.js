import request from "../utils/request"

// 获取我的标签
export function getMyTags() {
  return request.get("/tags/my")
}
