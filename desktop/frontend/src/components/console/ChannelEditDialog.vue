<script setup lang="ts">
import { ref, reactive, computed, watch } from 'vue'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Textarea } from '@/components/ui/textarea'
import { Switch } from '@/components/ui/switch'
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select'
import { ScrollArea } from '@/components/ui/scroll-area'
import { Loader2, X, ChevronDown, ChevronUp } from 'lucide-vue-next'
import { useConsoleChannels } from '@/composables/useConsoleChannels'
import { useLanguage } from '@/composables/useLanguage'
import { buildChannelPayload } from '@/utils/channel-payload'
import { parseQuickInput } from '@/utils/quick-input-parser'
import type { Channel } from '@/services/admin-api'

interface Props {
  channel?: Channel | null
  channelType: string
}

const props = defineProps<Props>()
const emit = defineEmits<{
  (e: 'close'): void
  (e: 'saved'): void
}>()

const { tf } = useLanguage()
const { saveChannel } = useConsoleChannels()

const isEditMode = computed(() => !!props.channel)
const saving = ref(false)
const error = ref('')

// 展开高级选项
const showAdvanced = ref(false)

// 表单状态
const form = reactive({
  name: '',
  description: '',
  serviceType: '' as 'openai' | 'claude' | 'gemini' | 'responses' | '',
  baseUrl: '',
  baseUrlsText: '',
  website: '',
  proxyUrl: '',
  routePrefix: '',
  insecureSkipVerify: false,
  apiKeysText: '',
  customHeadersText: '{}',
  modelMappingText: '{}',
  supportedModelsText: '',
  // Flags
  noVision: false,
  passbackReasoningContent: false,
  passbackThinkingBlocks: false,
  fastMode: false,
  lowQuality: false,
  codexToolCompat: false,
  stripEmptyTextBlocks: false,
  normalizeSystemRoleToTopLevel: false,
})

// 从编辑频道填充表单
watch(() => props.channel, (ch) => {
  if (!ch) return
  form.name = ch.name || ''
  form.description = ch.description || ''
  form.serviceType = ch.serviceType || ''
  form.baseUrl = ch.baseUrl || ''
  form.baseUrlsText = (ch.baseUrls || []).join('\n')
  form.website = ch.website || ''
  form.proxyUrl = ch.proxyUrl || ''
  form.routePrefix = ch.routePrefix || ''
  form.insecureSkipVerify = ch.insecureSkipVerify || false
  form.apiKeysText = (ch.apiKeys || []).join('\n')
  form.customHeadersText = ch.customHeaders ? JSON.stringify(ch.customHeaders, null, 2) : '{}'
  form.modelMappingText = ch.modelMapping ? JSON.stringify(ch.modelMapping, null, 2) : '{}'
  form.supportedModelsText = (ch.supportedModels || []).join('\n')
  form.noVision = ch.noVision || false
  form.passbackReasoningContent = ch.passbackReasoningContent || false
  form.passbackThinkingBlocks = ch.passbackThinkingBlocks || false
  form.fastMode = ch.fastMode || false
  form.lowQuality = ch.lowQuality || false
  form.codexToolCompat = ch.codexToolCompat || false
  form.stripEmptyTextBlocks = ch.stripEmptyTextBlocks || false
  form.normalizeSystemRoleToTopLevel = ch.normalizeSystemRoleToTopLevel || false
}, { immediate: true })

// 校验
const errors = computed(() => {
  const errs: Record<string, string> = {}
  if (!form.name.trim()) errs.name = tf('console.form.nameRequired', '频道名称必填')
  if (!form.serviceType) errs.serviceType = tf('console.form.serviceTypeRequired', '请选择服务类型')
  if (!form.baseUrl.trim() && !form.baseUrlsText.trim()) errs.baseUrl = tf('console.form.baseUrlRequired', '至少需要一个 Base URL')
  return errs
})

const isValid = computed(() => Object.keys(errors.value).length === 0)

