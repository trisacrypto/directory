import React, { useState, useEffect } from 'react';
import * as Sentry from '@sentry/react';
import { Box, Heading, VStack, Flex, Input, Stack, Text } from '@chakra-ui/react';
import Card from 'components/ui/Card';
import DashboardLayout from 'layouts/DashboardLayout';
import NeedsAttention from 'components/NeedsAttention';
import NetworkAnnouncements from 'components/NetworkAnnouncements';
import Metrics from 'components/Metrics';
import useAuth from 'hooks/useAuth';
import { getMetrics } from './overview.service';
const Overview: React.FC = () => {
  const [result, setResult] = React.useState<any>('');
  const { user } = useAuth();
  useEffect(() => {
    (async () => {
      try {
        const response = await getMetrics();
        setResult(response);
      } catch (e: any) {
        console.log(e);
      }
    })();
  }, []);
  return (
    <DashboardLayout>
      <Heading marginBottom="69px">Overview</Heading>
      <NeedsAttention />
      <NetworkAnnouncements />
      {/* <Sentry.ErrorBoundary
        fallback={<Text color={'red'}>An error has occurred to load testnet metric</Text>}> */}
      <Metrics data={result?.testnet} type="Testnet" />
      <Metrics data={result?.mainnet} type="Mainnet" />
      {/* </Sentry.ErrorBoundary> */}
    </DashboardLayout>
  );
};

export default Overview;
