<script setup lang="ts">
import type { FormInstance } from 'ant-design-vue';

import { computed, onMounted, reactive, ref, watch } from 'vue';
import { message, Modal } from 'ant-design-vue';
import { useRoute, useRouter } from 'vue-router';

import {
  createPainPoint,
  deletePainPoint,
  getMerchantList,
  getPainPoints,
  updatePainPoint,
  type Merchant,
  type PainPoint,
  type PainPointPayload,
} from '#/api/admin';

defineOptions({ name: 'PainPointList' });

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
const list = ref<PainPoint[]>([]);
const merchantOptions = ref<Merchant[]>([]);

const statusOptions = ['待拆', '已采纳', '已转化'];

const form = reactive<PainPointPayload>({
  category: '',
  content: '',
  demandLevel: '中',
  emotion: '',
  merchantId: 0,
  product: '',
  source: '',
  status: '待拆',
  suggestedTopic: '',
  userQuote: '',
});

const dialogTitle = computed(() => (editingId.value ? '编辑痛点' : '新增痛点'));

const activeMerchantName = computed(() => {
  if (routeMerchantName.value) return routeMerchantName.value;
  const current = merchantOptions.value.find((item) => item.id === routeMerchantId.value);
  return current?.name || '';
});

const rules = {
  content: [{ message: '请输入痛点描述', required: true, trigger: 'blur' }],
};

function statusColor(status: string) {
  if (status === '已转化') return 'gold';
  if (status === '已采纳') return 'green';
  return 'blue';
}

const columns = [
  { title: '商家', dataIndex: 'merchantName', key: 'merchantName', width: 150 },
  { title: '来源', dataIndex: 'source', key: 'source', width: 110 },
  { title: '分类', dataIndex: 'category', key: 'category', width: 130 },
  { title: '痛点描述', dataIndex: 'content', key: 'content', ellipsis: true },
  { title: '用户原话', dataIndex: 'userQuote', key: 'userQuote', ellipsis: true },
  { title: '情绪', dataIndex: 'emotion', key: 'emotion', width: 100 },
  { title: '对应产品', dataIndex: 'product', key: 'product', width: 120 },
  { title: '需求强度', dataIndex: 'demandLevel', key: 'demandLevel', width: 90 },
  { title: '建议选题', dataIndex: 'suggestedTopic', key: 'suggestedTopic', ellipsis: true },
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
    const result = await getPainPoints({
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
  router.replace('/knowledge/pain-points');
  search();
}

function resetForm() {
  Object.assign(form, {
    category: '',
    content: '',
    demandLevel: '中',
    emotion: '',
    merchantId: routeMerchantId.value || 0,
    product: '',
    source: '',
    status: '待拆',
    suggestedTopic: '',
    userQuote: '',
  });
  formRef.value?.clearValidate();
}

function rowToPayload(row: PainPoint): PainPointPayload {
  return {
    category: row.category,
    content: row.content,
    demandLevel: row.demandLevel,
    emotion: row.emotion,
    merchantId: row.merchantId,
    product: row.product,
    source: row.source,
    status: row.status,
    suggestedTopic: row.suggestedTopic,
    userQuote: row.userQuote,
  };
}

function openCreate() {
  editingId.value = null;
  resetForm();
  dialogVisible.value = true;
}

function openEdit(row: PainPoint) {
  editingId.value = row.id;
  Object.assign(form, rowToPayload(row));
  dialogVisible.value = true;
}

async function submit() {
  await formRef.value?.validate();
  if (!form.content) {
    message.warning('请输入痛点描述');
    return;
  }

  submitting.value = true;
  try {
    if (editingId.value) {
      await updatePainPoint(editingId.value, form);
      message.success('痛点已更新');
    } else {
      await createPainPoint(form);
      message.success('痛点已创建');
    }
    dialogVisible.value = false;
    await loadList();
  } finally {
    submitting.value = false;
  }
}

async function removeRow(row: PainPoint) {
  Modal.confirm({
    title: '删除痛点',
    content: `确认删除痛点「${row.content.slice(0, 20)}」？删除后列表不再保留这条记录。`,
    okText: '删除',
    okType: 'danger',
    cancelText: '取消',
    async onOk() {
      await deletePainPoint(row.id);
      message.success('痛点已删除');
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
            <div class="page-title">痛点库</div>
            <div class="page-desc">
              沉淀目标人群在各平台、各场景下的真实痛点，作为选题与文案的源头。
            </div>
          </div>
          <div class="page-actions">
            <a-button @click="showAll">显示全部</a-button>
            <a-button type="primary" @click="openCreate">新增痛点</a-button>
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
            placeholder="商家 / 痛点 / 来源"
            style="width: 320px"
            @press-enter="search"
          />
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
                title="确认删除这条痛点？"
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

      <a-empty v-if="!loading && list.length === 0" description="当前没有痛点">
        <a-button type="primary" @click="openCreate">新增痛点</a-button>
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
            <a-form-item label="来源">
              <a-input v-model:value="form.source" placeholder="私信 / 评论区 / 群聊 / 热点 / 对标" />
            </a-form-item>
          </a-col>
          <a-col :span="12">
            <a-form-item label="分类">
              <a-input v-model:value="form.category" placeholder="内容选题 / 拍摄剪辑 / 投放转化 / 账号定位" />
            </a-form-item>
          </a-col>
          <a-col :span="12">
            <a-form-item label="情绪">
              <a-input v-model:value="form.emotion" placeholder="焦虑 / 困惑 / 期待" />
            </a-form-item>
          </a-col>
          <a-col :span="24">
            <a-form-item label="痛点描述" name="content">
              <a-textarea
                v-model:value="form.content"
                :rows="3"
                :maxlength="1000"
                placeholder="一句说清用户卡在哪"
                show-count
              />
            </a-form-item>
          </a-col>
          <a-col :span="24">
            <a-form-item label="用户原话">
              <a-textarea
                v-model:value="form.userQuote"
                :rows="2"
                placeholder="用户原话记录"
              />
            </a-form-item>
          </a-col>
          <a-col :span="12">
            <a-form-item label="对应产品">
              <a-input v-model:value="form.product" placeholder="如：双人火锅套餐" />
            </a-form-item>
          </a-col>
          <a-col :span="12">
            <a-form-item label="需求强度">
              <a-select v-model:value="form.demandLevel" style="width: 100%">
                <a-select-option value="高">高</a-select-option>
                <a-select-option value="中">中</a-select-option>
                <a-select-option value="低">低</a-select-option>
              </a-select>
            </a-form-item>
          </a-col>
          <a-col :span="24">
            <a-form-item label="建议选题方向">
              <a-textarea
                v-model:value="form.suggestedTopic"
                :rows="2"
                placeholder="痛点背后的选题方向或角度"
              />
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
