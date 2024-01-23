import Vue from 'vue'
import VueRouter from 'vue-router'
import store from '@/store'
import login from "@/views/login/login.vue";
import home from "@/views/home.vue";
import welcome from "@/views/welcome/welcome.vue"
import search from "@/views/search/search.vue";
import photos from "@/views/photos/photos.vue";
import albums from "@/views/albums/albums.vue";
import scene from "@/views/scene/scene.vue"
import label from "@/views/scene/label.vue";
import faces from "@/views/faces/faces.vue";
import cluster from "@/views/faces/cluster.vue";
import collect from "@/views/collect/collect.vue";
import upload from "@/views/upload/upload.vue";
import recycle from "@/views/recycle/recycle.vue";
import about from "@/views/about/about.vue";
import blank from "@/views/about/blank.vue";
import register from "@/views/register/register.vue";

Vue.use(VueRouter)

const routes = [
  {
    path: '/',
    component: login
  },
  {
    name: "login",
    path: "/login",
    component: login,
  },
  {
    name: "register",
    path: "/register",
    component: register,
  },
  {
    name:"home",
    path: "/home",
    component: home,
    children:[
      { name:"welcome",path:"welcome",component: welcome},
      { name:"search",path:"search",component: search},
      { name:"photos",path:"photos",component:photos },
      { name:"collect",path:"collect",component:collect },
      { name:"albums",path:"albums/:albumName",component:albums },
      { name:"scene",path: "scene",component: scene},
      { name:"labels",path:"labels/:labelName",component:label },
      { name:"faces",path: "faces",component: faces},
      { name:"clusters",path:"clusters/:clusterName",component:cluster },
      { name:"upload",path:"upload",component:upload },
      { name:"recycle",path:"recycle",component:recycle },
      { name:"about",path:"about",component:about },
      { name:"blank",path:"blank",component:blank },
    ]
  }
]

const router = new VueRouter({
  mode: 'history',
  base: process.env.BASE_URL,
  routes
})

router.beforeEach((to, from, next) => {
  // 全局路由守卫
  // 根据store中存储的login信息确认是否放行
  // 登录情况
  if (from.name === 'login') {
    if (store.state.isLoggedIn === true) {
      return next()
    }
    else {
      if (to.name === 'register') {
        return next()
      }
    }
  }
  // 直接更改url打开页面的情况
  if (from.name === null) {
    if (store.state.isLoggedIn === true) {
      // 用户现在是登录了的
      if (to.name !== 'login') {
        // 是登录状态，但是用户直接更改url来跳转页面，因为是登录状态，所以去哪都可以
        return next()
      } else {
        // 先跳转回登录界面再跳到首页，因为用户没有登出过
        next()
        return next({ name: 'home' })
      }
    } else {
      if (to.name !== 'login') {
        // 是登出状态，用户直接更改url来跳转到非登录页面，因为是登出状态，所以强制跳转到登录页面
        return next({ name: 'login' })
      } else {
        // 是登出状态，用户直接更改url来跳转到登录页面，因为是登出状态，所以放行
        return next()
      }
    }
  }
  // 从其它页面打算回到登录页面的情况
  if (to.name === 'login' && from.name != null) {
    // 只有通过右上角登出按钮，状态才会变成false，否则不变
    if (store.state.isLoggedIn === false) {
      return next()
    } else {
      // 停留在原来的页面
      return next({ name: from.name })
    }
  }
  // 其它正常情况
  if (to.name !== 'login' && from.name !== 'login' && to.name !== null && from.name !== null) {
    return next()
  }
})

export default router
