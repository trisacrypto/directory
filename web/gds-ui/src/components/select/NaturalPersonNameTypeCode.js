import React from 'react';
import { i18n } from "@lingui/core";
import { t } from "@lingui/macro";


const NaturalPersonNameTypeCode = () => {
  
  return (
    <>
    <option value={1}>{i18n._(t`Alias Name`)}</option>
    <option value={2}>{i18n._(t`Name at Birth`)}</option>
    <option value={3}>{i18n._(t`Maiden Name`)}</option>
    <option value={4}>{i18n._(t`Legal Name`)}</option>
    <option value={0}>{i18n._(t`Unspecified`)}</option>
    </>
  )
}

export default NaturalPersonNameTypeCode;