<script setup lang="ts">
import type { FormInstance } from 'element-plus';

import { computed, onMounted, reactive, ref, watch } from 'vue';
import { ElMessage, ElMessageBox } from 'element-plus';
import { useRoute, useRouter } from 'vue-router';

import {
  createFollowUpLog,
  deleteFollowUpLog,
  generateFollowUpSuggestion,
  getFollowUpLogList,
  getMerchantList,
  updateFollowUpLog,
  type Merchant,
  type MerchantFollowUpLog,
  type MerchantFollowUpPayload,
  type MerchantFollowUpSuggestion,
} from '#/api/admin';

defineOptions({ name: 'DurationLogPage' });

const route = useRoute();
const router = useRouter();

const stages = ['沟通中', '方案确认', '已签约', '暂缓', '已流失'];

const loading = ref(false);
const submitting = ref(false);
const suggestionLoading = ref(false);
const dialogVisible = ref(false);
const suggestionVisible = ref(false);
const editingId = ref<number | null>(null);
const formRef = ref<FormInstance>();

const keyword = ref('');
const stage = ref('');
const routeMerchantId = ref<number | undefined>();
const routeMerchantName = ref('');
const pagination = reactive({ page: 1, size: 10, total: 0 });
const list = ref<MerchantFollowUpLog[]>([]);
const merchantOptions = ref<Merchant[]>([]);
const suggestion = ref<MerchantFollowUpSuggestion>({
  actions: [],
  talkScript: '',
});

const form = reactive<MerchantFollowUpPayload>({
  followTime: '',
  latestTalk: '',
  merchantId: 0,
  nextFollowTime: '',
  nextStep: '',
  objection: '',
  owner: '',
  stage: '沟通中',
});

const dialogTitle = computed(() =>
  editingId.value ? '编辑跟进记录' : '新增跟进记录',
);

const currentMerchantName = computed(() => {
  if (routeMerchantName.value) return routeMerchantName.value;
  const current = merchantOptions.value.find(
    (item) => item.id === routeMerchantId.value,
  );
  return current?.name || '';
});

const rules = {
  latestTalk: [{ message: '请输入最近沟通', required: true, trigger: 'blur' }],
  merchantId: [{ message: '请选择商家', required: true, trigger: 'change' }],
};

function applyRouteQuery() {
  const rawMerchantId = Number(route.query.merchantId || 0);
  routeMerchantId.value = rawMerchantId > 0 ? rawMerchantId : undefined;
  routeMerchantName.value = String(route.query.merchantName || '');
  if (routeMerchantName.value) {
    keyword.value = routeMerchantName.value;
  }
}

function formatTime(value?: string) {
  if (!value) return '-';
  return value.replace('T', ' ').slice(0, 16);
}

function stageTagType(rowStage: string) {
  if (rowStage === '已签约') return 'success';
  if (rowStage === '已流失') return 'danger';
  if (rowStage === '暂缓') return 'info';
  return 'warning';
}

function resetForm() {
  Object.assign(form, {
    followTime: '',
    latestTalk: '',
    merchantId: routeMerchantId.value || 0,
    nextFollowTime: '',
    nextStep: '',
    objection: '',
    owner: '',
    stage: '沟通中',
  });
  formRef.value?.clearValidate();
}

function rowToPayload(row: MerchantFollowUpLog): MerchantFollowUpPayload {
  return {
    followTime: formatFormTime(row.followTime),
    latestTalk: row.latestTalk,
    merchantId: row.merchantId,
    nextFollowTime: formatFormTime(row.nextFollowTime),
    nextStep: row.nextStep,
    objection: row.objection,
    owner: row.owner,
    stage: row.stage,
  };
}

function formatFormTime(value?: string) {
  if (!value) return '';
  return value.replace('T', ' ').slice(0, 19);
}

async function loadMerchants() {
  const result = await getMerchantList({ page: 1, size: 100 });
  merchantOptions.value = result.list;
}

