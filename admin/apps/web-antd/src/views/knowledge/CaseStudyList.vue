<script setup lang="ts">
import type { FormInstance } from 'ant-design-vue';

import { computed, onMounted, reactive, ref, watch } from 'vue';
import { message, Modal } from 'ant-design-vue';
import { useRoute, useRouter } from 'vue-router';

import {
  createCaseStudy,
  deleteCaseStudy,
  getCaseStudies,
  getMerchantList,
  updateCaseStudy,
  type CaseStudy,
  type CaseStudyPayload,
  type Merchant,
} from '#/api/admin';

defineOptions({ name: 'CaseStudyList' });

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
const list = ref<CaseStudy[]>([]);
const merchantOptions = ref<Merchant[]>([]);

const statusOptions = ['有效', '失效'];

const form = reactive<CaseStudyPayload>({
  accountName: '',
  conversionAction: '',
  form: '',
  hookType: '',
  industry: '',
  lesson: '',
  merchantId: 0,
  platform: '',
  reusablePoint: '',
  structure: '',
  title: '',
  dataJson: '{}',
  status: '有效',
});

const dialogTitle = computed(() => (editingId.value ? '编辑案例' : '新增案例'));

const activeMerchantName = computed(() => {
  if (routeMerchantName.value) return routeMerchantName.value;
  const current = merchantOptions.value.find((item) => item.id === routeMerchantId.value);
  return current?.name || '';
});

const rules = {
  title: [{ message: '请输入案例标题', required: true, trigger: 'blur' }],
};

function statusColor(status: string) {
  return status === '有效' ? 'green' : 'red';
}

const columns = [
  { title: '商家', dataIndex: 'merchantName', key: 'merchantName', width: 150 },
  { title: '标题', dataIndex: 'title', key: 'title', ellipsis: true },
  { title: '平台', dataIndex: 'platform', key: 'platform', width: 100 },
  { title: '来源账号', dataIndex: 'accountName', key: 'accountName', width: 130 },
  { title: '行业', dataIndex: 'industry', key: 'industry', width: 120 },
  { title: '形式', dataIndex: 'form', key: 'form', width: 100 },
  { title: '钩子类型', dataIndex: 'hookType', key: 'hookType', width: 110 },
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
    const result = await getCaseStudies({
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
  router.replace('/knowledge/case-studies');
  search();
}

function resetForm() {
  Object.assign(form, {
    accountName: '',
    conversionAction: '',
    form: '',
    hookType: '',
    industry: '',
    lesson: '',
    merchantId: routeMerchantId.value || 0,
    platform: '',
    reusablePoint: '',
    structure: '',
    title: '',
    dataJson: '{}',
    status: '有效',
  });
  formRef.value?.clearValidate();
}

function rowToPayload(row: CaseStudy): CaseStudyPayload {
  return {
    accountName: row.accountName,
    conversionAction: row.conversionAction,
    form: row.form,
    hookType: row.hookType,
    industry: row.industry,
    lesson: row.lesson,
    merchantId: row.merchantId,
    platform: row.platform,
    reusablePoint: row.reusablePoint,
    structure: row.structure,
    title: row.title,
    dataJson: row.dataJson,
    status: row.status,
  };
}

function openCreate() {
  editingId.value = null;
  resetForm();
  dialogVisible.value = true;
}

function openEdit(row: CaseStudy) {
  editingId.value = row.id;
  Object.assign(form, rowToPayload(row));
  dialogVisible.value = true;
}

async function submit() {
  await formRef.value?.validate();
  if (!form.title) {
    message.warning('请输入案例标题');
    return;
  }

  submitting.value = true;
  try {
    if (editingId.value) {
      await updateCaseStudy(editingId.value, form);
      message.success('案例已更新');
    } else {
      await createCaseStudy(form);
      message.success('案例已创建');
    }
    dialogVisible.value = false;
    await loadList();
  } finally {
    submitting.value = false;
  }
}

async function removeRow(row: CaseStudy) {
  Modal.confirm({
    title: '删除案例',
    content: `确认删除案例「${row.title.slice(0, 20)}」？删除后列表不再保留这条记录。`,
    okText: '删除',
    okType: 'danger',
    cancelText: '取消',
    async onOk() {
      await deleteCaseStudy(row.id);
      message.success('案例已删除');
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
            <div class="page-title">案例库</div>
            <div class="page-desc">
              拆解对标 / 自身爆款，提取可复用结构、钩子与转化动作。
            </div>
          </div>
          <div class="page-actions">
            <a-button @click="showAll">显示全部</a-button>
            <a-button type="primary" @click="openCreate">新增案例</a-button>
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
            placeholder="标题 / 账号 / 行业"
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
                title="确认删除这条案例？"
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

      <a-empty v-if="!loading && list.length === 0" description="当前没有案例">
        <a-button type="primary" @click="openCreate">新增案例</a-button>
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
            <a-form-item label="标题" name="title">
              <a-input v-model:value="form.title" placeholder="爆款视频标题" />
            </a-form-item>
          </a-col>
          <a-col :span="12">
            <a-form-item label="平台">
              <a-input v-model:value="form.platform" placeholder="抖音 / 小红书 / 视频号" />
            </a-form-item>
          </a-col>
          <a-col :span="12">
            <a-form-item label="来源账号">
              <a-input v-model:value="form.accountName" placeholder="对标账号名" />
            </a-form-item>
          </a-col>
          <a-col :span="12">
            <a-form-item label="行业">
              <a-input v-model:value="form.industry" placeholder="行业 / 赛道" />
            </a-form-item>
          </a-col>
          <a-col :span="12">
            <a-form-item label="形式">
              <a-input v-model:value="form.form" placeholder="口播 / 剧情 / 测评 / 干货" />
            </a-form-item>
          </a-col>
          <a-col :span="12">
            <a-form-item label="钩子类型">
              <a-input v-model:value="form.hookType" placeholder="钩子类型" />
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
          <a-col :span="24">
            <a-form-item label="结构（痛点-价值-行动）">
              <a-textarea v-model:value="form.structure" :rows="3" placeholder="拆解结构" />
            </a-form-item>
          </a-col>
          <a-col :span="24">
            <a-form-item label="转化动作">
              <a-textarea v-model:value="form.conversionAction" :rows="2" placeholder="转化动作" />
            </a-form-item>
          </a-col>
          <a-col :span="24">
            <a-form-item label="可复用点">
              <a-textarea v-model:value="form.reusablePoint" :rows="2" placeholder="可复用点" />
            </a-form-item>
          </a-col>
          <a-col :span="24">
            <a-form-item label="数据快照（JSON）">
              <a-textarea v-model:value="form.dataJson" :rows="2" placeholder='{}' />
            </a-form-item>
          </a-col>
          <a-col :span="24">
            <a-form-item label="经验沉淀">
              <a-textarea v-model:value="form.lesson" :rows="2" placeholder="经验沉淀" />
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
