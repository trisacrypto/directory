import React, { ReactNode } from 'react';
import { Stack, Box, Text, Heading, Flex, chakra } from '@chakra-ui/react';

export interface StatCardProps {
  title: string;
  children: ReactNode;
}
const StatCard = ({ title = 'Verified VASPs', children = 0 }: StatCardProps) => {
  return (
    <Box
      bg={'white'}
      border="1px solid #DFE0EB"
      fontFamily={'Open Sans'}
      color={'#252733'}
      textAlign={'center'}
      p={5}
      mt={10}
      px={5}>
      <Stack>
        <chakra.h1
          textAlign={'center'}
          fontSize={20}
          fontWeight={'bold'}
          textTransform="capitalize"
          data-testid="start-card__title">
          {title}
        </chakra.h1>
        <Text fontSize={40} pt={3} fontWeight={'bold'} data-testid="start-card__body">
          {children}
        </Text>
      </Stack>
    </Box>
  );
};

export default StatCard;
