import React from 'react';
import { Trans } from "@lingui/macro"


const AddressTypeCode = () => {
  return (
    <>
    <option value={1}><Trans>Residential</Trans></option>
    <option value={2}><Trans>Business</Trans></option>
    <option value={3}><Trans>Geographic</Trans></option>
    <option value={0}><Trans>Unspecified</Trans></option>
    </>
  )
}

export default AddressTypeCode;