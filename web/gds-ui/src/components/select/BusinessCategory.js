import React from 'react';
import { useLingui } from "@lingui/react";


const BusinessCategory = () => {
  const { i18n } = useLingui();

  return (
    <>
    <option value={0}></option>
    <option value={1}>{i18n._(t`Private Organization`)}</option>
    <option value={2}>{i18n._(t`Government Entity`)}</option>
    <option value={3}>{i18n._(t`Business Entity`)}</option>
    <option value={4}>{i18n._(t`Non-Commercial Entity`)}</option>
    </>
  )
}

export default BusinessCategory;