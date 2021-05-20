import React from 'react';

const NationalIdentifierTypeCode = () => {
  return (
    <>
    <option value={0}>Alien Registration Number</option>
    <option value={1}>Passport Number</option>
    <option value={2}>Registration Authority Identifier</option>
    <option value={3}>Driver's License Number</option>
    <option value={4}>Foreign Investment Identity Number</option>
    <option value={5}>Tax Identification Number</option>
    <option value={6}>Social Security Number</option>
    <option value={7}>Identity Card Number</option>
    <option value={8}>Legal Entity Identifier (LEI)</option>
    <option value={9}>Unspecified</option>
    </>
  )
}

export default NationalIdentifierTypeCode;