import React from 'react';
import { t } from "@lingui/macro";


const VASPCategory = () => {

  return (
    <>
    <option value="Exchange">{t`Centralized Exchange`}</option>
    <option value="DEX">{t`Decentralized Exchange`}</option>
    <option value="P2P">{t`Person-to-Person Exchange`}</option>
    <option value="Kiosk">{t`Kiosk / Crypto ATM Operator`}</option>
    <option value="Custodian">{t`Custody Provider`}</option>
    <option value="OTC">{t`Over-The-Counter Trading Desk`}</option>
    <option value="Fund">{t`Investment Fund - hedge funds, ETFs, and family offices`}</option>
    <option value="Project">{t`Token Project`}</option>
    <option value="Gambling">{t`Gambling or Gaming Site`}</option>
    <option value="Miner">{t`Mining Pool`}</option>
    <option value="Mixer">{t`Mixing Service`}</option>
    <option value="Individual">{t`Legal person`}</option>
    <option value="Other">{t`Other`}</option>
    </>
  )
}

export default VASPCategory;

