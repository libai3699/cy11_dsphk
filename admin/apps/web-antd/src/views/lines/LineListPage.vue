<script setup lang="ts">
import type { FormInstance } from 'element-plus';

import { computed, onMounted, reactive, ref, watch } from 'vue';
import { ElMessage, ElMessageBox } from 'element-plus';
import { useRoute, useRouter } from 'vue-router';

import {
  analyzeBenchmarkAccount,
  analyzeBenchmarkAccounts,
  createBenchmarkAccount,
  deleteBenchmarkAccount,
  getBenchmarkAccountList,
  getMerchantList,
  updateBenchmarkAccount,
  type BenchmarkAccount,
  type BenchmarkAccountPayload,
  type BenchmarkAnalysisTask,
  type Merchant,
} from '#/api/admin';

defineOptions({ name: 'LineListPage' });

const route = useRoute();
const router = useRouter();

const loading = ref(false);
const submitting = ref(false);
const analyzingId = ref<number | null>(null);
const dialogVisible = ref(false);
const resultVisible = ref(false);
const editingId = ref<number | null>(null);
const formRef = ref<FormInstance>();

const keyword = ref('');
const statusFilter = ref('');
const routeMerchantId = ref<number | undefined>();
const routeMerchantName = ref('');
const pagination = reactive({ page: 1, size: 10, total: 0 });
const list = ref<BenchmarkAccount[]>([]);
const merchantOptions = ref<Merchant[]>([]);
const activeAnalysis = ref<BenchmarkAnalysisTask>();
const selectedRows = ref<BenchmarkAccount[]>([]);

const statusOptions = ['待分析', '已分析', '停用'];
const platformOptions = ['抖音', '小红书', '视频号', '其他'];

const form = reactive<BenchmarkAccountPayload>({
  accountName: '',
  accountUrl: '',
  bestPlayCount: 0,
  city: '',
  followerCount: 0,
  industry: '',
  latestHitTitle: '',
  merchantId: 0,
  platform: '抖音',
  remark: '',
  risk: '',
  status: '待分析',
  takeaway: '',
});

const dialogTitle = computed(() =>
  editingId.value ? '编辑对标账号' : '新增对标账号',
);

const currentMerchantName = computed(() => {
  if (routeMerchantName.value) return routeMerchantName.value;
  const current = merchantOptions.value.find((item) => item.id === routeMerchantId.value);
  return current?.name || '';
});
const selectedMerchant = computed(() => {
  if (selectedRows.value.length === 0) return undefined;
  const first = selectedRows.value[0];
  if (!first) return undefined;
  const sameMerchant = selectedRows.value.every(
    (item) => item.merchantId === first.merchantId,
  );
  return sameMerchant ? first : undefined;
});

const rules = {
  accountName: [{ message: '请输入对标账号', required: true, trigger: 'blur' }],
  merchantId: [{ message: '请选择商家', required: true, trigger: 'change' }],
};

function applyRouteQuery() {
  const rawMerchantId = Number(route.query.merchantId || 0);
  routeMerchantId.value = rawMerchantId > 0 ? rawMerchantId : undefined;
  routeMerchantName.value = String(route.query.merchantName || '');
  keyword.value = routeMerchantName.value;
}

function getStatusTagType(status: string) {
  if (status === '已分析') return 'success';
  if (status === '停用') return 'info';
  return 'warning';
}

function formatNumber(value: number) {
  return Number(value || 0).toLocaleString('zh-CN');
}

function resetForm() {
  const current = merchantOptions.value.find((item) => item.id === routeMerchantId.value);
  Object.assign(form, {
    accountName: '',
    accountUrl: '',
    bestPlayCount: 0,
    city: current?.city || '',
    followerCount: 0,
    industry: current?.industry || '',
    latestHitTitle: '',
    merchantId: routeMerchantId.value || 0,
    platform: '抖音',
    remark: '',
    risk: '',
    status: '待分析',
    takeaway: '',
  });
  formRef.value?.clearValidate();
}

function accountToPayload(row: BenchmarkAccount): BenchmarkAccountPayload {
  return {
    accountName: row.accountName,
    accountUrl: row.accountUrl,
    bestPlayCount: row.bestPlayCount,
    city: row.city,
    followerCount: row.followerCount,
    industry: row.industry,
    latestHitTitle: row.latestHitTitle,
    merchantId: row.merchantId,
    platform: row.platform,
    remark: row.remark,
    risk: row.risk,
    status: row.status,
    takeaway: row.takeaway,
  };
}

