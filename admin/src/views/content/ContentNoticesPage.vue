<script setup lang="ts">
import { onMounted, reactive, ref } from 'vue';
import { ElMessage, ElMessageBox } from 'element-plus';
import { getNoticeList, createNotice, updateNotice, deleteNotice, type Notice } from '#/api/admin/notices';

const loading = ref(false);
const list = ref<Notice[]>([]);
const dialogVisible = ref(false);
const isEdit = ref(false);
const editId = ref(0);
const form = reactive({ content: '', target_user_id: null as number | null, type: 1, is_active: 1 as number, sort_order: 0 });

async function load() {
  loading.value = true;
  try { list.value = await getNoticeList() ?? []; } finally { loading.value = false; }
}

function openCreate() {
  isEdit.value = false;
  Object.assign(form, { content: '', target_user_id: null, type: 1, is_active: 1, sort_order: 0 });
  dialogVisible.value = true;
}

function openEdit(row: Notice) {
  isEdit.value = true; editId.value = row.id;
  Object.assign(form, { content: row.content, target_user_id: row.target_user_id ?? null, type: row.type, is_active: row.is_active, sort_order: row.sort_order });
  dialogVisible.value = true;
}

async function handleSubmit() {
  if (!form.content) { ElMessage.warning('请输入通知内容'); return; }
  const data = { ...form, target_user_id: form.target_user_id || null };
  if (isEdit.value) { await updateNotice(editId.value, data); ElMessage.success('更新成功'); }
  else { await createNotice(data); ElMessage.success('创建成功'); }
  dialogVisible.value = false; load();
}

async function handleDelete(row: Notice) {
  await ElMessageBox.confirm('确定删除该通知？', '提示', { type: 'warning' });
  await deleteNotice(row.id); ElMessage.success('删除成功'); load();
}

onMounted(load);
</script>

<template>
  <div class="p-4">
    <el-card>
      <div class="mb-4 flex justify-end">
        <el-button type="primary" @click="openCreate">新增通知</el-button>
      </div>
      <el-table :data="list" v-loading="loading" border stripe>
        <el-table-column prop="id" label="ID" width="70" />
        <el-table-column prop="content" label="内容" show-overflow-tooltip />
        <el-table-column label="目标用户" width="110">
          <template #default="{ row }">{{ row.target_user_id ?? '公共' }}</template>
        </el-table-column>
        <el-table-column label="类型" width="90">
          <template #default="{ row }"><el-tag :type="row.type === 2 ? 'danger' : 'primary'">{{ row.type === 2 ? '重要' : '普通' }}</el-tag></template>
        </el-table-column>
        <el-table-column label="状态" width="80">
          <template #default="{ row }"><el-tag :type="row.is_active === 1 ? 'success' : 'info'">{{ row.is_active === 1 ? '显示' : '隐藏' }}</el-tag></template>
        </el-table-column>
        <el-table-column prop="sort_order" label="排序" width="80" />
        <el-table-column prop="created_at" label="创建时间" width="170" />
        <el-table-column label="操作" width="150">
          <template #default="{ row }">
            <el-button size="small" @click="openEdit(row)">编辑</el-button>
            <el-button size="small" type="danger" @click="handleDelete(row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <el-dialog v-model="dialogVisible" :title="isEdit ? '编辑通知' : '新增通知'" width="520px">
      <el-form :model="form" label-width="90px">
        <el-form-item label="内容" required><el-input v-model="form.content" type="textarea" :rows="4" /></el-form-item>
        <el-form-item label="目标用户"><el-input-number v-model="form.target_user_id" style="width:100%" :min="1" placeholder="留空为公共通知" /></el-form-item>
        <el-form-item label="类型">
          <el-radio-group v-model="form.type"><el-radio :value="1">普通</el-radio><el-radio :value="2">重要</el-radio></el-radio-group>
        </el-form-item>
        <el-form-item label="状态">
          <el-radio-group v-model="form.is_active"><el-radio :value="1">显示</el-radio><el-radio :value="0">隐藏</el-radio></el-radio-group>
        </el-form-item>
        <el-form-item label="排序"><el-input-number v-model="form.sort_order" style="width:100%" /></el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleSubmit">确定</el-button>
      </template>
    </el-dialog>
  </div>
</template>
