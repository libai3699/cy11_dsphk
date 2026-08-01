<script setup lang="ts">
import { ElMessage } from 'element-plus';
import { useRouter } from 'vue-router';

import { shootTasks } from '#/mock/content-ops';

defineOptions({ name: 'LogUserPage' });

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
            <div class="page-title">拍摄任务</div>
            <div class="page-desc">运营定完分镜后，拍摄和剪辑就按任务流往下走。</div>
          </div>
          <div class="page-actions">
            <el-button @click="demoAction('导出今日拍摄单')">导出拍摄单</el-button>
            <el-button type="primary" @click="demoAction('新增拍摄任务')">新增任务</el-button>
          </div>
        </div>
      </template>

      <el-table :data="shootTasks" border stripe>
        <el-table-column prop="taskName" label="任务" min-width="220" />
        <el-table-column prop="merchant" label="商家" min-width="160" />
        <el-table-column prop="shotList" label="镜头清单" min-width="260" />
        <el-table-column prop="assignee" label="执行人" width="120" />
        <el-table-column prop="deadline" label="截止时间" width="120" />
        <el-table-column prop="status" label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="row.status === '已完成' ? 'success' : 'warning'">{{ row.status }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="300" fixed="right">
          <template #default="{ row }">
            <el-button size="small" @click="demoAction(`上传素材：${row.taskName}`)">上传素材</el-button>
            <el-button size="small" type="primary" @click="router.push('/content/payments')">排发布</el-button>
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
