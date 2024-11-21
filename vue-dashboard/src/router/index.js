import { createRouter, createWebHistory } from 'vue-router'
import LiveAppView from '../views/LiveAppView.vue'
import DeviceView from '@/views/DeviceView.vue'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    { path: '/', name: 'home', component: LiveAppView, },

    { path: '/device', name: 'device', component: DeviceView, },


    // {
    //   path: '/device', component: AuthView, children: [
    //     { path: 'login', name: 'viewLogin', component: LoginView },
    //     { path: 'reset', name: 'viewReset', component: ResetView },
    //     { path: 'logout', name: 'viewLogout', component: LogoutView },
    //   ]
    // },



    {
      path: '/about', name: 'about',
      // route level code-splitting
      // this generates a separate chunk (About.[hash].js) for this route
      // which is lazy-loaded when the route is visited.
      component: () => import('../views/AboutView.vue'),
    },
  ],
})

export default router
