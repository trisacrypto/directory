import React from 'react';
import { useLingui } from "@lingui/react";


const LegalPersonNameTypeCode = () => {
  const { i18n } = useLingui();

  return (
    <>
    <option value={0}>{i18n._("Unspecified")}</option>
    <option value={1}>{i18n._("Legal Name")}</option>
    <option value={2}>{i18n._("Short Name")}</option>
    <option value={3}>{i18n._("Trading Name")}</option>
    </>
  )
}

export default LegalPersonNameTypeCode;