import request from "../utils/request"

export function getRecords(params) {
  return request.get("/records", { params })
}

export function getRecord(id) {
  return request.get("/records/" + id)
}

export function createRecord(data) {
  return request.post("/records", data)
}

export function updateRecord(id, data) {
  return request.put("/records/" + id, data)
}

export function deleteRecord(id) {
  return request.delete("/records/" + id)
}

export function uploadFile(formData) {
  return request.post("/upload", formData, {
    headers: { "Content-Type": "multipart/form-data" }
  })
}