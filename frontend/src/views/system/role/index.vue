<template>
  <div class="page">
    <a-card :bordered="false" class="search-card">
      <a-form :model="queryForm" layout="inline">
        <a-form-item label="角色名称">
          <a-input v-model:value="queryForm.name" placeholder="请输入角色名称" allow-clear style="width: 180px" />
        </a-form-item>
        <a-form-item label="角色标识">
          <a-input v-model:value="queryForm.code" placeholder="请输入角色标识" allow-clear style="width: 180px" />
        </a-form-item>
        <a-form-item label="状态">
          <a-select v-model:value="queryForm.status" placeholder="全部" allow-clear style="width: 100px">
            <a-select-option :value="1">正常</a-select-option>
            <a-select-option :value="0">停用</a-select-option>
          </a-select>
        </a-form-item>
        <a-form-item>
          <a-space>
            <a-button type="primary" @click="handleQuery"><search-outlined />查询</a-button>
            <a-button @click="handleReset">重置</a-button>
          </a-space>
        </a-form-item>
      </a-form>
    </a-card>

    <a-card :bordered="false" style="margin-top: 12px">
      <div class="toolbar">
        <a-button type="primary" @click="handleCreate"><plus-outlined />新增</a-button>
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
              {{ record.status === 1 ? '正常' : '停用' }}
            </a-tag>
          </template>
          <template v-if="column.key === 'action'">
            <a-space size="small">
              <a-button type="link" size="small" @click="handleEdit(record)">编辑</a-button>
              <a-button type="link" size="small" @click="handleAssignMenus(record)">菜单权限</a-button>
              <a-popconfirm title="确定删除该角色吗？" @confirm="handleDelete(record.id)">
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
      @ok="handleSubmit"
      :confirm-loading="submitLoading"
      @cancel="formRef?.resetFields()"
    >
      <a-form :model="form" :rules="formRules" ref="formRef" :label-col="{ span: 6 }" :wrapper-col="{ span: 16 }">
        <a-form-item label="角色名称" name="name">
          <a-input v-model:value="form.name" placeholder="请输入角色名称" />
        </a-form-item>
        <a-form-item label="角色标识" name="code">
          <a-input v-model:value="form.code" placeholder="如：admin" />
        </a-form-item>
        <a-form-item label="显示排序" name="sort">
          <a-input-number v-model:value="form.sort" :min="0" style="width: 100%" />
        </a-form-item>
        <a-form-item label="状态">
          <a-radio-group v-model:value="form.status">
            <a-radio :value="1">正常</a-radio>
            <a-radio :value="0">停用</a-radio>
          </a-radio-group>
        </a-form-item>
        <a-form-item label="备注">
          <a-textarea v-model:value="form.remark" :rows="3" />
        </a-form-item>
      </a-form>
    </a-modal>

    <!-- 菜单权限弹窗 -->
    <a-modal
      v-model:open="menuModalVisible"
      title="菜单权限"
      width="480px"
      @ok="handleMenuSubmit"
      :confirm-loading="submitLoading"
    >
      <a-tree
        v-model:checkedKeys="checkedMenuKeys"
        :tree-data="menuTree"
        :field-names="{ title: 'name', key: 'id', children: 'children' }"
        checkable
        default-expand-all
      />
    </a-modal>
  </div>
</template>

<script setup>
  import { ref, reactive, onMounted } from 'vue'
  import { message } from 'ant-design-vue'
  import { roleApi } from '@/api/role'
  import { menuApi } from '@/api/menu'

  const loading = ref(false)
  const tableData = ref([])
  const queryForm = reactive({ name: '', code: '', status: undefined })
  const pagination = reactive({ current: 1, pageSize: 20, total: 0, showTotal: (t) => `共 ${t} 条` })

  const columns = [
    { title: '角色名称', dataIndex: 'name', key: 'name' },
    { title: '角色标识', dataIndex: 'code', key: 'code' },
    { title: '排序', dataIndex: 'sort', key: 'sort', width: 80 },
    { title: '状态', key: 'status', width: 80 },
    { title: '创建时间', dataIndex: 'created_at', key: 'created_at', width: 180 },
    { title: '操作', key: 'action', width: 200, fixed: 'right' },
  ]

  const modalVisible = ref(false)
  const modalTitle = ref('新增角色')
  const submitLoading = ref(false)
  const formRef = ref()
  const form = reactive({ id: null, name: '', code: '', sort: 0, status: 1, remark: '' })
  const formRules = {
    name: [{ required: true, message: '请输入角色名称' }],
    code: [{ required: true, message: '请输入角色标识' }],
  }

  // 菜单权限弹窗
  const menuModalVisible = ref(false)
  const menuTree = ref([])
  const checkedMenuKeys = ref([])
  const currentRoleId = ref(null)

  async function loadData() {
    loading.value = true
    try {
      const data = await roleApi.list({ page: pagination.current, page_size: pagination.pageSize, ...queryForm })
      tableData.value = data.list
      pagination.total = data.total
    } finally {
      loading.value = false
    }
  }

  function handleQuery() { pagination.current = 1; loadData() }
  function handleReset() { Object.assign(queryForm, { name: '', code: '', status: undefined }); handleQuery() }
  function handleTableChange(pag) { pagination.current = pag.current; loadData() }

  function handleCreate() {
    Object.assign(form, { id: null, name: '', code: '', sort: 0, status: 1, remark: '' })
    modalTitle.value = '新增角色'
    modalVisible.value = true
  }

  async function handleEdit(record) {
    Object.assign(form, record)
    modalTitle.value = '编辑角色'
    modalVisible.value = true
  }

  async function handleSubmit() {
    await formRef.value.validate()
    submitLoading.value = true
    try {
      form.id ? await roleApi.update(form.id, form) : await roleApi.create(form)
      message.success('操作成功')
      modalVisible.value = false
      loadData()
    } finally {
      submitLoading.value = false
    }
  }

  async function handleDelete(id) {
    await roleApi.delete(id)
    message.success('删除成功')
    loadData()
  }

  async function handleAssignMenus(record) {
    currentRoleId.value = record.id
    const [menuData, checkedIds] = await Promise.all([menuApi.list(), roleApi.getMenuIDs(record.id)])
    menuTree.value = menuData || []
    checkedMenuKeys.value = checkedIds || []
    menuModalVisible.value = true
  }

  async function handleMenuSubmit() {
    submitLoading.value = true
    try {
      await roleApi.assignMenus(currentRoleId.value, { menu_ids: checkedMenuKeys.value })
      message.success('保存成功')
      menuModalVisible.value = false
    } finally {
      submitLoading.value = false
    }
  }

  onMounted(loadData)
</script>

<style scoped>
  .page { padding: 0; }
  .search-card { border-radius: 8px; }
  .toolbar { display: flex; justify-content: space-between; align-items: center; margin-bottom: 12px; }
</style>
