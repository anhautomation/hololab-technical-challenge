import { createRouter, createWebHistory } from "vue-router";
import Home from "../pages/Home.vue";
import CreateBot from "../pages/CreateBot.vue";
import Chat from "../pages/Chat.vue";

export default createRouter({
  history: createWebHistory(),
  routes: [
    { path: "/", component: Home },
    { path: "/create", component: CreateBot },
    { path: "/chat/:id", component: Chat, props: true },
  ],
});
