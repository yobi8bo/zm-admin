<template>
  <div class="page">
    <a-card :bordered="false">
      <div class="toolbar">
        <a-popconfirm title="确定清空全部操作日志吗？" @confirm="handleClear">
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
        :scroll="{ x: 1200 }"
      >
        <template #bodyCell="{ column, record }">
          <template v-if="column.key === 'status'">
            <a-tag :color="isOperationSuccess(record) ? 'success' : 'error'">
              {{ isOperationSuccess(record) ? '成功' : '失败' }}
            </a-tag>
          </template>
          <template v-if="column.key === 'latency'">
            {{ record.latency }}ms
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
    { title: '用户名', dataIndex: 'username', key: 'username', width: 100 },
    { title: '模块', dataIndex: 'module', key: 'module', width: 100 },
    { title: '操作', dataIndex: 'action', key: 'action', width: 120 },
    { title: '执行结果', key: 'status', width: 100 },
    { title: '耗时', key: 'latency', width: 80 },
    { title: '时间', dataIndex: 'created_at', key: 'created_at', width: 170 },
  ]

  function isOperationSuccess(record) {
    return record.status === 200 && !record.error
  }

  async function loadData() {
    loading.value = true
    try {
      const data = await logApi.operationList({ page: pagination.current, page_size: pagination.pageSize })
      tableData.value = data.list
      pagination.total = data.total
    } finally {
      loading.value = false
    }
  }

  function handleTableChange(pag) { pagination.current = pag.current; loadData() }

  async function handleClear() {
    await logApi.clearOperationLog()
    message.success('清空成功')
    loadData()
  }

  onActivated(loadData)
</script>

<style scoped>
  .page { padding: 0; }
  .toolbar { display: flex; justify-content: space-between; align-items: center; margin-bottom: 12px; }
</style>
