import { requestClient } from '#/api/request';

export interface WorkspaceMetric {
  hint: string;
  key: string;
  label: string;
  path: string;
  value: string;
}

export interface WorkspaceMerchant {
  city: string;
  commissionRate: number;
  estimatedCommission: number;
  id: number;
  industry: string;
  name: string;
  nextAction: string;
  recentWriteOffAmount: number;
  stage: string;
}

export interface WorkspaceShootingTask {
  assignee: string;
  deadline?: string;
  id: number;
  merchantId: number;
  merchantName: string;
  shotCount: number;
  status: string;
  taskTitle: string;
}

export interface WorkspaceTopic {
  hook: string;
  id: number;
  merchantId: number;
  merchantName: string;
  publishWindow: string;
  status: string;
  title: string;
}

export interface WorkspaceReview {
  dealCount: number;
  id: number;
  merchantId: number;
  merchantName: string;
  playCount: number;
  status: string;
  videoTitle: string;
  writeOffAmount: number;
}

export interface WorkspaceOverview {
  merchants: WorkspaceMerchant[];
  metrics: WorkspaceMetric[];
  reviews: WorkspaceReview[];
  shootingTasks: WorkspaceShootingTask[];
  topics: WorkspaceTopic[];
}

interface ListResult<T> {
  list: T[];
  total: number;
}

interface MerchantRow {
  city: string;
  commissionRate: number;
  contactName: string;
  contactPhone: string;
  cooperationType: string;
  douyinAccount: string;
  douyinLaikeAccount: string;
  id: number;
  industry: string;
  name: string;
  stage: string;
  status: number;
}

interface PackageRow {
  merchantId: number;
  status: number;
}

interface ShootingTaskRow {
  assignee: string;
  deadline?: string;
  id: number;
  merchantId: number;
  merchantName: string;
  shotCount: number;
  status: string;
  taskTitle: string;
}

interface TopicRow {
  hook: string;
  id: number;
  merchantId: number;
  merchantName: string;
  publishWindow: string;
  status: string;
  title: string;
}

interface ScheduleRow {
  status: string;
}

interface ReviewRow {
  createdAt: string;
  dealCount: number;
  id: number;
  merchantId: number;
  merchantName: string;
  playCount: number;
  status: string;
  videoTitle: string;
  writeOffAmount: number;
}

