import { createRouter, createWebHistory } from 'vue-router'
import LiveAppView from '../views/LiveAppView.vue'
import DeviceView from '@/views/DeviceView.vue'
import UserView from '@/views/UserView.vue'
import AuthView from '@/views/AuthView.vue'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    { path: '/', name: 'homeView', component: LiveAppView, },
    { path: '/user', name: 'userView', component: UserView, },
    { path: '/device', name: 'deviceView', component: DeviceView, },
    { path: '/login', name: 'loginView', component: AuthView, },
    { path: '/forgot-password', name: 'forgotPasswordView', component: AuthView, },


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
