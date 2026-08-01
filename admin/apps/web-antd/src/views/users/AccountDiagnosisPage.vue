<script setup lang="ts">
import { computed, onMounted, reactive, ref, watch } from 'vue';
import { ElMessage } from 'element-plus';
import { useRoute, useRouter } from 'vue-router';

import {
  createAccountDiagnosis,
  getAccountDiagnosisList,
  getMerchantAccountAuthList,
  getMerchantList,
  getMerchantPackageList,
  type AccountDiagnosisPayload,
  type AccountDiagnosisTask,
  type Merchant,
  type MerchantAccountAuth,
  type MerchantPackage,
} from '#/api/admin';

defineOptions({ name: 'AccountDiagnosisPage' });

const route = useRoute();
const router = useRouter();

const loading = ref(false);
const submitting = ref(false);
const keyword = ref('');
const statusFilter = ref('');
const routeMerchantId = ref<number | undefined>();
const routeMerchantName = ref('');
const routeAccountAuthId = ref<number | undefined>();
const pagination = reactive({ page: 1, size: 10, total: 0 });

const merchantOptions = ref<Merchant[]>([]);
const accountAuthOptions = ref<MerchantAccountAuth[]>([]);
const packageOptions = ref<MerchantPackage[]>([]);
const list = ref<AccountDiagnosisTask[]>([]);
const activeTask = ref<AccountDiagnosisTask>();
const recentVideosText = ref('');

const statusOptions = ['待执行', '执行中', '已完成', '失败'];
const ownerAppearanceOptions = ['愿意出镜', '不愿出镜', '可偶尔出镜', '待确认'];

const form = reactive<AccountDiagnosisPayload>({
  accountAuthId: undefined,
  avgPlayCount: 0,
  bestVideoPlay: 0,
  bestVideoTitle: '',
  currentProblems: '',
  followerCount: 0,
  merchantId: 0,
  operatorGoal: '',
  ownerAppearance: '待确认',
  recentVideoCount: 0,
  recentVideos: [],
  remark: '',
  targetPackage: '',
});

const currentMerchantName = computed(() => {
  if (routeMerchantName.value) return routeMerchantName.value;
  const current = merchantOptions.value.find((item) => item.id === routeMerchantId.value);
  return current?.name || '';
});

const selectedAuth = computed(() =>
  accountAuthOptions.value.find((item) => item.id === form.accountAuthId),
);

const canSubmit = computed(() => Boolean(form.merchantId));

function applyRouteQuery() {
  const rawMerchantId = Number(route.query.merchantId || 0);
  const rawAccountAuthId = Number(route.query.accountAuthId || 0);
  routeMerchantId.value = rawMerchantId > 0 ? rawMerchantId : undefined;
  routeAccountAuthId.value = rawAccountAuthId > 0 ? rawAccountAuthId : undefined;
  routeMerchantName.value = String(route.query.merchantName || '');
  keyword.value = routeMerchantName.value;
  form.merchantId = routeMerchantId.value || 0;
  form.accountAuthId = routeAccountAuthId.value;
}

function getStatusTagType(status: string) {
  if (status === '已完成') return 'success';
  if (status === '执行中') return 'warning';
  if (status === '失败') return 'danger';
  return 'info';
}

function getScoreType(score?: number) {
  if (!score || score < 60) return 'exception';
  if (score < 80) return 'warning';
  return 'success';
}

function getResult(task?: AccountDiagnosisTask) {
  return task?.result || {};
}

function getProblems(task?: AccountDiagnosisTask) {
  const result = getResult(task);
  return result.problems || [];
}

function getActions(task?: AccountDiagnosisTask) {
  const result = getResult(task);
  return result.nextActions || result.suggestions || [];
}

async function loadMerchants() {
  const result = await getMerchantList({ page: 1, size: 100 });
  merchantOptions.value = result.list;
}

async function loadAccountAuths() {
  const result = await getMerchantAccountAuthList({
    merchantId: form.merchantId || routeMerchantId.value,
    page: 1,
    size: 100,
  });
  accountAuthOptions.value = result.list;
  if (routeAccountAuthId.value) {
    form.accountAuthId = routeAccountAuthId.value;
  } else if (!form.accountAuthId && result.list.length > 0) {
    form.accountAuthId = result.list[0]?.id;
  }
}

async function loadPackages() {
  const result = await getMerchantPackageList({
    merchantId: form.merchantId || routeMerchantId.value,
    page: 1,
    size: 100,
  });
  packageOptions.value = result.list;
  if (!form.targetPackage && result.list.length > 0) {
    form.targetPackage = result.list[0]?.name || '';
  }
}

