import React from 'react';
import { Trans } from "@lingui/macro";


const MainNet = () => {
  return (
    <header className="bg-gradient">
      <div className="container">
        <div className="text-center hero">
          <h1><Trans>TRISA Global Directory Service</Trans></h1>
          <p className="lead">
            <Trans>Become a TRISA certified Virtual Asset Service Provider.</Trans>
          </p>
          <small>
            <Trans>Looking for the <a href="https://testnet.directory">TestNet Directory Service</a>?</Trans>
          </small>
        </div>
      </div>
    </header>
  );
};

export default MainNet;