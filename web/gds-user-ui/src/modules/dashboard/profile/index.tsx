import React, { useState, Suspense } from 'react';
import Loader from 'components/Loader';
import Profile from 'components/UserProfile';

const UserProfile: React.FC = () => {
  useState(() => {
    // fetch user information
  });

  // const handleUpdate = () => {};
  return (
    <>
      <Suspense fallback={<Loader />}>
        <Profile />
      </Suspense>
    </>
  );
};

export default UserProfile;
