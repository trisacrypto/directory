import React from 'react';
import { useLingui } from "@lingui/react";
import { t } from "@lingui/macro";


const VASPCategory = () => {
  const { i18n } = useLingui();

  return (
    <>
    <option value="Exchange">{i18n._(t`Centralized Exchange`)}</option>
    <option value="DEX">{i18n._(t`Decentralized Exchange`)}</option>
    <option value="P2P">{i18n._(t`Person-to-Person Exchange`)}</option>
    <option value="Kiosk">{i18n._(t`Kiosk / Crypto ATM Operator`)}</option>
    <option value="Custodian">{i18n._(t`Custody Provider`)}</option>
    <option value="OTC">{i18n._(t`Over-The-Counter Trading Desk`)}</option>
    <option value="Fund">{i18n._(t`Investment Fund - hedge funds, ETFs, and family offices`)}</option>
    <option value="Project">{i18n._(t`Token Project`)}</option>
    <option value="Gambling">{i18n._(t`Gambling or Gaming Site`)}</option>
    <option value="Miner">{i18n._(t`Mining Pool`)}</option>
    <option value="Mixer">{i18n._(t`Mixing Service`)}</option>
    <option value="Individual">{i18n._(t`Legal person`)}</option>
    <option value="Other">{i18n._(t`Other`)}</option>
    </>
  )
}

export default VASPCategory;

