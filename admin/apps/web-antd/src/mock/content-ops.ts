export const dashboardMetrics = [
  { label: '在运营商家', value: '8', hint: '本周新增 2 家' },
  { label: '本周待拍视频', value: '14', hint: '3 条已排期' },
  { label: '近 7 天核销额', value: '¥28,640', hint: '样机静态数据' },
  { label: '预估本周分成', value: '¥3,212', hint: '按 12% 口径' },
] as const;

export const merchants = [
  {
    city: '贵阳',
    commissionRate: '12%',
    industry: '火锅',
    merchant: '黔炉鲜牛火锅',
    nextAction: '周一补拍老板出镜 + 新品锅底',
    packageName: '双人团购 128 元',
    recentGmv: '¥9,860',
    stage: '稳定运营',
  },
  {
    city: '遵义',
    commissionRate: '10%',
    industry: '烘焙',
    merchant: '山野面包研究所',
    nextAction: '补 3 条七夕预热视频文案',
    packageName: '下午茶双人套餐 69 元',
    recentGmv: '¥4,220',
    stage: '冷启动第 2 周',
  },
] as const;

export const accountAccesses = [
  {
    account: '抖音号：qianlu_hotpot',
    lastLogin: '2026-08-01 09:30',
    merchant: '黔炉鲜牛火锅',
    method: '店家手机验证码代登',
    note: '仅登录抖音来客与创作服务平台',
    status: '已授权',
  },
  {
    account: '抖音号：breadlab_zunyi',
    lastLogin: '2026-07-31 20:15',
    merchant: '山野面包研究所',
    method: '商家子账号协作',
    note: '发布需老板二次确认',
    status: '待续期',
  },
] as const;

export const followUpLogs = [
  {
    latestTalk: '老板认可先不收费，只按核销额提点',
    merchant: '黔炉鲜牛火锅',
    nextStep: '确认 8 月选题池与拍摄时间',
    objection: '担心团购价格打太低伤利润',
    owner: '你',
    stage: '已签约',
  },
  {
    latestTalk: '已经开放后台，但希望先看 1 周效果',
    merchant: '山野面包研究所',
    nextStep: '先跑 3 条种草短视频',
    objection: '不想频繁出镜',
    owner: '运营小刘',
    stage: '试运营',
  },
] as const;

export const packages = [
  {
    cost: '¥72',
    margin: '43.8%',
    merchant: '黔炉鲜牛火锅',
    packageName: '双人鲜牛火锅套餐',
    profitGuard: '不建议再降价',
    sellingPrice: '¥128',
    trafficLabel: '门店引流主推',
  },
  {
    cost: '¥31',
    margin: '55.1%',
    merchant: '山野面包研究所',
    packageName: '贝果 + 咖啡双人下午茶',
    profitGuard: '可做节日限定加赠',
    sellingPrice: '¥69',
    trafficLabel: '适合做同城种草',
  },
] as const;

export const settlementOrders = [
  {
    commission: '¥1,183',
    commissionRate: '12%',
    merchant: '黔炉鲜牛火锅',
    orderWindow: '07/25 - 07/31',
    redeemedAmount: '¥9,860',
    sourceVideo: '老板试吃 128 元双人锅',
    status: '待对账',
  },
  {
    commission: '¥422',
    commissionRate: '10%',
    merchant: '山野面包研究所',
    orderWindow: '07/25 - 07/31',
    redeemedAmount: '¥4,220',
    sourceVideo: '打工人下班面包补给',
    status: '已确认',
  },
] as const;

export const benchmarks = [
  {
    account: '@贵阳火锅局',
    city: '贵阳',
    lane: '火锅探店',
    latestHit: '90 元吃到撑的老火锅',
    risk: '高频低价，容易伤毛利',
    takeaway: '老板出镜 + 价格锚点 + 排队实拍',
  },
  {
    account: '@面包控阿木',
    city: '遵义',
    lane: '烘焙日常',
    latestHit: '这家贝果像在东京街角',
    risk: '调性强，复制时不能太像',
    takeaway: '清晨出炉镜头 + 生活方式旁白',
  },
] as const;

