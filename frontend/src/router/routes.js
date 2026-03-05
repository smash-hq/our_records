import Timeline from "../views/Timeline.vue"
import Upload from "../views/Upload.vue"
import Records from "../views/Records.vue"
import Login from "../views/Login.vue"

const routes = [
  { path: "/", redirect: "/login" },
  { path: "/login", name: "Login", component: Login },
  { path: "/timeline", name: "Timeline", component: Timeline, meta: { requiresAuth: true } },
  { path: "/upload", name: "Upload", component: Upload, meta: { requiresAuth: true } },
  { path: "/records", name: "Records", component: Records, meta: { requiresAuth: true } }
]

export default routes