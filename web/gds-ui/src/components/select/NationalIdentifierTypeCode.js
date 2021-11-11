import React from 'react';
import { t } from "@lingui/macro";


const NationalIdentifierTypeCode = () => {

  return (
    <>
    <option value={1}>{t`Alien Registration Number`}</option>
    <option value={2}>{t`Passport Number`}</option>
    <option value={3}>{t`Registration Authority Identifier`}</option>
    <option value={4}>{t`Driver's License Number`}</option>
    <option value={5}>{t`Foreign Investment Identity Number`}</option>
    <option value={6}>{t`Tax Identification Number`}</option>
    <option value={7}>{t`Social Security Number`}</option>
    <option value={8}>{t`Identity Card Number`}</option>
    <option value={9}>{t`Legal Entity Identifier (LEI)`}</option>
    <option value={0}>{t`Unspecified`}</option>
    </>
  )
}

export default NationalIdentifierTypeCode;