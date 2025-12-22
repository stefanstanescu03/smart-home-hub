import { createRouter, createWebHistory } from "vue-router";

import Home from "../pages/Home.vue";
import Login from "../pages/Login.vue";
import Signup from "../pages/Signup.vue";
import Other from "../pages/Other.vue";
import Account from "../pages/Account.vue";
import Dashboards from "../pages/Dashboards.vue";
import Dashboard from "../pages/Dashboard.vue";
import Alert from "../pages/Alert.vue";

const routes = [
  { path: "/", component: Home },
  { path: "/login", component: Login },
  { path: "/signup", component: Signup },
  { path: "/other", component: Other },
  { path: "/account", component: Account },
  { path: "/dashboards", component: Dashboards },
  { path: "/dashboard/:id", component: Dashboard },
  { path: "/alerts/:id", component: Alert },
];

export default createRouter({
  history: createWebHistory(),
  routes,
});
