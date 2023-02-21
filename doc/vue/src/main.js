import Vue from 'vue'
import App from './App.vue'
import router from './router'
import ElementUI from 'element-ui'
import 'element-ui/lib/theme-chalk/index.css'
import axios from 'axios'
import { Message } from 'element-ui'

Vue.use(ElementUI)

const http = axios.create({
    baseURL: ''
})

// http.interceptors.request.use(config => {
//     if (localStorage.elementToken) {
//         config.headers.Authorization = localStorage.elementToken
//     }
//     return config
// })

http.interceptors.response.use(res => {
    return res
}, err => {
    console.log(err.response)
    Message.error(err.response.data)
})

Vue.prototype.$axios = http
Vue.config.productionTip = false

new Vue({
    router,
    render: h => h(App)
}).$mount('#app')