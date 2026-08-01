<script setup lang="ts">
import { onMounted, reactive, ref } from 'vue';
import { ElMessage } from 'element-plus';
import { getOrderList, createOrder, type Order } from '#/api/admin/orders';
import { getPlanList, type Plan } from '#/api/admin/plans';

const loading = ref(false);
const list = ref<Order[]>([]);
const total = ref(0);
const page = reactive({ current: 1, size: 20 });
const dialogVisible = ref(false);
const plans = ref<Plan[]>([]);
const form = reactive({
  billing_cycle: 'month',
  plan_id: undefined as number | undefined,
  remark: '',
  user_id: undefined as number | undefined,
});
const cycleOptions = [
  { label: '月付', value: 'month' },
  { label: '季付', value: 'quarter' },
  { label: '半年付', value: 'half_year' },
  { label: '年付', value: 'year' },
];

async function load() {
  loading.value = true;
  try {
    const res = await getOrderList({ page: page.current, size: page.size });
    list.value = res?.list ?? []; total.value = res?.total ?? 0;
  } finally { loading.value = false; }
}

async function openCreate() {
  plans.value = await getPlanList() ?? [];
  Object.assign(form, { billing_cycle: 'month', user_id: undefined, plan_id: undefined, remark: '' });
  dialogVisible.value = true;
}

async function handleSubmit() {
  if (!form.user_id || !form.plan_id) { ElMessage.warning('请填写用户ID和套餐'); return; }
  await createOrder({
    billing_cycle: form.billing_cycle as 'half_year' | 'month' | 'quarter' | 'year',
    user_id: form.user_id,
    plan_id: form.plan_id,
    remark: form.remark,
  });
  ElMessage.success('开通成功'); dialogVisible.value = false; load();
}

onMounted(load);
</script>

<template>
  <div class="p-4">
    <el-card>
      <div class="mb-4 flex justify-end">
        <el-button type="primary" @click="openCreate">手动开通</el-button>
      </div>
      <el-table :data="list" v-loading="loading" border stripe>
        <el-table-column prop="id" label="ID" width="70" />
        <el-table-column prop="user_id" label="用户ID" width="90" />
        <el-table-column prop="plan_name" label="套餐" />
        <el-table-column prop="plan_price" label="价格" width="90" />
        <el-table-column label="流量" width="90">
          <template #default="{ row }">{{ row.traffic_gb ? `${row.traffic_gb} GB` : '无限' }}</template>
        </el-table-column>
        <el-table-column prop="duration_days" label="有效期(天)" width="100" />
        <el-table-column prop="started_at" label="生效时间" width="170" />
        <el-table-column prop="expired_at" label="到期时间" width="170" />
        <el-table-column prop="pay_method" label="方式" width="80" />
        <el-table-column prop="remark" label="备注" />
      </el-table>
      <div class="mt-4 flex justify-end">
        <el-pagination v-model:current-page="page.current" v-model:page-size="page.size" :total="total" layout="total, prev, pager, next" @change="load" />
      </div>
    </el-card>

    <el-dialog v-model="dialogVisible" title="手动开通套餐" width="480px">
      <el-form :model="form" label-width="100px">
        <el-form-item label="用户ID" required><el-input-number v-model="form.user_id" :min="1" style="width:100%" /></el-form-item>
        <el-form-item label="套餐" required>
          <el-select v-model="form.plan_id" style="width:100%">
            <el-option v-for="p in plans" :key="p.id" :value="p.id" :label="`${p.name} - ¥${p.price} / ${p.duration_days}天`" />
          </el-select>
        </el-form-item>
        <el-form-item label="开通周期" required>
          <el-select v-model="form.billing_cycle" style="width:100%">
            <el-option v-for="item in cycleOptions" :key="item.value" :label="item.label" :value="item.value" />
          </el-select>
        </el-form-item>
        <el-form-item label="备注"><el-input v-model="form.remark" /></el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleSubmit">确定</el-button>
      </template>
    </el-dialog>
  </div>
</template>
