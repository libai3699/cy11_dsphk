<script setup lang="ts">
import type { FormInstance } from 'element-plus';

import { computed, onMounted, reactive, ref, watch } from 'vue';
import { ElMessage, ElMessageBox } from 'element-plus';
import { useRoute, useRouter } from 'vue-router';

import {
  createMerchantAccountAuth,
  deleteMerchantAccountAuth,
  getMerchantAccountAuthList,
  getMerchantList,
  updateMerchantAccountAuth,
  type Merchant,
  type MerchantAccountAuth,
  type MerchantAccountAuthPayload,
} from '#/api/admin';

defineOptions({ name: 'UserDevicesPage' });

const route = useRoute();
const router = useRouter();

const loading = ref(false);
const submitting = ref(false);
const dialogVisible = ref(false);
const editingId = ref<number | null>(null);
const formRef = ref<FormInstance>();

const keyword = ref('');
const statusFilter = ref('');
const routeMerchantId = ref<number | undefined>();
const routeMerchantName = ref('');
const pagination = reactive({ page: 1, size: 10, total: 0 });
const list = ref<MerchantAccountAuth[]>([]);
const merchantOptions = ref<Merchant[]>([]);

const platformOptions = ['抖音号', '抖音来客', '创作服务平台', '其他'];
const methodOptions = ['验证码代登', '子账号协作', '人工登录', '老板自行发布'];
const statusOptions = ['待授权', '已授权', '待续期', '已失效'];

const form = reactive<MerchantAccountAuthPayload>({
  accountName: '',
  accountUid: '',
  authMethod: '验证码代登',
  authStatus: '待授权',
  lastLoginAt: '',
  merchantId: 0,
  platform: '抖音号',
  remark: '',
});

const dialogTitle = computed(() =>
  editingId.value ? '编辑账号授权' : '新增账号授权',
);

const currentMerchantName = computed(() => {
  if (routeMerchantName.value) return routeMerchantName.value;
  const current = merchantOptions.value.find((item) => item.id === routeMerchantId.value);
  return current?.name || '';
});

const rules = {
  merchantId: [{ message: '请选择商家', required: true, trigger: 'change' }],
};

function applyRouteQuery() {
  const rawMerchantId = Number(route.query.merchantId || 0);
  routeMerchantId.value = rawMerchantId > 0 ? rawMerchantId : undefined;
  routeMerchantName.value = String(route.query.merchantName || '');
  keyword.value = routeMerchantName.value;
}

function getStatusTagType(status: string) {
  if (status === '已授权') return 'success';
  if (status === '待续期') return 'warning';
  if (status === '已失效') return 'danger';
  return 'info';
}

function resetForm() {
  Object.assign(form, {
    accountName: '',
    accountUid: '',
    authMethod: '验证码代登',
    authStatus: '待授权',
    lastLoginAt: '',
    merchantId: routeMerchantId.value || 0,
    platform: '抖音号',
    remark: '',
  });
  formRef.value?.clearValidate();
}

function authToPayload(row: MerchantAccountAuth): MerchantAccountAuthPayload {
  return {
    accountName: row.accountName,
    accountUid: row.accountUid,
    authMethod: row.authMethod,
    authStatus: row.authStatus,
    lastLoginAt: row.lastLoginAt,
    merchantId: row.merchantId,
    platform: row.platform,
    remark: row.remark,
  };
}

async function loadMerchants() {
  const result = await getMerchantList({ page: 1, size: 100 });
  merchantOptions.value = result.list;
}

