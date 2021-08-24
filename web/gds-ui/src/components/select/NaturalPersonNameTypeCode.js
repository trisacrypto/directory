import React from 'react';

const NaturalPersonNameTypeCode = () => {
  return (
    <>
    <option value={1}>Alias Name</option>
    <option value={2}>Name at Birth</option>
    <option value={3}>Maiden Name</option>
    <option value={4}>Legal Name</option>
    <option value={0}>Unspecified</option>
    </>
  )
}

export default NaturalPersonNameTypeCode;