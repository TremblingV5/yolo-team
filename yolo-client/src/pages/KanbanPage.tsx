import { useState, useCallback } from 'react'
import { Drawer, message } from 'antd'
import { useSearchParams } from 'react-router-dom'
import { useListProjects, useListIssues, useListExecutors, useListDocuments } from '../generated'
import {
  YoloTeamYoloCliInternalModelIssue as Issue,
  ListIssuesQueryParams,
} from '../generated'
import TopBar from '../components/TopBar'
import KanbanBoard from '../components/KanbanBoard'
import ExecutorKanban from '../components/ExecutorKanban'
import IssueDrawer from '../components/IssueDrawer'
import CreateIssueForm from '../components/CreateIssueForm'
import CreateProjectModal from '../components/CreateProjectModal'
import ProjectEditModal from '../components/ProjectEditModal'

export default function KanbanPage() {
  const [searchParams, setSearchParams] = useSearchParams()
  const selectedProject = searchParams.get('project') || undefined

  const [drawerOpen, setDrawerOpen] = useState(false)
  const [drawerIssue, setDrawerIssue] = useState<Issue | null>(null)
  const [isCreating, setIsCreating] = useState(false)
  const [projectModalOpen, setProjectModalOpen] = useState(false)
  const [editProjectKey, setEditProjectKey] = useState<string | null>(null)
  const [localIssues, setLocalIssues] = useState<Issue[]>([])
  const [viewMode, setViewMode] = useState<'status' | 'executor'>('status')

  const { data: projectsData, refetch: refetchProjects } = useListProjects({})
  const projects = (projectsData as any)?.data || []

  const issueParams: ListIssuesQueryParams = {}
  if (selectedProject) issueParams.project_id = Number(selectedProject)

  const { data: issuesData, refetch: refetchIssues } = useListIssues({ queryParams: issueParams })
  const issues = (issuesData as any)?.data || []

  const { data: executorsData } = useListExecutors({})
  const executors = (executorsData as any)?.data || []

  const selectedProjectObj = projects.find((p: any) => p.id === Number(selectedProject))
  const { data: docsData } = useListDocuments({ key: selectedProjectObj?.key || '' })
  const projectDocs = selectedProject ? ((docsData as any)?.data || []) : []

  const displayIssues = localIssues.length > 0 ? localIssues : issues

  const setProject = useCallback((v: string | undefined) => {
    if (v) {
      setSearchParams({ project: v })
    } else {
      setSearchParams({})
    }
    setLocalIssues([])
  }, [setSearchParams])

  const load = useCallback(() => {
    refetchProjects()
    refetchIssues()
    setLocalIssues([])
  }, [refetchProjects, refetchIssues])

  const updateIssue = async (key: string, body: Record<string, any>) => {
    const resp = await fetch(`/api/v1/issues/${key}`, {
      method: 'PUT',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(body),
    })
    const json = await resp.json()
    if (json.code !== 0) throw new Error(json.message || '更新失败')
    return json
  }

  const handleDragEnd = async (issueKey: string, newStatus: string) => {
    const current = displayIssues
    const issue = current.find((i: Issue) => i.key === issueKey)
    if (!issue || issue.status === newStatus) return

    const previous = [...current]
    const updated = current.map((i: Issue) =>
      i.key === issueKey ? { ...i, status: newStatus } : i
    )
    setLocalIssues(updated)

    try {
      await updateIssue(issueKey, { status: newStatus })
      load()
    } catch (e: any) {
      setLocalIssues(previous)
      message.error(e.message || '更新失败，已回滚')
    }
  }

  const handleArchive = async (key: string) => {
    const current = displayIssues
    const previous = [...current]
    const updated = current.map((i: Issue) =>
      i.key === key ? { ...i, status: 'archived' } : i
    )
    setLocalIssues(updated)

    try {
      await updateIssue(key, { status: 'archived' })
      message.success('已归档')
      load()
    } catch (e: any) {
      setLocalIssues(previous)
      message.error(e.message || '归档失败，已回滚')
    }
  }

  const handleAssign = async (key: string, executorId: number) => {
    const current = displayIssues
    const previous = [...current]
    const updated = current.map((i: Issue) =>
      i.key === key ? { ...i, executor_id: executorId } : i
    )
    setLocalIssues(updated)

    try {
      await updateIssue(key, { executor_id: executorId })
      load()
    } catch (e: any) {
      setLocalIssues(previous)
      message.error(e.message || '指派失败，已回滚')
    }
  }

  const handleCardClick = async (issue: Issue) => {
    try {
      const resp = await fetch(`/api/v1/issues/${issue.key}`)
      const json = await resp.json()
      if (json.code === 0) {
        setDrawerIssue(json.data)
        setIsCreating(false)
        setDrawerOpen(true)
      }
    } catch (e: any) { message.error(e.message) }
  }

  const handleCreateClick = () => {
    setDrawerIssue(null)
    setIsCreating(true)
    setDrawerOpen(true)
  }

  const handleDrawerClose = () => {
    setDrawerOpen(false)
    setDrawerIssue(null)
    setIsCreating(false)
  }

  const handleSaved = () => {
    setDrawerOpen(false)
    setIsCreating(false)
    load()
  }

  return (
    <div>
      <TopBar
        projects={projects}
        selectedProject={selectedProject}
        viewMode={viewMode}
        onSelectProject={setProject}
        onCreateProject={() => setProjectModalOpen(true)}
        onEditProject={() => {
          if (selectedProjectObj) setEditProjectKey(selectedProjectObj.key)
        }}
        onToggleView={() => setViewMode(v => v === 'status' ? 'executor' : 'status')}
      />

      {viewMode === 'status' ? (
        <KanbanBoard
          issues={displayIssues}
          executors={executors}
          selectedProject={selectedProject}
          onCardClick={handleCardClick}
          onCreateClick={handleCreateClick}
          onArchive={handleArchive}
          onAssign={handleAssign}
          onDragEnd={handleDragEnd}
          activeId={null}
        />
      ) : (
        <ExecutorKanban
          issues={displayIssues}
          executors={executors}
          onCardClick={handleCardClick}
          onArchive={handleArchive}
          onAssign={handleAssign}
        />
      )}

      <Drawer
        title={isCreating ? '新建 Issue' : 'Issue 详情'}
        placement="left"
        width={520}
        open={drawerOpen}
        onClose={handleDrawerClose}
      >
        {isCreating ? (
          selectedProject ? (
            <CreateIssueForm
              projectId={Number(selectedProject)}
              executors={executors}
              onSuccess={handleSaved}
            />
          ) : (
            <div>请先选择一个项目</div>
          )
        ) : drawerIssue ? (
          <IssueDrawer
            issue={drawerIssue}
            executors={executors}
            projectDocs={projectDocs}
            projectKey={selectedProjectObj?.key || ''}
            open={drawerOpen}
            onClose={handleDrawerClose}
            onSaved={handleSaved}
          />
        ) : null}
      </Drawer>

      <CreateProjectModal
        open={projectModalOpen}
        onClose={() => setProjectModalOpen(false)}
        onDone={load}
      />

      {editProjectKey && (
        <ProjectEditModal
          projectKey={editProjectKey}
          open={!!editProjectKey}
          onClose={() => setEditProjectKey(null)}
        />
      )}
    </div>
  )
}
