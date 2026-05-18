# 小米 MiMo 配置指南

## 获取 API Key

MiMo 模型可通过以下平台访问：
- [硅基流动 SiliconFlow](https://cloud.siliconflow.cn/)（推荐）
- [小米 MiMo 官网](https://mimo.xiaomi.com/)

### 通过硅基流动获取

1. 访问 [硅基流动](https://cloud.siliconflow.cn/)
2. 注册并登录账号
3. 进入「API Keys」页面
4. 创建新的 API Key 并复制

## 在 CCX 中添加渠道

### 通过硅基流动访问

| 字段 | 值 |
|------|-----|
| 名称 | `MiMo (SiliconFlow)`（自定义） |
| 服务类型 | `openai` |
| Base URL | `https://api.siliconflow.cn/v1` |
| API Keys | 你的 SiliconFlow API Key |

#### 配置步骤

1. 进入 CCX 管理界面，选择 **Chat** 入口
2. 点击「添加渠道」
3. 填写以下信息：
   - **名称**：`MiMo`
   - **服务类型**：选择 `OpenAI Chat`
   - **Base URL**：`https://api.siliconflow.cn/v1`
   - **API Keys**：粘贴你的 API Key
4. 点击保存

### 模型白名单（可选）

```
XiaomiMiMo/MiMo-V2.5-Pro
XiaomiMiMo/MiMo-V2.5
XiaomiMiMo/MiMo-V2-Flash
```

### 模型映射（可选）

```json
{
  "mimo-pro": "XiaomiMiMo/MiMo-V2.5-Pro",
  "mimo": "XiaomiMiMo/MiMo-V2.5",
  "mimo-flash": "XiaomiMiMo/MiMo-V2-Flash"
}
```

## 可用模型

| 模型 | 说明 |
|------|------|
| `XiaomiMiMo/MiMo-V2.5-Pro` | 最新旗舰，1.02T 总参 / 42B 激活 |
| `XiaomiMiMo/MiMo-V2.5` | 310B 总参 / 15B 激活，原生多模态 |
| `XiaomiMiMo/MiMo-V2-Flash` | 309B 总参 / 15B 激活，高速推理 |

::: tip
硅基流动上的模型 ID 格式为 `组织名/模型名`，如 `XiaomiMiMo/MiMo-V2.5-Pro`。使用时需要填写完整标识。
:::

## 注意事项

- MiMo 通过兼容 OpenAI 协议的平台访问
- 硅基流动国内 Base URL：`https://api.siliconflow.cn/v1`
- 硅基流动国际 Base URL：`https://api.siliconflow.com/v1`
- MiMo 是推理模型，支持 `reasoning_content` 字段返回思考过程
- 硅基流动也提供 Anthropic 兼容端点（`/anthropic/v1/messages`）

### 视觉支持

MiMo 各模型的视觉支持情况：

| 模型 | 视觉支持 |
|------|----------|
| `MiMo-V2.5-Pro` | 不支持 |
| `MiMo-V2.5` | 支持（原生多模态） |
| `MiMo-V2-Flash` | 不支持 |

::: warning
`MiMo-V2.5-Pro` 不支持图片输入。如果需要处理包含图片的请求，必须配置**视觉回退模型**。
:::

**配置方式：** 编辑渠道，在「视觉回退模型」字段填入 `MiMo-V2.5`。当请求包含图片且目标模型（如 `MiMo-V2.5-Pro`）不支持视觉时，CCX 会自动使用 `MiMo-V2.5` 替代模型处理该请求。

如果留空视觉回退模型，包含图片的请求将跳过该渠道，failover 到下一个支持视觉的渠道。

### 回传思考内容

MiMo 作为推理模型，思考过程中产生的 `reasoning_content` 需要回传给 API，否则会返回 HTTP 400 错误。

**必须启用：** 编辑渠道时打开「回传思考内容」开关（`PassbackReasoningContent`）。

启用后 CCX 会自动处理：
- **请求方向：** 为缺少 `thinking` 块的 assistant 消息注入占位符，满足 MiMo 的回传要求
- **响应方向：** 将上游返回的 `reasoning_content` 转换为 Claude 原生的 `thinking` 内容块，下游客户端可正常解析

::: warning
不开启此开关会导致包含历史对话的推理请求失败（HTTP 400）。
:::

### 推荐的模型映射

将 Claude 模型名映射到 MiMo 模型：

| 请求模型 | 重定向到 | 说明 |
|----------|----------|------|
| `haiku` | `mimo-v2.5-pro` | 旗舰推理 |
| `opus` | `mimo-v2.5-pro` | 旗舰推理 |
| `sonnet` | `mimo-v2.5-pro` | 旗舰推理 |

::: tip
所有 Claude 模型统一映射到 `MiMo-V2.5-Pro` 以获得最佳推理能力。配合视觉回退模型 `MiMo-V2.5` 使用，图片请求会自动切换到支持视觉的模型。
:::
