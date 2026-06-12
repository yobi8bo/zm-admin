<template>
  <div class="profile-page">
    <a-row :gutter="16">
      <!-- 左侧：头像 + 基本信息卡片 -->
      <a-col :span="8">
        <a-card :bordered="false" class="info-card">
          <div class="avatar-section">
            <a-upload
              :show-upload-list="false"
              :before-upload="beforeUpload"
              accept="image/*"
            >
              <div class="avatar-wrap">
                <a-avatar :size="90" :src="userStore.info?.avatar" class="avatar">
                  {{ userStore.info?.nickname?.[0] || 'U' }}
                </a-avatar>
                <div class="avatar-mask">
                  <camera-outlined />
                  <span>更换头像</span>
                </div>
              </div>
            </a-upload>
            <div class="user-name">{{ userStore.info?.nickname }}</div>
            <div class="user-dept">{{ userStore.info?.dept_name || '暂无部门' }}</div>
          </div>

          <a-divider />

          <div class="meta-list">
            <div class="meta-item">
              <user-outlined class="meta-icon" />
              <span class="meta-label">用户名</span>
              <span class="meta-value">{{ userStore.info?.username }}</span>
            </div>
            <div class="meta-item">
              <phone-outlined class="meta-icon" />
              <span class="meta-label">手机号</span>
              <span class="meta-value">{{ userStore.info?.phone || '未填写' }}</span>
            </div>
            <div class="meta-item">
              <mail-outlined class="meta-icon" />
              <span class="meta-label">邮箱</span>
              <span class="meta-value">{{ userStore.info?.email || '未填写' }}</span>
            </div>
            <div class="meta-item">
              <clock-circle-outlined class="meta-icon" />
              <span class="meta-label">最后登录</span>
              <span class="meta-value">{{ formatDate(userStore.info?.last_login) }}</span>
            </div>
          </div>
        </a-card>
      </a-col>

      <!-- 右侧：Tab 编辑区 -->
      <a-col :span="16">
        <a-card :bordered="false">
          <a-tabs v-model:activeKey="activeTab">
            <!-- 基本信息 -->
            <a-tab-pane key="info" tab="基本信息">
              <a-form
                :model="infoForm"
                :rules="infoRules"
                ref="infoFormRef"
                :label-col="{ span: 5 }"
                :wrapper-col="{ span: 16 }"
                style="margin-top: 12px"
              >
                <a-form-item label="用户昵称" name="nickname">
                  <a-input v-model:value="infoForm.nickname" placeholder="请输入昵称" />
                </a-form-item>
                <a-form-item label="手机号" name="phone">
                  <a-input v-model:value="infoForm.phone" placeholder="请输入手机号" />
                </a-form-item>
                <a-form-item label="邮箱" name="email">
                  <a-input v-model:value="infoForm.email" placeholder="请输入邮箱" />
                </a-form-item>
                <a-form-item label="性别" name="gender">
                  <a-radio-group v-model:value="infoForm.gender">
                    <a-radio :value="0">保密</a-radio>
                    <a-radio :value="1">男</a-radio>
                    <a-radio :value="2">女</a-radio>
                  </a-radio-group>
                </a-form-item>
                <a-form-item :wrapper-col="{ offset: 5, span: 16 }">
                  <a-button type="primary" :loading="infoLoading" @click="handleInfoSubmit">
                    保存修改
                  </a-button>
                </a-form-item>
              </a-form>
            </a-tab-pane>

            <!-- 修改密码 -->
            <a-tab-pane key="password" tab="修改密码">
              <a-form
                :model="pwdForm"
                :rules="pwdRules"
                ref="pwdFormRef"
                :label-col="{ span: 5 }"
                :wrapper-col="{ span: 16 }"
                style="margin-top: 12px"
              >
                <a-form-item label="当前密码" name="old_password">
                  <a-input-password v-model:value="pwdForm.old_password" placeholder="请输入当前密码" />
                </a-form-item>
                <a-form-item label="新密码" name="new_password">
                  <a-input-password v-model:value="pwdForm.new_password" placeholder="请输入新密码，不少于6位" />
                </a-form-item>
                <a-form-item label="确认新密码" name="confirm_password">
                  <a-input-password v-model:value="pwdForm.confirm_password" placeholder="请再次输入新密码" />
                </a-form-item>
                <a-form-item :wrapper-col="{ offset: 5, span: 16 }">
                  <a-button type="primary" :loading="pwdLoading" @click="handlePwdSubmit">
                    修改密码
                  </a-button>
                </a-form-item>
              </a-form>
            </a-tab-pane>
          </a-tabs>
        </a-card>
      </a-col>
    </a-row>
  </div>
