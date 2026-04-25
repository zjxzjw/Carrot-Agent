import { useState } from 'react'
import { Card, List, Tag, Modal, Form, Input, Select, Empty, Typography, Popconfirm, message, Button } from 'antd'
import { DeleteOutlined, DatabaseOutlined, AimOutlined, BulbOutlined } from '@ant-design/icons'
import type { Memory } from '../store'
import './MemoryList.css'

const { Text, Paragraph } = Typography

interface MemoryListProps {
  memories: Memory[]
  loading: boolean
  onRefresh?: () => void
  onAdd: (type: string, content: string, metadata?: string) => Promise<void>
}

const MEMORY_TYPE_CONFIG = {
  snapshot: { color: 'blue', label: '快照', icon: <AimOutlined /> },
  session: { color: 'green', label: '会话', icon: <DatabaseOutlined /> },
  longterm: { color: 'purple', label: '长期', icon: <BulbOutlined /> }
}

export function MemoryList({ memories, loading, onAdd }: MemoryListProps) {
  const [modalVisible, setModalVisible] = useState(false)
  const [form] = Form.useForm()

  const handleAdd = async () => {
    try {
      const values = await form.validateFields()
      await onAdd(values.type, values.content, values.metadata)
      form.resetFields()
      setModalVisible(false)
      message.success('记忆添加成功')
    } catch (e) {
      if (e instanceof Error) {
        message.error(e.message)
      }
    }
  }

  const getMemoryTypeConfig = (type: string) => {
    return MEMORY_TYPE_CONFIG[type as keyof typeof MEMORY_TYPE_CONFIG] || {
      color: 'default',
      label: type,
      icon: null
    }
  }

  return (
    <div className="memory-list-container">
      <Card
        title={
          <span>
            <DatabaseOutlined style={{ color: '#1890ff', marginRight: 8 }} />
            记忆管理
          </span>
        }
        extra={
          <Select
            placeholder="筛选类型"
            style={{ width: 120 }}
            allowClear
            onChange={() => message.info('筛选功能开发中')}
          >
            <Select.Option value="snapshot">快照</Select.Option>
            <Select.Option value="session">会话</Select.Option>
            <Select.Option value="longterm">长期</Select.Option>
          </Select>
        }
        loading={loading}
      >
        {memories.length === 0 ? (
          <Empty description="暂无记忆" image={Empty.PRESENTED_IMAGE_SIMPLE} />
        ) : (
          <List
            dataSource={memories}
            renderItem={(memory) => {
              const config = getMemoryTypeConfig(memory.type)
              return (
                <List.Item
                  className="memory-item"
                  actions={[
                    <Popconfirm
                      key="delete"
                      title="确定删除此记忆？"
                      onConfirm={() => message.info('删除功能开发中')}
                      okText="确定"
                      cancelText="取消"
                    >
                      <DeleteOutlined style={{ color: '#ff4d4f', cursor: 'pointer' }} />
                    </Popconfirm>
                  ]}
                >
                  <List.Item.Meta
                    avatar={
                      <Tag color={config.color} icon={config.icon}>
                        {config.label}
                      </Tag>
                    }
                    title={
                      <div className="memory-title">
                        <Text type="secondary">
                          创建于: {new Date(memory.created_at).toLocaleString()}
                        </Text>
                      </div>
                    }
                    description={
                      <Paragraph
                        ellipsis={{ rows: 3 }}
                        className="memory-content"
                      >
                        {memory.content}
                      </Paragraph>
                    }
                  />
                </List.Item>
              )
            }}
          />
        )}
      </Card>

      <Modal
        title="添加新记忆"
        open={modalVisible}
        onOk={handleAdd}
        onCancel={() => {
          setModalVisible(false)
          form.resetFields()
        }}
        footer={null}
      >
        <Form form={form} layout="vertical">
          <Form.Item
            name="type"
            label="记忆类型"
            rules={[{ required: true, message: '请选择记忆类型' }]}
          >
            <Select placeholder="选择记忆类型">
              <Select.Option value="snapshot">快照 - 重要用户信息</Select.Option>
              <Select.Option value="session">会话 - 当前会话信息</Select.Option>
              <Select.Option value="longterm">长期 - 持久化信息</Select.Option>
            </Select>
          </Form.Item>
          <Form.Item
            name="content"
            label="记忆内容"
            rules={[{ required: true, message: '请输入记忆内容' }]}
          >
            <Input.TextArea placeholder="输入要记住的内容..." rows={4} />
          </Form.Item>
          <Form.Item name="metadata" label="元数据（可选）">
            <Input.TextArea placeholder='{"key": "value"}' rows={2} />
          </Form.Item>
          <Button type="primary" onClick={handleAdd} block>
            添加记忆
          </Button>
        </Form>
      </Modal>
    </div>
  )
}