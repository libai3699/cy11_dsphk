<script setup lang="ts">
import { computed, onMounted, reactive, ref, watch } from 'vue';
import { ElMessage } from 'element-plus';
import { useRoute, useRouter } from 'vue-router';

import {
  generateContentTopics,
  getBenchmarkAccountList,
  getContentTopicList,
  getMerchantList,
  updateContentTopicStatus,
  type BenchmarkAccount,
  type ContentTopic,
  type GenerateTopicsPayload,
  type Merchant,
} from '#/api/admin';

defineOptions({ name: 'ContentNoticesPage' });

const route = useRoute();
const router = useRouter();

const loading = ref(false);
const generating = ref(false);
const dialogVisible = ref(false);
const keyword = ref('');
const statusFilter = ref('');
const routeMerchantId = ref<number | undefined>();
const routeMerchantName = ref('');
const routeBenchmarkId = ref<number | undefined>();
const routeBenchmarkIds = ref<number[]>([]);
const routeBenchmarkName = ref('');
const pagination = reactive({ page: 1, size: 10, total: 0 });
const list = ref<ContentTopic[]>([]);
const merchantOptions = ref<Merchant[]>([]);
const benchmarkOptions = ref<BenchmarkAccount[]>([]);

const statusOptions = ['待确认', '已采用', '停用'];

const form = reactive<GenerateTopicsPayload>({
  benchmarkId: undefined,
  benchmarkIds: [],
  benchmarkName: '',
  cityHotspots: [],
  extraRequirement: '',
  industryHotspots: [],
  merchantId: 0,
  nationalHotspots: [],
});

const cityHotspotsText = ref('');
const industryHotspotsText = ref('');
const nationalHotspotsText = ref('');

const currentMerchantName = computed(() => {
  if (routeMerchantName.value) return routeMerchantName.value;
  const current = merchantOptions.value.find((item) => item.id === routeMerchantId.value);
  return current?.name || '';
});
const selectedBenchmarkNames = computed(() => {
  if (routeBenchmarkName.value) return routeBenchmarkName.value;
  const names = benchmarkOptions.value
    .filter((item) => form.benchmarkIds?.includes(item.id))
    .map((item) => item.accountName);
  return names.join('、');
});

function applyRouteQuery() {
  const rawMerchantId = Number(route.query.merchantId || 0);
  const rawBenchmarkId = Number(route.query.benchmarkAccountId || route.query.benchmarkId || 0);
  const rawBenchmarkIds = String(route.query.benchmarkAccountIds || route.query.benchmarkIds || '')
    .split(',')
    .map((item) => Number(item))
    .filter((item) => item > 0);
  routeMerchantId.value = rawMerchantId > 0 ? rawMerchantId : undefined;
  routeMerchantName.value = String(route.query.merchantName || '');
  routeBenchmarkId.value = rawBenchmarkId > 0 ? rawBenchmarkId : undefined;
  routeBenchmarkIds.value = rawBenchmarkIds.length
    ? rawBenchmarkIds
    : routeBenchmarkId.value
      ? [routeBenchmarkId.value]
      : [];
  routeBenchmarkName.value = String(route.query.benchmarkAccountNames || route.query.benchmarkAccountName || route.query.benchmarkName || '');
  keyword.value = routeMerchantName.value;
  form.merchantId = routeMerchantId.value || 0;
  form.benchmarkId = routeBenchmarkId.value;
  form.benchmarkIds = routeBenchmarkIds.value;
  form.benchmarkName = routeBenchmarkName.value;
}

function splitLines(value: string) {
  return value
    .split('\n')
    .map((item) => item.trim())
    .filter(Boolean);
}

function getRiskTagType(riskLevel: string) {
  if (riskLevel === 'high') return 'danger';
  if (riskLevel === 'medium') return 'warning';
  return 'success';
}

async function loadMerchants() {
  const result = await getMerchantList({ page: 1, size: 100 });
  merchantOptions.value = result.list;
}

