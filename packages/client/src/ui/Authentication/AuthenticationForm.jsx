import PropTypes from 'prop-types'
import styled from 'styled-components'
import { Button, Form } from 'antd'
import Link from '../common/Link'
// import { Form, Button, Link } from "../common";

const Wrapper = styled.div``

const SubmitButton = styled(Button)`
  width: 100%;
`

const Footer = styled.div`
  display: flex;
  justify-content: 'flex-end';
  margin-top: 32px;
`

const AlternativeAction = styled.div`
  font-weight: bold;

  > a {
    color: ${(props) => props.theme.colorText};
  }
`

const AuthenticationForm = (props) => {
  const {
    alternativeActionLabel,
    alternativeActionLink,
    className,
    children,

    loading,
    onSubmit,
    submitButtonLabel,
  } = props

  return (
    <Wrapper className={className}>
      <Form layout="vertical" onFinish={onSubmit}>
        {children}

        <SubmitButton htmlType="submit" loading={loading} type="primary">
          {submitButtonLabel}
        </SubmitButton>
      </Form>

      {!!alternativeActionLabel && (
        <Footer>
          <AlternativeAction>
            <Link to={alternativeActionLink}>{alternativeActionLabel}</Link>
          </AlternativeAction>
        </Footer>
      )}
    </Wrapper>
  )
}

AuthenticationForm.propTypes = {
  /** Function to run on form submit */
  onSubmit: PropTypes.func.isRequired,
  className: PropTypes.string,
  children: PropTypes.arrayOf(PropTypes.element),
  /** Text displayed at bottom right */
  alternativeActionLabel: PropTypes.string,
  /** Link to redirect to when clicking on alternative action */
  alternativeActionLink: PropTypes.string,
  /** Control waiting for response status */
  loading: PropTypes.bool,
  /** Text displayed inside submit button */
  submitButtonLabel: PropTypes.string,
}

// AuthenticationForm.defaultProps = {
//   alternativeActionLabel: null,
//   alternativeActionLink: null,
//   loading: false,
//   submitButtonLabel: "Submit",
// };

export default AuthenticationForm
