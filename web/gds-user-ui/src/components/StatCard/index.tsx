import React from 'react';
import { Stack, Box, Text, Heading, Flex, chakra } from '@chakra-ui/react';

interface StatCardProps {
  title: string;
  number: number;
}
const StatCard = ({ title, number }: StatCardProps) => {
  return (
    <Box
      bg={'white'}
      border="1px solid #DFE0EB"
      fontFamily={'Open Sans'}
      color={'#252733'}
      textAlign={'center'}
      // minWidth={250}
      // height={170}
      fontSize={18}
      p={5}
      mt={10}
      px={5}>
      <Stack>
        <chakra.h1 textAlign={'center'} fontSize={20} fontWeight={'bold'}>
          {title}
        </chakra.h1>
        <Text fontSize={40} pt={3} fontWeight={'bold'}>
          {number}
        </Text>
      </Stack>
    </Box>
  );
};
StatCard.defaultProps = {
  title: 'Verified VASPs',
  number: 0
};

export default StatCard;
