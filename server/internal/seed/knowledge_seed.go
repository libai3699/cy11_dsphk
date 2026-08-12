package seed

import (
	"cy11dsphk/server/internal/model"

	"gorm.io/gorm"
)

// SeedKnowledge 初始化「运营知识库」模块的种子数据。
// 自包含 AutoMigrate，不改动 database.MigrateAndSeed 现有流程；
// 使用 FirstOrCreate 按自然键幂等插入，重复执行不会重复写入。
func SeedKnowledge(db *gorm.DB) error {
	if err := db.AutoMigrate(
		&model.PainPoint{},
		&model.CaseStudy{},
		&model.AccountProfile{},
		&model.PlatformRule{},
		&model.ContentTemplate{},
	); err != nil {
		return err
	}

	if err := seedPainPoints(db); err != nil {
		return err
	}
	if err := seedCaseStudies(db); err != nil {
		return err
	}
	if err := seedAccountProfiles(db); err != nil {
		return err
	}
	if err := seedPlatformRules(db); err != nil {
		return err
	}
	return seedContentTemplates(db)
}

func seedPainPoints(db *gorm.DB) error {
	items := []model.PainPoint{
		{
			Source:         "评论区",
			Category:       "内容选题",
			Content:        "做了半年抖音，发了不少视频但播放一直卡在两三百，不知道问题出在哪。",
			UserQuote:      "为什么别人随便发都有几千播放，我精心剪辑的却没人看？",
			Emotion:        "困惑",
			Product:        "账号诊断服务",
			DemandLevel:    "高",
			SuggestedTopic: "新手账号播放低的5个常见原因（误区类）",
			Status:         model.PainPointStatusPending,
		},
		{
			Source:         "私信",
			Category:       "投放转化",
			Content:        "视频有播放有赞但没有咨询和成交，挂了车也没人点。",
			UserQuote:      "流量有了就是不转化，是不是我话术有问题？",
			Emotion:        "焦虑",
			Product:        "带货/知识付费",
			DemandLevel:    "高",
			SuggestedTopic: "口播如何把『看完』变成『下单』（转化话术拆解）",
			Status:         model.PainPointStatusAdopted,
		},
		{
			Source:         "群聊",
			Category:       "账号定位",
			Content:        "不知道账号该面向谁、卖什么，今天发美食明天发干货，粉丝很乱。",
			UserQuote:      "定位到底怎么定？感觉什么都能发又什么都不垂直。",
			Emotion:        "困惑",
			Product:        "账号定位策划",
			DemandLevel:    "中",
			SuggestedTopic: "用倒推法3步确定你的账号定位（实操类）",
			Status:         model.PainPointStatusPending,
		},
		{
			Source:         "热点",
			Category:       "拍摄剪辑",
			Content:        "开头总是拖沓，用户3秒就划走，完播率上不去。",
			UserQuote:      "怎么才能让开头不废话、直接抓住人？",
			Emotion:        "期待",
			Product:        "脚本服务",
			DemandLevel:    "高",
			SuggestedTopic: "8种开头钩子类型库+怎么选（揭秘/对比）",
			Status:         model.PainPointStatusAdopted,
		},
		{
			Source:         "对标评论区",
			Category:       "合规",
			Content:        "担心踩平台红线被限流，尤其涉及功效/对比词时不敢说。",
			UserQuote:      "说『最有效』会不会被判定夸大？医疗类要什么资质？",
			Emotion:        "焦虑",
			Product:        "合规咨询",
			DemandLevel:    "中",
			SuggestedTopic: "抖音内容合规红线清单（避坑类）",
			Status:         model.PainPointStatusPending,
		},
	}
	for _, it := range items {
		if err := db.Where("content = ? AND source = ?", it.Content, it.Source).
			Assign(it).FirstOrCreate(&model.PainPoint{}).Error; err != nil {
			return err
		}
	}
	return nil
}

