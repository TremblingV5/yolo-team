import { PlusOutlined } from '@ant-design/icons'
import { Button, Form, Input, message, Modal, Space, Tag } from 'antd'
import { useState } from 'react'
import { useNavigate } from 'react-router-dom'
import {
  YoloTeamYoloCliInternalModelDocument as Document,
  YoloTeamYoloCliInternalModelExecutor as Executor,
  YoloTeamYoloCliInternalModelIssue as Issue,
  useDeleteIssue, useUpdateIssue,
  useLinkDocument, useUnlinkDocument,
  useListTasks,
} from '../generated'
import IssueFormFields, { IssueFormValues } from './IssueFormFields'
import IssueTaskSection from './IssueTaskSection'

interface IssueDrawerProps {
  issue: Issue
  executors: Executor[]
  projectDocs: Document[]
  projectKey: string
  open: boolean
  onClose: () => void
  onSaved: () => void
}

export default function IssueDrawer({ issue, executors, projectDocs, projectKey, onClose, onSaved }: IssueDrawerProps) {
  const navigate = useNavigate()
  const [form, setForm] = useState<IssueFormValues>({
    title: issue.title || '',
    description: issue.description || '',
    status: issue.status || '',
    priority: issue.priority || '',
    executorId: issue.executor_id,
    repoUrl: issue.repo_url || '',
    branch: issue.branch_name || '',
  })
  const [docModal, setDocModal] = useState(false)
  const [newDocTitle, setNewDocTitle] = useState('')

  const { mutate: updateMutate } = useUpdateIssue({ key: issue.key || '' })
  const { mutate: deleteMutate } = useDeleteIssue({})
  const { mutate: linkDocMutate } = useLinkDocument({ key: issue.key || '' })
  const { mutate: unlinkDocMutate } = useUnlinkDocument({ key: issue.key || '' })
  const { data: tasksData, refetch: refetchTasks } = useListTasks({ key: issue.key || '' })
  const tasks = (tasksData as any)?.data || []

  const linkedDocs = (issue as any).documents || []

  const showError = (e: any) => message.error(e?.data?.message || e.message || '操作失败')

  const handleSave = async () => {
    try {
      await updateMutate({
        title: form.title, description: form.description, status: form.status,
        priority: form.priority, executor_id: form.executorId,
        repo_url: form.repoUrl, branch_name: form.branch,
      } as any)
      message.success('已保存')
      onSaved()
    } catch (e: any) { showError(e) }
  }

  const handleDelete = async () => {
    try {
      await deleteMutate(issue.key || '')
      message.success('已删除')
      onClose()
      onSaved()
    } catch (e: any) { showError(e) }
  }

  const handleLinkDoc = async (docId: number) => {
    try {
      await linkDocMutate({ document_id: docId } as any)
      message.success('已关联')
      onSaved()
    } catch (e: any) { showError(e) }
  }

  const handleUnlinkDoc = async (docId: number) => {
    try {
      await unlinkDocMutate(undefined, { document_id: docId })
      message.success('已取消关联')
      onSaved()
    } catch (e: any) { showError(e) }
  }

  const handleCreateAndEdit = () => {
    const docTitle = newDocTitle || '未命名文档'
    setDocModal(false)
    setNewDocTitle('')
    navigate(`/doc/new?linkTo=${issue.key}&projectKey=${projectKey}&title=${encodeURIComponent(docTitle)}`)
  }

  const availableDocs = projectDocs.filter((d: any) =>
    !linkedDocs.find((ld: any) => ld.id === d.id) &&
    (!newDocTitle || d.title?.toLowerCase().includes(newDocTitle.toLowerCase()))
  )

  return (
    <div style={{ display: 'flex', flexDirection: 'column', height: '100%' }}>
      <div style={{ flex: 1, overflow: 'auto', paddingBottom: 70 }}>
        <Form layout="vertical" size="middle">
          {/* Key badge + form fields */}
          <div style={{ display: 'flex', gap: 8, alignItems: 'center', marginBottom: 16 }}>
            <Tag color="default" style={{ fontSize: 13, padding: '2px 10px', margin: 0 }}>
              {issue.key}
            </Tag>
          </div>

          <IssueFormFields
            values={form}
            executors={executors}
            onChange={(patch) => setForm((prev) => ({ ...prev, ...patch }))}
          />

          {/* Task section */}
          <IssueTaskSection
            issueKey={issue.key || ''}
            tasks={tasks}
            onTaskChange={refetchTasks}
          />

          {/* Linked docs */}
          <Form.Item label="关联文档">
            <div style={{ marginBottom: 6 }}>
              <Button size="small" icon={<PlusOutlined />} onClick={() => setDocModal(true)}>添加文档</Button>
            </div>
            {linkedDocs.length === 0 ? (
              <span style={{ color: '#bbb', fontSize: 13 }}>暂无</span>
            ) : (
              <Space wrap>
                {linkedDocs.map((d: any) => (
                  <Tag key={d.id} closable onClose={() => handleUnlinkDoc(d.id)} color="processing">
                    {d.title}
                  </Tag>
                ))}
              </Space>
            )}
          </Form.Item>
        </Form>
      </div>

      {/* Sticky footer */}
      <div style={{
        position: 'absolute', bottom: 0, left: 0, right: 0,
        padding: '12px 24px', background: '#fff', borderTop: '1px solid #f0f0f0',
        display: 'flex', gap: 8,
      }}>
        <Button type="primary" onClick={handleSave} style={{ flex: 1 }}>保存</Button>
        <Button danger onClick={handleDelete}>删除</Button>
      </div>

      {/* Link document modal */}
      <Modal title="关联文档" open={docModal} onCancel={() => { setDocModal(false); setNewDocTitle('') }} footer={null} width={480}>
        <Input.Search
          placeholder="搜索已有文档..." value={newDocTitle} onChange={e => setNewDocTitle(e.target.value)}
          allowClear enterButton="检索" style={{ marginBottom: 12 }}
        />
        <Button type="dashed" icon={<PlusOutlined />} block onClick={handleCreateAndEdit} style={{ marginBottom: 12 }}>
          新建文档并关联
        </Button>
        <div style={{ maxHeight: 280, overflowY: 'auto' }}>
          {availableDocs.map((d: any) => (
            <div key={d.id} style={{
              padding: '8px 12px', cursor: 'pointer', borderRadius: 6,
              marginBottom: 4, background: '#fafafa', border: '1px solid #f0f0f0',
            }}
              onMouseEnter={e => (e.currentTarget.style.background = '#e6f7ff')}
              onMouseLeave={e => (e.currentTarget.style.background = '#fafafa')}
              onClick={() => { handleLinkDoc(d.id); setDocModal(false) }}>
              <span style={{ fontWeight: 500 }}>{d.title}</span>
              {d.creator && <span style={{ marginLeft: 8, color: '#999', fontSize: 12 }}>by {d.creator}</span>}
            </div>
          ))}
          {availableDocs.length === 0 && !newDocTitle && (
            <div style={{ color: '#bbb', textAlign: 'center', padding: 24 }}>所有文档已关联</div>
          )}
        </div>
      </Modal>
    </div>
  )
}
