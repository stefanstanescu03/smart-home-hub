import { createRouter, createWebHistory } from "vue-router";

import Home from "../pages/Home.vue";
import Login from "../pages/Login.vue";
import Signup from "../pages/Signup.vue";
import Other from "../pages/Other.vue";
import Account from "../pages/Account.vue";
import Dashboards from "../pages/Dashboards.vue";
import Dashboard from "../pages/Dashboard.vue";
import PublicDashboard from "../pages/PublicDashboard.vue";
import Alert from "../pages/Alert.vue";
import Admin from "../pages/Admin.vue";
import Automations from "../pages/Automations.vue";
import Analyze from "../pages/Analyze.vue";

const routes = [
  { path: "/", component: Home },
  { path: "/login", component: Login },
  { path: "/signup", component: Signup },
  { path: "/other", component: Other },
  { path: "/account", component: Account },
  { path: "/dashboards", component: Dashboards },
  { path: "/dashboard/:id", component: Dashboard },
  { path: "/dashboard/public/:id", component: PublicDashboard },
  { path: "/alerts/:id", component: Alert },
  { path: "/admin", component: Admin },
  { path: "/automations", component: Automations },
  { path: "/analyze/:id", component: Analyze },
];

export default createRouter({
  history: createWebHistory(),
  routes,
});