func seedCaseStudies(db *gorm.DB) error {
	items := []model.CaseStudy{
		{
			Title:            "对标账号『XX说运营』误区类口播爆款拆解",
			Platform:         "抖音",
			AccountName:      "XX说运营",
			Industry:         "自媒体运营/知识付费",
			Form:             "口播",
			HookType:         "误区类",
			Structure:        "痛点(你也这样?)→价值(其实错在X)→行动(关注领模板)",
			ConversionAction: "评论置顶引导私信领资料",
			ReusablePoint:    "用『90%的人都做错了』制造认知冲突，前3秒直接抛反常识结论",
			DataJSON:         `{"playCount":580000,"likeCount":32000,"commentCount":4100,"shareCount":2600,"favoriteCount":9800}`,
			Lesson:           "反常识开头+评论区置顶资料，私信量显著提升",
			Status:           model.CaseStudyStatusActive,
		},
		{
			Title:            "自身账号『新手起步』故事类口播复盘",
			Platform:         "抖音",
			AccountName:      "新手起步",
			Industry:         "个人成长",
			Form:             "口播",
			HookType:         "故事类",
			Structure:        "故事(我踩过的坑)→价值(3个方法)→行动(点主页看合集)",
			ConversionAction: "主页合集+私域引流",
			ReusablePoint:    "真实失败故事比干货更易引发共鸣，情绪曲线先抑后扬",
			DataJSON:         `{"playCount":210000,"likeCount":15000,"commentCount":2200,"shareCount":1300,"favoriteCount":5400}`,
			Lesson:           "真人出镜+真实故事信任度高，适合冷启动",
			Status:           model.CaseStudyStatusActive,
		},
		{
			Title:            "对标账号『选品实验室』测评类爆款拆解",
			Platform:         "抖音",
			AccountName:      "选品实验室",
			Industry:         "好物测评/带货",
			Form:             "测评",
			HookType:         "对比类",
			Structure:        "痛点(怕踩雷)→价值(实测3款)→行动(挂车下单)",
			ConversionAction: "挂车+直播间",
			ReusablePoint:    "用『实测前后对比』建立信任，结尾强CTA挂车",
			DataJSON:         `{"playCount":430000,"likeCount":28000,"commentCount":3600,"shareCount":1900,"favoriteCount":7200}`,
			Lesson:           "测评类适合带货，实测画面+数据增强可信度",
			Status:           model.CaseStudyStatusActive,
		},
	}
	for _, it := range items {
		if err := db.Where("title = ? AND account_name = ?", it.Title, it.AccountName).
			Assign(it).FirstOrCreate(&model.CaseStudy{}).Error; err != nil {
			return err
		}
	}
	return nil
}

func seedAccountProfiles(db *gorm.DB) error {
	items := []model.AccountProfile{
		{
			AccountName:     "示例·知识付费主理人",
			Platform:        "抖音",
			Industry:        "自媒体运营教学",
			TargetAudience:  "25-40岁、想做/刚做自媒体的个体创业者与小团队，常搜『抖音怎么做』『账号定位』『涨粉』",
			Product:         "自媒体运营教学课程+账号诊断服务+私域陪跑；变现方式：知识付费/引流私域；客单价299-1999",
			Positioning:     "用倒推法帮新手把账号做起来的实操派",
			Persona:         "专业+亲和，像带你的学长",
			ToneStyle:       "口语化、少术语、给步骤",
			Monetization:    "知识付费+私域引流",
			UpdateFrequency: "每周3-4条",
			Strengths:       "方法论清晰、案例多",
			Weaknesses:      "出镜表现力待提升",
			Status:          model.AccountProfileStatusActive,
		},
	}
	for _, it := range items {
		if err := db.Where("account_name = ? AND platform = ?", it.AccountName, it.Platform).
			Assign(it).FirstOrCreate(&model.AccountProfile{}).Error; err != nil {
			return err
		}
	}
	return nil
}

