<script setup lang="ts">
import type { FormInstance } from 'element-plus';

import { computed, onMounted, reactive, ref } from 'vue';
import { useRouter } from 'vue-router';

import { useUserStore } from '@vben/stores';

import { ElMessage, ElMessageBox } from 'element-plus';

import {
  createMerchant,
  deleteMerchant,
  getMerchantList,
  updateMerchant,
  type Merchant,
  type MerchantPayload,
} from '#/api/admin';

defineOptions({ name: 'UserListPage' });

const router = useRouter();
const userStore = useUserStore();

const loading = ref(false);
const dialogVisible = ref(false);
const submitting = ref(false);
const editingId = ref<number | null>(null);
const formRef = ref<FormInstance>();

const searchForm = reactive({ keyword: '' });
const pagination = reactive({ page: 1, size: 10, total: 0 });
const list = ref<Merchant[]>([]);
const canDeleteMerchant = computed(() =>
  (userStore.userInfo?.roles || []).includes('super_admin'),
);

const form = reactive<MerchantPayload>({
  address: '',
  city: '',
  commissionRate: 10,
  contactName: '',
  contactPhone: '',
  cooperationType: '成交提点',
  douyinAccount: '',
  douyinLaikeAccount: '',
  industry: '',
  name: '',
  remark: '',
  stage: '已建档',
  status: 1,
});

const dialogTitle = computed(() =>
  editingId.value ? '编辑商家档案' : '新建商家档案',
);

const rules = {
  name: [{ message: '请输入商家名称', required: true, trigger: 'blur' }],
};

function resetForm() {
  Object.assign(form, {
    address: '',
    city: '',
    commissionRate: 10,
    contactName: '',
    contactPhone: '',
    cooperationType: '成交提点',
    douyinAccount: '',
    douyinLaikeAccount: '',
    industry: '',
    name: '',
    remark: '',
    stage: '已建档',
    status: 1,
  });
  formRef.value?.clearValidate();
}

function merchantToPayload(row: Merchant): MerchantPayload {
  return {
    address: row.address,
    city: row.city,
    commissionRate: row.commissionRate,
    contactName: row.contactName,
    contactPhone: row.contactPhone,
    cooperationType: row.cooperationType,
    douyinAccount: row.douyinAccount,
    douyinLaikeAccount: row.douyinLaikeAccount,
    industry: row.industry,
    name: row.name,
    remark: row.remark,
    stage: row.stage,
    status: row.status,
  };
}

interface DirectAction {
  desc: string;
  route?: string;
  text: string;
}

function getDirectAction(row: Merchant): DirectAction {
  if (!row.name || !row.industry || !row.city || !row.contactName) {
    return {
      desc: '缺商家名称、行业、城市或联系人，先把建档信息补完整。',
      text: '补齐基础档案',
    };
  }
  if (!row.cooperationType || row.commissionRate <= 0) {
    return {
      desc: '先确认合作方式和成交提点比例，避免后面套餐和内容方向跑偏。',
      text: '确认合作规则',
    };
  }
  if (!row.douyinAccount && !row.douyinLaikeAccount) {
    return {
      desc: '先记录抖音账号或抖音来客账号，后面诊断和运营才有对象。',
      route: `/users/devices?merchantId=${row.id}&merchantName=${encodeURIComponent(row.name)}`,
      text: '记录账号授权',
    };
  }
  return {
    desc: '基础资料已够，先把要卖的团购套餐建起来。',
    route: `/plans/list?merchantId=${row.id}&merchantName=${encodeURIComponent(row.name)}`,
    text: '去建团购套餐',
  };
}

function getDirectActionText(row: Merchant) {
  return getDirectAction(row).text;
}

function getDirectActionDesc(row: Merchant) {
  return getDirectAction(row).desc;
}

