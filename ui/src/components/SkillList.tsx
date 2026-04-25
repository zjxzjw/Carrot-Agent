import { useState } from 'react'
import { Card, List, Button, Tag, Modal, Form, Input, Empty, Typography, Popconfirm, message } from 'antd'
import { PlusOutlined, DeleteOutlined, ThunderboltOutlined } from '@ant-design/icons'
import type { Skill } from '../store'
import './SkillList.css'

const { Text, Paragraph } = Typography

interface SkillListProps {
  skills: Skill[]
  loading: boolean
  onRefresh?: () => void
  onCreate: (name: string, description: string, content: string) => Promise<void>
}

export function SkillList({ skills, loading, onCreate }: SkillListProps) {
  const [modalVisible, setModalVisible] = useState(false)
  const [form] = Form.useForm()

  const handleCreate = async () => {
    try {
      const values = await form.validateFields()
      await onCreate(values.name, values.description, values.content)
      form.resetFields()
      setModalVisible(false)
      message.success('技能创建成功')
    } catch (e) {
      if (e instanceof Error) {
        message.error(e.message)
      }
    }
  }

  const parsePlatforms = (platforms: string): string[] => {
    try {
      return JSON.parse(platforms)
    } catch {
      return []
    }
  }

  return (
    <div className="skill-list-container">
      <Card
        title={
          <span>
            <ThunderboltOutlined style={{ color: '#faad14', marginRight: 8 }} />
            技能列表
          </span>
        }
        extra={
          <Button type="primary" icon={<PlusOutlined />} onClick={() => setModalVisible(true)}>
            创建技能
          </Button>
        }
        loading={loading}
      >
        {skills.length === 0 ? (
          <Empty description="暂无技能" image={Empty.PRESENTED_IMAGE_SIMPLE} />
        ) : (
          <List
            dataSource={skills}
            renderItem={(skill) => (
              <List.Item
                className="skill-item"
                actions={[
                  <Popconfirm
                    key="delete"
                    title="确定删除此技能？"
                    onConfirm={() => message.info('删除功能开发中')}
                    okText="确定"
                    cancelText="取消"
                  >
                    <Button type="text" danger icon={<DeleteOutlined />} size="small" />
                  </Popconfirm>
                ]}
              >
                <List.Item.Meta
                  title={
                    <div className="skill-title">
                      <Text strong>{skill.name}</Text>
                      <Tag color="blue">v{skill.version}</Tag>
                    </div>
                  }
                  description={
                    <div className="skill-description">
                      <Paragraph type="secondary" ellipsis={{ rows: 2 }}>
                        {skill.description}
                      </Paragraph>
                      <div className="skill-meta">
                        <Text type="secondary" className="skill-date">
                          更新于: {new Date(skill.updated_at).toLocaleDateString()}
                        </Text>
                        <div className="skill-platforms">
                          {parsePlatforms(skill.platforms).map((p) => (
                            <Tag key={p} className="platform-tag">{p}</Tag>
                          ))}
                        </div>
                      </div>
                    </div>
                  }
                />
              </List.Item>
            )}
          />
        )}
      </Card>

      <Modal
        title="创建新技能"
        open={modalVisible}
        onOk={handleCreate}
        onCancel={() => {
          setModalVisible(false)
          form.resetFields()
        }}
        okText="创建"
        cancelText="取消"
      >
        <Form form={form} layout="vertical">
          <Form.Item
            name="name"
            label="技能名称"
            rules={[{ required: true, message: '请输入技能名称' }]}
          >
            <Input placeholder="例如: Git 常用命令" />
          </Form.Item>
          <Form.Item
            name="description"
            label="技能描述"
            rules={[{ required: true, message: '请输入技能描述' }]}
          >
            <Input.TextArea placeholder="简要描述这个技能的用途..." rows={2} />
          </Form.Item>
          <Form.Item
            name="content"
            label="技能内容"
            rules={[{ required: true, message: '请输入技能内容' }]}
          >
            <Input.TextArea placeholder="详细的技能内容..." rows={8} />
          </Form.Item>
        </Form>
      </Modal>
    </div>
  )
}