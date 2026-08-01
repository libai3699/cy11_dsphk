<script setup lang="ts">
import { computed, ref, watch } from 'vue';
import { ElMessage } from 'element-plus';
import { useRoute, useRouter } from 'vue-router';

import { storyboards } from '#/mock/content-ops';

defineOptions({ name: 'ContentDiscoveriesPage' });

const route = useRoute();
const router = useRouter();
const keyword = ref('');

const merchantName = computed(() => String(route.query.merchantName || ''));
const merchantId = computed(() => String(route.query.merchantId || ''));

const filteredStoryboards = computed(() => {
  const searchText = keyword.value.trim();
  if (!searchText) return storyboards;
  return storyboards.filter((item) =>
    [
      item.merchant,
      item.scriptTitle,
      item.scene,
      item.location,
      item.lens,
      item.dialogue,
    ].some((value) => value.includes(searchText)),
  );
});

watch(
  () => route.query,
  () => {
    keyword.value = merchantName.value;
  },
  { immediate: true },
);

function demoAction(text: string) {
  ElMessage.success(`${text}：演示动作已触发`);
}

function resetSearch() {
  keyword.value = '';
}

function goShootTasks(rowMerchant?: string) {
  const targetMerchant = rowMerchant || merchantName.value;
  router.push({
    path: '/logs/user',
    query: targetMerchant
      ? {
          merchantId: merchantId.value,
          merchantName: targetMerchant,
        }
      : {},
  });
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
            <el-button type="primary" @click="goShootTasks()">派拍摄任务</el-button>
          </div>
        </div>
      </template>

      <el-alert
        v-if="merchantName"
        class="mb-4"
        type="info"
        :closable="false"
        show-icon
      >
        当前处理商家：{{ merchantName }}。已自动带入筛选条件，商家ID：{{ merchantId || '-' }}。
      </el-alert>

      <el-form :inline="true" class="search-form">
        <el-form-item label="关键词">
          <el-input
            v-model="keyword"
            clearable
            placeholder="商家 / 脚本 / 场景 / 地点 / 镜头"
            style="width: 320px"
          />
        </el-form-item>
        <el-form-item>
          <el-button @click="resetSearch">显示全部</el-button>
        </el-form-item>
      </el-form>

      <el-table :data="filteredStoryboards" border stripe>
        <el-table-column prop="merchant" label="商家" min-width="160" />
        <el-table-column prop="scriptTitle" label="脚本标题" min-width="200" />
        <el-table-column prop="scene" label="场景" min-width="160" />
        <el-table-column prop="location" label="时间 / 地点" min-width="200" />
        <el-table-column prop="lens" label="镜头设计" min-width="220" />
        <el-table-column prop="dialogue" label="台词提示" min-width="220" />
        <el-table-column label="操作" width="280" fixed="right">
          <template #default="{ row }">
            <el-button size="small" @click="demoAction(`导出拍摄清单：${row.scriptTitle}`)">导出清单</el-button>
            <el-button size="small" type="primary" @click="goShootTasks(row.merchant)">派任务</el-button>
          </template>
        </el-table-column>
      </el-table>

      <el-empty v-if="filteredStoryboards.length === 0" description="当前商家还没有分镜">
        <el-button type="primary" @click="demoAction(`AI 拆 ${merchantName || '商家'} 分镜`)">
          生成分镜
        </el-button>
      </el-empty>
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

.search-form {
  margin-bottom: 12px;
}
</style>
