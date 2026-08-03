<script setup lang="ts">
import { onMounted, reactive, ref, watch } from 'vue';
import { ElMessage } from 'element-plus';
import { useRoute, useRouter } from 'vue-router';

import {
  generateContentReview,
  getContentReviews,
  getPublishSchedules,
  type ContentReviewTask,
  type PublishSchedule,
} from '#/api/admin/content-production';

defineOptions({ name: 'LogAdminPage' });

const route = useRoute();
const router = useRouter();

const loading = ref(false);
const generating = ref(false);
const dialogVisible = ref(false);
const routeMerchantId = ref<number | undefined>();
const routeMerchantName = ref('');
const routeScheduleId = ref<number | undefined>();
const routeVideoTitle = ref('');
const scheduleOptions = ref<PublishSchedule[]>([]);
const list = ref<ContentReviewTask[]>([]);
const pagination = reactive({ page: 1, size: 10, total: 0 });
const form = reactive({
  commentCount: 0,
  dealCount: 0,
  likeCount: 0,
  periodEnd: '',
  periodStart: '',
  playCount: 0,
  scheduleId: 0,
  shareCount: 0,
  writeOffAmount: 0,
});

function applyRouteQuery() {
  const merchantId = Number(route.query.merchantId || 0);
  const scheduleId = Number(route.query.scheduleId || 0);
  routeMerchantId.value = merchantId > 0 ? merchantId : undefined;
  routeMerchantName.value = String(route.query.merchantName || '');
  routeScheduleId.value = scheduleId > 0 ? scheduleId : undefined;
  routeVideoTitle.value = String(route.query.videoTitle || '');
  form.scheduleId = routeScheduleId.value || form.scheduleId || 0;
}

async function loadSchedules() {
  const result = await getPublishSchedules({
    merchantId: routeMerchantId.value,
    page: 1,
    size: 100,
    status: '已发布',
  });
  scheduleOptions.value = result.list.length
    ? result.list
    : (
        await getPublishSchedules({
          merchantId: routeMerchantId.value,
          page: 1,
          size: 100,
        })
      ).list;
}

async function loadList() {
  loading.value = true;
  try {
    const result = await getContentReviews({
      merchantId: routeMerchantId.value,
      page: pagination.page,
      size: pagination.size,
    });
    list.value = result.list;
    pagination.total = result.total;
    pagination.page = result.page;
    pagination.size = result.size;
  } finally {
    loading.value = false;
  }
}

async function openReview() {
  await loadSchedules();
  form.scheduleId = routeScheduleId.value || form.scheduleId || scheduleOptions.value[0]?.id || 0;
  dialogVisible.value = true;
}

async function submitReview() {
  if (generating.value) return;
  if (!form.scheduleId) {
    ElMessage.warning('请选择已发布视频');
    return;
  }
  generating.value = true;
  try {
    await generateContentReview({
      commentCount: form.commentCount,
      dealCount: form.dealCount,
      likeCount: form.likeCount,
      periodEnd: form.periodEnd,
      periodStart: form.periodStart,
      playCount: form.playCount,
      scheduleId: form.scheduleId,
      shareCount: form.shareCount,
      writeOffAmount: form.writeOffAmount,
    });
    ElMessage.success('复盘生成完成');
    dialogVisible.value = false;
    await loadList();
  } catch (error) {
    const maybeError = error as { error?: string; message?: string };
    ElMessage.error(maybeError?.error || maybeError?.message || '复盘生成失败');
  } finally {
    generating.value = false;
  }
}

function goTopics(row: ContentReviewTask) {
  router.push({
    path: '/content/notices',
    query: {
      merchantId: row.merchantId,
      merchantName: row.merchantName,
    },
  });
}

function handlePageChange(page: number) {
  pagination.page = page;
  loadList();
}

function formatDate(value?: string) {
  if (!value) return '-';
  return value.slice(0, 10);
}

watch(
  () => route.query,
  async () => {
    applyRouteQuery();
    await loadSchedules();
    await loadList();
  },
);

onMounted(async () => {
  applyRouteQuery();
  await loadSchedules();
  await loadList();
});
</script>

