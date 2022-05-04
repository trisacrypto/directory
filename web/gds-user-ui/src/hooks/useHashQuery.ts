import React from 'react';
import { BrowserRouter as Router, Link, useLocation } from 'react-router-dom';
import queryString from 'query-string';

function useHashQuery() {
  const location = useLocation();
  const hash = location.hash;

  return React.useMemo(() => queryString.parse(hash), [hash]);
}

export default useHashQuery;
