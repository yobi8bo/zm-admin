import { defineStore } from 'pinia'

export const useAppStore = defineStore('app', {
  state: () => ({
    sidebarCollapsed: false,
    // 多标签页
    tabs: [{ key: '/dashboard', title: '首页', closable: false }],
    activeTab: '/dashboard',
  }),

  actions: {
    toggleSidebar() {
      this.sidebarCollapsed = !this.sidebarCollapsed
    },

    addTab(route) {
      const exists = this.tabs.find((t) => t.key === route.path)
      if (!exists) {
        this.tabs.push({
          key: route.path,
          title: route.meta?.title || route.name,
          closable: route.path !== '/dashboard',
        })
      }
      this.activeTab = route.path
    },

    removeTab(key) {
      const idx = this.tabs.findIndex((t) => t.key === key)
      if (idx === -1) return

      this.tabs.splice(idx, 1)

      if (this.activeTab === key) {
        const next = this.tabs[idx] || this.tabs[idx - 1]
        if (next) this.activeTab = next.key
      }
    },

    closeOtherTabs(key) {
      this.tabs = this.tabs.filter((t) => !t.closable || t.key === key)
      this.activeTab = key
    },

    closeAllTabs() {
      this.tabs = this.tabs.filter((t) => !t.closable)
      this.activeTab = '/dashboard'
    },
  },
})
