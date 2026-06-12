import { useUserStore } from '@/stores/user'

// v-permission="'system:user:add'"
export const permission = {
  mounted(el, binding) {
    const perm = binding.value
    if (!perm) return
    const userStore = useUserStore()
    if (!userStore.hasPermission(perm)) {
      el.parentNode?.removeChild(el)
    }
  },
}
