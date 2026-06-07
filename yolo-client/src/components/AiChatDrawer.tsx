import { LoadingOutlined } from '@ant-design/icons'
import { HistoryOutlined, PlusOutlined } from '@ant-design/icons'
import { Bubble, Sender, Conversations } from '@ant-design/x'
import XMarkdown from '@ant-design/x-markdown'
import '@ant-design/x-markdown/themes/light.css'
import '@ant-design/x-markdown/themes/dark.css'
import { Button, Collapse, Drawer, Popover, Space, Spin, Tag, message as antMsg } from 'antd'
import { useEffect, useRef, useState } from 'react'
import { useDeleteSession, useGetSessionMessages, useListSessions } from '../generated'

interface Session {
  session_id: string
  title?: string
}

interface ChatMessage {
  key: string
  role: 'user' | 'ai'
  content: string
  streaming?: boolean
  status?: 'loading' | 'success' | 'error' | 'abort'
}

function parseSSELine(line: string): { type: string; content?: string; session_id?: string } | null {
  if (!line.startsWith('data: ')) return null
  try {
    return JSON.parse(line.slice(6).trim())
  } catch {
    return null
  }
}

interface AiChatDrawerProps {
  open: boolean
  onClose: () => void
}

const LS_KEY = 'ai_active_session'

function loadSessionFromLS(): string | undefined {
  try {
    return localStorage.getItem(LS_KEY) || undefined
  } catch { return undefined }
}

function saveSessionToLS(sessionId: string | undefined) {
  try {
    if (sessionId) localStorage.setItem(LS_KEY, sessionId)
    else localStorage.removeItem(LS_KEY)
  } catch { /* ignore */ }
}

