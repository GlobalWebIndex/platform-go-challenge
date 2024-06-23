import styled from 'styled-components'
import List from './List'

const StyledList = styled(List)`
  height: 100%;
`
const Dashboard = ({
  items,
  isLoading,
  pageSize,
  paginationHandler,
  shouldDisplayRemove,
  mainActionHandler,
  secondaryActionHandler,
}) => (
  <StyledList
    items={items}
    isLoading={isLoading}
    paginationHandler={paginationHandler}
    pageSize={pageSize}
    shouldDisplayRemove={shouldDisplayRemove}
    mainActionHandler={mainActionHandler}
    secondaryActionHandler={secondaryActionHandler}
  />
)

export default Dashboard
