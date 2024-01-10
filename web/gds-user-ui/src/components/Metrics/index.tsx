import React from 'react';
import { Box, Text, Flex, SimpleGrid } from '@chakra-ui/react';
import StatCard from 'components/StatCard';
import StatusCard from 'components/StatusCard';
import * as Sentry from '@sentry/react';
import { t } from '@lingui/macro';

interface MetricsProps {
  data: any;
  type: string;
}
const Metrics = ({ data, type }: MetricsProps) => {
  return (
    <Flex>
      <Box textAlign={'center'} justifyContent="center" justifyItems={'center'} mx={'auto'}>
        <Sentry.ErrorBoundary
          fallback={
            <Text color={'red'} pt={20}>{t`An error has occurred to load ${type} metric`}</Text>
          }>
          <SimpleGrid columns={{ base: 1, sm: 2, lg: 4 }} spacingX="30px">
            <StatusCard isOnline={data?.status || 'UNKNOWN'} />
            <StatCard title={t`Verified VASPs`} number={data?.vasps} />
            <StatCard title={t`Identity Certificates`} number={data?.certificates_issued} />
            <StatCard title={t`New Members`} number={data?.new_members} />
          </SimpleGrid>
        </Sentry.ErrorBoundary>
      </Box>
    </Flex>
  );
};

export default Metrics;
