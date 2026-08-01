<script setup lang="ts">
defineOptions({ name: 'WorkspacePage' });
</script>

<template>
  <div class="workspace-page">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>欢迎使用 9点9 VPN 管理后台</span>
        </div>
      </template>
      <div class="welcome-content">
        <h2>系统功能</h2>
        <el-row :gutter="20">
          <el-col :span="8">
            <el-card shadow="hover">
              <h3>用户管理</h3>
              <p>管理用户账号、设备、时长等</p>
              <el-button type="primary" @click="$router.push('/users/list')">进入</el-button>
            </el-card>
          </el-col>
          <el-col :span="8">
            <el-card shadow="hover">
              <h3>套餐管理</h3>
              <p>管理套餐、订单等</p>
              <el-button type="primary" @click="$router.push('/plans/list')">进入</el-button>
            </el-card>
          </el-col>
          <el-col :span="8">
            <el-card shadow="hover">
              <h3>内容管理</h3>
              <p>管理语录、支付配置等</p>
              <el-button type="primary" @click="$router.push('/content/quotes')">进入</el-button>
            </el-card>
          </el-col>
        </el-row>
      </div>
    </el-card>
  </div>
</template>

<style scoped>
.workspace-page {
  padding: 20px;
}

.card-header {
  font-size: 18px;
  font-weight: bold;
}

.welcome-content {
  padding: 20px 0;
}

.welcome-content h2 {
  margin-bottom: 20px;
  color: #333;
}

.welcome-content h3 {
  margin-bottom: 10px;
  color: #666;
}

.welcome-content p {
  margin-bottom: 15px;
  color: #999;
}
</style>
