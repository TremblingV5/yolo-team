import { useState } from 'react'
import { Table, Button, Modal, Form, Input, Select, Popconfirm, Tag, message } from 'antd'
import { PlusOutlined } from '@ant-design/icons'
import { useNavigate } from 'react-router-dom'
import {
  useListExecutors,
  useCreateExecutor,
  useUpdateExecutor,
  useDeleteExecutor,
  YoloTeamYoloCliInternalModelExecutor as Executor,
} from '../generated'

const ROLE_NAMES: Record<string, string> = { leader: 'Leader', architect: 'Architect', developer: 'Developer', qa: 'QA' }
const ROLE_COLORS: Record<string, string> = { leader: 'red', architect: 'purple', developer: 'blue', qa: 'orange' }

export default function ExecutorManage() {
  const navigate = useNavigate()
  const [createOpen, setCreateOpen] = useState(false)
  const [editOpen, setEditOpen] = useState(false)
  const [editExec, setEditExec] = useState<Executor | null>(null)
  const [newName, setNewName] = useState('')
  const [newSoul, setNewSoul] = useState('')
  const [editRole, setEditRole] = useState('')
  const [editSoul, setEditSoul] = useState('')

  const { data, refetch } = useListExecutors({})
  const executors = (data as any)?.data || []

  const { mutate: createMutate } = useCreateExecutor({})
  const { mutate: deleteMutate } = useDeleteExecutor({})

  const handleCreate = async () => {
    try {
      await createMutate({ name: newName, soul: newSoul } as any)
      message.success('创建成功')
      setCreateOpen(false)
      setNewName('')
      setNewSoul('')
      refetch()
    } catch (e: any) { message.error(e.message) }
  }

  const handleDelete = async (name: string) => {
    try {
      await deleteMutate(name)
      message.success('已删除')
      refetch()
    } catch (e: any) { message.error(e.message) }
  }

  const openEdit = (exec: Executor) => {
    setEditExec(exec)
    setEditRole(exec.role || '')
    setEditSoul(exec.soul || '')
    setEditOpen(true)
  }

  const handleEditSave = async () => {
    if (!editExec) return
    try {
      const resp = await fetch(`/api/v1/executors/${editExec.name}`, {
        method: 'PUT',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ role: editRole, soul: editSoul }),
      })
      const json = await resp.json()
      if (json.code === 0) {
        message.success('已保存')
        setEditOpen(false)
        refetch()
      } else {
        message.error(json.message || '保存失败')
      }
    } catch (e: any) { message.error(e.message) }
  }

  const cols = [
    { title: 'ID', dataIndex: 'id', width: 60 },
    { title: 'Name', dataIndex: 'name' },
    {
      title: 'Role', dataIndex: 'role',
      render: (v: string) => <Tag color={ROLE_COLORS[v] || 'default'}>{ROLE_NAMES[v] || v}</Tag>,
    },
    { title: 'Soul', dataIndex: 'soul', ellipsis: true },
    { title: 'Created', dataIndex: 'created_at', render: (v: string) => v?.slice(0, 19) },
    {
      title: '操作', render: (_: any, r: Executor) => (
        <>
          <a onClick={() => openEdit(r)}>编辑</a>
          <Popconfirm title="确定删除？" onConfirm={() => handleDelete(r.name || '')}>
            <a style={{ marginLeft: 12, color: 'red' }}>删除</a>
          </Popconfirm>
        </>
      )
    },
  ]

  return (
    <div style={{ padding: 24 }}>
      <h2>执行人管理</h2>
      <div style={{ marginBottom: 16 }}>
        <Button type="primary" icon={<PlusOutlined />} onClick={() => setCreateOpen(true)}>新建执行人</Button>
        <Button style={{ marginLeft: 8 }} onClick={() => navigate(-1)}>返回</Button>
      </div>

      <Table columns={cols} dataSource={executors} rowKey="id" size="small" />

      <Modal title="新建执行人" open={createOpen} onCancel={() => setCreateOpen(false)} footer={null}>
        <Form layout="vertical">
          <Form.Item label="Name" required>
            <Input value={newName} onChange={e => setNewName(e.target.value)} placeholder="不含空格的唯一名称" />
          </Form.Item>
          <Form.Item label="Soul">
            <Input.TextArea rows={4} value={newSoul} onChange={e => setNewSoul(e.target.value)}
              placeholder="描述执行人的能力与行为准则" />
          </Form.Item>
          <Button type="primary" onClick={handleCreate}>创建</Button>
        </Form>
      </Modal>

      <Modal title="编辑执行人" open={editOpen} onCancel={() => setEditOpen(false)} footer={null}>
        <Form layout="vertical">
          <Form.Item label="Name"><Input value={editExec?.name} disabled /></Form.Item>
          <Form.Item label="Role">
            <Select value={editRole} onChange={setEditRole}
              options={Object.entries(ROLE_NAMES).map(([k, v]) => ({ label: v, value: k }))} />
          </Form.Item>
          <Form.Item label="Soul">
            <Input.TextArea rows={4} value={editSoul} onChange={e => setEditSoul(e.target.value)} />
          </Form.Item>
          <Button type="primary" onClick={handleEditSave}>保存</Button>
        </Form>
      </Modal>
    </div>
  )
}
