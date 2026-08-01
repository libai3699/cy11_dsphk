<script setup lang="ts">
import { ElMessage } from 'element-plus';
import { useRouter } from 'vue-router';

import {
  dashboardMetrics,
  merchants,
  reviewRows,
  shootTasks,
  topics,
} from '#/mock/content-ops';

defineOptions({ name: 'WorkspacePage' });

const router = useRouter();

const workflowSteps = [
  { label: '商家建档', path: '/users/list' },
  { label: '账号授权', path: '/users/devices' },
  { label: '对标分析', path: '/lines/list' },
  { label: '生成选题', path: '/content/notices' },
  { label: '生成文案', path: '/content/quotes' },
  { label: '生成分镜', path: '/content/discoveries' },
  { label: '派拍摄', path: '/logs/user' },
  { label: '排发布', path: '/content/payments' },
  { label: '做复盘', path: '/logs/admin' },
] as const;

function demoAction(text: string) {
  ElMessage.success(`${text}：演示动作已触发`);
}
</script>

<template>
  <div class="workspace-page">
    <el-card class="hero-card" shadow="never">
      <div class="hero-head">
        <div>
          <div class="eyebrow">内容生产后台样机</div>
          <h2>本地商家抖音获客工作台</h2>
          <p>按文档里的运营流程做成可点击样机：先建商家，再走账号、对标、选题、文案、分镜、拍摄、发布、复盘。</p>
        </div>
        <div class="hero-actions">
          <el-button type="primary" @click="router.push('/users/list')">新增商家</el-button>
          <el-button @click="demoAction('AI 生成本周运营计划')">生成本周计划</el-button>
        </div>
      </div>
      <div class="stats-grid">
        <div v-for="item in dashboardMetrics" :key="item.label" class="stat-card">
          <div class="stat-label">{{ item.label }}</div>
          <div class="stat-value">{{ item.value }}</div>
          <div class="stat-hint">{{ item.hint }}</div>
        </div>
      </div>
    </el-card>

    <el-card shadow="never">
      <template #header>
        <div class="card-header">
          <span>运营流程</span>
          <el-tag>按文档核对</el-tag>
        </div>
      </template>
      <div class="workflow-strip">
        <el-button
          v-for="(step, index) in workflowSteps"
          :key="step.label"
          plain
          type="primary"
          @click="router.push(step.path)"
        >
          {{ index + 1 }}. {{ step.label }}
        </el-button>
      </div>
    </el-card>

    <el-row :gutter="16" class="workspace-row">
      <el-col :span="14">
        <el-card shadow="never">
          <template #header>
            <div class="card-header">
              <span>重点商家</span>
              <el-tag type="success">2 家样例</el-tag>
            </div>
          </template>
          <el-table :data="merchants" border stripe>
            <el-table-column prop="merchant" label="商家" min-width="160" />
            <el-table-column prop="industry" label="行业" width="90" />
            <el-table-column prop="stage" label="阶段" width="120" />
            <el-table-column prop="recentGmv" label="近 7 天成交" width="130" />
            <el-table-column prop="nextAction" label="下一动作" min-width="220" />
            <el-table-column label="操作" width="220" fixed="right">
              <template #default="{ row }">
                <el-button size="small" @click="router.push('/content/notices')">做选题</el-button>
                <el-button size="small" type="primary" @click="demoAction(`诊断 ${row.merchant}`)">诊断</el-button>
              </template>
            </el-table-column>
          </el-table>
        </el-card>
      </el-col>
      <el-col :span="10">
        <el-card shadow="never">
          <template #header>
            <div class="card-header">
              <span>本周待执行</span>
              <el-tag type="warning">高优先级</el-tag>
            </div>
          </template>
          <div class="todo-list">
            <div v-for="task in shootTasks" :key="task.taskName" class="todo-item">
              <div class="todo-title">{{ task.taskName }}</div>
              <div class="todo-meta">{{ task.merchant }} · {{ task.deadline }} · {{ task.assignee }}</div>
              <el-button class="todo-button" size="small" type="primary" @click="demoAction(`推进 ${task.taskName}`)">
                推进任务
              </el-button>
            </div>
          </div>
        </el-card>
      </el-col>
    </el-row>

    <el-row :gutter="16" class="workspace-row">
      <el-col :span="12">
        <el-card shadow="never">
          <template #header>
            <div class="card-header">
              <span>选题池</span>
              <el-tag>待拍 / 待定稿</el-tag>
            </div>
          </template>
          <div class="topic-list">
            <div v-for="item in topics" :key="item.topic" class="topic-item">
              <div class="topic-title">{{ item.topic }}</div>
              <div class="topic-hook">{{ item.hook }}</div>
              <div class="topic-meta">{{ item.merchant }} · {{ item.publishWindow }} · {{ item.status }}</div>
              <el-button class="todo-button" size="small" type="primary" @click="router.push('/content/quotes')">
                生成文案
              </el-button>
            </div>
          </div>
        </el-card>
      </el-col>
      <el-col :span="12">
        <el-card shadow="never">
          <template #header>
            <div class="card-header">
              <span>最近复盘</span>
              <el-tag type="info">2 条样例</el-tag>
            </div>
          </template>
          <div class="review-list">
            <div v-for="item in reviewRows" :key="item.video" class="review-item">
              <div class="review-title">{{ item.video }}</div>
              <div class="review-meta">{{ item.merchant }} · {{ item.views }} · {{ item.deals }}</div>
              <div class="review-next">下一步：{{ item.nextAction }}</div>
              <el-button class="todo-button" size="small" @click="demoAction(`生成复盘建议：${item.video}`)">
                生成复盘建议
              </el-button>
            </div>
          </div>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<style scoped>
