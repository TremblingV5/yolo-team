import { Form, Input, Select } from 'antd'
import { STATUSES, STATUS_NAMES } from '../constants'

export interface IssueFormValues {
  title: string
  description: string
  status: string
  priority: string
  executorId: number | null | undefined
  repoUrl: string
  branch: string
}

interface IssueFormFieldsProps {
  values: IssueFormValues
  executors: { id: number; name: string }[]
  onChange: (patch: Partial<IssueFormValues>) => void
}

export default function IssueFormFields({ values, executors, onChange }: IssueFormFieldsProps) {
  return (
    <>
      {/* Row 1: status + priority */}
      <div style={{ display: 'flex', gap: 8, alignItems: 'center', marginBottom: 16 }}>
        <Select size="small" value={values.status} onChange={(v) => onChange({ status: v })} style={{ width: 110 }}
          options={STATUSES.map(s => ({ label: STATUS_NAMES[s], value: s }))} />
        <Select size="small" value={values.priority} onChange={(v) => onChange({ priority: v })} style={{ width: 100 }}
          options={['critical', 'high', 'medium', 'low'].map(p => ({ label: p, value: p }))} />
      </div>

      {/* Row 2: title */}
      <Form.Item label="标题">
        <Input value={values.title} onChange={(e) => onChange({ title: e.target.value })} />
      </Form.Item>

      {/* Row 3: description */}
      <Form.Item label="描述">
        <Input.TextArea rows={3} value={values.description} onChange={(e) => onChange({ description: e.target.value })} />
      </Form.Item>

      {/* Row 4: executor */}
      <Form.Item label="执行人">
        <Select value={values.executorId} onChange={(v) => onChange({ executorId: v })} allowClear placeholder="选择执行人"
          options={executors.map(e => ({ label: e.name, value: e.id }))} />
      </Form.Item>

      {/* Row 5: repo URL */}
      <Form.Item label="仓库 URL">
        <Input value={values.repoUrl} onChange={(e) => onChange({ repoUrl: e.target.value })} placeholder="https://github.com/..." />
      </Form.Item>

      {/* Row 6: branch */}
      <Form.Item label="分支">
        <Input value={values.branch} onChange={(e) => onChange({ branch: e.target.value })} placeholder="main" />
      </Form.Item>
    </>
  )
}
