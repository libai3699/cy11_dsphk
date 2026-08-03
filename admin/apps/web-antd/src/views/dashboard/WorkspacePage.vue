<script setup lang="ts">
import { computed, onMounted, ref } from 'vue';
import { ElMessage } from 'element-plus';
import { useRouter } from 'vue-router';

import {
  getContentReviews,
  getContentTopicList,
  getMerchantList,
  getMerchantPackageList,
  getPublishSchedules,
  getShootingTasks,
  type Merchant,
  type WorkspaceOverview,
  type WorkspaceShootingTask,
  type WorkspaceTopic,
} from '#/api/admin';

defineOptions({ name: 'WorkspacePage' });

const router = useRouter();

const loading = ref(false);
const overview = ref<WorkspaceOverview>({
  merchants: [],
  metrics: [],
  reviews: [],
  shootingTasks: [],
  topics: [],
});

const workflowSteps = [
  { label: '商家建档', path: '/users/list' },
  { label: '账号授权', path: '/users/devices' },
  { label: '对标分析', path: '/lines/list' },
  { label: '生成选题', path: '/content/notices' },
  { label: '生成文案', path: '/content/quotes' },
  { label: '生成分镜', path: '/content/discoveries' },
  { label: '派拍摄', path: '/logs/user' },
  { label: '排发布', path: '/content/payments' },
  { label: '做复盘', path: '/logs/admin' },
] as const;

const firstMerchant = computed(() => overview.value.merchants[0]);

async function loadOverview() {
  loading.value = true;
  try {
    overview.value = await buildOverviewFromExistingApis();
  } catch (error) {
    const maybeError = error as { error?: string; message?: string };
    ElMessage.error(maybeError?.error || maybeError?.message || '读取工作台数据失败');
  } finally {
    loading.value = false;
  }
}

async function buildOverviewFromExistingApis(): Promise<WorkspaceOverview> {
  const [
    merchantResult,
    packageResult,
    shootingResult,
    topicResult,
    scheduleResult,
    reviewResult,
  ] = await Promise.all([
    getMerchantList({ page: 1, size: 100 }),
    getMerchantPackageList({ page: 1, size: 100 }),
    getShootingTasks({ page: 1, size: 100 }),
    getContentTopicList({ page: 1, size: 6 }),
    getPublishSchedules({ page: 1, size: 100 }),
    getContentReviews({ page: 1, size: 100 }),
  ]);

  const merchants = merchantResult.list || [];
  const packages = packageResult.list || [];
  const shootingTasks = shootingResult.list || [];
  const schedules = scheduleResult.list || [];
  const reviews = reviewResult.list || [];
  const recentStart = Date.now() - 7 * 24 * 60 * 60 * 1000;
  const merchantMap = new Map(merchants.map((item) => [item.id, item]));
  const recentReviews = reviews.filter((item) => {
    const time = new Date(item.createdAt || '').getTime();
    return Number.isFinite(time) && time >= recentStart;
  });
  const recentWriteOffAmount = recentReviews.reduce(
    (sum, item) => sum + Number(item.writeOffAmount || 0),
    0,
  );
  const estimatedCommission = recentReviews.reduce((sum, item) => {
    const merchant = merchantMap.get(item.merchantId);
    const rate = merchant?.commissionRate || 10;
    return sum + Number(item.writeOffAmount || 0) * rate / 100;
  }, 0);
  const activeMerchantCount = merchants.filter((item) => item.status === 1).length;
  const enabledPackageCount = packages.filter((item) => item.status === 1).length;
  const pendingShootingCount = shootingTasks.filter((item) =>
    ['待拍摄', '拍摄中', '已拍摄', '已剪辑'].includes(item.status),
  ).length;
  const pendingScheduleCount = schedules.filter((item) => item.status === '待发布').length;

  return {
    merchants: merchants.slice(0, 6).map((item) =>
      buildWorkspaceMerchant(item, packages, reviews),
    ),
    metrics: [
      {
        hint: `总商家 ${merchantResult.total} 家，启用套餐 ${enabledPackageCount} 个`,
        key: 'merchants',
        label: '在运营商家',
        path: '/users/list',
        value: String(activeMerchantCount),
      },
      {
        hint: `待发布视频 ${pendingScheduleCount} 条`,
        key: 'shooting',
        label: '待拍/待剪任务',
        path: '/logs/user',
        value: String(pendingShootingCount),
      },
      {
        hint: `已复盘 ${reviewResult.total} 条视频`,
        key: 'writeoff',
        label: '近 7 天核销额',
        path: '/logs/admin',
        value: formatMoney(recentWriteOffAmount),
      },
      {
        hint: '按各商家分成比例估算',
        key: 'commission',
        label: '预估分成',
        path: '/plans/orders',
        value: formatMoney(estimatedCommission),
      },
    ],
    reviews: reviews.slice(0, 6).map((item) => ({
      dealCount: item.dealCount,
      id: item.id,
      merchantId: item.merchantId,
      merchantName: item.merchantName,
      playCount: item.playCount,
      status: item.status,
      videoTitle: item.videoTitle,
      writeOffAmount: item.writeOffAmount,
    })),
    shootingTasks: shootingTasks
      .filter((item) => item.status !== '已完成')
      .slice(0, 6)
      .map((item) => ({
        assignee: item.assignee,
        deadline: item.deadline,
        id: item.id,
        merchantId: item.merchantId,
        merchantName: item.merchantName,
        shotCount: item.shotCount,
        status: item.status,
        taskTitle: item.taskTitle,
      })),
    topics: (topicResult.list || []).map((item) => ({
      hook: item.hook,
      id: item.id,
      merchantId: item.merchantId,
      merchantName: item.merchantName,
      publishWindow: item.publishWindow,
      status: item.status,
      title: item.title,
    })),
  };
}

