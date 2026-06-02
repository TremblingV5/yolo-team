import { Button, Form, Input, message, Select } from 'antd'
import { useState } from 'react'
import { YoloTeamYoloCliInternalModelExecutor as Executor, useCreateIssue } from '../generated'

interface CreateIssueFormProps {
  projectId: number
  executors: Executor[]
  onSuccess: () => void
}

export default function CreateIssueForm({ projectId, executors, onSuccess }: CreateIssueFormProps) {
  const [title, setTitle] = useState('')
  const [desc, setDesc] = useState('')
  const [execId, setExecId] = useState<number | null>(null)
  const { mutate } = useCreateIssue({})

  const handleSubmit = async () => {
    try {
      await mutate({ project_id: projectId, title, description: desc, executor_id: execId } as any)
      message.success('创建成功')
      onSuccess()
    } catch (e: any) { message.error(e.message) }
  }

  return (
    <Form layout="vertical">
      <Form.Item label="标题" required>
        <Input value={title} onChange={e => setTitle(e.target.value)} />
      </Form.Item>
      <Form.Item label="描述">
        <Input.TextArea rows={3} value={desc} onChange={e => setDesc(e.target.value)} />
      </Form.Item>
      <Form.Item label="执行人">
        <Select value={execId} allowClear onChange={setExecId as any}
          options={executors.map(e => ({ label: e.name, value: e.id }))} />
      </Form.Item>
      <Button type="primary" onClick={handleSubmit}>创建</Button>
    </Form>
  )
}
