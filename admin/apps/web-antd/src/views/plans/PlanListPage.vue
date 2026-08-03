<script setup lang="ts">
import type { FormInstance } from 'element-plus';

import { computed, onMounted, reactive, ref, watch } from 'vue';
import { ElMessage, ElMessageBox } from 'element-plus';
import { useRoute, useRouter } from 'vue-router';

import {
  createMerchantPackage,
  deleteMerchantPackage,
  getMerchantList,
  getMerchantPackageList,
  updateMerchantPackage,
  type Merchant,
  type MerchantPackage,
  type MerchantPackagePayload,
} from '#/api/admin';

defineOptions({ name: 'PlanListPage' });

const route = useRoute();
const router = useRouter();

const loading = ref(false);
const submitting = ref(false);
const dialogVisible = ref(false);
const editingId = ref<number | null>(null);
const formRef = ref<FormInstance>();

const keyword = ref('');
const routeMerchantId = ref<number | undefined>();
const routeMerchantName = ref('');
const pagination = reactive({ page: 1, size: 10, total: 0 });
const list = ref<MerchantPackage[]>([]);
const merchantOptions = ref<Merchant[]>([]);

const form = reactive<MerchantPackagePayload>({
  commissionRate: 10,
  costPrice: 0,
  merchantId: 0,
  name: '',
  originalPrice: 0,
  profitGuard: '',
  remark: '',
  sellingPrice: 0,
  status: 1,
  trafficLabel: '',
});

const dialogTitle = computed(() =>
  editingId.value ? '编辑团购套餐' : '新增团购套餐',
);

const currentMerchantName = computed(() => {
  if (routeMerchantName.value) return routeMerchantName.value;
  const current = merchantOptions.value.find((item) => item.id === routeMerchantId.value);
  return current?.name || '';
});

const grossProfitPreview = computed(() =>
  Number(((form.sellingPrice || 0) - (form.costPrice || 0)).toFixed(2)),
);

const estimatedCommissionPreview = computed(() =>
  Number((((form.sellingPrice || 0) * (form.commissionRate || 0)) / 100).toFixed(2)),
);

const netAfterCommissionPreview = computed(() =>
  Number((grossProfitPreview.value - estimatedCommissionPreview.value).toFixed(2)),
);

const marginRatePreview = computed(() => {
  if (!form.sellingPrice) return 0;
  return Number(((grossProfitPreview.value / form.sellingPrice) * 100).toFixed(2));
});

const rules = {
  merchantId: [{ message: '请选择商家', required: true, trigger: 'change' }],
  name: [{ message: '请输入套餐名称', required: true, trigger: 'blur' }],
  sellingPrice: [{ message: '请输入售价', required: true, trigger: 'blur' }],
};

function applyRouteQuery() {
  const rawMerchantId = Number(route.query.merchantId || 0);
  routeMerchantId.value = rawMerchantId > 0 ? rawMerchantId : undefined;
  routeMerchantName.value = String(route.query.merchantName || '');
  keyword.value = routeMerchantName.value;
}

function formatMoney(value: number) {
  return `¥${Number(value || 0).toFixed(2)}`;
}

function formatPercent(value: number) {
  return `${Number(value || 0).toFixed(2)}%`;
}

function getProfitTagType(row: MerchantPackage) {
  if (row.netAfterCommission <= 0) return 'danger';
  if (row.marginRate < 20) return 'warning';
  return 'success';
}

function resetForm() {
  const currentMerchant = merchantOptions.value.find(
    (item) => item.id === routeMerchantId.value,
  );
  Object.assign(form, {
    commissionRate: currentMerchant?.commissionRate || 10,
    costPrice: 0,
    merchantId: routeMerchantId.value || 0,
    name: '',
    originalPrice: 0,
    profitGuard: '',
    remark: '',
    sellingPrice: 0,
    status: 1,
    trafficLabel: '',
  });
  formRef.value?.clearValidate();
}

