<script setup lang="ts">
import { computed, onMounted, ref } from 'vue';
import { ElMessage } from 'element-plus';
import { useRoute, useRouter } from 'vue-router';

import {
  getMerchantWorkspace,
  type Merchant,
  type MerchantWorkspace,
  type MerchantWorkspaceMetrics,
} from '#/api/admin';

defineOptions({ name: 'MerchantProcessPage' });

type StepStatus = 'done' | 'locked' | 'todo';

interface ProcessStep {
  actionText: string;
  desc: string;
  key: string;
  metricText: string;
  missing: string[];
  path: string;
  query?: Record<string, number | string>;
  status: StepStatus;
  title: string;
}

const route = useRoute();
const router = useRouter();
const loading = ref(false);
const workspace = ref<MerchantWorkspace>();

const merchantId = computed(() => Number(route.params.id || 0));
const merchant = computed<Merchant | undefined>(() => workspace.value?.merchant);
const metrics = computed<MerchantWorkspaceMetrics | undefined>(() => workspace.value?.metrics);

const baseQuery = computed<Record<string, number | string> | undefined>(() => {
  const row = merchant.value;
  return row
    ? {
        merchantId: row.id,
        merchantName: row.name,
      }
    : undefined;
});

const requirementMap = computed(() => {
  const result: Record<string, string[]> = {};
  for (const item of workspace.value?.requirements || []) {
    result[item.key] = item.missing;
  }
  return result;
});

const profileDone = computed(() => (requirementMap.value.profile || []).length === 0);
const cooperationDone = computed(() => (requirementMap.value.cooperation || []).length === 0);
const packageDone = computed(() => (requirementMap.value.package || []).length === 0);
const accountDone = computed(() => (requirementMap.value.account || []).length === 0);
const baseReady = computed(() => profileDone.value && cooperationDone.value);
const contentReady = computed(() => baseReady.value && packageDone.value);

