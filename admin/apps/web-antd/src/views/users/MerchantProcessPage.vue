<script setup lang="ts">
import { computed, onMounted, ref } from 'vue';
import { useRoute, useRouter } from 'vue-router';

import { ElMessage } from 'element-plus';

import { getMerchant, type Merchant } from '#/api/admin';

defineOptions({ name: 'MerchantProcessPage' });

type StepStatus = 'done' | 'locked' | 'todo';

interface ProcessStep {
  actionText: string;
  desc: string;
  key: string;
  missing: string[];
  status: StepStatus;
  title: string;
}

const route = useRoute();
const router = useRouter();
const loading = ref(false);
const merchant = ref<Merchant>();

const merchantId = computed(() => Number(route.params.id || 0));

const basicMissing = computed(() => {
  const row = merchant.value;
  if (!row) return ['商家档案'];
  const missing: string[] = [];
  if (!row.name) missing.push('商家名称');
  if (!row.industry) missing.push('行业');
  if (!row.city) missing.push('城市');
  if (!row.contactName) missing.push('联系人');
  if (!row.contactPhone) missing.push('联系电话');
  return missing;
});

const cooperationMissing = computed(() => {
  const row = merchant.value;
  if (!row) return ['合作规则'];
  const missing: string[] = [];
  if (!row.cooperationType) missing.push('合作方式');
  if (!row.commissionRate || row.commissionRate <= 0) missing.push('分成比例');
  return missing;
});

const accountMissing = computed(() => {
  const row = merchant.value;
  if (!row) return ['账号授权'];
  const missing: string[] = [];
  if (!row.douyinAccount) missing.push('抖音账号');
  if (!row.douyinLaikeAccount) missing.push('抖音来客账号');
  return missing;
});

const profileDone = computed(() => basicMissing.value.length === 0);
const cooperationDone = computed(() => cooperationMissing.value.length === 0);
const accountDone = computed(() => accountMissing.value.length === 0);

const steps = computed<ProcessStep[]>(() => {
  const canDoPackage = profileDone.value && cooperationDone.value;
  const canDoAccount = profileDone.value;
  const canDoDiagnosis = profileDone.value && cooperationDone.value && accountDone.value;

  return [
    {
      actionText: '回列表编辑档案',
      desc: '确认商家是谁、在哪、做什么行业、谁负责对接。',
      key: 'basic',
      missing: basicMissing.value,
      status: profileDone.value ? 'done' : 'todo',
      title: '1. 商家基础档案',
    },
    {
      actionText: '回列表补合作规则',
      desc: '确认合作方式、分成比例和当前合作阶段。',
      key: 'cooperation',
      missing: cooperationMissing.value,
      status: cooperationDone.value ? 'done' : profileDone.value ? 'todo' : 'locked',
      title: '2. 合作规则确认',
    },
    {
      actionText: '去建团购套餐',
      desc: '没有套餐，就不知道视频到底推什么、利润空间够不够。',
      key: 'packages',
      missing: ['团购套餐模块待开发'],
      status: canDoPackage ? 'todo' : 'locked',
      title: '3. 团购套餐建档',
    },
    {
      actionText: '去记录账号授权',
      desc: '记录抖音账号、抖音来客账号和协作方式，先不做扫码登录。',
      key: 'account',
      missing: accountMissing.value,
      status: accountDone.value ? 'done' : canDoAccount ? 'todo' : 'locked',
      title: '4. 账号授权记录',
    },
    {
      actionText: '开始账号诊断',
      desc: '资料、合作、账号都齐了，再让 Agent 做账号诊断。',
      key: 'diagnosis',
      missing: canDoDiagnosis ? ['账号诊断 Agent 尚未接入'] : ['先补齐前置资料'],
      status: canDoDiagnosis ? 'todo' : 'locked',
      title: '5. 账号诊断',
    },
    {
      actionText: '整理对标账号',
      desc: '诊断之后再找同城、同行、全国对标账号。',
      key: 'benchmark',
      missing: ['等待账号诊断结果'],
      status: 'locked',
      title: '6. 对标账号整理',
    },
    {
      actionText: '生成选题池',
      desc: '只有套餐、诊断和对标齐了，选题才有依据。',
      key: 'topic',
      missing: ['等待套餐、诊断、对标'],
      status: 'locked',
      title: '7. 生成选题池',
    },
    {
      actionText: '生成文案',
      desc: '选题确认后，再生成开头、主体、结尾引导。',
      key: 'copywriting',
      missing: ['等待选题确认'],
      status: 'locked',
      title: '8. 生成文案脚本',
    },
    {
      actionText: '生成分镜',
      desc: '文案确认后，再拆成可拍摄镜头。',
      key: 'storyboard',
      missing: ['等待文案确认'],
      status: 'locked',
      title: '9. 生成分镜脚本',
    },
    {
      actionText: '安排拍摄',
      desc: '分镜确认后，才分配拍摄任务。',
      key: 'shooting',
      missing: ['等待分镜确认'],
      status: 'locked',
      title: '10. 拍摄任务',
    },
    {
      actionText: '安排发布',
      desc: '视频素材完成后，再进入发布排期。',
      key: 'publish',
      missing: ['等待拍摄剪辑完成'],
      status: 'locked',
      title: '11. 发布排期',
    },
    {
      actionText: '发布后复盘',
      desc: '发布后看播放、互动、成交、核销，再优化下一轮。',
      key: 'review',
      missing: ['等待视频发布和数据回收'],
      status: 'locked',
      title: '12. 数据复盘',
    },
  ];
});

const currentStep = computed(
  () =>
    steps.value.find((item) => item.status !== 'done') ||
    steps.value[steps.value.length - 1],
);
const doneCount = computed(() => steps.value.filter((item) => item.status === 'done').length);
const progressPercent = computed(() => Math.round((doneCount.value / steps.value.length) * 100));

async function loadMerchant() {
  if (!merchantId.value) return;
  loading.value = true;
  try {
    merchant.value = await getMerchant(merchantId.value);
  } finally {
    loading.value = false;
  }
}

function handleStepAction(step: ProcessStep) {
  if (step.status === 'locked') {
    ElMessage.warning(`还不能做：${step.missing.join('、')}`);
    return;
  }

  if (step.key === 'packages') {
    router.push('/plans/list');
    return;
  }
  if (step.key === 'account') {
    router.push('/users/devices');
    return;
  }
  if (step.key === 'basic' || step.key === 'cooperation') {
    router.push('/users/list');
    return;
  }

  ElMessage.info('这个步骤的具体页面和 Agent 接口后续开发。');
}

onMounted(loadMerchant);
</script>

<template>
  <div class="p-4">
    <el-skeleton v-if="loading" animated :rows="8" />

    <template v-else-if="merchant">
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
            <div class="stage-label">当前系统建议</div>
            <div class="stage-title">{{ currentStep?.title }}</div>
            <div class="stage-desc">{{ currentStep?.desc }}</div>
          </div>
        </div>

        <el-progress :percentage="progressPercent" :stroke-width="12" />
      </el-card>

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
                  <div class="step-title">{{ step.title }}</div>
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
              <div class="section-title">现在不要做什么</div>
            </template>
            <el-alert type="warning" :closable="false" show-icon>
              资料没齐之前，不应该直接做选题、文案、分镜。否则 Agent 没上下文，输出就会变成泛泛而谈。
            </el-alert>
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
  width: 360px;
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
