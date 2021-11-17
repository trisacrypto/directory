import React from 'react';
import { t } from "@lingui/macro";


const NaturalPersonNameTypeCode = () => {
  
  return (
    <>
    <option value={1}>{t`Alias Name`}</option>
    <option value={2}>{t`Name at Birth`}</option>
    <option value={3}>{t`Maiden Name`}</option>
    <option value={4}>{t`Legal Name`}</option>
    <option value={0}>{t`Unspecified`}</option>
    </>
  )
}

export default NaturalPersonNameTypeCode;