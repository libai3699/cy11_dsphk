<script setup lang="ts">
import { onMounted, reactive, ref, watch } from 'vue';
import { ElMessage } from 'element-plus';
import { useRoute, useRouter } from 'vue-router';

import {
  getContentTopicList,
  type ContentTopic,
} from '#/api/admin';
import {
  generateContentScript,
  getContentScripts,
  updateContentScriptStatus,
  type ContentScript,
} from '#/api/admin/content-production';

defineOptions({ name: 'ContentQuotesPage' });

const route = useRoute();
const router = useRouter();

const loading = ref(false);
const generating = ref(false);
const dialogVisible = ref(false);
const keyword = ref('');
const statusFilter = ref('');
const routeMerchantId = ref<number | undefined>();
const routeMerchantName = ref('');
const routeTopicId = ref<number | undefined>();
const routeTopicTitle = ref('');
const topicOptions = ref<ContentTopic[]>([]);
const list = ref<ContentScript[]>([]);
const pagination = reactive({ page: 1, size: 10, total: 0 });
const form = reactive({
  extraRequirement: '',
  topicId: 0,
});

function applyRouteQuery() {
  const merchantId = Number(route.query.merchantId || 0);
  const topicId = Number(route.query.topicId || 0);
  routeMerchantId.value = merchantId > 0 ? merchantId : undefined;
  routeMerchantName.value = String(route.query.merchantName || '');
  routeTopicId.value = topicId > 0 ? topicId : undefined;
  routeTopicTitle.value = String(route.query.topicTitle || '');
  keyword.value = routeTopicTitle.value || routeMerchantName.value;
  form.topicId = routeTopicId.value || form.topicId || 0;
}

async function loadTopics() {
  const result = await getContentTopicList({
    merchantId: routeMerchantId.value,
    page: 1,
    size: 100,
    status: '已采用',
  });
  topicOptions.value = result.list.length
    ? result.list
    : (
        await getContentTopicList({
          merchantId: routeMerchantId.value,
          page: 1,
          size: 100,
        })
      ).list;
}

async function loadList() {
  loading.value = true;
  try {
    const result = await getContentScripts({
      keyword: keyword.value.trim(),
      merchantId: routeMerchantId.value,
      page: pagination.page,
      size: pagination.size,
      status: statusFilter.value,
      topicId: routeTopicId.value,
    });
    list.value = result.list;
    pagination.total = result.total;
    pagination.page = result.page;
    pagination.size = result.size;
  } finally {
    loading.value = false;
  }
}

function search() {
  pagination.page = 1;
  loadList();
}

function showAll() {
  keyword.value = '';
  statusFilter.value = '';
  routeMerchantId.value = undefined;
  routeMerchantName.value = '';
  routeTopicId.value = undefined;
  routeTopicTitle.value = '';
  router.replace('/content/quotes');
  search();
}

async function openGenerate(row?: ContentTopic) {
  await loadTopics();
  form.topicId = row?.id || routeTopicId.value || form.topicId || topicOptions.value[0]?.id || 0;
  dialogVisible.value = true;
}

async function submitGenerate() {
  if (generating.value) return;
  if (!form.topicId) {
    ElMessage.warning('请选择已确认的选题');
    return;
  }
  generating.value = true;
  try {
    const result = await generateContentScript({
      extraRequirement: form.extraRequirement,
      topicId: form.topicId,
    });
    ElMessage.success('文案生成完成');
    dialogVisible.value = false;
    routeMerchantId.value = result.merchantId;
    routeMerchantName.value = result.merchantName;
    routeTopicId.value = result.topicId;
    routeTopicTitle.value = result.topicTitle;
    await loadList();
  } catch (error) {
    ElMessage.error(getErrorMessage(error, '文案生成失败'));
  } finally {
    generating.value = false;
  }
}

async function confirmScript(row: ContentScript) {
  await updateContentScriptStatus(row.id, '已确认');
  ElMessage.success('文案已确认');
  await loadList();
}

function goStoryboards(row: ContentScript) {
  router.push({
    path: '/content/discoveries',
    query: {
      merchantId: row.merchantId,
      merchantName: row.merchantName,
      scriptId: row.id,
      scriptTitle: row.title,
      topicId: row.topicId,
      topicTitle: row.topicTitle,
    },
  });
}

