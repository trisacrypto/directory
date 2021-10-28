import React from 'react';
import { useLingui } from "@lingui/react";


const NationalIdentifierTypeCode = () => {
  const { i18n } = useLingui();

  return (
    <>
    <option value={1}>{i18n._(t`Alien Registration Number`)}</option>
    <option value={2}>{i18n._(t`Passport Number`)}</option>
    <option value={3}>{i18n._(t`Registration Authority Identifier`)}</option>
    <option value={4}>{i18n._(t`Driver's License Number`)}</option>
    <option value={5}>{i18n._(t`Foreign Investment Identity Number`)}</option>
    <option value={6}>{i18n._(t`Tax Identification Number`)}</option>
    <option value={7}>{i18n._(t`Social Security Number`)}</option>
    <option value={8}>{i18n._(t`Identity Card Number`)}</option>
    <option value={9}>{i18n._(t`Legal Entity Identifier (LEI)`)}</option>
    <option value={0}>{i18n._(t`Unspecified`)}</option>
    </>
  )
}

export default NationalIdentifierTypeCode;