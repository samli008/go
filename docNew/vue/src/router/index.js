import Vue from 'vue'
import VueRouter from 'vue-router'
import linuxPage from '../components/linuxPage'
import pageEdit from '../components/pageEdit'
import pageCreate from '../components/pageCreate'
import pageView from '../components/pageView'
import netappPage from '../components/netappPage'
import dellPage from '../components/dellPage'
import privatePage from '../components/privatePage'
import vmwarePage from '../components/vmwarePage'
import windowsPage from '../components/windowsPage'

Vue.use(VueRouter)

const routes = [
    { path: '/linux', name: "linux", component: linuxPage },
    { path: '/netapp', name: "netappPage", component: netappPage },
    { path: '/dell', name: "dellPage", component: dellPage },
    { path: '/vmware', name: "vmwarePage", component: vmwarePage },
    { path: '/windows', name: "windowsPage", component: windowsPage },
    { path: '/private', name: "privateView", component: privatePage },
    { path: '/pageEdit/:id', name: "pageEdit", component: pageEdit },
    { path: '/pageCreate', name: "pageCreate", component: pageCreate },
    { path: '/pageView/:id', name: "pageView", component: pageView },
]

const router = new VueRouter({
    routes
})

export default router