export default function AiChatDrawer({ open, onClose }: AiChatDrawerProps) {
  const [activeKey, setActiveKey] = useState<string | undefined>(loadSessionFromLS)
  const [messagesMap, setMessagesMap] = useState<Record<string, ChatMessage[]>>({})
  const [isLoading, setIsLoading] = useState(false)
  const [inputValue, setInputValue] = useState('')
  const [sessionPopoverOpen, setSessionPopoverOpen] = useState(false)
  const abortRef = useRef<AbortController | null>(null)
  const [isDark] = useState(() =>
    typeof window !== 'undefined' && window.matchMedia('(prefers-color-scheme: dark)').matches,
  )
  const xmClassName = isDark ? 'x-markdown-dark' : 'x-markdown-light'

  const deleteSession = useDeleteSession({})
  const doDeleteSession = deleteSession.mutate

  // 会话列表：组件挂载时自动获取，通过 refetch 手动刷新
  const { data: sessionsData, refetch: refetchSessions } = useListSessions({
    lazy: false,
  })
  const sessionsFromHook = (sessionsData as any)?.data || []

  useEffect(() => {
    if (open) refetchSessions()
  }, [open])

  // 持久化当前会话到 localStorage
  useEffect(() => { saveSessionToLS(activeKey) }, [activeKey])

  // 读取会话消息：基于 pendingSessionId 驱动
  const [pendingSessionId, setPendingSessionId] = useState<string | null>(() => {
    const cached = loadSessionFromLS()
    return cached || null
  })
  const [fetchKey, setFetchKey] = useState(0)
  const { data: messagesData, refetch: refetchMessages } = useGetSessionMessages({
    session_id: pendingSessionId || '',
    lazy: true,
  })
  useEffect(() => {
    if (pendingSessionId) {
      refetchMessages().then(() => setFetchKey((k) => k + 1))
    }
  }, [pendingSessionId])
  useEffect(() => {
    if (messagesData && pendingSessionId) {
      const raw = (messagesData as any)?.data
      if (Array.isArray(raw)) {
        const msgs: ChatMessage[] = raw.map((m: any, i: number) => ({
          key: `${pendingSessionId}-${i}`,
          role: m.role === 'user' ? 'user' : 'ai',
          content: m.content || '',
          streaming: false,
          status: 'success' as const,
        }))
        setMessagesMap((prev) => ({ ...prev, [pendingSessionId]: msgs }))
      }
      setPendingSessionId(null)
    }
  }, [messagesData, fetchKey])

  const messages = activeKey ? (messagesMap[activeKey] || []) : []

  const handleSubmit = async (value: string) => {
    if (!value.trim() || isLoading) return
    setInputValue('')

    const isNew = !activeKey
    const sessionId = isNew ? '' : activeKey!
    const userKey = `user-${Date.now()}`
    const aiKey = `ai-${Date.now()}`
    const displayKey = isNew ? `pending-${Date.now()}` : sessionId

    setMessagesMap((prev) => ({
      ...prev,
      [displayKey]: [
        ...(prev[displayKey] || []),
        { key: userKey, role: 'user', content: value.trim() },
        { key: aiKey, role: 'ai', content: '', streaming: true, status: 'loading' },
      ],
    }))
    if (isNew) setActiveKey(displayKey)
    setIsLoading(true)

    const controller = new AbortController()
    abortRef.current = controller

    try {
      // SSE 流式聊天 restful-react 不支持，使用原生 fetch
      const res = await fetch('/api/v1/ai/chat', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ session_id: sessionId, message: value.trim() }),
        signal: controller.signal,
      })

      if (!res.body) {
        console.error('[SSE] response body is null')
        setMessagesMap((prev) => ({
          ...prev, [displayKey]: (prev[displayKey] || []).map((m) =>
            m.key === aiKey ? { ...m, status: 'error' as const, content: '请求失败：响应体为空' } : m
          ),
        }))
        return
      }

      const reader = res.body.getReader()
      const decoder = new TextDecoder()
      let buffer = '', fullContent = '', resolvedSessionId: string | undefined

      while (true) {
        const { done, value: chunk } = await reader.read()
        if (done) break
        buffer += decoder.decode(chunk, { stream: true })
        const lines = buffer.split('\n')
        buffer = lines.pop() || ''
        for (const line of lines) {
          const event = parseSSELine(line)
          if (!event) continue
          if (event.type === 'token' && event.content) {
            fullContent += event.content
            console.log('[SSE] token', { chunkLen: event.content.length, totalLen: fullContent.length, streaming: true })
            setMessagesMap((prev) => ({
              ...prev, [displayKey]: (prev[displayKey] || []).map((m) =>
                m.key === aiKey ? { ...m, content: fullContent } : m),
            }))
          } else if (event.type === 'done') {
              console.log('[SSE] done', { session_id: event.session_id, isNew })
              resolvedSessionId = event.session_id
              if (isNew && resolvedSessionId) {
                setMessagesMap((prev) => {
                const next = { ...prev }
                const msgs = next[displayKey] || []
                delete next[displayKey]
                next[resolvedSessionId!] = msgs.map((m) =>
                  m.key === aiKey ? { ...m, streaming: false, status: 'success' as const } : m)
                return next
              })
              setActiveKey(resolvedSessionId)
            } else {
              setMessagesMap((prev) => ({
                ...prev, [sessionId]: (prev[sessionId] || []).map((m) =>
                  m.key === aiKey ? { ...m, streaming: false, status: 'success' as const } : m),
              }))
            }
            refetchSessions()
        } else if (event.type === 'error') {
            setMessagesMap((prev) => {
              const msgs = prev[displayKey] || []
              return {
                ...prev,
                [displayKey]: msgs.map((m) =>
                  m.key === aiKey
                    ? { ...m, content: `⚠️ 模型调用异常：${event.content || '出错了'}`, streaming: false, status: 'error' as const }
                    : m),
              }
            })
          }
        }
      }
    } catch (e: any) {
      if (e?.name !== 'AbortError') {
        antMsg.error('请求失败')
        setMessagesMap((prev) => {
          const msgs = prev[displayKey] || []
          return {
            ...prev,
            [displayKey]: msgs.map((m) =>
              m.key === aiKey
                ? { ...m, content: '⚠️ 网络请求失败', streaming: false, status: 'error' as const }
                : m),
          }
        })
      }
    } finally {
      setIsLoading(false)
      abortRef.current = null
    }
  }

  const handleCancel = () => {
    abortRef.current?.abort()
    abortRef.current = null
    setMessagesMap((prev) => {
      if (!activeKey) return prev
      return {
        ...prev,
        [activeKey]: (prev[activeKey] || []).map((m) =>
          m.role === 'ai' && m.streaming ? { ...m, streaming: false, status: 'abort' as const } : m),
      }
    })
  }

  const handleDeleteSession = (sessionId: string) => {
    doDeleteSession({ pathParams: { session_id: sessionId } }).then((res: any) => {
      if (res?.code === 0) {
        setMessagesMap((prev) => {
          const next = { ...prev }
          delete next[sessionId]
          return next
        })
        if (activeKey === sessionId) setActiveKey(undefined)
        refetchSessions()
      }
    }).catch(() => {})
  }

  const handleNewChat = () => setActiveKey(undefined)

  const handleSelectSession = (sessionId: string) => {
    setActiveKey(sessionId)
    setSessionPopoverOpen(false)
    if (!messagesMap[sessionId]) {
      setPendingSessionId(sessionId)
    }
  }

  // 通用高阶组件：标签未闭合时显示加载文本
  const withCompleteOnly = (Component: any, loadingText = '加载中...') =>
    ({ streamStatus, ...props }: any) => {
      if (streamStatus === 'loading') {
        return (
          <div style={{ margin: '8px 0', color: '#8c8c8c' }}>
            <Space>
              <Spin indicator={<LoadingOutlined spin />} size="small" />
              <span>{loadingText}</span>
            </Space>
          </div>
        )
      }
      return <Component {...props} />
    }

  const bubbleItems = messages.map((m) => ({
    key: m.key,
    role: m.role,
    content: m.content || (m.role === 'ai' && m.streaming ? '...' : ''),
    loading: m.role === 'ai' && m.streaming && !m.content,
    status: m.status,
    contentRender: m.role === 'ai'
      ? (content: any) => {
          const text = String(content ?? '')
          return (
            <XMarkdown className={xmClassName} content={text} paragraphTag="div"
              streaming={{ hasNextChunk: m.streaming ?? false, enableAnimation: true, tail: { content: '▋' } }}
              components={{
                tool_call: withCompleteOnly(({ children, ...props }: any) => {
                  console.log('[XMarkdown] tool_call', { props, childLen: String(children).length })
                  return (
                    <Collapse
                      size="small"
                      defaultActiveKey={[]}
                      style={{ margin: '8px 0', background: '#fafafa', borderRadius: 6 }}
                      items={[{
                        key: 'tool_call',
                        label: <span style={{ fontSize: 13 }}>🛠️ 调用工具 <Tag style={{ fontSize: 11 }}>{props['data-tool-name'] || '未知'}</Tag></span>,
                        children: <pre style={{ fontSize: 12, margin: 0, whiteSpace: 'pre-wrap' }}>{children}</pre>,
                      }]}
                    />
                  )
                }, '正在调用工具...'),

                tool_result: withCompleteOnly(({ children }: any) => {
                  console.log('[XMarkdown] tool_result', { childLen: String(children).length })
                  return (
                    <Collapse
                      size="small"
                      defaultActiveKey={[]}
                      style={{ margin: '8px 0', background: '#fafafa', borderRadius: 6 }}
                      items={[{
                        key: 'tool_result',
                        label: <span style={{ fontSize: 13 }}>📎 调用工具结果</span>,
                        children: <pre style={{ fontSize: 12, margin: 0, whiteSpace: 'pre-wrap', color: '#666' }}>{children}</pre>,
                      }]}
                    />
                  )
                }, '正在等待工具结果...'),

                tool_call_error: withCompleteOnly(({ children }: any) => {
                  console.log('[XMarkdown] tool_call_error', { childLen: String(children).length })
                  return (
                    <div style={{ margin: '8px 0', padding: '8px 12px', background: '#fff2f0', border: '1px solid #ffccc7', borderRadius: 6, fontSize: 13, color: '#cf1322' }}>
                      ⚠️ 调用异常：{children}
                    </div>
                  )
                }, ''),
              }}
            />
          )
        }
      : undefined,
  }))

  const roles = {
    user: { placement: 'end' as const },
    ai: { placement: 'start' as const },
  }

  const conversationItems = sessionsFromHook.slice(0, 50).map((s: any) => ({
    key: s.session_id,
    label: s.title || s.session_id.slice(0, 8),
  }))

  return (
    <Drawer
      title={
        <div style={{ display: 'flex', alignItems: 'center', justifyContent: 'space-between' }}>
          <span>AI 智能助手</span>
          <div style={{ display: 'flex', alignItems: 'center', gap: 4 }}>
            <Button type="primary" size="small" icon={<PlusOutlined />} onClick={handleNewChat}>新建</Button>
            <Popover
              open={sessionPopoverOpen}
              onOpenChange={(open) => { setSessionPopoverOpen(open); if (open) refetchSessions() }}
              trigger="click"
              placement="bottomRight"
              content={
                <div style={{ width: 260, maxHeight: 300, overflow: 'auto' }}>
                  <Conversations
                    items={conversationItems}
                    activeKey={activeKey}
                    onActiveChange={handleSelectSession}
                    menu={(item) => ({
                      items: [{ key: 'delete', label: '删除', danger: true, onClick: () => handleDeleteSession(item.key) }],
                    })}
                  />
                </div>
              }
            >
              <Button type="text" size="small" icon={<HistoryOutlined />}>历史对话</Button>
            </Popover>
          </div>
        </div>
      }
      placement="right"
      width="33.33vw"
      open={open}
      onClose={onClose}
      styles={{ body: { padding: 0, display: 'flex', flexDirection: 'column' } }}
    >
      <div style={{ flex: 1, display: 'flex', flexDirection: 'column', minHeight: 0 }}>
        <Bubble.List
          items={bubbleItems}
          role={roles}
          autoScroll
          style={{ flex: 1, padding: 16, overflow: 'auto' }}
        />
        <div style={{ padding: '8px 16px 16px', borderTop: '1px solid #f0f0f0' }}>
          <Sender
            value={inputValue}
            onChange={(v) => setInputValue(v ?? '')}
            onSubmit={handleSubmit}
            loading={isLoading}
            onCancel={handleCancel}
            placeholder="输入消息，Enter 发送..."
          />
        </div>
      </div>
    </Drawer>
  )
}
