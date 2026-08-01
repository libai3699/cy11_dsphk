<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue';
import { getDeviceList, type Device } from '#/api/admin/devices';

const loading = ref(false);
const list = ref<Device[]>([]);
const total = ref(0);
const keyword = ref('');
const page = reactive({ current: 1, size: 20 });

async function load() {
  loading.value = true;
  try {
    const res = await getDeviceList({ page: page.current, size: page.size, keyword: keyword.value });
    list.value = res?.list ?? [];
    total.value = res?.total ?? 0;
  } finally { loading.value = false; }
}

function handleSearch() { page.current = 1; load(); }

function fmtTime(value?: string | null) {
  if (!value) return '-';
  const date = new Date(value);
  if (Number.isNaN(date.getTime())) return value;
  const pad = (n: number) => String(n).padStart(2, '0');
  return `${date.getFullYear()}-${pad(date.getMonth() + 1)}-${pad(date.getDate())} ${pad(date.getHours())}:${pad(date.getMinutes())}:${pad(date.getSeconds())}`;
}

onMounted(load);
</script>

<template>
  <div class="p-4">
    <el-card>
      <div class="mb-4">
        <el-input v-model="keyword" placeholder="搜索设备ID/品牌/型号" style="width:280px" clearable @keyup.enter="handleSearch">
          <template #append><el-button @click="handleSearch">搜索</el-button></template>
        </el-input>
      </div>
      <el-table :data="list" v-loading="loading" border stripe>
        <el-table-column prop="id" label="ID" width="70" />
        <el-table-column prop="display_id" label="展示ID" width="100" />
        <el-table-column prop="device_id" label="设备ID" show-overflow-tooltip />
        <el-table-column prop="user_id" label="用户ID" width="90" />
        <el-table-column prop="brand" label="品牌" width="100" />
        <el-table-column prop="model" label="型号" width="150" />
        <el-table-column prop="os_version" label="Android版本" width="110" />
        <el-table-column prop="app_version" label="App版本" width="100" />
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
        <el-table-column label="最后活跃" width="170">
          <template #default="{ row }">{{ fmtTime(row.last_seen_at) }}</template>
        </el-table-column>
        <el-table-column label="创建时间" width="170">
          <template #default="{ row }">{{ fmtTime(row.created_at) }}</template>
        </el-table-column>
        <el-table-column label="更新时间" width="170">
          <template #default="{ row }">{{ fmtTime(row.updated_at) }}</template>
        </el-table-column>
      </el-table>
      <div class="mt-4 flex justify-end">
        <el-pagination v-model:current-page="page.current" v-model:page-size="page.size" :total="total" layout="total, prev, pager, next" @change="load" />
      </div>
    </el-card>
  </div>
</template>
