import React from 'react';
import { useLingui } from "@lingui/react";


const NaturalPersonNameTypeCode = () => {
  const { i18n } = useLingui();
  return (
    <>
    <option value={1}>{i18n._("Alias Name")}</option>
    <option value={2}>{i18n._("Name at Birth")}</option>
    <option value={3}>{i18n._("Maiden Name")}</option>
    <option value={4}>{i18n._("Legal Name")}</option>
    <option value={0}>{i18n._("Unspecified")}</option>
    </>
  )
}

export default NaturalPersonNameTypeCode;