import { Card, Row, Col, Statistic, Progress, Typography, Tag } from 'antd'
import {
  ToolOutlined,
  StarOutlined,
  MessageOutlined,
  DatabaseOutlined
} from '@ant-design/icons'
import type { Stats as StatsType } from '../store'
import './StatsPanel.css'

const { Title } = Typography

interface StatsPanelProps {
  stats: StatsType | null
  loading: boolean
}

export function StatsPanel({ stats, loading }: StatsPanelProps) {
  if (loading) {
    return <Card loading />;
  }

  if (!stats) {
    return (
      <Card>
        <Title level={5}>暂无统计数据</Title>
      </Card>
    );
  }

  const memoryTypes = Object.entries(stats.memory_stats || {});
  const totalMemory = memoryTypes.reduce((sum, [, count]) => sum + count, 0);

  return (
    <div className="stats-panel">
      <Row gutter={[16, 16]}>
        <Col xs={24} sm={12} md={6}>
          <Card className="stat-card">
            <Statistic
              title="工具调用次数"
              value={stats.tool_call_count}
              prefix={<ToolOutlined />}
              valueStyle={{ color: '#1890ff' }}
            />
          </Card>
        </Col>
        <Col xs={24} sm={12} md={6}>
          <Card className="stat-card">
            <Statistic
              title="技能数量"
              value={stats.skill_count}
              prefix={<StarOutlined />}
              valueStyle={{ color: '#faad14' }}
            />
          </Card>
        </Col>
        <Col xs={24} sm={12} md={6}>
          <Card className="stat-card">
            <Statistic
              title="对话长度"
              value={stats.conversation_len}
              prefix={<MessageOutlined />}
              valueStyle={{ color: '#52c41a' }}
            />
          </Card>
        </Col>
        <Col xs={24} sm={12} md={6}>
          <Card className="stat-card">
            <Statistic
              title="记忆总数"
              value={totalMemory}
              prefix={<DatabaseOutlined />}
              valueStyle={{ color: '#722ed1' }}
            />
          </Card>
        </Col>
      </Row>

      <Row gutter={[16, 16]} style={{ marginTop: 16 }}>
        <Col xs={24} md={12}>
          <Card title="记忆分布" className="memory-distribution">
            {memoryTypes.length > 0 ? (
              <div className="memory-bars">
                {memoryTypes.map(([type, count]) => (
                  <div key={type} className="memory-bar-item">
                    <div className="memory-bar-label">
                      <Tag color={
                        type === 'snapshot' ? 'blue' :
                        type === 'session' ? 'green' : 'purple'
                      }>
                        {type}
                      </Tag>
                      <span>{count} 条</span>
                    </div>
                    <Progress
                      percent={totalMemory > 0 ? (count / totalMemory) * 100 : 0}
                      strokeColor={
                        type === 'snapshot' ? '#1890ff' :
                        type === 'session' ? '#52c41a' : '#722ed1'
                      }
                      showInfo={false}
                    />
                  </div>
                ))}
              </div>
            ) : (
              <p>暂无记忆数据</p>
            )}
          </Card>
        </Col>
        <Col xs={24} md={12}>
          <Card title="系统状态" className="system-status">
            <div className="status-item">
              <span>Agent 版本</span>
              <Tag>v0.1.0</Tag>
            </div>
            <div className="status-item">
              <span>状态</span>
              <Tag color="success">运行中</Tag>
            </div>
            <div className="status-item">
              <span>模型</span>
              <Tag>GPT-4</Tag>
            </div>
          </Card>
        </Col>
      </Row>
    </div>
  )
}