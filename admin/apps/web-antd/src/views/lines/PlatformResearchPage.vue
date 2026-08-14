<script setup lang="ts">
import type {
  PlatformCaseStudy,
  PlatformResearchTask,
  PlatformSearchResult,
} from '#/api/admin/platform-research';

import { ElMessage } from 'element-plus';
import { computed, onMounted, reactive, ref } from 'vue';
import { useRoute } from 'vue-router';

import {
  generatePlatformResearch,
  getMerchantList,
  getPlatformResearchList,
} from '#/api/admin';

const route = useRoute();

const loading = ref(false);
const generating = ref(false);
const detailVisible = ref(false);
const activeTask = ref<PlatformResearchTask>();
const list = ref<PlatformResearchTask[]>([]);
const total = ref(0);
const merchants = ref<{ id: number; name: string; city: string; industry: string }[]>([]);

const filters = reactive({
  keyword: '',
  merchantId: undefined as number | undefined,
  page: 1,
  size: 10,
});

const form = reactive({
  extraRequirement: '',
  keywordsText: '',
  limit: 5,
  merchantId: undefined as number | undefined,
  sources: ['douyin', 'xiaohongshu', 'meituan', 'eleme'] as string[],
});

const sourceOptions = [
  { label: '抖音', value: 'douyin' },
  { label: '小红书', value: 'xiaohongshu' },
  { label: '美团/大众点评', value: 'meituan' },
  { label: '饿了么', value: 'eleme' },
];

const platformMap: Record<string, string> = {
  douyin: '抖音',
  eleme: '饿了么',
  meituan: '美团',
  xiaohongshu: '小红书',
};

const activeSearchResults = computed<PlatformSearchResult[]>(() => {
  return activeTask.value?.searchResults || [];
});

const groupedResults = computed(() => {
  const groups: Record<string, PlatformSearchResult[]> = {};
  for (const item of activeSearchResults.value) {
    const key = item.platform || 'unknown';
    groups[key] ||= [];
    groups[key].push(item);
  }
  return groups;
});

function parseKeywords() {
  return form.keywordsText
    .split(/[\n,，]/)
    .map((item) => item.trim())
    .filter(Boolean);
}

async function loadMerchants() {
  const result = await getMerchantList({ page: 1, size: 100 });
  merchants.value = result.list.map((item) => ({
    city: item.city,
    id: item.id,
    industry: item.industry,
    name: item.name,
  }));
}

async function loadList() {
  loading.value = true;
  try {
    const result = await getPlatformResearchList({
      keyword: filters.keyword || undefined,
      merchantId: filters.merchantId,
      page: filters.page,
      size: filters.size,
    });
    list.value = result.list;
    total.value = result.total;
  } finally {
    loading.value = false;
  }
}

function resetSearch() {
  filters.keyword = '';
  filters.merchantId = undefined;
  filters.page = 1;
  loadList();
}

function openGenerate() {
  const merchantId = Number(route.query.merchantId || 0);
  if (merchantId > 0) {
    form.merchantId = merchantId;
  }
}

async function submitGenerate() {
  if (!form.merchantId) {
    ElMessage.warning('请选择商家');
    return;
  }
  if (form.sources.length === 0) {
    ElMessage.warning('至少选择一个平台');
    return;
  }

  generating.value = true;
  try {
    const task = await generatePlatformResearch({
      extraRequirement: form.extraRequirement,
      keywords: parseKeywords(),
      limit: form.limit,
      merchantId: form.merchantId,
      sources: form.sources,
    });
    activeTask.value = task;
    detailVisible.value = true;
    ElMessage.success('平台调研完成');
    await loadList();
  } catch (error) {
    const message =
      error instanceof Error ? error.message : '平台调研失败，请查看后端日志';
    ElMessage.error(message);
  } finally {
    generating.value = false;
  }
}

function showDetail(row: PlatformResearchTask) {
  activeTask.value = row;
  detailVisible.value = true;
}

function statusText(status: string) {
  if (status === 'completed') return '完成';
  if (status === 'failed') return '失败';
  return '处理中';
}

function statusType(status: string) {
  if (status === 'completed') return 'success';
  if (status === 'failed') return 'danger';
  return 'warning';
}

function platformName(value: string) {
  return platformMap[value] || value || '-';
}

function caseList(value?: PlatformCaseStudy[]) {
  return value || [];
}

onMounted(async () => {
  await loadMerchants();
  const merchantId = Number(route.query.merchantId || 0);
  if (merchantId > 0) {
    filters.merchantId = merchantId;
    form.merchantId = merchantId;
  }
  openGenerate();
  await loadList();
});
</script>

