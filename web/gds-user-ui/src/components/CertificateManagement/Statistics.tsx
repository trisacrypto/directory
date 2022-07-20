import { Stack } from '@chakra-ui/react';
import StatisticCard from './StatisticCard';

const STATISTICS = { current: 0, expired: 0, revoked: 0, total: 0 };

function Statistics() {
  return (
    <Stack direction={'row'} flexWrap={'wrap'} spacing={4}>
      {Object.entries(STATISTICS).map((statistic) => (
        <StatisticCard key={statistic[0]} title={statistic[0]} total={statistic[1]} />
      ))}
    </Stack>
  );
}

export default Statistics;
