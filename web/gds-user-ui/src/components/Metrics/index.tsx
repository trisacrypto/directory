import { Text, SimpleGrid } from '@chakra-ui/react';
import StatCard from 'components/StatCard';
import StatusCard from 'components/StatusCard';
import * as Sentry from '@sentry/react';
interface MetricsProps {
  data: any;
  type: string;
}
const Metrics = ({ data, type }: MetricsProps) => {
  console.log('[Metrics] data', data.status);
  return (
    <Sentry.ErrorBoundary
      fallback={
        <Text color={'red'} pt={20}>{`An error has occurred to load ${type} metric`}</Text>
      }>
      <SimpleGrid columns={{ base: 4, sm: 2, lg: 4, md: 4 }} spacingX="20px" spacingY="20px">
        <StatCard title="Network Status">
          <StatusCard isOnline={data?.status || 'UNKNOWN'} />
        </StatCard>
        <StatCard title="Verified VASPs">{data?.vasps_count}</StatCard>
        <StatCard title="Identity Certificates">{data?.certificates_issued}</StatCard>
        <StatCard title="New Members">{data?.new_members}</StatCard>
      </SimpleGrid>
    </Sentry.ErrorBoundary>
  );
};

export default Metrics;
