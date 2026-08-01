<script setup lang="ts">
import { computed, ref, watch } from 'vue';
import { ElMessage } from 'element-plus';
import { useRoute, useRouter } from 'vue-router';

import { scripts } from '#/mock/content-ops';

defineOptions({ name: 'ContentQuotesPage' });

const route = useRoute();
const router = useRouter();
const keyword = ref('');

const merchantName = computed(() => String(route.query.merchantName || ''));
const merchantId = computed(() => String(route.query.merchantId || ''));

const filteredScripts = computed(() => {
  const searchText = keyword.value.trim();
  if (!searchText) return scripts;
  return scripts.filter((item) =>
    [item.merchant, item.scriptTitle, item.opening, item.cta, item.status].some((value) =>
      value.includes(searchText),
    ),
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

function goStoryboards(rowMerchant?: string) {
  const targetMerchant = rowMerchant || merchantName.value;
  router.push({
    path: '/content/discoveries',
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
            <div class="page-title">文案脚本</div>
            <div class="page-desc">运营人员最终要交付的，不是“AI 提示词”，而是可直接拍的脚本。</div>
          </div>
          <div class="page-actions">
            <el-button @click="demoAction('批量改写文案')">批量改写</el-button>
            <el-button
              type="primary"
              @click="demoAction(`AI 生成 ${merchantName || '商家'} 文案`)"
            >
              生成文案
            </el-button>
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
            placeholder="商家 / 脚本 / 开场 / 引导 / 状态"
            style="width: 320px"
          />
        </el-form-item>
        <el-form-item>
          <el-button @click="resetSearch">显示全部</el-button>
        </el-form-item>
      </el-form>

      <el-table :data="filteredScripts" border stripe>
        <el-table-column prop="merchant" label="商家" min-width="160" />
        <el-table-column prop="scriptTitle" label="脚本标题" min-width="200" />
        <el-table-column prop="opening" label="开场前三秒" min-width="240" />
        <el-table-column prop="cta" label="结尾引导" min-width="220" />
        <el-table-column prop="status" label="状态" width="120">
          <template #default="{ row }">
            <el-tag :type="row.status === '已过审' ? 'success' : 'warning'">{{ row.status }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="300" fixed="right">
          <template #default="{ row }">
            <el-button size="small" @click="demoAction(`送审脚本：${row.scriptTitle}`)">送老板确认</el-button>
            <el-button size="small" type="primary" @click="goStoryboards(row.merchant)">生成分镜</el-button>
          </template>
        </el-table-column>
      </el-table>

      <el-empty v-if="filteredScripts.length === 0" description="当前商家还没有文案">
        <el-button
          type="primary"
          @click="demoAction(`AI 生成 ${merchantName || '商家'} 文案`)"
        >
          生成文案
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
