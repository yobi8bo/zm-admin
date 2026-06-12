<template>
  <div class="tabs-bar">
    <a-tabs
      v-model:activeKey="appStore.activeTab"
      type="editable-card"
      hide-add
      size="small"
      class="tabs"
      @change="handleTabChange"
      @edit="handleTabEdit"
    >
      <a-tab-pane
        v-for="tab in appStore.tabs"
        :key="tab.key"
        :tab="tab.title"
        :closable="tab.closable"
      />
    </a-tabs>

    <!-- 右键菜单触发按钮 -->
    <a-dropdown :trigger="['click']">
      <down-outlined class="tabs-action" />
      <template #overlay>
        <a-menu>
          <a-menu-item @click="appStore.closeOtherTabs(appStore.activeTab)">关闭其他</a-menu-item>
          <a-menu-item @click="appStore.closeAllTabs">关闭全部</a-menu-item>
        </a-menu>
      </template>
    </a-dropdown>
  </div>
</template>

<script setup>
  import { watch } from 'vue'
  import { useRouter, useRoute } from 'vue-router'
  import { useAppStore } from '@/stores/app'

  const router = useRouter()
  const route = useRoute()
  const appStore = useAppStore()

  // 路由变化时同步添加 tab
  watch(
    () => route.path,
    () => {
      if (route.meta?.title) {
        appStore.addTab(route)
      }
    },
    { immediate: true }
  )

  function handleTabChange(key) {
    router.push(key)
  }

  function handleTabEdit(key, action) {
    if (action === 'remove') {
      const current = appStore.activeTab
      appStore.removeTab(key)
      if (current === key) {
        router.push(appStore.activeTab)
      }
    }
  }
</script>

<style scoped>
  .tabs-bar {
    display: flex;
    align-items: center;
    background: #fff;
    border-bottom: 1px solid var(--border-color);
    padding: 0 12px 0 0;
    height: var(--tabs-height);
  }

  .tabs {
    flex: 1;
    overflow: hidden;
  }

  :deep(.ant-tabs-nav) {
    margin: 0 !important;
    height: var(--tabs-height);
  }

  :deep(.ant-tabs-nav::before) {
    border: none !important;
  }

  :deep(.ant-tabs-tab) {
    padding: 6px 12px !important;
    font-size: 13px;
    border-radius: 4px 4px 0 0 !important;
  }

  .tabs-action {
    font-size: 13px;
    padding: 4px 6px;
    cursor: pointer;
    color: var(--text-secondary);
    border-radius: 4px;
    transition: background 0.2s;
  }

  .tabs-action:hover {
    background: rgba(0, 0, 0, 0.04);
  }
</style>
