import { useEffect, useCallback } from 'react'
import { useDispatch, useSelector } from 'react-redux'
import { Table, Button, Popconfirm, message } from 'antd'
import { DeleteOutlined, ReloadOutlined } from '@ant-design/icons'
import { useTranslation } from 'react-i18next'
import { AppDispatch, RootState } from '../store'
import { fetchSessions, deleteSession } from '../store/sessionsSlice'
import { PageHeader, Loading, EmptyState } from '../components'

const SessionsPage: React.FC = () => {
  const dispatch = useDispatch<AppDispatch>()
  const { sessions, loading, error } = useSelector((state: RootState) => state.sessions)
  const sessionsData = sessions || []
  const { t } = useTranslation()

  useEffect(() => {
    dispatch(fetchSessions())
  }, [dispatch])

  useEffect(() => {
    if (error) {
      message.error(error)
    }
  }, [error])

  const handleDeleteSession = useCallback((sessionId: string) => {
    dispatch(deleteSession(sessionId))
      .unwrap()
      .then(() => {
        message.success(t('sessions.success'))
      })
      .catch((err: string) => {
        message.error(err || t('common.error'))
      })
  }, [dispatch, t])

  const handleRefresh = useCallback(() => {
    dispatch(fetchSessions())
  }, [dispatch])

  const columns = [
    {
      title: 'Session ID',
      dataIndex: 'id',
      key: 'id',
      ellipsis: true,
      width: '60%',
    },
    {
      title: 'Created At',
      dataIndex: 'created_at',
      key: 'created_at',
      width: '25%',
    },
    {
      title: t('common.delete'),
      key: 'action',
      width: '15%',
      render: (_: unknown, record: { id: string }) => (
        <Popconfirm
          title={t('sessions.confirmDelete')}
          description={t('sessions.deleteDesc')}
          onConfirm={() => handleDeleteSession(record.id)}
          okText={t('common.confirm')}
          cancelText={t('common.cancel')}
          okButtonProps={{ danger: true }}
        >
          <Button danger size="small" icon={<DeleteOutlined />}>
            {t('common.delete')}
          </Button>
        </Popconfirm>
      ),
    },
  ]

  if (loading && sessionsData.length === 0) {
    return <Loading size="large" />
  }

  return (
    <div>
      <PageHeader
        title={t('sessions.title')}
        showAction={false}
        extra={
          <Button
            icon={<ReloadOutlined />}
            onClick={handleRefresh}
            loading={loading}
          >
            {t('common.refresh')}
          </Button>
        }
      />

      {sessionsData.length === 0 ? (
        <EmptyState
          title={t('sessions.empty')}
          description={t('sessions.emptyDesc')}
        />
      ) : (
        <Table
          columns={columns}
          dataSource={sessionsData}
          rowKey="id"
          loading={loading}
          pagination={{ pageSize: 10, showSizeChanger: true, showTotal: (total: number) => `Total: ${total}` }}
        />
      )}
    </div>
  )
}

export default SessionsPage