</template>

<script setup>
  import { ref, reactive, onMounted } from 'vue'
  import { message } from 'ant-design-vue'
  import dayjs from 'dayjs'
  import { userApi } from '@/api/user'
  import { useUserStore } from '@/stores/user'

  const userStore = useUserStore()
  const activeTab = ref('info')

  // 基本信息表单
  const infoFormRef = ref()
  const infoLoading = ref(false)
  const infoForm = reactive({ nickname: '', phone: '', email: '', gender: 0 })
  const infoRules = {
    nickname: [{ required: true, message: '请输入昵称' }],
    email: [{ type: 'email', message: '邮箱格式不正确' }],
  }

  // 修改密码表单
  const pwdFormRef = ref()
  const pwdLoading = ref(false)
  const pwdForm = reactive({ old_password: '', new_password: '', confirm_password: '' })
  const pwdRules = {
    old_password: [{ required: true, message: '请输入当前密码' }],
    new_password: [{ required: true, min: 6, message: '新密码不少于6位' }],
    confirm_password: [
      { required: true, message: '请确认新密码' },
      {
        validator: (_, value) =>
          value === pwdForm.new_password ? Promise.resolve() : Promise.reject('两次密码不一致'),
      },
    ],
  }

  function formatDate(date) {
    if (!date) return '从未登录'
    return dayjs(date).format('YYYY-MM-DD HH:mm:ss')
  }

  function beforeUpload(file) {
    const isImage = file.type.startsWith('image/')
    const isLt2M = file.size / 1024 / 1024 < 2
    if (!isImage) { message.error('只能上传图片文件'); return false }
    if (!isLt2M) { message.error('图片不能超过 2MB'); return false }
    // TODO: 接入 rustfs 后实现上传
    message.info('文件存储暂未配置')
    return false
  }

  async function handleInfoSubmit() {
    await infoFormRef.value.validate()
    infoLoading.value = true
    try {
      await userApi.updateMe(infoForm)
      await userStore.fetchUserInfo()
      message.success('保存成功')
    } finally {
      infoLoading.value = false
    }
  }

  async function handlePwdSubmit() {
    await pwdFormRef.value.validate()
    pwdLoading.value = true
    try {
      await userApi.updateMyPassword({
        old_password: pwdForm.old_password,
        new_password: pwdForm.new_password,
      })
      message.success('密码修改成功，请重新登录')
      pwdFormRef.value.resetFields()
    } finally {
      pwdLoading.value = false
    }
  }

  onMounted(() => {
    if (userStore.info) {
      infoForm.nickname = userStore.info.nickname
      infoForm.phone = userStore.info.phone || ''
      infoForm.email = userStore.info.email || ''
      infoForm.gender = userStore.info.gender
    }
  })
</script>

<style scoped>
  .profile-page {
    max-width: 1000px;
  }

  .info-card {
    border-radius: 8px;
  }

  .avatar-section {
    display: flex;
    flex-direction: column;
    align-items: center;
    padding: 16px 0 8px;
    gap: 10px;
  }

  .avatar-wrap {
    position: relative;
    cursor: pointer;
    border-radius: 50%;
    overflow: hidden;
  }

  .avatar {
    font-size: 32px;
    background: var(--primary-color);
    display: block;
  }

  .avatar-mask {
    position: absolute;
    inset: 0;
    background: rgba(0, 0, 0, 0.45);
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    gap: 4px;
    color: #fff;
    font-size: 12px;
    opacity: 0;
    transition: opacity 0.2s;
    border-radius: 50%;
  }

  .avatar-wrap:hover .avatar-mask {
    opacity: 1;
  }

  .user-name {
    font-size: 18px;
    font-weight: 600;
    color: var(--text-primary);
  }

  .user-dept {
    font-size: 13px;
    color: var(--text-secondary);
  }

  .meta-list {
    display: flex;
    flex-direction: column;
    gap: 14px;
    padding: 0 8px;
  }

  .meta-item {
    display: flex;
    align-items: center;
    gap: 8px;
    font-size: 13px;
  }

  .meta-icon {
    color: var(--text-secondary);
    font-size: 15px;
    width: 16px;
  }

  .meta-label {
    color: var(--text-secondary);
    width: 56px;
    flex-shrink: 0;
  }

  .meta-value {
    color: var(--text-primary);
    flex: 1;
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
  }
</style>
