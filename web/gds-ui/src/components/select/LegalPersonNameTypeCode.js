import React from 'react';
import { i18n } from "@lingui/core";
import { t } from "@lingui/macro";


const LegalPersonNameTypeCode = () => {

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