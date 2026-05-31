import { useState, useEffect, useCallback } from 'react'
import { Routes, Route, useNavigate } from 'react-router-dom'
import { Select, Button, Modal, Input, Drawer, Form, Tag, message, Table } from 'antd'
import { PlusOutlined, SettingOutlined } from '@ant-design/icons'
import {
  useListProjects,
  useListIssues,
  useListExecutors,
  useCreateProject,
  useCreateIssue,
  useUpdateIssue,
  useDeleteIssue,
  useGetProject,
  useUpdateProject,
  useListDocuments,
  useCreateDocument,
  useGetDocument,
  useUpdateDocument,
  YoloTeamYoloCliInternalModelProject,
  YoloTeamYoloCliInternalModelIssue,
  YoloTeamYoloCliInternalModelExecutor,
  ListIssuesQueryParams,
} from './generated'
import dayjs from 'dayjs'

const STATUSES = ['created', 'design', 'review', 'implementation', 'qa', 'pending_review', 'done', 'archived']
const STATUS_NAMES: Record<string, string> = {
  created: '已创建', design: '方案设计', review: '方案评审',
  implementation: '代码实现', qa: 'QA质检', pending_review: '待审查',
  done: '已完成', archived: '已归档',
}
const PRIORITY_COLORS: Record<string, string> = { critical: 'red', high: 'orange', medium: 'blue', low: 'green' }

type Project = YoloTeamYoloCliInternalModelProject
type Issue = YoloTeamYoloCliInternalModelIssue
type Executor = YoloTeamYoloCliInternalModelExecutor

function KanbanPage() {
  const navigate = useNavigate()
  const [selectedProject, setSelectedProject] = useState<string | undefined>()
  const [drawerOpen, setDrawerOpen] = useState(false)
  const [drawerIssue, setDrawerIssue] = useState<Issue | null>(null)
  const [projectModalOpen, setProjectModalOpen] = useState(false)

  const { data: projectsData, refetch: refetchProjects } = useListProjects({})
  const projects = (projectsData as any)?.data || []

  const issueParams: ListIssuesQueryParams = {}
  if (selectedProject) issueParams.project_id = Number(selectedProject)

  const { data: issuesData, refetch: refetchIssues } = useListIssues({ queryParams: issueParams })
  const issues = (issuesData as any)?.data || []

  const { data: executorsData } = useListExecutors({})
  const executors = (executorsData as any)?.data || []

  const load = useCallback(() => {
    refetchProjects()
    refetchIssues()
  }, [refetchProjects, refetchIssues])

  const handleCardClick = async (issue: Issue) => {
    try {
      const resp = await fetch(`/api/v1/issues/${issue.key}`)
      const json = await resp.json()
      if (json.code === 0) {
        setDrawerIssue(json.data)
        setDrawerOpen(true)
      }
    } catch (e: any) { message.error(e.message) }
  }

  const columns = STATUSES.map(status => issues.filter((i: Issue) => i.status === status))

  return (
    <div>
      <div className="topbar">
        <Select
          style={{ width: 220 }}
          placeholder="全部项目"
          allowClear
          onClear={() => setSelectedProject(undefined)}
          onChange={(v) => setSelectedProject(v)}
          options={projects.map((p: Project) => ({ label: p.name, value: p.id as any }))}
        />
        <Button type="primary" icon={<PlusOutlined />} onClick={() => setProjectModalOpen(true)}>新建项目</Button>
        {selectedProject && (
          <Button icon={<SettingOutlined />} onClick={() => {
            const proj = projects.find((p: Project) => p.id === Number(selectedProject))
            if (proj) navigate(`/project/${proj.key}/manage`)
          }}>管理</Button>
        )}
        {selectedProject && (
          <Button onClick={() => navigate(`/project/${projects.find((p: Project) => p.id === Number(selectedProject))?.key}/docs`)}>文档</Button>
        )}
      </div>
      <div className="kanban-container">
        {STATUSES.map((status, idx) => (
          <div key={status} className="kanban-col">
            <div className={`kanban-col-header col-${status}`}>{STATUS_NAMES[status]} ({columns[idx].length})</div>
            {status === 'created' && selectedProject && (
              <Button block size="small" style={{ marginBottom: 8 }} onClick={() => {
                setDrawerIssue({} as Issue)
                setDrawerOpen(true)
              }}>+ 新建</Button>
            )}
            {columns[idx].map((issue: Issue) => (
              <div key={issue.key} className="issue-card" onClick={() => handleCardClick(issue)}>
                <div className="title">{issue.title}</div>
                <div className="meta">
                  <Tag color={PRIORITY_COLORS[issue.priority || '']}>{issue.priority}</Tag>
                  {issue.deadline && <span>{dayjs(issue.deadline).format('MM-DD')}</span>}
                  {issue.repo_name && <span>🔗</span>}
                </div>
              </div>
            ))}
          </div>
        ))}
      </div>

      <Drawer
        title="Issue 详情"
        placement="left"
        width={520}
        open={drawerOpen}
        onClose={() => { setDrawerOpen(false); setDrawerIssue(null) }}
      >
        {drawerIssue && drawerIssue.key ? (
          <IssueDetail issue={drawerIssue} executors={executors} onSave={load} onClose={() => setDrawerOpen(false)} />
        ) : selectedProject ? (
          <CreateIssueForm
            projectId={Number(selectedProject!)}
            executors={executors}
            onSave={() => { setDrawerOpen(false); load() }}
          />
        ) : (
          <div>请先选择一个项目</div>
        )}
      </Drawer>

      <CreateProjectModal open={projectModalOpen} onClose={() => setProjectModalOpen(false)} onDone={load} />
    </div>
  )
}

