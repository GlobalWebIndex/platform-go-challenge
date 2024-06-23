import PropTypes from 'prop-types'
import { LockOutlined, UserOutlined } from '@ant-design/icons'
// import styled from 'styled-components'

import { Form, Input } from 'antd'
import AuthenticationForm from './AuthenticationForm'
import AuthenticationHeader from './AuthenticationHeader'
import AuthenticationWrapper from './AuthenticationWrapper'

const Login = (props) => {
  const { loading, onSubmit } = props

  return (
    <AuthenticationWrapper>
      <AuthenticationHeader>Login</AuthenticationHeader>

      <AuthenticationForm
        alternativeActionLabel="Do you want to sign up instead?"
        alternativeActionLink="/signup"
        loading={loading}
        onSubmit={onSubmit}
        showForgotPassword
        submitButtonLabel="Log in"
        title="Login"
      >
        <Form.Item
          label="Email"
          name="email"
          rules={[
            {
              required: true,
              message: 'Email is required',
            },
            {
              type: 'email',
              message: 'This is not a valid email address',
            },
          ]}
        >
          <Input
            autoComplete="on"
            placeholder="Please enter your email"
            prefix={<UserOutlined className="site-form-item-icon" />}
            type="email"
          />
        </Form.Item>

        <Form.Item
          label="Password"
          name="password"
          rules={[{ required: true, message: 'Password is required' }]}
        >
          <Input
            autoComplete="on"
            placeholder="Please enter your password"
            prefix={<LockOutlined className="site-form-item-icon" />}
            type="password"
          />
        </Form.Item>
      </AuthenticationForm>
    </AuthenticationWrapper>
  )
}

Login.propTypes = {
  loading: PropTypes.bool,
  onSubmit: PropTypes.func.isRequired,
}

export default Login
