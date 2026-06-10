import { Drawer, message } from 'antd'
import { useCallback, useState } from 'react'
import { useSearchParams } from 'react-router-dom'
import { useGet, useMutate } from 'restful-react'
import CreateIssueForm from '../components/CreateIssueForm'
import CreateProjectModal from '../components/CreateProjectModal'
import ExecutorKanban from '../components/ExecutorKanban'
import IssueDrawer from '../components/IssueDrawer'
import KanbanBoard from '../components/KanbanBoard'
import ProjectEditModal from '../components/ProjectEditModal'
import TopBar from '../components/TopBar'
import {
    YoloTeamYoloCliInternalCommonResponse as CommonResponse,
    YoloTeamYoloCliInternalModelIssue as Issue,
    ListIssuesQueryParams, useListDocuments, useListExecutors, useListIssues, useListProjects
} from '../generated'

interface UpdateIssueResponse extends CommonResponse { data?: Issue }
interface GetIssueResponse extends CommonResponse { data?: Issue }

export default function KanbanPage() {
  const [searchParams, setSearchParams] = useSearchParams()
  const selectedProject = searchParams.get('project') || undefined

  const [drawerOpen, setDrawerOpen] = useState(false)
  const [drawerIssue, setDrawerIssue] = useState<Issue | null>(null)
  const [isCreating, setIsCreating] = useState(false)
  const [issueNavStack, setIssueNavStack] = useState<string[]>([])
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

  // Dynamic-key API helpers using restful-react's useMutate / useGet
  const { mutate: updateIssueMutate } = useMutate<
    UpdateIssueResponse, CommonResponse, void, Record<string, any>, { key: string }
  >('PUT', (params) => `/api/v1/issues/${params.key}`)

  const { refetch: getIssue } = useGet<GetIssueResponse, CommonResponse, void, { key: string }>(
    (params) => `/api/v1/issues/${params.key}`,
    { lazy: true }
  )

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
    try {
      await updateIssueMutate(body, { pathParams: { key } })
    } catch (e: any) {
      throw new Error(e?.data?.message || e.message || '更新失败')
    }
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
      const result = await getIssue({ pathParams: { key: issue.key || '' } })
      const json = (result as any) || {}
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
    if (issueNavStack.length > 0) {
      // Go back to parent issue
      const prevKey = issueNavStack[issueNavStack.length - 1]
      setIssueNavStack((s) => s.slice(0, -1))
      getIssue({ pathParams: { key: prevKey } }).then((result: any) => {
        const json = (result as any) || {}
        if (json.code === 0) setDrawerIssue(json.data)
      }).catch(() => {
        setDrawerOpen(false)
        setDrawerIssue(null)
      })
    } else {
      setDrawerOpen(false)
      setDrawerIssue(null)
    }
    setIsCreating(false)
  }

  const handleSaved = () => {
    // Refresh drawer issue if drawer is open
    if (drawerIssue?.key) {
      getIssue({ pathParams: { key: drawerIssue.key } }).then((result: any) => {
        const json = (result as any) || {}
        if (json.code === 0) setDrawerIssue(json.data)
      })
    }
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
            onChildClick={async (key) => {
              setIssueNavStack((s) => [...s, drawerIssue?.key || ''])
              // Close drawer first for animation, then open with new issue
              setDrawerOpen(false)
              setTimeout(async () => {
                try {
                  const result = await getIssue({ pathParams: { key } })
                  const json = (result as any) || {}
                  if (json.code === 0) {
                    setDrawerIssue(json.data)
                    setDrawerOpen(true)
                  }
                } catch {}
              }, 150)
            }}
            onLinkChild={async (childKey, parentKey) => {
              try {
                const result = await updateIssueMutate({ parent_key: parentKey } as any, { pathParams: { key: childKey } })
                return (result as any)?.code === 0
              } catch { return false }
            }}
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
