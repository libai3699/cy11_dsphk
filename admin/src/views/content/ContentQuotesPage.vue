<script setup lang="ts">
import { onMounted, reactive, ref } from 'vue';
import { ElMessage, ElMessageBox } from 'element-plus';
import { getQuoteList, createQuote, updateQuote, deleteQuote, type Quote } from '#/api/admin/quotes';

const loading = ref(false);
const list = ref<Quote[]>([]);
const dialogVisible = ref(false);
const isEdit = ref(false);
const editId = ref(0);
const form = reactive({ content: '', author: '', is_active: 1 as number, sort_order: 0 });

async function load() {
  loading.value = true;
  try { list.value = await getQuoteList() ?? []; } finally { loading.value = false; }
}

function openCreate() {
  isEdit.value = false;
  Object.assign(form, { content: '', author: '', is_active: 1, sort_order: 0 });
  dialogVisible.value = true;
}

function openEdit(row: Quote) {
  isEdit.value = true; editId.value = row.id;
  Object.assign(form, { content: row.content, author: row.author, is_active: row.is_active, sort_order: row.sort_order });
  dialogVisible.value = true;
}

async function handleSubmit() {
  if (!form.content) { ElMessage.warning('请输入语录内容'); return; }
  const data = { ...form };
  if (isEdit.value) { await updateQuote(editId.value, data); ElMessage.success('更新成功'); }
  else { await createQuote(data); ElMessage.success('创建成功'); }
  dialogVisible.value = false; load();
}

async function handleDelete(row: Quote) {
  await ElMessageBox.confirm('确定删除该语录？', '提示', { type: 'warning' });
  await deleteQuote(row.id); ElMessage.success('删除成功'); load();
}

onMounted(load);
</script>

<template>
  <div class="p-4">
    <el-card>
      <div class="mb-4 flex justify-end">
        <el-button type="primary" @click="openCreate">新增语录</el-button>
      </div>
      <el-table :data="list" v-loading="loading" border stripe>
        <el-table-column prop="id" label="ID" width="70" />
        <el-table-column prop="content" label="内容" show-overflow-tooltip />
        <el-table-column prop="author" label="作者" width="120" />
        <el-table-column label="状态" width="80">
          <template #default="{ row }"><el-tag :type="row.is_active === 1 ? 'success' : 'info'">{{ row.is_active === 1 ? '启用' : '禁用' }}</el-tag></template>
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

    <el-dialog v-model="dialogVisible" :title="isEdit ? '编辑语录' : '新增语录'" width="520px">
      <el-form :model="form" label-width="90px">
        <el-form-item label="内容" required><el-input v-model="form.content" type="textarea" :rows="4" /></el-form-item>
        <el-form-item label="作者"><el-input v-model="form.author" placeholder="可选" /></el-form-item>
        <el-form-item label="状态">
          <el-radio-group v-model="form.is_active"><el-radio :value="1">启用</el-radio><el-radio :value="0">禁用</el-radio></el-radio-group>
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
