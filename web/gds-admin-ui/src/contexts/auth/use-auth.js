import React from 'react';

import AuthContext from './auth-context';

const useAuth = () => {
  const context = React.useContext(AuthContext);

  if (!context) {
    throw new Error('useAuth should be used within an AuthProvider');
  }
  return context;
};

export default useAuth;