function buildWorkspaceMerchant(
  item: Merchant,
  packages: { merchantId: number; status: number }[],
  reviews: { createdAt: string; merchantId: number; writeOffAmount: number }[],
) {
  const recentStart = Date.now() - 7 * 24 * 60 * 60 * 1000;
  const recentWriteOffAmount = reviews
    .filter((review) => {
      const time = new Date(review.createdAt || '').getTime();
      return (
        review.merchantId === item.id &&
        Number.isFinite(time) &&
        time >= recentStart
      );
    })
    .reduce((sum, review) => sum + Number(review.writeOffAmount || 0), 0);
  return {
    city: item.city,
    commissionRate: item.commissionRate,
    estimatedCommission: recentWriteOffAmount * (item.commissionRate || 10) / 100,
    id: item.id,
    industry: item.industry,
    name: item.name,
    nextAction: buildMerchantNextAction(item, packages),
    recentWriteOffAmount,
    stage: item.stage,
  };
}

function buildMerchantNextAction(
  item: Merchant,
  packages: { merchantId: number; status: number }[],
) {
  if (!item.industry || !item.city || !item.contactName || !item.contactPhone) {
    return '补齐基础档案';
  }
  if (!item.cooperationType || item.commissionRate <= 0) {
    return '确认合作规则';
  }
  if (!packages.some((pkg) => pkg.merchantId === item.id && pkg.status === 1)) {
    return '新增启用套餐';
  }
  if (!item.douyinAccount && !item.douyinLaikeAccount) {
    return '记录账号授权';
  }
  return '推进选题和内容生产';
}

function goMetric(path: string) {
  router.push(path);
}

function goMerchantProcess(id: number) {
  router.push(`/users/detail/${id}`);
}

function goFirstMerchantProcess() {
  if (!firstMerchant.value) return;
  goMerchantProcess(firstMerchant.value.id);
}

function goTopics(row?: { id?: number; merchantId?: number; name?: string }) {
  router.push({
    path: '/content/notices',
    query: row?.merchantId
      ? { merchantId: row.merchantId, merchantName: row.name || '' }
      : undefined,
  });
}

