import PropTypes from 'prop-types'
import { useLocation } from 'react-router-dom'

import { App, Layout, Menu, Spin, Button, Avatar } from 'antd'
const { Header, Footer, Content } = Layout
import {
  HomeOutlined,
  AreaChartOutlined,
  UsergroupAddOutlined,
  ExclamationCircleOutlined,
  StarOutlined,
  UserOutlined,
} from '@ant-design/icons'

import { useAuth } from './hooks/useAuth'
import { useCurrentUser } from './hooks/useCurrentUser'

import AppRoutes from './AppRoutes'
import LinkWithoutStyles from './ui/common/Link'

const layoutStyle = {
  overflow: 'hidden',
  width: '100%',
  height: '100%',
}

const mainContentStyle = {
  background: 'white',
  padding: '10px 48px',
  height: 'calc(100% - 64px - 70xp - 10px)',
  margin: '0 auto',
  width: '95%',
}

const generateMenuItems = (currentUser, logoutHandler) => {
  const menuItems = [
    {
      key: 'assets',
      label: <LinkWithoutStyles to="/assets">All Assets</LinkWithoutStyles>,
      icon: <HomeOutlined />,
    },
    {
      key: 'charts',
      label: <LinkWithoutStyles to="/charts">Charts</LinkWithoutStyles>,
      icon: <AreaChartOutlined />,
    },
    {
      key: 'audiences',
      label: <LinkWithoutStyles to="/audiences">Audiences</LinkWithoutStyles>,
      icon: <UsergroupAddOutlined />,
    },
    {
      key: 'insights',
      label: <LinkWithoutStyles to="/insights">Insights</LinkWithoutStyles>,
      icon: <ExclamationCircleOutlined />,
    },
    {
      key: 'favorites',
      label: <LinkWithoutStyles to="/favorites">Favorites</LinkWithoutStyles>,
      icon: <StarOutlined />,
    },
  ]

  if (currentUser) {
    menuItems.push({
      label: `${currentUser.givenName} ${currentUser.surname} `,
      key: 'userMenu',
      icon: <Avatar size={16} icon={<UserOutlined />} />,
      children: [
        {
          type: 'group',
          label: 'User Actions',
          children: [
            {
              label: (
                <Button onClick={logoutHandler} type="primary" danger>
                  Logout
                </Button>
              ),
              key: 'logout',
            },
          ],
        },
      ],
    })
  }
  return menuItems
}

export const Application = () => {
  const auth = useAuth()
  const { currentUser, loadingUser } = useCurrentUser()
  const location = useLocation()
  const { pathname } = location
  const activeItem = pathname.replace(/^\/+/g, '')

  const items = generateMenuItems(currentUser, auth.logout)

  return (
    <App style={{ height: '100vh' }}>
      {loadingUser ? (
        <Spin size="large" />
      ) : (
        <Layout style={layoutStyle}>
          {currentUser ? (
            <Header
              style={{
                display: 'flex',
                alignItems: 'center',
                marginBottom: '10px',
              }}
            >
              <Menu
                theme="dark"
                mode="horizontal"
                defaultSelectedKeys={['assets']}
                selectedKeys={[activeItem]}
                items={items}
                style={{ flex: 1, minWidth: 0 }}
              />
            </Header>
          ) : null}
          <Content style={mainContentStyle}>
            <AppRoutes />
          </Content>
          <Footer style={{ textAlign: 'center' }}>
            GWI ©{new Date().getFullYear()} Created by Alexandros Georgantas
          </Footer>
        </Layout>
      )}
    </App>
  )
}

Application.propTypes = {
  children: PropTypes.any,
}
