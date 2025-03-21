package model

var (
	PromptSummarizeText = `请用markdown列表格式**详细**总结我发送给你的文本。充分合理使用缩进和子列表，如果有需要可以使用多层子列表，或是在子列表中包含多个条目（3个或以上）。在每个总结项开头，用简短的词语描述该项。忽略和文章主体无关的内容（如广告）。无论原文语言为何，总是使用中文进行总结。

示例如下：

1. **Trello简介**：Trello是Fog Creek Software推出的100%基于云端的团队协作工具，自发布以来获得了积极的反馈和强劲的用户增长。

2. **开发模式转变**：Trello的开发标志着Fog Creek转向完全的云服务，不再提供安装版软件，开发过程中未使用Visual Basic，体现了开发流程的现代化。

3. **产品定位**：Trello是一款横跨多行业的产品，与之前主要针对软件开发者的垂直产品不同，它适用于各行各业的用户。

4. **横纵对比**：
   - **横向产品**：适用于广泛用户群体，如Word处理器和Web浏览器，难以定价过高，风险与回报并存。
   - **垂直产品**：针对特定行业，如牙医软件，用户定位明确，利润空间大，适合初创企业。

5. **Excel故事**：通过Excel的使用案例说明，大多数用户使用Excel实际上是作为列表工具，而非复杂的计算，引出“杀手级应用实际上是高级数据结构”的观点。

6. **Trello的核心**：Trello是一个高度灵活的数据结构应用，不仅限于敏捷开发的Kanban板，适用于规划婚礼、管理招聘流程等多种场景。

7. **产品特性**：
   - **持续交付**：新功能不断推出，无重大或次要版本的区别。
   - **快速迭代与修复**：测试不求面面俱到，但快速响应修复。
   - **公共透明**：开发过程公开，用户可参与反馈和投票。
   - **快速扩张策略**：目标是大规模用户增长，初期免费，优先消除采用障碍。
   - **API优先**：鼓励通过API和插件扩展功能，用户和第三方参与建设。

8. **技术选择**：采用前沿技术如MongoDB、WebSockets、CoffeeScript和Node.js，虽然有挑战，但有利于吸引顶尖程序员并为长期发展做准备。

9. **总结**：Trello及其开发策略体现了现代互联网产品的开发趋势，注重用户基础的快速扩展，技术的前沿性，以及通过社区参与和反馈来不断优化产品。`

	PromptOneSentenceSummary = "以下是对一篇长文的列表形式总结。请基于此输出对该文章的简短总结，长度不超过100个字。总是使用简体中文输出。"

	PromptPodcast = `请生成一份适合语音播报的Hacker News晨间新闻简报，要求如下：

播报风格
采用BBC式专业新闻播报口吻
语句简洁流畅，避免复杂术语
重要新闻优先排序（按原文列表顺序）
每条新闻包含：标题核心内容、关键细节
适当加入过渡词如"接下来关注..." "最后来看..."
内容结构
[开场白]
"早上好，今天是[日期]，现在是Hacker News技术晨报时间。以下是今日值得关注的要闻："
[新闻主体]

对每个条目按此模板处理：

"• [核心标题]：用1句话概括新闻主题（例：政府技术团队18F突遭解散）

详细内容：提炼原文3个关键要素（人物/事件/影响/技术亮点）+ 补充必要背景说明"

[结束语]

"以上就是今日精选的技术要闻，您可通过Hacker News获取完整讨论。祝您今日工作顺利，我们明天见。"

特别要求
10条新闻完整保留，不可遗漏
将技术细节转化为通俗表达（例：将"GLP-1药物经济影响"转化为"糖尿病药物对全球经济产生连锁反应"）
涉及机构名称首次出现时加注解（例：18F——美国政府数字服务团队）
中文播报但保留必要英文缩写（如AI/GPT/GSA）
每条新闻时长控制在20-40秒
示例输出
（以下为第一条新闻处理示例）
"• 政府技术团队18F突遭解散：非党派组织18F在3月1日被总务管理局突然关闭，90名员工遭解雇。这个被前特斯拉工程师称为'公民技术黄金标准'的团队，曾主导多项政府数字化项目。此举被视为削减IT支持体系的一部分，引发技术社区强烈反响。"`
)

type OpenAIResponse struct {
	Choices []*Choice `json:"choices"`
}

type Choice struct {
	Message *Message `json:"message"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}
