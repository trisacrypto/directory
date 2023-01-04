import classNames from 'classnames';
import PropTypes from 'prop-types';
import React from 'react';
import { Button } from 'react-bootstrap';

import Oval from './Oval';

function ButtonSpinner({ variant, isLoading, label, loadingMessage, ...rest }) {
  return (
    <Button
      data-testid="deleteBtn"
      variant={variant}
      className={classNames('d-flex align-items-center justify-content-between gap-1', {
        'py-0': isLoading,
      })}
      disabled={isLoading}
      {...rest}
    >
      {!isLoading ? (
        <span>{label}</span>
      ) : (
        <>
          <span>
            <Oval width="20" />
          </span>
          <span>{loadingMessage}</span>
        </>
      )}
    </Button>
  );
}

ButtonSpinner.defaultProps = {
  variant: 'primary',
  isLoading: false,
  loadingMessage: 'Loading...',
};

ButtonSpinner.propTypes = {
  label: PropTypes.string.isRequired,
  isLoading: PropTypes.bool,
  loadingMessage: PropTypes.string,
  variant: PropTypes.oneOf(['danger', 'success', 'primary', 'info', 'warning']),
};

export default ButtonSpinner;
