export const STATUSES = ['created', 'in_progress', 'done', 'archived']
export const STATUS_NAMES: Record<string, string> = {
  created: '已创建', in_progress: '执行中', done: '已完成', archived: '已归档',
}
export const PRIORITY_COLORS: Record<string, string> = { critical: 'red', high: 'orange', medium: 'blue', low: 'green' }
