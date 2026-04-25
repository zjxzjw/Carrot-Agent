import { useEffect, useRef, useState, useReducer } from 'react'
import { Layout, Menu, Button, Input, Tabs, Typography, Alert, Tooltip, Space } from 'antd'
import { MessageOutlined, SettingOutlined, PlusOutlined, SendOutlined, SyncOutlined } from '@ant-design/icons'
import {
  initialState, 
  agentReducer,
  checkConnection,
  sendMessage,
  fetchSkills,
  fetchMemories,
  fetchStats,
  createSkill,
  addMemory,
  selectConversation,
  addConversation,
  clearError,
  getSortedMemories,
  getSortedSkills
} from './store'
import { ChatList, ConnectionStatus } from './components/ChatMessage'
import { SkillList } from './components/SkillList'
import { MemoryList } from './components/MemoryList'
import { StatsPanel } from './components/StatsPanel'
import './App.css'

const { Header, Content, Sider } = Layout
const { Title } = Typography
const { Search, TextArea } = Input
const { TabPane } = Tabs

function App() {
  const [state, dispatch] = useReducer(agentReducer, initialState)
  const [messageInput, setMessageInput] = useState('')
  const chatListRef = useRef<HTMLDivElement>(null)

  useEffect(() => {
    checkConnection(dispatch)
    const interval = setInterval(() => checkConnection(dispatch), 30000)
    return () => clearInterval(interval)
  }, [])

  useEffect(() => {
    if (chatListRef.current) {
      chatListRef.current.scrollTop = chatListRef.current.scrollHeight
    }
  }, [state.messages])

  const handleSendMessage = () => {
    if (messageInput.trim()) {
      sendMessage(dispatch, messageInput, state.currentConversationId)
      setMessageInput('')
    }
  }

  const handleKeyPress = (e: React.KeyboardEvent) => {
    if (e.key === 'Enter' && !e.shiftKey) {
      e.preventDefault()
      handleSendMessage()
    }
  }

  const handleTabChange = async (key: string) => {
    switch (key) {
      case 'skills':
        await fetchSkills(dispatch)
        break
      case 'memory':
        await fetchMemories(dispatch)
        break
      case 'stats':
        await fetchStats(dispatch)
        break
    }
  }

  const sortedSkills = getSortedSkills(state)
  const sortedMemories = getSortedMemories(state)

  return (
    <Layout style={{ minHeight: '100vh' }} className="app-layout">
      <Header className="app-header">
        <div className="header-left">
            <Title level={3} style={{ margin: 0, fontWeight: 700 }}>Carrot Agent</Title>
            <ConnectionStatus connected={state.connected} />
          </div>
        <Space size="middle">
          <Tooltip title="刷新连接状态">
            <Button 
              icon={<SyncOutlined spin={!state.connected} />} 
              onClick={() => checkConnection(dispatch)}
              shape="circle"
            />
          </Tooltip>
          <Button 
            type="primary" 
            icon={<PlusOutlined />} 
            onClick={() => addConversation(dispatch)}
            size="large"
          >
            新建对话
          </Button>
        </Space>
      </Header>

      <Layout className="app-body">
        <Sider width={280} className="app-sider" collapsible>
          <div className="sider-search">
            <Search placeholder="搜索对话..." allowClear />
          </div>
          <Menu
            mode="inline"
            selectedKeys={[state.currentConversationId]}
            onClick={({ key }) => selectConversation(dispatch, key)}
            className="conversation-menu"
            items={state.conversations.map((conv) => ({
              key: conv.id,
              icon: <MessageOutlined />,
              label: conv.title
            }))}
          />
        </Sider>

        <Content className="app-content">
          <Tabs
            defaultActiveKey="chat"
            onChange={handleTabChange}
            className="main-tabs"
            tabBarExtraContent={
              <Button type="text" icon={<SettingOutlined />}>设置</Button>
            }
          >
            <TabPane tab="对话" key="chat">
              <div className="chat-container">
                <div className="chat-messages" ref={chatListRef}>
                  <ChatList messages={state.messages} />
                </div>

                {state.error && (
                  <Alert
                    message={state.error}
                    type="error"
                    showIcon
                    closable
                    onClose={() => clearError(dispatch)}
                    style={{ marginBottom: 16 }}
                  />
                )}

                <div className="chat-input-area">
                  <TextArea
                    value={messageInput}
                    onChange={(e) => setMessageInput(e.target.value)}
                    onKeyDown={handleKeyPress}
                    placeholder="输入消息... (Enter 发送, Shift+Enter 换行)"
                    autoSize={{ minRows: 1, maxRows: 4 }}
                    disabled={state.loading}
                    bordered={false}
                  />
                  <Button
                    type="primary"
                    icon={<SendOutlined />}
                    onClick={handleSendMessage}
                    loading={state.loading}
                    disabled={!messageInput?.trim()}
                    size="large"
                    shape="round"
                  >
                    发送
                  </Button>
                </div>
              </div>
            </TabPane>

            <TabPane tab="技能" key="skills">
              <SkillList
                skills={sortedSkills}
                loading={state.loading}
                onCreate={(name, description, content) => createSkill(dispatch, name, description, content)}
              />
            </TabPane>

            <TabPane tab="记忆" key="memory">
              <MemoryList
                memories={sortedMemories}
                loading={state.loading}
                onAdd={(type, content, metadata) => addMemory(dispatch, type, content, metadata)}
              />
            </TabPane>

            <TabPane tab="统计" key="stats">
              <StatsPanel stats={state.stats} loading={state.loading} />
            </TabPane>
          </Tabs>
        </Content>
      </Layout>
    </Layout>
  )
}

export default App