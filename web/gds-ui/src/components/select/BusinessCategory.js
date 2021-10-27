import React from 'react';
import { useLingui } from "@lingui/react";


const BusinessCategory = () => {
  const { i18n } = useLingui();

  return (
    <>
    <option value={0}></option>
    <option value={1}>{i18n._("Private Organization")}</option>
    <option value={2}>{i18n._("Government Entity")}</option>
    <option value={3}>{i18n._("Business Entity")}</option>
    <option value={4}>{i18n._("Non-Commercial Entity")}</option>
    </>
  )
}

export default BusinessCategory;