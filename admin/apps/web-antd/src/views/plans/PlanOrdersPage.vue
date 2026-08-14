<script setup lang="ts">
import type { FormInstance } from 'element-plus';

import { computed, onMounted, reactive, ref, watch } from 'vue';
import { ElMessage, ElMessageBox } from 'element-plus';
import { useRoute, useRouter } from 'vue-router';

import {
  createSettlementOrder,
  deleteSettlementOrder,
  generateSettlementOrders,
  getMerchantList,
  getPublishSchedules,
  getSettlementOrderList,
  updateSettlementOrder,
  updateSettlementOrderStatus,
  type Merchant,
  type PublishSchedule,
  type SettlementOrder,
  type SettlementOrderPayload,
} from '#/api/admin';

defineOptions({ name: 'PlanOrdersPage' });

const route = useRoute();
const router = useRouter();

const statuses = ['待核对', '已确认', '已结算'];

const loading = ref(false);
const submitting = ref(false);
const generating = ref(false);
const dialogVisible = ref(false);
const editingId = ref<number | null>(null);
const formRef = ref<FormInstance>();

const keyword = ref('');
const status = ref('');
const routeMerchantId = ref<number | undefined>();
const routeMerchantName = ref('');
const pagination = reactive({ page: 1, size: 10, total: 0 });
const list = ref<SettlementOrder[]>([]);
const merchantOptions = ref<Merchant[]>([]);
const scheduleOptions = ref<PublishSchedule[]>([]);

const form = reactive<SettlementOrderPayload>({
  commissionRate: 10,
  merchantId: 0,
  orderWindow: '',
  periodEnd: '',
  periodStart: '',
  redeemedAmount: 0,
  remark: '',
  scheduleId: 0,
  sourceVideo: '',
  status: '待核对',
  videoTitle: '',
});

const dialogTitle = computed(() =>
  editingId.value ? '编辑分成订单' : '新增分成订单',
);

const currentMerchantName = computed(() => {
  if (routeMerchantName.value) return routeMerchantName.value;
  const current = merchantOptions.value.find(
    (item) => item.id === routeMerchantId.value,
  );
  return current?.name || '';
});

const commissionPreview = computed(() =>
  Number(
    ((((form.redeemedAmount || 0) * (form.commissionRate || 0)) / 100)).toFixed(
      2,
    ),
  ),
);

const rules = {
  merchantId: [{ message: '请选择商家', required: true, trigger: 'change' }],
  redeemedAmount: [{ message: '请输入核销额', required: true, trigger: 'blur' }],
};

function applyRouteQuery() {
  const rawMerchantId = Number(route.query.merchantId || 0);
  routeMerchantId.value = rawMerchantId > 0 ? rawMerchantId : undefined;
  routeMerchantName.value = String(route.query.merchantName || '');
  if (routeMerchantName.value) {
    keyword.value = routeMerchantName.value;
  }
}

function formatMoney(value: number) {
  return `¥${Number(value || 0).toFixed(2)}`;
}

function formatPercent(value: number) {
  return `${Number(value || 0).toFixed(2)}%`;
}

function formatDate(value?: string) {
  if (!value) return '-';
  return value.slice(0, 10);
}

function statusTagType(rowStatus: string) {
  if (rowStatus === '已结算') return 'success';
  if (rowStatus === '已确认') return 'primary';
  return 'warning';
}

function resetForm() {
  const currentMerchant = merchantOptions.value.find(
    (item) => item.id === routeMerchantId.value,
  );
  Object.assign(form, {
    commissionRate: currentMerchant?.commissionRate || 10,
    merchantId: routeMerchantId.value || 0,
    orderWindow: '',
    periodEnd: '',
    periodStart: '',
    redeemedAmount: 0,
    remark: '',
    scheduleId: 0,
    sourceVideo: '',
    status: '待核对',
    videoTitle: '',
  });
  formRef.value?.clearValidate();
}

function rowToPayload(row: SettlementOrder): SettlementOrderPayload {
  return {
    commissionRate: row.commissionRate,
    merchantId: row.merchantId,
    orderWindow: row.orderWindow,
    periodEnd: formatDateForForm(row.periodEnd),
    periodStart: formatDateForForm(row.periodStart),
    redeemedAmount: row.redeemedAmount,
    remark: row.remark,
    scheduleId: row.scheduleId,
    sourceVideo: row.sourceVideo,
    status: row.status,
    videoTitle: row.videoTitle,
  };
}

function formatDateForForm(value?: string) {
  if (!value) return '';
  return value.slice(0, 10);
}

async function loadMerchants() {
  const result = await getMerchantList({ page: 1, size: 100 });
  merchantOptions.value = result.list;
}

async function loadSchedules() {
  const result = await getPublishSchedules({
    merchantId: form.merchantId || routeMerchantId.value,
    page: 1,
    size: 100,
  });
  scheduleOptions.value = result.list;
}

