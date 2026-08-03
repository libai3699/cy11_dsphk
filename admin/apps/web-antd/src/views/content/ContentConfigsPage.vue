<script setup lang="ts">
import { onMounted, ref } from 'vue';
import { ElMessage } from 'element-plus';

import {
  getAgentConfigs,
  type AgentConfig,
} from '#/api/admin/content-production';

defineOptions({ name: 'ContentConfigsPage' });

const loading = ref(false);
const configs = ref<AgentConfig[]>([]);

const agentNames: Record<string, string> = {
  account_diagnosis: '账号诊断 Agent',
  benchmark: '对标分析 Agent',
  copywriting: '文案脚本 Agent',
  hotspot_topic: '找爆款 Agent',
  review: '数据复盘 Agent',
  storyboard: '分镜脚本 Agent',
};

const envRows = [
  { key: 'STEP_API_KEY', value: '公共账号 Key，所有 Agent 共用；后端不回传明文' },
  { key: 'STEP_BASE_URL', value: '公共 StepFun Base URL' },
  { key: 'STEP_MODEL', value: '公共默认模型' },
  { key: 'STEP_TIMEOUT_SECONDS', value: '公共默认超时秒数' },
  { key: 'AGENT_{NAME}_MODEL', value: '单个 Agent 覆盖模型，例如 AGENT_COPYWRITING_MODEL' },
  { key: 'AGENT_{NAME}_TIMEOUT_SECONDS', value: '单个 Agent 覆盖超时' },
  { key: 'AGENT_{NAME}_BASE_URL', value: '单个 Agent 覆盖 Base URL，一般不需要' },
  { key: 'AGENT_{NAME}_ENABLED', value: '单个 Agent 开关：true/false' },
];

async function loadConfigs() {
  loading.value = true;
  try {
    configs.value = await getAgentConfigs();
  } finally {
    loading.value = false;
  }
}

function saveTip() {
  ElMessage.info('当前配置从 server/.env 读取。修改 .env 后重启后端生效。');
}

onMounted(loadConfigs);
</script>

<template>
  <div class="p-4">
    <el-card>
      <template #header>
        <div class="page-head">
          <div>
            <div class="page-title">系统规则</div>
            <div class="page-desc">这里看 Agent 配置和生产流程边界。Key 只放后端环境变量，不在前端暴露。</div>
          </div>
          <div class="page-actions">
            <el-button @click="loadConfigs">刷新</el-button>
            <el-button type="primary" @click="saveTip">如何修改配置</el-button>
          </div>
        </div>
      </template>

      <el-alert class="mb-4" type="info" :closable="false" show-icon>
        当前策略：所有 Agent 共用同一个 StepFun 账号 Key；每个 Agent 可以单独覆盖模型、超时、Base URL 和启停状态。
      </el-alert>

      <el-table v-loading="loading" :data="configs" border stripe class="mb-4">
        <el-table-column label="Agent" min-width="170">
          <template #default="{ row }">{{ agentNames[row.agent] || row.agent }}</template>
        </el-table-column>
        <el-table-column prop="provider" label="Provider" width="110" />
        <el-table-column prop="model" label="模型" min-width="160" />
        <el-table-column prop="baseUrl" label="Base URL" min-width="260" show-overflow-tooltip />
        <el-table-column prop="timeoutSeconds" label="超时秒数" width="110" />
        <el-table-column label="Key" width="120">
          <template #default="{ row }">
            <el-tag :type="row.keyConfigured ? 'success' : 'danger'">
              {{ row.keyConfigured ? '已配置' : '未配置' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="启用" width="100">
          <template #default="{ row }">
            <el-tag :type="row.enabled ? 'success' : 'info'">
              {{ row.enabled ? '启用' : '关闭' }}
            </el-tag>
          </template>
        </el-table-column>
      </el-table>

      <el-descriptions :column="1" border>
        <el-descriptions-item v-for="item in envRows" :key="item.key" :label="item.key">
          {{ item.value }}
        </el-descriptions-item>
      </el-descriptions>
    </el-card>
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
</style>
