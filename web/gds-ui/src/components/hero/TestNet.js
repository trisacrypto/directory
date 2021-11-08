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
            {i18n._(t`Get started with the TRISA TestNet to implement your Travel Rule compliance service.`)}
          </p>
          <small>
            {i18n._(t`Looking for the <a href="https://vaspdirectory.net">Production Directory Service</a>?`)}
          </small>
        </div>
      </div>
    </header>
  );
};

export default TestNet;