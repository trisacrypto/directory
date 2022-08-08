import React from 'react';
import { Stack, Box, Text, Heading, HStack, Flex, chakra } from '@chakra-ui/react';
import { IoEllipse } from 'react-icons/io5';
interface StatusCardProps {
  isOnline: string;
}

const StatusCard = ({ isOnline }: StatusCardProps) => {
  const status = !!(typeof isOnline === 'string' && isOnline.toUpperCase() === 'HEALTHY');

  return (
    <Stack
      fontSize={40}
      pt={5}
      alignItems={'center'}
      textAlign={'center'}
      justifyContent={'center'}
      mx={'auto'}>
      {
        <IoEllipse
          fontSize="3rem"
          fill={status ? '#60C4CA' : '#C4C4C4'}
          data-testid="status__color"
        />
      }
    </Stack>
  );
};
StatusCard.defaultProps = {
  isOnline: 'HEALTH'
};
export default StatusCard;