async function loadList() {
  loading.value = true;
  try {
    const result = await getMerchantList({
      keyword: searchForm.keyword.trim(),
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

function resetSearch() {
  searchForm.keyword = '';
  search();
}

function openCreate() {
  editingId.value = null;
  resetForm();
  dialogVisible.value = true;
}

function openEdit(row: Merchant) {
  editingId.value = row.id;
  Object.assign(form, merchantToPayload(row));
  dialogVisible.value = true;
}

function handleDirectAction(row: Merchant) {
  const action = getDirectAction(row);
  if (!action.route) {
    openEdit(row);
    return;
  }
  router.push(action.route);
}

async function submit() {
  await formRef.value?.validate();
  submitting.value = true;
  try {
    if (editingId.value) {
      await updateMerchant(editingId.value, form);
      ElMessage.success('商家档案已更新');
    } else {
      await createMerchant(form);
      ElMessage.success('商家档案已创建');
    }
    dialogVisible.value = false;
    await loadList();
  } finally {
    submitting.value = false;
  }
}

async function toggleStatus(row: Merchant) {
  const nextStatus = row.status === 1 ? 0 : 1;
  await updateMerchant(row.id, {
    ...merchantToPayload(row),
    status: nextStatus,
  });
  ElMessage.success(nextStatus === 1 ? '商家已启用' : '商家已停用');
  await loadList();
}

async function handleDelete(row: Merchant) {
  if (!canDeleteMerchant.value) {
    ElMessage.error('只有超级管理员可以删除商家');
    return;
  }

  try {
    await ElMessageBox.confirm(
      `确定删除「${row.name}」吗？确认后会同时物理删除该商家的套餐、账号授权、账号诊断、对标分析、选题、文案、分镜、拍摄任务、发布排期、数据复盘等全部关联数据，无法恢复。`,
      '删除商家及全部关联数据',
      {
        cancelButtonText: '取消',
        confirmButtonText: '确认删除',
        type: 'warning',
      },
    );

    const result = await deleteMerchant(row.id);
    const totalDeleted = Object.values(result.deleted || {}).reduce(
      (sum, value) => sum + Number(value || 0),
      0,
    );
    ElMessage.success(
      `已删除「${result.name || row.name}」及 ${totalDeleted} 条关联数据`,
    );
    await loadList();
  } catch (error) {
    if (error === 'cancel' || error === 'close') {
      return;
    }
    const maybeError = error as {
      error?: string;
      message?: string;
      response?: { data?: { error?: string; message?: string } };
    };
    ElMessage.error(
      maybeError?.response?.data?.error ||
        maybeError?.response?.data?.message ||
        maybeError?.error ||
        maybeError?.message ||
        '删除商家失败',
    );
  }
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

onMounted(loadList);
</script>

<template>
  <div class="p-4">
    <el-card>
      <template #header>
        <div class="page-head">
          <div>
            <div class="page-title">商家建档</div>
            <div class="page-desc">
              这里只做商家基础建档和动作引导。运营人员根据系统给出的动作直接跳到该处理的页面。
            </div>
          </div>
          <el-button type="primary" @click="openCreate">新建商家档案</el-button>
        </div>
      </template>

      <el-alert class="mb-4" type="info" :closable="false" show-icon>
        系统根据商家资料完整度给出当前动作，点按钮直接进入对应页面。
      </el-alert>

      <el-form :inline="true" :model="searchForm" class="search-form">
        <el-form-item label="关键词">
          <el-input
            v-model="searchForm.keyword"
            clearable
            placeholder="商家名 / 城市 / 行业 / 联系人"
            style="width: 280px"
            @keyup.enter="search"
          />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="search">搜索</el-button>
          <el-button @click="resetSearch">重置</el-button>
        </el-form-item>
      </el-form>

      <el-table v-loading="loading" :data="list" border stripe>
        <el-table-column label="商家" min-width="220">
          <template #default="{ row }">
            <div class="merchant-name">{{ row.name }}</div>
            <div class="merchant-sub">
              {{ row.industry || '未填行业' }} · {{ row.city || '未填城市' }}
            </div>
          </template>
        </el-table-column>
        <el-table-column label="联系人" width="160">
          <template #default="{ row }">
            <div>{{ row.contactName || '-' }}</div>
            <div class="merchant-sub">{{ row.contactPhone || '-' }}</div>
          </template>
        </el-table-column>
        <el-table-column label="合作方式" width="150">
          <template #default="{ row }">
            <div>{{ row.cooperationType }}</div>
            <div class="merchant-sub">提点 {{ row.commissionRate }}%</div>
          </template>
        </el-table-column>
        <el-table-column label="当前阶段" width="130">
          <template #default="{ row }">
            <el-tag>{{ row.stage }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column label="系统建议" min-width="220">
          <template #default="{ row }">
            <strong>{{ getDirectActionText(row) }}</strong>
            <div class="merchant-sub">{{ getDirectActionDesc(row) }}</div>
          </template>
        </el-table-column>
        <el-table-column label="状态" width="90">
          <template #default="{ row }">
            <el-tag :type="row.status === 1 ? 'success' : 'info'">
              {{ row.status === 1 ? '启用' : '停用' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="360" fixed="right">
          <template #default="{ row }">
            <el-button size="small" @click="openEdit(row)">编辑档案</el-button>
            <el-button size="small" @click="toggleStatus(row)">
              {{ row.status === 1 ? '停用' : '启用' }}
            </el-button>
            <el-button size="small" type="primary" @click="handleDirectAction(row)">
              {{ getDirectActionText(row) }}
            </el-button>
            <el-button
              v-if="canDeleteMerchant"
              size="small"
              type="danger"
              @click="handleDelete(row)"
            >
              删除
            </el-button>
          </template>
        </el-table-column>
      </el-table>

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
            <el-form-item label="商家名称" prop="name">
              <el-input v-model="form.name" placeholder="例如：贵阳某某火锅店" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="行业">
              <el-input v-model="form.industry" placeholder="餐饮 / 教培 / 门店" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="城市">
              <el-input v-model="form.city" placeholder="例如：贵阳" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="当前阶段">
              <el-select v-model="form.stage" style="width: 100%">
                <el-option label="已建档" value="已建档" />
                <el-option label="已签约" value="已签约" />
                <el-option label="账号诊断中" value="账号诊断中" />
                <el-option label="内容运营中" value="内容运营中" />
                <el-option label="暂停合作" value="暂停合作" />
              </el-select>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="联系人">
              <el-input v-model="form.contactName" placeholder="老板 / 店长" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="联系电话">
              <el-input v-model="form.contactPhone" placeholder="手机号 / 微信号" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="抖音账号">
              <el-input v-model="form.douyinAccount" placeholder="抖音号 / 昵称" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="抖音来客">
              <el-input
                v-model="form.douyinLaikeAccount"
                placeholder="来客后台账号或备注"
              />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="合作方式">
              <el-select v-model="form.cooperationType" style="width: 100%">
                <el-option label="成交提点" value="成交提点" />
                <el-option label="固定服务费" value="固定服务费" />
                <el-option label="服务费 + 提点" value="服务费 + 提点" />
                <el-option label="待确认" value="待确认" />
              </el-select>
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
          <el-col :span="24">
            <el-form-item label="门店地址">
              <el-input v-model="form.address" placeholder="详细地址" />
            </el-form-item>
          </el-col>
          <el-col :span="24">
            <el-form-item label="备注">
              <el-input
                v-model="form.remark"
                maxlength="1000"
                placeholder="合作边界、老板关注点、当前动作等"
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
          保存档案
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
  max-width: 780px;
  margin-top: 6px;
  color: #64748b;
  font-size: 13px;
  line-height: 1.6;
}

.search-form {
  margin-bottom: 12px;
}

.merchant-name {
  color: #0f172a;
  font-weight: 600;
}

.merchant-sub {
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