function handlePageChange(page: number) {
  pagination.page = page;
  loadList();
}

function handleSizeChange(size: number) {
  pagination.size = size;
  pagination.page = 1;
  loadList();
}

function getErrorMessage(error: unknown, fallback: string) {
  const maybeError = error as { error?: string; message?: string };
  const message = maybeError?.error || maybeError?.message;
  return message ? `${fallback}：${message}` : fallback;
}

watch(
  () => route.query,
  async () => {
    applyRouteQuery();
    pagination.page = 1;
    await loadTopics();
    await loadList();
  },
);

onMounted(async () => {
  applyRouteQuery();
  await loadTopics();
  await loadList();
});
</script>

<template>
  <div class="p-4">
    <el-card>
      <template #header>
        <div class="page-head">
          <div>
            <div class="page-title">文案脚本</div>
            <div class="page-desc">从已确认选题生成可拍、可口播、可转化的短视频脚本。</div>
          </div>
          <div class="page-actions">
            <el-button @click="showAll">显示全部</el-button>
            <el-button type="primary" @click="openGenerate()">生成文案</el-button>
          </div>
        </div>
      </template>

      <el-alert
        v-if="routeMerchantName || routeTopicTitle"
        class="mb-4"
        type="info"
        :closable="false"
        show-icon
      >
        当前处理：{{ routeMerchantName || '-' }}；选题：{{ routeTopicTitle || '-' }}。
      </el-alert>

      <el-form :inline="true" class="search-form">
        <el-form-item label="关键词">
          <el-input
            v-model="keyword"
            clearable
            placeholder="商家 / 选题 / 标题 / 开场 / 脚本"
            style="width: 360px"
            @keyup.enter="search"
          />
        </el-form-item>
        <el-form-item label="状态">
          <el-select v-model="statusFilter" clearable placeholder="全部" style="width: 140px">
            <el-option label="草稿" value="草稿" />
            <el-option label="已确认" value="已确认" />
            <el-option label="停用" value="停用" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="search">搜索</el-button>
        </el-form-item>
      </el-form>

      <el-table v-loading="loading" :data="list" border stripe>
        <el-table-column prop="merchantName" label="商家" min-width="150" />
        <el-table-column prop="topicTitle" label="来源选题" min-width="220" />
        <el-table-column prop="title" label="脚本标题" min-width="220" />
        <el-table-column prop="opening" label="开场前三秒" min-width="260" show-overflow-tooltip />
        <el-table-column prop="cta" label="结尾引导" min-width="220" show-overflow-tooltip />
        <el-table-column prop="status" label="状态" width="110">
          <template #default="{ row }">
            <el-tag :type="row.status === '已确认' ? 'success' : 'warning'">{{ row.status }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="260" fixed="right">
          <template #default="{ row }">
            <el-button size="small" @click="confirmScript(row)">确认</el-button>
            <el-button size="small" type="primary" @click="goStoryboards(row)">生成分镜</el-button>
          </template>
        </el-table-column>
      </el-table>

      <el-pagination
        class="mt-4"
        layout="total, sizes, prev, pager, next"
        :current-page="pagination.page"
        :page-size="pagination.size"
        :page-sizes="[10, 20, 50]"
        :total="pagination.total"
        @current-change="handlePageChange"
        @size-change="handleSizeChange"
      />

      <el-empty v-if="!loading && list.length === 0" description="还没有文案">
        <el-button type="primary" @click="openGenerate()">生成文案</el-button>
      </el-empty>
    </el-card>

    <el-dialog v-model="dialogVisible" title="生成文案" width="720px">
      <el-form label-width="120px">
        <el-form-item label="选择选题" required>
          <el-select v-model="form.topicId" filterable placeholder="选择已确认或已有选题" style="width: 100%">
            <el-option
              v-for="item in topicOptions"
              :key="item.id"
              :label="`${item.merchantName}｜${item.title}`"
              :value="item.id"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="补充要求">
          <el-input
            v-model="form.extraRequirement"
            type="textarea"
            :rows="4"
            placeholder="例如：老板出镜、突出 89 元套餐、不要太硬广、控制在 45 秒内"
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" :disabled="generating" :loading="generating" @click="submitGenerate">
          调用文案 Agent
        </el-button>
      </template>
    </el-dialog>
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
