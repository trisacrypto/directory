import React from 'react';
import { Trans } from "@lingui/macro"


const NationalIdentifierTypeCode = () => {
  return (
    <>
    <option value={1}><Trans>Alien Registration Number</Trans></option>
    <option value={2}><Trans>Passport Number</Trans></option>
    <option value={3}><Trans>Registration Authority Identifier</Trans></option>
    <option value={4}><Trans>Driver's License Number</Trans></option>
    <option value={5}><Trans>Foreign Investment Identity Number</Trans></option>
    <option value={6}><Trans>Tax Identification Number</Trans></option>
    <option value={7}><Trans>Social Security Number</Trans></option>
    <option value={8}><Trans>Identity Card Number</Trans></option>
    <option value={9}><Trans>Legal Entity Identifier (LEI)</Trans></option>
    <option value={0}><Trans>Unspecified</Trans></option>
    </>
  )
}

export default NationalIdentifierTypeCode;