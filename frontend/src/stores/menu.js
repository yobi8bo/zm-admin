import { defineStore } from 'pinia'
import { userApi } from '@/api/user'
import { router } from '@/router'
import Layout from '@/layout/index.vue'

// 预加载所有页面组件
const modules = import.meta.glob('@/views/**/*.vue')

function resolveComponent(component) {
  if (!component || component === 'Layout') return null
  const key = `/src/views/${component}.vue`
  return modules[key] || null
}

// 将菜单树转换为 vue-router 路由配置
// 目录(type=1)：注册为带 Layout 的嵌套路由父级
// 菜单(type=2)：注册为实际页面路由
function buildRoutes(menus) {
  const routes = []

  menus.forEach((menu) => {
    // 目录：作为二级路由父级，children 是真正的页面
    if (menu.children && menu.children.length) {
      const children = menu.children
        .filter((child) => resolveComponent(child.component))
        .map((child) => ({
          path: child.path,
          name: child.name,
          component: resolveComponent(child.component),
          meta: { title: child.name, icon: child.icon, menuId: child.id },
        }))

      if (children.length) {
        routes.push({
          path: menu.path,
          component: Layout,
          meta: { title: menu.name, icon: menu.icon },
          children,
        })
      }
    } else {
      // 独立菜单页面（无子项）
      const component = resolveComponent(menu.component)
      if (component) {
        routes.push({
          path: menu.path,
          name: menu.name,
          component,
          meta: { title: menu.name, icon: menu.icon, menuId: menu.id },
        })
      }
    }
  })

  return routes
}

export const useMenuStore = defineStore('menu', {
  state: () => ({
    menus: [],           // 菜单树，用于侧边栏渲染
    routesAdded: false,
  }),

  actions: {
    async generateRoutes() {
      if (this.routesAdded) return

      const menus = await userApi.getMyMenus()
      this.menus = menus || []

      const dynamicRoutes = buildRoutes(this.menus)
      dynamicRoutes.forEach((route) => {
        router.addRoute(route)
      })

      this.routesAdded = true
    },

    $reset() {
      this.menus = []
      this.routesAdded = false
    },
  },
})
