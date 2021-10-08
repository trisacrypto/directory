import React from 'react';
import { Trans } from "@lingui/macro"


const NaturalPersonNameTypeCode = () => {
  return (
    <>
    <option value={1}><Trans>Alias Name</Trans></option>
    <option value={2}><Trans>Name at Birth</Trans></option>
    <option value={3}><Trans>Maiden Name</Trans></option>
    <option value={4}><Trans>Legal Name</Trans></option>
    <option value={0}><Trans>Unspecified</Trans></option>
    </>
  )
}

export default NaturalPersonNameTypeCode;