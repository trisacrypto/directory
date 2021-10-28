import React from 'react';
import { i18n } from "@lingui/core";
import { t } from "@lingui/macro";


const AddressTypeCode = () => {

  return (
    <>
    <option value={1}>{i18n._(t`Residential`)}</option>
    <option value={2}>{i18n._(t`Business`)}</option>
    <option value={3}>{i18n._(t`Geographic`)}</option>
    <option value={0}>{i18n._(t`Unspecified`)}</option>
    </>
  )
}

export default AddressTypeCode;