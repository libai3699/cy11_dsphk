<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue';
import { ElMessage, ElMessageBox } from 'element-plus';
import { getUserList, createUser, updateUser, deleteUser, addUserDuration, type User } from '#/api/admin/users';

const loading = ref(false);
const list = ref<User[]>([]);
const total = ref(0);
const keyword = ref('');
const page = reactive({ current: 1, size: 20 });

const dialogVisible = ref(false);
const isEdit = ref(false);
const editId = ref(0);
const form = reactive({ username: '', password: '', phone: '', status: 1, free_used_seconds: 0, free_limit_seconds: 2700 });

const durationDialogVisible = ref(false);
const durationUserId = ref(0);
const durationDays = ref(30);

async function load() {
  loading.value = true;
  try {
    const res = await getUserList({ page: page.current, size: page.size, keyword: keyword.value });
    list.value = res?.list ?? [];
    total.value = res?.total ?? 0;
  } finally { loading.value = false; }
}

function handleSearch() { page.current = 1; load(); }

function openCreate() {
  isEdit.value = false;
  Object.assign(form, { username: '', password: '', phone: '', status: 1, free_used_seconds: 0, free_limit_seconds: 2700 });
  dialogVisible.value = true;
}

function openEdit(row: User) {
  isEdit.value = true;
  editId.value = row.id;
  Object.assign(form, { username: row.username, password: '', phone: row.phone || '', status: row.status, free_used_seconds: row.free_used_seconds, free_limit_seconds: row.free_limit_seconds });
  dialogVisible.value = true;
}

async function handleSubmit() {
  if (!isEdit.value) {
    if (!form.username || !form.password) { ElMessage.warning('用户名和密码必填'); return; }
    await createUser({ username: form.username, password: form.password, phone: form.phone, free_limit_seconds: form.free_limit_seconds });
    ElMessage.success('创建成功');
  } else {
    const data: any = { status: form.status, phone: form.phone, free_used_seconds: form.free_used_seconds, free_limit_seconds: form.free_limit_seconds };
    if (form.password) data.password = form.password;
    await updateUser(editId.value, data);
    ElMessage.success('更新成功');
  }
  dialogVisible.value = false;
  load();
}

async function handleToggle(row: User) {
  await updateUser(row.id, { status: row.status === 1 ? 0 : 1 });
  ElMessage.success(row.status === 1 ? '已禁用' : '已启用');
  load();
}

async function handleDelete(row: User) {
  await ElMessageBox.confirm(`确定删除用户 ${row.username}？`, '提示', { type: 'warning' });
  await deleteUser(row.id);
  ElMessage.success('删除成功');
  load();
}

function openAddDuration(row: User) {
  durationUserId.value = row.id;
  durationDays.value = 30;
  durationDialogVisible.value = true;
}

async function handleAddDuration() {
  if (durationDays.value < 1) {
    ElMessage.warning('追加天数必须大于0');
    return;
  }
  await addUserDuration(durationUserId.value, durationDays.value);
  ElMessage.success('追加成功');
  durationDialogVisible.value = false;
  load();
}

const fmtSec = (s: number) => `${Math.floor(s / 60)} 分钟`;
const fmtTime = (value?: string | null) => {
  if (!value) return '-';
  const date = new Date(value);
  if (Number.isNaN(date.getTime())) return value;
  const pad = (n: number) => String(n).padStart(2, '0');
  return `${date.getFullYear()}-${pad(date.getMonth() + 1)}-${pad(date.getDate())} ${pad(date.getHours())}:${pad(date.getMinutes())}:${pad(date.getSeconds())}`;
};
const fmtRemaining = (expiredAt?: string | null) => {
  if (!expiredAt) return '未开通';
  const sec = Math.max(0, Math.floor((new Date(expiredAt).getTime() - Date.now()) / 1000));
  if (sec <= 0) return '已过期';
  const days = Math.floor(sec / 86400);
  const hours = Math.floor((sec % 86400) / 3600);
  return days > 0 ? `${days}天 ${hours}小时` : `${Math.floor(sec / 60)}分钟`;
};
const fmtTraffic = (used: number, limit: number | null) => {
  if (limit === null || limit === undefined) return '不限';
  const remain = Math.max(0, limit - used);
  return `${(remain / 1024 / 1024 / 1024).toFixed(2)} GB`;
};

onMounted(load);
</script>

