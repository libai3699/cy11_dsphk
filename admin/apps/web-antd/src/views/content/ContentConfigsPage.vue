<script setup lang="ts">
import { ElMessage } from 'element-plus';

import { roleCards, systemRules } from '#/mock/content-ops';

defineOptions({ name: 'ContentConfigsPage' });

function demoAction(text: string) {
  ElMessage.success(`${text}：演示动作已触发`);
}
</script>

<template>
  <div class="p-4">
    <el-card>
      <template #header>
        <div class="page-head">
          <div>
            <div class="page-title">系统规则</div>
            <div class="page-desc">这页对应的是后续后台里最不该交给 AI 自主决定的业务规则。</div>
          </div>
          <div class="page-actions">
            <el-button @click="demoAction('编辑审批节点')">编辑审批</el-button>
            <el-button type="primary" @click="demoAction('保存系统规则')">保存规则</el-button>
          </div>
        </div>
      </template>

      <el-descriptions :column="1" border class="mb-4">
        <el-descriptions-item
          v-for="item in systemRules"
          :key="item.label"
          :label="item.label"
        >
          {{ item.value }}
        </el-descriptions-item>
      </el-descriptions>

      <el-row :gutter="16">
        <el-col v-for="item in roleCards" :key="item.role" :span="8">
          <el-card shadow="never" class="role-card">
            <div class="role-name">{{ item.role }}</div>
            <div class="role-focus">{{ item.focus }}</div>
            <el-button class="role-button" size="small" @click="demoAction(`配置 ${item.role} 权限`)">配置权限</el-button>
          </el-card>
        </el-col>
      </el-row>
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

.role-card {
  min-height: 140px;
  border-radius: 16px;
}

.role-name {
  color: #0f172a;
  font-size: 16px;
  font-weight: 700;
}

.role-focus {
  margin-top: 10px;
  color: #475569;
  line-height: 1.7;
}

.role-button {
  margin-top: 14px;
}
</style>
