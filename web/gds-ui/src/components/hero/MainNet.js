import React from 'react';
import { i18n } from "@lingui/core";
import { t } from "@lingui/macro";


const MainNet = () => {
  return (
    <header className="bg-gradient">
      <div className="container">
        <div className="text-center hero">
          <h1>{i18n._(t`TRISA Global Directory Service`)}</h1>
          <p className="lead">
           {i18n._(t`Become a TRISA certified Virtual Asset Service Provider.`)}
          </p>
          <small>
            {i18n._(t`Looking for the <a href="https://trisatest.net">TestNet Directory Service</a>?`)}
          </small>
        </div>
      </div>
    </header>
  );
};

export default MainNet;