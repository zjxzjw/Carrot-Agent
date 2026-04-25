import { useEffect, useCallback } from 'react'
import { useDispatch, useSelector } from 'react-redux'
import { Card, Row, Col, Statistic, Button, message } from 'antd'
import { ReloadOutlined, AimOutlined, BulbOutlined, HistoryOutlined, DatabaseOutlined } from '@ant-design/icons'
import { useTranslation } from 'react-i18next'
import { AppDispatch, RootState } from '../store'
import { fetchStats } from '../store/statsSlice'
import { Loading } from '../components'

const StatsPage: React.FC = () => {
  const dispatch = useDispatch<AppDispatch>()
  const { tool_call_count, skill_count, memory_stats, conversation_len, loading, error } = useSelector(
    (state: RootState) => state.stats
  )
  const { t } = useTranslation()

  useEffect(() => {
    dispatch(fetchStats())
  }, [dispatch])

  useEffect(() => {
    if (error) {
      message.error(error)
    }
  }, [error])

  const handleRefresh = useCallback(() => {
    dispatch(fetchStats())
  }, [dispatch])

  const totalMemories = memory_stats.snapshot + memory_stats.session + memory_stats.longterm

  const statsCards = [
    {
      title: t('stats.toolCalls'),
      value: tool_call_count,
      icon: <AimOutlined style={{ fontSize: 24, color: '#1890ff' }} />,
      color: '#1890ff',
    },
    {
      title: t('stats.skillCount'),
      value: skill_count,
      icon: <BulbOutlined style={{ fontSize: 24, color: '#52c41a' }} />,
      color: '#52c41a',
    },
    {
      title: t('stats.conversationLen'),
      value: conversation_len,
      icon: <HistoryOutlined style={{ fontSize: 24, color: '#faad14' }} />,
      color: '#faad14',
    },
    {
      title: t('stats.totalMemories'),
      value: totalMemories,
      icon: <DatabaseOutlined style={{ fontSize: 24, color: '#722ed1' }} />,
      color: '#722ed1',
    },
  ]

  const memoryCards = [
    { title: t('stats.snapshot'), value: memory_stats.snapshot, color: '#1890ff' },
    { title: t('stats.session'), value: memory_stats.session, color: '#52c41a' },
    { title: t('stats.longterm'), value: memory_stats.longterm, color: '#faad14' },
  ]

  if (loading && !tool_call_count) {
    return <Loading size="large" />
  }

  return (
    <div>
      <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: 24 }}>
        <h2 style={{ margin: 0 }}>{t('stats.title')}</h2>
        <Button 
          icon={<ReloadOutlined />} 
          onClick={handleRefresh} 
          loading={loading}
        >
          {t('stats.refresh')}
        </Button>
      </div>

      <Row gutter={[16, 16]}>
        {statsCards.map((stat, index) => (
          <Col xs={24} sm={12} md={6} key={index}>
            <Card bordered={false} hoverable>
              <div style={{ display: 'flex', alignItems: 'center', gap: 16 }}>
                <div
                  style={{
                    width: 56,
                    height: 56,
                    borderRadius: 8,
                    backgroundColor: `${stat.color}15`,
                    display: 'flex',
                    alignItems: 'center',
                    justifyContent: 'center',
                  }}
                >
                  {stat.icon}
                </div>
                <div>
                  <Statistic title={<span style={{ fontSize: 14, color: '#666' }}>{stat.title}</span>} value={stat.value} />
                </div>
              </div>
            </Card>
          </Col>
        ))}
      </Row>

      <h3 style={{ marginTop: 32, marginBottom: 16 }}>{t('stats.memoryDistribution')}</h3>
      <Row gutter={[16, 16]}>
        {memoryCards.map((mem, index) => (
          <Col xs={24} sm={8} key={index}>
            <Card bordered={false} hoverable>
              <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
                <div>
                  <div style={{ fontSize: 14, color: '#666', marginBottom: 8 }}>{mem.title}</div>
                  <div style={{ fontSize: 28, fontWeight: 'bold', color: mem.color }}>{mem.value}</div>
                </div>
                <div
                  style={{
                    width: 64,
                    height: 64,
                    borderRadius: '50%',
                    border: `4px solid ${mem.color}`,
                    display: 'flex',
                    alignItems: 'center',
                    justifyContent: 'center',
                    fontSize: 20,
                    fontWeight: 'bold',
                    color: mem.color,
                  }}
                >
                  {totalMemories > 0 ? Math.round((mem.value / totalMemories) * 100) : 0}%
                </div>
              </div>
            </Card>
          </Col>
        ))}
      </Row>
    </div>
  )
}

export default StatsPage