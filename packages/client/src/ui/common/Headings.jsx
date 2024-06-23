import PropTypes from 'prop-types'

import { Typography } from 'antd'

const { Title } = Typography

const Heading = (props) => {
  const { className, children, level } = props

  return (
    // eslint-disable-next-line react/jsx-props-no-spreading
    <Title className={className} level={level}>
      {children}
    </Title>
  )
}

Heading.propTypes = {
  level: PropTypes.number,
}

export const H1 = ({ children, className }) => (
  <Heading className={className} level={1}>
    {children}
  </Heading>
)

export const H2 = ({ children, className }) => (
  <Heading className={className} level={2}>
    {children}
  </Heading>
)

export const H3 = ({ children, className }) => (
  <Heading className={className} level={3}>
    {children}
  </Heading>
)

export const H4 = ({ children, className }) => (
  <Heading className={className} level={4}>
    {children}
  </Heading>
)

export const H5 = ({ children, className }) => (
  <Heading className={className} level={5}>
    {children}
  </Heading>
)

Heading.propTypes = {
  className: PropTypes.string,
  children: PropTypes.any,
}
H1.propTypes = {
  className: PropTypes.string,
  children: PropTypes.any,
}
H2.propTypes = {
  className: PropTypes.string,
  children: PropTypes.any,
}
H3.propTypes = {
  className: PropTypes.string,
  children: PropTypes.any,
}
H4.propTypes = {
  className: PropTypes.string,
  children: PropTypes.any,
}
H5.propTypes = {
  className: PropTypes.string,
  children: PropTypes.any,
}
