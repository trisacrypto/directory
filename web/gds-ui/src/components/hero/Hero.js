import React from 'react';
import MainNet from './MainNet';
import TestNet from './TestNet';
import NetworkContext from '../../contexts/NetworkContext';

const Hero = () => {
  const chooseHero = (isTestNet) => {
    if (isTestNet) {
      return <TestNet />
    }
    return <MainNet />
  }

  return (
    <NetworkContext.Consumer>
      {chooseHero}
    </NetworkContext.Consumer>
  );
};

export default Hero;