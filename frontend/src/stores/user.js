import { defineStore } from 'pinia'
import { userApi } from '@/api/user'

export const useUserStore = defineStore('user', {
  state: () => ({
    info: null,
    permissions: [],
  }),

  getters: {
    hasPermission: (state) => (perm) => state.permissions.includes(perm),
  },

  actions: {
    async fetchUserInfo() {
      const data = await userApi.getMe()
      this.info = data
    },

    async fetchPermissions() {
      const data = await userApi.getMyPermissions()
      this.permissions = data || []
    },
  },
})