async function loadList() {
  loading.value = true;
  try {
    const result = await getAccountDiagnosisList({
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
    activeTask.value = result.list[0];
  } finally {
    loading.value = false;
  }
}

async function reloadContext() {
  await loadAccountAuths();
  await loadPackages();
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
  routeAccountAuthId.value = undefined;
  router.replace('/users/account-diagnosis');
  search();
}

async function submitDiagnosis() {
  if (!canSubmit.value) {
    ElMessage.warning('请选择商家');
    return;
  }
  submitting.value = true;
  try {
    const payload: AccountDiagnosisPayload = {
      ...form,
      recentVideos: recentVideosText.value
        .split('\n')
        .map((item) => item.trim())
        .filter(Boolean),
    };
    const task = await createAccountDiagnosis(payload);
    activeTask.value = task;
    ElMessage.success('账号诊断已完成');
    await loadList();
  } finally {
    submitting.value = false;
  }
}

function goBenchmark() {
  const merchantName = currentMerchantName.value;
  router.push({
    path: '/lines/list',
    query: merchantName
      ? {
          merchantId: form.merchantId || routeMerchantId.value,
          merchantName,
        }
      : {},
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
    await reloadContext();
    await loadList();
  },
);

watch(
  () => form.merchantId,
  async (value) => {
    if (!value) return;
    routeMerchantId.value = value;
    await reloadContext();
  },
);

onMounted(async () => {
  applyRouteQuery();
  await loadMerchants();
  await reloadContext();
  await loadList();
});
</script>

<template>
  <div class="p-4">
    <el-row :gutter="16">
      <el-col :span="9">
        <el-card>
          <template #header>
            <div class="section-head">
              <div>
                <div class="page-title">账号诊断</div>
                <div class="page-desc">先录入账号现状，生成诊断报告，再进入对标账号。</div>
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
            当前处理商家：{{ currentMerchantName }}
          </el-alert>

          <el-form label-width="120px">
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
            <el-form-item label="授权账号">
              <el-select
                v-model="form.accountAuthId"
                clearable
                placeholder="选择授权记录"
                style="width: 100%"
              >
                <el-option
                  v-for="auth in accountAuthOptions"
                  :key="auth.id"
                  :label="`${auth.platform} · ${auth.accountName || auth.accountUid || auth.authStatus}`"
                  :value="auth.id"
                />
              </el-select>
              <div v-if="selectedAuth" class="form-tip">
                {{ selectedAuth.authMethod }} · {{ selectedAuth.authStatus }}
              </div>
            </el-form-item>
            <el-form-item label="主推套餐">
              <el-select
                v-model="form.targetPackage"
                allow-create
                clearable
                filterable
                placeholder="选择或输入主推套餐"
                style="width: 100%"
              >
                <el-option
                  v-for="item in packageOptions"
                  :key="item.id"
                  :label="item.name"
                  :value="item.name"
                />
              </el-select>
            </el-form-item>
            <el-row :gutter="12">
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
                <el-form-item label="平均播放">
                  <el-input-number
                    v-model="form.avgPlayCount"
                    :min="0"
                    style="width: 100%"
                  />
                </el-form-item>
              </el-col>
              <el-col :span="12">
                <el-form-item label="最高播放">
                  <el-input-number
                    v-model="form.bestVideoPlay"
                    :min="0"
                    style="width: 100%"
                  />
                </el-form-item>
              </el-col>
              <el-col :span="12">
                <el-form-item label="近期视频数">
                  <el-input-number
                    v-model="form.recentVideoCount"
                    :min="0"
                    style="width: 100%"
                  />
                </el-form-item>
              </el-col>
            </el-row>
            <el-form-item label="最高播放标题">
              <el-input v-model="form.bestVideoTitle" placeholder="账号里表现最好的一条视频" />
            </el-form-item>
            <el-form-item label="老板出镜">
              <el-select v-model="form.ownerAppearance" style="width: 100%">
                <el-option
                  v-for="item in ownerAppearanceOptions"
                  :key="item"
                  :label="item"
                  :value="item"
                />
              </el-select>
            </el-form-item>
            <el-form-item label="当前问题">
              <el-input
                v-model="form.currentProblems"
                placeholder="例如：播放低、没有成交、老板不愿出镜"
                type="textarea"
                :rows="3"
              />
            </el-form-item>
            <el-form-item label="运营目标">
              <el-input
                v-model="form.operatorGoal"
                placeholder="例如：先跑通 5 条内容，带动双人套餐核销"
                type="textarea"
                :rows="3"
              />
            </el-form-item>
            <el-form-item label="热视频标题">
              <el-input
                v-model="recentVideosText"
                placeholder="一行一条，手动粘最近表现好的视频标题"
                type="textarea"
                :rows="4"
              />
            </el-form-item>
            <el-form-item label="备注">
              <el-input v-model="form.remark" type="textarea" :rows="3" />
            </el-form-item>
            <el-form-item>
              <el-button
                type="primary"
                :disabled="!canSubmit"
                :loading="submitting"
                @click="submitDiagnosis"
              >
                开始诊断
              </el-button>
              <el-button @click="goBenchmark">去对标账号</el-button>
            </el-form-item>
          </el-form>
        </el-card>
      </el-col>

      <el-col :span="15">
        <el-card class="mb-4">
          <template #header>
            <div class="section-head">
              <div>
                <div class="page-title">诊断结果</div>
                <div class="page-desc">结果会保存为任务记录，后面给对标和选题使用。</div>
              </div>
            </div>
          </template>

          <template v-if="activeTask">
            <div class="result-head">
              <div>
                <div class="result-title">
                  {{ activeTask.merchantName }} · {{ activeTask.accountName || '账号诊断' }}
                </div>
                <div class="page-desc">{{ activeTask.createdAt }}</div>
              </div>
              <el-tag :type="getStatusTagType(activeTask.status)">
                {{ activeTask.status }}
              </el-tag>
            </div>

            <div class="score-grid">
              <div class="score-card">
                <span>账号基础</span>
                <el-progress
                  :percentage="Number(getResult(activeTask).accountScore || 0)"
                  :status="getScoreType(Number(getResult(activeTask).accountScore || 0))"
                />
              </div>
              <div class="score-card">
                <span>内容表现</span>
                <el-progress
                  :percentage="Number(getResult(activeTask).contentScore || 0)"
                  :status="getScoreType(Number(getResult(activeTask).contentScore || 0))"
                />
              </div>
              <div class="score-card">
                <span>转化准备</span>
                <el-progress
                  :percentage="Number(getResult(activeTask).convertScore || 0)"
                  :status="getScoreType(Number(getResult(activeTask).convertScore || 0))"
                />
              </div>
            </div>

            <el-alert
              class="mb-4"
              type="success"
              :closable="false"
              :title="String(getResult(activeTask).summary || '暂无摘要')"
            />

            <el-row :gutter="16">
              <el-col :span="12">
                <div class="block-title">发现的问题</div>
                <div class="result-list">
                  <div v-for="item in getProblems(activeTask)" :key="item" class="result-item warning">
                    {{ item }}
                  </div>
                </div>
              </el-col>
              <el-col :span="12">
                <div class="block-title">建议动作</div>
                <div class="result-list">
                  <div v-for="item in getActions(activeTask)" :key="item" class="result-item">
                    {{ item }}
                  </div>
                </div>
              </el-col>
            </el-row>
          </template>

          <el-empty v-else description="还没有诊断结果，先在左侧录入账号现状并开始诊断" />
        </el-card>

        <el-card>
          <template #header>
            <div class="section-head">
              <div class="page-title">诊断历史</div>
              <div class="history-actions">
                <el-input
                  v-model="keyword"
                  clearable
                  placeholder="商家 / 账号 / 状态"
                  style="width: 220px"
                  @keyup.enter="search"
                />
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
                <el-button type="primary" @click="search">搜索</el-button>
                <el-button @click="showAll">显示全部</el-button>
              </div>
            </div>
          </template>

          <el-table v-loading="loading" :data="list" border stripe>
            <el-table-column prop="merchantName" label="商家" min-width="150" />
            <el-table-column prop="accountName" label="账号" min-width="150" />
            <el-table-column label="状态" width="100">
              <template #default="{ row }">
                <el-tag :type="getStatusTagType(row.status)">{{ row.status }}</el-tag>
              </template>
            </el-table-column>
            <el-table-column label="评分" width="180">
              <template #default="{ row }">
                {{ getResult(row).accountScore || 0 }} /
                {{ getResult(row).contentScore || 0 }} /
                {{ getResult(row).convertScore || 0 }}
              </template>
            </el-table-column>
            <el-table-column prop="createdAt" label="创建时间" width="180" />
            <el-table-column label="操作" width="120" fixed="right">
              <template #default="{ row }">
                <el-button size="small" type="primary" @click="activeTask = row">
                  查看
                </el-button>
              </template>
            </el-table-column>
          </el-table>

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
      </el-col>
    </el-row>
  </div>
</template>

<style scoped>
.section-head {
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

.page-desc,
.form-tip {
  margin-top: 4px;
  color: #64748b;
  font-size: 13px;
}

.result-head {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  gap: 16px;
  margin-bottom: 16px;
}

.result-title {
  color: #0f172a;
  font-size: 18px;
  font-weight: 700;
}

.score-grid {
  display: grid;
  grid-template-columns: repeat(3, minmax(0, 1fr));
  gap: 12px;
  margin-bottom: 16px;
}

.score-card {
  border: 1px solid #e2e8f0;
  border-radius: 12px;
  padding: 12px;
}

.score-card span,
.block-title {
  display: block;
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

.history-actions {
  display: flex;
  align-items: center;
  gap: 8px;
}

.pagination-wrap {
  display: flex;
  justify-content: flex-end;
  margin-top: 16px;
}
</style>
