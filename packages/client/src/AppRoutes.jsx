import { Route, Routes, Navigate } from 'react-router-dom'

import LoginPage from './pages/Login.page'
import SignUpPage from './pages/SignUp.page'
import AssetsPage from './pages/Assets.page'
import AudiencesPage from './pages/Audience.page'
import InsightsPage from './pages/Insights.page'
import FavoritesPage from './pages/Favorites.page'
import ChartsPage from './pages/Charts.page'
import PrivateRoute from './helpers/PrivateRoute'

import { useAuth } from './hooks/useAuth'

const AppRoutes = () => {
  const auth = useAuth()
  return (
    <Routes>
      <Route exact path="/" element={<Navigate to="/assets" replace />} />
      <Route
        element={!auth.token ? <SignUpPage /> : <Navigate replace to={'/'} />}
        exact
        path="/signup"
      />
      <Route
        element={!auth.token ? <LoginPage /> : <Navigate replace to={'/'} />}
        exact
        path="/login"
      />
      <Route element={<PrivateRoute />}>
        <Route path="/assets" element={<AssetsPage />} />
      </Route>
      <Route element={<PrivateRoute />}>
        <Route path="/charts" element={<ChartsPage />} />
      </Route>
      <Route element={<PrivateRoute />}>
        <Route path="/insights" element={<InsightsPage />} />
      </Route>
      <Route element={<PrivateRoute />}>
        <Route path="/audiences" element={<AudiencesPage />} />
      </Route>
      <Route element={<PrivateRoute />}>
        <Route path="/favorites" element={<FavoritesPage />} />
      </Route>
    </Routes>
  )
}

export default AppRoutes
