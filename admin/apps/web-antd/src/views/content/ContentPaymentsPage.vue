<script setup lang="ts">
import { onMounted, reactive, ref, watch } from 'vue';
import { ElMessage } from 'element-plus';
import { useRoute, useRouter } from 'vue-router';

import {
  createPublishSchedule,
  getContentStoryboards,
  getPublishSchedules,
  updatePublishScheduleStatus,
  type ContentStoryboard,
  type PublishSchedule,
} from '#/api/admin/content-production';

defineOptions({ name: 'ContentPaymentsPage' });

const route = useRoute();
const router = useRouter();

const loading = ref(false);
const saving = ref(false);
const dialogVisible = ref(false);
const keyword = ref('');
const statusFilter = ref('');
const routeMerchantId = ref<number | undefined>();
const routeMerchantName = ref('');
const routeStoryboardId = ref<number | undefined>();
const routeScriptTitle = ref('');
const storyboardOptions = ref<ContentStoryboard[]>([]);
const list = ref<PublishSchedule[]>([]);
const pagination = reactive({ page: 1, size: 10, total: 0 });
const form = reactive({
  douyinAccount: '',
  materialStatus: '待拍摄',
  owner: '',
  publishTime: '',
  remark: '',
  status: '待发布',
  storyboardId: 0,
  videoTitle: '',
});

function applyRouteQuery() {
  const merchantId = Number(route.query.merchantId || 0);
  const storyboardId = Number(route.query.storyboardId || 0);
  routeMerchantId.value = merchantId > 0 ? merchantId : undefined;
  routeMerchantName.value = String(route.query.merchantName || '');
  routeStoryboardId.value = storyboardId > 0 ? storyboardId : undefined;
  routeScriptTitle.value = String(route.query.scriptTitle || '');
  keyword.value = routeScriptTitle.value || routeMerchantName.value;
  form.storyboardId = routeStoryboardId.value || form.storyboardId || 0;
  form.videoTitle = routeScriptTitle.value || form.videoTitle;
}

async function loadStoryboards() {
  const result = await getContentStoryboards({
    merchantId: routeMerchantId.value,
    page: 1,
    size: 100,
    status: '已确认',
  });
  storyboardOptions.value = result.list.length
    ? result.list
    : (
        await getContentStoryboards({
          merchantId: routeMerchantId.value,
          page: 1,
          size: 100,
        })
      ).list;
}

