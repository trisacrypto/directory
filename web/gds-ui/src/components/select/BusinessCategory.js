import React from 'react';
import { t } from "@lingui/macro";


const BusinessCategory = () => {

  return (
    <>
    <option value={0}></option>
    <option value={1}>{t`Private Organization`}</option>
    <option value={2}>{t`Government Entity`}</option>
    <option value={3}>{t`Business Entity`}</option>
    <option value={4}>{t`Non-Commercial Entity`}</option>
    </>
  )
}

export default BusinessCategory;