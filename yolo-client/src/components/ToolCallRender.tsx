import { LoadingOutlined } from '@ant-design/icons'
import { Collapse, Space, Spin, Tag } from 'antd'

/**
 * 通用高阶组件：XMarkdown 的 streamStatus = 'loading' 时显示加载文本。
 */
export function withCompleteOnly(Component: any, loadingText = '加载中...') {
  return ({ streamStatus, ...props }: any) => {
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
}

// ---------- 工具调用渲染 ----------

function ToolCallContent({ children, ...props }: any) {
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
}

function ToolResultContent({ children }: any) {
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
}

function ToolCallErrorContent({ children }: any) {
  return (
    <div style={{ margin: '8px 0', padding: '8px 12px', background: '#fff2f0', border: '1px solid #ffccc7', borderRadius: 6, fontSize: 13, color: '#cf1322' }}>
      ⚠️ 调用异常：{children}
    </div>
  )
}

/** XMarkdown 使用的 tool_call / tool_result / tool_call_error 自定义组件。 */
export const toolCallComponents = {
  tool_call: withCompleteOnly(ToolCallContent, '正在调用工具...'),
  tool_result: withCompleteOnly(ToolResultContent, '正在等待工具结果...'),
  tool_call_error: withCompleteOnly(ToolCallErrorContent, ''),
}
