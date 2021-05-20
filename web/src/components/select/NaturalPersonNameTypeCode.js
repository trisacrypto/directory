import React from 'react';

const NaturalPersonNameTypeCode = () => {
  return (
    <>
    <option value={0}>Alias Name</option>
    <option value={1}>Name at Birth</option>
    <option value={2}>Maiden Name</option>
    <option value={3}>Legal Name</option>
    <option value={4}>Unspecified</option>
    </>
  )
}

export default NaturalPersonNameTypeCode;