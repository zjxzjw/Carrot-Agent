import { useState, useCallback } from 'react'
import { Layout, Menu, ConfigProvider, Select, Spin } from 'antd'
import zhCN from 'antd/locale/zh_CN'
import enUS from 'antd/locale/en_US'
import {
  MessageOutlined,
  ToolOutlined,
  DatabaseOutlined,
  HistoryOutlined,
  BarChartOutlined,
  SettingOutlined,
  GlobalOutlined
} from '@ant-design/icons'
import { useTranslation } from 'react-i18next'
import i18n from './i18n'
import { ChatPage, SkillsPage, MemoryPage, SessionsPage, StatsPage, ConfigPage } from './pages'
import logo from './assets/logo.png'

const { Header, Content, Sider } = Layout

const App: React.FC = () => {
  const [current, setCurrent] = useState('chat')
  const [loading, setLoading] = useState(false)
  const { t, i18n: i18nInstance } = useTranslation()

  const items = [
    { key: 'chat', label: t('menu.chat'), icon: <MessageOutlined /> },
    { key: 'skills', label: t('menu.skills'), icon: <ToolOutlined /> },
    { key: 'memory', label: t('menu.memory'), icon: <DatabaseOutlined /> },
    { key: 'sessions', label: t('menu.sessions'), icon: <HistoryOutlined /> },
    { key: 'stats', label: t('menu.stats'), icon: <BarChartOutlined /> },
    { key: 'config', label: t('menu.config'), icon: <SettingOutlined /> },
  ]

  const handleMenuClick = useCallback((e: { key: string }) => {
    setCurrent(e.key)
  }, [])

  const handleLanguageChange = useCallback((value: string) => {
    setLoading(true)
    i18nInstance.changeLanguage(value)
      .finally(() => setLoading(false))
  }, [i18nInstance])

  const renderContent = () => {
    switch (current) {
      case 'chat':
        return <ChatPage />
      case 'skills':
        return <SkillsPage />
      case 'memory':
        return <MemoryPage />
      case 'sessions':
        return <SessionsPage />
      case 'stats':
        return <StatsPage />
      case 'config':
        return <ConfigPage />
      default:
        return <ChatPage />
    }
  }

  const getAntdLocale = () => {
    return i18n.language === 'zh-CN' ? zhCN : enUS
  }

  return (
    <ConfigProvider locale={getAntdLocale()}>
      <Layout style={{ minHeight: '100vh' }}>
        <Sider 
          theme="light" 
          width={220} 
          style={{ 
            borderRight: '1px solid #f0f0f0',
            boxShadow: '2px 0 8px rgba(0,0,0,0.05)',
          }}
        >
          <div 
            style={{ 
              height: 72, 
              display: 'flex', 
              alignItems: 'center', 
              justifyContent: 'center',
              borderBottom: '1px solid #f0f0f0',
            }}
          >
            <img 
              src={logo} 
              alt="Carrot Agent Logo" 
              style={{ width: 32, height: 32, marginRight: 12 }}
            />
            <div>
              <div style={{ fontSize: 18, fontWeight: 'bold', color: '#333' }}>Carrot Agent</div>
              <div style={{ fontSize: 12, color: '#999' }}>{t('app.subtitle')}</div>
            </div>
          </div>
          <Menu
            mode="inline"
            selectedKeys={[current]}
            onClick={handleMenuClick}
            items={items}
            style={{ height: 'calc(100% - 72px)', borderRight: 0 }}
          />
        </Sider>
        <Layout>
          <Header 
            style={{ 
              padding: '0 24px', 
              background: '#fff',
              borderBottom: '1px solid #f0f0f0',
              display: 'flex',
              alignItems: 'center',
              justifyContent: 'space-between',
            }}
          >
            <span style={{ fontSize: 14, color: '#666' }}>
              {t(`menu.${current}`)}
            </span>
            <div style={{ display: 'flex', alignItems: 'center', gap: 8 }}>
              <GlobalOutlined style={{ color: '#666' }} />
              <Select
                value={i18n.language}
                onChange={handleLanguageChange}
                options={[
                  { value: 'zh-CN', label: '中文' },
                  { value: 'en-US', label: 'English' },
                ]}
                style={{ width: 120 }}
                disabled={loading}
              />
              {loading && <Spin size="small" />}
            </div>
          </Header>
          <Content 
            style={{ 
              margin: 24, 
              padding: 24, 
              background: '#fff', 
              minHeight: 280,
              borderRadius: 8,
              boxShadow: '0 1px 3px rgba(0,0,0,0.1)',
            }}
          >
            {renderContent()}
          </Content>
        </Layout>
      </Layout>
    </ConfigProvider>
  )
}

export default App