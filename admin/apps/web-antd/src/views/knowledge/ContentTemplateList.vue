<script setup lang="ts">
import type { FormInstance } from 'ant-design-vue';

import { computed, onMounted, reactive, ref, watch } from 'vue';
import { message, Modal } from 'ant-design-vue';
import { useRoute, useRouter } from 'vue-router';

import {
  createContentTemplate,
  deleteContentTemplate,
  getContentTemplates,
  getMerchantList,
  updateContentTemplate,
  type ContentTemplate,
  type ContentTemplatePayload,
  type Merchant,
} from '#/api/admin';

defineOptions({ name: 'ContentTemplateList' });

const route = useRoute();
const router = useRouter();

const loading = ref(false);
const submitting = ref(false);
const dialogVisible = ref(false);
const editingId = ref<number | null>(null);
const formRef = ref<FormInstance>();

const keyword = ref('');
const statusFilter = ref('');
const routeMerchantId = ref<number | undefined>();
const routeMerchantName = ref('');
const pagination = reactive({ page: 1, size: 10, total: 0 });
const list = ref<ContentTemplate[]>([]);
const merchantOptions = ref<Merchant[]>([]);

const statusOptions = ['启用', '停用'];

const form = reactive<ContentTemplatePayload>({
  category: '',
  content: '',
  merchantId: 0,
  name: '',
  structureJson: '{}',
  type: '',
  usageNote: '',
  status: '启用',
});

const dialogTitle = computed(() => (editingId.value ? '编辑内容模板' : '新增内容模板'));

const activeMerchantName = computed(() => {
  if (routeMerchantName.value) return routeMerchantName.value;
  const current = merchantOptions.value.find((item) => item.id === routeMerchantId.value);
  return current?.name || '';
});

const rules = {
  name: [{ message: '请输入模板名称', required: true, trigger: 'blur' }],
};

function statusColor(status: string) {
  return status === '启用' ? 'green' : 'default';
}

const columns = [
  { title: '商家', dataIndex: 'merchantName', key: 'merchantName', width: 150 },
  { title: '名称', dataIndex: 'name', key: 'name', width: 180 },
  { title: '类型', dataIndex: 'type', key: 'type', width: 140 },
  { title: '分类', dataIndex: 'category', key: 'category', width: 140 },
  { title: '内容', dataIndex: 'content', key: 'content', ellipsis: true },
  { title: '用法说明', dataIndex: 'usageNote', key: 'usageNote', ellipsis: true },
  { title: '状态', dataIndex: 'status', key: 'status', width: 90, fixed: 'right' as const },
  { title: '操作', key: 'action', width: 160, fixed: 'right' as const },
];

function applyRouteQuery() {
  const rawMerchantId = Number(route.query.merchantId || 0);
  routeMerchantId.value = rawMerchantId > 0 ? rawMerchantId : undefined;
  routeMerchantName.value = String(route.query.merchantName || '');
  keyword.value = routeMerchantName.value;
}

async function loadMerchants() {
  const result = await getMerchantList({ page: 1, size: 100 });
  merchantOptions.value = result.list;
}

async function loadList() {
  loading.value = true;
  try {
    const result = await getContentTemplates({
      keyword: keyword.value.trim(),
      merchantId: routeMerchantId.value,
      page: pagination.page,
      size: pagination.size,
      status: statusFilter.value || undefined,
    });
    list.value = result.list;
    pagination.total = result.total;
    pagination.page = result.page;
    pagination.size = result.size;
  } finally {
    loading.value = false;
  }
}

function search() {
  pagination.page = 1;
  loadList();
}

function showAll() {
  keyword.value = '';
  statusFilter.value = '';
  routeMerchantId.value = undefined;
  routeMerchantName.value = '';
  router.replace('/knowledge/content-templates');
  search();
}

function resetForm() {
  Object.assign(form, {
    category: '',
    content: '',
    merchantId: routeMerchantId.value || 0,
    name: '',
    structureJson: '{}',
    type: '',
    usageNote: '',
    status: '启用',
  });
  formRef.value?.clearValidate();
}

function rowToPayload(row: ContentTemplate): ContentTemplatePayload {
  return {
    category: row.category,
    content: row.content,
    merchantId: row.merchantId,
    name: row.name,
    structureJson: row.structureJson,
    type: row.type,
    usageNote: row.usageNote,
    status: row.status,
  };
}

function openCreate() {
  editingId.value = null;
  resetForm();
  dialogVisible.value = true;
}

function openEdit(row: ContentTemplate) {
  editingId.value = row.id;
  Object.assign(form, rowToPayload(row));
  dialogVisible.value = true;
}

async function submit() {
  await formRef.value?.validate();
  if (!form.name) {
    message.warning('请输入模板名称');
    return;
  }

  submitting.value = true;
  try {
    if (editingId.value) {
      await updateContentTemplate(editingId.value, form);
      message.success('内容模板已更新');
    } else {
      await createContentTemplate(form);
      message.success('内容模板已创建');
    }
    dialogVisible.value = false;
    await loadList();
  } finally {
    submitting.value = false;
  }
}

async function removeRow(row: ContentTemplate) {
  Modal.confirm({
    title: '删除内容模板',
    content: `确认删除模板「${row.name.slice(0, 20)}」？删除后列表不再保留这条记录。`,
    okText: '删除',
    okType: 'danger',
    cancelText: '取消',
    async onOk() {
      await deleteContentTemplate(row.id);
      message.success('内容模板已删除');
      await loadList();
    },
  });
}

function handlePageChange(page: number) {
  pagination.page = page;
  loadList();
}

function handleSizeChange(size: number) {
  pagination.size = size;
  pagination.page = 1;
  loadList();
}

