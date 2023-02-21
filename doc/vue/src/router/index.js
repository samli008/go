import Vue from 'vue'
import VueRouter from 'vue-router'
import labPage from '../components/labPage'
import labEdit from '../components/labEdit'
import labCreate from '../components/labCreate'
import labView from '../components/labView'

Vue.use(VueRouter)

const routes = [
    { path: '/lab', component: labPage },
    { path: '/labEdit/:id', component: labEdit },
    { path: '/labCreate', component: labCreate },
    { path: '/labView/:id', component: labView },
]

const router = new VueRouter({
    routes
})

export default router