async function loadList() {
  loading.value = true;
  try {
    const result = await getSettlementOrderList({
      keyword: keyword.value.trim(),
      merchantId: routeMerchantId.value,
      page: pagination.page,
      size: pagination.size,
      status: status.value,
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
  status.value = '';
  routeMerchantId.value = undefined;
  routeMerchantName.value = '';
  router.replace('/plans/orders');
  search();
}

async function openCreate() {
  editingId.value = null;
  resetForm();
  await loadSchedules();
  dialogVisible.value = true;
}

async function openEdit(row: SettlementOrder) {
  editingId.value = row.id;
  Object.assign(form, rowToPayload(row));
  await loadSchedules();
  dialogVisible.value = true;
}

async function handleMerchantChange() {
  const merchant = merchantOptions.value.find((item) => item.id === form.merchantId);
  if (merchant && !editingId.value) {
    form.commissionRate = merchant.commissionRate || 10;
  }
  form.scheduleId = 0;
  await loadSchedules();
}

function handleScheduleChange() {
  const schedule = scheduleOptions.value.find((item) => item.id === form.scheduleId);
  if (!schedule) return;
  form.videoTitle = schedule.videoTitle;
  form.sourceVideo = schedule.videoTitle;
}

async function submit() {
  await formRef.value?.validate();
  if ((form.redeemedAmount || 0) < 0) {
    ElMessage.warning('核销额不能小于 0');
    return;
  }
  submitting.value = true;
  try {
    if (editingId.value) {
      await updateSettlementOrder(editingId.value, form);
      ElMessage.success('分成订单已更新');
    } else {
      await createSettlementOrder(form);
      ElMessage.success('分成订单已创建');
    }
    dialogVisible.value = false;
    await loadList();
  } finally {
    submitting.value = false;
  }
}

async function generateOrders() {
  generating.value = true;
  try {
    const result = await generateSettlementOrders({
      merchantId: routeMerchantId.value,
    });
    ElMessage.success(`已生成 ${result.created} 条待核对分成订单`);
    await loadList();
  } finally {
    generating.value = false;
  }
}

async function confirmOrder(row: SettlementOrder) {
  await updateSettlementOrderStatus(row.id, { status: '已确认' });
  ElMessage.success('分成订单已确认');
  await loadList();
}

async function markPaid(row: SettlementOrder) {
  await updateSettlementOrderStatus(row.id, { status: '已结算' });
  ElMessage.success('分成订单已结算');
  await loadList();
}

async function removeOrder(row: SettlementOrder) {
  await ElMessageBox.confirm(
    `确认删除「${row.merchantName}」这条分成订单？`,
    '删除分成订单',
    { type: 'warning' },
  );
  await deleteSettlementOrder(row.id);
  ElMessage.success('分成订单已删除');
  await loadList();
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

watch(
  () => route.query,
  () => {
    applyRouteQuery();
    pagination.page = 1;
    loadList();
  },
);

onMounted(async () => {
  applyRouteQuery();
  await loadMerchants();
  await loadList();
});
</script>

<template>
  <div class="p-4">
    <el-card>
      <template #header>
        <div class="page-head">
          <div>
            <div class="page-title">分成订单</div>
            <div class="page-desc">
              当前先按“发布视频 / 复盘核销额”生成待核对订单；后续接抖音来客核销单后再自动归因。
            </div>
          </div>
          <div class="page-actions">
            <el-button @click="showAll">显示全部</el-button>
            <el-button @click="openCreate">新增订单</el-button>
            <el-button type="primary" :loading="generating" @click="generateOrders">
              从发布/复盘生成
            </el-button>
          </div>
        </div>
      </template>

      <el-alert
        v-if="currentMerchantName"
        class="mb-4"
        type="info"
        :closable="false"
        show-icon
      >
        当前处理商家：{{ currentMerchantName }}。生成订单时只会处理该商家的已发布/已复盘视频。
      </el-alert>

      <el-form :inline="true" class="search-form">
        <el-form-item label="关键词">
          <el-input
            v-model="keyword"
            clearable
            placeholder="商家 / 视频 / 统计周期 / 备注"
            style="width: 320px"
            @keyup.enter="search"
          />
        </el-form-item>
        <el-form-item label="状态">
          <el-select
            v-model="status"
            clearable
            placeholder="全部状态"
            style="width: 160px"
          >
            <el-option
              v-for="item in statuses"
              :key="item"
              :label="item"
              :value="item"
            />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="search">搜索</el-button>
        </el-form-item>
      </el-form>

      <el-table v-loading="loading" :data="list" border stripe>
        <el-table-column prop="merchantName" label="商家" min-width="160" />
        <el-table-column prop="sourceVideo" label="来源视频" min-width="240">
          <template #default="{ row }">
            <div>{{ row.sourceVideo || row.videoTitle || '-' }}</div>
            <div class="table-sub" v-if="row.scheduleId">排期ID：{{ row.scheduleId }}</div>
          </template>
        </el-table-column>
        <el-table-column label="统计周期" width="220">
          <template #default="{ row }">
            <div>{{ row.orderWindow || '-' }}</div>
            <div class="table-sub">
              {{ formatDate(row.periodStart) }} / {{ formatDate(row.periodEnd) }}
            </div>
          </template>
        </el-table-column>
        <el-table-column label="核销额" width="120">
          <template #default="{ row }">{{ formatMoney(row.redeemedAmount) }}</template>
        </el-table-column>
        <el-table-column label="分成比" width="100">
          <template #default="{ row }">{{ formatPercent(row.commissionRate) }}</template>
        </el-table-column>
        <el-table-column label="应分成" width="120">
          <template #default="{ row }">{{ formatMoney(row.commission) }}</template>
        </el-table-column>
        <el-table-column prop="status" label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="statusTagType(row.status)">{{ row.status }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="remark" label="备注" min-width="220" />
        <el-table-column label="操作" width="360" fixed="right">
          <template #default="{ row }">
            <el-button size="small" @click="openEdit(row)">编辑</el-button>
            <el-button
              v-if="row.status === '待核对'"
              size="small"
              type="primary"
              @click="confirmOrder(row)"
            >
              确认
            </el-button>
            <el-button
              v-if="row.status !== '已结算'"
              size="small"
              type="success"
              @click="markPaid(row)"
            >
              已结算
            </el-button>
            <el-button size="small" type="danger" @click="removeOrder(row)">
              删除
            </el-button>
          </template>
        </el-table-column>
      </el-table>

      <el-empty v-if="!loading && list.length === 0" description="当前没有分成订单">
        <el-button type="primary" @click="generateOrders">从发布/复盘生成</el-button>
      </el-empty>

      <div class="pagination-wrap">
        <el-pagination
          v-model:current-page="pagination.page"
          v-model:page-size="pagination.size"
          background
          layout="total, sizes, prev, pager, next"
          :page-sizes="[10, 20, 50]"
          :total="pagination.total"
          @current-change="handlePageChange"
          @size-change="handleSizeChange"
        />
      </div>
    </el-card>

    <el-dialog v-model="dialogVisible" :title="dialogTitle" width="880px">
      <el-form ref="formRef" :model="form" :rules="rules" label-width="120px">
        <el-row :gutter="16">
          <el-col :span="12">
            <el-form-item label="商家" prop="merchantId">
              <el-select
                v-model="form.merchantId"
                filterable
                placeholder="选择商家"
                style="width: 100%"
                @change="handleMerchantChange"
              >
                <el-option
                  v-for="merchant in merchantOptions"
                  :key="merchant.id"
                  :label="merchant.name"
                  :value="merchant.id"
                />
              </el-select>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="关联排期">
              <el-select
                v-model="form.scheduleId"
                clearable
                filterable
                placeholder="可选，选择已发布视频"
                style="width: 100%"
                @change="handleScheduleChange"
              >
                <el-option
                  v-for="item in scheduleOptions"
                  :key="item.id"
                  :label="`${item.videoTitle}｜${item.status}`"
                  :value="item.id"
                />
              </el-select>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="视频标题">
              <el-input v-model="form.videoTitle" placeholder="例如：9元单人火锅套餐" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="来源视频">
              <el-input v-model="form.sourceVideo" placeholder="视频标题或链接" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="统计开始">
              <el-date-picker
                v-model="form.periodStart"
                placeholder="选择日期"
                style="width: 100%"
                type="date"
                value-format="YYYY-MM-DD"
              />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="统计结束">
              <el-date-picker
                v-model="form.periodEnd"
                placeholder="选择日期"
                style="width: 100%"
                type="date"
                value-format="YYYY-MM-DD"
              />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="核销额" prop="redeemedAmount">
              <el-input-number
                v-model="form.redeemedAmount"
                :min="0"
                :precision="2"
                :step="100"
                style="width: 100%"
              />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="分成比例">
              <el-input-number
                v-model="form.commissionRate"
                :max="100"
                :min="0"
                :precision="2"
                :step="1"
                style="width: 100%"
              />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="状态">
              <el-select v-model="form.status" style="width: 100%">
                <el-option
                  v-for="item in statuses"
                  :key="item"
                  :label="item"
                  :value="item"
                />
              </el-select>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="应分成预览">
              <el-input :model-value="formatMoney(commissionPreview)" disabled />
            </el-form-item>
          </el-col>
          <el-col :span="24">
            <el-form-item label="统计周期文案">
              <el-input
                v-model="form.orderWindow"
                placeholder="不填时后端按开始/结束日期自动生成"
              />
            </el-form-item>
          </el-col>
          <el-col :span="24">
            <el-form-item label="备注">
              <el-input
                v-model="form.remark"
                maxlength="1000"
                placeholder="核销数据来源、待老板确认事项、结算凭证说明等"
                show-word-limit
                type="textarea"
                :rows="4"
              />
            </el-form-item>
          </el-col>
        </el-row>
      </el-form>

      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="submitting" @click="submit">
          保存订单
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

.table-sub {
  margin-top: 4px;
  color: #64748b;
  font-size: 12px;
}

.pagination-wrap {
  display: flex;
  justify-content: flex-end;
  margin-top: 16px;
}
</style>
