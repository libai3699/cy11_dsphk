<script setup lang="ts">
import { computed, onMounted, reactive, ref, watch } from 'vue';
import { ElMessage } from 'element-plus';
import { useRoute, useRouter } from 'vue-router';

import {
  getContentStoryboards,
  type ContentStoryboard,
} from '#/api/admin/content-production';
import {
  createShootingTask,
  getShootingTasks,
  updateShootingTaskStatus,
  type ShootingTask,
} from '#/api/admin/shooting-tasks';

defineOptions({ name: 'LogUserPage' });

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
const list = ref<ShootingTask[]>([]);
const pagination = reactive({ page: 1, size: 10, total: 0 });
const expandedRowKeys = computed(() => list.value.map((item) => item.id));
const form = reactive({
  assignee: '',
  deadline: '',
  materialUrl: '',
  remark: '',
  shootTime: '',
  status: '待拍摄',
  storyboardId: 0,
  taskTitle: '',
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
  form.taskTitle = routeScriptTitle.value ? `${routeScriptTitle.value}｜拍摄任务` : form.taskTitle;
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
    const result = await getShootingTasks({
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
  router.replace('/logs/user');
  search();
}

async function openCreate() {
  await loadStoryboards();
  form.storyboardId = routeStoryboardId.value || form.storyboardId || storyboardOptions.value[0]?.id || 0;
  const selected = storyboardOptions.value.find((item) => item.id === form.storyboardId);
  form.taskTitle = selected ? `${selected.scriptTitle}｜拍摄任务` : form.taskTitle;
  dialogVisible.value = true;
}

async function submitCreate() {
  if (saving.value) return;
  if (!form.storyboardId) {
    ElMessage.warning('请选择分镜');
    return;
  }
  saving.value = true;
  try {
    const result = await createShootingTask({
      assignee: form.assignee,
      deadline: form.deadline,
      materialUrl: form.materialUrl,
      remark: form.remark,
      shootTime: form.shootTime,
      status: form.status,
      storyboardId: form.storyboardId,
      taskTitle: form.taskTitle,
    });
    ElMessage.success('拍摄任务已创建');
    dialogVisible.value = false;
    routeMerchantId.value = result.merchantId;
    routeMerchantName.value = result.merchantName;
    routeStoryboardId.value = result.storyboardId;
    await loadList();
  } catch (error) {
    const maybeError = error as { error?: string; message?: string };
    ElMessage.error(maybeError?.error || maybeError?.message || '创建拍摄任务失败');
  } finally {
    saving.value = false;
  }
}

async function updateStatus(row: ShootingTask, status: string) {
  await updateShootingTaskStatus(row.id, { status });
  ElMessage.success(`已更新为：${status}`);
  await loadList();
}

function goPublish(row: ShootingTask) {
  router.push({
    path: '/content/payments',
    query: {
      merchantId: row.merchantId,
      merchantName: row.merchantName,
      scriptId: row.scriptId,
      scriptTitle: row.scriptTitle,
      storyboardId: row.storyboardId,
      topicId: row.topicId,
      topicTitle: row.topicTitle,
    },
  });
}

function formatTime(value?: string) {
  if (!value) return '-';
  return value.replace('T', ' ').slice(0, 16);
}

function handlePageChange(page: number) {
  pagination.page = page;
  loadList();
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
            <div class="page-title">拍摄任务</div>
            <div class="page-desc">分镜确认后，在这里派拍摄、跟进拍摄/剪辑状态，再进入发布排期。</div>
          </div>
          <div class="page-actions">
            <el-button @click="showAll">显示全部</el-button>
            <el-button type="primary" @click="openCreate">新增任务</el-button>
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
            placeholder="商家 / 任务 / 镜头 / 执行人 / 状态"
            style="width: 360px"
            @keyup.enter="search"
          />
        </el-form-item>
        <el-form-item label="状态">
          <el-select v-model="statusFilter" clearable placeholder="全部" style="width: 140px">
            <el-option label="待拍摄" value="待拍摄" />
            <el-option label="拍摄中" value="拍摄中" />
            <el-option label="已拍摄" value="已拍摄" />
            <el-option label="已剪辑" value="已剪辑" />
            <el-option label="已完成" value="已完成" />
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
              <div class="shots-title">拍摄清单：共 {{ row.shots?.length || row.shotCount || 0 }} 个镜头</div>
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
        <el-table-column prop="taskTitle" label="任务" min-width="240" />
        <el-table-column prop="merchantName" label="商家" min-width="150" />
        <el-table-column prop="shotCount" label="镜头数" width="90" />
        <el-table-column prop="assignee" label="执行人" width="120" />
        <el-table-column label="拍摄时间" width="170">
          <template #default="{ row }">{{ formatTime(row.shootTime) }}</template>
        </el-table-column>
        <el-table-column prop="status" label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="row.status === '已完成' || row.status === '已剪辑' ? 'success' : 'warning'">
              {{ row.status }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="320" fixed="right">
          <template #default="{ row }">
            <el-button size="small" @click="updateStatus(row, '拍摄中')">开拍</el-button>
            <el-button size="small" @click="updateStatus(row, '已剪辑')">已剪辑</el-button>
            <el-button size="small" type="primary" @click="goPublish(row)">排发布</el-button>
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

    <el-dialog v-model="dialogVisible" title="新增拍摄任务" width="760px">
      <el-form label-width="120px">
        <el-row :gutter="16">
          <el-col :span="12">
            <el-form-item label="选择分镜" required>
              <el-select v-model="form.storyboardId" filterable placeholder="选择分镜" style="width: 100%">
                <el-option
                  v-for="item in storyboardOptions"
                  :key="item.id"
                  :label="`${item.merchantName}｜${item.scriptTitle}`"
                  :value="item.id"
                />
              </el-select>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="执行人">
              <el-input v-model="form.assignee" placeholder="拍摄/剪辑负责人" />
            </el-form-item>
          </el-col>
        </el-row>
        <el-form-item label="任务标题">
          <el-input v-model="form.taskTitle" placeholder="默认按文案标题生成" />
        </el-form-item>
        <el-row :gutter="16">
          <el-col :span="12">
            <el-form-item label="拍摄时间">
              <el-date-picker
                v-model="form.shootTime"
                type="datetime"
                value-format="YYYY-MM-DD HH:mm:ss"
                placeholder="选择拍摄时间"
                style="width: 100%"
              />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="截止时间">
              <el-date-picker
                v-model="form.deadline"
                type="datetime"
                value-format="YYYY-MM-DD HH:mm:ss"
                placeholder="选择截止时间"
                style="width: 100%"
              />
            </el-form-item>
          </el-col>
        </el-row>
        <el-form-item label="备注">
          <el-input v-model="form.remark" type="textarea" :rows="3" placeholder="拍摄注意事项、老板确认要求等" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" :disabled="saving" :loading="saving" @click="submitCreate">
          保存任务
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
