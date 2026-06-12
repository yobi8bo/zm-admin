<template>
  <div class="page">
    <a-card :bordered="false">
      <div class="toolbar">
        <a-button type="primary" @click="handleCreate"><plus-outlined />新增菜单</a-button>
        <a-button @click="loadData"><reload-outlined /></a-button>
      </div>

      <a-table
        :columns="columns"
        :data-source="tableData"
        :loading="loading"
        row-key="id"
        :pagination="false"
        :default-expand-all-rows="true"
        size="middle"
      >
        <template #bodyCell="{ column, record }">
          <template v-if="column.key === 'type'">
            <a-tag :color="typeColor[record.type]">{{ typeLabel[record.type] }}</a-tag>
          </template>
          <template v-if="column.key === 'icon'">
            <component v-if="record.icon" :is="getIcon(record.icon)" />
            <span v-else>—</span>
          </template>
          <template v-if="column.key === 'visible'">
            <a-tag :color="record.visible === 1 ? 'success' : 'default'">
              {{ record.visible === 1 ? '显示' : '隐藏' }}
            </a-tag>
          </template>
          <template v-if="column.key === 'status'">
            <a-tag :color="record.status === 1 ? 'success' : 'error'">
              {{ record.status === 1 ? '正常' : '停用' }}
            </a-tag>
          </template>
          <template v-if="column.key === 'action'">
            <a-space size="small">
              <a-button type="link" size="small" @click="handleCreate(record.id)">新增子项</a-button>
              <a-button type="link" size="small" @click="handleEdit(record)">编辑</a-button>
              <a-popconfirm title="确定删除吗？" @confirm="handleDelete(record.id)">
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
      :confirm-loading="submitLoading"
      @cancel="formRef?.resetFields()"
    >
      <a-form :model="form" :rules="formRules" ref="formRef" :label-col="{ span: 6 }" :wrapper-col="{ span: 16 }">
        <a-form-item label="上级菜单">
          <a-tree-select
            v-model:value="form.parent_id"
            :tree-data="tableData"
            :field-names="{ label: 'name', value: 'id', children: 'children' }"
            placeholder="顶级菜单"
            allow-clear
          />
        </a-form-item>
        <a-form-item label="菜单类型" name="type">
          <a-radio-group v-model:value="form.type">
            <a-radio :value="1">目录</a-radio>
            <a-radio :value="2">菜单</a-radio>
            <a-radio :value="3">按钮</a-radio>
          </a-radio-group>
        </a-form-item>
        <a-form-item label="菜单名称" name="name">
          <a-input v-model:value="form.name" placeholder="请输入菜单名称" />
        </a-form-item>
        <a-form-item v-if="form.type !== 3" label="路由地址">
          <a-input v-model:value="form.path" placeholder="如：/system/user 或 user" />
        </a-form-item>
        <a-form-item v-if="form.type === 2" label="组件路径">
          <a-input v-model:value="form.component" placeholder="如：system/user/index" />
        </a-form-item>
        <a-form-item v-if="form.type !== 1" label="权限标识">
          <a-input v-model:value="form.permission" placeholder="如：system:user:list" />
        </a-form-item>
        <a-form-item v-if="form.type !== 3" label="菜单图标">
          <a-input v-model:value="form.icon" placeholder="如：UserOutlined" />
        </a-form-item>
        <a-form-item label="显示排序">
          <a-input-number v-model:value="form.sort" :min="0" style="width: 100%" />
        </a-form-item>
        <a-form-item v-if="form.type !== 3" label="显示状态">
          <a-radio-group v-model:value="form.visible">
            <a-radio :value="1">显示</a-radio>
            <a-radio :value="0">隐藏</a-radio>
          </a-radio-group>
        </a-form-item>
        <a-form-item label="菜单状态">
          <a-radio-group v-model:value="form.status">
            <a-radio :value="1">正常</a-radio>
            <a-radio :value="0">停用</a-radio>
          </a-radio-group>
        </a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>

<script setup>
  import { ref, reactive, onMounted } from 'vue'
  import { message } from 'ant-design-vue'
  import * as Icons from '@ant-design/icons-vue'
  import { menuApi } from '@/api/menu'

  const loading = ref(false)
  const tableData = ref([])

  const typeLabel = { 1: '目录', 2: '菜单', 3: '按钮' }
  const typeColor = { 1: 'blue', 2: 'green', 3: 'orange' }

  const columns = [
    { title: '菜单名称', dataIndex: 'name', key: 'name', width: 180 },
    { title: '类型', key: 'type', width: 80 },
    { title: '图标', key: 'icon', width: 60 },
    { title: '路由地址', dataIndex: 'path', key: 'path' },
    { title: '权限标识', dataIndex: 'permission', key: 'permission' },
    { title: '排序', dataIndex: 'sort', key: 'sort', width: 70 },
    { title: '显示', key: 'visible', width: 80 },
    { title: '状态', key: 'status', width: 80 },
    { title: '操作', key: 'action', width: 200, fixed: 'right' },
  ]

  const modalVisible = ref(false)
  const modalTitle = ref('新增菜单')
  const submitLoading = ref(false)
  const formRef = ref()
  const form = reactive({
    id: null, parent_id: undefined, name: '', type: 1,
    path: '', component: '', permission: '', icon: '',
    sort: 0, visible: 1, status: 1,
  })
  const formRules = { name: [{ required: true, message: '请输入菜单名称' }], type: [{ required: true }] }

  function getIcon(name) { return Icons[name] || Icons['AppstoreOutlined'] }

  async function loadData() {
    loading.value = true
    try {
      tableData.value = await menuApi.list() || []
    } finally {
      loading.value = false
    }
  }

  function handleCreate(parentId = undefined) {
    Object.assign(form, { id: null, parent_id: parentId, name: '', type: 1, path: '', component: '', permission: '', icon: '', sort: 0, visible: 1, status: 1 })
    modalTitle.value = '新增菜单'
    modalVisible.value = true
  }

  async function handleEdit(record) {
    Object.assign(form, record)
    modalTitle.value = '编辑菜单'
    modalVisible.value = true
  }

  async function handleSubmit() {
    await formRef.value.validate()
    submitLoading.value = true
    try {
      form.id ? await menuApi.update(form.id, form) : await menuApi.create(form)
      message.success('操作成功')
      modalVisible.value = false
      loadData()
    } finally {
      submitLoading.value = false
    }
  }

  async function handleDelete(id) {
    await menuApi.delete(id)
    message.success('删除成功')
    loadData()
  }

  onMounted(loadData)
</script>

<style scoped>
  .page { padding: 0; }
  .toolbar { display: flex; justify-content: space-between; align-items: center; margin-bottom: 12px; }
</style>