func seedPlatformRules(db *gorm.DB) error {
	items := []model.PlatformRule{
		{
			Platform:  "抖音",
			Category:  "合规红线",
			Title:     "内容合规红线：禁止虚假宣传与绝对化用语",
			Content:   "禁止虚假宣传、夸大功效、使用绝对化用语（『最』『第一』『100%有效』等）；医疗、金融、减肥、祛斑等类目需资质与额外审查；诱导互动（『点赞才能看』）与诱导站外引流会被限流；搬运、未授权素材存在版权风险。",
			RiskLevel: "高",
			Source:    "抖音社区自律公约/平台规则",
			Status:    model.PlatformRuleStatusActive,
		},
	}
	for _, it := range items {
		if err := db.Where("title = ? AND category = ?", it.Title, it.Category).
			Assign(it).FirstOrCreate(&model.PlatformRule{}).Error; err != nil {
			return err
		}
	}
	return nil
}

func seedContentTemplates(db *gorm.DB) error {
	items := []model.ContentTemplate{
		{
			Name:     "选题卡模板",
			Type:     "选题卡",
			Category: "选题中心",
			StructureJSON: `{"fields":["选题ID","标题","人群","产品/变现","核心痛点","痛点来源","内容角度","开头钩子","转化动作","预估转化潜力","预估播放潜力","状态"]}`,
			Content:  "# 选题卡模板\n\n- 选题ID：T001\n- 标题（暂定）：\n- 人群：（来自 account-profile）\n- 产品/变现：（挂钩什么产品）\n- 核心痛点：（用户愿意掏钱的痛点）\n- 痛点来源：热点/对标评论区/本账号评论区/私信/群聊\n- 内容角度：误区/对比/揭秘/实操/故事\n- 开头钩子（3秒）：\n- 转化动作：（评论置顶/私信/挂车/引流）\n- 预估转化潜力：高/中/低\n- 预估播放潜力：高/中/低\n- 状态：待脚本/已拍/已发/废弃",
			UsageNote: "由 01 选题策划填写，必须含【人群/产品/痛点/转化动作/钩子】。",
			Status:    model.ContentTemplateStatusActive,
		},
		{
			Name:     "口播脚本模板",
			Type:     "口播脚本",
			Category: "文案脚本",
			StructureJSON: `{"sections":["钩子0-3s","痛点","价值/方法","信任背书","行动号召"]}`,
			Content:  "# 口播脚本模板\n\n- 选题ID：T001\n- 标题（A版）：\n- 标题（B版）：\n- 话题标签：# # #\n- 正文脚本：\n  - 【钩子 0-3s】：\n  - 【痛点】：\n  - 【价值/方法】：\n  - 【信任背书】：\n  - 【行动号召】：（对应转化动作）\n- 文案（发布用）：\n- 备注/拍摄提示：",
			UsageNote: "由 02 脚本创作填写，结构：钩子→痛点→价值→信任→行动。",
			Status:    model.ContentTemplateStatusActive,
		},
		{
			Name:     "分镜/剪辑脚本模板",
			Type:     "分镜脚本",
			Category: "分镜脚本",
			StructureJSON: `{"shotTable":["镜头","景别","机位","动作/画面","台词","时长"],"editing":["转场","字幕","BGM","节奏点/黄金3秒"]}`,
			Content:  "# 分镜 / 剪辑脚本模板\n\n- 选题ID/脚本ID：T001 / S001\n- 总时长：\n- 分镜表：\n  | 镜头 | 景别 | 机位 | 动作/画面 | 台词 | 时长 |\n  |---|---|---|---|---|---|\n  | 1 | | | | | |\n- 剪辑脚本：\n  - 转场：\n  - 字幕（与口播对齐）：\n  - BGM：\n  - 节奏点/黄金3秒：\n- 封面规范：（构图/大字标题/人物表情）\n- 备注：",
			UsageNote: "由 03 拍摄剪辑填写，主理人据此出镜拍摄。",
			Status:    model.ContentTemplateStatusActive,
		},
	}
	for _, it := range items {
		if err := db.Where("name = ? AND type = ?", it.Name, it.Type).
			Assign(it).FirstOrCreate(&model.ContentTemplate{}).Error; err != nil {
			return err
		}
	}
	return nil
}