async function loadList() {
  loading.value = true;
  try {
    const result = await getMerchantAccountAuthList({
      keyword: keyword.value.trim(),
      merchantId: routeMerchantId.value,
      page: pagination.page,
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
  router.replace('/users/devices');
  search();
}

function openCreate() {
  editingId.value = null;
  resetForm();
  dialogVisible.value = true;
}

function openEdit(row: MerchantAccountAuth) {
  editingId.value = row.id;
  Object.assign(form, authToPayload(row));
  dialogVisible.value = true;
}

async function submit() {
  await formRef.value?.validate();
  if (!form.merchantId) {
    ElMessage.warning('请选择商家');
    return;
  }

  submitting.value = true;
  try {
    if (editingId.value) {
      await updateMerchantAccountAuth(editingId.value, form);
      ElMessage.success('授权记录已更新');
    } else {
      await createMerchantAccountAuth(form);
      ElMessage.success('授权记录已创建');
    }
    dialogVisible.value = false;
    await loadList();
  } finally {
    submitting.value = false;
  }
}

async function updateStatus(row: MerchantAccountAuth, authStatus: string) {
  await updateMerchantAccountAuth(row.id, {
    ...authToPayload(row),
    authStatus,
  });
  ElMessage.success(`授权状态已更新为${authStatus}`);
  await loadList();
}

async function removeAuth(row: MerchantAccountAuth) {
  await ElMessageBox.confirm(
    `确认删除「${row.merchantName}」的 ${row.platform} 授权记录？`,
    '删除授权记录',
    { type: 'warning' },
  );
  await deleteMerchantAccountAuth(row.id);
  ElMessage.success('授权记录已删除');
  await loadList();
}

function startDiagnosis(row: MerchantAccountAuth) {
  router.push({
    path: '/users/account-diagnosis',
    query: {
      accountAuthId: row.id,
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
            <div class="page-title">账号授权</div>
            <div class="page-desc">
              记录抖音号、抖音来客和创作服务平台的协作方式，不保存密码。
            </div>
          </div>
          <div class="page-actions">
            <el-button @click="showAll">显示全部</el-button>
            <el-button type="primary" @click="openCreate">新增授权</el-button>
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

      <el-alert class="mb-4" type="warning" :closable="false" show-icon>
        安全边界：只记录协作授权、验证码代登、人工登录和老板确认发布，不托管密码。
      </el-alert>

      <el-form :inline="true" class="search-form">
        <el-form-item label="关键词">
          <el-input
            v-model="keyword"
            clearable
            placeholder="商家 / 平台 / 授权方式 / 账号 / 状态"
            style="width: 320px"
            @keyup.enter="search"
          />
        </el-form-item>
        <el-form-item label="状态">
          <el-select
            v-model="statusFilter"
            clearable
            placeholder="全部状态"
            style="width: 140px"
          >
            <el-option
              v-for="status in statusOptions"
              :key="status"
              :label="status"
              :value="status"
            />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="search">搜索</el-button>
        </el-form-item>
      </el-form>

      <el-table v-loading="loading" :data="list" border stripe>
        <el-table-column prop="merchantName" label="商家" min-width="160" />
        <el-table-column prop="platform" label="平台" width="130" />
        <el-table-column prop="authMethod" label="授权方式" min-width="150" />
        <el-table-column label="账号信息" min-width="210">
          <template #default="{ row }">
            <div>{{ row.accountName || '-' }}</div>
            <div class="table-sub">{{ row.accountUid || '-' }}</div>
          </template>
        </el-table-column>
        <el-table-column label="状态" width="110">
          <template #default="{ row }">
            <el-tag :type="getStatusTagType(row.authStatus)">
              {{ row.authStatus }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="lastLoginAt" label="最近登录" width="170">
          <template #default="{ row }">{{ row.lastLoginAt || '-' }}</template>
        </el-table-column>
        <el-table-column prop="remark" label="备注" min-width="220" />
        <el-table-column label="操作" width="390" fixed="right">
          <template #default="{ row }">
            <el-button size="small" @click="openEdit(row)">编辑</el-button>
            <el-dropdown class="ml-2" trigger="click">
              <el-button size="small">改状态</el-button>
              <template #dropdown>
                <el-dropdown-menu>
                  <el-dropdown-item
                    v-for="status in statusOptions"
                    :key="status"
                    @click="updateStatus(row, status)"
                  >
                    {{ status }}
                  </el-dropdown-item>
                </el-dropdown-menu>
              </template>
            </el-dropdown>
            <el-button size="small" type="primary" @click="startDiagnosis(row)">
              账号诊断
            </el-button>
            <el-button size="small" type="danger" @click="removeAuth(row)">
              删除
            </el-button>
          </template>
        </el-table-column>
      </el-table>

      <el-empty v-if="!loading && list.length === 0" description="当前没有授权记录">
        <el-button type="primary" @click="openCreate">新增授权</el-button>
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

    <el-dialog v-model="dialogVisible" :title="dialogTitle" width="760px">
      <el-form ref="formRef" :model="form" :rules="rules" label-width="110px">
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
            <el-form-item label="平台">
              <el-select v-model="form.platform" style="width: 100%">
                <el-option
                  v-for="platform in platformOptions"
                  :key="platform"
                  :label="platform"
                  :value="platform"
                />
              </el-select>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="授权方式">
              <el-select v-model="form.authMethod" style="width: 100%">
                <el-option
                  v-for="method in methodOptions"
                  :key="method"
                  :label="method"
                  :value="method"
                />
              </el-select>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="授权状态">
              <el-select v-model="form.authStatus" style="width: 100%">
                <el-option
                  v-for="status in statusOptions"
                  :key="status"
                  :label="status"
                  :value="status"
                />
              </el-select>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="账号名称">
              <el-input v-model="form.accountName" placeholder="抖音昵称 / 来客账号备注" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="账号标识">
              <el-input v-model="form.accountUid" placeholder="抖音号 / 手机号后四位 / 子账号" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="最近登录">
              <el-input v-model="form.lastLoginAt" placeholder="例如：2026-08-01 18:30" />
            </el-form-item>
          </el-col>
          <el-col :span="24">
            <el-form-item label="备注">
              <el-input
                v-model="form.remark"
                maxlength="1000"
                placeholder="登录边界、发布是否需老板确认、验证码找谁要等"
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
          保存授权
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
