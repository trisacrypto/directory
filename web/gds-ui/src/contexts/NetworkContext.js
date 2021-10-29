import React, { useState, useEffect } from 'react';
import { isTestNet as getTestNetEnv } from '../lib/testnet';

const Context = React.createContext();

export const NetworkStore = ({children}) => {
  const [ isTestNet ] = useState(getTestNetEnv());

  // Report the network the application is running on.
  useEffect(() => {
    if (isTestNet) {
      console.log("TestNet GDS UI context loaded");
    } else {
      console.log("MainNet (Production) GDS UI context loaded");
    }
  }, [isTestNet])

  return (
    <Context.Provider value={isTestNet}>
      {children}
    </Context.Provider>
  );
}

export default Context;