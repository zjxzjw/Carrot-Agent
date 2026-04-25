import { useState, useEffect, useCallback } from 'react'
import { useDispatch, useSelector } from 'react-redux'
import { Table, Modal, Form, Input, Select, message } from 'antd'
import { useTranslation } from 'react-i18next'
import { AppDispatch, RootState } from '../store'
import { fetchMemories, addMemory } from '../store/memorySlice'
import { PageHeader, Loading, EmptyState } from '../components'
import { MEMORY_TYPE_OPTIONS } from '../types'

const { TextArea } = Input

const MemoryPage: React.FC = () => {
  const [isModalOpen, setIsModalOpen] = useState(false)
  const [form] = Form.useForm()
  const [memoryType, setMemoryType] = useState('')
  const dispatch = useDispatch<AppDispatch>()
  const { memories, loading, error } = useSelector((state: RootState) => state.memory)
  const { t } = useTranslation()

  useEffect(() => {
    dispatch(fetchMemories(memoryType))
  }, [dispatch, memoryType])

  useEffect(() => {
    if (error) {
      message.error(error)
    }
  }, [error])

  const handleAddMemory = useCallback(() => {
    form.validateFields().then((values) => {
      dispatch(addMemory(values))
        .unwrap()
        .then(() => {
          message.success(t('memory.success'))
          setIsModalOpen(false)
          form.resetFields()
          dispatch(fetchMemories(memoryType))
        })
        .catch((err: string) => {
          message.error(err || t('common.error'))
        })
    })
  }, [form, dispatch, memoryType, t])

  const handleOpenModal = useCallback(() => {
    setIsModalOpen(true)
  }, [])

  const handleCloseModal = useCallback(() => {
    setIsModalOpen(false)
    form.resetFields()
  }, [form])

  const columns = [
    {
      title: t('memory.type'),
      dataIndex: 'type',
      key: 'type',
      width: 120,
      filters: [
        { text: '快照记忆', value: 'snapshot' },
        { text: '会话记忆', value: 'session' },
        { text: '长期记忆', value: 'longterm' },
      ],
      onFilter: (value: unknown, record: { type: string }) => record.type === value,
      render: (type: string) => {
        const typeMap: Record<string, { text: string; color: string }> = {
          snapshot: { text: '快照', color: 'blue' },
          session: { text: '会话', color: 'green' },
          longterm: { text: '长期', color: 'orange' },
        }
        const config = typeMap[type] || { text: type, color: 'default' }
        return <span style={{ color: config.color }}>{config.text}</span>
      },
    },
    {
      title: t('memory.content'),
      dataIndex: 'content',
      key: 'content',
      ellipsis: true,
    },
    {
      title: t('memory.metadata'),
      dataIndex: 'metadata',
      key: 'metadata',
      width: 150,
      ellipsis: true,
      render: (metadata: string) => metadata || '-',
    },
    {
      title: 'Created At',
      dataIndex: 'created_at',
      key: 'created_at',
      width: 180,
    },
  ]

  if (loading && memories.length === 0) {
    return <Loading size="large" />
  }

  return (
    <div>
      <PageHeader
        title={t('memory.title')}
        actionText={t('memory.add')}
        onAction={handleOpenModal}
        extra={
          <Select
            value={memoryType}
            onChange={setMemoryType}
            style={{ width: 160 }}
            options={MEMORY_TYPE_OPTIONS}
          />
        }
      />

      {memories.length === 0 ? (
        <EmptyState
          title={t('memory.empty')}
          description={t('memory.emptyDesc')}
          actionText={t('memory.add')}
          onAction={handleOpenModal}
        />
      ) : (
        <Table
          columns={columns}
          dataSource={memories}
          rowKey="id"
          loading={loading}
          pagination={{ pageSize: 10, showSizeChanger: true, showQuickJumper: true }}
        />
      )}

      <Modal
        title={t('memory.add')}
        open={isModalOpen}
        onOk={handleAddMemory}
        onCancel={handleCloseModal}
        okText={t('common.create')}
        cancelText={t('common.cancel')}
        destroyOnClose
      >
        <Form form={form} layout="vertical" preserve={false}>
          <Form.Item
            name="type"
            label={t('memory.type')}
            rules={[{ required: true, message: t('memory.type') + t('common.required') }]}
          >
            <Select
              placeholder={t('memory.type')}
              options={MEMORY_TYPE_OPTIONS.filter(opt => opt.value !== '')}
            />
          </Form.Item>
          <Form.Item
            name="content"
            label={t('memory.content')}
            rules={[{ required: true, message: t('memory.content') + t('common.required') }]}
          >
            <TextArea rows={4} placeholder={t('memory.content') + '...'} />
          </Form.Item>
          <Form.Item
            name="metadata"
            label={t('memory.metadata')}
          >
            <Input placeholder="JSON格式" />
          </Form.Item>
        </Form>
      </Modal>
    </div>
  )
}

export default MemoryPage