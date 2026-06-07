import { useState } from 'react'
import { Button } from 'antd'
import { RobotOutlined } from '@ant-design/icons'
import { Routes, Route } from 'react-router-dom'
import KanbanPage from './pages/KanbanPage'
import DocumentList from './pages/DocumentList'
import DocumentDetail from './pages/DocumentDetail'
import ExecutorManage from './pages/ExecutorManage'
import AiChatDrawer from './components/AiChatDrawer'

export default function App() {
  const [drawerOpen, setDrawerOpen] = useState(false)

  return (
    <>
      <Routes>
        <Route path="/" element={<KanbanPage />} />
        <Route path="/project/:key/docs" element={<DocumentList />} />
        <Route path="/doc/:key" element={<DocumentDetail />} />
        <Route path="/executors" element={<ExecutorManage />} />
      </Routes>

      <Button
        type="primary"
        shape="circle"
        size="large"
        icon={<RobotOutlined style={{ fontSize: 22 }} />}
        onClick={() => setDrawerOpen(true)}
        style={{
          position: 'fixed',
          right: 32,
          bottom: 32,
          width: 56,
          height: 56,
          boxShadow: '0 4px 12px rgba(0,0,0,0.25)',
          zIndex: 1000,
        }}
      />

      <AiChatDrawer open={drawerOpen} onClose={() => setDrawerOpen(false)} />
    </>
  )
}