async function loadBenchmarks() {
  const result = await getBenchmarkAccountList({
    merchantId: form.merchantId || routeMerchantId.value,
    page: 1,
    size: 100,
  });
  benchmarkOptions.value = result.list;
}

async function loadList() {
  loading.value = true;
  try {
    const result = await getContentTopicList({
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
  routeBenchmarkId.value = undefined;
  routeBenchmarkName.value = '';
  router.replace('/content/notices');
  search();
}

async function openGenerate() {
  form.merchantId = routeMerchantId.value || form.merchantId || 0;
  form.benchmarkId = routeBenchmarkId.value;
  form.benchmarkIds = routeBenchmarkIds.value;
  form.benchmarkName = routeBenchmarkName.value;
  await loadBenchmarks();
  dialogVisible.value = true;
}

async function submitGenerate() {
  if (generating.value) return;
  if (!form.merchantId) {
    ElMessage.warning('请选择商家');
    return;
  }
  generating.value = true;
  try {
    const selectedBenchmarkNames = benchmarkOptions.value
      .filter((item) => form.benchmarkIds?.includes(item.id))
      .map((item) => item.accountName)
      .join('、');
    const result = await generateContentTopics({
      benchmarkId: form.benchmarkIds?.[0] || form.benchmarkId,
      benchmarkIds: form.benchmarkIds,
      benchmarkName: selectedBenchmarkNames || form.benchmarkName,
      cityHotspots: splitLines(cityHotspotsText.value),
      extraRequirement: form.extraRequirement,
      industryHotspots: splitLines(industryHotspotsText.value),
      merchantId: form.merchantId,
      nationalHotspots: splitLines(nationalHotspotsText.value),
    });
    ElMessage.success(`找爆款完成，生成 ${result.topics?.length || 0} 条选题`);
    dialogVisible.value = false;
    routeMerchantId.value = form.merchantId;
    await loadList();
  } catch (error) {
    ElMessage.error(getGenerateErrorMessage(error));
  } finally {
    generating.value = false;
  }
}

function getGenerateErrorMessage(error: unknown) {
  const maybeError = error as {
    error?: string;
    message?: string;
    response?: { data?: { error?: string; message?: string } };
  };
  const responseMessage =
    maybeError?.response?.data?.error || maybeError?.response?.data?.message;
  const message = maybeError?.error || maybeError?.message || responseMessage;
  if (message?.includes('timeout')) {
    return '找爆款超时：模型接口超过 120 秒未返回，请稍后重试';
  }
  return message ? `找爆款失败：${message}` : '找爆款失败：模型接口返回错误';
}

async function acceptTopic(row: ContentTopic) {
  await updateContentTopicStatus(row.id, '已采用');
  ElMessage.success('选题已采用');
  await loadList();
}

function goScripts(row: ContentTopic) {
  router.push({
    path: '/content/quotes',
    query: {
      merchantId: row.merchantId,
      merchantName: row.merchantName,
      topicId: row.id,
      topicTitle: row.title,
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
  async () => {
    applyRouteQuery();
    pagination.page = 1;
    await loadBenchmarks();
    await loadList();
  },
);

watch(
  () => form.merchantId,
  async (value) => {
    if (!value) return;
    await loadBenchmarks();
  },
);

onMounted(async () => {
  applyRouteQuery();
  await loadMerchants();
  await loadBenchmarks();
  await loadList();
});
</script>

<template>
  <div class="p-4">
    <el-card>
      <template #header>
        <div class="page-head">
          <div>
            <div class="page-title">选题中心</div>
            <div class="page-desc">
              这里开始真正找爆款：基于商家、套餐、对标账号和热点，生成可执行选题。
            </div>
          </div>
          <div class="page-actions">
            <el-button @click="showAll">显示全部</el-button>
            <el-button type="primary" @click="openGenerate">找爆款</el-button>
          </div>
        </div>
      </template>

      <el-alert
        v-if="currentMerchantName || selectedBenchmarkNames"
        class="mb-4"
        type="info"
        :closable="false"
        show-icon
      >
        当前处理商家：{{ currentMerchantName || '-' }}。参考对标账号：{{ selectedBenchmarkNames || '-' }}。
      </el-alert>

      <el-form :inline="true" class="search-form">
        <el-form-item label="关键词">
          <el-input
            v-model="keyword"
            clearable
            placeholder="商家 / 对标 / 选题 / 钩子 / 角度 / 转化目标"
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

      <el-table v-loading="loading" :data="list" border stripe>
        <el-table-column prop="merchantName" label="商家" min-width="150" />
        <el-table-column label="选题" min-width="260">
          <template #default="{ row }">
            <div class="topic-title">{{ row.title }}</div>
            <div class="table-sub">{{ row.hook }}</div>
          </template>
        </el-table-column>
        <el-table-column prop="angle" label="内容角度" min-width="180" />
        <el-table-column prop="target" label="转化目标" min-width="140" />
        <el-table-column label="风险" width="100">
          <template #default="{ row }">
            <el-tag :type="getRiskTagType(row.riskLevel)">
              {{ row.riskLevel }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="benchmarkName" label="参考对标" min-width="140" />
        <el-table-column prop="recommendReason" label="推荐理由" min-width="240" />
        <el-table-column label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="row.status === '已采用' ? 'success' : 'warning'">
              {{ row.status }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="260" fixed="right">
          <template #default="{ row }">
            <el-button size="small" @click="acceptTopic(row)">采用</el-button>
            <el-button size="small" type="primary" @click="goScripts(row)">
              生成文案
            </el-button>
          </template>
        </el-table-column>
      </el-table>

      <el-empty v-if="!loading && list.length === 0" description="当前没有选题">
        <el-button type="primary" @click="openGenerate">找爆款</el-button>
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

    <el-dialog v-model="dialogVisible" title="找爆款" width="820px">
      <el-form label-width="120px">
        <el-row :gutter="16">
          <el-col :span="12">
            <el-form-item label="商家">
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
            <el-form-item label="参考对标">
              <el-select
                v-model="form.benchmarkIds"
                clearable
                filterable
                multiple
                collapse-tags
                collapse-tags-tooltip
                placeholder="可多选，建议选 5-10 个已分析对标"
                style="width: 100%"
              >
                <el-option
                  v-for="item in benchmarkOptions"
                  :key="item.id"
                  :label="item.accountName"
                  :value="item.id"
                />
              </el-select>
              <div class="form-tip">
                已选：{{ selectedBenchmarkNames || '未选择，Agent 将只根据商家、套餐和热点保守生成' }}
              </div>
            </el-form-item>
          </el-col>
          <el-col :span="24">
            <el-form-item label="同城热点">
              <el-input
                v-model="cityHotspotsText"
                placeholder="一行一个，例如：贵阳降温第一顿火锅 / 花果园夜宵"
                type="textarea"
                :rows="3"
              />
            </el-form-item>
          </el-col>
          <el-col :span="24">
            <el-form-item label="行业热点">
              <el-input
                v-model="industryHotspotsText"
                placeholder="一行一个，例如：牛油锅底 / 低价套餐 / 老板出镜"
                type="textarea"
                :rows="3"
              />
            </el-form-item>
          </el-col>
          <el-col :span="24">
            <el-form-item label="全国热点">
              <el-input
                v-model="nationalHotspotsText"
                placeholder="一行一个，可以不填"
                type="textarea"
                :rows="3"
              />
            </el-form-item>
          </el-col>
          <el-col :span="24">
            <el-form-item label="额外要求">
              <el-input
                v-model="form.extraRequirement"
                placeholder="例如：不要低价伤利润，优先老板出镜，围绕双人套餐"
                type="textarea"
                :rows="3"
              />
            </el-form-item>
          </el-col>
        </el-row>
      </el-form>

      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button
          type="primary"
          :disabled="generating"
          :loading="generating"
          @click="submitGenerate"
        >
          开始找爆款
        </el-button>
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

.topic-title {
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

.form-tip {
  margin-top: 6px;
  color: #64748b;
  font-size: 12px;
  line-height: 1.5;
}
</style>
