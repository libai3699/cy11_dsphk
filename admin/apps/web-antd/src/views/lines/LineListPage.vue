<script setup lang="ts">
import { ElMessage } from 'element-plus';
import { useRouter } from 'vue-router';

import { benchmarks } from '#/mock/content-ops';

defineOptions({ name: 'LineListPage' });

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
            <div class="page-title">对标账号库</div>
            <div class="page-desc">运营做选题前，先知道抄谁、抄哪一段、哪些不能抄。</div>
          </div>
          <div class="page-actions">
            <el-button @click="demoAction('AI 搜索同城对标账号')">搜索对标</el-button>
            <el-button type="primary" @click="demoAction('新增对标账号')">新增对标</el-button>
          </div>
        </div>
      </template>

      <el-table :data="benchmarks" border stripe>
        <el-table-column prop="account" label="对标账号" min-width="180" />
        <el-table-column prop="city" label="城市" width="90" />
        <el-table-column prop="lane" label="赛道" width="130" />
        <el-table-column prop="latestHit" label="最近爆款" min-width="220" />
        <el-table-column prop="takeaway" label="可抄点" min-width="220" />
        <el-table-column prop="risk" label="风险提醒" min-width="180" />
        <el-table-column label="操作" width="280" fixed="right">
          <template #default="{ row }">
            <el-button size="small" @click="demoAction(`分析 ${row.account}`)">分析爆款</el-button>
            <el-button size="small" type="primary" @click="router.push('/content/notices')">生成选题</el-button>
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
