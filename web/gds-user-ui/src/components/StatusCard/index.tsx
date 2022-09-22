import React from 'react';
import { Stack, Box, chakra } from '@chakra-ui/react';
import { IoEllipse } from 'react-icons/io5';
import { Trans } from '@lingui/react';
interface StatusCardProps {
  isOnline: string;
}

const StatusCard = ({ isOnline }: StatusCardProps) => {
  const status = !!(typeof isOnline === 'string' && isOnline.toUpperCase() === 'HEALTHY');

  return (
    <Box
      bg={'white'}
      border="1px solid #DFE0EB"
      fontFamily={'Open Sans'}
      color={'#252733'}
      // minWidth={250}
      // height={170}
      fontSize={18}
      p={5}
      mt={10}
      px={5}>
      <Stack textAlign={'center'}>
        <chakra.h1 textAlign={'center'} fontSize={20} fontWeight={'bold'}>
          <Trans id="Network Status">Network Status</Trans>
        </chakra.h1>
        <Stack
          fontSize={40}
          pt={5}
          alignItems={'center'}
          textAlign={'center'}
          justifyContent={'center'}
          mx={'auto'}>
          {status ? (
            <IoEllipse fontSize="3rem" fill={'#60C4CA'} />
          ) : (
            <IoEllipse fontSize="3rem" fill={'#C4C4C4'} />
          )}
        </Stack>
      </Stack>
    </Box>
  );
};
StatusCard.defaultProps = {
  isOnline: 'HEALTH'
};
export default StatusCard;
