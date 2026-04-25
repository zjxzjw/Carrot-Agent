import { useState, useEffect } from 'react'
import { Card, Table, Tag, Typography, Spin, Alert, Empty } from 'antd'
import { useTranslation } from 'react-i18next'
import { ToolOutlined } from '@ant-design/icons'
import axios from 'axios'

const { Title, Text } = Typography

interface Tool {
  name: string
  description: string
  parameters: any
  toolset: string
  version: string
  enabled: boolean
}

interface ToolsPageProps {}

const ToolsPage: React.FC<ToolsPageProps> = () => {
  const [tools, setTools] = useState<Tool[]>([])
  const [toolsByToolset, setToolsByToolset] = useState<Record<string, Tool[]>>({})
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)
  const { t } = useTranslation()

  useEffect(() => {
    const fetchTools = async () => {
      try {
        setLoading(true)
        setError(null)
        const response = await axios.get('/api/tools', {
          headers: {
            Authorization: localStorage.getItem('sessionId') || ''
          }
        })
        setTools(response.data.tools || [])

        const grouped: Record<string, Tool[]> = {}
        const toolList = response.data.tools || []

        // Group tools by toolset
        for (const tool of toolList) {
          const ts = tool.toolset || 'default'
          if (!grouped[ts]) {
            grouped[ts] = []
          }
          grouped[ts].push(tool)
        }
        setToolsByToolset(grouped)
      } catch (err) {
        setError(t('tools.error.fetch'))
        console.error('Failed to fetch tools:', err)
      } finally {
        setLoading(false)
      }
    }

    fetchTools()
  }, [t])

  const columns = [
    {
      title: t('tools.table.name'),
      dataIndex: 'name',
      key: 'name',
      render: (text: string) => <Text strong>{text}</Text>
    },
    {
      title: t('tools.table.description'),
      dataIndex: 'description',
      key: 'description',
    },
    {
      title: t('tools.table.toolset'),
      dataIndex: 'toolset',
      key: 'toolset',
      render: (text: string) => (
        <Tag color={text === 'default' ? 'blue' : 'default'}>
          {text}
        </Tag>
      )
    },
    {
      title: t('tools.table.parameters'),
      dataIndex: 'parameters',
      key: 'parameters',
      render: (parameters: any) => {
        if (!parameters || !parameters.properties) return null
        const props = Object.keys(parameters.properties)
        return props.length > 0 ? (
          <div>
            {props.map((prop, index) => (
              <div key={index} style={{ fontSize: '12px' }}>
                <Text type="secondary">{prop}</Text>
              </div>
            ))}
          </div>
        ) : null
      }
    }
  ]

  if (loading) {
    return (
      <div style={{ display: 'flex', justifyContent: 'center', alignItems: 'center', height: '400px' }}>
        <Spin size="large" />
      </div>
    )
  }

  if (error) {
    return (
      <Alert
        message={t('tools.error.title')}
        description={error}
        type="error"
        showIcon
      />
    )
  }

  return (
    <div>
      <div style={{ marginBottom: 24 }}>
        <Title level={4}>{t('tools.title')}</Title>
        <Text type="secondary">{t('tools.description')}</Text>
      </div>

      <Card
        title={
          <div style={{ display: 'flex', alignItems: 'center', gap: 8 }}>
            <ToolOutlined />
            {t('tools.all_tools')}
          </div>
        }
        style={{ marginBottom: 24 }}
      >
        {tools.length > 0 ? (
          <Table
            columns={columns}
            dataSource={tools}
            rowKey="name"
            pagination={{ pageSize: 10 }}
          />
        ) : (
          <Empty description={t('tools.empty')} />
        )}
      </Card>

      {Object.entries(toolsByToolset).map(([toolset, toolList]) => (
        <Card
          key={toolset}
          title={
            <div style={{ display: 'flex', alignItems: 'center', gap: 8 }}>
              <Tag color={toolset === 'default' ? 'blue' : 'default'}>
                {toolset}
              </Tag>
              <Text type="secondary">({toolList.length} {t('tools.tools')})</Text>
            </div>
          }
          style={{ marginBottom: 16 }}
        >
          <Table
            columns={columns.slice(0, 3)} // Exclude parameters column for toolset view
            dataSource={toolList}
            rowKey="name"
            pagination={false}
          />
        </Card>
      ))}
    </div>
  )
}

export default ToolsPage
