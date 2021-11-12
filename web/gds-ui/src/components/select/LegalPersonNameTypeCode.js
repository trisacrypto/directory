import React from 'react';
import { t } from "@lingui/macro";


const LegalPersonNameTypeCode = () => {

  return (
    <>
    <option value={0}>{t`Unspecified`}</option>
    <option value={1}>{t`Legal Name`}</option>
    <option value={2}>{t`Short Name`}</option>
    <option value={3}>{t`Trading Name`}</option>
    </>
  )
}

export default LegalPersonNameTypeCode;