export async function getWorkspaceOverview() {
  const [
    merchantResult,
    packageResult,
    shootingResult,
    topicResult,
    scheduleResult,
    reviewResult,
  ] = await Promise.all([
    requestClient.get<ListResult<MerchantRow>>('/merchants', {
      params: { page: 1, size: 100 },
    }),
    requestClient.get<ListResult<PackageRow>>('/packages', {
      params: { page: 1, size: 100 },
    }),
    requestClient.get<ListResult<ShootingTaskRow>>('/shooting-tasks', {
      params: { page: 1, size: 100 },
    }),
    requestClient.get<ListResult<TopicRow>>('/topics', {
      params: { page: 1, size: 6 },
    }),
    requestClient.get<ListResult<ScheduleRow>>('/schedules', {
      params: { page: 1, size: 100 },
    }),
    requestClient.get<ListResult<ReviewRow>>('/reviews', {
      params: { page: 1, size: 100 },
    }),
  ]);

  const merchants = merchantResult.list || [];
  const packages = packageResult.list || [];
  const shootingTasks = shootingResult.list || [];
  const schedules = scheduleResult.list || [];
  const reviews = reviewResult.list || [];
  const recentStart = Date.now() - 7 * 24 * 60 * 60 * 1000;
  const merchantMap = new Map(merchants.map((item) => [item.id, item]));
  const recentReviews = reviews.filter((item) => {
    const time = new Date(item.createdAt || '').getTime();
    return Number.isFinite(time) && time >= recentStart;
  });
  const recentWriteOffAmount = recentReviews.reduce(
    (sum, item) => sum + Number(item.writeOffAmount || 0),
    0,
  );
  const estimatedCommission = recentReviews.reduce((sum, item) => {
    const merchant = merchantMap.get(item.merchantId);
    const rate = merchant?.commissionRate || 10;
    return sum + (Number(item.writeOffAmount || 0) * rate) / 100;
  }, 0);

  return {
    merchants: merchants
      .slice(0, 6)
      .map((item) => buildWorkspaceMerchant(item, packages, reviews)),
    metrics: [
      {
        hint: `总商家 ${merchantResult.total || 0} 家，启用套餐 ${
          packages.filter((item) => item.status === 1).length
        } 个`,
        key: 'merchants',
        label: '在运营商家',
        path: '/users/list',
        value: String(merchants.filter((item) => item.status === 1).length),
      },
      {
        hint: `待发布视频 ${
          schedules.filter((item) => item.status === '待发布').length
        } 条`,
        key: 'shooting',
        label: '待拍/待剪任务',
        path: '/logs/user',
        value: String(
          shootingTasks.filter((item) =>
            ['待拍摄', '拍摄中', '已拍摄', '已剪辑'].includes(item.status),
          ).length,
        ),
      },
      {
        hint: `已复盘 ${reviewResult.total || 0} 条视频`,
        key: 'writeoff',
        label: '近 7 天核销额',
        path: '/logs/admin',
        value: formatMoney(recentWriteOffAmount),
      },
      {
        hint: '按各商家分成比例估算',
        key: 'commission',
        label: '预估分成',
        path: '/plans/orders',
        value: formatMoney(estimatedCommission),
      },
    ],
    reviews: reviews.slice(0, 6).map((item) => ({
      dealCount: item.dealCount,
      id: item.id,
      merchantId: item.merchantId,
      merchantName: item.merchantName,
      playCount: item.playCount,
      status: item.status,
      videoTitle: item.videoTitle,
      writeOffAmount: item.writeOffAmount,
    })),
    shootingTasks: shootingTasks
      .filter((item) => item.status !== '已完成')
      .slice(0, 6)
      .map((item) => ({
        assignee: item.assignee,
        deadline: item.deadline,
        id: item.id,
        merchantId: item.merchantId,
        merchantName: item.merchantName,
        shotCount: item.shotCount,
        status: item.status,
        taskTitle: item.taskTitle,
      })),
    topics: (topicResult.list || []).map((item) => ({
      hook: item.hook,
      id: item.id,
      merchantId: item.merchantId,
      merchantName: item.merchantName,
      publishWindow: item.publishWindow,
      status: item.status,
      title: item.title,
    })),
  };
}

function buildWorkspaceMerchant(
  item: MerchantRow,
  packages: PackageRow[],
  reviews: ReviewRow[],
) {
  const recentStart = Date.now() - 7 * 24 * 60 * 60 * 1000;
  const recentWriteOffAmount = reviews
    .filter((review) => {
      const time = new Date(review.createdAt || '').getTime();
      return (
        review.merchantId === item.id &&
        Number.isFinite(time) &&
        time >= recentStart
      );
    })
    .reduce((sum, review) => sum + Number(review.writeOffAmount || 0), 0);

  return {
    city: item.city,
    commissionRate: item.commissionRate,
    estimatedCommission:
      (recentWriteOffAmount * (item.commissionRate || 10)) / 100,
    id: item.id,
    industry: item.industry,
    name: item.name,
    nextAction: buildMerchantNextAction(item, packages),
    recentWriteOffAmount,
    stage: item.stage,
  };
}

function buildMerchantNextAction(item: MerchantRow, packages: PackageRow[]) {
  if (!item.industry || !item.city || !item.contactName || !item.contactPhone) {
    return '补齐基础档案';
  }
  if (!item.cooperationType || item.commissionRate <= 0) {
    return '确认合作规则';
  }
  if (!packages.some((pkg) => pkg.merchantId === item.id && pkg.status === 1)) {
    return '新增启用套餐';
  }
  if (!item.douyinAccount && !item.douyinLaikeAccount) {
    return '记录账号授权';
  }
  return '推进选题和内容生产';
}

function formatMoney(value: number) {
  return `¥${Number(value || 0).toFixed(2)}`;
}
