<script setup lang="ts">
import { ElMessage } from 'element-plus';
import { useRouter } from 'vue-router';

import { reviewRows } from '#/mock/content-ops';

defineOptions({ name: 'LogAdminPage' });

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
            <div class="page-title">数据复盘</div>
            <div class="page-desc">每条视频发出去后，最终都要回到“有没有成交、下一条怎么改”。</div>
          </div>
          <div class="page-actions">
            <el-button @click="demoAction('同步视频数据')">同步数据</el-button>
            <el-button type="primary" @click="demoAction('AI 生成周复盘')">生成周复盘</el-button>
          </div>
        </div>
      </template>

      <el-table :data="reviewRows" border stripe>
        <el-table-column prop="merchant" label="商家" min-width="160" />
        <el-table-column prop="video" label="视频" min-width="220" />
        <el-table-column prop="hook" label="开场钩子" min-width="220" />
        <el-table-column prop="views" label="播放" width="100" />
        <el-table-column prop="deals" label="成交单数" width="110" />
        <el-table-column prop="ordersAmount" label="成交额" width="110" />
        <el-table-column prop="nextAction" label="下一步动作" min-width="220" />
        <el-table-column label="操作" width="300" fixed="right">
          <template #default="{ row }">
            <el-button size="small" @click="demoAction(`复盘 ${row.video}`)">单条复盘</el-button>
            <el-button size="small" type="primary" @click="router.push('/content/notices')">生成下轮选题</el-button>
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
