import { useState } from 'react'

import { useMutation, useQueryClient } from 'react-query'
import axios from 'axios'
import { App } from 'antd'

import Login from '../ui/Authentication/Login'

import { useAuth } from '../hooks/useAuth'

const logInUser = async (formData) => {
  try {
    return axios.post('http://localhost:3000/api/v1/tokens', formData, {
      headers: { 'Content-Type': 'application/json' },
    })
  } catch (e) {
    throw new Error(e.message)
  }
}

const LoginPage = () => {
  const { notification } = App.useApp()
  const queryClient = useQueryClient()
  const auth = useAuth()
  const [loading, setLoading] = useState(false)

  const showNotification = (eMsg) => {
    notification.error({
      message: 'Oops, something went wrong',
      description: eMsg,
      placement: 'topRight',
    })
  }

  const mutation = useMutation({
    mutationFn: logInUser,
    onError: (error) => {
      showNotification(error.message)
      queryClient.invalidateQueries('currentUser')
      setLoading(false)
    },
    onSuccess: async ({ data }) => {
      const { token } = data
      queryClient.invalidateQueries('currentUser')
      auth.login(token)
    },
    onSettled: async () => {
      setLoading(false)
    },
  })

  const onLogin = (formData) => {
    setLoading(true)
    mutation.mutate(formData)
  }

  return <Login loading={loading} onSubmit={onLogin} />
}

export default LoginPage
