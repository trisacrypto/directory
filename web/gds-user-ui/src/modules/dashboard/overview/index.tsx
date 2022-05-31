import React, { useState, useEffect } from 'react';

import { Box, Heading, VStack, Flex, Input, Stack } from '@chakra-ui/react';
import Card from 'components/ui/Card';
import DashboardLayout from 'layouts/DashboardLayout';
import NeedsAttention from 'components/NeedsAttention';
import NetworkAnnouncements from 'components/NetworkAnnouncements';
import Metrics from 'components/Metrics';
import useAuth from 'hooks/useAuth';
const Overview: React.FC = () => {
  const [userId, setUserId] = React.useState('');
  const { user } = useAuth();
  useState(() => {
    // fetch user information
  });
  return (
    <DashboardLayout>
      <Heading marginBottom="69px">Overview</Heading>
      <NeedsAttention />
      <NetworkAnnouncements />
      <Metrics status={false} />
    </DashboardLayout>
  );
};

export default Overview;
