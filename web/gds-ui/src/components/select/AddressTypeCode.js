import React from 'react';

const AddressTypeCode = () => {
  return (
    <>
    <option value={1}>Residential</option>
    <option value={2}>Business</option>
    <option value={3}>Geographic</option>
    <option value={0}>Unspecified</option>
    </>
  )
}

export default AddressTypeCode;