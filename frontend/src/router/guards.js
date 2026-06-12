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
      // 初次刷新动态页面时，to.name 可能已经被匹配为 NotFound。
      // 仅使用原始 URL 重新导航，确保新注册的动态路由重新参与匹配。
      return { path: to.fullPath, replace: true }
    } catch {
      const authStore = useAuthStore()
      await authStore.logout()
      return { path: '/login' }
    }
  }

  return true
})
