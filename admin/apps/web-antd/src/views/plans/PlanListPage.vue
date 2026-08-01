<script setup lang="ts">
import { ElMessage } from 'element-plus';
import { useRouter } from 'vue-router';

import { packages } from '#/mock/content-ops';

defineOptions({ name: 'PlanListPage' });

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
            <div class="page-title">团购套餐</div>
            <div class="page-desc">这里的重点不是“卖得便宜”，而是先算明白商家是否还有利润空间。</div>
          </div>
          <div class="page-actions">
            <el-button @click="demoAction('AI 检查利润空间')">利润检查</el-button>
            <el-button type="primary" @click="demoAction('新增团购套餐')">新增套餐</el-button>
          </div>
        </div>
      </template>

      <el-table :data="packages" border stripe>
        <el-table-column prop="merchant" label="商家" min-width="160" />
        <el-table-column prop="packageName" label="套餐名称" min-width="180" />
        <el-table-column prop="sellingPrice" label="售价" width="100" />
        <el-table-column prop="cost" label="成本" width="100" />
        <el-table-column prop="margin" label="毛利率" width="100" />
        <el-table-column prop="trafficLabel" label="投放定位" min-width="160" />
        <el-table-column prop="profitGuard" label="利润保护建议" min-width="200" />
        <el-table-column label="操作" width="260" fixed="right">
          <template #default="{ row }">
            <el-button size="small" @click="demoAction(`优化 ${row.packageName}`)">优化套餐</el-button>
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
