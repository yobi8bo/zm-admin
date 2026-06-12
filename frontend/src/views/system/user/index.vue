<template>
  <div class="page">
    <!-- 搜索区 -->
    <a-card :bordered="false" class="search-card">
      <a-form :model="queryForm" layout="inline">
        <a-form-item label="用户名">
          <a-input v-model:value="queryForm.username" placeholder="请输入用户名" allow-clear style="width: 180px" />
        </a-form-item>
        <a-form-item label="手机号">
          <a-input v-model:value="queryForm.phone" placeholder="请输入手机号" allow-clear style="width: 160px" />
        </a-form-item>
        <a-form-item label="状态">
          <a-select v-model:value="queryForm.status" placeholder="全部" allow-clear style="width: 100px">
            <a-select-option :value="1">正常</a-select-option>
            <a-select-option :value="0">禁用</a-select-option>
          </a-select>
        </a-form-item>
        <a-form-item label="部门">
          <a-tree-select
            v-model:value="queryForm.dept_id"
            :tree-data="deptTree"
            :field-names="{ label: 'name', value: 'id', children: 'children' }"
            placeholder="请选择部门"
            allow-clear
            style="width: 180px"
          />
        </a-form-item>
        <a-form-item>
          <a-space>
            <a-button type="primary" @click="handleQuery"><search-outlined />查询</a-button>
            <a-button @click="handleReset">重置</a-button>
          </a-space>
        </a-form-item>
      </a-form>
    </a-card>

    <!-- 表格区 -->
    <a-card :bordered="false" style="margin-top: 12px">
      <!-- 工具栏 -->
      <div class="toolbar">
        <a-space>
          <a-button type="primary" @click="handleCreate"><plus-outlined />新增</a-button>
          <a-button danger :disabled="!selectedRowKeys.length" @click="handleBatchDelete">
            <delete-outlined />批量删除
          </a-button>
        </a-space>
        <a-button @click="loadData"><reload-outlined /></a-button>
      </div>

      <a-table
        :columns="columns"
        :data-source="tableData"
        :loading="loading"
        :pagination="pagination"
        row-key="id"
        :row-selection="rowSelection"
        size="middle"
        @change="handleTableChange"
      >
        <template #bodyCell="{ column, record }">
          <template v-if="column.key === 'avatar'">
            <a-avatar :src="record.avatar" size="small" style="background:#1677ff">
              {{ record.nickname?.[0] || 'U' }}
            </a-avatar>
          </template>
          <template v-if="column.key === 'status'">
            <a-tag :color="record.status === 1 ? 'success' : 'error'">
              {{ record.status === 1 ? '正常' : '禁用' }}
            </a-tag>
          </template>
          <template v-if="column.key === 'action'">
            <a-space size="small">
              <a-button type="link" size="small" @click="handleEdit(record)">编辑</a-button>
              <a-button v-if="!record.is_admin" type="link" size="small" @click="handleAssignRoles(record)">分配角色</a-button>
              <a-button type="link" size="small" @click="handleResetPwd(record)">重置密码</a-button>
              <a-popconfirm v-if="!record.is_admin" title="确定删除该用户吗？" @confirm="handleDelete(record.id)">
                <a-button type="link" size="small" danger>删除</a-button>
              </a-popconfirm>
            </a-space>
          </template>
        </template>
      </a-table>
    </a-card>

    <!-- 新增/编辑弹窗 -->
    <a-modal
      v-model:open="modalVisible"
      :title="modalTitle"
      width="600px"
      @ok="handleSubmit"
      @cancel="handleModalClose"
      :confirm-loading="submitLoading"
    >
      <a-form :model="form" :rules="formRules" ref="formRef" :label-col="{ span: 6 }" :wrapper-col="{ span: 16 }">
        <a-form-item label="所属部门" name="dept_id">
          <a-tree-select
            v-model:value="form.dept_id"
            :tree-data="deptTree"
            :field-names="{ label: 'name', value: 'id', children: 'children' }"
            placeholder="请选择部门"
            allow-clear
          />
        </a-form-item>
        <a-form-item label="用户名" name="username">
          <a-input v-model:value="form.username" placeholder="请输入用户名" :disabled="!!form.id" />
        </a-form-item>
        <a-form-item label="用户昵称" name="nickname">
          <a-input v-model:value="form.nickname" placeholder="请输入昵称" />
        </a-form-item>
        <a-form-item v-if="!form.id" label="密码" name="password">
          <a-input-password v-model:value="form.password" placeholder="请输入密码" />
        </a-form-item>
        <a-form-item label="手机号" name="phone">
          <a-input v-model:value="form.phone" placeholder="请输入手机号" />
        </a-form-item>
        <a-form-item label="邮箱" name="email">
          <a-input v-model:value="form.email" placeholder="请输入邮箱" />
        </a-form-item>
        <a-form-item label="性别" name="gender">
          <a-radio-group v-model:value="form.gender">
            <a-radio :value="0">未知</a-radio>
            <a-radio :value="1">男</a-radio>
            <a-radio :value="2">女</a-radio>
          </a-radio-group>
        </a-form-item>
        <a-form-item label="状态" name="status">
          <a-radio-group v-model:value="form.status">
            <a-radio :value="1">正常</a-radio>
            <a-radio :value="0">禁用</a-radio>
          </a-radio-group>
        </a-form-item>
        <a-form-item label="备注">
          <a-textarea v-model:value="form.remark" :rows="3" placeholder="请输入备注" />
        </a-form-item>
      </a-form>
    </a-modal>

    <!-- 分配角色弹窗 -->
    <a-modal
      v-model:open="roleModalVisible"
      title="分配角色"
      @ok="handleRoleSubmit"
      :confirm-loading="submitLoading"
    >
      <a-checkbox-group v-model:value="selectedRoles" style="display:flex;flex-direction:column;gap:8px">
        <a-checkbox v-for="role in allRoles" :key="role.id" :value="role.id">
          {{ role.name }}（{{ role.code }}）
        </a-checkbox>
      </a-checkbox-group>
    </a-modal>

    <!-- 重置密码弹窗 -->
    <a-modal v-model:open="resetPwdVisible" title="重置密码" @ok="handleResetPwdSubmit" :confirm-loading="submitLoading">
      <a-form :model="resetPwdForm" ref="resetPwdRef" :label-col="{ span: 6 }" :wrapper-col="{ span: 16 }">
        <a-form-item label="新密码" name="password" :rules="[{ required: true, min: 6, message: '密码不少于6位' }]">
          <a-input-password v-model:value="resetPwdForm.password" placeholder="请输入新密码" />
        </a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>

