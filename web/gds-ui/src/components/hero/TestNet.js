import React from 'react';
import { Trans } from "@lingui/macro";


const TestNet = () => {
  return (
    <header className="bg-testnet-gradient">
      <div className="container">
        <div className="text-center hero">
          <h1><Trans>TRISA TestNet Directory</Trans></h1>
          <p className="lead">
            <Trans>Get started with the TRISA TestNet to implement your Travel Rule compliance service.</Trans>
          </p>
          <small>
            <Trans>Looking for the <a href="https://trisa.directory">Production Directory Service</a>?</Trans>
          </small>
        </div>
      </div>
    </header>
  );
};

export default TestNet;