import { useState } from 'react'
import { useMutation } from 'react-query'
import { useNavigate } from 'react-router-dom'

import axios from 'axios'
import { App } from 'antd'

import SignUp from '../ui/Authentication/SignUp'

const signUpUser = async (formData) => {
  try {
    return axios.post('http://localhost:3000/api/v1/users', formData, {
      headers: { 'Content-Type': 'application/json' },
    })
  } catch (e) {
    throw new Error(e.message)
  }
}

const SignupPage = () => {
  const { notification } = App.useApp()
  const [loading, setLoading] = useState(false)
  const [hasSuccess, setHasSuccess] = useState(false)
  const navigate = useNavigate()

  const redirectToLogin = () => navigate('/login')

  const showNotification = (eMsg) => {
    notification.error({
      message: 'Oops, something went wrong',
      description: eMsg,
      placement: 'topRight',
    })
  }

  const mutation = useMutation({
    mutationFn: signUpUser,
    onError: (error) => {
      showNotification(error.message)
      setLoading(false)
    },
    onSuccess: async () => {
      setHasSuccess(true)
      setTimeout(redirectToLogin, 1000)
    },
    onSettled: async () => {
      setLoading(false)
    },
  })

  const onSignUp = (formData) => {
    setLoading(true)
    mutation.mutate(formData)
  }

  return (
    <SignUp hasSuccess={hasSuccess} loading={loading} onSubmit={onSignUp} />
  )
}

export default SignupPage