<template>
  <div class="p-4">
    <el-card>
      <div class="mb-4 flex items-center justify-between">
        <el-input v-model="keyword" placeholder="搜索用户名/手机号/设备ID" style="width:280px" clearable @keyup.enter="handleSearch">
          <template #append><el-button @click="handleSearch">搜索</el-button></template>
        </el-input>
        <el-button type="primary" @click="openCreate">新增用户</el-button>
      </div>

      <el-table :data="list" v-loading="loading" border stripe>
        <el-table-column prop="id" label="ID" width="70" />
        <el-table-column prop="username" label="用户名" width="120" />
        <!-- <el-table-column prop="phone" label="手机号" width="130" /> -->
        <el-table-column prop="device_id" label="设备ID" width="110" show-overflow-tooltip />
        <el-table-column label="展示ID" width="90">
          <template #default="{ row }">
            <span style="font-weight:600;letter-spacing:1px">{{ row.display_id || '-' }}</span>
          </template>
        </el-table-column>
        <el-table-column label="套餐到期" width="155">
          <template #default="{ row }">{{ row.plan_expired_at ? fmtTime(row.plan_expired_at) : '未开通' }}</template>
        </el-table-column>
        <el-table-column label="剩余时长" width="130">
          <template #default="{ row }">
            <el-tag :type="!row.plan_expired_at ? 'info' : new Date(row.plan_expired_at) > new Date() ? 'success' : 'danger'" size="small">
              {{ fmtRemaining(row.plan_expired_at) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="剩余流量" width="120">
          <template #default="{ row }">
            <span :style="{ color: row.traffic_limit_bytes !== null && (row.traffic_limit_bytes - row.traffic_used_bytes) <= 0 ? '#f56c6c' : '' }">
              {{ fmtTraffic(row.traffic_used_bytes, row.traffic_limit_bytes) }}
            </span>
          </template>
        </el-table-column>
        <el-table-column label="免费时长" width="150">
          <template #default="{ row }">{{ fmtSec(row.free_used_seconds) }} / {{ fmtSec(row.free_limit_seconds) }}</template>
        </el-table-column>
        <el-table-column label="线路ID" width="90">
          <template #default="{ row }">{{ row.current_line_id ?? '未分配' }}</template>
        </el-table-column>
        <el-table-column label="状态" width="80">
          <template #default="{ row }">
            <el-tag :type="row.status === 1 ? 'success' : 'danger'">{{ row.status === 1 ? '正常' : '禁用' }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="created_at" label="注册时间" width="170">
          <template #default="{ row }">{{ fmtTime(row.created_at) }}</template>
        </el-table-column>
        <el-table-column prop="last_ip" label="最后IP" width="140" />
        <el-table-column label="IP 详情" min-width="220">
          <template #default="{ row }">
            <div>{{ row.last_ip_detail?.location || '未知' }}</div>
            <div class="text-xs text-gray-500">
              {{ row.last_ip_detail?.type || '-' }}
              <span v-if="row.last_ip_detail?.is_private"> / 内网</span>
            </div>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="280" fixed="right">
          <template #default="{ row }">
            <el-button size="small" @click="openEdit(row)">编辑</el-button>
            <el-button size="small" type="primary" @click="openAddDuration(row)">追加时长</el-button>
            <el-button size="small" :type="row.status === 1 ? 'warning' : 'success'" @click="handleToggle(row)">
              {{ row.status === 1 ? '禁用' : '启用' }}
            </el-button>
            <el-button size="small" type="danger" @click="handleDelete(row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>

      <div class="mt-4 flex justify-end">
        <el-pagination
          v-model:current-page="page.current"
          v-model:page-size="page.size"
          :total="total"
          :page-sizes="[20, 50, 100]"
          layout="total, sizes, prev, pager, next"
          @change="load"
        />
      </div>
    </el-card>

    <el-dialog v-model="dialogVisible" :title="isEdit ? '编辑用户' : '新增用户'" width="480px">
      <el-form :model="form" label-width="140px">
        <el-form-item label="用户名" required>
          <el-input v-model="form.username" :disabled="isEdit" />
        </el-form-item>
        <el-form-item :label="isEdit ? '新密码（留空不改）' : '密码'" :required="!isEdit">
          <el-input v-model="form.password" type="password" show-password />
        </el-form-item>
        <!-- <el-form-item label="手机号">
          <el-input v-model="form.phone" />
        </el-form-item> -->
        <el-form-item v-if="isEdit" label="状态">
          <el-radio-group v-model="form.status">
            <el-radio :value="1">正常</el-radio>
            <el-radio :value="0">禁用</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item v-if="isEdit" label="已用免费时长(秒)">
          <el-input-number v-model="form.free_used_seconds" :min="0" style="width:100%" />
        </el-form-item>
        <el-form-item label="免费时长上限(秒)">
          <el-input-number v-model="form.free_limit_seconds" :min="0" style="width:100%" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleSubmit">确定</el-button>
      </template>
    </el-dialog>

    <el-dialog v-model="durationDialogVisible" title="追加时长" width="400px">
      <el-form label-width="100px">
        <el-form-item label="追加天数" required>
          <el-input-number v-model="durationDays" :min="1" :max="3650" style="width:100%" />
        </el-form-item>
        <el-alert type="info" :closable="false" show-icon>
          <template #default>
            <div>如果用户当前有剩余时长，将在原有基础上追加</div>
            <div>如果已过期或未开通，将从当前时间开始计算</div>
          </template>
        </el-alert>
      </el-form>
      <template #footer>
        <el-button @click="durationDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleAddDuration">确定</el-button>
      </template>
    </el-dialog>
  </div>
</template>
