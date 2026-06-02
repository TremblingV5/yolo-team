import { useState, useEffect } from 'react'
import { Modal, Form, Input, Button, message } from 'antd'
import { useGetProject, useUpdateProject } from '../generated'

interface ProjectEditModalProps {
  projectKey: string
  open: boolean
  onClose: () => void
}

export default function ProjectEditModal({ projectKey, open, onClose }: ProjectEditModalProps) {
  const [name, setName] = useState('')
  const [desc, setDesc] = useState('')

  const { data } = useGetProject({ key: projectKey })
  const project = (data as any)?.data
  const { mutate: updateMutate } = useUpdateProject({ key: projectKey })

  useEffect(() => {
    if (project && open) {
      setName(project.name || '')
      setDesc(project.description || '')
    }
  }, [project, open])

  const handleSave = async () => {
    try {
      await updateMutate({ name, description: desc } as any)
      message.success('已保存')
      onClose()
    } catch (e: any) { message.error(e.message) }
  }

  return (
    <Modal title="编辑项目" open={open} onCancel={onClose} footer={null}>
      <Form layout="vertical">
        <Form.Item label="名称"><Input value={name} onChange={e => setName(e.target.value)} /></Form.Item>
        <Form.Item label="描述"><Input.TextArea value={desc} onChange={e => setDesc(e.target.value)} /></Form.Item>
        <Button type="primary" onClick={handleSave}>保存</Button>
        <Button style={{ marginLeft: 8 }} onClick={onClose}>取消</Button>
      </Form>
    </Modal>
  )
}
