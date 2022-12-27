import React, { useState } from 'react';
import { Stack } from '@chakra-ui/react';
import Profile from 'components/UserProfile';

const UserProfile: React.FC = () => {
  useState(() => {
    // fetch user information
  });

  // const handleUpdate = () => {};
  return (
    <>
      <Stack my={5}>
        <Profile />
      </Stack>
    </>
  );
};

export default UserProfile;
