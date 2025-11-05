import { createRouter, createWebHistory } from "vue-router";

import Home from "../pages/Home.vue";
import Login from "../pages/Login.vue";
import Signup from "../pages/Signup.vue";
import Discovered from "../pages/Discovered.vue";
import Other from "../pages/Other.vue";

const routes = [
  { path: "/", component: Home },
  { path: "/login", component: Login },
  { path: "/signup", component: Signup },
  { path: "/discovered", component: Discovered },
  { path: "/other", component: Other },
];

export default createRouter({
  history: createWebHistory(),
  routes,
});
