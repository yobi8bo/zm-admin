<template>
  <div class="dashboard">
    <a-row :gutter="[16, 16]">
      <a-col :span="6" v-for="card in statCards" :key="card.title">
        <a-card class="stat-card" :bordered="false" :loading="loading">
          <div class="stat-inner">
            <div class="stat-info">
              <p class="stat-title">{{ card.title }}</p>
              <p class="stat-value">{{ card.value }}</p>
              <p class="stat-desc">{{ card.desc }}</p>
            </div>
            <div class="stat-icon" :style="{ background: card.color + '15' }">
              <component :is="card.icon" :style="{ color: card.color, fontSize: '28px' }" />
            </div>
          </div>
        </a-card>
      </a-col>
    </a-row>

    <a-row :gutter="[16, 16]" style="margin-top: 16px">
      <a-col :span="24">
        <a-card title="系统信息" :bordered="false">
          <a-descriptions :column="3">
            <a-descriptions-item label="系统名称">栈序管理平台</a-descriptions-item>
            <a-descriptions-item label="版本">v1.0.0</a-descriptions-item>
            <a-descriptions-item label="后端框架">Go + Gin</a-descriptions-item>
            <a-descriptions-item label="前端框架">Vue 3 + Ant Design Vue</a-descriptions-item>
            <a-descriptions-item label="数据库">MySQL 8.0</a-descriptions-item>
            <a-descriptions-item label="缓存">Redis</a-descriptions-item>
          </a-descriptions>
        </a-card>
      </a-col>
    </a-row>
  </div>
</template>

<script setup>
  import { ref, reactive, onMounted } from 'vue'
  import {
    TeamOutlined,
    SafetyCertificateOutlined,
    AppstoreOutlined,
    ApartmentOutlined,
  } from '@ant-design/icons-vue'
  import { userApi } from '@/api/user'
  import { roleApi } from '@/api/role'
  import { menuApi } from '@/api/menu'
  import { deptApi } from '@/api/dept'

  const loading = ref(false)

  const statCards = reactive([
    { title: '用户总数', value: 0, desc: '系统注册用户', icon: TeamOutlined, color: '#1677ff' },
    { title: '角色数量', value: 0, desc: '权限角色配置', icon: SafetyCertificateOutlined, color: '#52c41a' },
    { title: '菜单数量', value: 0, desc: '系统功能菜单', icon: AppstoreOutlined, color: '#faad14' },
    { title: '部门数量', value: 0, desc: '组织架构部门', icon: ApartmentOutlined, color: '#722ed1' },
  ])

  async function loadStats() {
    loading.value = true
    try {
      const [userData, roleData, menuData, deptData] = await Promise.all([
        userApi.list({ page: 1, page_size: 1 }),
        roleApi.list({ page: 1, page_size: 1 }),
        menuApi.list(),
        deptApi.list(),
      ])
      statCards[0].value = userData.total
      statCards[1].value = roleData.total
      statCards[2].value = countMenus(menuData)
      statCards[3].value = countDepts(deptData)
    } finally {
      loading.value = false
    }
  }

  // 递归统计菜单总数
  function countMenus(menus) {
    if (!menus || !menus.length) return 0
    return menus.reduce((sum, m) => sum + 1 + countMenus(m.children), 0)
  }

  // 递归统计部门总数
  function countDepts(depts) {
    if (!depts || !depts.length) return 0
    return depts.reduce((sum, d) => sum + 1 + countDepts(d.children), 0)
  }

  onMounted(loadStats)
</script>

<style scoped>
  .dashboard {
    padding: 0;
  }

  .stat-card {
    border-radius: 8px;
    box-shadow: 0 1px 4px rgba(0, 0, 0, 0.06);
  }

  .stat-inner {
    display: flex;
    align-items: center;
    justify-content: space-between;
  }

  .stat-title {
    font-size: 13px;
    color: var(--text-secondary);
    margin-bottom: 8px;
  }

  .stat-value {
    font-size: 28px;
    font-weight: 600;
    color: var(--text-primary);
    margin-bottom: 4px;
  }

  .stat-desc {
    font-size: 12px;
    color: var(--text-disabled);
  }

  .stat-icon {
    width: 56px;
    height: 56px;
    border-radius: 12px;
    display: flex;
    align-items: center;
    justify-content: center;
  }
</style>