async function loadList() {
  loading.value = true;
  try {
    const result = await getFollowUpLogList({
      keyword: keyword.value.trim(),
      merchantId: routeMerchantId.value,
      page: pagination.page,
      size: pagination.size,
      stage: stage.value,
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
  stage.value = '';
  routeMerchantId.value = undefined;
  routeMerchantName.value = '';
  router.replace('/users/duration-logs');
  search();
}

function openCreate() {
  editingId.value = null;
  resetForm();
  dialogVisible.value = true;
}

function openEdit(row: MerchantFollowUpLog) {
  editingId.value = row.id;
  Object.assign(form, rowToPayload(row));
  dialogVisible.value = true;
}

async function submit() {
  await formRef.value?.validate();
  submitting.value = true;
  try {
    if (editingId.value) {
      await updateFollowUpLog(editingId.value, form);
      ElMessage.success('跟进记录已更新');
    } else {
      await createFollowUpLog(form);
      ElMessage.success('跟进记录已创建');
    }
    dialogVisible.value = false;
    await loadList();
  } finally {
    submitting.value = false;
  }
}

async function removeLog(row: MerchantFollowUpLog) {
  await ElMessageBox.confirm(
    `确认删除「${row.merchantName}」这条跟进记录？`,
    '删除跟进记录',
    { type: 'warning' },
  );
  await deleteFollowUpLog(row.id);
  ElMessage.success('跟进记录已删除');
  await loadList();
}

async function openSuggestion(row: MerchantFollowUpLog) {
  suggestionVisible.value = true;
  suggestion.value = { actions: [], talkScript: '' };
  suggestionLoading.value = true;
  try {
    suggestion.value = await generateFollowUpSuggestion(row.id);
  } finally {
    suggestionLoading.value = false;
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
            <div class="page-title">跟进记录</div>
            <div class="page-desc">
              记录谈商家过程中的顾虑、推进动作和下次跟进时间，后面签约后直接进入账号和内容流程。
            </div>
          </div>
          <div class="page-actions">
            <el-button @click="showAll">显示全部</el-button>
            <el-button type="primary" @click="openCreate">新增跟进</el-button>
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
        当前处理商家：{{ currentMerchantName }}。已自动带入筛选条件。
      </el-alert>

      <el-form :inline="true" class="search-form">
        <el-form-item label="关键词">
          <el-input
            v-model="keyword"
            clearable
            placeholder="商家 / 沟通内容 / 异议 / 后续动作 / 负责人"
            style="width: 360px"
            @keyup.enter="search"
          />
        </el-form-item>
        <el-form-item label="阶段">
          <el-select
            v-model="stage"
            clearable
            placeholder="全部阶段"
            style="width: 180px"
          >
            <el-option
              v-for="item in stages"
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
        <el-table-column prop="stage" label="阶段" width="110">
          <template #default="{ row }">
            <el-tag :type="stageTagType(row.stage)">{{ row.stage }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="latestTalk" label="最近沟通" min-width="260" />
        <el-table-column prop="objection" label="关键异议" min-width="220" />
        <el-table-column prop="nextStep" label="后续动作" min-width="220" />
        <el-table-column prop="owner" label="负责人" width="110" />
        <el-table-column label="跟进时间" width="160">
          <template #default="{ row }">{{ formatTime(row.followTime) }}</template>
        </el-table-column>
        <el-table-column label="下次跟进" width="160">
          <template #default="{ row }">{{ formatTime(row.nextFollowTime) }}</template>
        </el-table-column>
        <el-table-column label="操作" width="300" fixed="right">
          <template #default="{ row }">
            <el-button size="small" @click="openSuggestion(row)">生成话术</el-button>
            <el-button size="small" type="primary" @click="openEdit(row)">编辑</el-button>
            <el-button size="small" type="danger" @click="removeLog(row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>

      <el-empty v-if="!loading && list.length === 0" description="当前没有跟进记录">
        <el-button type="primary" @click="openCreate">新增跟进</el-button>
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

    <el-dialog v-model="dialogVisible" :title="dialogTitle" width="860px">
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
            <el-form-item label="阶段">
              <el-select v-model="form.stage" style="width: 100%">
                <el-option
                  v-for="item in stages"
                  :key="item"
                  :label="item"
                  :value="item"
                />
              </el-select>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="负责人">
              <el-input v-model="form.owner" placeholder="例如：张三" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="沟通时间">
              <el-date-picker
                v-model="form.followTime"
                placeholder="选择时间"
                style="width: 100%"
                type="datetime"
                value-format="YYYY-MM-DD HH:mm:ss"
              />
            </el-form-item>
          </el-col>
          <el-col :span="24">
            <el-form-item label="最近沟通" prop="latestTalk">
              <el-input
                v-model="form.latestTalk"
                maxlength="2000"
                placeholder="这次聊了什么，老板态度、价格顾虑、是否愿意授权账号等"
                show-word-limit
                type="textarea"
                :rows="4"
              />
            </el-form-item>
          </el-col>
          <el-col :span="24">
            <el-form-item label="关键异议">
              <el-input
                v-model="form.objection"
                maxlength="1000"
                placeholder="例如：担心低价套餐亏损、担心账号给出去不安全、担心视频没效果"
                show-word-limit
                type="textarea"
                :rows="3"
              />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="后续动作">
              <el-input
                v-model="form.nextStep"
                placeholder="例如：明天核算套餐成本，后天确定授权"
              />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="下次跟进">
              <el-date-picker
                v-model="form.nextFollowTime"
                placeholder="选择时间"
                style="width: 100%"
                type="datetime"
                value-format="YYYY-MM-DD HH:mm:ss"
              />
            </el-form-item>
          </el-col>
        </el-row>
      </el-form>

      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="submitting" @click="submit">
          保存跟进
        </el-button>
      </template>
    </el-dialog>

    <el-dialog v-model="suggestionVisible" title="跟进话术" width="760px">
      <el-skeleton v-if="suggestionLoading" :rows="5" animated />
      <template v-else>
        <el-alert type="success" :closable="false" show-icon>
          {{ suggestion.talkScript }}
        </el-alert>
        <div class="suggestion-title">动作清单</div>
        <el-timeline>
          <el-timeline-item v-for="item in suggestion.actions" :key="item">
            {{ item }}
          </el-timeline-item>
        </el-timeline>
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

.pagination-wrap {
  display: flex;
  justify-content: flex-end;
  margin-top: 16px;
}

.suggestion-title {
  margin: 18px 0 10px;
  color: #0f172a;
  font-weight: 700;
}
</style>
