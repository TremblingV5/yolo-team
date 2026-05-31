const API = '/api/v1'

export interface Project {
  id: number; key: string; name: string; description: string
  created_at: string; updated_at: string;
}
export interface Issue {
  id: number; key: string; project_id: number; parent_id: number | null
  status: string; title: string; description: string; deadline: string | null
  priority: string; executor_id: number | null; repo_url: string; repo_name: string; branch_name: string
  sort_order: number; created_at: string; updated_at: string; children: Issue[]
}
export interface Executor {
  id: number; name: string; role: string; soul: string
  created_at: string; updated_at: string;
}
export interface Document {
  id: number; key: string; project_id: number; title: string
  content: string; file_path: string
  sort_order: number; created_at: string; updated_at: string;
}

async function req<T>(path: string, options?: RequestInit): Promise<T> {
  const res = await fetch(`${API}${path}`, {
    headers: { 'Content-Type': 'application/json', ...options?.headers },
    ...options,
  })
  const json = await res.json()
  if (json.code !== 0) throw new Error(json.message || 'API error')
  return json.data as T
}

export const api = {
  projects: {
    list: () => req<Project[]>('/projects'),
    get: (key: string) => req<Project>(`/projects/${key}`),
    create: (name: string, description: string) =>
      req<Project>('/projects', { method: 'POST', body: JSON.stringify({ name, description }) }),
    update: (key: string, data: Partial<Project>) =>
      req<Project>(`/projects/${key}`, { method: 'PUT', body: JSON.stringify(data) }),
  },
  executors: {
    list: () => req<Executor[]>('/executors'),
  },
  issues: {
    list: (params?: Record<string, string | number>) => {
      const qs = params ? '?' + new URLSearchParams(params as any).toString() : ''
      return req<Issue[]>(`/issues${qs}`)
    },
    get: (key: string) => req<Issue>(`/issues/${key}`),
    create: (data: any) =>
      req<Issue>('/issues', { method: 'POST', body: JSON.stringify(data) }),
    update: (key: string, data: any) =>
      req<Issue>(`/issues/${key}`, { method: 'PUT', body: JSON.stringify(data) }),
    delete: (key: string) =>
      req<null>(`/issues/${key}`, { method: 'DELETE' }),
    todo: (executorId: number, limit = 20) =>
      req<Issue[]>(`/issues/todo?executor_id=${executorId}&limit=${limit}`),
  },
  documents: {
    list: (projectKey: string) => req<Document[]>(`/projects/${projectKey}/documents`),
    get: (key: string) => req<Document>(`/documents/${key}`),
    create: (projectKey: string, title: string) =>
      req<Document>(`/projects/${projectKey}/documents`, { method: 'POST', body: JSON.stringify({ title, content: '' }) }),
    update: (key: string, data: any) =>
      req<Document>(`/documents/${key}`, { method: 'PUT', body: JSON.stringify(data) }),
    delete: (key: string) =>
      req<null>(`/documents/${key}`, { method: 'DELETE' }),
  },
}
