import PropTypes from 'prop-types'

import AuthenticationForm from './AuthenticationForm'
import AuthenticationHeader from './AuthenticationHeader'
import AuthenticationWrapper from './AuthenticationWrapper'

import { Form, Input, Result, Typography } from 'antd'
const { Paragraph } = Typography

const SignUp = (props) => {
  const { hasSuccess, loading, onSubmit } = props

  return (
    <AuthenticationWrapper>
      <AuthenticationHeader>Sign up</AuthenticationHeader>

      {hasSuccess && (
        <div role="alert">
          <Result
            status="success"
            subTitle={
              <Paragraph>You will be redirected to the login screen</Paragraph>
            }
            title="Sign up successful!"
          />
        </div>
      )}

      {!hasSuccess && (
        <AuthenticationForm
          alternativeActionLabel="Do you want to log in instead?"
          alternativeActionLink="/login"
          loading={loading}
          onSubmit={onSubmit}
          showForgotPassword={false}
          submitButtonLabel="Sign up"
          title="Sign up"
        >
          <Form.Item
            label="Given Name"
            name="givenName"
            rules={[{ required: true, message: 'Given name is required' }]}
          >
            <Input placeholder="Fill in your first name" />
          </Form.Item>

          <Form.Item
            label="Surname"
            name="surname"
            rules={[{ required: true, message: 'Surname is required' }]}
          >
            <Input placeholder="Fill in your last name" />
          </Form.Item>

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
            <Input placeholder="Fill in your email" type="email" />
          </Form.Item>

          <Form.Item
            label="Password"
            name="password"
            rules={[
              { required: true, message: 'Password is required' },
              () => ({
                validator(_, value) {
                  if (value && value.length >= 8) {
                    return Promise.resolve()
                  }

                  return Promise.reject(
                    new Error(
                      'Password should not be shorter than 8 characters'
                    )
                  )
                },
              }),
            ]}
          >
            <Input placeholder="Fill in your password" type="password" />
          </Form.Item>

          <Form.Item
            dependencies={['password']}
            label="Confirm Password"
            name="confirmPassword"
            rules={[
              {
                required: true,
                message: 'Please confirm your password!',
              },
              ({ getFieldValue }) => ({
                validator(_, value) {
                  if (!value || getFieldValue('password') === value) {
                    return Promise.resolve()
                  }

                  return Promise.reject(
                    new Error(
                      'The two passwords that you entered do not match!'
                    )
                  )
                },
              }),
            ]}
          >
            <Input placeholder="Fill in your password again" type="password" />
          </Form.Item>
        </AuthenticationForm>
      )}
    </AuthenticationWrapper>
  )
}

SignUp.propTypes = {
  onSubmit: PropTypes.func.isRequired,
  hasSuccess: PropTypes.bool,
  loading: PropTypes.bool,
}

export default SignUp