<template>
  <div class="platform-research-page">
    <el-card shadow="never">
      <template #header>
        <div class="card-header">
          <div>
            <div class="title">平台调研</div>
            <div class="sub-title">
              输入少量商家信息后，系统搜索抖音、小红书、美团、饿了么公开结果，再由 Agent 汇总好案例、坏案例和可执行洞察。
            </div>
          </div>
          <el-button :loading="generating" type="primary" @click="submitGenerate">
            开始调研
          </el-button>
        </div>
      </template>

      <el-form class="research-form" label-width="96px">
        <el-row :gutter="16">
          <el-col :lg="12" :md="24">
            <el-form-item label="商家">
              <el-select
                v-model="form.merchantId"
                clearable
                filterable
                placeholder="选择要调研的商家"
                style="width: 100%"
              >
                <el-option
                  v-for="item in merchants"
                  :key="item.id"
                  :label="`${item.name}｜${item.city || '-'}｜${item.industry || '-'}`"
                  :value="item.id"
                />
              </el-select>
            </el-form-item>
          </el-col>
          <el-col :lg="12" :md="24">
            <el-form-item label="平台">
              <el-checkbox-group v-model="form.sources">
                <el-checkbox-button
                  v-for="item in sourceOptions"
                  :key="item.value"
                  :label="item.value"
                >
                  {{ item.label }}
                </el-checkbox-button>
              </el-checkbox-group>
            </el-form-item>
          </el-col>
          <el-col :lg="12" :md="24">
            <el-form-item label="补充关键词">
              <el-input
                v-model="form.keywordsText"
                :autosize="{ minRows: 4, maxRows: 8 }"
                placeholder="可不填。每行一个，比如：花果园火锅 9.9 团购"
                type="textarea"
              />
            </el-form-item>
          </el-col>
          <el-col :lg="12" :md="24">
            <el-form-item label="调研要求">
              <el-input
                v-model="form.extraRequirement"
                :autosize="{ minRows: 4, maxRows: 8 }"
                placeholder="比如：不要低价伤利润，重点找能带到店核销的打法"
                type="textarea"
              />
            </el-form-item>
          </el-col>
          <el-col :lg="12" :md="24">
            <el-form-item label="单词结果数">
              <el-input-number v-model="form.limit" :max="10" :min="1" />
              <span class="form-tip">每个平台、每个关键词最多返回多少条公开搜索结果。</span>
            </el-form-item>
          </el-col>
        </el-row>
      </el-form>
    </el-card>

    <el-card class="list-card" shadow="never">
      <template #header>
        <div class="card-header">
          <div class="title">调研记录</div>
          <div class="filters">
            <el-select
              v-model="filters.merchantId"
              clearable
              filterable
              placeholder="全部商家"
              style="width: 260px"
              @change="loadList"
            >
              <el-option
                v-for="item in merchants"
                :key="item.id"
                :label="item.name"
                :value="item.id"
              />
            </el-select>
            <el-input
              v-model="filters.keyword"
              clearable
              placeholder="搜商家/行业/结论"
              style="width: 260px"
              @keyup.enter="loadList"
            />
            <el-button @click="loadList">查询</el-button>
            <el-button @click="resetSearch">重置</el-button>
          </div>
        </div>
      </template>

      <el-table v-loading="loading" :data="list" border>
        <el-table-column label="商家" min-width="180" prop="merchantName" />
        <el-table-column label="城市/行业" min-width="160">
          <template #default="{ row }">
            {{ row.city || '-' }} / {{ row.industry || '-' }}
          </template>
        </el-table-column>
        <el-table-column label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="statusType(row.status)">
              {{ statusText(row.status) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="搜索结果" width="100">
          <template #default="{ row }">
            {{ row.searchResults?.length || 0 }}
          </template>
        </el-table-column>
        <el-table-column label="结论" min-width="260" prop="summary" show-overflow-tooltip />
        <el-table-column label="时间" min-width="170" prop="createdAt" />
        <el-table-column fixed="right" label="操作" width="110">
          <template #default="{ row }">
            <el-button link type="primary" @click="showDetail(row)">查看</el-button>
          </template>
        </el-table-column>
      </el-table>

      <div class="pager">
        <el-pagination
          v-model:current-page="filters.page"
          v-model:page-size="filters.size"
          :page-sizes="[10, 20, 50]"
          :total="total"
          layout="total, sizes, prev, pager, next"
          @current-change="loadList"
          @size-change="loadList"
        />
      </div>
    </el-card>

    <el-dialog v-model="detailVisible" title="平台调研详情" width="960px">
      <template v-if="activeTask">
        <el-alert
          :closable="false"
          :title="activeTask.summary || activeTask.errorMessage || '暂无结论'"
          class="detail-alert"
          show-icon
          type="info"
        />

        <el-tabs>
          <el-tab-pane label="Agent 结论">
            <el-row :gutter="16">
              <el-col :md="12" :sm="24">
                <el-card shadow="never">
                  <template #header>好案例</template>
                  <el-empty v-if="caseList(activeTask.goodCases).length === 0" description="暂无" />
                  <div
                    v-for="(item, index) in caseList(activeTask.goodCases)"
                    :key="`${item.title}-${index}`"
                    class="case-item good"
                  >
                    <div class="case-title">
                      {{ platformName(item.platform) }}｜{{ item.title }}
                    </div>
                    <div class="case-line">原因：{{ item.reason }}</div>
                    <div class="case-line">可学：{{ item.takeaway }}</div>
                    <a v-if="item.url" :href="item.url" target="_blank">打开来源</a>
                  </div>
                </el-card>
              </el-col>
              <el-col :md="12" :sm="24">
                <el-card shadow="never">
                  <template #header>坏案例 / 避坑</template>
                  <el-empty v-if="caseList(activeTask.badCases).length === 0" description="暂无" />
                  <div
                    v-for="(item, index) in caseList(activeTask.badCases)"
                    :key="`${item.title}-${index}`"
                    class="case-item bad"
                  >
                    <div class="case-title">
                      {{ platformName(item.platform) }}｜{{ item.title }}
                    </div>
                    <div class="case-line">原因：{{ item.reason }}</div>
                    <div class="case-line">避开：{{ item.takeaway }}</div>
                    <a v-if="item.url" :href="item.url" target="_blank">打开来源</a>
                  </div>
                </el-card>
              </el-col>
            </el-row>

            <el-card class="detail-card" shadow="never">
              <template #header>可执行洞察</template>
              <el-tag
                v-for="item in activeTask.insights || []"
                :key="item"
                class="insight-tag"
                type="success"
              >
                {{ item }}
              </el-tag>
              <el-empty v-if="(activeTask.insights || []).length === 0" description="暂无" />
            </el-card>
          </el-tab-pane>

          <el-tab-pane label="搜索结果">
            <el-collapse :model-value="Object.keys(groupedResults)">
              <el-collapse-item
                v-for="(items, platform) in groupedResults"
                :key="platform"
                :name="platform"
                :title="`${platformName(String(platform))}（${items.length}）`"
              >
                <div
                  v-for="(item, index) in items"
                  :key="`${item.url}-${index}`"
                  class="search-item"
                >
                  <div class="search-title">
                    <a :href="item.url" target="_blank">{{ item.title }}</a>
                  </div>
                  <div class="search-meta">
                    关键词：{{ item.keyword }}｜来源：{{ item.source }}｜相关分：{{ item.score || 0 }}
                  </div>
                  <div v-if="item.query" class="search-meta">搜索式：{{ item.query }}</div>
                  <div class="search-snippet">{{ item.snippet }}</div>
                </div>
              </el-collapse-item>
            </el-collapse>
          </el-tab-pane>

          <el-tab-pane label="关键词">
            <el-tag v-for="item in activeTask.keywords || []" :key="item" class="keyword-tag">
              {{ item }}
            </el-tag>
          </el-tab-pane>
        </el-tabs>
      </template>
    </el-dialog>
  </div>
</template>

<style scoped>
.platform-research-page {
  display: flex;
  flex-direction: column;
  gap: 16px;
  padding: 16px;
}

.card-header {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  gap: 16px;
}

.title {
  font-size: 18px;
  font-weight: 600;
}

.sub-title {
  margin-top: 6px;
  color: var(--el-text-color-secondary);
  line-height: 1.6;
}

.research-form {
  max-width: 1280px;
}

.form-tip {
  margin-left: 12px;
  color: var(--el-text-color-secondary);
}

.list-card {
  margin-top: 16px;
}

.filters {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
  justify-content: flex-end;
}

.pager {
  display: flex;
  justify-content: flex-end;
  margin-top: 16px;
}

.detail-alert,
.detail-card {
  margin-bottom: 16px;
}

.case-item {
  padding: 12px;
  margin-bottom: 12px;
  border: 1px solid var(--el-border-color-lighter);
  border-radius: 8px;
}

.case-item.good {
  background: #f0f9eb;
}

.case-item.bad {
  background: #fef0f0;
}

.case-title {
  margin-bottom: 8px;
  font-weight: 600;
}

.case-line {
  margin-bottom: 6px;
  color: var(--el-text-color-regular);
  line-height: 1.5;
}

.insight-tag,
.keyword-tag {
  max-width: 100%;
  height: auto;
  margin: 0 8px 8px 0;
  padding: 8px 10px;
  white-space: normal;
}

.search-item {
  padding: 12px 0;
  border-bottom: 1px solid var(--el-border-color-lighter);
}

.search-title {
  margin-bottom: 6px;
  font-weight: 600;
}

.search-meta {
  margin-bottom: 6px;
  color: var(--el-text-color-secondary);
  font-size: 12px;
}

.search-snippet {
  color: var(--el-text-color-regular);
  line-height: 1.6;
}
</style>