// 快速粘贴解析
function handleQuickPaste(text: string) {
  const result = parseQuickInput(text, form.serviceType || undefined)
  if (result.detectedBaseUrl) form.baseUrl = result.detectedBaseUrl
  if (result.detectedBaseUrls.length > 1) {
    form.baseUrlsText = result.detectedBaseUrls.join('\n')
  }
  if (result.detectedApiKeys.length) {
    form.apiKeysText = result.detectedApiKeys.join('\n')
  }
  if (result.detectedServiceType && !form.serviceType) {
    form.serviceType = result.detectedServiceType
  }
}

// 提交
async function handleSubmit() {
  if (!isValid.value) return
  saving.value = true
  error.value = ''

  try {
    // 解析 JSON 字段
    let customHeaders: Record<string, string> = {}
    let modelMapping: Record<string, string> = {}
    try { customHeaders = JSON.parse(form.customHeadersText) } catch { /* ignore */ }
    try { modelMapping = JSON.parse(form.modelMappingText) } catch { /* ignore */ }

    const baseUrls = form.baseUrlsText
      .split('\n')
      .map(s => s.trim())
      .filter(Boolean)

    const apiKeys = form.apiKeysText
      .split('\n')
      .map(s => s.trim())
      .filter(Boolean)

    const supportedModels = form.supportedModelsText
      .split('\n')
      .map(s => s.trim())
      .filter(Boolean)

    const payload = buildChannelPayload({
      name: form.name,
      serviceType: form.serviceType as any,
      baseUrl: form.baseUrl,
      baseUrls,
      website: form.website,
      insecureSkipVerify: form.insecureSkipVerify,
      lowQuality: form.lowQuality,
      injectDummyThoughtSignature: false,
      stripThoughtSignature: false,
      passbackReasoningContent: form.passbackReasoningContent,
      passbackThinkingBlocks: form.passbackThinkingBlocks,
      description: form.description,
      apiKeys,
      modelMapping,
      reasoningMapping: {},
      reasoningParamStyle: 'reasoning',
      textVerbosity: '',
      fastMode: form.fastMode,
      customHeaders,
      proxyUrl: form.proxyUrl,
      routePrefix: form.routePrefix,
      supportedModels,
      autoBlacklistBalance: true,
      normalizeMetadataUserId: true,
      stripEmptyTextBlocks: form.stripEmptyTextBlocks,
      normalizeSystemRoleToTopLevel: form.normalizeSystemRoleToTopLevel,
      codexNativeToolPassthrough: false,
      codexToolCompat: form.codexToolCompat,
      noVision: form.noVision,
      noVisionModels: [],
      visionFallbackModel: '',
    })

    await saveChannel(payload, props.channel?.index ?? null)
    emit('saved')
    emit('close')
  } catch (e) {
    error.value = e instanceof Error ? e.message : String(e)
  } finally {
    saving.value = false
  }
}

function onKeyDown(e: KeyboardEvent) {
  if (e.key === 'Escape') emit('close')
  if (e.key === 'Enter' && (e.metaKey || e.ctrlKey)) handleSubmit()
}
</script>

