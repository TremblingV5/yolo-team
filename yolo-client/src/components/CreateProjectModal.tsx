import { useState } from 'react'
import { Modal, Form, Input, Button, message } from 'antd'
import { useCreateProject } from '../generated'

interface CreateProjectModalProps {
  open: boolean
  onClose: () => void
  onDone: () => void
}

export default function CreateProjectModal({ open, onClose, onDone }: CreateProjectModalProps) {
  const [name, setName] = useState('')
  const [desc, setDesc] = useState('')
  const { mutate } = useCreateProject({})

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
        <Form.Item label="名称" required>
          <Input value={name} onChange={e => setName(e.target.value)} />
        </Form.Item>
        <Form.Item label="描述">
          <Input.TextArea value={desc} onChange={e => setDesc(e.target.value)} />
        </Form.Item>
        <Button type="primary" onClick={handleSubmit}>创建</Button>
      </Form>
    </Modal>
  )
}
