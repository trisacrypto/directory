import React from 'react';
import { useLingui } from "@lingui/react";


const AddressTypeCode = () => {
  const { i18n } = useLingui();

  return (
    <>
    <option value={1}>{i18n._("Residential")}</option>
    <option value={2}>{i18n._("Business")}</option>
    <option value={3}>{i18n._("Geographic")}</option>
    <option value={0}>{i18n._("Unspecified")}</option>
    </>
  )
}

export default AddressTypeCode;