watch(
  () => route.query,
  () => {
    applyRouteQuery();
    pagination.page = 1;
    loadList();
  },
);

onMounted(async () => {
  applyRouteQuery();
  await loadMerchants();
  await loadList();
});
</script>

<template>
  <div class="p-4">
    <a-card>
      <template #title>
        <div class="page-head">
          <div>
            <div class="page-title">内容模板库</div>
            <div class="page-desc">
              沉淀选题卡 / 口播脚本 / 分镜脚本 / 复盘报告等结构化模板。
            </div>
          </div>
          <div class="page-actions">
            <a-button @click="showAll">显示全部</a-button>
            <a-button type="primary" @click="openCreate">新增内容模板</a-button>
          </div>
        </div>
      </template>

      <a-alert
        v-if="activeMerchantName"
        class="mb-4"
        type="info"
        :closable="false"
        show-icon
      >
        当前处理商家：{{ activeMerchantName }}。已自动带入筛选条件，商家ID：{{ routeMerchantId || '-' }}。
      </a-alert>

      <a-form layout="inline" class="search-form">
        <a-form-item label="关键词">
          <a-input
            v-model:value="keyword"
            allow-clear
            placeholder="名称 / 类型 / 分类"
            style="width: 320px"
            @press-enter="search"
          />
        </a-input>
        </a-form-item>
        <a-form-item label="状态">
          <a-select
            v-model:value="statusFilter"
            allow-clear
            placeholder="全部"
            style="width: 140px"
          >
            <a-select-option v-for="item in statusOptions" :key="item" :value="item">
              {{ item }}
            </a-select-option>
          </a-select>
        </a-form-item>
        <a-form-item>
          <a-button type="primary" @click="search">搜索</a-button>
        </a-form-item>
      </a-form>

      <a-table
        :data-source="list"
        :columns="columns"
        :loading="loading"
        :pagination="false"
        row-key="id"
        bordered
        size="middle"
      >
        <template #bodyCell="{ column, record }">
          <template v-if="column.key === 'status'">
            <a-tag :color="statusColor(record.status)">{{ record.status }}</a-tag>
          </template>
          <template v-else-if="column.key === 'action'">
            <a-space>
              <a-button size="small" @click="openEdit(record)">编辑</a-button>
              <a-popconfirm
                title="确认删除这条内容模板？"
                ok-text="删除"
                cancel-text="取消"
                @confirm="removeRow(record)"
              >
                <a-button size="small" danger>删除</a-button>
              </a-popconfirm>
            </a-space>
          </template>
        </template>
      </a-table>

      <a-empty v-if="!loading && list.length === 0" description="当前没有内容模板">
        <a-button type="primary" @click="openCreate">新增内容模板</a-button>
      </a-empty>

      <div class="pagination-wrap">
        <a-pagination
          v-model:current="pagination.page"
          v-model:page-size="pagination.size"
          :total="pagination.total"
          :show-size-changer="true"
          :page-size-options="['10', '20', '50']"
          show-quick-jumper
          @change="handlePageChange"
          @showSizeChange="handleSizeChange"
        />
      </div>
    </a-card>

    <a-modal
      v-model:open="dialogVisible"
      :title="dialogTitle"
      width="820px"
      :confirm-loading="submitting"
      @ok="submit"
      @cancel="dialogVisible = false"
    >
      <a-form ref="formRef" :model="form" :rules="rules" layout="vertical">
        <a-row :gutter="16">
          <a-col :span="12">
            <a-form-item label="商家">
              <a-select
                v-model:value="form.merchantId"
                allow-clear
                show-search
                placeholder="选择商家（可留空）"
                :options="merchantOptions.map((m) => ({ label: m.name, value: m.id }))"
                style="width: 100%"
              />
            </a-form-item>
          </a-col>
          <a-col :span="12">
            <a-form-item label="名称" name="name">
              <a-input v-model:value="form.name" placeholder="模板名称" />
            </a-form-item>
          </a-col>
          <a-col :span="12">
            <a-form-item label="类型">
              <a-input v-model:value="form.type" placeholder="选题卡 / 口播脚本 / 分镜脚本 / 复盘报告 / 对标报告" />
            </a-form-item>
          </a-col>
          <a-col :span="12">
            <a-form-item label="分类">
              <a-input v-model:value="form.category" placeholder="栏目 / 场景" />
            </a-form-item>
          </a-col>
          <a-col :span="24">
            <a-form-item label="内容">
              <a-textarea v-model:value="form.content" :rows="3" placeholder="模板内容" />
            </a-form-item>
          </a-col>
          <a-col :span="24">
            <a-form-item label="结构（JSON）">
              <a-textarea v-model:value="form.structureJson" :rows="2" placeholder='{}' />
            </a-form-item>
          </a-col>
          <a-col :span="24">
            <a-form-item label="用法说明">
              <a-textarea v-model:value="form.usageNote" :rows="2" placeholder="用法说明" />
            </a-form-item>
          </a-col>
          <a-col :span="12">
            <a-form-item label="状态">
              <a-select v-model:value="form.status" style="width: 100%">
                <a-select-option v-for="item in statusOptions" :key="item" :value="item">
                  {{ item }}
                </a-select-option>
              </a-select>
            </a-form-item>
          </a-col>
        </a-row>
      </a-form>
    </a-modal>
  </div>
</template>

<style scoped>
.page-head {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 16px;
}

.page-title {
  color: #0f172a;
  font-size: 18px;
  font-weight: 700;
}

.page-desc {
  color: #64748b;
  font-size: 13px;
}

.page-actions {
  display: flex;
  flex-shrink: 0;
  gap: 10px;
}

.search-form {
  margin-bottom: 12px;
}

.pagination-wrap {
  display: flex;
  justify-content: flex-end;
  margin-top: 16px;
}
</style>
