import Vue from 'vue'
import Vuex from 'vuex'

Vue.use(Vuex)

export default new Vuex.Store({
  state: {
    username: '',
    isLoggedIn: false, // 初始状态为未登录
    token: '',
    albumArray: [], // 这是你的数组状态属性
  },
  mutations: {
    // 定义一个 mutation 来设置登录状态
    SET_LOGIN_STATUS(state, status) {
      state.isLoggedIn = status;
    },
    SET_TOKEN(state, token) {
      state.token = token;
    },
    SET_USER_NAME(state, username) {
      state.username = username;
    },
    SET_Album_Array(state, albumArray) {
      state.albumArray = albumArray;
    },
  },
  actions: {
    // 定义一个 action 来触发设置登录状态的 mutation
    login({ commit }) {
      // 这里可以进行登录逻辑，例如向后端发送登录请求
      // 登录成功后调用 commit('SET_LOGIN_STATUS', true); 来设置登录状态
      // 登录失败则根据需要进行处理
    },
    logout({ commit }) {
      // 这里可以进行注销逻辑，例如向后端发送注销请求
      // 注销成功后调用 commit('SET_LOGIN_STATUS', false); 来设置登录状态
      // 注销失败则根据需要进行处理
    },
  },
  getters: {
    // 定义一个 getter 来获取登录状态
    isLoggedIn: (state) => state.isLoggedIn,
    token: (state) => state.token,
    username: (state) => state.username,
  },
  modules: {
  }
})
