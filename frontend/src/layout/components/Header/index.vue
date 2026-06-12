<template>
  <div class="header-inner">
    <!-- 左侧：折叠按钮 + 面包屑 -->
    <div class="header-left">
      <menu-fold-outlined
        v-if="!appStore.sidebarCollapsed"
        class="collapse-btn"
        @click="appStore.toggleSidebar"
      />
      <menu-unfold-outlined v-else class="collapse-btn" @click="appStore.toggleSidebar" />

      <a-breadcrumb class="breadcrumb">
        <a-breadcrumb-item v-for="item in breadcrumbs" :key="item.path">
          {{ item.title }}
        </a-breadcrumb-item>
      </a-breadcrumb>
    </div>

    <!-- 右侧：操作区 -->
    <div class="header-right">
      <!-- 全屏 -->
      <a-tooltip title="全屏">
        <fullscreen-outlined v-if="!isFullscreen" class="action-icon" @click="toggleFullscreen" />
        <fullscreen-exit-outlined v-else class="action-icon" @click="toggleFullscreen" />
      </a-tooltip>

      <!-- 用户头像下拉 -->
      <a-dropdown>
        <div class="user-info">
          <a-avatar size="small" :src="userStore.info?.avatar" style="background: #1677ff">
            {{ userStore.info?.nickname?.[0] || 'U' }}
          </a-avatar>
          <span class="username">{{ userStore.info?.nickname || userStore.info?.username }}</span>
        </div>
        <template #overlay>
          <a-menu>
            <a-menu-item key="profile" @click="$router.push('/profile')">
              <user-outlined />
              个人中心
            </a-menu-item>
            <a-menu-divider />
            <a-menu-item key="logout" @click="handleLogout">
              <logout-outlined />
              退出登录
            </a-menu-item>
          </a-menu>
        </template>
      </a-dropdown>
    </div>
  </div>
</template>

<script setup>
  import { ref, computed } from 'vue'
  import { useRoute, useRouter } from 'vue-router'
  import { useAppStore } from '@/stores/app'
  import { useUserStore } from '@/stores/user'
  import { useAuthStore } from '@/stores/auth'

  const route = useRoute()
  const router = useRouter()
  const appStore = useAppStore()
  const userStore = useUserStore()
  const authStore = useAuthStore()

  const isFullscreen = ref(false)

  const breadcrumbs = computed(() => {
    return route.matched
      .filter((r) => r.meta?.title)
      .map((r) => ({ path: r.path, title: r.meta.title }))
  })

  function toggleFullscreen() {
    if (!document.fullscreenElement) {
      document.documentElement.requestFullscreen()
      isFullscreen.value = true
    } else {
      document.exitFullscreen()
      isFullscreen.value = false
    }
  }

  async function handleLogout() {
    await authStore.logout()
    router.push('/login')
  }
</script>

<style scoped>
  .header-inner {
    display: flex;
    align-items: center;
    justify-content: space-between;
    height: 100%;
    padding: 0 16px;
  }

  .header-left {
    display: flex;
    align-items: center;
    gap: 12px;
  }

  .collapse-btn {
    font-size: 18px;
    cursor: pointer;
    color: var(--text-secondary);
    padding: 4px;
    border-radius: 4px;
    transition: all 0.2s;
  }

  .collapse-btn:hover {
    color: var(--primary-color);
    background: rgba(22, 119, 255, 0.06);
  }

  .breadcrumb {
    font-size: 14px;
  }

  .header-right {
    display: flex;
    align-items: center;
    gap: 16px;
  }

  .action-icon {
    font-size: 17px;
    cursor: pointer;
    color: var(--text-secondary);
    transition: color 0.2s;
  }

  .action-icon:hover {
    color: var(--primary-color);
  }

  .user-info {
    display: flex;
    align-items: center;
    gap: 8px;
    cursor: pointer;
    padding: 4px 8px;
    border-radius: var(--border-radius);
    transition: background 0.2s;
  }

  .user-info:hover {
    background: rgba(0, 0, 0, 0.04);
  }

  .username {
    font-size: 14px;
    color: var(--text-primary);
    max-width: 100px;
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
  }
</style>
