import PropTypes from 'prop-types'
import styled, { createGlobalStyle } from 'styled-components'

// TO DO -- Remove div > div when you clean up client from pubsweet
const GlobalStyle = createGlobalStyle`
  html {
    height: 100%;
  }

  body {
    height: 100vh;
    overflow: hidden;
  }

  #root {
    height: 100%;
  }
`

const PageContainer = styled.div`
  flex: auto;

  height: 100%;
  overflow-y: auto;
  padding: 8px 8px 50px 8px;
  width: 1200px;
`

// TO DO -- move global style to root when you export that from this client
const Page = ({ children, className }) => (
  <>
    <GlobalStyle />
    <PageContainer className={className}>{children}</PageContainer>
  </>
)

Page.propTypes = {
  className: PropTypes.string,
  children: PropTypes.node,
}

export default Page