const steps = computed<ProcessStep[]>(() => {
  const m = metrics.value;
  if (!merchant.value || !m) return [];

  return [
    {
      actionText: '编辑档案',
      desc: '确认商家是谁、在哪、做什么行业、谁负责对接。',
      key: 'profile',
      metricText: profileDone.value ? '资料已齐' : '资料未齐',
      missing: requirementMap.value.profile || [],
      path: '/users/list',
      status: profileDone.value ? 'done' : 'todo',
      title: '1. 商家基础档案',
    },
    {
      actionText: '补合作规则',
      desc: '确认合作方式、分成比例和当前合作阶段。',
      key: 'cooperation',
      metricText: cooperationDone.value ? '合作规则已齐' : '合作规则未齐',
      missing: requirementMap.value.cooperation || [],
      path: '/users/list',
      status: cooperationDone.value ? 'done' : profileDone.value ? 'todo' : 'locked',
      title: '2. 合作规则确认',
    },
    {
      actionText: '建套餐',
      desc: '没有套餐，就不知道视频到底推什么、利润空间够不够。',
      key: 'packages',
      metricText: `${m.enabledPackageCount}/${m.packageCount} 个启用套餐`,
      missing: requirementMap.value.package || [],
      path: '/plans/list',
      query: baseQuery.value,
      status: packageDone.value ? 'done' : baseReady.value ? 'todo' : 'locked',
      title: '3. 团购套餐建档',
    },
    {
      actionText: '记录授权',
      desc: '记录抖音账号、来客账号和代登/协作方式。',
      key: 'account',
      metricText: `${m.activeAccountAuthCount}/${m.accountAuthCount} 条已授权`,
      missing: requirementMap.value.account || [],
      path: '/users/devices',
      query: baseQuery.value,
      status: accountDone.value ? 'done' : profileDone.value ? 'todo' : 'locked',
      title: '4. 账号授权记录',
    },
    {
      actionText: '账号诊断',
      desc: '资料、合作、套餐和账号齐了，再让 Agent 做账号诊断。',
      key: 'diagnosis',
      metricText: `${m.completedDiagnosisCount}/${m.diagnosisCount} 次已完成`,
      missing: contentReady.value && accountDone.value ? [] : ['先补齐资料、套餐、账号授权'],
      path: '/users/account-diagnosis',
      query: baseQuery.value,
      status: m.completedDiagnosisCount > 0 ? 'done' : contentReady.value && accountDone.value ? 'todo' : 'locked',
      title: '5. 账号诊断',
    },
    {
      actionText: '找对标',
      desc: '找同城、同行、全国对标账号，拆内容结构和转化路径。',
      key: 'benchmark',
      metricText: `${m.analyzedBenchmarkCount}/${m.benchmarkCount} 个已分析`,
      missing: contentReady.value ? [] : ['先补齐基础资料和套餐'],
      path: '/lines/list',
      query: baseQuery.value,
      status: m.analyzedBenchmarkCount > 0 ? 'done' : contentReady.value ? 'todo' : 'locked',
      title: '6. 对标账号整理',
    },
    {
      actionText: '生成选题',
      desc: '基于商家、套餐、对标和热点生成选题池。',
      key: 'topic',
      metricText: `${m.acceptedTopicCount}/${m.topicCount} 条已采用`,
      missing: contentReady.value ? [] : ['先补齐基础资料和套餐'],
      path: '/content/notices',
      query: baseQuery.value,
      status: m.acceptedTopicCount > 0 ? 'done' : contentReady.value ? 'todo' : 'locked',
      title: '7. 生成选题池',
    },
    {
      actionText: '生成文案',
      desc: '选题采用后，再生成开头、主体、结尾引导。',
      key: 'copywriting',
      metricText: `${m.confirmedScriptCount}/${m.scriptCount} 条已确认`,
      missing: m.acceptedTopicCount > 0 ? [] : ['先采用至少 1 条选题'],
      path: '/content/quotes',
      query: baseQuery.value,
      status: m.confirmedScriptCount > 0 ? 'done' : m.acceptedTopicCount > 0 ? 'todo' : 'locked',
      title: '8. 生成文案脚本',
    },
    {
      actionText: '生成分镜',
      desc: '文案确认后，再拆成可执行镜头清单。',
      key: 'storyboard',
      metricText: `${m.confirmedStoryboardCount}/${m.storyboardCount} 条已确认`,
      missing: m.confirmedScriptCount > 0 ? [] : ['先确认至少 1 条文案'],
      path: '/content/discoveries',
      query: baseQuery.value,
      status: m.confirmedStoryboardCount > 0 ? 'done' : m.confirmedScriptCount > 0 ? 'todo' : 'locked',
      title: '9. 生成分镜脚本',
    },
    {
      actionText: '派拍摄',
      desc: '分镜确认后，分配拍摄/剪辑任务。',
      key: 'shooting',
      metricText: `${m.readyShootingTaskCount}/${m.shootingTaskCount} 个已剪辑/完成`,
      missing: m.confirmedStoryboardCount > 0 ? [] : ['先确认至少 1 条分镜'],
      path: '/logs/user',
      query: baseQuery.value,
      status: m.readyShootingTaskCount > 0 ? 'done' : m.confirmedStoryboardCount > 0 ? 'todo' : 'locked',
      title: '10. 拍摄任务',
    },
    {
      actionText: '排发布',
      desc: '视频素材完成后，再进入发布排期。',
      key: 'publish',
      metricText: `${m.publishedScheduleCount}/${m.scheduleCount} 条已发布`,
      missing: m.readyShootingTaskCount > 0 ? [] : ['先完成拍摄/剪辑'],
      path: '/content/payments',
      query: baseQuery.value,
      status: m.publishedScheduleCount > 0 ? 'done' : m.readyShootingTaskCount > 0 ? 'todo' : 'locked',
      title: '11. 发布排期',
    },
    {
      actionText: '做复盘',
      desc: '发布后看播放、互动、成交、核销，再优化下一轮。',
      key: 'review',
      metricText: `${m.reviewCount} 条复盘`,
      missing: m.publishedScheduleCount > 0 ? [] : ['先发布至少 1 条视频'],
      path: '/logs/admin',
      query: baseQuery.value,
      status: m.reviewCount > 0 ? 'done' : m.publishedScheduleCount > 0 ? 'todo' : 'locked',
      title: '12. 数据复盘',
    },
  ];
});

const currentStep = computed(
  () =>
    steps.value.find((item) => item.status === 'todo') ||
    steps.value.find((item) => item.status === 'locked') ||
    steps.value[steps.value.length - 1],
);
const doneCount = computed(() => steps.value.filter((item) => item.status === 'done').length);
const progressPercent = computed(() =>
  steps.value.length > 0 ? Math.round((doneCount.value / steps.value.length) * 100) : 0,
);