export const topics = [
  {
    hook: '128 元双人火锅到底能吃到什么程度',
    merchant: '黔炉鲜牛火锅',
    publishWindow: '周二晚 19:30',
    source: '同城热门 + 老板招牌锅底',
    status: '待拍摄',
    topic: '价格锚点型探店',
  },
  {
    hook: '打工人下班 15 分钟买到的治愈系面包',
    merchant: '山野面包研究所',
    publishWindow: '周四 17:45',
    source: '七夕预热 + 下班场景',
    status: '待定稿',
    topic: '情绪价值型种草',
  },
] as const;

export const scripts = [
  {
    cta: '评论区回复“火锅”领取到店提醒',
    merchant: '黔炉鲜牛火锅',
    opening: '先给你看账单，再决定这锅值不值',
    scriptTitle: '老板试吃 128 元双人锅',
    status: '已过审',
  },
  {
    cta: '私信“贝果”发你当天出炉时间',
    merchant: '山野面包研究所',
    opening: '下班路过这家店，空气都是黄油味',
    scriptTitle: '打工人下班面包补给',
    status: '待老板确认',
  },
] as const;

export const storyboards = [
  {
    dialogue: '第一句先报价格，第二句给份量反差',
    lens: '门头远景 -> 锅底特写 -> 老板出镜',
    location: '门口、明档、二楼卡座',
    scene: '价格锚点型探店',
    scriptTitle: '老板试吃 128 元双人锅',
  },
  {
    dialogue: '用第一视角说“今天真的太累了”',
    lens: '出炉拉开 -> 切面特写 -> 打包离店',
    location: '烘焙台、收银台、门口橱窗',
    scene: '情绪价值型种草',
    scriptTitle: '打工人下班面包补给',
  },
] as const;

export const publishSchedule = [
  {
    merchant: '黔炉鲜牛火锅',
    owner: '运营小刘',
    publishTime: '2026-08-05 19:30',
    status: '待发布',
    title: '老板试吃 128 元双人锅',
  },
  {
    merchant: '山野面包研究所',
    owner: '你',
    publishTime: '2026-08-07 17:45',
    status: '已排期',
    title: '打工人下班面包补给',
  },
] as const;

export const systemRules = [
  {
    label: '合作模式',
    value: '前期不收服务费，按抖音相关核销额提 10% - 15%',
  },
  {
    label: '结算口径',
    value: '优先按核销额结算，退款单与未核销单不计入分成',
  },
  {
    label: '人工兜底',
    value: '价格调整、发布确认、账号登录、结算确认都必须人工审批',
  },
] as const;

export const roleCards = [
  {
    focus: '谈商家、签约、定套餐边界、收素材',
    role: '前线 BD / 主理人',
  },
  {
    focus: '账号诊断、对标分析、选题、文案、分镜、排期',
    role: '运营人员',
  },
  {
    focus: '按分镜执行拍摄、粗剪、复核、发布',
    role: '拍摄剪辑',
  },
] as const;

export const shootTasks = [
  {
    assignee: '拍摄阿豪',
    deadline: '08-04 15:00',
    merchant: '黔炉鲜牛火锅',
    shotList: '门头、锅底、老板试吃、翻台镜头',
    status: '待执行',
    taskName: '补拍火锅双人锅素材',
  },
  {
    assignee: '拍摄小米',
    deadline: '08-06 10:00',
    merchant: '山野面包研究所',
    shotList: '出炉、切面、打包、街景',
    status: '已完成',
    taskName: '七夕面包预热素材',
  },
] as const;

export const reviewRows = [
  {
    deals: '18 单',
    hook: '先给你看账单，再决定值不值',
    merchant: '黔炉鲜牛火锅',
    nextAction: '补老板口播，强化“份量感”镜头',
    ordersAmount: '¥3,296',
    video: '老板试吃 128 元双人锅',
    views: '1.8 万',
  },
  {
    deals: '9 单',
    hook: '下班路过这家店，空气都是黄油味',
    merchant: '山野面包研究所',
    nextAction: '放大“下班治愈”情绪标签',
    ordersAmount: '¥621',
    video: '打工人下班面包补给',
    views: '6,400',
  },
] as const;
