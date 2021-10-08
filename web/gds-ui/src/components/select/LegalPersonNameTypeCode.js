import React from 'react';
import { Trans } from "@lingui/macro"


const LegalPersonNameTypeCode = () => {
  return (
    <>
    <option value={0}><Trans>Unspecified</Trans></option>
    <option value={1}><Trans>Legal Name</Trans></option>
    <option value={2}><Trans>Short Name</Trans></option>
    <option value={3}><Trans>Trading Name</Trans></option>
    </>
  )
}

export default LegalPersonNameTypeCode;