const statCards = computed(() => {
  const m = metrics.value;
  if (!m) return [];
  return [
    { label: '启用套餐', value: m.enabledPackageCount, total: m.packageCount },
    { label: '已授权账号', value: m.activeAccountAuthCount, total: m.accountAuthCount },
    { label: '已采用选题', value: m.acceptedTopicCount, total: m.topicCount },
    { label: '已确认文案', value: m.confirmedScriptCount, total: m.scriptCount },
    { label: '已确认分镜', value: m.confirmedStoryboardCount, total: m.storyboardCount },
    { label: '待发布素材', value: m.readyShootingTaskCount, total: m.shootingTaskCount },
    { label: '已发布视频', value: m.publishedScheduleCount, total: m.scheduleCount },
    { label: '复盘记录', value: m.reviewCount, total: m.reviewCount },
  ];
});

async function loadWorkspace() {
  if (!merchantId.value) return;
  loading.value = true;
  try {
    workspace.value = await getMerchantWorkspace(merchantId.value);
  } finally {
    loading.value = false;
  }
}

function handleStepAction(step: ProcessStep) {
  if (step.status === 'locked') {
    ElMessage.warning(`还不能做：${step.missing.join('、')}`);
    return;
  }
  router.push({
    path: step.path,
      query: step.query,
  });
}

onMounted(loadWorkspace);
</script>

<template>
  <div class="p-4">
    <el-skeleton v-if="loading" animated :rows="8" />

    <template v-else-if="merchant && metrics">
      <el-card class="mb-4">
        <div class="process-head">
          <div>
            <div class="breadcrumb-link" @click="router.push('/users/list')">
              ← 返回商家列表
            </div>
            <h2>{{ merchant.name }}</h2>
            <div class="merchant-meta">
              {{ merchant.industry || '未填行业' }} · {{ merchant.city || '未填城市' }} ·
              {{ merchant.cooperationType }} {{ merchant.commissionRate }}%
            </div>
          </div>
          <div class="stage-card">
            <div class="stage-label">当前应该处理</div>
            <div class="stage-title">{{ currentStep?.title }}</div>
            <div class="stage-desc">{{ currentStep?.desc }}</div>
            <el-button
              class="stage-button"
              type="primary"
              :disabled="!currentStep || currentStep.status !== 'todo'"
              @click="currentStep && handleStepAction(currentStep)"
            >
              {{ currentStep?.actionText || '暂无动作' }}
            </el-button>
          </div>
        </div>

        <el-progress :percentage="progressPercent" :stroke-width="12" />
      </el-card>

      <el-row :gutter="16" class="mb-4">
        <el-col v-for="item in statCards" :key="item.label" :span="6">
          <el-card shadow="never" class="stat-card">
            <div class="stat-label">{{ item.label }}</div>
            <div class="stat-value">{{ item.value }}<span>/ {{ item.total }}</span></div>
          </el-card>
        </el-col>
      </el-row>

      <el-row :gutter="16">
        <el-col :span="16">
          <el-card>
            <template #header>
              <div class="section-title">商家推进流程</div>
            </template>

            <div class="step-list">
              <div
                v-for="step in steps"
                :key="step.key"
                class="step-item"
                :class="step.status"
              >
                <div class="step-status">
                  <span v-if="step.status === 'done'">✓</span>
                  <span v-else-if="step.status === 'locked'">锁</span>
                  <span v-else>做</span>
                </div>
                <div class="step-main">
                  <div class="step-line">
                    <div class="step-title">{{ step.title }}</div>
                    <el-tag size="small" :type="step.status === 'done' ? 'success' : step.status === 'todo' ? 'primary' : 'info'">
                      {{ step.metricText }}
                    </el-tag>
                  </div>
                  <div class="step-desc">{{ step.desc }}</div>
                  <div v-if="step.missing.length > 0" class="step-missing">
                    缺：{{ step.missing.join('、') }}
                  </div>
                </div>
                <el-button
                  size="small"
                  :disabled="step.status === 'done'"
                  :type="step.status === 'todo' ? 'primary' : 'default'"
                  @click="handleStepAction(step)"
                >
                  {{ step.status === 'done' ? '已完成' : step.actionText }}
                </el-button>
              </div>
            </div>
          </el-card>
        </el-col>

        <el-col :span="8">
          <el-card class="mb-4">
            <template #header>
              <div class="section-title">前置条件检查</div>
            </template>
            <div
              v-for="item in workspace?.requirements || []"
              :key="item.key"
              class="requirement-item"
            >
              <el-tag :type="item.done ? 'success' : 'warning'">
                {{ item.done ? '已齐' : '缺' }}
              </el-tag>
              <div>
                <div class="requirement-title">{{ item.title }}</div>
                <div class="requirement-missing">
                  {{ item.done ? '可以进入后续流程' : item.missing.join('、') }}
                </div>
              </div>
            </div>
          </el-card>

          <el-card>
            <template #header>
              <div class="section-title">商家基础信息</div>
            </template>
            <div class="info-row">
              <span>联系人</span>
              <strong>{{ merchant.contactName || '-' }}</strong>
            </div>
            <div class="info-row">
              <span>联系电话</span>
              <strong>{{ merchant.contactPhone || '-' }}</strong>
            </div>
            <div class="info-row">
              <span>抖音账号</span>
              <strong>{{ merchant.douyinAccount || '-' }}</strong>
            </div>
            <div class="info-row">
              <span>抖音来客</span>
              <strong>{{ merchant.douyinLaikeAccount || '-' }}</strong>
            </div>
            <div class="info-row">
              <span>地址</span>
              <strong>{{ merchant.address || '-' }}</strong>
            </div>
          </el-card>
        </el-col>
      </el-row>
    </template>

    <el-empty v-else description="商家不存在或参数错误">
      <el-button type="primary" @click="router.push('/users/list')">
        返回商家列表
      </el-button>
    </el-empty>
  </div>