<script setup>
  import { ref, reactive, onMounted } from 'vue'
  import { message } from 'ant-design-vue'
  import { userApi } from '@/api/user'
  import { roleApi } from '@/api/role'
  import { deptApi } from '@/api/dept'

  const loading = ref(false)
  const tableData = ref([])
  const selectedRowKeys = ref([])
  const deptTree = ref([])
  const allRoles = ref([])

  const queryForm = reactive({ username: '', phone: '', status: undefined, dept_id: undefined })
  const pagination = reactive({ current: 1, pageSize: 20, total: 0, showTotal: (t) => `共 ${t} 条` })
  const rowSelection = computed(() => ({
    selectedRowKeys: selectedRowKeys.value,
    onChange: (keys) => (selectedRowKeys.value = keys),
    getCheckboxProps: (record) => ({ disabled: record.is_admin }),
  }))

  const columns = [
    { title: '头像', key: 'avatar', width: 60 },
    { title: '用户名', dataIndex: 'username', key: 'username' },
    { title: '昵称', dataIndex: 'nickname', key: 'nickname' },
    { title: '部门', dataIndex: 'dept_name', key: 'dept_name' },
    { title: '手机号', dataIndex: 'phone', key: 'phone' },
    { title: '状态', key: 'status', width: 80 },
    { title: '创建时间', dataIndex: 'created_at', key: 'created_at', width: 180 },
    { title: '操作', key: 'action', width: 240, fixed: 'right' },
  ]

  // 表单
  const modalVisible = ref(false)
  const modalTitle = ref('新增用户')
  const submitLoading = ref(false)
  const formRef = ref()
  const form = reactive({
    id: null, dept_id: undefined, username: '', nickname: '', password: '',
    phone: '', email: '', gender: 0, status: 1, remark: '',
  })
  const formRules = {
    username: [{ required: true, message: '请输入用户名' }],
    nickname: [{ required: true, message: '请输入昵称' }],
    password: [{ required: true, min: 6, message: '密码不少于6位' }],
  }

  // 角色弹窗
  const roleModalVisible = ref(false)
  const selectedRoles = ref([])
  const currentUserId = ref(null)

  // 重置密码弹窗
  const resetPwdVisible = ref(false)
  const resetPwdRef = ref()
  const resetPwdForm = reactive({ password: '' })
  const resetUserId = ref(null)

  async function loadData() {
    loading.value = true
    try {
      const params = {
        page: pagination.current,
        page_size: pagination.pageSize,
        ...queryForm,
      }
      const data = await userApi.list(params)
      tableData.value = data.list
      pagination.total = data.total
    } finally {
      loading.value = false
    }
  }

  async function loadDeptTree() {
    const data = await deptApi.list()
    deptTree.value = data || []
  }

  async function loadAllRoles() {
    const data = await roleApi.all()
    allRoles.value = data || []
  }

  function handleQuery() {
    pagination.current = 1
    loadData()
  }

  function handleReset() {
    Object.assign(queryForm, { username: '', phone: '', status: undefined, dept_id: undefined })
    handleQuery()
  }

  function handleTableChange(pag) {
    pagination.current = pag.current
    pagination.pageSize = pag.pageSize
    loadData()
  }

  function handleCreate() {
    Object.assign(form, { id: null, dept_id: undefined, username: '', nickname: '', password: '', phone: '', email: '', gender: 0, status: 1, remark: '' })
    modalTitle.value = '新增用户'
    modalVisible.value = true
  }

  async function handleEdit(record) {
    const data = await userApi.get(record.id)
    Object.assign(form, { ...data, password: '' })
    modalTitle.value = '编辑用户'
    modalVisible.value = true
  }

  async function handleSubmit() {
    await formRef.value.validate()
    submitLoading.value = true
    try {
      if (form.id) {
        await userApi.update(form.id, form)
      } else {
        await userApi.create(form)
      }
      message.success('操作成功')
      modalVisible.value = false
      loadData()
    } finally {
      submitLoading.value = false
    }
  }

  function handleModalClose() {
    formRef.value?.resetFields()
  }

  async function handleDelete(id) {
    await userApi.delete(id)
    message.success('删除成功')
    loadData()
  }

  async function handleBatchDelete() {
    // 批量删除逐个调用
    await Promise.all(selectedRowKeys.value.map((id) => userApi.delete(id)))
    message.success('删除成功')
    selectedRowKeys.value = []
    loadData()
  }

  async function handleAssignRoles(record) {
    currentUserId.value = record.id
    await loadAllRoles()
    // 获取当前用户已有角色
    const data = await userApi.get(record.id)
    selectedRoles.value = data.role_ids || []
    roleModalVisible.value = true
  }

  async function handleRoleSubmit() {
    submitLoading.value = true
    try {
      await userApi.assignRoles(currentUserId.value, { role_ids: selectedRoles.value })
      message.success('分配成功')
      roleModalVisible.value = false
    } finally {
      submitLoading.value = false
    }
  }

  function handleResetPwd(record) {
    resetUserId.value = record.id
    resetPwdForm.password = ''
    resetPwdVisible.value = true
  }

  async function handleResetPwdSubmit() {
    await resetPwdRef.value.validate()
    submitLoading.value = true
    try {
      await userApi.resetPassword(resetUserId.value, { password: resetPwdForm.password })
      message.success('重置成功')
      resetPwdVisible.value = false
    } finally {
      submitLoading.value = false
    }
  }

  onMounted(() => {
    loadData()
    loadDeptTree()
  })
</script>

<style scoped>
  .page { padding: 0; }
  .search-card { border-radius: 8px; }
  .toolbar {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 12px;
  }
</style>