function CreateProjectModal({ open, onClose, onDone }: { open: boolean; onClose: () => void; onDone: () => void }) {
  const [name, setName] = useState('')
  const [desc, setDesc] = useState('')
  const [, mutate] = useCreateProject({})

  const handleSubmit = async () => {
    try {
      await mutate({ name, description: desc } as any)
      message.success('项目创建成功')
      setName('')
      setDesc('')
      onClose()
      onDone()
    } catch (e: any) { message.error(e.message) }
  }

  return (
    <Modal title="新建项目" open={open} onCancel={onClose} footer={null}>
      <Form layout="vertical">
        <Form.Item label="名称" required><Input value={name} onChange={e => setName(e.target.value)} /></Form.Item>
        <Form.Item label="描述"><Input.TextArea value={desc} onChange={e => setDesc(e.target.value)} /></Form.Item>
        <Button type="primary" onClick={handleSubmit}>创建</Button>
      </Form>
    </Modal>
  )
}

function IssueDetail({ issue, executors, onSave, onClose }: {
  issue: Issue; executors: Executor[]; onSave: () => void; onClose: () => void
}) {
  const [title, setTitle] = useState(issue.title)
  const [desc, setDesc] = useState(issue.description)
  const [status, setStatus] = useState(issue.status)
  const [priority, setPriority] = useState(issue.priority)
  const [execId, setExecId] = useState(issue.executor_id)
  const [repoUrl, setRepoUrl] = useState(issue.repo_url)
  const [repoName, setRepoName] = useState(issue.repo_name)
  const [branch, setBranch] = useState(issue.branch_name)
  const [, updateMutate] = useUpdateIssue({ key: issue.key || '' })
  const [, deleteMutate] = useDeleteIssue({})

  const handleSave = async () => {
    try {
      await updateMutate({
        title, description: desc, status, priority,
        executor_id: execId, repo_url: repoUrl, repo_name: repoName, branch_name: branch,
      } as any)
      message.success('已保存')
      onSave()
    } catch (e: any) { message.error(e.message) }
  }

  const handleDelete = async () => {
    try {
      await deleteMutate(issue.key)
      message.success('已删除')
      onClose()
      onSave()
    } catch (e: any) { message.error(e.message) }
  }

  return (
    <div>
      <p><Tag>{issue.key}</Tag> <Tag color="blue">{STATUS_NAMES[issue.status || '']}</Tag></p>
      <Form layout="vertical">
        <Form.Item label="标题"><Input value={title} onChange={e => setTitle(e.target.value)} /></Form.Item>
        <Form.Item label="描述"><Input.TextArea rows={4} value={desc} onChange={e => setDesc(e.target.value)} /></Form.Item>
        <Form.Item label="状态">
          <Select value={status} onChange={setStatus} options={STATUSES.map(s => ({ label: STATUS_NAMES[s], value: s }))} />
        </Form.Item>
        <Form.Item label="优先级">
          <Select value={priority} onChange={setPriority}
            options={['critical', 'high', 'medium', 'low'].map(p => ({ label: p, value: p }))} />
        </Form.Item>
        <Form.Item label="执行人">
          <Select value={execId} onChange={setExecId} allowClear
            options={executors.map((e: Executor) => ({ label: e.name, value: e.id }))} />
        </Form.Item>
        <Form.Item label="仓库 URL"><Input value={repoUrl} onChange={e => setRepoUrl(e.target.value)} /></Form.Item>
        <Form.Item label="仓库名称"><Input value={repoName} onChange={e => setRepoName(e.target.value)} /></Form.Item>
        <Form.Item label="分支"><Input value={branch} onChange={e => setBranch(e.target.value)} /></Form.Item>
      </Form>
      <Button type="primary" onClick={handleSave}>保存</Button>
      <Button danger style={{ marginLeft: 8 }} onClick={handleDelete}>删除</Button>
    </div>
  )
}

