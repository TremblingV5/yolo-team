import { useState, useEffect } from 'react'
import { Button, Input, message } from 'antd'
import { useParams, useNavigate, useSearchParams } from 'react-router-dom'
import MDEditor from '@uiw/react-md-editor'
import { useGetDocument, useUpdateDocument } from '../generated'

export default function DocumentDetail() {
  const { key } = useParams<{ key: string }>()
  const navigate = useNavigate()
  const [searchParams] = useSearchParams()
  const linkTo = searchParams.get('linkTo')
  const projectKey = searchParams.get('projectKey')
  const isNew = key === 'new'

  const [title, setTitle] = useState(searchParams.get('title') || '')
  const [content, setContent] = useState('')
  const [editing, setEditing] = useState(!!isNew || !!linkTo)
  const [saving, setSaving] = useState(false)

  const { data } = useGetDocument({ key: isNew ? '' : (key || '') })
  const doc = (data as any)?.data
  const { mutate: updateMutate } = useUpdateDocument({ key: isNew ? '' : (key || '') })

  useEffect(() => {
    if (doc && !isNew) {
      setTitle(doc.title || '')
      setContent(doc.content || '')
    }
  }, [doc, isNew])

  const handleSave = async () => {
    setSaving(true)
    try {
      let docKey = key
      let newDoc: any = null

      if (isNew && projectKey) {
        const resp = await fetch(`/api/v1/projects/${projectKey}/documents`, {
          method: 'POST',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify({ title: title || '未命名文档', content, creator: '人类' }),
        })
        const json = await resp.json()
        if (json.code !== 0) throw new Error(json.message || '创建失败')
        docKey = json.data.key
        newDoc = json.data
      } else if (!isNew) {
        await updateMutate({ content } as any)
      }

      if (linkTo && docKey) {
        const docId = newDoc ? newDoc.id : (doc?.id)
        if (docId) {
          await fetch(`/api/v1/issues/${linkTo}/documents`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ document_id: docId }),
          })
        }
        message.success('已保存并关联到 Issue')
      } else {
        message.success('已保存')
      }

      setEditing(false)
      setSaving(false)
      if (isNew && docKey) {
        navigate(`/doc/${docKey}${linkTo ? `?linkTo=${linkTo}` : ''}`, { replace: true })
      }
    } catch (e: any) {
      message.error(e.message)
      setSaving(false)
    }
  }

  const displayTitle = isNew ? (title || '新建文档') : (doc?.title || '')
  const displayPath = isNew ? '新建中...' : (doc?.file_path || '')
  const displayCreator = isNew ? '人类' : (doc?.creator || '')

  if (!isNew && !doc) return <div>Loading...</div>

  return (
    <div style={{ padding: 24 }}>
      <div style={{ display: 'flex', alignItems: 'center', marginBottom: 12 }}>
        <div style={{ flex: 1 }}>
          {(isNew || editing) ? (
            <Input
              size="large"
              value={title}
              onChange={e => setTitle(e.target.value)}
              placeholder="文档标题"
              style={{ width: 400, fontWeight: 600, fontSize: 20 }}
            />
          ) : (
            <h2 style={{ margin: 0 }}>{displayTitle}</h2>
          )}
          <p style={{ color: '#999', margin: '4px 0 0' }}>
            路径: {displayPath}
            {displayCreator && <span style={{ marginLeft: 16 }}>创建人: {displayCreator}</span>}
          </p>
        </div>
        <div style={{ display: 'flex', gap: 8 }}>
          {editing ? (
            <>
              <Button type="primary" onClick={handleSave} loading={saving}>保存</Button>
              <Button onClick={() => { setEditing(false); if (isNew) navigate(-1) }}>取消</Button>
            </>
          ) : (
            <Button type="primary" onClick={() => setEditing(true)}>编辑</Button>
          )}
          <Button onClick={() => navigate(-1)}>返回</Button>
        </div>
      </div>

      {editing ? (
        <div data-color-mode="light">
          <MDEditor value={content} onChange={v => setContent(v || '')} height={600} />
        </div>
      ) : (
        <div data-color-mode="light">
          <MDEditor.Markdown source={content || '_暂无内容_'} />
        </div>
      )}
    </div>
  )
}
