import React from 'react';
import { i18n } from "@lingui/core";
import { t } from "@lingui/macro";


const TestNet = () => {
  return (
    <header className="bg-testnet-gradient">
      <div className="container">
        <div className="text-center hero">
          <h1>{i18n._(t`TRISA TestNet Directory`)}</h1>
          <p className="lead">
            <Trans>Get started with the TRISA TestNet to implement your Travel Rule compliance service.</Trans>
          </p>
          <small>
            <Trans>Looking for the <a href="https://vaspdirectory.net">Production Directory Service</a>?</Trans>
          </small>
        </div>
      </div>
    </header>
  );
};

export default TestNet;