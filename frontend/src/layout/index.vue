<template>
  <a-layout class="layout">
    <!-- 侧边栏 -->
    <a-layout-sider
      v-model:collapsed="appStore.sidebarCollapsed"
      :width="220"
      :collapsed-width="48"
      :trigger="null"
      collapsible
      class="sider"
    >
      <div class="logo">
        <img src="@/assets/logo.svg" alt="logo" class="logo-img" />
        <span v-show="!appStore.sidebarCollapsed" class="logo-title">栈序管理平台</span>
      </div>
      <Sidebar />
    </a-layout-sider>

    <a-layout class="main">
      <!-- 顶部导航 -->
      <a-layout-header class="header">
        <Header />
      </a-layout-header>

      <!-- 多标签页 -->
      <TabsView />

      <!-- 内容区 -->
      <a-layout-content class="content">
        <router-view v-slot="{ Component, route }">
          <transition name="fade" mode="out-in">
            <keep-alive :max="10">
              <component :is="Component" :key="route.path" />
            </keep-alive>
          </transition>
        </router-view>
      </a-layout-content>
    </a-layout>
  </a-layout>
</template>

<script setup>
  import { useAppStore } from '@/stores/app'
  import Sidebar from './components/Sidebar/index.vue'
  import Header from './components/Header/index.vue'
  import TabsView from './components/TabsView/index.vue'

  const appStore = useAppStore()
</script>

<style scoped>
  .layout {
    height: 100vh;
  }

  .sider {
    background: var(--sider-bg) !important;
    overflow: hidden;
    position: fixed;
    left: 0;
    top: 0;
    bottom: 0;
    z-index: 100;
    box-shadow: 2px 0 8px rgba(0, 0, 0, 0.15);
  }

  .logo {
    height: var(--header-height);
    display: flex;
    align-items: center;
    padding: 0 16px;
    gap: 10px;
    border-bottom: 1px solid rgba(255, 255, 255, 0.08);
    overflow: hidden;
    white-space: nowrap;
  }

  .logo-img {
    width: 28px;
    height: 28px;
    flex-shrink: 0;
  }

  .logo-title {
    color: #fff;
    font-size: 16px;
    font-weight: 600;
    letter-spacing: 0.5px;
  }

  .main {
    margin-left: var(--sider-width);
    transition: margin-left 0.2s;
    min-height: 100vh;
    display: flex;
    flex-direction: column;
  }

  .header {
    height: var(--header-height);
    background: #fff;
    padding: 0;
    position: sticky;
    top: 0;
    z-index: 99;
    box-shadow: 0 1px 4px rgba(0, 0, 0, 0.08);
    line-height: var(--header-height);
  }

  .content {
    flex: 1;
    padding: var(--spacing-lg);
    overflow: auto;
  }

  .fade-enter-active,
  .fade-leave-active {
    transition: opacity 0.15s ease;
  }
  .fade-enter-from,
  .fade-leave-to {
    opacity: 0;
  }
</style>
