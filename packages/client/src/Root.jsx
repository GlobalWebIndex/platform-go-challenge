import { QueryClient, QueryClientProvider } from 'react-query'
import { BrowserRouter } from 'react-router-dom'
import { Normalize } from 'styled-normalize'

import AuthProvider from './providers/AuthProvider'
import CurrentUserProvider from './providers/CurrentUserProvider'

import { Application } from './App'

const queryClient = new QueryClient()

const Root = () => {
  return (
    <QueryClientProvider client={queryClient}>
      <BrowserRouter>
        <Normalize />
        <AuthProvider>
          <CurrentUserProvider>
            <Application />
          </CurrentUserProvider>
        </AuthProvider>
      </BrowserRouter>
    </QueryClientProvider>
  )
}

export default Root
