import React from 'react';
import { t } from "@lingui/macro";


const AddressTypeCode = () => {

  return (
    <>
    <option value={1}>{t`Residential`}</option>
    <option value={2}>{t`Business`}</option>
    <option value={3}>{t`Geographic`}</option>
    <option value={0}>{t`Unspecified`}</option>
    </>
  )
}

export default AddressTypeCode;