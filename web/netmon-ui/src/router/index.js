import { createRouter, createWebHistory } from 'vue-router'
import Home from '../views/Home.vue'
// import Packets from '../views/Packets.vue'
// import Hosts from '../views/Hosts.vue'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      name: 'home',
      component: Home
    },
    // {
    //   path: '/packets',
    //   name: 'packets',
    //   component: Packets
    // },
    // {
    //   path: '/hosts',
    //   name: 'hosts',
    //   component: Hosts
    // }
  ]
})

export default router
