import React, { useState } from 'react';

import { Heading, Stack } from '@chakra-ui/react';
import Card from 'components/ui/Card';
import DashboardLayout from 'layouts/DashboardLayout';
import InputFormControl from 'components/ui/InputFormControl';
import UserDetails from 'components/UserDetails';

const UserProfile: React.FC = () => {
  useState(() => {
    // fetch user information
  });

  const handleUpdate = () => {};
  return (
    <DashboardLayout>
      <Heading marginBottom="69px">User Account</Heading>
      <Card>
        <Card.Header>
          {' '}
          <Heading as="h4" size="md" marginBottom="20px">
            User Settings
          </Heading>
        </Card.Header>
        <Card.Body>
          <Stack spacing={4}>
            <InputFormControl
              controlId="name"
              label="Name"
              value={''}
              hasBtn
              handleFn={handleUpdate}
            />
            <InputFormControl
              controlId="email"
              label="Email Address"
              type="email"
              value={''}
              hasBtn
              handleFn={handleUpdate}
            />
            <InputFormControl
              controlId="password"
              label="Password"
              type="password"
              value={''}
              hasBtn
              setBtnName="Change"
              handleFn={handleUpdate}
            />
          </Stack>
          <UserDetails
            userId="C0000213"
            createdDate="01/01/2020"
            status="Active"
            permissions="Admin"
            lastLogin="01/01/2020"
          />
        </Card.Body>
      </Card>
    </DashboardLayout>
  );
};

export default UserProfile;
