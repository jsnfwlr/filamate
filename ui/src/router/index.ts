import { createRouter, createWebHistory } from 'vue-router'
import HomeView from '../views/dashboard.vue'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      // landing page / dashboard
      path: '/',
      name: 'dashboard',
      component: HomeView,
      meta: { title: 'Filamate - Dashboard' },
    },

    {
      // spools
      path: '/spools',
      name: 'spools',
      component: () => import('../views/spools.vue'),
      meta: { title: 'Filamate - Spools' },
    },
    {
      // brands
      path: '/brands',
      name: 'brands',
      component: () => import('../views/brands.vue'),
      meta: { title: 'Filamate - Brands' },
    },
    {
      // colors
      path: '/colors',
      name: 'colors',
      component: () => import('../views/colors.vue'),
      meta: { title: 'Filamate - Colors' },
    },
    {
      // locations
      path: '/locations',
      name: 'locations',
      component: () => import('../views/locations.vue'),
      meta: { title: 'Filamate - Locations' },
    },
    {
      // materials
      path: '/materials',
      name: 'materials',
      component: () => import('../views/materials.vue'),
      meta: { title: 'Filamate - Materials' },
    },
    {
      // stores
      path: '/stores',
      name: 'stores',
      component: () => import('../views/stores.vue'),
      meta: { title: 'Filamate - Stores' },
    },
  ],
})

export default router
