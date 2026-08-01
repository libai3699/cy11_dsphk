<script setup lang="ts">
import { ElMessage } from 'element-plus';
import { useRouter } from 'vue-router';

import { merchants } from '#/mock/content-ops';

defineOptions({ name: 'UserListPage' });

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
            <div class="page-title">商家列表</div>
            <div class="page-desc">这里先看你签下来的商家、行业、套餐、近 7 天成交和下一动作。</div>
          </div>
          <div class="page-actions">
            <el-button @click="demoAction('导入商家资料')">导入资料</el-button>
            <el-button type="primary" @click="demoAction('新增商家')">新增商家</el-button>
          </div>
        </div>
      </template>

      <el-alert class="mb-4" type="info" :closable="false" show-icon>
        正式版这里会接入商家资料、合同、账号诊断结果和负责人分配。现在先只保留你判断业务节奏最关键的字段。
      </el-alert>

      <el-table :data="merchants" border stripe>
        <el-table-column prop="merchant" label="商家" min-width="180" />
        <el-table-column prop="industry" label="行业" width="100" />
        <el-table-column prop="city" label="城市" width="90" />
        <el-table-column prop="stage" label="当前阶段" width="130">
          <template #default="{ row }">
            <el-tag :type="row.stage.includes('冷启动') ? 'warning' : 'success'">{{ row.stage }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="packageName" label="主推套餐" min-width="180" />
        <el-table-column prop="recentGmv" label="近 7 天成交" width="130" />
        <el-table-column prop="commissionRate" label="分成比" width="90" />
        <el-table-column prop="nextAction" label="下一动作" min-width="240" />
        <el-table-column label="操作" width="300" fixed="right">
          <template #default="{ row }">
            <el-button size="small" @click="demoAction(`生成 ${row.merchant} 账号诊断`)">账号诊断</el-button>
            <el-button size="small" @click="router.push('/lines/list')">找对标</el-button>
            <el-button size="small" type="primary" @click="router.push('/content/notices')">做选题</el-button>
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
  margin-top: 6px;
  color: #64748b;
  font-size: 13px;
}

.page-actions {
  display: flex;
  gap: 10px;
}
</style>
