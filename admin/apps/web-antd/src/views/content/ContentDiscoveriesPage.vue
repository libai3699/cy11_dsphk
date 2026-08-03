<script setup lang="ts">
import { computed, onMounted, reactive, ref, watch } from 'vue';
import { ElMessage } from 'element-plus';
import { useRoute, useRouter } from 'vue-router';

import {
  generateContentStoryboard,
  getContentScripts,
  getContentStoryboards,
  updateContentStoryboardStatus,
  type ContentScript,
  type ContentStoryboard,
} from '#/api/admin/content-production';

defineOptions({ name: 'ContentDiscoveriesPage' });

const route = useRoute();
const router = useRouter();

const loading = ref(false);
const generating = ref(false);
const dialogVisible = ref(false);
const keyword = ref('');
const statusFilter = ref('');
const routeMerchantId = ref<number | undefined>();
const routeMerchantName = ref('');
const routeScriptId = ref<number | undefined>();
const routeScriptTitle = ref('');
const scriptOptions = ref<ContentScript[]>([]);
const list = ref<ContentStoryboard[]>([]);
const pagination = reactive({ page: 1, size: 10, total: 0 });
const expandedRowKeys = computed(() => list.value.map((item) => item.id));
const form = reactive({
  locationsText: '门店门头\n门店环境\n产品/套餐展示区\n收银/团购核销处',
  scriptId: 0,
});

function applyRouteQuery() {
  const merchantId = Number(route.query.merchantId || 0);
  const scriptId = Number(route.query.scriptId || 0);
  routeMerchantId.value = merchantId > 0 ? merchantId : undefined;
  routeMerchantName.value = String(route.query.merchantName || '');
  routeScriptId.value = scriptId > 0 ? scriptId : undefined;
  routeScriptTitle.value = String(route.query.scriptTitle || '');
  keyword.value = routeScriptTitle.value || routeMerchantName.value;
  form.scriptId = routeScriptId.value || form.scriptId || 0;
}

function splitLines(value: string) {
  return value
    .split('\n')
    .map((item) => item.trim())
    .filter(Boolean);
}

async function loadScripts() {
  const result = await getContentScripts({
    merchantId: routeMerchantId.value,
    page: 1,
    size: 100,
    status: '已确认',
  });
  scriptOptions.value = result.list.length
    ? result.list
    : (
        await getContentScripts({
          merchantId: routeMerchantId.value,
          page: 1,
          size: 100,
        })
      ).list;
}

