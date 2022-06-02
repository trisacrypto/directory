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
import { useLocation, useNavigate } from 'react-router-dom';
const Overview: React.FC = () => {
  const [result, setResult] = React.useState<any>('');
  const { user, getUser } = useAuth();
  const navigate = useNavigate();

  useEffect(() => {
    (async () => {
      try {
        const response = await getMetrics();
        console.log('response', response);
        setResult(response.data);
      } catch (e: any) {
        if (e.response.status === 401) {
          navigate('/auth/login?redirect=/dashboard/overview&q=unauthorized');
        }
        if (e.response.status === 403) {
          navigate('/auth/login?redirect=/dashboard/overview&q=token_expired');
        }

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
