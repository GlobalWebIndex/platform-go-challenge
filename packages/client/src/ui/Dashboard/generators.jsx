import { Avatar, Tag, Typography, Divider, Space } from 'antd'
const { Text } = Typography
import {
  AreaChartOutlined,
  UsergroupAddOutlined,
  ExclamationCircleOutlined,
} from '@ant-design/icons'

const generateItemAvatar = (type) => {
  switch (type) {
    case 'charts':
      return <Avatar size="large" icon={<AreaChartOutlined />} />
    case 'insights':
      return <Avatar size="large" icon={<ExclamationCircleOutlined />} />
    case 'audiences':
      return <Avatar size="large" icon={<UsergroupAddOutlined />} />
  }
}

const generateItemTitle = (type) => {
  switch (type) {
    case 'charts':
      return (
        <>
          <Tag color="magenta">Chart</Tag>
          <span>Awesome Chart</span>
        </>
      )
    case 'insights':
      return (
        <>
          <Tag color="purple">Insight</Tag>
          <span>Awesome Insight</span>
        </>
      )
    case 'audiences':
      return (
        <>
          <Tag color="green">Audience</Tag>
          <span>Awesome Audience</span>
        </>
      )
  }
}

const generateItemContent = (type, item) => {
  switch (type) {
    case 'charts':
      return (
        <Space>
          <Text strong>Title:</Text>
          <Text>{item.Chart.Title}</Text>
          <Divider type="vertical" />
          <Text strong>HorizontalAxisLabel:</Text>
          <Text>{item.Chart.HorizontalAxisLabel}</Text>
          <Divider type="vertical" />
          <Text strong>VerticalAxisLabel:</Text>
          <Text>{item.Chart.VerticalAxisLabel}</Text>
          <Divider type="vertical" />
          <Text strong>Data:</Text>
          <Text>{item.Chart.Data}</Text>
        </Space>
      )
    case 'insights':
      return (
        <Space>
          <Text strong>Text:</Text>
          <Text>{item.Insight.Text}</Text>
        </Space>
      )
    case 'audiences':
      return (
        <Space>
          <Space>
            <Text strong>Gender:</Text>
            <Text>{item.Audience.Gender}</Text>
            <Divider type="vertical" />
            <Text strong>Age Group:</Text>
            <Text>{item.Audience.AgeGroup}</Text>
          </Space>
          <Space>
            <Divider type="vertical" />
            <Text strong>Country of Birth:</Text>
            <Text>{item.Audience.CountryOfBirth}</Text>
            <Divider type="vertical" />
            <Text strong>Hours spent on social media:</Text>
            <Text>{item.Audience.DailyHoursOnSocialMedia}</Text>
            <Divider type="vertical" />
            <Text strong>Number of purchases last month:</Text>
            <Text>{item.Audience.LastMonthNumberOfPurchases}</Text>
          </Space>
        </Space>
      )
  }
}

export { generateItemAvatar, generateItemContent, generateItemTitle }
