<script setup lang="ts">
import { onMounted, reactive, ref } from 'vue';
import { ElMessage, ElMessageBox } from 'element-plus';
import { Delete, Plus } from '@element-plus/icons-vue';
import { uploadPaymentImage } from '#/api/admin/files';
import { getPaymentConfigList, createPaymentConfig, updatePaymentConfig, deletePaymentConfig, type PaymentConfig } from '#/api/admin/payments';

const loading = ref(false);
const list = ref<PaymentConfig[]>([]);
const dialogVisible = ref(false);
const isEdit = ref(false);
const editId = ref(0);
const form = reactive({ type: '', label: '', address: '', qr_code: '', is_active: 1 as number, sort_order: 0, remark: '' });
const uploading = ref(false);

async function load() {
  loading.value = true;
  try { list.value = await getPaymentConfigList() ?? []; } finally { loading.value = false; }
}

function openCreate() {
  isEdit.value = false;
  Object.assign(form, { type: '', label: '', address: '', qr_code: '', is_active: 1, sort_order: 0, remark: '' });
  dialogVisible.value = true;
}

function openEdit(row: PaymentConfig) {
  isEdit.value = true; editId.value = row.id;
  Object.assign(form, { type: row.type, label: row.label, address: row.address, qr_code: row.qr_code, is_active: row.is_active, sort_order: row.sort_order, remark: row.remark });
  dialogVisible.value = true;
}

async function handleUpload(file: File) {
  if (!file.type.startsWith('image/')) {
    ElMessage.error('只能上传图片文件');
    return;
  }

  if (file.size > 5 * 1024 * 1024) {
    ElMessage.error('图片大小不能超过 5MB');
    return;
  }

  uploading.value = true;
  try {
    const data = new FormData();
    data.append('file', file);
    const result = await uploadPaymentImage(data);
    form.qr_code = result.url;
    ElMessage.success('上传成功');
  } finally {
    uploading.value = false;
  }
}

function handleRemoveQrCode() {
  form.qr_code = '';
}

async function handleSubmit() {
  if (!form.type || !form.label) { ElMessage.warning('请填写必填项'); return; }
  const data = { ...form };
  if (isEdit.value) { await updatePaymentConfig(editId.value, data); ElMessage.success('更新成功'); }
  else { await createPaymentConfig(data); ElMessage.success('创建成功'); }
  dialogVisible.value = false; load();
}

async function handleDelete(row: PaymentConfig) {
  await ElMessageBox.confirm('确定删除该支付配置？', '提示', { type: 'warning' });
  await deletePaymentConfig(row.id); ElMessage.success('删除成功'); load();
}

onMounted(load);
</script>

<template>
  <div class="p-4">
    <el-card>
      <div class="mb-4 flex justify-end">
        <el-button type="primary" @click="openCreate">新增支付方式</el-button>
      </div>
      <el-table :data="list" v-loading="loading" border stripe>
        <el-table-column prop="id" label="ID" width="70" />
        <el-table-column prop="type" label="类型" width="130" />
        <el-table-column prop="label" label="显示名称" width="150" />
        <el-table-column prop="address" label="地址/账号" show-overflow-tooltip />
        <el-table-column prop="qr_code" label="二维码" width="100">
          <template #default="{ row }">
            <el-image v-if="row.qr_code" :src="row.qr_code" style="width: 40px; height: 40px" :preview-src-list="[row.qr_code]" fit="cover" />
            <span v-else class="text-gray-400">-</span>
          </template>
        </el-table-column>
        <el-table-column label="状态" width="80">
          <template #default="{ row }"><el-tag :type="row.is_active === 1 ? 'success' : 'info'">{{ row.is_active === 1 ? '启用' : '禁用' }}</el-tag></template>
        </el-table-column>
        <el-table-column prop="sort_order" label="排序" width="80" />
        <el-table-column label="操作" width="150">
          <template #default="{ row }">
            <el-button size="small" @click="openEdit(row)">编辑</el-button>
            <el-button size="small" type="danger" @click="handleDelete(row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <el-dialog v-model="dialogVisible" :title="isEdit ? '编辑支付配置' : '新增支付配置'" width="600px">
      <el-form :model="form" label-width="100px">
        <el-form-item label="类型" required>
          <el-select v-model="form.type" placeholder="请选择" style="width:100%" :disabled="isEdit">
            <el-option label="USDT (TRC20)" value="usdt_trc20" />
            <el-option label="USDT (BEP20)" value="usdt_bep20" />
            <el-option label="USDT (ERC20)" value="usdt_erc20" />
            <el-option label="微信支付" value="wechat" />
            <el-option label="支付宝" value="alipay" />
          </el-select>
        </el-form-item>
        <el-form-item label="显示名称" required><el-input v-model="form.label" /></el-form-item>
        <el-form-item label="地址/账号"><el-input v-model="form.address" placeholder="USDT地址或收款账号" /></el-form-item>
        <el-form-item label="二维码">
          <div class="flex flex-col gap-2">
            <div v-if="form.qr_code" class="relative inline-block">
              <el-image
                :src="form.qr_code"
                style="width: 200px; height: 200px"
                fit="contain"
                :preview-src-list="[form.qr_code]"
              />
              <el-button
                type="danger"
                size="small"
                circle
                :icon="Delete"
                class="absolute top-0 right-0"
                @click="handleRemoveQrCode"
              />
            </div>
            <el-upload
              :show-file-list="false"
              :before-upload="(file) => { handleUpload(file); return false; }"
              accept="image/*"
              :disabled="uploading"
            >
              <el-button :loading="uploading" :icon="Plus">
                {{ form.qr_code ? '更换二维码' : '上传二维码' }}
              </el-button>
            </el-upload>
            <div class="text-xs text-gray-500">
              支持 JPG、PNG、WebP、GIF 格式，大小不超过 5MB
            </div>
            <el-input v-model="form.qr_code" placeholder="输入二维码图片URL">
              <template #prepend>图片URL</template>
            </el-input>
          </div>
        </el-form-item>
        <el-form-item label="状态">
          <el-radio-group v-model="form.is_active"><el-radio :value="1">启用</el-radio><el-radio :value="0">禁用</el-radio></el-radio-group>
        </el-form-item>
        <el-form-item label="排序"><el-input-number v-model="form.sort_order" style="width:100%" /></el-form-item>
        <el-form-item label="备注"><el-input v-model="form.remark" type="textarea" :rows="2" /></el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleSubmit" :loading="uploading">确定</el-button>
      </template>
    </el-dialog>
  </div>
</template>