function getAnalysisResult() {
  return activeAnalysis.value?.result || {};
}

function getAnalysisPatterns() {
  return getAnalysisResult().patterns || [];
}

function getAnalysisRisks() {
  return getAnalysisResult().risks || [];
}

function getAnalysisActions() {
  return getAnalysisResult().suggestions || [];
}

async function loadMerchants() {
  const result = await getMerchantList({ page: 1, size: 100 });
  merchantOptions.value = result.list;
}

async function loadList() {
  loading.value = true;
  try {
    const result = await getBenchmarkAccountList({
      keyword: keyword.value.trim(),
      merchantId: routeMerchantId.value,
      page: pagination.page,
      size: pagination.size,
      status: statusFilter.value,
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
  router.replace('/lines/list');
  search();
}

function openCreate() {
  editingId.value = null;
  resetForm();
  dialogVisible.value = true;
}

function openEdit(row: BenchmarkAccount) {
  editingId.value = row.id;
  Object.assign(form, accountToPayload(row));
  dialogVisible.value = true;
}

async function submit() {
  await formRef.value?.validate();
  if (!form.merchantId) {
    ElMessage.warning('请选择商家');
    return;
  }
  submitting.value = true;
  try {
    if (editingId.value) {
      await updateBenchmarkAccount(editingId.value, form);
      ElMessage.success('对标账号已更新');
    } else {
      await createBenchmarkAccount(form);
      ElMessage.success('对标账号已创建');
    }
    dialogVisible.value = false;
    await loadList();
  } finally {
    submitting.value = false;
  }
}

async function analyze(row: BenchmarkAccount) {
  analyzingId.value = row.id;
  try {
    activeAnalysis.value = await analyzeBenchmarkAccount(row.id);
    resultVisible.value = true;
    ElMessage.success('对标分析已完成');
    await loadList();
  } finally {
    analyzingId.value = null;
  }
}

async function analyzeSelected() {
  if (selectedRows.value.length === 0) {
    ElMessage.warning('请先勾选对标账号');
    return;
  }
  if (!selectedMerchant.value) {
    ElMessage.warning('一次只能分析同一个商家的对标账号');
    return;
  }
  analyzingId.value = -1;
  try {
    activeAnalysis.value = await analyzeBenchmarkAccounts({
      benchmarkIds: selectedRows.value.map((item) => item.id),
      merchantId: selectedMerchant.value.merchantId,
    });
    resultVisible.value = true;
    ElMessage.success(`已分析 ${selectedRows.value.length} 个对标账号`);
    await loadList();
  } finally {
    analyzingId.value = null;
  }
}

async function removeAccount(row: BenchmarkAccount) {
  await ElMessageBox.confirm(
    `确认删除对标账号「${row.accountName}」？`,
    '删除对标账号',
    { type: 'warning' },
  );
  await deleteBenchmarkAccount(row.id);
  ElMessage.success('对标账号已删除');
  await loadList();
}

function goTopics(row: BenchmarkAccount) {
  router.push({
    path: '/content/notices',
    query: {
      benchmarkAccountId: row.id,
      benchmarkAccountName: row.accountName,
      merchantId: row.merchantId,
      merchantName: row.merchantName,
    },
  });
}

function goSelectedTopics() {
  if (selectedRows.value.length === 0) {
    ElMessage.warning('请先勾选对标账号');
    return;
  }
  if (!selectedMerchant.value) {
    ElMessage.warning('一次只能基于同一个商家的对标账号生成选题');
    return;
  }
  router.push({
    path: '/content/notices',
    query: {
      benchmarkAccountIds: selectedRows.value.map((item) => item.id).join(','),
      benchmarkAccountNames: selectedRows.value
        .map((item) => item.accountName)
        .join('、'),
      merchantId: selectedMerchant.value.merchantId,
      merchantName: selectedMerchant.value.merchantName,
    },
  });
}

function handleSelectionChange(rows: BenchmarkAccount[]) {
  selectedRows.value = rows;
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
    <el-card>
      <template #header>
        <div class="page-head">
          <div>
            <div class="page-title">对标账号库</div>
            <div class="page-desc">
              先人工录入同城、同行或全国优秀账号，分析后再生成选题。
            </div>
          </div>
          <div class="page-actions">
            <el-button @click="showAll">显示全部</el-button>
            <el-button
              :disabled="selectedRows.length === 0"
              :loading="analyzingId === -1"
              @click="analyzeSelected"
            >
              选中分析
            </el-button>
            <el-button
              :disabled="selectedRows.length === 0"
              type="success"
              @click="goSelectedTopics"
            >
              选中生成选题
            </el-button>
            <el-button type="primary" @click="openCreate">新增对标</el-button>
          </div>
        </div>
      </template>

      <el-alert
        v-if="currentMerchantName"
        class="mb-4"
        type="info"
        :closable="false"
        show-icon
      >
        当前处理商家：{{ currentMerchantName }}。已自动带入筛选条件，商家ID：{{ routeMerchantId || '-' }}。
      </el-alert>

      <el-form :inline="true" class="search-form">
        <el-form-item label="关键词">
          <el-input
            v-model="keyword"
            clearable
            placeholder="商家 / 对标账号 / 城市 / 行业 / 爆款 / 可抄点"
            style="width: 360px"
            @keyup.enter="search"
          />
        </el-form-item>
        <el-form-item label="状态">
          <el-select
            v-model="statusFilter"
            clearable
            placeholder="全部状态"
            style="width: 130px"
          >
            <el-option
              v-for="status in statusOptions"
              :key="status"
              :label="status"
              :value="status"
            />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="search">搜索</el-button>
        </el-form-item>
      </el-form>

      <el-table
        v-loading="loading"
        :data="list"
        border
        stripe
        @selection-change="handleSelectionChange"
      >
        <el-table-column type="selection" width="48" />
        <el-table-column label="对标账号" min-width="220">
          <template #default="{ row }">
            <div class="account-name">{{ row.accountName }}</div>
            <div class="table-sub">
              {{ row.platform }} · {{ row.city || '-' }} · {{ row.industry || '-' }}
            </div>
          </template>
        </el-table-column>
        <el-table-column prop="merchantName" label="对应商家" min-width="160" />
        <el-table-column label="数据" width="170">
          <template #default="{ row }">
            <div>粉丝：{{ formatNumber(row.followerCount) }}</div>
            <div class="table-sub">最高播放：{{ formatNumber(row.bestPlayCount) }}</div>
          </template>
        </el-table-column>
        <el-table-column prop="latestHitTitle" label="最近爆款" min-width="220" />
        <el-table-column prop="takeaway" label="可抄点" min-width="220" />
        <el-table-column prop="risk" label="风险提醒" min-width="200" />
        <el-table-column label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="getStatusTagType(row.status)">{{ row.status }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="390" fixed="right">
          <template #default="{ row }">
            <el-button size="small" @click="openEdit(row)">编辑</el-button>
            <el-button
              size="small"
              :loading="analyzingId === row.id"
              @click="analyze(row)"
            >
              分析爆款
            </el-button>
            <el-button size="small" type="primary" @click="goTopics(row)">
              生成选题
            </el-button>
            <el-button size="small" type="danger" @click="removeAccount(row)">
              删除
            </el-button>
          </template>
        </el-table-column>
      </el-table>

      <el-empty v-if="!loading && list.length === 0" description="当前没有对标账号">
        <el-button type="primary" @click="openCreate">新增对标</el-button>
      </el-empty>

      <div class="pagination-wrap">
        <el-pagination
          v-model:current-page="pagination.page"
          v-model:page-size="pagination.size"
          background
          layout="total, sizes, prev, pager, next"
          :page-sizes="[10, 20, 50]"
          :total="pagination.total"
          @current-change="handlePageChange"
          @size-change="handleSizeChange"
        />
      </div>
    </el-card>

    <el-dialog v-model="dialogVisible" :title="dialogTitle" width="860px">
      <el-form ref="formRef" :model="form" :rules="rules" label-width="110px">
        <el-row :gutter="16">
          <el-col :span="12">
            <el-form-item label="对应商家" prop="merchantId">
              <el-select
                v-model="form.merchantId"
                filterable
                placeholder="选择商家"
                style="width: 100%"
              >
                <el-option
                  v-for="merchant in merchantOptions"
                  :key="merchant.id"
                  :label="merchant.name"
                  :value="merchant.id"
                />
              </el-select>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="对标账号" prop="accountName">
              <el-input v-model="form.accountName" placeholder="例如：@贵阳火锅局" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="平台">
              <el-select v-model="form.platform" style="width: 100%">
                <el-option
                  v-for="platform in platformOptions"
                  :key="platform"
                  :label="platform"
                  :value="platform"
                />
              </el-select>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="城市">
              <el-input v-model="form.city" placeholder="例如：贵阳" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="行业/赛道">
              <el-input v-model="form.industry" placeholder="火锅 / 烘焙 / 本地生活" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="账号链接">
              <el-input v-model="form.accountUrl" placeholder="主页链接，可先不填" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="粉丝数">
              <el-input-number
                v-model="form.followerCount"
                :min="0"
                style="width: 100%"
              />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="最高播放">
              <el-input-number
                v-model="form.bestPlayCount"
                :min="0"
                style="width: 100%"
              />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="状态">
              <el-select v-model="form.status" style="width: 100%">
                <el-option
                  v-for="status in statusOptions"
                  :key="status"
                  :label="status"
                  :value="status"
                />
              </el-select>
            </el-form-item>
          </el-col>
          <el-col :span="24">
            <el-form-item label="最近爆款">
              <el-input
                v-model="form.latestHitTitle"
                placeholder="例如：90 元吃到撑的老火锅"
              />
            </el-form-item>
          </el-col>
          <el-col :span="24">
            <el-form-item label="可抄点">
              <el-input
                v-model="form.takeaway"
                placeholder="只写结构和方法，不写照搬文案"
                type="textarea"
                :rows="3"
              />
            </el-form-item>
          </el-col>
          <el-col :span="24">
            <el-form-item label="风险提醒">
              <el-input
                v-model="form.risk"
                placeholder="例如：靠低价爆量，不适合直接复制"
                type="textarea"
                :rows="3"
              />
            </el-form-item>
          </el-col>
          <el-col :span="24">
            <el-form-item label="备注">
              <el-input v-model="form.remark" type="textarea" :rows="3" />
            </el-form-item>
          </el-col>
        </el-row>
      </el-form>

      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="submitting" @click="submit">
          保存对标
        </el-button>
      </template>
    </el-dialog>

    <el-dialog v-model="resultVisible" title="对标分析结果" width="780px">
      <template v-if="activeAnalysis">
        <el-alert
          class="mb-4"
          type="success"
          :closable="false"
          :title="String(activeAnalysis.result?.summary || '分析完成')"
        />
        <el-row :gutter="16">
          <el-col :span="12">
            <div class="block-title">可复用结构</div>
            <div class="result-list">
              <div v-for="item in getAnalysisPatterns()" :key="item" class="result-item">
                {{ item }}
              </div>
            </div>
          </el-col>
          <el-col :span="12">
            <div class="block-title">风险提醒</div>
            <div class="result-list">
              <div v-for="item in getAnalysisRisks()" :key="item" class="result-item warning">
                {{ item }}
              </div>
            </div>
          </el-col>
        </el-row>
        <div class="block-title mt-4">建议动作</div>
        <div class="result-list">
          <div v-for="item in getAnalysisActions()" :key="item" class="result-item action">
            {{ item }}
          </div>
        </div>
      </template>
    </el-dialog>
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

.account-name {
  color: #0f172a;
  font-weight: 700;
}

.table-sub {
  margin-top: 4px;
  color: #64748b;
  font-size: 12px;
}

.pagination-wrap {
  display: flex;
  justify-content: flex-end;
  margin-top: 16px;
}

.block-title {
  margin-bottom: 10px;
  color: #0f172a;
  font-weight: 700;
}

.result-list {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.result-item {
  border-radius: 10px;
  background: #eff6ff;
  color: #1d4ed8;
  line-height: 1.6;
  padding: 10px 12px;
}

.result-item.warning {
  background: #fff7ed;
  color: #c2410c;
}

.result-item.action {
  background: #f0fdf4;
  color: #15803d;
}
</style>