<template>
  <Teleport to="body">
    <Transition name="fade">
      <div
        v-if="true"
        class="fixed inset-0 z-50 flex items-center justify-center"
        @keydown="onKeyDown"
      >
        <div class="absolute inset-0 bg-black/60 backdrop-blur-sm" @click="emit('close')" />

        <div class="relative z-10 w-[90vw] max-w-2xl max-h-[85vh] rounded-2xl border border-border bg-card shadow-2xl flex flex-col">
          <!-- Header -->
          <div class="flex items-center justify-between p-4 border-b border-border shrink-0">
            <h3 class="text-sm font-semibold">
              {{ isEditMode
                ? tf('console.form.editChannel', '编辑频道')
                : tf('console.form.addChannel', '添加频道')
              }}
            </h3>
            <Button variant="ghost" size="icon-sm" @click="emit('close')">
              <X class="w-4 h-4" />
            </Button>
          </div>

          <!-- Body -->
          <ScrollArea class="flex-1 min-h-0">
            <form @submit.prevent="handleSubmit" class="p-4 space-y-5">
              <!-- Error -->
              <div v-if="error" class="text-sm text-destructive bg-destructive/10 rounded-lg p-3">
                {{ error }}
              </div>

              <!-- Section: Basic Info -->
              <div class="space-y-3">
                <h4 class="text-xs font-semibold text-muted-foreground uppercase tracking-wider">
                  {{ tf('console.form.basicInfo', '基础信息') }}
                </h4>
                <div class="grid grid-cols-2 gap-3">
                  <div class="space-y-1.5">
                    <Label>{{ tf('console.form.name', '名称') }} *</Label>
                    <Input v-model="form.name" :class="{ 'border-destructive': errors.name }" />
                    <p v-if="errors.name" class="text-[10px] text-destructive">{{ errors.name }}</p>
                  </div>
                  <div class="space-y-1.5">
                    <Label>{{ tf('console.form.serviceType', '服务类型') }} *</Label>
                    <Select v-model="form.serviceType">
                      <SelectTrigger :class="{ 'border-destructive': errors.serviceType }">
                        <SelectValue :placeholder="tf('console.form.selectServiceType', '选择服务类型')" />
                      </SelectTrigger>
                      <SelectContent>
                        <SelectItem value="claude">Claude</SelectItem>
                        <SelectItem value="openai">OpenAI</SelectItem>
                        <SelectItem value="gemini">Gemini</SelectItem>
                        <SelectItem value="responses">Responses</SelectItem>
                      </SelectContent>
                    </Select>
                    <p v-if="errors.serviceType" class="text-[10px] text-destructive">{{ errors.serviceType }}</p>
                  </div>
                </div>
                <div class="space-y-1.5">
                  <Label>{{ tf('console.form.description', '描述') }}</Label>
                  <Textarea v-model="form.description" rows="2" />
                </div>
              </div>

              <!-- Section: Connection -->
              <div class="space-y-3">
                <h4 class="text-xs font-semibold text-muted-foreground uppercase tracking-wider">
                  {{ tf('console.form.connection', '连接') }}
                </h4>
                <div class="space-y-1.5">
                  <Label>{{ tf('console.form.baseUrl', 'Base URL') }} *</Label>
                  <Input v-model="form.baseUrl" placeholder="https://api.example.com" :class="{ 'border-destructive': errors.baseUrl }" />
                  <p v-if="errors.baseUrl" class="text-[10px] text-destructive">{{ errors.baseUrl }}</p>
                </div>
                <div class="space-y-1.5">
                  <Label>{{ tf('console.form.additionalUrls', '额外 URL（每行一个）') }}</Label>
                  <Textarea v-model="form.baseUrlsText" rows="2" placeholder="https://backup.example.com" />
                </div>
                <div class="grid grid-cols-2 gap-3">
                  <div class="space-y-1.5">
                    <Label>{{ tf('console.form.proxyUrl', '代理 URL') }}</Label>
                    <Input v-model="form.proxyUrl" placeholder="socks5://..." />
                  </div>
                  <div class="space-y-1.5">
                    <Label>{{ tf('console.form.routePrefix', '路由前缀') }}</Label>
                    <Input v-model="form.routePrefix" placeholder="kimi" />
                  </div>
                </div>
                <div class="flex items-center gap-2">
                  <Switch v-model="form.insecureSkipVerify" />
                  <Label class="text-xs">{{ tf('console.form.insecureSkipVerify', '跳过 TLS 验证') }}</Label>
                </div>
              </div>

              <!-- Section: Auth -->
              <div class="space-y-3">
                <h4 class="text-xs font-semibold text-muted-foreground uppercase tracking-wider">
                  {{ tf('console.form.authentication', '认证') }}
                </h4>
                <div class="space-y-1.5">
                  <Label>{{ tf('console.form.apiKeys', 'API Keys（每行一个）') }} *</Label>
                  <Textarea v-model="form.apiKeysText" rows="3" placeholder="sk-xxx&#10;sk-yyy" class="font-mono text-xs" />
                </div>
                <div class="space-y-1.5">
                  <Label>{{ tf('console.form.customHeaders', '自定义 Headers（JSON）') }}</Label>
                  <Textarea v-model="form.customHeadersText" rows="2" class="font-mono text-xs" />
                </div>
              </div>

              <!-- Section: Models -->
              <div class="space-y-3">
                <h4 class="text-xs font-semibold text-muted-foreground uppercase tracking-wider">
                  {{ tf('console.form.models', '模型') }}
                </h4>
                <div class="space-y-1.5">
                  <Label>{{ tf('console.form.modelMapping', '模型映射（JSON）') }}</Label>
                  <Textarea v-model="form.modelMappingText" rows="2" class="font-mono text-xs" />
                </div>
                <div class="space-y-1.5">
                  <Label>{{ tf('console.form.supportedModels', '支持的模型（每行一个，留空=全部）') }}</Label>
                  <Textarea v-model="form.supportedModelsText" rows="2" placeholder="gpt-4*&#10;claude-3*" class="font-mono text-xs" />
                </div>
              </div>

              <!-- Section: Flags (collapsible) -->
              <div class="space-y-3">
                <button
                  type="button"
                  class="flex items-center gap-1.5 text-xs font-semibold text-muted-foreground uppercase tracking-wider hover:text-foreground transition-colors"
                  @click="showAdvanced = !showAdvanced"
                >
                  <ChevronDown v-if="!showAdvanced" class="w-3.5 h-3.5" />
                  <ChevronUp v-else class="w-3.5 h-3.5" />
                  {{ tf('console.form.advancedFlags', '高级选项') }}
                </button>
                <div v-if="showAdvanced" class="grid grid-cols-2 gap-2 pl-2">
                  <div v-for="flag in [
                    { key: 'noVision', label: tf('console.form.noVision', '禁用视觉') },
                    { key: 'passbackReasoningContent', label: tf('console.form.passbackReasoning', '回传推理内容') },
                    { key: 'passbackThinkingBlocks', label: tf('console.form.passbackThinking', '回传思考块') },
                    { key: 'fastMode', label: tf('console.form.fastMode', '快速模式') },
                    { key: 'lowQuality', label: tf('console.form.lowQuality', '低质量标记') },
                    { key: 'codexToolCompat', label: tf('console.form.codexCompat', 'Codex 工具兼容') },
                    { key: 'stripEmptyTextBlocks', label: tf('console.form.stripEmptyBlocks', '移除空文本块') },
                    { key: 'normalizeSystemRoleToTopLevel', label: tf('console.form.normalizeSystem', '规范化系统角色') },
                  ]" :key="flag.key" class="flex items-center gap-2">
                    <Switch :model-value="(form as any)[flag.key]" @update:model-value="(v: boolean) => (form as any)[flag.key] = v" />
                    <Label class="text-xs">{{ flag.label }}</Label>
                  </div>
                </div>
              </div>
            </form>
          </ScrollArea>

          <!-- Footer -->
          <div class="flex items-center justify-end gap-2 p-4 border-t border-border shrink-0">
            <Button variant="ghost" @click="emit('close')">
              {{ tf('common.cancel', '取消') }}
            </Button>
            <Button :disabled="!isValid || saving" @click="handleSubmit">
              <Loader2 v-if="saving" class="w-4 h-4 mr-2 animate-spin" />
              {{ isEditMode
                ? tf('console.form.save', '保存')
                : tf('console.form.create', '创建')
              }}
            </Button>
          </div>
        </div>
      </div>
    </Transition>
  </Teleport>
</template>

<style scoped>
.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.15s ease;
}
.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}
</style>
