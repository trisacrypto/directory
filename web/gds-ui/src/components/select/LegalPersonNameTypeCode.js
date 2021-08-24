import React from 'react';

const LegalPersonNameTypeCode = () => {
  return (
    <>
    <option value={0}>Unspecified</option>
    <option value={1}>Legal Name</option>
    <option value={2}>Short Name</option>
    <option value={3}>Trading Name</option>
    </>
  )
}

export default LegalPersonNameTypeCode;