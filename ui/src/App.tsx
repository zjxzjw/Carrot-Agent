import { useState, useCallback, useEffect } from 'react'
import { Layout, Menu, ConfigProvider, Select, Spin, Button, Switch, Tooltip, theme as antTheme } from 'antd'
import zhCN from 'antd/locale/zh_CN'
import enUS from 'antd/locale/en_US'
import {
  MessageOutlined,
  ToolOutlined,
  DatabaseOutlined,
  HistoryOutlined,
  BarChartOutlined,
  SettingOutlined,
  GlobalOutlined,
  LogoutOutlined,
  MoonOutlined,
  SunOutlined
} from '@ant-design/icons'
import { useTranslation } from 'react-i18next'
import { useNavigate } from 'react-router-dom'
import { useDispatch, useSelector } from 'react-redux'
import { logout } from './store/authSlice'
import { toggleTheme } from './store/themeSlice'
import { RootState } from './store'
import i18n from './i18n'
import { ChatPage, SkillsPage, MemoryPage, SessionsPage, StatsPage, ConfigPage, ToolsPage } from './pages'
import logo from './assets/logo.png'

const { Header, Content, Sider } = Layout

const App: React.FC = () => {
  const [current, setCurrent] = useState('chat')
  const [loading, setLoading] = useState(false)
  const { t, i18n: i18nInstance } = useTranslation()
  const navigate = useNavigate()
  const dispatch = useDispatch()
  const currentTheme = useSelector((state: RootState) => state.theme.theme)

  useEffect(() => {
    if (currentTheme === 'dark') {
      document.body.style.backgroundColor = '#0f172a'
      document.documentElement.setAttribute('data-theme', 'dark')
    } else {
      document.body.style.backgroundColor = '#f5f7fa'
      document.documentElement.setAttribute('data-theme', 'light')
    }
  }, [currentTheme])

  useEffect(() => {
    document.documentElement.setAttribute('data-theme', currentTheme)
  }, [currentTheme])

  const handleLogout = () => {
    dispatch(logout())
    navigate('/login')
  }

  const items = [
    { key: 'chat', label: t('menu.chat'), icon: <MessageOutlined /> },
    { key: 'tools', label: t('menu.tools'), icon: <ToolOutlined /> },
    { key: 'skills', label: t('menu.skills'), icon: <DatabaseOutlined /> },
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

  const handleThemeToggle = () => {
    dispatch(toggleTheme())
  }

  const renderContent = () => {
    switch (current) {
      case 'chat':
        return <ChatPage />
      case 'tools':
        return <ToolsPage />
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

  const isDark = currentTheme === 'dark'

  return (
    <ConfigProvider 
      locale={getAntdLocale()}
      theme={{
        algorithm: isDark ? antTheme.darkAlgorithm : antTheme.defaultAlgorithm,
        token: {
          colorPrimary: '#3b82f6',
          borderRadius: 8,
        },
      }}
    >
      <Layout style={{ minHeight: '100vh' }}>
        <Sider 
          theme={isDark ? 'dark' : 'light'} 
          width={220} 
          style={{ 
            borderRight: isDark ? '1px solid #1e293b' : '1px solid #e2e8f0',
            boxShadow: '2px 0 12px rgba(0,0,0,0.08)',
          }}
        >
          <div 
            style={{ 
              height: 72, 
              display: 'flex', 
              alignItems: 'center', 
              justifyContent: 'center',
              borderBottom: isDark ? '1px solid #1e293b' : '1px solid #e2e8f0',
              background: isDark ? 'linear-gradient(135deg, #1e293b 0%, #0f172a 100%)' : 'linear-gradient(135deg, #f8fafc 0%, #e2e8f0 100%)',
            }}
          >
            <img 
              src={logo} 
              alt="Carrot Agent Logo" 
              style={{ width: 36, height: 36, marginRight: 12, borderRadius: 10 }}
            />
            <div>
              <div style={{ fontSize: 18, fontWeight: '700', color: isDark ? '#f1f5f9' : '#1e293b' }}>Carrot Agent</div>
              <div style={{ fontSize: 12, color: isDark ? '#94a3b8' : '#64748b' }}>{t('app.subtitle')}</div>
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
        <Layout style={{ background: isDark ? '#0f172a' : '#f5f7fa' }}>
          <Header 
            style={{ 
              padding: '0 24px', 
              background: isDark ? '#1e293b' : '#fff',
              borderBottom: isDark ? '1px solid #334155' : '1px solid #e2e8f0',
              display: 'flex',
              alignItems: 'center',
              justifyContent: 'space-between',
              height: 64,
            }}
          >
            <span style={{ fontSize: 15, color: isDark ? '#e2e8f0' : '#475569', fontWeight: 500 }}>
              {t(`menu.${current}`)}
            </span>
            <div style={{ display: 'flex', alignItems: 'center', gap: 20 }}>
              <div style={{ display: 'flex', alignItems: 'center', gap: 8 }}>
                <Tooltip title={isDark ? '切换到亮色模式' : '切换到暗色模式'}>
                  <Switch
                    checked={isDark}
                    onChange={handleThemeToggle}
                    checkedChildren={<MoonOutlined />}
                    unCheckedChildren={<SunOutlined />}
                  />
                </Tooltip>
              </div>
              <div style={{ display: 'flex', alignItems: 'center', gap: 8 }}>
                <GlobalOutlined style={{ color: isDark ? '#94a3b8' : '#64748b' }} />
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
              </div>
              {loading && <Spin size="small" />}
              <Button 
                type="text" 
                icon={<LogoutOutlined />} 
                onClick={handleLogout}
                style={{ color: isDark ? '#94a3b8' : '#64748b' }}
              >
                {t('header.logout')}
              </Button>
            </div>
          </Header>
          <Content 
            style={{ 
              margin: 24, 
              padding: 24, 
              background: isDark ? '#1e293b' : '#fff', 
              minHeight: 280,
              borderRadius: 16,
              boxShadow: isDark 
                ? '0 4px 20px rgba(0, 0, 0, 0.3)' 
                : '0 4px 20px rgba(0, 0, 0, 0.06)',
            }}
          >
            <div className="fade-in-up">
              {renderContent()}
            </div>
          </Content>
        </Layout>
      </Layout>
    </ConfigProvider>
  )
}

export default App