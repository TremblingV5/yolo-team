import { Select, Button, Switch } from 'antd'
import { PlusOutlined, EditOutlined, UserOutlined } from '@ant-design/icons'
import { useNavigate } from 'react-router-dom'
import {
  YoloTeamYoloCliInternalModelProject as Project,
} from '../generated'

interface TopBarProps {
  projects: Project[]
  selectedProject: string | undefined
  viewMode: 'status' | 'executor'
  onSelectProject: (v: string | undefined) => void
  onCreateProject: () => void
  onEditProject: () => void
  onToggleView: () => void
}

export default function TopBar({ projects, selectedProject, viewMode, onSelectProject, onCreateProject, onEditProject, onToggleView }: TopBarProps) {
  const navigate = useNavigate()

  return (
    <div className="topbar">
      <Select
        style={{ width: 220 }}
        placeholder="全部项目"
        allowClear
        value={selectedProject ? Number(selectedProject) : undefined}
        onClear={() => onSelectProject(undefined)}
        onChange={(v) => onSelectProject(v ? String(v) : undefined)}
        options={projects.map(p => ({ label: p.name, value: p.id }))}
      />
      <Button type="primary" icon={<PlusOutlined />} onClick={onCreateProject}>新建项目</Button>
      {selectedProject && (
        <Button icon={<EditOutlined />} onClick={onEditProject}>编辑</Button>
      )}
      {selectedProject && (
        <Button onClick={() => {
          const proj = projects.find(p => p.id === Number(selectedProject))
          if (proj) navigate(`/project/${proj.key}/docs`)
        }}>文档</Button>
      )}
      <Button icon={<UserOutlined />} onClick={() => navigate('/executors')}>执行人</Button>

      <div style={{ flex: 1 }} />

      <span style={{ marginRight: 8, color: '#666', fontSize: 13 }}>
        {viewMode === 'status' ? '状态视角' : '执行人视角'}
      </span>
      <Switch
        checked={viewMode === 'executor'}
        onChange={onToggleView}
        checkedChildren="执行人"
        unCheckedChildren="状态"
      />
    </div>
  )
}
