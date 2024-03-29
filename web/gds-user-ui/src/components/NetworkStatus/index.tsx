import { Stack, Text } from '@chakra-ui/react';
import { Trans } from '@lingui/react';

import { IoEllipse } from 'react-icons/io5';

interface NetworkStatusProps {
  isOnline: boolean;
}
const NetworkStatus = (props: NetworkStatusProps) => {
  const { isOnline } = props;
  return (
    <Stack minHeight={82} bg={'white'} p={5} border="1px solid #C4C4C4">
      <Stack direction={'row'} justifyContent="space-between" alignItems="center">
        <Text fontWeight={'bold'}>
          <Trans id="Network Status">Network Status</Trans>
        </Text>
        {isOnline ? (
          <IoEllipse fontSize="2rem" fill={'#60C4CA'} height={'45px'} />
        ) : (
          <IoEllipse fontSize="2rem" fill={'#C4C4C4'} />
        )}
      </Stack>
    </Stack>
  );
};
NetworkStatus.defaultProps = {
  isOnline: true
};
export default NetworkStatus;
