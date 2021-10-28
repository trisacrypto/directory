import React from 'react';
import { useLingui } from "@lingui/react";


const LegalPersonNameTypeCode = () => {
  const { i18n } = useLingui();

  return (
    <>
    <option value={0}>{i18n._(t`Unspecified`)}</option>
    <option value={1}>{i18n._(t`Legal Name`)}</option>
    <option value={2}>{i18n._(t`Short Name`)}</option>
    <option value={3}>{i18n._(t`Trading Name`)}</option>
    </>
  )
}

export default LegalPersonNameTypeCode;