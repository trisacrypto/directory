import React from 'react';
import { useLingui } from "@lingui/react";


const VASPCategory = () => {
  const { i18n } = useLingui();

  return (
    <>
    <option value="Exchange">{i18n._("Centralized Exchange")}</option>
    <option value="DEX">{i18n._("Decentralized Exchange")}</option>
    <option value="P2P">{i18n._("Person-to-Person Exchange")}</option>
    <option value="Kiosk">{i18n._("Kiosk / Crypto ATM Operator")}</option>
    <option value="Custodian">{i18n._("Custody Provider")}</option>
    <option value="OTC">{i18n._("Over-The-Counter Trading Desk")}</option>
    <option value="Fund">{i18n._("Investment Fund - hedge funds, ETFs, and family offices")}</option>
    <option value="Project">{i18n._("Token Project")}</option>
    <option value="Gambling">{i18n._("Gambling or Gaming Site")}</option>
    <option value="Miner">{i18n._("Mining Pool")}</option>
    <option value="Mixer">{i18n._("Mixing Service")}</option>
    <option value="Individual">{i18n._("Legal person")}</option>
    <option value="Other">{i18n._("Other")}</option>
    </>
  )
}

export default VASPCategory;

