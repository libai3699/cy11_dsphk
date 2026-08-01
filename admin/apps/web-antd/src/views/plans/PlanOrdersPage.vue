<script setup lang="ts">
import { ElMessage } from 'element-plus';

import { settlementOrders } from '#/mock/content-ops';

defineOptions({ name: 'PlanOrdersPage' });

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
            <div class="page-title">分成订单</div>
            <div class="page-desc">正式版这里会接“视频 -> 团购订单 -> 核销额 -> 分成”的归因链路。</div>
          </div>
          <div class="page-actions">
            <el-button @click="demoAction('导入抖音来客核销单')">导入核销</el-button>
            <el-button type="primary" @click="demoAction('生成本周结算单')">生成结算单</el-button>
          </div>
        </div>
      </template>

      <el-table :data="settlementOrders" border stripe>
        <el-table-column prop="merchant" label="商家" min-width="160" />
        <el-table-column prop="sourceVideo" label="来源视频" min-width="220" />
        <el-table-column prop="orderWindow" label="统计周期" width="150" />
        <el-table-column prop="redeemedAmount" label="核销额" width="110" />
        <el-table-column prop="commissionRate" label="分成比" width="90" />
        <el-table-column prop="commission" label="应分成" width="110" />
        <el-table-column prop="status" label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="row.status === '已确认' ? 'success' : 'warning'">{{ row.status }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="260" fixed="right">
          <template #default="{ row }">
            <el-button size="small" @click="demoAction(`核对 ${row.merchant} 分成`)">核对</el-button>
            <el-button size="small" type="primary" @click="demoAction(`确认 ${row.merchant} 结算`)">确认结算</el-button>
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
