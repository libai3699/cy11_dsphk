<script setup lang="ts">
import { ElMessage } from 'element-plus';
import { useRouter } from 'vue-router';

import { storyboards } from '#/mock/content-ops';

defineOptions({ name: 'ContentDiscoveriesPage' });

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
            <div class="page-title">分镜脚本</div>
            <div class="page-desc">拍摄的人最需要的是“到哪拍、拍什么、镜头里说什么”。</div>
          </div>
          <div class="page-actions">
            <el-button @click="demoAction('AI 拆分镜头')">AI 拆分镜</el-button>
            <el-button type="primary" @click="router.push('/logs/user')">派拍摄任务</el-button>
          </div>
        </div>
      </template>

      <el-table :data="storyboards" border stripe>
        <el-table-column prop="scriptTitle" label="脚本标题" min-width="200" />
        <el-table-column prop="scene" label="场景" min-width="160" />
        <el-table-column prop="location" label="时间 / 地点" min-width="200" />
        <el-table-column prop="lens" label="镜头设计" min-width="220" />
        <el-table-column prop="dialogue" label="台词提示" min-width="220" />
        <el-table-column label="操作" width="280" fixed="right">
          <template #default="{ row }">
            <el-button size="small" @click="demoAction(`导出拍摄清单：${row.scriptTitle}`)">导出清单</el-button>
            <el-button size="small" type="primary" @click="router.push('/logs/user')">派任务</el-button>
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
