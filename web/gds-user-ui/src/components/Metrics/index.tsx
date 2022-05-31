import React from 'react';
import { Stack, Box, Text, Heading, HStack, Flex } from '@chakra-ui/react';
import StatCard from 'components/StatCard';
import StatusCard from 'components/StatusCard';
interface MetricsProps {
  datas?: any;
  status: boolean;
}
const Metrics = ({ datas, status }: MetricsProps) => {
  return (
    <Flex pt={4}>
      <HStack spacing="24">
        <StatusCard isOnline={status} />

        <StatCard title="Verified VASPs" number={248} />
        <StatCard title="TestNet Cerificates" number={248} />
        <StatCard title="MainNet Cerificates" number={248} />
      </HStack>
    </Flex>
  );
};
Metrics.defaultProps = {
  title: 'Verified VASPs',
  number: 248
};

export default Metrics;