.workspace-page {
  padding: 20px;
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.hero-card {
  border-radius: 20px;
  background:
    radial-gradient(circle at top right, rgba(37, 99, 235, 0.18), transparent 28%),
    linear-gradient(135deg, #0f172a 0%, #1d4ed8 100%);
  border: 0;
  color: #fff;
}

.hero-head {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 16px;
}

.hero-actions {
  display: flex;
  flex-shrink: 0;
  gap: 10px;
}

.eyebrow {
  margin-bottom: 8px;
  color: rgba(255, 255, 255, 0.72);
  font-size: 12px;
  letter-spacing: 0.12em;
  text-transform: uppercase;
}

.hero-head h2 {
  margin: 0;
  font-size: 30px;
}

.hero-head p {
  margin: 8px 0 0;
  max-width: 760px;
  color: rgba(255, 255, 255, 0.82);
}

.stats-grid {
  margin-top: 24px;
  display: grid;
  grid-template-columns: repeat(4, minmax(0, 1fr));
  gap: 14px;
}

.workspace-row {
  margin: 0;
}

.card-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  font-weight: 600;
}

.workflow-strip {
  display: flex;
  flex-wrap: wrap;
  gap: 10px;
}

.stat-card {
  min-width: 0;
  padding: 18px;
  background: rgba(255, 255, 255, 0.12);
  border: 1px solid rgba(255, 255, 255, 0.16);
  border-radius: 16px;
  backdrop-filter: blur(6px);
}

.stat-label {
  color: rgba(255, 255, 255, 0.72);
  font-size: 13px;
}

.stat-value {
  margin-top: 10px;
  color: #fff;
  font-size: 30px;
  font-weight: 700;
}

.stat-hint {
  margin-top: 8px;
  color: rgba(255, 255, 255, 0.68);
  font-size: 12px;
}

.todo-list,
.topic-list,
.review-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.todo-item,
.topic-item,
.review-item {
  border: 1px solid #e5e7eb;
  border-radius: 14px;
  padding: 14px 16px;
  background: #fff;
}

.todo-title,
.topic-title,
.review-title {
  color: #0f172a;
  font-weight: 600;
}

.todo-meta,
.topic-meta,
.review-meta {
  margin-top: 6px;
  color: #64748b;
  font-size: 13px;
}

.topic-hook,
.review-next {
  margin-top: 8px;
  color: #1e293b;
  font-size: 14px;
}

.todo-button {
  margin-top: 10px;
}

@media (max-width: 1200px) {
  .stats-grid {
    grid-template-columns: repeat(2, minmax(0, 1fr));
  }
}
</style>
