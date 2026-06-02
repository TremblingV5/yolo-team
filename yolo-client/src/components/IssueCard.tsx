import { Tag, Select, Button, message } from 'antd'
import dayjs from 'dayjs'
import { PRIORITY_COLORS } from '../constants'
import {
  YoloTeamYoloCliInternalModelIssue as Issue,
  YoloTeamYoloCliInternalModelExecutor as Executor,
} from '../generated'

interface IssueCardProps {
  issue: Issue
  executors: Executor[]
  onClick: () => void
  onArchive: (key: string) => void
  onAssign: (key: string, executorId: number) => void
  isDragging?: boolean
}

export default function IssueCard({ issue, executors, onClick, onArchive, onAssign, isDragging }: IssueCardProps) {
  return (
    <div
      className={`issue-card${isDragging ? ' dragging' : ''}`}
      onClick={onClick}
      style={{ opacity: isDragging ? 0.5 : 1 }}
    >
      <div className="title">{issue.title}</div>
      <div className="meta">
        <Tag color={PRIORITY_COLORS[issue.priority || '']}>{issue.priority}</Tag>
        {issue.deadline && <span>{dayjs(issue.deadline).format('MM-DD')}</span>}
        {issue.repo_name && <span>🔗</span>}
      </div>
      <div className="actions" onClick={e => e.stopPropagation()}>
        <Select
          size="small"
          style={{ width: 120 }}
          placeholder="指派"
          allowClear
          value={issue.executor_id}
          onChange={(v) => onAssign(issue.key || '', v)}
          options={executors.map(e => ({ label: e.name, value: e.id }))}
        />
        {issue.status === 'done' && (
          <Button
            size="small"
            type="link"
            danger
            style={{ marginLeft: 4 }}
            onClick={() => onArchive(issue.key || '')}
          >
            归档
          </Button>
        )}
      </div>
    </div>
  )
}
