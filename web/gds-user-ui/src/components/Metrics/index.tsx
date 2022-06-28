import React from 'react';
import { Stack, Box, Text, Heading, HStack, Flex, SimpleGrid } from '@chakra-ui/react';
import StatCard from 'components/StatCard';
import StatusCard from 'components/StatusCard';
import * as Sentry from '@sentry/react';
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
            <Text color={'red'} pt={20}>{`An error has occurred to load ${type} metric`}</Text>
          }>
          <SimpleGrid columns={{ base: 4, sm: 2, lg: 4, md: 4 }} spacingX="20px" spacingY="20px">
            <StatusCard isOnline={data?.status || 'UNKNOW'} />
            <StatCard title="Verified VASPs" number={data?.vasps_count} />
            <StatCard title="Identity Certificates" number={data?.certificates_issued} />
            <StatCard title="New Members" number={data?.new_members} />
          </SimpleGrid>
        </Sentry.ErrorBoundary>
      </Box>
    </Flex>
  );
};

export default Metrics;