async function loadList() {
  loading.value = true;
  try {
    const result = await getContentStoryboards({
      keyword: keyword.value.trim(),
      merchantId: routeMerchantId.value,
      page: pagination.page,
      scriptId: routeScriptId.value,
      size: pagination.size,
      status: statusFilter.value,
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
  routeScriptId.value = undefined;
  routeScriptTitle.value = '';
  router.replace('/content/discoveries');
  search();
}

async function openGenerate(row?: ContentScript) {
  await loadScripts();
  form.scriptId = row?.id || routeScriptId.value || form.scriptId || scriptOptions.value[0]?.id || 0;
  dialogVisible.value = true;
}

async function submitGenerate() {
  if (generating.value) return;
  if (!form.scriptId) {
    ElMessage.warning('请选择文案');
    return;
  }
  generating.value = true;
  try {
    const result = await generateContentStoryboard({
      locations: splitLines(form.locationsText),
      scriptId: form.scriptId,
    });
    ElMessage.success('分镜生成完成');
    dialogVisible.value = false;
    routeMerchantId.value = result.merchantId;
    routeMerchantName.value = result.merchantName;
    routeScriptId.value = result.scriptId;
    routeScriptTitle.value = result.scriptTitle;
    await loadList();
  } catch (error) {
    ElMessage.error(getErrorMessage(error, '分镜生成失败'));
  } finally {
    generating.value = false;
  }
}

async function confirmStoryboard(row: ContentStoryboard) {
  await updateContentStoryboardStatus(row.id, '已确认');
  ElMessage.success('分镜已确认');
  await loadList();
}

function goShootingTask(row: ContentStoryboard) {
  router.push({
    path: '/logs/user',
    query: {
      merchantId: row.merchantId,
      merchantName: row.merchantName,
      scriptId: row.scriptId,
      scriptTitle: row.scriptTitle,
      storyboardId: row.id,
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
    await loadScripts();
    await loadList();
  },
);

onMounted(async () => {
  applyRouteQuery();
  await loadScripts();
  await loadList();
});
</script>

<template>
  <div class="p-4">
    <el-card>
      <template #header>
        <div class="page-head">
          <div>
            <div class="page-title">分镜脚本</div>
            <div class="page-desc">把文案拆成拍摄人员能直接执行的镜头清单。</div>
          </div>
          <div class="page-actions">
            <el-button @click="showAll">显示全部</el-button>
            <el-button type="primary" @click="openGenerate()">生成分镜</el-button>
          </div>
        </div>
      </template>

      <el-alert
        v-if="routeMerchantName || routeScriptTitle"
        class="mb-4"
        type="info"
        :closable="false"
        show-icon
      >
        当前处理：{{ routeMerchantName || '-' }}；文案：{{ routeScriptTitle || '-' }}。
      </el-alert>

      <el-form :inline="true" class="search-form">
        <el-form-item label="关键词">
          <el-input
            v-model="keyword"
            clearable
            placeholder="商家 / 选题 / 文案 / 镜头内容"
            style="width: 360px"
            @keyup.enter="search"
          />
        </el-form-item>
        <el-form-item label="状态">
          <el-select v-model="statusFilter" clearable placeholder="全部" style="width: 140px">
            <el-option label="草稿" value="草稿" />
            <el-option label="已确认" value="已确认" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="search">搜索</el-button>
        </el-form-item>
      </el-form>

      <el-table
        v-loading="loading"
        :data="list"
        border
        :expand-row-keys="expandedRowKeys"
        row-key="id"
        stripe
      >
        <el-table-column type="expand" width="48">
          <template #default="{ row }">
            <div class="shots-panel">
              <div class="shots-title">完整分镜：共 {{ row.shots?.length || 0 }} 个镜头</div>
              <el-table :data="row.shots || []" border size="small">
                <el-table-column prop="index" label="#" width="60" />
                <el-table-column prop="duration" label="时长" width="100" />
                <el-table-column prop="location" label="地点" min-width="160" />
                <el-table-column prop="camera" label="镜头" min-width="160" />
                <el-table-column prop="content" label="画面内容" min-width="240" />
                <el-table-column prop="line" label="台词/字幕" min-width="240" />
                <el-table-column prop="note" label="注意事项" min-width="220" />
              </el-table>
            </div>
          </template>
        </el-table-column>
        <el-table-column prop="merchantName" label="商家" min-width="150" />
        <el-table-column prop="scriptTitle" label="来源文案" min-width="220" />
        <el-table-column label="镜头数" width="90">
          <template #default="{ row }">{{ row.shots?.length || 0 }}</template>
        </el-table-column>
        <el-table-column label="镜头概览" min-width="360">
          <template #default="{ row }">
            <span v-if="row.shots?.length">
              已默认展开全部 {{ row.shots.length }} 个镜头；首镜头：{{ row.shots[0].duration }}｜{{ row.shots[0].location }}
            </span>
            <span v-else>-</span>
          </template>
        </el-table-column>
        <el-table-column prop="status" label="状态" width="110">
          <template #default="{ row }">
            <el-tag :type="row.status === '已确认' ? 'success' : 'warning'">{{ row.status }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="250" fixed="right">
          <template #default="{ row }">
            <el-button size="small" @click="confirmStoryboard(row)">确认</el-button>
            <el-button size="small" type="primary" @click="goShootingTask(row)">派拍摄</el-button>
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

      <el-empty v-if="!loading && list.length === 0" description="还没有分镜">
        <el-button type="primary" @click="openGenerate()">生成分镜</el-button>
      </el-empty>
    </el-card>

    <el-dialog v-model="dialogVisible" title="生成分镜" width="720px">
      <el-form label-width="120px">
        <el-form-item label="选择文案" required>
          <el-select v-model="form.scriptId" filterable placeholder="选择已确认或已有文案" style="width: 100%">
            <el-option
              v-for="item in scriptOptions"
              :key="item.id"
              :label="`${item.merchantName}｜${item.title}`"
              :value="item.id"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="可拍地点">
          <el-input
            v-model="form.locationsText"
            type="textarea"
            :rows="5"
            placeholder="一行一个地点，例如：门头、厨房、包间、收银台"
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" :disabled="generating" :loading="generating" @click="submitGenerate">
          调用分镜 Agent
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

.shots-panel {
  padding: 12px 24px;
  background: #f8fafc;
}

.shots-title {
  margin-bottom: 10px;
  color: #0f172a;
  font-weight: 700;
}
</style>
