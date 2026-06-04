import { PlusOutlined } from '@ant-design/icons'
import { Button, Form, Input, Modal, Select, Space, Tag, message } from 'antd'
import { useState } from 'react'
import { useNavigate } from 'react-router-dom'
import { STATUSES, STATUS_NAMES, TASK_STATUSES, TASK_STATUS_NAMES } from '../constants'
import {
    YoloTeamYoloCliInternalModelDocument as Document,
    YoloTeamYoloCliInternalModelExecutor as Executor,
    YoloTeamYoloCliInternalModelIssue as Issue,
    YoloTeamYoloCliInternalModelTask as Task,
    useDeleteIssue, useUpdateIssue,
    useLinkDocument, useUnlinkDocument,
    useListTasks, useCreateTask, useUpdateTask, useDeleteTask,
} from '../generated'

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
  const [title, setTitle] = useState(issue.title || '')
  const [desc, setDesc] = useState(issue.description || '')
  const [status, setStatus] = useState(issue.status || '')
  const [priority, setPriority] = useState(issue.priority || '')
  const [execId, setExecId] = useState(issue.executor_id)
  const [repoUrl, setRepoUrl] = useState(issue.repo_url || '')
  const [branch, setBranch] = useState(issue.branch_name || '')
  const [docModal, setDocModal] = useState(false)
  const [newDocTitle, setNewDocTitle] = useState('')

  // Task state
  const [taskModal, setTaskModal] = useState(false)
  const [taskTitle, setTaskTitle] = useState('')
  const [taskDesc, setTaskDesc] = useState('')
  const [taskDetail, setTaskDetail] = useState<Task | null>(null)

  const { mutate: updateMutate } = useUpdateIssue({ key: issue.key || '' })
  const { mutate: deleteMutate } = useDeleteIssue({})

  const { mutate: linkDocMutate } = useLinkDocument({ key: issue.key || '' })
  const { mutate: unlinkDocMutate } = useUnlinkDocument({ key: issue.key || '' })

  // Task hooks
  const { data: tasksData, refetch: refetchTasks } = useListTasks({ key: issue.key || '' })
  const tasks = (tasksData as any)?.data || ([] as Task[])
  const { mutate: createTaskMutate } = useCreateTask({ key: issue.key || '' })
  const { mutate: updateTaskMutate } = useUpdateTask({ key: issue.key || '', task_key: '' })
  const { mutate: deleteTaskMutate } = useDeleteTask({ key: issue.key || '', task_key: '' })

  const linkedDocs = (issue as any).documents || []

  const showError = (e: any) => message.error(e?.data?.message || e.message || '操作失败')

  const handleSave = async () => {
    try {
      await updateMutate({
        title, description: desc, status, priority,
        executor_id: execId, repo_url: repoUrl, branch_name: branch,
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

  // Task handlers
  const handleCreateTask = async () => {
    if (!taskTitle.trim()) { message.warning('请输入任务标题'); return }
    try {
      await createTaskMutate({ title: taskTitle, description: taskDesc })
      message.success('任务创建成功')
      setTaskModal(false)
      setTaskTitle('')
      setTaskDesc('')
      refetchTasks()
    } catch (e: any) { showError(e) }
  }

  const handleUpdateTaskStatus = async (taskKey: string, newStatus: string) => {
    try {
      await updateTaskMutate(
        { status: newStatus },
        { pathParams: { key: issue.key || '', task_key: taskKey } }
      )
      refetchTasks()
    } catch (e: any) { showError(e) }
  }

  const handleDeleteTask = async (taskKey: string) => {
    try {
      await deleteTaskMutate(undefined, { pathParams: { key: issue.key || '', task_key: taskKey } })
      message.success('任务已删除')
      refetchTasks()
    } catch (e: any) { showError(e) }
  }

  return (
    <div style={{ display: 'flex', flexDirection: 'column', height: '100%' }}>
      <div style={{ flex: 1, overflow: 'auto', paddingBottom: 70 }}>
        <Form layout="vertical" size="middle">

          {/* Row 1: key + status + priority */}
          <div style={{ display: 'flex', gap: 8, alignItems: 'center', marginBottom: 16 }}>
            <Tag color="default" style={{ fontSize: 13, padding: '2px 10px', margin: 0 }}>
              {issue.key}
            </Tag>
            <Select size="small" value={status} onChange={setStatus} style={{ width: 110 }}
              options={STATUSES.map(s => ({ label: STATUS_NAMES[s], value: s }))} />
            <Select size="small" value={priority} onChange={setPriority} style={{ width: 100 }}
              options={['critical', 'high', 'medium', 'low'].map(p => ({ label: p, value: p }))} />
          </div>

          {/* Row 2: title */}
          <Form.Item label="标题">
            <Input value={title} onChange={e => setTitle(e.target.value)} />
          </Form.Item>

          {/* Row 3: description */}
          <Form.Item label="描述">
            <Input.TextArea rows={3} value={desc} onChange={e => setDesc(e.target.value)} />
          </Form.Item>

          {/* Row 4: executor */}
          <Form.Item label="执行人">
            <Select value={execId} onChange={setExecId} allowClear placeholder="选择执行人"
              options={executors.map(e => ({ label: e.name, value: e.id }))} />
          </Form.Item>

          {/* Row 5: repo URL */}
          <Form.Item label="仓库 URL">
            <Input value={repoUrl} onChange={e => setRepoUrl(e.target.value)} placeholder="https://github.com/..." />
          </Form.Item>

          {/* Row 6: branch */}
          <Form.Item label="分支">
            <Input value={branch} onChange={e => setBranch(e.target.value)} placeholder="main" />
          </Form.Item>

          {/* Row 7: tasks */}
          <Form.Item label="任务">
            <div style={{ marginBottom: 8 }}>
              <Button size="small" icon={<PlusOutlined />} onClick={() => setTaskModal(true)}>添加任务</Button>
            </div>
            {tasks.length === 0 ? (
              <span style={{ color: '#bbb', fontSize: 13 }}>暂无任务</span>
            ) : (
              <div style={{ display: 'flex', flexDirection: 'column', gap: 6 }}>
                {tasks.map(t => (
                  <div key={t.id} style={{
                    display: 'flex', alignItems: 'center', gap: 8,
                    padding: '6px 10px', borderRadius: 6,
                    background: t.status === 'done' ? '#f6ffed' : '#fafafa',
                    border: '1px solid', borderColor: t.status === 'done' ? '#b7eb8f' : '#f0f0f0',
                  }}>
                    <Select size="small" value={t.status} style={{ width: 90, flexShrink: 0 }}
                      onChange={v => handleUpdateTaskStatus(t.key!, v)}
                      options={TASK_STATUSES.map(s => ({ label: TASK_STATUS_NAMES[s], value: s }))} />
                    <Tag style={{ fontSize: 11, flexShrink: 0, maxWidth: 120, overflow: 'hidden', textOverflow: 'ellipsis' }}>{t.key}</Tag>
                    <span style={{ flex: 1, fontSize: 13, overflow: 'hidden', textOverflow: 'ellipsis', whiteSpace: 'nowrap', minWidth: 0 }}>{t.title}</span>
                    <Button size="small" type="link"
                      onClick={() => setTaskDetail(t)}>详情</Button>
                    <Button size="small" type="link" danger
                      onClick={() => handleDeleteTask(t.key!)}>删除</Button>
                  </div>
                ))}
              </div>
            )}
          </Form.Item>

          {/* Row 8: linked docs */}
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

      {/* Create Task Modal */}
      <Modal title="新建任务" open={taskModal} onCancel={() => { setTaskModal(false); setTaskTitle(''); setTaskDesc('') }}
        onOk={handleCreateTask} okText="创建" width={420}>
        <Form layout="vertical">
          <Form.Item label="任务标题" required>
            <Input value={taskTitle} onChange={e => setTaskTitle(e.target.value)} placeholder="输入任务标题" />
          </Form.Item>
          <Form.Item label="描述">
            <Input.TextArea rows={2} value={taskDesc} onChange={e => setTaskDesc(e.target.value)} />
          </Form.Item>
        </Form>
      </Modal>

      {/* Task Detail Modal */}
      <Modal title="任务详情" open={!!taskDetail} onCancel={() => setTaskDetail(null)} footer={null} width={480}>
        {taskDetail && (
          <div>
            <p><strong>Key:</strong> {taskDetail.key}</p>
            <p><strong>标题:</strong> {taskDetail.title}</p>
            <p><strong>描述:</strong> {taskDetail.description || '无'}</p>
            <p><strong>状态:</strong> {TASK_STATUS_NAMES[taskDetail.status] || taskDetail.status}</p>
            <p><strong>创建时间:</strong> {taskDetail.created_at}</p>
            <p><strong>更新时间:</strong> {taskDetail.updated_at}</p>
          </div>
        )}
      </Modal>

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
