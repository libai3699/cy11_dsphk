<script setup lang="ts">
import { ElMessage } from 'element-plus';

import { accountAccesses } from '#/mock/content-ops';

defineOptions({ name: 'UserDevicesPage' });

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
            <div class="page-title">账号授权</div>
            <div class="page-desc">这里只保留你真正会用到的远程登录和协作信息，不碰任何自动登录。</div>
          </div>
          <div class="page-actions">
            <el-button @click="demoAction('记录验证码代登')">记录代登</el-button>
            <el-button type="primary" @click="demoAction('新增协作账号')">新增授权</el-button>
          </div>
        </div>
      </template>

      <el-alert class="mb-4" type="warning" :closable="false" show-icon>
        后续正式版建议只保留“协作授权、验证码代登、人工确认发布”这几种安全流程，不要把密码托管给系统。
      </el-alert>

      <el-table :data="accountAccesses" border stripe>
        <el-table-column prop="merchant" label="商家" min-width="160" />
        <el-table-column prop="method" label="授权方式" min-width="160" />
        <el-table-column prop="account" label="账号信息" min-width="180" />
        <el-table-column prop="status" label="状态" width="110">
          <template #default="{ row }">
            <el-tag :type="row.status === '已授权' ? 'success' : 'warning'">{{ row.status }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="lastLogin" label="最近登录" width="170" />
        <el-table-column prop="note" label="备注" min-width="220" />
        <el-table-column label="操作" width="260" fixed="right">
          <template #default="{ row }">
            <el-button size="small" @click="demoAction(`刷新 ${row.merchant} 授权状态`)">刷新状态</el-button>
            <el-button size="small" type="primary" @click="demoAction(`进入 ${row.merchant} 账号诊断`)">账号诊断</el-button>
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
