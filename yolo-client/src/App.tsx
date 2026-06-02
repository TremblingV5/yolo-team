import { Routes, Route } from 'react-router-dom'
import KanbanPage from './pages/KanbanPage'
import DocumentList from './pages/DocumentList'
import DocumentDetail from './pages/DocumentDetail'
import ExecutorManage from './pages/ExecutorManage'

export default function App() {
  return (
    <Routes>
      <Route path="/" element={<KanbanPage />} />
      <Route path="/project/:key/docs" element={<DocumentList />} />
      <Route path="/doc/:key" element={<DocumentDetail />} />
      <Route path="/executors" element={<ExecutorManage />} />
    </Routes>
  )
}
