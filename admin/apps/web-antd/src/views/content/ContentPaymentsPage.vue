<script setup lang="ts">
import { ElMessage } from 'element-plus';

import { publishSchedule } from '#/mock/content-ops';

defineOptions({ name: 'ContentPaymentsPage' });

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
            <div class="page-title">发布排期</div>
            <div class="page-desc">这里先看每条视频什么时候发、谁负责、当前卡在哪一步。</div>
          </div>
          <div class="page-actions">
            <el-button @click="demoAction('检查发布素材')">检查素材</el-button>
            <el-button type="primary" @click="demoAction('新增发布排期')">新增排期</el-button>
          </div>
        </div>
      </template>

      <el-table :data="publishSchedule" border stripe>
        <el-table-column prop="merchant" label="商家" min-width="160" />
        <el-table-column prop="title" label="视频标题" min-width="220" />
        <el-table-column prop="publishTime" label="发布时间" width="170" />
        <el-table-column prop="owner" label="负责人" width="120" />
        <el-table-column prop="status" label="状态" width="110">
          <template #default="{ row }">
            <el-tag :type="row.status === '待发布' ? 'warning' : 'success'">{{ row.status }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="300" fixed="right">
          <template #default="{ row }">
            <el-button size="small" @click="demoAction(`老板确认：${row.title}`)">老板确认</el-button>
            <el-button size="small" type="primary" @click="demoAction(`发布：${row.title}`)">发布</el-button>
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

.page-actions {
  display: flex;
  flex-shrink: 0;
  gap: 10px;
}
</style>