function CreateIssueForm({ projectId, executors, onSave }: {
  projectId: number; executors: Executor[]; onSave: () => void
}) {
  const [title, setTitle] = useState('')
  const [desc, setDesc] = useState('')
  const [execId, setExecId] = useState<number | null>(null)
  const [, mutate] = useCreateIssue({})

  const handleSubmit = async () => {
    try {
      await mutate({ project_id: projectId, title, description: desc, executor_id: execId } as any)
      message.success('创建成功')
      onSave()
    } catch (e: any) { message.error(e.message) }
  }

  return (
    <Form layout="vertical">
      <Form.Item label="标题" required><Input value={title} onChange={e => setTitle(e.target.value)} /></Form.Item>
      <Form.Item label="描述"><Input.TextArea rows={3} value={desc} onChange={e => setDesc(e.target.value)} /></Form.Item>
      <Form.Item label="执行人">
        <Select value={execId} allowClear onChange={setExecId as any}
          options={executors.map((e: Executor) => ({ label: e.name, value: e.id }))} />
      </Form.Item>
      <Button type="primary" onClick={handleSubmit}>创建</Button>
    </Form>
  )
}

function ProjectManage() {
  const key = window.location.pathname.split('/')[2]
  const [name, setName] = useState('')
  const [desc, setDesc] = useState('')

  const { data } = useGetProject({ key })
  const project = (data as any)?.data
  const [, updateMutate] = useUpdateProject({ key })

  useEffect(() => {
    if (project) { setName(project.name || ''); setDesc(project.description || '') }
  }, [project])

  const handleSave = async () => {
    try {
      await updateMutate({ name, description: desc } as any)
      message.success('已保存')
    } catch (e: any) { message.error(e.message) }
  }

  return (
    <div style={{ padding: 24, maxWidth: 500 }}>
      <h2>项目管理</h2>
      {project && <p>Key: {project.key}</p>}
      <Form layout="vertical">
        <Form.Item label="名称"><Input value={name} onChange={e => setName(e.target.value)} /></Form.Item>
        <Form.Item label="描述"><Input.TextArea value={desc} onChange={e => setDesc(e.target.value)} /></Form.Item>
        <Button type="primary" onClick={handleSave}>保存</Button>
        <Button style={{ marginLeft: 8 }} onClick={() => window.history.back()}>返回</Button>
      </Form>
    </div>
  )
}

function DocumentList() {
  const key = window.location.pathname.split('/')[2]
  const [title, setTitle] = useState('')

  const { data, refetch } = useListDocuments({ key })
  const docs = (data as any)?.data || []
  const [, createMutate] = useCreateDocument({ key })

  const handleCreate = async () => {
    try {
      await createMutate({ title, content: '' } as any)
      message.success('创建成功')
      setTitle('')
      refetch()
    } catch (e: any) { message.error(e.message) }
  }

  const cols = [
    { title: 'Key', dataIndex: 'key' },
    { title: '标题', dataIndex: 'title' },
    { title: '更新时间', dataIndex: 'updated_at', render: (v: string) => v?.slice(0, 19) },
    {
      title: '操作', render: (_: any, r: any) => (
        <a href={`/doc/${r.key}`}>查看</a>
      )
    },
  ]

  return (
    <div style={{ padding: 24 }}>
      <h2>文档列表</h2>
      <div style={{ display: 'flex', gap: 8, marginBottom: 16 }}>
        <Input style={{ width: 300 }} placeholder="文档标题" value={title} onChange={e => setTitle(e.target.value)} />
        <Button type="primary" onClick={handleCreate}>新建文档</Button>
        <Button onClick={() => window.history.back()}>返回</Button>
      </div>
      <Table columns={cols} dataSource={docs} rowKey="key" size="small" />
    </div>
  )
}

function DocumentDetail() {
  const key = window.location.pathname.split('/')[2]
  const [content, setContent] = useState('')
  const [editing, setEditing] = useState(false)

  const { data } = useGetDocument({ key })
  const doc = (data as any)?.data
  const [, updateMutate] = useUpdateDocument({ key })

  useEffect(() => {
    if (doc) setContent(doc.content || '')
  }, [doc])

  const handleSave = async () => {
    try {
      await updateMutate({ content } as any)
      message.success('已保存')
      setEditing(false)
    } catch (e: any) { message.error(e.message) }
  }

  if (!doc) return <div>Loading...</div>

  return (
    <div style={{ padding: 24, maxWidth: 900 }}>
      <h2>{doc.title}</h2>
      <p style={{ color: '#999' }}>路径: {doc.file_path}</p>
      {editing ? (
        <div>
          <Input.TextArea rows={20} value={content} onChange={e => setContent(e.target.value)} />
          <Button type="primary" onClick={handleSave} style={{ marginTop: 8 }}>保存</Button>
          <Button onClick={() => setEditing(false)} style={{ marginTop: 8, marginLeft: 8 }}>取消</Button>
        </div>
      ) : (
        <div>
          <pre style={{ whiteSpace: 'pre-wrap' }}>{content}</pre>
          <Button onClick={() => setEditing(true)}>编辑</Button>
          <Button style={{ marginLeft: 8 }} onClick={() => window.history.back()}>返回</Button>
        </div>
      )}
    </div>
  )
}

export default function App() {
  return (
    <Routes>
      <Route path="/" element={<KanbanPage />} />
      <Route path="/project/:key/manage" element={<ProjectManage />} />
      <Route path="/project/:key/docs" element={<DocumentList />} />
      <Route path="/doc/:key" element={<DocumentDetail />} />
    </Routes>
  )
}
