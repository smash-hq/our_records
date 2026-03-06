import Timeline from "../views/Timeline.vue"
import Upload from "../views/Upload.vue"
import Records from "../views/Records.vue"
import Login from "../views/Login.vue"
import Profile from "../views/Profile.vue"
import Groups from "../views/Groups.vue"

const routes = [
  { path: "/", redirect: "/login" },
  { path: "/login", name: "Login", component: Login },
  { path: "/timeline", name: "Timeline", component: Timeline, meta: { requiresAuth: true } },
  { path: "/upload", name: "Upload", component: Upload, meta: { requiresAuth: true } },
  { path: "/records", name: "Records", component: Records, meta: { requiresAuth: true } },
  { path: "/profile", name: "Profile", component: Profile, meta: { requiresAuth: true } },
  { path: "/groups", name: "Groups", component: Groups, meta: { requiresAuth: true } }
]

export default routes