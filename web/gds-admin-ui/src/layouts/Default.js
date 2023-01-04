import PropTypes from 'prop-types';
import React, { Suspense, useEffect } from 'react';

const loading = () => <div className="" />;

const DefaultLayout = (props) => {
  useEffect(() => {
    if (document.body) document.body.classList.add('authentication-bg');

    return () => {
      if (document.body) document.body.classList.remove('authentication-bg');
    };
  }, []);

  // get the child view which we would like to render
  const children = props.children || null;

  return <Suspense fallback={loading()}>{children}</Suspense>;
};

DefaultLayout.propTypes = {
  children: PropTypes.node.isRequired,
};

export default DefaultLayout;
