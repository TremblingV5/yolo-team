import { Tag } from 'antd'
import { STATUS_NAMES } from '../constants'
import IssueCard from './IssueCard'
import {
  YoloTeamYoloCliInternalModelIssue as Issue,
  YoloTeamYoloCliInternalModelExecutor as Executor,
} from '../generated'

const EXEC_COLORS = ['#1890FF', '#52C41A', '#FA8C16', '#722ED1', '#EB2F96', '#13C2C2', '#F5222D', '#2F54EB']

interface ExecutorKanbanProps {
  issues: Issue[]
  executors: Executor[]
  onCardClick: (issue: Issue) => void
  onArchive: (key: string) => void
  onAssign: (key: string, executorId: number) => void
}

export default function ExecutorKanban({
  issues, executors, onCardClick, onArchive, onAssign,
}: ExecutorKanbanProps) {
  const activeIssues = issues.filter(i => i.status === 'created' || i.status === 'in_progress')

  return (
    <div className="kanban-container">
      {executors.map((exec, idx) => {
        const execIssues = activeIssues.filter(i => i.executor_id === exec.id)
        const color = EXEC_COLORS[idx % EXEC_COLORS.length]
        return (
          <div key={exec.id} className="kanban-col">
            <div className="kanban-col-header" style={{ background: color }}>
              {exec.name} ({execIssues.length})
            </div>
            {execIssues.map(issue => (
              <div key={issue.key}>
                <IssueCard
                  issue={issue}
                  executors={executors}
                  onClick={() => onCardClick(issue)}
                  onArchive={onArchive}
                  onAssign={onAssign}
                />
              </div>
            ))}
            {execIssues.length === 0 && (
              <div style={{ color: '#bbb', fontSize: 12, padding: 12, textAlign: 'center' }}>无任务</div>
            )}
          </div>
        )
      })}
    </div>
  )
}
