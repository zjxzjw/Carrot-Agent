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
              prefix={<ToolOutlined style={{ color: '#007AFF' }} />}
              valueStyle={{ 
                color: '#007AFF',
                fontWeight: 700,
                letterSpacing: '-2px'
              }}
            />
          </Card>
        </Col>
        <Col xs={24} sm={12} md={6}>
          <Card className="stat-card">
            <Statistic
              title="技能数量"
              value={stats.skill_count}
              prefix={<StarOutlined style={{ color: '#FF9500' }} />}
              valueStyle={{ 
                color: '#FF9500',
                fontWeight: 700,
                letterSpacing: '-2px'
              }}
            />
          </Card>
        </Col>
        <Col xs={24} sm={12} md={6}>
          <Card className="stat-card">
            <Statistic
              title="对话长度"
              value={stats.conversation_len}
              prefix={<MessageOutlined style={{ color: '#34C759' }} />}
              valueStyle={{ 
                color: '#34C759',
                fontWeight: 700,
                letterSpacing: '-2px'
              }}
            />
          </Card>
        </Col>
        <Col xs={24} sm={12} md={6}>
          <Card className="stat-card">
            <Statistic
              title="记忆总数"
              value={totalMemory}
              prefix={<DatabaseOutlined style={{ color: '#AF52DE' }} />}
              valueStyle={{ 
                color: '#AF52DE',
                fontWeight: 700,
                letterSpacing: '-2px'
              }}
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
                      strokeColor={{
                        '0%': type === 'snapshot' ? '#007AFF' :
                              type === 'session' ? '#34C759' : '#AF52DE',
                        '100%': type === 'snapshot' ? '#5856D6' :
                                type === 'session' ? '#30D158' : '#BF5AF2'
                      }}
                      showInfo={false}
                      strokeWidth={8}
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