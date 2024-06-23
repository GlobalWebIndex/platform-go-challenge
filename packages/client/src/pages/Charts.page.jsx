import { useState } from 'react'

import { useQuery, useMutation, useQueryClient } from 'react-query'
import { useLocation } from 'react-router-dom'
import { App } from 'antd'

import Dashboard from '../ui/Dashboard/Dashboard'
import { useCurrentUser } from '../hooks/useCurrentUser'

import {
  getAssets,
  addToFavorites,
  updateAssetDescription,
} from '../utils/network'

const ChartsPage = () => {
  const { notification } = App.useApp()
  const location = useLocation()
  const { pathname } = location
  const assetType = pathname.replace(/^\/+/g, '')
  const { currentUser } = useCurrentUser()
  const queryClient = useQueryClient()

  const [paginationSpecs, setPaginationSpecs] = useState({
    page: 1,
    pageSize: 10,
  })

  const showErrorNotification = (eMsg) => {
    notification.error({
      message: 'Oops, something went wrong',
      description: eMsg,
      placement: 'topRight',
    })
  }

  const showSuccessNotification = (message) => {
    notification.success({
      message,
      placement: 'topRight',
    })
  }

  const showAlreadyNotification = () => {
    notification.warning({
      message: 'Already there',
      placement: 'topRight',
    })
  }

  const addToFavoritesMutation = useMutation({
    mutationFn: addToFavorites,
    onError: ({ message }) => {
      if (message !== 'already in favorites of user') {
        return showErrorNotification(message)
      }
      showAlreadyNotification()
    },
    onSuccess: async () => {
      showSuccessNotification('Added to favorites')
    },
  })

  const { isLoading, isError, data, error } = useQuery({
    queryKey: [assetType, paginationSpecs.page, paginationSpecs.pageSize],
    queryFn: () =>
      getAssets(paginationSpecs.page, paginationSpecs.pageSize, assetType),
  })

  const updateAssetDescriptionMutation = useMutation({
    mutationFn: updateAssetDescription,
    onError: ({ message }) => {
      showErrorNotification(message)
    },
    onSuccess: async () => {
      showSuccessNotification('Description updated')
      queryClient.invalidateQueries({
        queryKey: [assetType, paginationSpecs.page, paginationSpecs.pageSize],
      })
    },
  })

  const onMainAction = (assetId) => {
    const payload = { assetId, userId: currentUser.id }
    return addToFavoritesMutation.mutateAsync(payload)
  }

  const onSecondaryAction = (assetId, description) => {
    const payload = { assetId, description }
    return updateAssetDescriptionMutation.mutateAsync(payload)
  }

  const onPaginationChange = (page, pageSize) => {
    setPaginationSpecs({ page, pageSize })
  }

  if (isError) {
    showErrorNotification(error)
  }

  return (
    <Dashboard
      items={data?.assets}
      isLoading={isLoading}
      paginationHandler={onPaginationChange}
      pageSize={paginationSpecs.pageSize}
      shouldDisplayRemove={false}
      mainActionHandler={onMainAction}
      secondaryActionHandler={onSecondaryAction}
    />
  )
}

export default ChartsPage
