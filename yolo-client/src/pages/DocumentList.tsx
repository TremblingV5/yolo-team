import { PlusOutlined } from '@ant-design/icons'
import { Button, Form, Input, message, Modal, Popconfirm, Table } from 'antd'
import { useState } from 'react'
import { useNavigate, useParams } from 'react-router-dom'
import { useCreateDocument, useDeleteDocument, useListDocuments } from '../generated'

export default function DocumentList() {
  const { key } = useParams<{ key: string }>()
  const navigate = useNavigate()
  const [search, setSearch] = useState('')
  const [createOpen, setCreateOpen] = useState(false)
  const [newTitle, setNewTitle] = useState('')

  const { data, refetch } = useListDocuments({ key: key || '' })
  const docs = ((data as any)?.data || []).filter((d: any) =>
    !search || d.title?.toLowerCase().includes(search.toLowerCase()) || d.key?.toLowerCase().includes(search.toLowerCase())
  )

  const { mutate: createMutate } = useCreateDocument({ key: key || '' })
  const { mutate: deleteMutate } = useDeleteDocument({ key: '' })

  const handleCreate = async () => {
    if (!newTitle) return
    try {
      const result = await createMutate({ title: newTitle, content: '', creator: '人类' } as any)
      const docKey = (result as any)?.data?.key
      if (docKey) {
        message.success('创建成功')
        setNewTitle('')
        setCreateOpen(false)
        navigate(`/doc/${docKey}`)
      } else {
        message.success('创建成功')
        setNewTitle('')
        setCreateOpen(false)
        refetch()
      }
    } catch (e: any) { message.error(e.message) }
  }

  const handleDelete = async (docKey: string) => {
    try {
      await deleteMutate(undefined, { pathParams: { key: docKey } })
      message.success('已删除')
      refetch()
    } catch (e: any) { message.error(e.message) }
  }

  const cols = [
    { title: 'Key', dataIndex: 'key', width: 150 },
    { title: '标题', dataIndex: 'title' },
    { title: '创建人', dataIndex: 'creator', width: 100 },
    { title: '更新时间', dataIndex: 'updated_at', width: 180, render: (v: string) => v?.slice(0, 19) },
    {
      title: '操作', width: 120, render: (_: any, r: any) => (
        <>
          <a onClick={() => navigate(`/doc/${r.key}`)}>查看</a>
          <Popconfirm title="确定删除？" onConfirm={() => handleDelete(r.key)}>
            <a style={{ marginLeft: 12, color: 'red' }}>删除</a>
          </Popconfirm>
        </>
      )
    },
  ]

  return (
    <div style={{ padding: 24 }}>
      <h2>文档列表</h2>
      <div style={{ display: 'flex', gap: 8, marginBottom: 16 }}>
        <Input.Search
          style={{ width: 300 }}
          placeholder="检索文档"
          allowClear
          value={search}
          onChange={e => setSearch(e.target.value)}
        />
        <Button type="primary" icon={<PlusOutlined />} onClick={() => setCreateOpen(true)}>新建文档</Button>
        <Button onClick={() => navigate(-1)}>返回</Button>
      </div>
      <Table columns={cols} dataSource={docs} rowKey="key" size="small" />

      <Modal title="新建文档" open={createOpen} onCancel={() => setCreateOpen(false)} footer={null}>
        <Form layout="vertical">
          <Form.Item label="标题" required>
            <Input value={newTitle} onChange={e => setNewTitle(e.target.value)}
              onPressEnter={handleCreate} placeholder="输入文档标题" />
          </Form.Item>
          <Button type="primary" onClick={handleCreate} disabled={!newTitle}>创建并编辑</Button>
          <Button style={{ marginLeft: 8 }} onClick={() => setCreateOpen(false)}>取消</Button>
        </Form>
      </Modal>
    </div>
  )
}
