import React from 'react';
import { useLingui } from "@lingui/react";


const NationalIdentifierTypeCode = () => {
  const { i18n } = useLingui();

  return (
    <>
    <option value={1}>{i18n._("Alien Registration Number")}</option>
    <option value={2}>{i18n._("Passport Number")}</option>
    <option value={3}>{i18n._("Registration Authority Identifier")}</option>
    <option value={4}>{i18n._("Driver's License Number")}</option>
    <option value={5}>{i18n._("Foreign Investment Identity Number")}</option>
    <option value={6}>{i18n._("Tax Identification Number")}</option>
    <option value={7}>{i18n._("Social Security Number")}</option>
    <option value={8}>{i18n._("Identity Card Number")}</option>
    <option value={9}>{i18n._("Legal Entity Identifier (LEI)")}</option>
    <option value={0}>{i18n._("Unspecified")}</option>
    </>
  )
}

export default NationalIdentifierTypeCode;