function packageToPayload(row: MerchantPackage): MerchantPackagePayload {
  return {
    commissionRate: row.commissionRate,
    costPrice: row.costPrice,
    merchantId: row.merchantId,
    name: row.name,
    originalPrice: row.originalPrice,
    profitGuard: row.profitGuard,
    remark: row.remark,
    sellingPrice: row.sellingPrice,
    status: row.status,
    trafficLabel: row.trafficLabel,
  };
}

async function loadMerchants() {
  const result = await getMerchantList({ page: 1, size: 100 });
  merchantOptions.value = result.list;
}

async function loadList() {
  loading.value = true;
  try {
    const result = await getMerchantPackageList({
      keyword: keyword.value.trim(),
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

function search() {
  pagination.page = 1;
  loadList();
}

function showAll() {
  keyword.value = '';
  routeMerchantId.value = undefined;
  routeMerchantName.value = '';
  router.replace('/plans/list');
  search();
}

function openCreate() {
  editingId.value = null;
  resetForm();
  dialogVisible.value = true;
}

function openEdit(row: MerchantPackage) {
  editingId.value = row.id;
  Object.assign(form, packageToPayload(row));
  dialogVisible.value = true;
}

async function submit() {
  await formRef.value?.validate();
  if (!form.merchantId) {
    ElMessage.warning('请选择商家');
    return;
  }
  if (form.sellingPrice <= 0) {
    ElMessage.warning('售价必须大于 0');
    return;
  }
  if ((form.costPrice || 0) > form.sellingPrice) {
    ElMessage.warning('成本不能高于售价');
    return;
  }

  submitting.value = true;
  try {
    if (editingId.value) {
      await updateMerchantPackage(editingId.value, form);
      ElMessage.success('套餐已更新');
    } else {
      await createMerchantPackage(form);
      ElMessage.success('套餐已创建');
    }
    dialogVisible.value = false;
    await loadList();
  } finally {
    submitting.value = false;
  }
}

async function toggleStatus(row: MerchantPackage) {
  const nextStatus = row.status === 1 ? 0 : 1;
  await updateMerchantPackage(row.id, {
    ...packageToPayload(row),
    status: nextStatus,
  });
  ElMessage.success(nextStatus === 1 ? '套餐已启用' : '套餐已停用');
  await loadList();
}

async function removePackage(row: MerchantPackage) {
  await ElMessageBox.confirm(
    `确认删除套餐「${row.name}」？删除后列表不再保留这条套餐。`,
    '删除套餐',
    { type: 'warning' },
  );
  await deleteMerchantPackage(row.id);
  ElMessage.success('套餐已删除');
  await loadList();
}

function goTopics(row: MerchantPackage) {
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
            <div class="page-title">团购套餐</div>
            <div class="page-desc">
              先把商家要卖的套餐、售价、成本、提点和利润空间算清楚，再进入选题。
            </div>
          </div>
          <div class="page-actions">
            <el-button @click="showAll">显示全部</el-button>
            <el-button type="primary" @click="openCreate">新增套餐</el-button>
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
        当前处理商家：{{ currentMerchantName }}。已自动带入筛选条件，商家ID：{{ routeMerchantId || '-' }}。
      </el-alert>

      <el-form :inline="true" class="search-form">
        <el-form-item label="关键词">
          <el-input
            v-model="keyword"
            clearable
            placeholder="商家 / 套餐 / 投放定位 / 利润建议"
            style="width: 320px"
            @keyup.enter="search"
          />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="search">搜索</el-button>
        </el-form-item>
      </el-form>

      <el-table v-loading="loading" :data="list" border stripe>
        <el-table-column prop="merchantName" label="商家" min-width="160" />
        <el-table-column prop="name" label="套餐名称" min-width="180" />
        <el-table-column label="价格" width="170">
          <template #default="{ row }">
            <div>售价：{{ formatMoney(row.sellingPrice) }}</div>
            <div class="table-sub">成本：{{ formatMoney(row.costPrice) }}</div>
          </template>
        </el-table-column>
        <el-table-column label="利润" width="190">
          <template #default="{ row }">
            <el-tag :type="getProfitTagType(row)">
              毛利 {{ formatPercent(row.marginRate) }}
            </el-tag>
            <div class="table-sub">
              扣提点后：{{ formatMoney(row.netAfterCommission) }}
            </div>
          </template>
        </el-table-column>
        <el-table-column prop="trafficLabel" label="投放定位" min-width="160" />
        <el-table-column prop="profitGuard" label="利润保护建议" min-width="220" />
        <el-table-column label="状态" width="90">
          <template #default="{ row }">
            <el-tag :type="row.status === 1 ? 'success' : 'info'">
              {{ row.status === 1 ? '启用' : '停用' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="340" fixed="right">
          <template #default="{ row }">
            <el-button size="small" @click="openEdit(row)">编辑</el-button>
            <el-button size="small" @click="toggleStatus(row)">
              {{ row.status === 1 ? '停用' : '启用' }}
            </el-button>
            <el-button size="small" type="primary" @click="goTopics(row)">
              生成选题
            </el-button>
            <el-button size="small" type="danger" @click="removePackage(row)">
              删除
            </el-button>
          </template>
        </el-table-column>
      </el-table>

      <el-empty v-if="!loading && list.length === 0" description="当前没有套餐">
        <el-button type="primary" @click="openCreate">新增套餐</el-button>
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

    <el-dialog v-model="dialogVisible" :title="dialogTitle" width="820px">
      <el-form ref="formRef" :model="form" :rules="rules" label-width="120px">
        <el-row :gutter="16">
          <el-col :span="12">
            <el-form-item label="商家" prop="merchantId">
              <el-select
                v-model="form.merchantId"
                filterable
                placeholder="选择商家"
                style="width: 100%"
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
            <el-form-item label="套餐名称" prop="name">
              <el-input v-model="form.name" placeholder="例如：双人鲜牛火锅套餐" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="门市价">
              <el-input-number
                v-model="form.originalPrice"
                :min="0"
                :precision="2"
                :step="1"
                style="width: 100%"
              />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="团购售价" prop="sellingPrice">
              <el-input-number
                v-model="form.sellingPrice"
                :min="0"
                :precision="2"
                :step="1"
                style="width: 100%"
              />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="商家成本">
              <el-input-number
                v-model="form.costPrice"
                :min="0"
                :precision="2"
                :step="1"
                style="width: 100%"
              />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="提点比例">
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
                <el-option label="启用" :value="1" />
                <el-option label="停用" :value="0" />
              </el-select>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="投放定位">
              <el-input v-model="form.trafficLabel" placeholder="门店引流主推 / 同城种草" />
            </el-form-item>
          </el-col>
          <el-col :span="24">
            <div class="profit-preview">
              <div>
                <span>毛利</span>
                <strong>{{ formatMoney(grossProfitPreview) }}</strong>
              </div>
              <div>
                <span>毛利率</span>
                <strong>{{ formatPercent(marginRatePreview) }}</strong>
              </div>
              <div>
                <span>预估提点</span>
                <strong>{{ formatMoney(estimatedCommissionPreview) }}</strong>
              </div>
              <div>
                <span>扣提点后</span>
                <strong>{{ formatMoney(netAfterCommissionPreview) }}</strong>
              </div>
            </div>
          </el-col>
          <el-col :span="24">
            <el-form-item label="利润保护建议">
              <el-input
                v-model="form.profitGuard"
                placeholder="不填时后端会按售价、成本和提点自动生成"
              />
            </el-form-item>
          </el-col>
          <el-col :span="24">
            <el-form-item label="备注">
              <el-input
                v-model="form.remark"
                maxlength="1000"
                placeholder="套餐包含内容、不可降价边界、老板特殊要求等"
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
          保存套餐
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

.profit-preview {
  display: grid;
  grid-template-columns: repeat(4, minmax(0, 1fr));
  gap: 12px;
  margin: 0 0 18px 120px;
}

.profit-preview div {
  border: 1px solid #dbeafe;
  border-radius: 12px;
  background: #eff6ff;
  padding: 12px;
}

.profit-preview span {
  display: block;
  color: #64748b;
  font-size: 12px;
}

.profit-preview strong {
  display: block;
  margin-top: 6px;
  color: #1d4ed8;
  font-size: 16px;
}
</style>
