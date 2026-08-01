<script setup lang="ts">
import { computed, ref, watch } from 'vue';
import { ElMessage } from 'element-plus';

import { publishSchedule } from '#/mock/content-ops';
import { useRoute } from 'vue-router';

defineOptions({ name: 'ContentPaymentsPage' });

const route = useRoute();
const keyword = ref('');

const merchantName = computed(() => String(route.query.merchantName || ''));
const merchantId = computed(() => String(route.query.merchantId || ''));

const filteredPublishSchedule = computed(() => {
  const searchText = keyword.value.trim();
  if (!searchText) return publishSchedule;
  return publishSchedule.filter((item) =>
    [item.merchant, item.title, item.publishTime, item.owner, item.status].some((value) =>
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
</script>

<template>
  <div class="p-4">
    <el-card>
      <template #header>
        <div class="page-head">
          <div>
            <div class="page-title">发布排期</div>
            <div class="page-desc">这里先看每条视频什么时候发、谁负责、当前卡在哪一步。</div>
          </div>
          <div class="page-actions">
            <el-button @click="demoAction('检查发布素材')">检查素材</el-button>
            <el-button
              type="primary"
              @click="demoAction(`新增 ${merchantName || '商家'} 发布排期`)"
            >
              新增排期
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
            placeholder="商家 / 视频 / 发布时间 / 负责人 / 状态"
            style="width: 320px"
          />
        </el-form-item>
        <el-form-item>
          <el-button @click="resetSearch">显示全部</el-button>
        </el-form-item>
      </el-form>

      <el-table :data="filteredPublishSchedule" border stripe>
        <el-table-column prop="merchant" label="商家" min-width="160" />
        <el-table-column prop="title" label="视频标题" min-width="220" />
        <el-table-column prop="publishTime" label="发布时间" width="170" />
        <el-table-column prop="owner" label="负责人" width="120" />
        <el-table-column prop="status" label="状态" width="110">
          <template #default="{ row }">
            <el-tag :type="row.status === '待发布' ? 'warning' : 'success'">{{ row.status }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="300" fixed="right">
          <template #default="{ row }">
            <el-button size="small" @click="demoAction(`老板确认：${row.title}`)">老板确认</el-button>
            <el-button size="small" type="primary" @click="demoAction(`发布：${row.title}`)">发布</el-button>
          </template>
        </el-table-column>
      </el-table>

      <el-empty v-if="filteredPublishSchedule.length === 0" description="当前商家还没有发布排期">
        <el-button
          type="primary"
          @click="demoAction(`新增 ${merchantName || '商家'} 发布排期`)"
        >
          新增排期
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