async function loadList() {
  loading.value = true;
  try {
    const result = await getPublishSchedules({
      keyword: keyword.value.trim(),
      merchantId: routeMerchantId.value,
      page: pagination.page,
      size: pagination.size,
      status: statusFilter.value,
      storyboardId: routeStoryboardId.value,
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
  routeStoryboardId.value = undefined;
  routeScriptTitle.value = '';
  router.replace('/content/payments');
  search();
}

async function openCreate() {
  await loadStoryboards();
  form.storyboardId = routeStoryboardId.value || form.storyboardId || storyboardOptions.value[0]?.id || 0;
  const selected = storyboardOptions.value.find((item) => item.id === form.storyboardId);
  form.videoTitle = selected?.scriptTitle || routeScriptTitle.value || form.videoTitle;
  dialogVisible.value = true;
}

async function submitCreate() {
  if (saving.value) return;
  if (!form.storyboardId) {
    ElMessage.warning('请选择分镜');
    return;
  }
  if (!form.videoTitle.trim()) {
    ElMessage.warning('请输入视频标题');
    return;
  }
  saving.value = true;
  try {
    const result = await createPublishSchedule({
      douyinAccount: form.douyinAccount,
      materialStatus: form.materialStatus,
      owner: form.owner,
      publishTime: form.publishTime,
      remark: form.remark,
      status: form.status,
      storyboardId: form.storyboardId,
      videoTitle: form.videoTitle,
    });
    ElMessage.success('发布排期已创建');
    dialogVisible.value = false;
    routeMerchantId.value = result.merchantId;
    routeMerchantName.value = result.merchantName;
    routeStoryboardId.value = result.storyboardId;
    await loadList();
  } catch (error) {
    ElMessage.error(getErrorMessage(error, '创建排期失败'));
  } finally {
    saving.value = false;
  }
}

async function markPublished(row: PublishSchedule) {
  await updatePublishScheduleStatus(row.id, {
    materialStatus: '已剪辑',
    status: '已发布',
  });
  ElMessage.success('已标记发布');
  await loadList();
}

function goReview(row: PublishSchedule) {
  router.push({
    path: '/logs/admin',
    query: {
      merchantId: row.merchantId,
      merchantName: row.merchantName,
      scheduleId: row.id,
      videoTitle: row.videoTitle,
    },
  });
}

function formatPublishTime(value?: string) {
  if (!value) return '-';
  return value.replace('T', ' ').slice(0, 16);
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
    await loadStoryboards();
    await loadList();
  },
);

onMounted(async () => {
  applyRouteQuery();
  await loadStoryboards();
  await loadList();
});
</script>

<template>
  <div class="p-4">
    <el-card>
      <template #header>
        <div class="page-head">
          <div>
            <div class="page-title">发布排期</div>
            <div class="page-desc">把已确认分镜排到拍摄、剪辑、发布和复盘。</div>
          </div>
          <div class="page-actions">
            <el-button @click="showAll">显示全部</el-button>
            <el-button type="primary" @click="openCreate">新增排期</el-button>
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
        当前处理：{{ routeMerchantName || '-' }}；脚本：{{ routeScriptTitle || '-' }}。
      </el-alert>

      <el-form :inline="true" class="search-form">
        <el-form-item label="关键词">
          <el-input
            v-model="keyword"
            clearable
            placeholder="商家 / 视频 / 负责人 / 抖音账号"
            style="width: 360px"
            @keyup.enter="search"
          />
        </el-form-item>
        <el-form-item label="状态">
          <el-select v-model="statusFilter" clearable placeholder="全部" style="width: 140px">
            <el-option label="待发布" value="待发布" />
            <el-option label="已发布" value="已发布" />
            <el-option label="已复盘" value="已复盘" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="search">搜索</el-button>
        </el-form-item>
      </el-form>

      <el-table v-loading="loading" :data="list" border stripe>
        <el-table-column prop="merchantName" label="商家" min-width="150" />
        <el-table-column prop="videoTitle" label="视频标题" min-width="240" />
        <el-table-column label="发布时间" width="170">
          <template #default="{ row }">{{ formatPublishTime(row.publishTime) }}</template>
        </el-table-column>
        <el-table-column prop="owner" label="负责人" width="110" />
        <el-table-column prop="materialStatus" label="素材状态" width="120" />
        <el-table-column prop="status" label="发布状态" width="110">
          <template #default="{ row }">
            <el-tag :type="row.status === '已发布' || row.status === '已复盘' ? 'success' : 'warning'">
              {{ row.status }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="260" fixed="right">
          <template #default="{ row }">
            <el-button size="small" type="primary" @click="markPublished(row)">标记发布</el-button>
            <el-button size="small" @click="goReview(row)">去复盘</el-button>
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

      <el-empty v-if="!loading && list.length === 0" description="还没有发布排期">
        <el-button type="primary" @click="openCreate">新增排期</el-button>
      </el-empty>
    </el-card>

    <el-dialog v-model="dialogVisible" title="新增发布排期" width="760px">
      <el-form label-width="120px">
        <el-form-item label="选择分镜" required>
          <el-select v-model="form.storyboardId" filterable placeholder="选择已确认或已有分镜" style="width: 100%">
            <el-option
              v-for="item in storyboardOptions"
              :key="item.id"
              :label="`${item.merchantName}｜${item.scriptTitle}`"
              :value="item.id"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="视频标题" required>
          <el-input v-model="form.videoTitle" placeholder="发布时使用的视频标题" />
        </el-form-item>
        <el-form-item label="发布时间">
          <el-date-picker
            v-model="form.publishTime"
            type="datetime"
            value-format="YYYY-MM-DD HH:mm:ss"
            placeholder="选择发布时间"
            style="width: 100%"
          />
        </el-form-item>
        <el-form-item label="负责人">
          <el-input v-model="form.owner" placeholder="拍摄/剪辑/发布负责人" />
        </el-form-item>
        <el-form-item label="抖音账号">
          <el-input v-model="form.douyinAccount" placeholder="商家抖音账号" />
        </el-form-item>
        <el-form-item label="素材状态">
          <el-select v-model="form.materialStatus" style="width: 100%">
            <el-option label="待拍摄" value="待拍摄" />
            <el-option label="已拍摄" value="已拍摄" />
            <el-option label="已剪辑" value="已剪辑" />
          </el-select>
        </el-form-item>
        <el-form-item label="备注">
          <el-input v-model="form.remark" type="textarea" :rows="3" placeholder="老板确认、发布注意事项等" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" :disabled="saving" :loading="saving" @click="submitCreate">
          保存排期
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
