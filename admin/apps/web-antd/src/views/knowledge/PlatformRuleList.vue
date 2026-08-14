<script setup lang="ts">
import type { FormInstance } from 'ant-design-vue';

import { computed, onMounted, reactive, ref, watch } from 'vue';
import { message, Modal } from 'ant-design-vue';
import { useRoute, useRouter } from 'vue-router';

import {
  createPlatformRule,
  deletePlatformRule,
  getMerchantList,
  getPlatformRules,
  updatePlatformRule,
  type Merchant,
  type PlatformRule,
  type PlatformRulePayload,
} from '#/api/admin';

defineOptions({ name: 'PlatformRuleList' });

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
const list = ref<PlatformRule[]>([]);
const merchantOptions = ref<Merchant[]>([]);

const statusOptions = ['生效', '失效'];
const riskOptions = ['高', '中', '低'];

const form = reactive<PlatformRulePayload>({
  category: '',
  content: '',
  effectiveDate: '',
  merchantId: 0,
  platform: '抖音',
  riskLevel: '中',
  source: '',
  status: '生效',
  title: '',
});

const dialogTitle = computed(() => (editingId.value ? '编辑平台规则' : '新增平台规则'));

const activeMerchantName = computed(() => {
  if (routeMerchantName.value) return routeMerchantName.value;
  const current = merchantOptions.value.find((item) => item.id === routeMerchantId.value);
  return current?.name || '';
});

const rules = {
  title: [{ message: '请输入规则标题', required: true, trigger: 'blur' }],
  platform: [{ message: '请选择平台', required: true, trigger: 'change' }],
};

function statusColor(status: string) {
  return status === '生效' ? 'green' : 'red';
}

function riskColor(level: string) {
  if (level === '高') return 'red';
  if (level === '中') return 'orange';
  return 'default';
}

const columns = [
  { title: '商家', dataIndex: 'merchantName', key: 'merchantName', width: 150 },
  { title: '平台', dataIndex: 'platform', key: 'platform', width: 100 },
  { title: '分类', dataIndex: 'category', key: 'category', width: 130 },
  { title: '标题', dataIndex: 'title', key: 'title', ellipsis: true },
  { title: '风险等级', dataIndex: 'riskLevel', key: 'riskLevel', width: 100 },
  { title: '来源', dataIndex: 'source', key: 'source', width: 120 },
  { title: '生效日期', dataIndex: 'effectiveDate', key: 'effectiveDate', width: 120 },
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
    const result = await getPlatformRules({
      keyword: keyword.value.trim(),
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
  router.replace('/knowledge/platform-rules');
  search();
}

function resetForm() {
  Object.assign(form, {
    category: '',
    content: '',
    effectiveDate: '',
    merchantId: routeMerchantId.value || 0,
    platform: '抖音',
    riskLevel: '中',
    source: '',
    status: '生效',
    title: '',
  });
  formRef.value?.clearValidate();
}

function rowToPayload(row: PlatformRule): PlatformRulePayload {
  return {
    category: row.category,
    content: row.content,
    effectiveDate: row.effectiveDate,
    merchantId: row.merchantId,
    platform: row.platform,
    riskLevel: row.riskLevel,
    source: row.source,
    status: row.status,
    title: row.title,
  };
}

function openCreate() {
  editingId.value = null;
  resetForm();
  dialogVisible.value = true;
}

function openEdit(row: PlatformRule) {
  editingId.value = row.id;
  Object.assign(form, rowToPayload(row));
  dialogVisible.value = true;
}

async function submit() {
  await formRef.value?.validate();
  if (!form.title) {
    message.warning('请输入规则标题');
    return;
  }
  if (!form.platform) {
    message.warning('请选择平台');
    return;
  }

  submitting.value = true;
  try {
    if (editingId.value) {
      await updatePlatformRule(editingId.value, form);
      message.success('平台规则已更新');
    } else {
      await createPlatformRule(form);
      message.success('平台规则已创建');
    }
    dialogVisible.value = false;
    await loadList();
  } finally {
    submitting.value = false;
  }
}

async function removeRow(row: PlatformRule) {
  Modal.confirm({
    title: '删除平台规则',
    content: `确认删除规则「${row.title.slice(0, 20)}」？删除后列表不再保留这条记录。`,
    okText: '删除',
    okType: 'danger',
    cancelText: '取消',
    async onOk() {
      await deletePlatformRule(row.id);
      message.success('平台规则已删除');
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
            <div class="page-title">平台规则库</div>
            <div class="page-desc">
              沉淀抖音等平台的算法分发、合规红线与转化规则，用于选题与脚本判断。
            </div>
          </div>
          <div class="page-actions">
            <a-button @click="showAll">显示全部</a-button>
            <a-button type="primary" @click="openCreate">新增平台规则</a-button>
          </div>
        </div>
      </template>

      <a-form layout="inline" class="search-form">
        <a-form-item label="关键词">
          <a-input
            v-model:value="keyword"
            allow-clear
            placeholder="标题 / 分类 / 来源"
            style="width: 320px"
            @press-enter="search"
          />
        </a-input>
        </a-form-item>
        <a-form-item label="状态">
          <a-select
            class="status-select"
            style="width: 140px"
            placeholder="全部"
            allow-clear
            v-model:value="statusFilter"
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
          <template v-else-if="column.key === 'riskLevel'">
            <a-tag :color="riskColor(record.riskLevel)">{{ record.riskLevel }}</a-tag>
          </template>
          <template v-else-if="column.key === 'status'">
            <a-tag :color="statusColor(record.status)">{{ record.status }}</a-tag>
          </template>
          <template v-else-if="column.key === 'action'">
            <a-space>
              <a-button size="small" @click="openEdit(record)">编辑</a-button>
              <a-popconfirm
                title="确认删除这条平台规则？"
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

      <a-empty v-if="!loading && list.length === 0" description="当前没有平台规则">
        <a-button type="primary" @click="openCreate">新增平台规则</a-button>
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
            <a-form-item label="平台" name="platform">
              <a-input v-model:value="form.platform" placeholder="抖音 / 小红书 / 视频号" />
            </a-form-item>
          </a-col>
          <a-col :span="12">
            <a-form-item label="分类">
              <a-input v-model:value="form.category" placeholder="算法分发 / 合规红线 / 转化规则" />
            </a-form-item>
          </a-col>
          <a-col :span="12">
            <a-form-item label="风险等级">
              <a-select v-model:value="form.riskLevel" style="width: 100%">
                <a-select-option v-for="item in riskOptions" :key="item" :value="item">
                  {{ item }}
                </a-select-option>
              </a-select>
            </a-form-item>
          </a-col>
          <a-col :span="24">
            <a-form-item label="标题" name="title">
              <a-input v-model:value="form.title" placeholder="规则标题" />
            </a-form-item>
          </a-col>
          <a-col :span="24">
            <a-form-item label="内容">
              <a-textarea v-model:value="form.content" :rows="3" placeholder="规则内容" />
            </a-form-item>
          </a-col>
          <a-col :span="12">
            <a-form-item label="来源">
              <a-input v-model:value="form.source" placeholder="来源" />
            </a-form-item>
          </a-col>
          <a-col :span="12">
            <a-form-item label="生效日期">
              <a-input v-model:value="form.effectiveDate" placeholder="如：2025-01-01" />
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
