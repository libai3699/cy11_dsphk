<script setup lang="ts">
import { onMounted, reactive, ref } from 'vue';

import { ElMessage } from 'element-plus';

import {
  assignUserLine,
  createLine,
  deleteLine,
  getLineList,
  updateLine,
  type VpnLine,
} from '#/api/admin/lines';

const loading = ref(false);
const list = ref<VpnLine[]>([]);
const editOpen = ref(false);
const assignOpen = ref(false);
const editing = ref<VpnLine | null>(null);
const form = reactive<Partial<VpnLine>>({
  address: '',
  description: '',
  is_active: 1,
  is_default: 0,
  name: '',
  protocol: 'VMESS',
  raw_uri: '',
  region: '',
  sort_order: 0,
});
const assignForm = reactive({ line_id: 0, notice: '', user_id: undefined as number | undefined });

async function load() {
  loading.value = true;
  try {
    list.value = (await getLineList()) ?? [];
  } finally {
    loading.value = false;
  }
}

function openCreate() {
  editing.value = null;
  Object.assign(form, {
    address: '',
    description: '',
    is_active: 1,
    is_default: 0,
    name: '',
    protocol: 'VMESS',
    raw_uri: '',
    region: '',
    sort_order: 0,
  });
  editOpen.value = true;
}

function openEdit(row: VpnLine) {
  editing.value = row;
  Object.assign(form, row);
  editOpen.value = true;
}

async function saveLine() {
  if (editing.value) await updateLine(editing.value.id, form);
  else await createLine(form);
  ElMessage.success('保存成功');
  editOpen.value = false;
  await load();
}

async function removeLine(row: VpnLine) {
  await deleteLine(row.id);
  ElMessage.success('删除成功');
  await load();
}

function openAssign(row: VpnLine) {
  Object.assign(assignForm, {
    line_id: row.id,
    notice: `你的 VPN 线路已切换为：${row.name}`,
    user_id: undefined,
  });
  assignOpen.value = true;
}

async function saveAssign() {
  if (!assignForm.user_id) {
    ElMessage.error('请输入用户ID');
    return;
  }
  await assignUserLine(assignForm as { line_id: number; notice: string; user_id: number });
  ElMessage.success('已切换用户线路并发送通知');
  assignOpen.value = false;
}

onMounted(load);
</script>

<template>
  <div class="p-4">
    <el-card v-loading="loading">
      <template #header>
        <div class="flex items-center justify-between">
          <span>线路管理</span>
          <el-button type="primary" @click="openCreate">新增线路</el-button>
        </div>
      </template>

      <el-table :data="list" border stripe>
        <el-table-column prop="id" label="ID" width="70" />
        <el-table-column prop="name" label="线路名称" width="160" />
        <el-table-column prop="region" label="区域" width="100" />
        <el-table-column prop="protocol" label="协议" width="100" />
        <el-table-column prop="address" label="地址" min-width="180" show-overflow-tooltip />
        <el-table-column label="默认" width="80">
          <template #default="{ row }">
            <el-tag v-if="row.is_default === 1" type="success">默认</el-tag>
          </template>
        </el-table-column>
        <el-table-column label="状态" width="90">
          <template #default="{ row }">
            <el-tag :type="row.is_active === 1 ? 'success' : 'info'">
              {{ row.is_active === 1 ? '启用' : '停用' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="250" fixed="right">
          <template #default="{ row }">
            <el-button size="small" @click="openEdit(row)">编辑</el-button>
            <el-button size="small" type="warning" @click="openAssign(row)">切给用户</el-button>
            <el-button size="small" type="danger" @click="removeLine(row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <el-dialog v-model="editOpen" :title="editing ? '编辑线路' : '新增线路'" width="720px">
      <el-form label-position="top">
        <el-form-item label="线路名称"><el-input v-model="form.name" /></el-form-item>
        <el-form-item label="区域"><el-input v-model="form.region" /></el-form-item>
        <el-form-item label="协议"><el-input v-model="form.protocol" /></el-form-item>
        <el-form-item label="地址"><el-input v-model="form.address" /></el-form-item>
        <el-form-item label="原始链接 raw_uri">
          <el-input v-model="form.raw_uri" type="textarea" :rows="5" />
        </el-form-item>
        <el-form-item label="描述"><el-input v-model="form.description" /></el-form-item>
        <el-form-item label="排序"><el-input-number v-model="form.sort_order" :min="0" /></el-form-item>
        <el-form-item label="开关">
          <el-switch v-model="form.is_active" :active-value="1" :inactive-value="0" active-text="启用" inactive-text="停用" />
          <el-switch class="ml-4" v-model="form.is_default" :active-value="1" :inactive-value="0" active-text="默认" inactive-text="非默认" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="editOpen = false">取消</el-button>
        <el-button type="primary" @click="saveLine">保存</el-button>
      </template>
    </el-dialog>

    <el-dialog v-model="assignOpen" title="给用户切换线路" width="520px">
      <el-form label-position="top">
        <el-form-item label="用户ID"><el-input-number v-model="assignForm.user_id" :min="1" /></el-form-item>
        <el-form-item label="通知内容"><el-input v-model="assignForm.notice" type="textarea" :rows="3" /></el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="assignOpen = false">取消</el-button>
        <el-button type="primary" @click="saveAssign">确认切换</el-button>
      </template>
    </el-dialog>
  </div>
</template>
