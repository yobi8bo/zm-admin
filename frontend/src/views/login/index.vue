<template>
  <div class="login-page">
    <div class="login-bg" />
    <div class="login-box">
      <div class="login-header">
        <img src="@/assets/logo.svg" class="login-logo" alt="logo" />
        <h1 class="login-title">栈序管理平台</h1>
        <p class="login-desc">企业级中后台管理平台</p>
      </div>

      <a-form :model="form" :rules="rules" ref="formRef" class="login-form" @finish="handleLogin">
        <a-form-item name="username">
          <a-input
            v-model:value="form.username"
            size="large"
            placeholder="请输入用户名"
            allow-clear
          >
            <template #prefix>
              <user-outlined class="form-icon" />
            </template>
          </a-input>
        </a-form-item>

        <a-form-item name="password">
          <a-input-password
            v-model:value="form.password"
            size="large"
            placeholder="请输入密码"
          >
            <template #prefix>
              <lock-outlined class="form-icon" />
            </template>
          </a-input-password>
        </a-form-item>

        <a-form-item name="captcha_code">
          <div class="captcha-row">
            <a-input
              v-model:value="form.captcha_code"
              size="large"
              placeholder="请输入验证码"
              style="flex: 1"
            />
            <img
              :src="captchaImg"
              class="captcha-img"
              alt="验证码"
              title="点击刷新"
              @click="loadCaptcha"
            />
          </div>
        </a-form-item>

        <a-form-item>
          <a-button
            type="primary"
            html-type="submit"
            size="large"
            block
            :loading="loading"
            class="login-btn"
          >
            登录
          </a-button>
        </a-form-item>
      </a-form>
    </div>
  </div>
</template>

<script setup>
  import { ref, onMounted } from 'vue'
  import { useRouter, useRoute } from 'vue-router'
  import { message } from 'ant-design-vue'
  import { authApi } from '@/api/auth'
  import { useAuthStore } from '@/stores/auth'

  const router = useRouter()
  const route = useRoute()
  const authStore = useAuthStore()

  const formRef = ref()
  const loading = ref(false)
  const captchaImg = ref('')
  const captchaId = ref('')

  const form = ref({
    username: '',
    password: '',
    captcha_id: '',
    captcha_code: '',
  })

  const rules = {
    username: [{ required: true, message: '请输入用户名' }],
    password: [{ required: true, message: '请输入密码' }],
    captcha_code: [{ required: true, message: '请输入验证码' }],
  }

  async function loadCaptcha() {
    try {
      const data = await authApi.getCaptcha()
      captchaId.value = data.captcha_id
      captchaImg.value = data.image
      form.value.captcha_id = data.captcha_id
      form.value.captcha_code = ''
    } catch {
      message.error('获取验证码失败')
    }
  }

  async function handleLogin() {
    loading.value = true
    try {
      form.value.captcha_id = captchaId.value
      await authStore.login(form.value)
      const redirect = route.query.redirect || '/dashboard'
      router.push(redirect)
    } catch {
      loadCaptcha()
    } finally {
      loading.value = false
    }
  }

  onMounted(loadCaptcha)
</script>

<style scoped>
  .login-page {
    min-height: 100vh;
    display: flex;
    align-items: center;
    justify-content: center;
    background: linear-gradient(135deg, #e8f4ff 0%, #f5f6fa 50%, #e8f0ff 100%);
    position: relative;
    overflow: hidden;
  }

  .login-bg {
    position: absolute;
    inset: 0;
    background: url("data:image/svg+xml,%3Csvg width='60' height='60' viewBox='0 0 60 60' xmlns='http://www.w3.org/2000/svg'%3E%3Cg fill='none' fill-rule='evenodd'%3E%3Cg fill='%231677ff' fill-opacity='0.03'%3E%3Ccircle cx='30' cy='30' r='20'/%3E%3C/g%3E%3C/g%3E%3C/svg%3E")
      repeat;
  }

  .login-box {
    width: 400px;
    background: #fff;
    border-radius: 12px;
    padding: 48px 40px;
    box-shadow: 0 8px 40px rgba(0, 0, 0, 0.1);
    position: relative;
    z-index: 1;
  }

  .login-header {
    text-align: center;
    margin-bottom: 36px;
  }

  .login-logo {
    width: 52px;
    height: 52px;
    margin-bottom: 12px;
  }

  .login-title {
    font-size: 22px;
    font-weight: 600;
    color: var(--text-primary);
    margin-bottom: 6px;
  }

  .login-desc {
    font-size: 13px;
    color: var(--text-secondary);
  }

  .login-form {
    margin-top: 8px;
  }

  .form-icon {
    color: var(--text-disabled);
  }

  .captcha-row {
    display: flex;
    gap: 12px;
    align-items: center;
  }

  .captcha-img {
    width: 110px;
    height: 40px;
    border-radius: var(--border-radius);
    border: 1px solid #d9d9d9;
    cursor: pointer;
    object-fit: cover;
    flex-shrink: 0;
    transition: opacity 0.2s;
  }

  .captcha-img:hover {
    opacity: 0.8;
  }

  .login-btn {
    height: 42px;
    font-size: 15px;
    margin-top: 4px;
    border-radius: var(--border-radius);
  }
</style>
