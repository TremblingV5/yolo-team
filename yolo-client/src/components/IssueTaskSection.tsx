import { PlusOutlined } from '@ant-design/icons'
import { Button, Form, Input, Modal, Select, Tag, message } from 'antd'
import { useState } from 'react'
import { TASK_STATUSES, TASK_STATUS_NAMES } from '../constants'
import {
  YoloTeamYoloCliInternalModelTask as Task,
  useCreateTask, useUpdateTask, useDeleteTask,
} from '../generated'

interface IssueTaskSectionProps {
  issueKey: string
  tasks: Task[]
  onTaskChange: () => void
}

export default function IssueTaskSection({ issueKey, tasks, onTaskChange }: IssueTaskSectionProps) {
  const [taskModal, setTaskModal] = useState(false)
  const [taskTitle, setTaskTitle] = useState('')
  const [taskDesc, setTaskDesc] = useState('')
  const [taskDetail, setTaskDetail] = useState<Task | null>(null)

  const { mutate: createTaskMutate } = useCreateTask({ key: issueKey })
  const { mutate: updateTaskMutate } = useUpdateTask({ key: issueKey, task_key: '' })
  const { mutate: deleteTaskMutate } = useDeleteTask({ key: issueKey, task_key: '' })

  const handleCreateTask = async () => {
    if (!taskTitle.trim()) { message.warning('请输入任务标题'); return }
    try {
      await createTaskMutate({ title: taskTitle, description: taskDesc })
      message.success('任务创建成功')
      setTaskModal(false)
      setTaskTitle('')
      setTaskDesc('')
      onTaskChange()
    } catch (e: any) { message.error(e?.data?.message || e.message || '操作失败') }
  }

  const handleUpdateTaskStatus = async (taskKey: string, newStatus: string) => {
    try {
      await updateTaskMutate(
        { status: newStatus },
        { pathParams: { key: issueKey, task_key: taskKey } },
      )
      onTaskChange()
    } catch (e: any) { message.error(e?.data?.message || e.message || '操作失败') }
  }

  const handleDeleteTask = async (taskKey: string) => {
    try {
      await deleteTaskMutate(undefined, { pathParams: { key: issueKey, task_key: taskKey } })
      message.success('任务已删除')
      onTaskChange()
    } catch (e: any) { message.error(e?.data?.message || e.message || '操作失败') }
  }

  return (
    <>
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
                <Button size="small" type="link" onClick={() => setTaskDetail(t)}>详情</Button>
                <Button size="small" type="link" danger onClick={() => handleDeleteTask(t.key!)}>删除</Button>
              </div>
            ))}
          </div>
        )}
      </Form.Item>

      {/* Create Task Modal */}
      <Modal title="新建任务" open={taskModal}
        onCancel={() => { setTaskModal(false); setTaskTitle(''); setTaskDesc('') }}
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
    </>
  )
}
