import { Button } from 'antd'
import { useDroppable, useDraggable, DndContext, DragEndEvent, PointerSensor, useSensor, useSensors } from '@dnd-kit/core'
import { STATUSES, STATUS_NAMES } from '../constants'
import IssueCard from './IssueCard'
import {
  YoloTeamYoloCliInternalModelIssue as Issue,
  YoloTeamYoloCliInternalModelExecutor as Executor,
} from '../generated'

interface KanbanBoardProps {
  issues: Issue[]
  executors: Executor[]
  selectedProject: string | undefined
  onCardClick: (issue: Issue) => void
  onCreateClick: () => void
  onArchive: (key: string) => void
  onAssign: (key: string, executorId: number) => void
  onDragEnd: (issueKey: string, newStatus: string) => void
  activeId: string | null
}

function DraggableIssueCard({
  issue,
  executors,
  onCardClick,
  onArchive,
  onAssign,
}: {
  issue: Issue
  executors: Executor[]
  onCardClick: (i: Issue) => void
  onArchive: (k: string) => void
  onAssign: (k: string, id: number) => void
}) {
  const { attributes, listeners, setNodeRef, transform, isDragging } = useDraggable({
    id: issue.key || '',
    data: { issue },
  })

  const style = transform ? {
    transform: `translate(${transform.x}px, ${transform.y}px)`,
    zIndex: 999,
  } : undefined

  return (
    <div ref={setNodeRef} style={style} {...listeners} {...attributes}>
      <IssueCard
        issue={issue}
        executors={executors}
        onClick={() => onCardClick(issue)}
        onArchive={onArchive}
        onAssign={onAssign}
        isDragging={isDragging}
      />
    </div>
  )
}

function KanbanColumn({
  status,
  issues,
  executors,
  onCardClick,
  onArchive,
  onAssign,
  onCreateClick,
  selectedProject,
}: {
  status: string
  issues: Issue[]
  executors: Executor[]
  onCardClick: (i: Issue) => void
  onArchive: (k: string) => void
  onAssign: (k: string, id: number) => void
  onCreateClick: () => void
  selectedProject: string | undefined
}) {
  const { setNodeRef, isOver } = useDroppable({ id: status })

  return (
    <div ref={setNodeRef} className={`kanban-col${isOver ? ' over' : ''}`}>
      <div className={`kanban-col-header col-${status}`}>
        {STATUS_NAMES[status]} ({issues.length})
      </div>
      {status === 'created' && selectedProject && (
        <Button block size="small" style={{ marginBottom: 8 }} onClick={onCreateClick}>
          + 新建
        </Button>
      )}
      {issues.map(issue => (
        <DraggableIssueCard
          key={issue.key}
          issue={issue}
          executors={executors}
          onCardClick={onCardClick}
          onArchive={onArchive}
          onAssign={onAssign}
        />
      ))}
    </div>
  )
}

export default function KanbanBoard(props: KanbanBoardProps) {
  const sensors = useSensors(
    useSensor(PointerSensor, { activationConstraint: { distance: 5 } }),
  )
  const columns = STATUSES.map(status =>
    props.issues.filter(i => i.status === status)
  )

  const handleDragEnd = (event: DragEndEvent) => {
    const { active, over } = event
    if (!over) return
    const newStatus = over.id as string
    if (STATUSES.includes(newStatus)) {
      props.onDragEnd(active.id as string, newStatus)
    }
  }

  return (
    <DndContext onDragEnd={handleDragEnd} sensors={sensors}>
      <div className="kanban-container">
        {STATUSES.map((status, idx) => (
          <KanbanColumn
            key={status}
            status={status}
            issues={columns[idx]}
            executors={props.executors}
            onCardClick={props.onCardClick}
            onArchive={props.onArchive}
            onAssign={props.onAssign}
            onCreateClick={props.onCreateClick}
            selectedProject={props.selectedProject}
          />
        ))}
      </div>
    </DndContext>
  )
}
