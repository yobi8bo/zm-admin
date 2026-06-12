<template>
  <a-menu
    v-model:selectedKeys="selectedKeys"
    v-model:openKeys="openKeys"
    mode="inline"
    theme="dark"
    class="sidebar-menu"
    :inline-collapsed="appStore.sidebarCollapsed"
  >
    <!-- 首页单独处理 -->
    <a-menu-item key="/dashboard" @click="router.push('/dashboard')">
      <template #icon><home-outlined /></template>
      首页
    </a-menu-item>

    <template v-for="menu in sideMenus" :key="menu.path">
      <!-- 有子菜单的目录 -->
      <a-sub-menu v-if="menu.children?.length" :key="menu.path">
        <template #icon>
          <component :is="getIcon(menu.icon)" />
        </template>
        <template #title>{{ menu.name }}</template>
        <a-menu-item
          v-for="child in menu.children"
          :key="fullPath(menu.path, child.path)"
          @click="router.push(fullPath(menu.path, child.path))"
        >
          <template #icon>
            <component :is="getIcon(child.icon)" />
          </template>
          {{ child.name }}
        </a-menu-item>
      </a-sub-menu>

      <!-- 无子菜单 -->
      <a-menu-item v-else :key="menu.path" @click="router.push(menu.path)">
        <template #icon>
          <component :is="getIcon(menu.icon)" />
        </template>
        {{ menu.name }}
      </a-menu-item>
    </template>
  </a-menu>
</template>

<script setup>
  import { ref, computed, watch } from 'vue'
  import { useRouter, useRoute } from 'vue-router'
  import { useMenuStore } from '@/stores/menu'
  import { useAppStore } from '@/stores/app'
  import * as Icons from '@ant-design/icons-vue'

  const router = useRouter()
  const route = useRoute()
  const menuStore = useMenuStore()
  const appStore = useAppStore()

  const selectedKeys = ref([route.path])
  const openKeys = ref([])

  // 过滤掉首页（单独渲染），只渲染目录和非首页菜单
  const sideMenus = computed(() =>
    menuStore.menus.filter((m) => m.path !== '/dashboard')
  )

  watch(
    () => route.path,
    (path) => {
      selectedKeys.value = [path]
      // 自动展开当前路由对应的父级目录
      const parent = menuStore.menus.find((m) =>
        m.children?.some((c) => fullPath(m.path, c.path) === path)
      )
      if (parent && !openKeys.value.includes(parent.path)) {
        openKeys.value = [...openKeys.value, parent.path]
      }
    },
    { immediate: true }
  )

  function getIcon(iconName) {
    return iconName && Icons[iconName] ? Icons[iconName] : Icons['AppstoreOutlined']
  }

  // 拼接完整路径：/system + user => /system/user
  function fullPath(parent, child) {
    if (child.startsWith('/')) return child
    return `${parent}/${child}`
  }
</script>

<style scoped>
  .sidebar-menu {
    border-right: none;
    overflow-y: auto;
    overflow-x: hidden;
    height: calc(100vh - var(--header-height));
  }

  :deep(.ant-menu-item),
  :deep(.ant-menu-submenu-title) {
    border-radius: 0 !important;
    margin: 0 !important;
    width: 100% !important;
  }

  :deep(.ant-menu-item-selected) {
    background: var(--primary-color) !important;
  }
</style>
