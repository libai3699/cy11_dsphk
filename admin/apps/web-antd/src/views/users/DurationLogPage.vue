<script setup lang="ts">
import { ElMessage } from 'element-plus';

import { followUpLogs } from '#/mock/content-ops';

defineOptions({ name: 'DurationLogPage' });

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
            <div class="page-title">跟进记录</div>
            <div class="page-desc">前线谈商家时最有价值的信息，是对方卡在哪里、下一步怎么推进。</div>
          </div>
          <el-button type="primary" @click="demoAction('新增跟进记录')">新增跟进</el-button>
        </div>
      </template>

      <el-table :data="followUpLogs" border stripe>
        <el-table-column prop="merchant" label="商家" min-width="160" />
        <el-table-column prop="stage" label="阶段" width="110">
          <template #default="{ row }">
            <el-tag :type="row.stage === '已签约' ? 'success' : 'warning'">{{ row.stage }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="latestTalk" label="最近沟通" min-width="220" />
        <el-table-column prop="objection" label="关键异议" min-width="200" />
        <el-table-column prop="nextStep" label="下一步" min-width="200" />
        <el-table-column prop="owner" label="负责人" width="110" />
        <el-table-column label="操作" width="260" fixed="right">
          <template #default="{ row }">
            <el-button size="small" @click="demoAction(`生成 ${row.merchant} 跟进话术`)">生成话术</el-button>
            <el-button size="small" type="primary" @click="demoAction(`推进 ${row.merchant} 下一步`)">推进</el-button>
          </template>
        </el-table-column>
      </el-table>
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
</style>
