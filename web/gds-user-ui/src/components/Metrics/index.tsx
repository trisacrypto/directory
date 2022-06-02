import React from 'react';
import { Stack, Box, Text, Heading, HStack, Flex } from '@chakra-ui/react';
import StatCard from 'components/StatCard';
import StatusCard from 'components/StatusCard';
import * as Sentry from '@sentry/react';
interface MetricsProps {
  data: any;
  type: string;
}
const Metrics = ({ data, type }: MetricsProps) => {
  return (
    <Flex my={5}>
      {/* <Stack bg={'#E5EDF1'} h="55px" justifyItems={'center'} p={4} my={5}>
        <Stack>
          <Heading fontSize={20}>{`${type} Network Metrics`}</Heading>
        </Stack>
      </Stack> */}
      <Box textAlign={'center'} justifyContent="center" justifyItems={'center'}>
        <Sentry.ErrorBoundary
          fallback={<Text color={'red'}>{`An error has occurred to load ${type} metric`}</Text>}>
          <HStack spacing="24" textAlign={'center'}>
            <StatusCard isOnline={data?.status || false} />
            <StatCard title="Verified VASPs" number={data?.vasps_count} />
            <StatCard title="Certificates" number={data?.certificates_issued} />
            <StatCard title="New Members" number={data?.new_members} />
          </HStack>
        </Sentry.ErrorBoundary>
      </Box>
    </Flex>
  );
};

export default Metrics;
