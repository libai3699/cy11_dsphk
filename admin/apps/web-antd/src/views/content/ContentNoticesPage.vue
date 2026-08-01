<script setup lang="ts">
import { ElMessage } from 'element-plus';
import { useRouter } from 'vue-router';

import { topics } from '#/mock/content-ops';

defineOptions({ name: 'ContentNoticesPage' });

const router = useRouter();

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
            <div class="page-title">选题中心</div>
            <div class="page-desc">选题不是一堆灵感，而是“商家目标 + 热点来源 + 发布时间”的可执行池子。</div>
          </div>
          <div class="page-actions">
            <el-button @click="demoAction('同步同城热点')">同步热点</el-button>
            <el-button type="primary" @click="demoAction('AI 生成选题池')">生成选题池</el-button>
          </div>
        </div>
      </template>

      <el-table :data="topics" border stripe>
        <el-table-column prop="merchant" label="商家" min-width="160" />
        <el-table-column prop="topic" label="选题" min-width="180" />
        <el-table-column prop="hook" label="开场钩子" min-width="260" />
        <el-table-column prop="source" label="灵感来源" min-width="180" />
        <el-table-column prop="publishWindow" label="推荐发布时间" width="160" />
        <el-table-column prop="status" label="状态" width="120">
          <template #default="{ row }">
            <el-tag :type="row.status === '待拍摄' ? 'warning' : 'success'">{{ row.status }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="300" fixed="right">
          <template #default="{ row }">
            <el-button size="small" @click="demoAction(`采用选题：${row.topic}`)">采用</el-button>
            <el-button size="small" type="primary" @click="router.push('/content/quotes')">生成文案</el-button>
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