<template>
  <div class="p-4">
    <el-card>
      <template #header>
        <div class="page-head">
          <div>
            <div class="page-title">数据复盘</div>
            <div class="page-desc">视频发出去后，回到播放、互动、成交和下一轮选题优化。</div>
          </div>
          <div class="page-actions">
            <el-button @click="loadList">刷新</el-button>
            <el-button type="primary" @click="openReview">单条复盘</el-button>
          </div>
        </div>
      </template>

      <el-alert
        v-if="routeMerchantName || routeVideoTitle"
        class="mb-4"
        type="info"
        :closable="false"
        show-icon
      >
        当前处理：{{ routeMerchantName || '-' }}；视频：{{ routeVideoTitle || '-' }}。
      </el-alert>

      <el-table v-loading="loading" :data="list" border stripe>
        <el-table-column prop="merchantName" label="商家" min-width="150" />
        <el-table-column prop="videoTitle" label="视频" min-width="240" />
        <el-table-column label="复盘周期" width="190">
          <template #default="{ row }">
            {{ formatDate(row.periodStart) }} 至 {{ formatDate(row.periodEnd) }}
          </template>
        </el-table-column>
        <el-table-column prop="playCount" label="播放" width="100" />
        <el-table-column prop="likeCount" label="点赞" width="100" />
        <el-table-column prop="dealCount" label="成交单数" width="110" />
        <el-table-column prop="writeOffAmount" label="核销金额" width="120" />
        <el-table-column label="复盘结论" min-width="300" show-overflow-tooltip>
          <template #default="{ row }">{{ row.result?.conclusion || '-' }}</template>
        </el-table-column>
        <el-table-column label="操作" width="180" fixed="right">
          <template #default="{ row }">
            <el-button size="small" type="primary" @click="goTopics(row)">生成下轮选题</el-button>
          </template>
        </el-table-column>
      </el-table>

      <el-pagination
        class="mt-4"
        layout="total, prev, pager, next"
        :current-page="pagination.page"
        :page-size="pagination.size"
        :total="pagination.total"
        @current-change="handlePageChange"
      />
    </el-card>

    <el-dialog v-model="dialogVisible" title="单条视频复盘" width="760px">
      <el-form label-width="120px">
        <el-form-item label="选择视频" required>
          <el-select v-model="form.scheduleId" filterable placeholder="选择已发布视频" style="width: 100%">
            <el-option
              v-for="item in scheduleOptions"
              :key="item.id"
              :label="`${item.merchantName}｜${item.videoTitle}`"
              :value="item.id"
            />
          </el-select>
        </el-form-item>
        <el-row :gutter="16">
          <el-col :span="12">
            <el-form-item label="开始日期">
              <el-date-picker
                v-model="form.periodStart"
                type="date"
                value-format="YYYY-MM-DD"
                placeholder="复盘开始日期"
                style="width: 100%"
              />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="结束日期">
              <el-date-picker
                v-model="form.periodEnd"
                type="date"
                value-format="YYYY-MM-DD"
                placeholder="复盘结束日期"
                style="width: 100%"
              />
            </el-form-item>
          </el-col>
        </el-row>
        <el-row :gutter="16">
          <el-col :span="12">
            <el-form-item label="播放">
              <el-input-number v-model="form.playCount" :min="0" style="width: 100%" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="点赞">
              <el-input-number v-model="form.likeCount" :min="0" style="width: 100%" />
            </el-form-item>
          </el-col>
        </el-row>
        <el-row :gutter="16">
          <el-col :span="12">
            <el-form-item label="评论">
              <el-input-number v-model="form.commentCount" :min="0" style="width: 100%" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="分享">
              <el-input-number v-model="form.shareCount" :min="0" style="width: 100%" />
            </el-form-item>
          </el-col>
        </el-row>
        <el-row :gutter="16">
          <el-col :span="12">
            <el-form-item label="成交单数">
              <el-input-number v-model="form.dealCount" :min="0" style="width: 100%" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="核销金额">
              <el-input-number v-model="form.writeOffAmount" :min="0" :precision="2" style="width: 100%" />
            </el-form-item>
          </el-col>
        </el-row>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" :disabled="generating" :loading="generating" @click="submitReview">
          调用复盘 Agent
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
</style>
