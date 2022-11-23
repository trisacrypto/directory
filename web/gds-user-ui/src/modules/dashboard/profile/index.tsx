import React, { useState } from 'react';

import Profile from 'components/UserProfile';

const UserProfile: React.FC = () => {
  useState(() => {
    // fetch user information
  });

  // const handleUpdate = () => {};
  return (
    <>
      <Profile />
    </>
  );
};

export default UserProfile;
