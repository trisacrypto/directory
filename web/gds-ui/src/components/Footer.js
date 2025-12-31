import React from 'react';
import { Trans } from "@lingui/macro";


const Footer = () => {
  return (
  <footer className="footer bg-dark">
    <div className="container text-center text-white">
      <p className="mt-3 mb-0 py-0">
        <Trans>A component of the <a href="https://travelrule.io/">TRISA</a> architecture for Cryptocurrency Travel Rule compliance.</Trans>
      </p>
      <p className="my-0 py-0 text-muted">
        <small><Trans>Created and maintained by <a href="https://rotational.io">Rotational Labs</a> in partnership with <a href="https://ciphertrace.com">CipherTrace</a> on behalf of <a href="https://travelrule.io">TRISA</a>.</Trans></small>
      </p>
    </div>
  </footer>
  );
}

export default Footer;
