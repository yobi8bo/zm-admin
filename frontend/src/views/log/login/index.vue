<template>
  <div class="page">
    <a-card :bordered="false">
      <div class="toolbar">
        <a-popconfirm title="确定清空全部登录日志吗？" @confirm="handleClear">
          <a-button danger><delete-outlined />清空日志</a-button>
        </a-popconfirm>
        <a-button @click="loadData"><reload-outlined /></a-button>
      </div>

      <a-table
        :columns="columns"
        :data-source="tableData"
        :loading="loading"
        :pagination="pagination"
        row-key="id"
        size="middle"
        @change="handleTableChange"
      >
        <template #bodyCell="{ column, record }">
          <template v-if="column.key === 'status'">
            <a-tag :color="record.status === 1 ? 'success' : 'error'">
              {{ record.status === 1 ? '成功' : '失败' }}
            </a-tag>
          </template>
        </template>
      </a-table>
    </a-card>
  </div>
</template>

<script setup>
  import { ref, reactive, onActivated } from 'vue'
  import { message } from 'ant-design-vue'
  import { logApi } from '@/api/log'

  const loading = ref(false)
  const tableData = ref([])
  const pagination = reactive({ current: 1, pageSize: 20, total: 0, showTotal: (t) => `共 ${t} 条` })

  const columns = [
    { title: '用户名', dataIndex: 'username', key: 'username', width: 120 },
    { title: 'IP地址', dataIndex: 'ip', key: 'ip', width: 140 },
    { title: '归属地', dataIndex: 'location', key: 'location', width: 120 },
    { title: '浏览器', dataIndex: 'browser', key: 'browser' },
    { title: '操作系统', dataIndex: 'os', key: 'os', width: 120 },
    { title: '状态', key: 'status', width: 80 },
    { title: '失败原因', dataIndex: 'message', key: 'message' },
    { title: '登录时间', dataIndex: 'created_at', key: 'created_at', width: 170 },
  ]

  async function loadData() {
    loading.value = true
    try {
      const data = await logApi.loginList({ page: pagination.current, page_size: pagination.pageSize })
      tableData.value = data.list
      pagination.total = data.total
    } finally {
      loading.value = false
    }
  }

  function handleTableChange(pag) { pagination.current = pag.current; loadData() }

  async function handleClear() {
    await logApi.clearLoginLog()
    message.success('清空成功')
    loadData()
  }

  onActivated(loadData)
</script>

<style scoped>
  .page { padding: 0; }
  .toolbar { display: flex; justify-content: space-between; align-items: center; margin-bottom: 12px; }
</style>
