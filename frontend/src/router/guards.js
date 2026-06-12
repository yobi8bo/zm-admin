import { router } from './index'
import { getToken } from '@/utils/auth'
import { useAuthStore } from '@/stores/auth'
import { useUserStore } from '@/stores/user'
import { useMenuStore } from '@/stores/menu'

const whiteList = ['/login']

router.beforeEach(async (to) => {
  const token = getToken()

  if (!token) {
    if (whiteList.includes(to.path)) return true
    return { path: '/login', query: { redirect: to.fullPath } }
  }

  if (to.path === '/login') {
    return { path: '/' }
  }

  const menuStore = useMenuStore()
  if (!menuStore.routesAdded) {
    const userStore = useUserStore()
    try {
      await Promise.all([userStore.fetchUserInfo(), userStore.fetchPermissions()])
      await menuStore.generateRoutes()
      // 动态路由注册后重新导航，确保路由匹配
      return { ...to, replace: true }
    } catch {
      const authStore = useAuthStore()
      await authStore.logout()
      return { path: '/login' }
    }
  }

  return true
})