function goScripts(row: WorkspaceTopic) {
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

function goShooting(row?: WorkspaceShootingTask) {
  router.push({
    path: '/logs/user',
    query: row
      ? {
          merchantId: row.merchantId,
          merchantName: row.merchantName,
        }
      : undefined,
  });
}

function goReview(merchantId?: number, merchantName?: string) {
  router.push({
    path: '/logs/admin',
    query: merchantId ? { merchantId, merchantName } : undefined,
  });
}

function formatMoney(value: number) {
  return `¥${Number(value || 0).toFixed(2)}`;
}

function formatNumber(value: number) {
  return Number(value || 0).toLocaleString('zh-CN');
}

function formatTime(value?: string) {
  if (!value) return '未设置截止时间';
  return value.replace('T', ' ').slice(0, 16);
}

onMounted(loadOverview);
</script>

<template>
  <div v-loading="loading" class="workspace-page">
    <el-card class="hero-card" shadow="never">
      <div class="hero-head">
        <div>
          <div class="eyebrow">真实运营总览</div>
          <h2>短视频获客工作台</h2>
          <p>从数据库汇总当前商家、内容生产、拍摄、发布和复盘状态。这里负责看全局，不再展示写死样例。</p>
        </div>
        <div class="hero-actions">
          <el-button type="primary" @click="router.push('/users/list')">新增商家</el-button>
          <el-button :disabled="!firstMerchant" @click="goFirstMerchantProcess">
            查看首个商家流程
          </el-button>
        </div>
      </div>
      <div class="stats-grid">
        <button
          v-for="item in overview.metrics"
          :key="item.key"
          class="stat-card"
          type="button"
          @click="goMetric(item.path)"
        >
          <div class="stat-label">{{ item.label }}</div>
          <div class="stat-value">{{ item.value }}</div>
          <div class="stat-hint">{{ item.hint }}</div>
        </button>
      </div>
    </el-card>

    <el-card shadow="never">
      <template #header>
        <div class="card-header">
          <span>运营流程</span>
          <el-tag>按真实页面跳转</el-tag>
        </div>
      </template>
      <div class="workflow-strip">
        <el-button
          v-for="(step, index) in workflowSteps"
          :key="step.label"
          plain
          type="primary"
          @click="router.push(step.path)"
        >
          {{ index + 1 }}. {{ step.label }}
        </el-button>
      </div>
    </el-card>

    <el-row :gutter="16" class="workspace-row">
      <el-col :span="14">
        <el-card shadow="never">
          <template #header>
            <div class="card-header">
              <span>重点商家</span>
              <el-tag type="success">{{ overview.merchants.length }} 家</el-tag>
            </div>
          </template>
          <el-table :data="overview.merchants" border stripe>
            <el-table-column prop="name" label="商家" min-width="160" />
            <el-table-column prop="industry" label="行业" width="100" />
            <el-table-column prop="stage" label="阶段" min-width="120" />
            <el-table-column label="近 7 天核销" width="130">
              <template #default="{ row }">{{ formatMoney(row.recentWriteOffAmount) }}</template>
            </el-table-column>
            <el-table-column label="预估分成" width="120">
              <template #default="{ row }">{{ formatMoney(row.estimatedCommission) }}</template>
            </el-table-column>
            <el-table-column prop="nextAction" label="当前应处理" min-width="180" />
            <el-table-column label="操作" width="250" fixed="right">
              <template #default="{ row }">
                <el-button size="small" @click="goMerchantProcess(row.id)">流程</el-button>
                <el-button
                  size="small"
                  type="primary"
                  @click="goTopics({ merchantId: row.id, name: row.name })"
                >
                  做选题
                </el-button>
              </template>
            </el-table-column>
          </el-table>
          <el-empty
            v-if="overview.merchants.length === 0"
            description="还没有商家，先去商家建档"
          >
            <el-button type="primary" @click="router.push('/users/list')">去建档</el-button>
          </el-empty>
        </el-card>
      </el-col>
      <el-col :span="10">
        <el-card shadow="never">
          <template #header>
            <div class="card-header">
              <span>待处理拍摄</span>
              <el-tag type="warning">{{ overview.shootingTasks.length }} 条</el-tag>
            </div>
          </template>
          <div v-if="overview.shootingTasks.length > 0" class="todo-list">
            <div v-for="task in overview.shootingTasks" :key="task.id" class="todo-item">
              <div class="todo-title">{{ task.taskTitle }}</div>
              <div class="todo-meta">
                {{ task.merchantName }} · {{ formatTime(task.deadline) }} · {{ task.assignee || '未分配' }}
              </div>
              <div class="todo-meta">镜头数：{{ task.shotCount }} · 状态：{{ task.status }}</div>
              <el-button class="todo-button" size="small" type="primary" @click="goShooting(task)">
                查看任务
              </el-button>
            </div>
          </div>
          <el-empty v-else description="暂无待处理拍摄任务">
            <el-button @click="goShooting()">去拍摄任务</el-button>
          </el-empty>
        </el-card>
      </el-col>
    </el-row>

    <el-row :gutter="16" class="workspace-row">
      <el-col :span="12">
        <el-card shadow="never">
          <template #header>
            <div class="card-header">
              <span>选题池</span>
              <el-tag>{{ overview.topics.length }} 条</el-tag>
            </div>
          </template>
          <div v-if="overview.topics.length > 0" class="topic-list">
            <div v-for="item in overview.topics" :key="item.id" class="topic-item">
              <div class="topic-title">{{ item.title }}</div>
              <div class="topic-hook">{{ item.hook }}</div>
              <div class="topic-meta">
                {{ item.merchantName }} · {{ item.publishWindow || '未设置发布窗口' }} · {{ item.status }}
              </div>
              <el-button class="todo-button" size="small" type="primary" @click="goScripts(item)">
                生成文案
              </el-button>
            </div>
          </div>
          <el-empty v-else description="暂无选题">
            <el-button type="primary" @click="goTopics()">去找爆款</el-button>
          </el-empty>
        </el-card>
      </el-col>
      <el-col :span="12">
        <el-card shadow="never">
          <template #header>
            <div class="card-header">
              <span>最近复盘</span>
              <el-tag type="info">{{ overview.reviews.length }} 条</el-tag>
            </div>
          </template>
          <div v-if="overview.reviews.length > 0" class="review-list">
            <div v-for="item in overview.reviews" :key="item.id" class="review-item">
              <div class="review-title">{{ item.videoTitle }}</div>
              <div class="review-meta">
                {{ item.merchantName }} · 播放 {{ formatNumber(item.playCount) }} · 成交 {{ item.dealCount }} 单
              </div>
              <div class="review-next">核销额：{{ formatMoney(item.writeOffAmount) }} · {{ item.status }}</div>
              <el-button
                class="todo-button"
                size="small"
                @click="goReview(item.merchantId, item.merchantName)"
              >
                查看复盘
              </el-button>
            </div>
          </div>
          <el-empty v-else description="暂无复盘数据">
            <el-button @click="goReview()">去复盘</el-button>
          </el-empty>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<style scoped>
.workspace-page {
  padding: 20px;
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.hero-card {
  border-radius: 20px;
  background:
    radial-gradient(circle at top right, rgba(37, 99, 235, 0.18), transparent 28%),
    linear-gradient(135deg, #0f172a 0%, #1d4ed8 100%);
  border: 0;
  color: #fff;
}

.hero-head {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 16px;
}

.hero-actions {
  display: flex;
  flex-shrink: 0;
  gap: 10px;
}

.eyebrow {
  margin-bottom: 8px;
  color: rgba(255, 255, 255, 0.72);
  font-size: 12px;
  letter-spacing: 0.12em;
}

.hero-head h2 {
  margin: 0;
  font-size: 30px;
}

.hero-head p {
  margin: 8px 0 0;
  max-width: 760px;
  color: rgba(255, 255, 255, 0.82);
}

.stats-grid {
  margin-top: 24px;
  display: grid;
  grid-template-columns: repeat(4, minmax(0, 1fr));
  gap: 14px;
}

.workspace-row {
  margin: 0;
}

.card-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  font-weight: 600;
}

.workflow-strip {
  display: flex;
  flex-wrap: wrap;
  gap: 10px;
}

.stat-card {
  min-width: 0;
  padding: 18px;
  text-align: left;
  background: rgba(255, 255, 255, 0.12);
  border: 1px solid rgba(255, 255, 255, 0.16);
  border-radius: 16px;
  cursor: pointer;
  backdrop-filter: blur(6px);
}

.stat-label {
  color: rgba(255, 255, 255, 0.72);
  font-size: 13px;
}

.stat-value {
  margin-top: 10px;
  color: #fff;
  font-size: 30px;
  font-weight: 700;
}

.stat-hint {
  margin-top: 8px;
  color: rgba(255, 255, 255, 0.68);
  font-size: 12px;
}

.todo-list,
.topic-list,
.review-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.todo-item,
.topic-item,
.review-item {
  border: 1px solid #e5e7eb;
  border-radius: 14px;
  padding: 14px 16px;
  background: #fff;
}

.todo-title,
.topic-title,
.review-title {
  color: #0f172a;
  font-weight: 600;
}

.todo-meta,
.topic-meta,
.review-meta {
  margin-top: 6px;
  color: #64748b;
  font-size: 13px;
}

.topic-hook,
.review-next {
  margin-top: 8px;
  color: #1e293b;
  font-size: 14px;
}

.todo-button {
  margin-top: 10px;
}

@media (max-width: 1200px) {
  .stats-grid {
    grid-template-columns: repeat(2, minmax(0, 1fr));
  }
}
</style>
