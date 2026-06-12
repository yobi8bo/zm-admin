import { defineStore } from 'pinia'
import { authApi } from '@/api/auth'
import { setToken, setRefreshToken, clearAuth, getToken } from '@/utils/auth'
import { useUserStore } from './user'
import { useMenuStore } from './menu'

export const useAuthStore = defineStore('auth', {
  state: () => ({
    token: getToken() || '',
  }),

  getters: {
    isLoggedIn: (state) => !!state.token,
  },

  actions: {
    async login(loginForm) {
      const data = await authApi.login(loginForm)
      this.token = data.access_token
      setToken(data.access_token)
      setRefreshToken(data.refresh_token)
    },

    async logout() {
      try {
        await authApi.logout()
      } finally {
        this.token = ''
        clearAuth()
        const userStore = useUserStore()
        const menuStore = useMenuStore()
        userStore.$reset()
        menuStore.$reset()
      }
    },
  },
})