</template>

<style scoped>
.process-head {
  display: flex;
  justify-content: space-between;
  gap: 24px;
  margin-bottom: 20px;
}

.breadcrumb-link {
  margin-bottom: 8px;
  color: #2563eb;
  cursor: pointer;
  font-size: 13px;
}

h2 {
  margin: 0;
  color: #0f172a;
  font-size: 24px;
}

.merchant-meta {
  margin-top: 8px;
  color: #64748b;
}

.stage-card {
  width: 380px;
  border-radius: 16px;
  background: #eef4ff;
  padding: 16px;
}

.stage-label {
  color: #64748b;
  font-size: 12px;
}

.stage-title {
  margin-top: 8px;
  color: #1d4ed8;
  font-size: 18px;
  font-weight: 700;
}

.stage-desc {
  margin-top: 8px;
  color: #475569;
  font-size: 13px;
  line-height: 1.6;
}

.stage-button {
  margin-top: 12px;
}

.stat-card {
  margin-bottom: 16px;
  border-radius: 14px;
}

.stat-label {
  color: #64748b;
  font-size: 13px;
}

.stat-value {
  margin-top: 8px;
  color: #0f172a;
  font-size: 24px;
  font-weight: 800;
}

.stat-value span {
  margin-left: 4px;
  color: #94a3b8;
  font-size: 13px;
  font-weight: 500;
}

.section-title {
  color: #0f172a;
  font-weight: 700;
}

.step-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.step-item {
  display: flex;
  align-items: center;
  gap: 14px;
  border: 1px solid #e2e8f0;
  border-radius: 14px;
  padding: 14px;
}

.step-item.done {
  background: #f0fdf4;
}

.step-item.todo {
  border-color: #93c5fd;
  background: #eff6ff;
}

.step-item.locked {
  background: #f8fafc;
  opacity: 0.78;
}

.step-status {
  display: flex;
  width: 36px;
  height: 36px;
  align-items: center;
  justify-content: center;
  border-radius: 50%;
  background: #dbeafe;
  color: #1d4ed8;
  flex-shrink: 0;
  font-size: 13px;
  font-weight: 700;
}

.step-item.done .step-status {
  background: #dcfce7;
  color: #16a34a;
}

.step-item.locked .step-status {
  background: #e2e8f0;
  color: #64748b;
}

.step-main {
  flex: 1;
}

.step-line {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
}

.step-title {
  color: #0f172a;
  font-weight: 700;
}

.step-desc,
.step-missing {
  margin-top: 4px;
  color: #64748b;
  font-size: 13px;
  line-height: 1.5;
}

.step-missing {
  color: #b45309;
}

.requirement-item {
  display: flex;
  gap: 10px;
  border-bottom: 1px solid #f1f5f9;
  padding: 12px 0;
}

.requirement-title {
  color: #0f172a;
  font-weight: 700;
}

.requirement-missing {
  margin-top: 4px;
  color: #64748b;
  font-size: 13px;
}

.info-row {
  display: flex;
  justify-content: space-between;
  gap: 12px;
  border-bottom: 1px solid #f1f5f9;
  padding: 12px 0;
}

.info-row span {
  color: #64748b;
}

.info-row strong {
  color: #0f172a;
  text-align: right;
}
</style>
